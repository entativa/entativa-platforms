import Foundation
import SwiftUI
import Combine
import LocalAuthentication

/// Authentication ViewModel for Vignette
@MainActor
class VignetteAuthViewModel: ObservableObject {
    // MARK: - Published Properties
    
    @Published var isAuthenticated = false
    @Published var currentUser: VignetteAuthAPIClient.AuthResponse.User?
    @Published var isLoading = false
    @Published var errorMessage: String?
    @Published var showError = false
    
    // Sign up form fields
    @Published var signUpUsername = ""
    @Published var signUpEmail = ""
    @Published var signUpFullName = ""
    @Published var signUpPassword = ""
    
    // Login form fields
    @Published var loginUsernameOrEmail = ""
    @Published var loginPassword = ""
    
    // Validation
    @Published var usernameError: String?
    @Published var emailError: String?
    @Published var fullNameError: String?
    @Published var passwordError: String?
    
    // Biometric auth
    @Published var biometricAuthAvailable = false
    @Published var biometricType: LABiometryType = .none
    
    private let apiClient = VignetteAuthAPIClient.shared
    private let context = LAContext()
    private var cancellables = Set<AnyCancellable>()
    
    // MARK: - Initialization
    
    init() {
        checkBiometricAvailability()
        checkAuthenticationStatus()
    }
    
    // MARK: - Authentication Status
    
    func checkAuthenticationStatus() {
        Task {
            do {
                let user = try await apiClient.getCurrentUser()
                self.currentUser = user
                self.isAuthenticated = true
            } catch {
                self.isAuthenticated = false
                self.currentUser = nil
            }
        }
    }
    
    // MARK: - Sign Up
    
    func signUp() async {
        // Validate all fields
        guard validateSignUpForm() else {
            return
        }
        
        isLoading = true
        errorMessage = nil
        showError = false
        
        do {
            let response = try await apiClient.signUp(
                username: signUpUsername.trimmingCharacters(in: .whitespaces).lowercased(),
                email: signUpEmail.trimmingCharacters(in: .whitespaces).lowercased(),
                fullName: signUpFullName.trimmingCharacters(in: .whitespaces),
                password: signUpPassword
            )
            
            self.currentUser = response.data?.user
            self.isAuthenticated = true
            self.clearSignUpForm()
            
        } catch let error as VignetteAuthError {
            self.errorMessage = error.errorDescription
            self.showError = true
        } catch {
            self.errorMessage = "An unexpected error occurred. Please try again."
            self.showError = true
        }
        
        isLoading = false
    }
    
    func validateSignUpForm() -> Bool {
        var isValid = true
        
        // Reset errors
        usernameError = nil
        emailError = nil
        fullNameError = nil
        passwordError = nil
        
        // Validate username (Instagram-style rules)
        let trimmedUsername = signUpUsername.trimmingCharacters(in: .whitespaces).lowercased()
        if trimmedUsername.isEmpty {
            usernameError = "Username is required"
            isValid = false
        } else if trimmedUsername.count < 3 {
            usernameError = "Username must be at least 3 characters"
            isValid = false
        } else if trimmedUsername.count > 30 {
            usernameError = "Username must be 30 characters or less"
            isValid = false
        } else if !isValidUsername(trimmedUsername) {
            usernameError = "Username can only contain letters, numbers, periods, and underscores"
            isValid = false
        } else if trimmedUsername.hasPrefix(".") || trimmedUsername.hasSuffix(".") {
            usernameError = "Username cannot start or end with a period"
            isValid = false
        } else if trimmedUsername.contains("..") {
            usernameError = "Username cannot have consecutive periods"
            isValid = false
        }
        
        // Validate email
        let trimmedEmail = signUpEmail.trimmingCharacters(in: .whitespaces)
        if trimmedEmail.isEmpty {
            emailError = "Email is required"
            isValid = false
        } else if !isValidEmail(trimmedEmail) {
            emailError = "Please enter a valid email address"
            isValid = false
        }
        
        // Validate full name
        let trimmedFullName = signUpFullName.trimmingCharacters(in: .whitespaces)
        if trimmedFullName.isEmpty {
            fullNameError = "Full name is required"
            isValid = false
        } else if trimmedFullName.count < 2 {
            fullNameError = "Full name must be at least 2 characters"
            isValid = false
        }
        
        // Validate password
        if signUpPassword.isEmpty {
            passwordError = "Password is required"
            isValid = false
        } else if signUpPassword.count < 8 {
            passwordError = "Password must be at least 8 characters"
            isValid = false
        } else if !signUpPassword.contains(where: { $0.isUppercase }) {
            passwordError = "Password must contain at least one uppercase letter"
            isValid = false
        } else if !signUpPassword.contains(where: { $0.isLowercase }) {
            passwordError = "Password must contain at least one lowercase letter"
            isValid = false
        } else if !signUpPassword.contains(where: { $0.isNumber }) {
            passwordError = "Password must contain at least one number"
            isValid = false
        }
        
        return isValid
    }
    
