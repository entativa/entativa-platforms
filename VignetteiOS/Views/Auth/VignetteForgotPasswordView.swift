import SwiftUI

/// Vignette Forgot Password Screen
struct VignetteForgotPasswordView: View {
    @Environment(\.dismiss) private var dismiss
    @State private var email = ""
    @State private var isLoading = false
    @State private var showSuccess = false
    @State private var errorMessage: String?
    @State private var showError = false
    @FocusState private var isEmailFocused: Bool
    
    var body: some View {
        NavigationView {
            ZStack {
                VignetteColors.backgroundPrimary
                    .ignoresSafeArea()
                
                ScrollView {
                    VStack(spacing: 0) {
                        Spacer()
                            .frame(height: 40)
                        
                        // Icon
                        Image(systemName: "lock.circle.fill")
                            .font(.system(size: 80))
                            .foregroundColor(VignetteColors.gunmetal)
                        
                        Spacer()
                            .frame(height: 32)
                        
                        // Title
                        Text("Forgot password?")
                            .vignetteHeadlineSmall()
                            .foregroundColor(VignetteColors.textPrimary)
                        
                        Spacer()
                            .frame(height: 16)
                        
                        // Description
                        Text("Enter your email and we'll send you a link to reset your password.")
                            .vignetteBodyMedium()
                            .foregroundColor(VignetteColors.textSecondary)
                            .multilineTextAlignment(.center)
                            .padding(.horizontal, 40)
                        
                        Spacer()
                            .frame(height: 32)
                        
                        // Email field
                        VignetteTextField(
                            text: $email,
                            placeholder: "Email",
                            keyboardType: .emailAddress,
                            textContentType: .emailAddress,
                            autocapitalization: .never
                        )
                        .focused($isEmailFocused)
                        .padding(.horizontal, 32)
                        
                        Spacer()
                            .frame(height: 16)
                        
                        // Send button
                        Button {
                            Task {
                                await sendResetLink()
                            }
                        } label: {
                            Text("Send Reset Link")
                                .vignetteButtonLarge()
                                .foregroundColor(VignetteColors.textOnPrimary)
                                .frame(maxWidth: .infinity)
                                .frame(height: 44)
                                .background(
                                    isValidEmail(email) && !isLoading ?
                                    VignetteColors.buttonPrimary :
                                    VignetteColors.buttonPrimaryDisabled
                                )
                                .cornerRadius(8)
                        }
                        .disabled(!isValidEmail(email) || isLoading)
                        .padding(.horizontal, 32)
                        
                        Spacer()
                            .frame(height: 24)
                        
                        // Divider
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
                        
                        Spacer()
                            .frame(height: 24)
                        
                        // Create new account
                        Button {
                            dismiss()
                        } label: {
                            Text("Create New Account")
                                .vignetteButtonMedium()
                                .foregroundColor(VignetteColors.buttonPrimaryDeemphText)
                                .frame(maxWidth: .infinity)
                                .frame(height: 44)
                                .background(VignetteColors.buttonPrimaryDeemph)
                                .cornerRadius(8)
                        }
                        .padding(.horizontal, 32)
                        
                        Spacer()
                            .frame(height: 32)
                        
                        // Back to login
                        Button {
                            dismiss()
                        } label: {
                            HStack(spacing: 8) {
                                Image(systemName: "arrow.left")
                                    .font(.system(size: 13, weight: .semibold))
                                Text("Back to Login")
                                    .vignetteLabelLarge()
                            }
                            .foregroundColor(VignetteColors.textLink)
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
                            .font(.system(size: 16, weight: .medium))
                            .foregroundColor(VignetteColors.textPrimary)
                    }
                }
            }
            .alert("Success", isPresented: $showSuccess) {
                Button("OK") {
                    dismiss()
                }
            } message: {
                Text("We've sent a password reset link to \(email). Please check your inbox.")
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
    
    private func sendResetLink() async {
        isEmailFocused = false
        isLoading = true
        errorMessage = nil
        
        // Simulate API call
        do {
            try await Task.sleep(nanoseconds: 2_000_000_000) // 2 seconds
            
            // Call forgot password API endpoint
            let endpoint = "http://localhost:8002/api/v1/auth/forgot-password"
            guard let url = URL(string: endpoint) else {
                throw NSError(domain: "", code: -1, userInfo: [NSLocalizedDescriptionKey: "Invalid URL"])
            }
            
            var request = URLRequest(url: url)
            request.httpMethod = "POST"
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")
            
            let body = ["email": email]
            request.httpBody = try? JSONEncoder().encode(body)
            
            let (_, response) = try await URLSession.shared.data(for: request)
            
            guard let httpResponse = response as? HTTPURLResponse,
                  httpResponse.statusCode == 200 else {
                throw NSError(domain: "", code: -1, userInfo: [NSLocalizedDescriptionKey: "Failed to send reset link"])
            }
            
            showSuccess = true
        } catch {
            errorMessage = error.localizedDescription
            showError = true
        }
        
        isLoading = false
    }
    
    private func isValidEmail(_ email: String) -> Bool {
        let emailRegex = "[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,64}"
        let emailPredicate = NSPredicate(format:"SELF MATCHES %@", emailRegex)
        return emailPredicate.evaluate(with: email)
    }
}

#Preview {
    VignetteForgotPasswordView()
}
