import Foundation
import SwiftUI
import Combine
import LocalAuthentication

/// Authentication ViewModel for Entativa
@MainActor
class AuthViewModel: ObservableObject {
    // MARK: - Published Properties
    
    @Published var isAuthenticated = false
    @Published var currentUser: AuthAPIClient.AuthResponse.User?
    @Published var isLoading = false
    @Published var errorMessage: String?
    @Published var showError = false
    
    // Sign up form fields
    @Published var signUpFirstName = ""
    @Published var signUpLastName = ""
    @Published var signUpEmail = ""
    @Published var signUpPassword = ""
    @Published var signUpBirthday = Date()
    @Published var signUpGender = "prefer_not_to_say"
    
    // Login form fields
    @Published var loginEmailOrUsername = ""
    @Published var loginPassword = ""
    
    // Validation
    @Published var firstNameError: String?
    @Published var lastNameError: String?
    @Published var emailError: String?
    @Published var passwordError: String?
    @Published var birthdayError: String?
    
    // Biometric auth
    @Published var biometricAuthAvailable = false
    @Published var biometricType: LABiometryType = .none
    
    private let apiClient = AuthAPIClient.shared
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
                firstName: signUpFirstName.trimmingCharacters(in: .whitespaces),
                lastName: signUpLastName.trimmingCharacters(in: .whitespaces),
                email: signUpEmail.trimmingCharacters(in: .whitespaces).lowercased(),
                password: signUpPassword,
                birthday: signUpBirthday,
                gender: signUpGender
            )
            
            self.currentUser = response.data?.user
            self.isAuthenticated = true
            self.clearSignUpForm()
            
        } catch let error as AuthError {
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
        firstNameError = nil
        lastNameError = nil
        emailError = nil
        passwordError = nil
        birthdayError = nil
        
        // Validate first name
        let trimmedFirstName = signUpFirstName.trimmingCharacters(in: .whitespaces)
        if trimmedFirstName.isEmpty {
            firstNameError = "First name is required"
            isValid = false
        } else if trimmedFirstName.count < 2 {
            firstNameError = "First name must be at least 2 characters"
            isValid = false
        } else if !trimmedFirstName.allSatisfy({ $0.isLetter || $0.isWhitespace || $0 == "-" || $0 == "'" }) {
            firstNameError = "First name can only contain letters"
            isValid = false
        }
        
        // Validate last name
        let trimmedLastName = signUpLastName.trimmingCharacters(in: .whitespaces)
        if trimmedLastName.isEmpty {
            lastNameError = "Last name is required"
            isValid = false
        } else if trimmedLastName.count < 2 {
            lastNameError = "Last name must be at least 2 characters"
            isValid = false
        } else if !trimmedLastName.allSatisfy({ $0.isLetter || $0.isWhitespace || $0 == "-" || $0 == "'" }) {
            lastNameError = "Last name can only contain letters"
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
        
        // Validate age (must be 13+)
        let age = Calendar.current.dateComponents([.year], from: signUpBirthday, to: Date()).year ?? 0
        if age < 13 {
            birthdayError = "You must be at least 13 years old to sign up"
            isValid = false
        } else if age > 120 {
            birthdayError = "Please enter a valid birthday"
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
                emailOrUsername: loginEmailOrUsername.trimmingCharacters(in: .whitespaces),
                password: loginPassword
            )
            
            self.currentUser = response.data?.user
            self.isAuthenticated = true
            self.clearLoginForm()
            
        } catch let error as AuthError {
            self.errorMessage = error.errorDescription
            self.showError = true
        } catch {
            self.errorMessage = "An unexpected error occurred. Please try again."
            self.showError = true
        }
        
        isLoading = false
    }
    
    func validateLoginForm() -> Bool {
        if loginEmailOrUsername.trimmingCharacters(in: .whitespaces).isEmpty {
            errorMessage = "Please enter your email or username"
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
        signUpFirstName = ""
        signUpLastName = ""
        signUpEmail = ""
        signUpPassword = ""
        signUpBirthday = Date()
        signUpGender = "prefer_not_to_say"
        
        firstNameError = nil
        lastNameError = nil
        emailError = nil
        passwordError = nil
        birthdayError = nil
    }
    
    func clearLoginForm() {
        loginEmailOrUsername = ""
        loginPassword = ""
    }
    
    // MARK: - Validation Helpers
    
    func isValidEmail(_ email: String) -> Bool {
        let emailRegex = "[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,64}"
        let emailPredicate = NSPredicate(format:"SELF MATCHES %@", emailRegex)
        return emailPredicate.evaluate(with: email)
    }
    
    // MARK: - Gender Options
    
    var genderOptions: [(value: String, label: String)] {
        [
            ("male", "Male"),
            ("female", "Female"),
            ("non_binary", "Non-binary"),
            ("prefer_not_to_say", "Prefer not to say"),
            ("custom", "Custom")
        ]
    }
}
