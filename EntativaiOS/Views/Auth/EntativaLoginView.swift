import SwiftUI
import LocalAuthentication

/// Entativa Login Screen - Facebook-inspired design
struct EntativaLoginView: View {
    @StateObject private var viewModel = AuthViewModel()
    @State private var showSignUp = false
    @State private var showPassword = false
    @FocusState private var focusedField: Field?
    
    enum Field: Hashable {
        case emailOrUsername
        case password
    }
    
    var body: some View {
        NavigationView {
            ZStack {
                // Background
                EntativaColors.backgroundPrimary
                    .ignoresSafeArea()
                
                ScrollView {
                    VStack(spacing: 0) {
                        // Logo and branding
                        brandingSection
                            .padding(.top, 60)
                            .padding(.bottom, 40)
                        
                        // Login form
                        loginFormSection
                            .padding(.horizontal, 24)
                        
                        // Divider
                        HStack(spacing: 12) {
                            Rectangle()
                                .fill(EntativaColors.borderDefault)
                                .frame(height: 1)
                            
                            Text("OR")
                                .entativaLabelSmall()
                                .foregroundColor(EntativaColors.textSecondary)
                            
                            Rectangle()
                                .fill(EntativaColors.borderDefault)
                                .frame(height: 1)
                        }
                        .padding(.horizontal, 24)
                        .padding(.vertical, 28)
                        
                        // Biometric login (if available)
                        if viewModel.biometricAuthAvailable {
                            biometricLoginButton
                                .padding(.horizontal, 24)
                                .padding(.bottom, 20)
                        }
                        
                        // Sign up prompt
                        signUpPrompt
                            .padding(.horizontal, 24)
                            .padding(.top, 20)
                    }
                    .padding(.bottom, 40)
                }
                
                // Loading overlay
                if viewModel.isLoading {
                    loadingOverlay
                }
            }
            .alert("Error", isPresented: $viewModel.showError) {
                Button("OK", role: .cancel) {}
            } message: {
                if let errorMessage = viewModel.errorMessage {
                    Text(errorMessage)
                }
            }
            .fullScreenCover(isPresented: $showSignUp) {
                EntativaSignUpView()
            }
        }
    }
    
    // MARK: - Branding Section
    
    private var brandingSection: some View {
        VStack(spacing: 12) {
            // Logo
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
            
            Text("Connect with friends and the world around you")
                .entativaBodyMedium()
                .foregroundColor(EntativaColors.textSecondary)
                .multilineTextAlignment(.center)
                .padding(.horizontal, 40)
        }
    }
    
    // MARK: - Login Form Section
    
    private var loginFormSection: some View {
        VStack(spacing: 12) {
            // Email/Username field
            EntativaTextField(
                text: $viewModel.loginEmailOrUsername,
                placeholder: "Email address or username",
                keyboardType: .emailAddress,
                textContentType: .username,
                autocapitalization: .never
            )
            .focused($focusedField, equals: .emailOrUsername)
            .onSubmit {
                focusedField = .password
            }
            
            // Password field
            EntativaSecureField(
                text: $viewModel.loginPassword,
                placeholder: "Password",
                showPassword: $showPassword
            )
            .focused($focusedField, equals: .password)
            .onSubmit {
                Task {
                    await viewModel.login()
                }
            }
            
            // Forgot password link
            HStack {
                Spacer()
                Button {
                    // TODO: Implement forgot password
                } label: {
                    Text("Forgotten password?")
                        .entativaLabelMedium()
                        .foregroundColor(EntativaColors.textLink)
                }
            }
            .padding(.top, 4)
            
            // Login button
            Button {
                Task {
                    await viewModel.login()
                }
            } label: {
                Text("Log In")
                    .entativaButtonLarge()
                    .foregroundColor(EntativaColors.textOnPrimary)
                    .frame(maxWidth: .infinity)
                    .frame(height: 48)
                    .background(
                        viewModel.isLoading ?
                        EntativaColors.buttonPrimaryDisabled :
                        EntativaColors.buttonPrimary
                    )
                    .cornerRadius(8)
            }
            .disabled(viewModel.isLoading)
            .padding(.top, 12)
        }
    }
    
    // MARK: - Biometric Login Button
    
