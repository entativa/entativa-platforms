import SwiftUI

/// Vignette Sign Up Screen - Instagram-inspired design
struct VignetteSignUpView: View {
    @StateObject private var viewModel = VignetteAuthViewModel()
    @Environment(\.dismiss) private var dismiss
    @State private var showPassword = false
    @FocusState private var focusedField: Field?
    
    enum Field: Hashable {
        case email
        case fullName
        case username
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
                        // Logo
                        logoSection
                            .padding(.top, 40)
                            .padding(.bottom, 32)
                        
                        // Sign up form
                        signUpFormSection
                            .padding(.horizontal, 32)
                        
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
                        .padding(.horizontal, 32)
                        .padding(.vertical, 24)
                        
                // Sign in with Entativa
                signInWithEntativaButton
                    .padding(.horizontal, 32)
                        
                        // Terms
                        termsSection
                            .padding(.horizontal, 40)
                            .padding(.top, 32)
                        
                        // Login link
                        Spacer()
                            .frame(minHeight: 40)
                        
                        loginSection
                    }
                    .padding(.bottom, 20)
                }
                
                // Loading overlay
                if viewModel.isLoading {
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
                            .font(.system(size: 16, weight: .medium))
                            .foregroundColor(VignetteColors.textPrimary)
                    }
                }
            }
            .alert("Error", isPresented: $viewModel.showError) {
                Button("OK", role: .cancel) {}
            } message: {
                if let errorMessage = viewModel.errorMessage {
                    Text(errorMessage)
                }
            }
        }
    }
    
    // MARK: - Logo Section
    
    private var logoSection: some View {
        VStack(spacing: 16) {
            Text("Vignette")
                .font(.custom("Snell Roundhand", size: 48))
                .fontWeight(.medium)
                .foregroundColor(VignetteColors.textPrimary)
            
            Text("Sign up to see photos and videos from your friends.")
                .vignetteBodyMedium()
                .foregroundColor(VignetteColors.textSecondary)
                .multilineTextAlignment(.center)
        }
    }
    
    // MARK: - Sign Up Form Section
    
    private var signUpFormSection: some View {
        VStack(spacing: 12) {
            // Email field
            VignetteTextField(
                text: $viewModel.signUpEmail,
                placeholder: "Email",
                keyboardType: .emailAddress,
                textContentType: .emailAddress,
                autocapitalization: .never,
                error: viewModel.emailError
            )
            .focused($focusedField, equals: .email)
            .onSubmit {
                focusedField = .fullName
            }
            
            // Full name field
            VignetteTextField(
                text: $viewModel.signUpFullName,
                placeholder: "Full Name",
                keyboardType: .namePhonePad,
                textContentType: .name,
                autocapitalization: .words,
                error: viewModel.fullNameError
            )
            .focused($focusedField, equals: .fullName)
            .onSubmit {
                focusedField = .username
            }
            
            // Username field
            VignetteTextField(
                text: $viewModel.signUpUsername,
                placeholder: "Username",
                keyboardType: .asciiCapable,
                textContentType: .username,
                autocapitalization: .never,
                error: viewModel.usernameError
            )
            .focused($focusedField, equals: .username)
            .onSubmit {
                focusedField = .password
            }
            .onChange(of: viewModel.signUpUsername) { oldValue, newValue in
                // Convert to lowercase automatically
                viewModel.signUpUsername = newValue.lowercased()
            }
            
            // Username hint
            if viewModel.usernameError == nil && !viewModel.signUpUsername.isEmpty {
                Text("Can contain letters, numbers, periods, and underscores")
                    .vignetteCaptionSmall()
                    .foregroundColor(VignetteColors.textSecondary)
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .padding(.horizontal, 4)
            }
            
            // Password field
            VignetteSecureField(
                text: $viewModel.signUpPassword,
                placeholder: "Password",
                showPassword: $showPassword,
                error: viewModel.passwordError
            )
            .focused($focusedField, equals: .password)
            
            // Password requirements (only show if typing)
            if !viewModel.signUpPassword.isEmpty {
                VStack(alignment: .leading, spacing: 4) {
                    PasswordRequirementRow(
                        text: "8+ characters",
                        isMet: viewModel.signUpPassword.count >= 8
                    )
                    PasswordRequirementRow(
                        text: "Uppercase & lowercase",
                        isMet: viewModel.signUpPassword.contains(where: { $0.isUppercase }) &&
                               viewModel.signUpPassword.contains(where: { $0.isLowercase })
                    )
                    PasswordRequirementRow(
                        text: "Contains a number",
                        isMet: viewModel.signUpPassword.contains(where: { $0.isNumber })
                    )
                }
                .padding(.horizontal, 4)
            }
            
            // Sign up button
            Button {
                Task {
                    await viewModel.signUp()
                    if viewModel.isAuthenticated {
                        dismiss()
                    }
                }
            } label: {
                Text("Sign Up")
                    .vignetteButtonLarge()
                    .foregroundColor(VignetteColors.textOnPrimary)
                    .frame(maxWidth: .infinity)
                    .frame(height: 44)
                    .background(
                        canSignUp() ?
                        VignetteColors.buttonPrimary :
                        VignetteColors.buttonPrimaryDisabled
                    )
                    .cornerRadius(8)
            }
            .disabled(!canSignUp() || viewModel.isLoading)
            .padding(.top, 8)
        }
    }
    
    // MARK: - Sign in with Entativa Button
    
    private var signInWithEntativaButton: some View {
        Button {
            // Navigate to Entativa sign-in flow
        } label: {
            HStack(spacing: 8) {
                Text("e")
                    .font(.system(size: 20, weight: .bold, design: .rounded))
                    .italic()
                
                Text("Sign in with Entativa")
                    .vignetteButtonMedium()
            }
            .foregroundColor(EntativaColors.primaryBlue)
            .frame(maxWidth: .infinity)
            .frame(height: 44)
            .background(EntativaColors.buttonPrimaryDeemph)
            .cornerRadius(8)
        }
    }
    
    // MARK: - Terms Section
    
    private var termsSection: some View {
        Group {
            Text("By signing up, you agree to our ")
                .vignetteCaptionMedium()
                .foregroundColor(VignetteColors.textSecondary)
            +
            Text("Terms")
                .vignetteCaptionMedium()
                .foregroundColor(VignetteColors.moonstone)
            +
            Text(", ")
                .vignetteCaptionMedium()
                .foregroundColor(VignetteColors.textSecondary)
            +
            Text("Privacy Policy")
                .vignetteCaptionMedium()
                .foregroundColor(VignetteColors.moonstone)
            +
            Text(" and ")
                .vignetteCaptionMedium()
                .foregroundColor(VignetteColors.textSecondary)
            +
            Text("Cookies Policy")
                .vignetteCaptionMedium()
                .foregroundColor(VignetteColors.moonstone)
            +
            Text(".")
                .vignetteCaptionMedium()
                .foregroundColor(VignetteColors.textSecondary)
        }
        .multilineTextAlignment(.center)
    }
    
    // MARK: - Login Section
    
    private var loginSection: some View {
        VStack(spacing: 16) {
            Rectangle()
                .fill(VignetteColors.borderDefault)
                .frame(height: 1)
                .padding(.horizontal, 32)
            
            HStack(spacing: 4) {
                Text("Have an account?")
                    .vignetteBodyMedium()
                    .foregroundColor(VignetteColors.textSecondary)
                
                Button {
                    dismiss()
                } label: {
                    Text("Log in")
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
    
    // MARK: - Helper Methods
    
    private func canSignUp() -> Bool {
        let trimmedEmail = viewModel.signUpEmail.trimmingCharacters(in: .whitespaces)
        let trimmedFullName = viewModel.signUpFullName.trimmingCharacters(in: .whitespaces)
        let trimmedUsername = viewModel.signUpUsername.trimmingCharacters(in: .whitespaces)
        
        return !trimmedEmail.isEmpty &&
               !trimmedFullName.isEmpty &&
               !trimmedUsername.isEmpty &&
               !viewModel.signUpPassword.isEmpty &&
               viewModel.signUpPassword.count >= 8
    }
}

// MARK: - Password Requirement Row

struct PasswordRequirementRow: View {
    let text: String
    let isMet: Bool
    
    var body: some View {
        HStack(spacing: 6) {
            Image(systemName: isMet ? "checkmark.circle.fill" : "circle")
                .font(.system(size: 12))
                .foregroundColor(isMet ? VignetteColors.success : VignetteColors.textTertiary)
            
            Text(text)
                .vignetteCaptionSmall()
                .foregroundColor(isMet ? VignetteColors.textPrimary : VignetteColors.textSecondary)
        }
    }
}

#Preview {
    VignetteSignUpView()
}
