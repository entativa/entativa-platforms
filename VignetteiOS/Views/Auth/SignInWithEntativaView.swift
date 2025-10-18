import SwiftUI

/// Sign in with Entativa - Cross-platform authentication
struct SignInWithEntativaView: View {
    @Environment(\.dismiss) private var dismiss
    @State private var emailOrUsername = ""
    @State private var password = ""
    @State private var showPassword = false
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var showError = false
    @State private var showSuccess = false
    @State private var isNewAccount = false
    @FocusState private var focusedField: Field?
    
    var onSuccess: (String) -> Void
    
    enum Field: Hashable {
        case emailOrUsername
        case password
    }
    
    var body: some View {
        NavigationView {
            ZStack {
                EntativaColors.backgroundPrimary
                    .ignoresSafeArea()
                
                ScrollView {
                    VStack(spacing: 0) {
                        Spacer()
                            .frame(height: 40)
                        
                        // Entativa branding
                        VStack(spacing: 16) {
                            Text("entativa")
                                .font(.system(size: 48, weight: .bold, design: .rounded))
                                .italic()
                                .foregroundStyle(
                                    LinearGradient(
                                        colors: [
                                            EntativaColors.gradientStart,
                                            EntativaColors.gradientMiddle,
                                            EntativaColors.gradientEnd
                                        ],
                                        startPoint: .leading,
                                        endPoint: .trailing
                                    )
                                )
                            
                            Text("Sign in with your Entativa account")
                                .entativaBodyMedium()
                                .foregroundColor(EntativaColors.textSecondary)
                                .multilineTextAlignment(.center)
                        }
                        
                        Spacer()
                            .frame(height: 40)
                        
                        // Login form
                        VStack(spacing: 12) {
                            EntativaTextField(
                                text: $emailOrUsername,
                                placeholder: "Email or username",
                                keyboardType: .emailAddress,
                                textContentType: .username,
                                autocapitalization: .never
                            )
                            .focused($focusedField, equals: .emailOrUsername)
                            .onSubmit {
                                focusedField = .password
                            }
                            
                            EntativaSecureField(
                                text: $password,
                                placeholder: "Password",
                                showPassword: $showPassword
                            )
                            .focused($focusedField, equals: .password)
                            .onSubmit {
                                Task {
                                    await signInWithEntativa()
                                }
                            }
                            
                            Button {
                                Task {
                                    await signInWithEntativa()
                                }
                            } label: {
                                Text("Continue")
                                    .entativaButtonLarge()
                                    .foregroundColor(EntativaColors.textOnPrimary)
                                    .frame(maxWidth: .infinity)
                                    .frame(height: 48)
                                    .background(
                                        isLoading ?
                                        EntativaColors.buttonPrimaryDisabled :
                                        EntativaColors.buttonPrimary
                                    )
                                    .cornerRadius(8)
                            }
                            .disabled(isLoading)
                            .padding(.top, 8)
                        }
                        .padding(.horizontal, 32)
                        
                        Spacer()
                            .frame(height: 32)
                        
                        // Info
                        VStack(spacing: 12) {
                            Image(systemName: "info.circle.fill")
                                .font(.system(size: 24))
                                .foregroundColor(EntativaColors.info)
                            
                            Text("Don't have an Entativa account yet?")
                                .entativaLabelMedium()
                                .foregroundColor(EntativaColors.textSecondary)
                            
                            Text("We'll automatically create a Vignette account for you using your Entativa profile information.")
                                .entativaCaptionMedium()
                                .foregroundColor(EntativaColors.textTertiary)
                                .multilineTextAlignment(.center)
                                .padding(.horizontal, 40)
                        }
                        
                        Spacer()
                            .frame(height: 40)
                    }
                }
                
                // Loading overlay
                if isLoading {
                    loadingOverlay
                }
            }
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button {
                        dismiss()
                    } label: {
                        Image(systemName: "xmark")
                            .font(.system(size: 16, weight: .semibold))
                            .foregroundColor(EntativaColors.textPrimary)
                    }
                }
            }
            .alert(isNewAccount ? "Welcome!" : "Success", isPresented: $showSuccess) {
                Button("Continue") {
                    dismiss()
                }
            } message: {
                if isNewAccount {
                    Text("Your Vignette account has been created using your Entativa profile!")
                } else {
                    Text("Successfully signed in with Entativa!")
                }
            }
            .alert("Error", isPresented: $showError) {
                Button("OK", role: .cancel) {}
            } message: {
                if let errorMessage = errorMessage {
                    Text(errorMessage)
                }
            }
        }
    }
    
    private var loadingOverlay: some View {
        ZStack {
            Color.black.opacity(0.3)
                .ignoresSafeArea()
            
            VStack(spacing: 16) {
                ProgressView()
                    .scaleEffect(1.5)
                    .tint(EntativaColors.buttonPrimary)
                
                Text("Signing in...")
                    .entativaBodyMedium()
                    .foregroundColor(.white)
            }
            .padding(32)
            .background(
                RoundedRectangle(cornerRadius: 16)
                    .fill(Color.white)
                    .shadow(radius: 10)
            )
        }
    }
    
    private func signInWithEntativa() async {
        focusedField = nil
        
        guard !emailOrUsername.isEmpty else {
            errorMessage = "Please enter your email or username"
            showError = true
            return
        }
        
        guard !password.isEmpty else {
            errorMessage = "Please enter your password"
            showError = true
            return
        }
        
        isLoading = true
        errorMessage = nil
        
        do {
            // First, authenticate with Entativa
            let entativaClient = AuthAPIClient.shared
            let entativaAuth = try await entativaClient.login(
                emailOrUsername: emailOrUsername,
                password: password
            )
            
            guard let entativaToken = entativaAuth.data?.accessToken else {
                throw CrossPlatformAuthErrorVignette.serverError("Failed to get Entativa token")
            }
            
            // Then use Entativa token to sign in to Vignette
            let crossPlatformClient = CrossPlatformAuthClientVignette.shared
            let vignetteAuth = try await crossPlatformClient.signInWithEntativa(
                entativaToken: entativaToken
            )
            
            isNewAccount = vignetteAuth.data?.isNewAccount ?? false
            
            // Store the Vignette token
            if let token = vignetteAuth.data?.accessToken {
                try VignetteKeychainManager.shared.save(token: token)
                onSuccess(token)
            }
            
            showSuccess = true
        } catch {
            errorMessage = error.localizedDescription
            showError = true
        }
        
        isLoading = false
    }
}

#Preview {
    SignInWithEntativaView { _ in }
}