    private var biometricLoginButton: some View {
        Button {
            Task {
                _ = await viewModel.authenticateWithBiometrics()
            }
        } label: {
            HStack(spacing: 12) {
                Image(systemName: viewModel.biometricType == .faceID ? "faceid" : "touchid")
                    .font(.system(size: 20))
                
                Text("Log in with \(viewModel.biometricType == .faceID ? "Face ID" : "Touch ID")")
                    .entativaButtonMedium()
            }
            .foregroundColor(EntativaColors.buttonPrimaryDeemphText)
            .frame(maxWidth: .infinity)
            .frame(height: 48)
            .background(EntativaColors.buttonPrimaryDeemph)
            .cornerRadius(8)
        }
    }
    
    // MARK: - Sign Up Prompt
    
    private var signUpPrompt: some View {
        VStack(spacing: 16) {
            Rectangle()
                .fill(EntativaColors.borderDefault)
                .frame(height: 1)
            
            Text("Don't have an account?")
                .entativaBodyMedium()
                .foregroundColor(EntativaColors.textSecondary)
            
            Button {
                showSignUp = true
            } label: {
                Text("Create New Account")
                    .entativaButtonMedium()
                    .foregroundColor(EntativaColors.buttonPrimaryDeemphText)
                    .frame(maxWidth: .infinity)
                    .frame(height: 48)
                    .background(EntativaColors.buttonPrimaryDeemph)
                    .cornerRadius(8)
            }
        }
    }
    
    // MARK: - Loading Overlay
    
    private var loadingOverlay: some View {
        ZStack {
            Color.black.opacity(0.3)
                .ignoresSafeArea()
            
            VStack(spacing: 16) {
                ProgressView()
                    .scaleEffect(1.5)
                    .tint(EntativaColors.buttonPrimary)
                
                Text("Logging in...")
                    .entativaBodyMedium()
                    .foregroundColor(.white)
            }
            .padding(32)
            .background(
                RoundedRectangle(cornerRadius: 16)
                    .fill(Color.white.opacity(0.95))
            )
        }
    }
}

// MARK: - Custom Text Field

struct EntativaTextField: View {
    @Binding var text: String
    let placeholder: String
    var keyboardType: UIKeyboardType = .default
    var textContentType: UITextContentType? = nil
    var autocapitalization: TextInputAutocapitalization = .sentences
    var error: String? = nil
    
    var body: some View {
        VStack(alignment: .leading, spacing: 6) {
            TextField(placeholder, text: $text)
                .entativaBodyLarge()
                .foregroundColor(EntativaColors.textPrimary)
                .keyboardType(keyboardType)
                .textContentType(textContentType)
                .textInputAutocapitalization(autocapitalization)
                .padding(.horizontal, 16)
                .padding(.vertical, 14)
                .background(EntativaColors.backgroundSecondary)
                .cornerRadius(8)
                .overlay(
                    RoundedRectangle(cornerRadius: 8)
                        .stroke(error != nil ? EntativaColors.borderError : Color.clear, lineWidth: 1.5)
                )
            
            if let error = error {
                Text(error)
                    .entativaCaptionMedium()
                    .foregroundColor(EntativaColors.error)
            }
        }
    }
}

struct EntativaSecureField: View {
    @Binding var text: String
    let placeholder: String
    @Binding var showPassword: Bool
    var error: String? = nil
    
    var body: some View {
        VStack(alignment: .leading, spacing: 6) {
            HStack(spacing: 0) {
                if showPassword {
                    TextField(placeholder, text: $text)
                        .entativaBodyLarge()
                        .foregroundColor(EntativaColors.textPrimary)
                        .textContentType(.password)
                        .textInputAutocapitalization(.never)
                } else {
                    SecureField(placeholder, text: $text)
                        .entativaBodyLarge()
                        .foregroundColor(EntativaColors.textPrimary)
                        .textContentType(.password)
                        .textInputAutocapitalization(.never)
                }
                
                Button {
                    showPassword.toggle()
                } label: {
                    Image(systemName: showPassword ? "eye.slash.fill" : "eye.fill")
                        .font(.system(size: 16))
                        .foregroundColor(EntativaColors.textSecondary)
                        .frame(width: 44, height: 44)
                }
            }
            .padding(.leading, 16)
            .padding(.trailing, 4)
            .background(EntativaColors.backgroundSecondary)
            .cornerRadius(8)
            .overlay(
                RoundedRectangle(cornerRadius: 8)
                    .stroke(error != nil ? EntativaColors.borderError : Color.clear, lineWidth: 1.5)
            )
            
            if let error = error {
                Text(error)
                    .entativaCaptionMedium()
                    .foregroundColor(EntativaColors.error)
            }
        }
    }
}

#Preview {
    EntativaLoginView()
}
