import SwiftUI

/// Sign in with Vignette - Cross-platform authentication
struct SignInWithVignetteView: View {
    @Environment(\.dismiss) private var dismiss
    @State private var usernameOrEmail = ""
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
        case usernameOrEmail
        case password
    }
    
    var body: some View {
        NavigationView {
            ZStack {
                VignetteColors.backgroundPrimary
                    .ignoresSafeArea()
                
                ScrollView {
                    VStack(spacing: 0) {
                        Spacer()
                            .frame(height: 40)
                        
                        // Vignette branding
                        VStack(spacing: 16) {
                            Text("Vignette")
                                .font(.custom("Snell Roundhand", size: 48))
                                .fontWeight(.medium)
                                .foregroundColor(VignetteColors.textPrimary)
                            
                            Text("Sign in with your Vignette account")
                                .vignetteBodyMedium()
                                .foregroundColor(VignetteColors.textSecondary)
                                .multilineTextAlignment(.center)
                        }
                        
                        Spacer()
                            .frame(height: 40)
                        
                        // Login form
                        VStack(spacing: 12) {
                            VignetteTextField(
                                text: $usernameOrEmail,
                                placeholder: "Username or email",
                                keyboardType: .emailAddress,
                                textContentType: .username,
                                autocapitalization: .never
                            )
                            .focused($focusedField, equals: .usernameOrEmail)
                            .onSubmit {
                                focusedField = .password
                            }
                            
                            VignetteSecureField(
                                text: $password,
                                placeholder: "Password",
                                showPassword: $showPassword
                            )
                            .focused($focusedField, equals: .password)
                            .onSubmit {
                                Task {
                                    await signInWithVignette()
                                }
                            }
                            
                            Button {
                                Task {
                                    await signInWithVignette()
                                }
                            } label: {
                                Text("Continue")
                                    .vignetteButtonLarge()
                                    .foregroundColor(VignetteColors.textOnPrimary)
                                    .frame(maxWidth: .infinity)
                                    .frame(height: 44)
                                    .background(
                                        isLoading ?
                                        VignetteColors.buttonPrimaryDisabled :
                                        VignetteColors.buttonPrimary
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
                                .foregroundColor(VignetteColors.info)
                            
                            Text("Don't have a Vignette account yet?")
                                .vignetteLabelMedium()
                                .foregroundColor(VignetteColors.textSecondary)
                            
                            Text("We'll automatically create an Entativa account for you using your Vignette profile information.")
                                .vignetteCaptionMedium()
                                .foregroundColor(VignetteColors.textTertiary)
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
                            .foregroundColor(VignetteColors.textPrimary)
                    }
                }
            }
            .alert(isNewAccount ? "Welcome!" : "Success", isPresented: $showSuccess) {
                Button("Continue") {
                    dismiss()
                }
            } message: {
                if isNewAccount {
                    Text("Your Entativa account has been created using your Vignette profile!")
                } else {
                    Text("Successfully signed in with Vignette!")
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
                    .tint(VignetteColors.buttonPrimary)
                
                Text("Signing in...")
                    .vignetteBodyMedium()
                    .foregroundColor(.white)
            }
            .padding(32)
            .background(
                RoundedRectangle(cornerRadius: 12)
                    .fill(Color.white)
                    .shadow(radius: 10)
            )
        }
    }
    
    private func signInWithVignette() async {
        focusedField = nil
        
        guard !usernameOrEmail.isEmpty else {
            errorMessage = "Please enter your username or email"
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
            // First, authenticate with Vignette
            let vignetteClient = VignetteAuthAPIClient.shared
            let vignetteAuth = try await vignetteClient.login(
                usernameOrEmail: usernameOrEmail,
                password: password
            )
            
            guard let vignetteToken = vignetteAuth.data?.accessToken else {
                throw CrossPlatformAuthError.serverError("Failed to get Vignette token")
            }
            
            // Then use Vignette token to sign in to Entativa
            let crossPlatformClient = CrossPlatformAuthClient.shared
            let entativaAuth = try await crossPlatformClient.signInWithVignette(
                vignetteToken: vignetteToken
            )
            
            isNewAccount = entativaAuth.data?.isNewAccount ?? false
            
            // Store the Entativa token
            if let token = entativaAuth.data?.accessToken {
                try KeychainManager.shared.save(token: token)
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
    SignInWithVignetteView { _ in }
}
