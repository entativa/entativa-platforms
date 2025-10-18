import SwiftUI

/// Entativa Forgot Password Screen
struct EntativaForgotPasswordView: View {
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
                EntativaColors.backgroundPrimary
                    .ignoresSafeArea()
                
                ScrollView {
                    VStack(spacing: 0) {
                        Spacer()
                            .frame(height: 40)
                        
                        // Icon
                        Image(systemName: "lock.circle.fill")
                            .font(.system(size: 96))
                            .foregroundStyle(
                                LinearGradient(
                                    colors: [
                                        EntativaColors.gradientStart,
                                        EntativaColors.gradientMiddle
                                    ],
                                    startPoint: .topLeading,
                                    endPoint: .bottomTrailing
                                )
                            )
                        
                        Spacer()
                            .frame(height: 32)
                        
                        // Title
                        Text("Trouble logging in?")
                            .entativaHeadlineSmall()
                            .foregroundColor(EntativaColors.textPrimary)
                        
                        Spacer()
                            .frame(height: 16)
                        
                        // Description
                        Text("Enter your email and we'll send you a link to get back into your account.")
                            .entativaBodyMedium()
                            .foregroundColor(EntativaColors.textSecondary)
                            .multilineTextAlignment(.center)
                            .padding(.horizontal, 40)
                        
                        Spacer()
                            .frame(height: 32)
                        
                        // Email field
                        EntativaTextField(
                            text: $email,
                            placeholder: "Email address",
                            keyboardType: .emailAddress,
                            textContentType: .emailAddress,
                            autocapitalization: .never
                        )
                        .focused($isEmailFocused)
                        .padding(.horizontal, 24)
                        
                        Spacer()
                            .frame(height: 16)
                        
                        // Send button
                        Button {
                            Task {
                                await sendResetLink()
                            }
                        } label: {
                            Text("Send Reset Link")
                                .entativaButtonLarge()
                                .foregroundColor(EntativaColors.textOnPrimary)
                                .frame(maxWidth: .infinity)
                                .frame(height: 48)
                                .background(
                                    isValidEmail(email) && !isLoading ?
                                    EntativaColors.buttonPrimary :
                                    EntativaColors.buttonPrimaryDisabled
                                )
                                .cornerRadius(8)
                        }
                        .disabled(!isValidEmail(email) || isLoading)
                        .padding(.horizontal, 24)
                        
                        Spacer()
                            .frame(height: 24)
                        
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
                        
                        Spacer()
                            .frame(height: 24)
                        
                        // Create new account
                        Button {
                            dismiss()
                        } label: {
                            Text("Create New Account")
                                .entativaButtonMedium()
                                .foregroundColor(EntativaColors.buttonPrimaryDeemphText)
                                .frame(maxWidth: .infinity)
                                .frame(height: 48)
                                .background(EntativaColors.buttonPrimaryDeemph)
                                .cornerRadius(8)
                        }
                        .padding(.horizontal, 24)
                        
                        Spacer()
                            .frame(height: 32)
                        
                        // Back to login
                        Button {
                            dismiss()
                        } label: {
                            HStack(spacing: 8) {
                                Image(systemName: "arrow.left")
                                    .font(.system(size: 14, weight: .semibold))
                                Text("Back to Login")
                                    .entativaLabelLarge()
                            }
                            .foregroundColor(EntativaColors.textLink)
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
            
            VStack(spacing: 16) {
                ProgressView()
                    .scaleEffect(1.5)
                    .tint(EntativaColors.buttonPrimary)
                
                Text("Sending reset link...")
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
    
    private func sendResetLink() async {
        isEmailFocused = false
        isLoading = true
        errorMessage = nil
        
        // Simulate API call
        do {
            try await Task.sleep(nanoseconds: 2_000_000_000) // 2 seconds
            
            // Call forgot password API endpoint
            let endpoint = "http://localhost:8001/api/v1/auth/forgot-password"
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
    EntativaForgotPasswordView()
}