    // MARK: - Login
    
    func login() async {
        // Validate login form
        guard validateLoginForm() else {
            return
        }
        
        isLoading = true
        errorMessage = nil
        showError = false
        
        do {
            let response = try await apiClient.login(
                usernameOrEmail: loginUsernameOrEmail.trimmingCharacters(in: .whitespaces),
                password: loginPassword
            )
            
            self.currentUser = response.data?.user
            self.isAuthenticated = true
            self.clearLoginForm()
            
        } catch let error as VignetteAuthError {
            self.errorMessage = error.errorDescription
            self.showError = true
        } catch {
            self.errorMessage = "An unexpected error occurred. Please try again."
            self.showError = true
        }
        
        isLoading = false
    }
    
    func validateLoginForm() -> Bool {
        if loginUsernameOrEmail.trimmingCharacters(in: .whitespaces).isEmpty {
            errorMessage = "Please enter your username or email"
            showError = true
            return false
        }
        
        if loginPassword.isEmpty {
            errorMessage = "Please enter your password"
            showError = true
            return false
        }
        
        return true
    }
    
    // MARK: - Biometric Authentication
    
    func checkBiometricAvailability() {
        var error: NSError?
        biometricAuthAvailable = context.canEvaluatePolicy(.deviceOwnerAuthenticationWithBiometrics, error: &error)
        biometricType = context.biometryType
    }
    
    func authenticateWithBiometrics() async -> Bool {
        let context = LAContext()
        var error: NSError?
        
        guard context.canEvaluatePolicy(.deviceOwnerAuthenticationWithBiometrics, error: &error) else {
            return false
        }
        
        do {
            let reason = biometricType == .faceID ? "Log in with Face ID" : "Log in with Touch ID"
            let success = try await context.evaluatePolicy(.deviceOwnerAuthenticationWithBiometrics, localizedReason: reason)
            
            if success {
                // Check if we have stored credentials
                checkAuthenticationStatus()
            }
            
            return success
        } catch {
            return false
        }
    }
    
    // MARK: - Logout
    
    func logout() async {
        isLoading = true
        
        do {
            try await apiClient.logout()
            self.isAuthenticated = false
            self.currentUser = nil
        } catch {
            print("Logout error: \(error)")
            // Still log out locally even if server request fails
            self.isAuthenticated = false
            self.currentUser = nil
        }
        
        isLoading = false
    }
    
    // MARK: - Form Management
    
    func clearSignUpForm() {
        signUpUsername = ""
        signUpEmail = ""
        signUpFullName = ""
        signUpPassword = ""
        
        usernameError = nil
        emailError = nil
        fullNameError = nil
        passwordError = nil
    }
    
    func clearLoginForm() {
        loginUsernameOrEmail = ""
        loginPassword = ""
    }
    
    // MARK: - Validation Helpers
    
    func isValidUsername(_ username: String) -> Bool {
        // Instagram-style username validation
        // Only letters, numbers, periods, and underscores
        let usernameRegex = "^[a-zA-Z0-9._]+$"
        let usernamePredicate = NSPredicate(format:"SELF MATCHES %@", usernameRegex)
        return usernamePredicate.evaluate(with: username)
    }
    
    func isValidEmail(_ email: String) -> Bool {
        let emailRegex = "[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,64}"
        let emailPredicate = NSPredicate(format:"SELF MATCHES %@", emailRegex)
        return emailPredicate.evaluate(with: email)
    }
}
