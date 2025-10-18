import SwiftUI
import LocalAuthentication

/// Vignette Login Screen - Instagram-inspired design
struct VignetteLoginView: View {
    @StateObject private var viewModel = VignetteAuthViewModel()
    @State private var showSignUp = false
    @State private var showPassword = false
    @FocusState private var focusedField: Field?
    
    enum Field: Hashable {
        case usernameOrEmail
        case password
    }
    
    var body: some View {
        NavigationView {
            ZStack {
                // Background
                VignetteColors.backgroundPrimary
                    .ignoresSafeArea()
                
                ScrollView {
                    VStack(spacing: 0) {
                        // Logo and branding
                        brandingSection
                            .padding(.top, 80)
                            .padding(.bottom, 48)
                        
                        // Login form
                        loginFormSection
                            .padding(.horizontal, 32)
                        
                        // Forgot password
                        forgotPasswordLink
                            .padding(.top, 20)
                        
                        // Sign up section
                        Spacer()
                            .frame(minHeight: 60)
                        
                        signUpSection
                    }
                    .padding(.bottom, 20)
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
                VignetteSignUpView()
            }
        }
    }
    
    // MARK: - Branding Section
    
    private var brandingSection: some View {
        VStack(spacing: 16) {
            // Logo text with Instagram-style script
            Text("Vignette")
                .font(.custom("Snell Roundhand", size: 52))
                .fontWeight(.medium)
                .foregroundColor(VignetteColors.textPrimary)
        }
    }
    
    // MARK: - Login Form Section
    
    private var loginFormSection: some View {
        VStack(spacing: 12) {
            // Username/Email field
            VignetteTextField(
                text: $viewModel.loginUsernameOrEmail,
                placeholder: "Username or email",
                keyboardType: .emailAddress,
                textContentType: .username,
                autocapitalization: .never
            )
            .focused($focusedField, equals: .usernameOrEmail)
            .onSubmit {
                focusedField = .password
            }
            
            // Password field
            VignetteSecureField(
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
            
            // Login button
            Button {
                Task {
                    await viewModel.login()
                }
            } label: {
                Text("Log In")
                    .vignetteButtonLarge()
                    .foregroundColor(VignetteColors.textOnPrimary)
                    .frame(maxWidth: .infinity)
                    .frame(height: 44)
                    .background(
                        viewModel.isLoading ?
                        VignetteColors.buttonPrimaryDisabled :
                        VignetteColors.buttonPrimary
                    )
                    .cornerRadius(8)
            }
            .disabled(viewModel.isLoading)
            .padding(.top, 8)
            
            // Biometric login (if available)
            if viewModel.biometricAuthAvailable {
                Button {
                    Task {
                        _ = await viewModel.authenticateWithBiometrics()
                    }
                } label: {
                    HStack(spacing: 8) {
                        Image(systemName: viewModel.biometricType == .faceID ? "faceid" : "touchid")
                            .font(.system(size: 18))
                        
                        Text("Log in with \(viewModel.biometricType == .faceID ? "Face ID" : "Touch ID")")
                            .vignetteButtonMedium()
                    }
                    .foregroundColor(VignetteColors.buttonPrimaryDeemphText)
                    .frame(maxWidth: .infinity)
                    .frame(height: 44)
                    .background(VignetteColors.buttonPrimaryDeemph)
                    .cornerRadius(8)
                }
                .padding(.top, 8)
            }
            
            // OR divider
            HStack(spacing: 16) {
                Rectangle()
                    .fill(VignetteColors.borderDefault)
                    .frame(height: 1)
                
                Text("OR")
                    .vignetteLabelSmall()
                    .foregroundColor(VignetteColors.textSecondary)
                
                Rectangle()
                    .fill(VignetteColors.borderDefault)
                    .frame(height: 1)
            }
            .padding(.vertical, 24)
            
            // Facebook login (placeholder)
            Button {
                // TODO: Implement Facebook OAuth
            } label: {
                HStack(spacing: 8) {
                    Image(systemName: "f.square.fill")
                        .font(.system(size: 20))
                    
                    Text("Log in with Facebook")
                        .vignetteButtonMedium()
                }
                .foregroundColor(VignetteColors.moonstone)
            }
        }
    }
    
    // MARK: - Forgot Password Link
    
    private var forgotPasswordLink: some View {
        Button {
            // TODO: Implement forgot password
        } label: {
            Text("Forgot password?")
                .vignetteLabelMedium()
                .foregroundColor(VignetteColors.moonstone)
        }
    }
    
    // MARK: - Sign Up Section
    
    private var signUpSection: some View {
        VStack(spacing: 16) {
            Rectangle()
                .fill(VignetteColors.borderDefault)
                .frame(height: 1)
                .padding(.horizontal, 32)
            
            HStack(spacing: 4) {
                Text("Don't have an account?")
                    .vignetteBodyMedium()
                    .foregroundColor(VignetteColors.textSecondary)
                
                Button {
                    showSignUp = true
                } label: {
                    Text("Sign up")
                        .vignetteBodyMedium()
                        .fontWeight(.semibold)
                        .foregroundColor(VignetteColors.moonstone)
                }
            }
        }
    }
    
    // MARK: - Loading Overlay
    
    private var loadingOverlay: some View {
        ZStack {
            Color.black.opacity(0.3)
                .ignoresSafeArea()
            
            ProgressView()
                .scaleEffect(1.5)
                .tint(VignetteColors.buttonPrimary)
                .padding(40)
                .background(
                    RoundedRectangle(cornerRadius: 12)
                        .fill(Color.white)
                        .shadow(radius: 10)
                )
        }
    }
}

// MARK: - Custom Text Field

struct VignetteTextField: View {
    @Binding var text: String
    let placeholder: String
    var keyboardType: UIKeyboardType = .default
    var textContentType: UITextContentType? = nil
    var autocapitalization: TextInputAutocapitalization = .sentences
    var error: String? = nil
    
    var body: some View {
        VStack(alignment: .leading, spacing: 6) {
            TextField(placeholder, text: $text)
                .vignetteBodyLarge()
                .foregroundColor(VignetteColors.textPrimary)
                .keyboardType(keyboardType)
                .textContentType(textContentType)
                .textInputAutocapitalization(autocapitalization)
                .padding(.horizontal, 16)
                .padding(.vertical, 12)
                .background(VignetteColors.backgroundSecondary)
                .cornerRadius(6)
                .overlay(
                    RoundedRectangle(cornerRadius: 6)
                        .stroke(
                            error != nil ? VignetteColors.borderError :
                            VignetteColors.borderDefault,
                            lineWidth: 1
                        )
                )
            
            if let error = error {
                Text(error)
                    .vignetteCaptionMedium()
                    .foregroundColor(VignetteColors.error)
            }
        }
    }
}

struct VignetteSecureField: View {
    @Binding var text: String
    let placeholder: String
    @Binding var showPassword: Bool
    var error: String? = nil
    
    var body: some View {
        VStack(alignment: .leading, spacing: 6) {
            HStack(spacing: 0) {
                if showPassword {
                    TextField(placeholder, text: $text)
                        .vignetteBodyLarge()
                        .foregroundColor(VignetteColors.textPrimary)
                        .textContentType(.password)
                        .textInputAutocapitalization(.never)
                } else {
                    SecureField(placeholder, text: $text)
                        .vignetteBodyLarge()
                        .foregroundColor(VignetteColors.textPrimary)
                        .textContentType(.password)
                        .textInputAutocapitalization(.never)
                }
                
                Button {
                    showPassword.toggle()
                } label: {
                    Image(systemName: showPassword ? "eye.slash.fill" : "eye.fill")
                        .font(.system(size: 14))
                        .foregroundColor(VignetteColors.textSecondary)
                        .frame(width: 40, height: 40)
                }
            }
            .padding(.leading, 16)
            .padding(.trailing, 4)
            .background(VignetteColors.backgroundSecondary)
            .cornerRadius(6)
            .overlay(
                RoundedRectangle(cornerRadius: 6)
                    .stroke(
                        error != nil ? VignetteColors.borderError :
                        VignetteColors.borderDefault,
                        lineWidth: 1
                    )
            )
            
            if let error = error {
                Text(error)
                    .vignetteCaptionMedium()
                    .foregroundColor(VignetteColors.error)
            }
        }
    }
}

#Preview {
    VignetteLoginView()
}
