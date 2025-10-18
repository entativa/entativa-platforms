import SwiftUI

/// Entativa Sign Up Screen - Facebook-inspired design
struct EntativaSignUpView: View {
    @StateObject private var viewModel = AuthViewModel()
    @Environment(\.dismiss) private var dismiss
    @State private var showPassword = false
    @State private var currentStep = 1
    @FocusState private var focusedField: Field?
    
    enum Field: Hashable {
        case firstName
        case lastName
        case email
        case password
    }
    
    var body: some View {
        NavigationView {
            ZStack {
                // Background
                EntativaColors.backgroundPrimary
                    .ignoresSafeArea()
                
                ScrollView {
                    VStack(spacing: 24) {
                        // Header
                        headerSection
                            .padding(.top, 24)
                        
                        // Progress indicator
                        progressIndicator
                            .padding(.horizontal, 24)
                        
                        // Form sections
                        if currentStep == 1 {
                            nameSection
                                .padding(.horizontal, 24)
                                .transition(.asymmetric(
                                    insertion: .move(edge: .trailing),
                                    removal: .move(edge: .leading)
                                ))
                        } else if currentStep == 2 {
                            emailPasswordSection
                                .padding(.horizontal, 24)
                                .transition(.asymmetric(
                                    insertion: .move(edge: .trailing),
                                    removal: .move(edge: .leading)
                                ))
                        } else if currentStep == 3 {
                            birthdayGenderSection
                                .padding(.horizontal, 24)
                                .transition(.asymmetric(
                                    insertion: .move(edge: .trailing),
                                    removal: .move(edge: .leading)
                                ))
                        }
                        
                        // Navigation buttons
                        navigationButtons
                            .padding(.horizontal, 24)
                            .padding(.top, 20)
                        
                        // Terms and privacy
                        termsSection
                            .padding(.horizontal, 24)
                            .padding(.top, 12)
                    }
                    .padding(.bottom, 40)
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
                            .font(.system(size: 16, weight: .semibold))
                            .foregroundColor(EntativaColors.textPrimary)
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
    
    // MARK: - Header Section
    
    private var headerSection: some View {
        VStack(spacing: 8) {
            Text("Create Account")
                .entativaHeadlineMedium()
                .foregroundColor(EntativaColors.textPrimary)
            
            Text(stepDescription)
                .entativaBodyMedium()
                .foregroundColor(EntativaColors.textSecondary)
                .multilineTextAlignment(.center)
        }
    }
    
    private var stepDescription: String {
        switch currentStep {
        case 1:
            return "What's your name?"
        case 2:
            return "Enter your email and password"
        case 3:
            return "Tell us about yourself"
        default:
            return ""
        }
    }
    
    // MARK: - Progress Indicator
    
    private var progressIndicator: some View {
        HStack(spacing: 8) {
            ForEach(1...3, id: \.self) { step in
                Capsule()
                    .fill(step <= currentStep ? EntativaColors.buttonPrimary : EntativaColors.borderDefault)
                    .frame(height: 4)
            }
        }
    }
    
    // MARK: - Name Section (Step 1)
    
    private var nameSection: some View {
        VStack(spacing: 16) {
            EntativaTextField(
                text: $viewModel.signUpFirstName,
                placeholder: "First name",
                keyboardType: .namePhonePad,
                textContentType: .givenName,
                autocapitalization: .words,
                error: viewModel.firstNameError
            )
            .focused($focusedField, equals: .firstName)
            .onSubmit {
                focusedField = .lastName
            }
            
            EntativaTextField(
                text: $viewModel.signUpLastName,
                placeholder: "Last name",
                keyboardType: .namePhonePad,
                textContentType: .familyName,
                autocapitalization: .words,
                error: viewModel.lastNameError
            )
            .focused($focusedField, equals: .lastName)
            .onSubmit {
                if canProceedToStep2() {
                    withAnimation {
                        currentStep = 2
                        focusedField = .email
                    }
                }
            }
        }
    }
    
    // MARK: - Email & Password Section (Step 2)
    
    private var emailPasswordSection: some View {
        VStack(spacing: 16) {
            EntativaTextField(
                text: $viewModel.signUpEmail,
                placeholder: "Email address",
                keyboardType: .emailAddress,
                textContentType: .emailAddress,
                autocapitalization: .never,
                error: viewModel.emailError
            )
            .focused($focusedField, equals: .email)
            .onSubmit {
                focusedField = .password
            }
            
            EntativaSecureField(
                text: $viewModel.signUpPassword,
                placeholder: "Password",
                showPassword: $showPassword,
                error: viewModel.passwordError
            )
            .focused($focusedField, equals: .password)
            
            // Password requirements
            VStack(alignment: .leading, spacing: 6) {
                PasswordRequirement(
                    text: "At least 8 characters",
                    isMet: viewModel.signUpPassword.count >= 8
                )
                PasswordRequirement(
                    text: "Contains uppercase letter",
                    isMet: viewModel.signUpPassword.contains(where: { $0.isUppercase })
                )
                PasswordRequirement(
                    text: "Contains lowercase letter",
                    isMet: viewModel.signUpPassword.contains(where: { $0.isLowercase })
                )
                PasswordRequirement(
                    text: "Contains number",
                    isMet: viewModel.signUpPassword.contains(where: { $0.isNumber })
                )
            }
            .padding(.horizontal, 4)
        }
    }
    
    // MARK: - Birthday & Gender Section (Step 3)
    
    private var birthdayGenderSection: some View {
        VStack(spacing: 16) {
            // Birthday picker
            VStack(alignment: .leading, spacing: 8) {
                Text("Birthday")
                    .entativaLabelLarge()
                    .foregroundColor(EntativaColors.textPrimary)
                
                DatePicker(
                    "",
                    selection: $viewModel.signUpBirthday,
                    in: ...Date(),
                    displayedComponents: .date
                )
                .datePickerStyle(.wheel)
                .labelsHidden()
                .frame(maxWidth: .infinity)
                .background(EntativaColors.backgroundSecondary)
                .cornerRadius(12)
                
                if let error = viewModel.birthdayError {
                    Text(error)
                        .entativaCaptionMedium()
                        .foregroundColor(EntativaColors.error)
                }
                
                Text("You must be at least 13 years old")
                    .entativaCaptionMedium()
                    .foregroundColor(EntativaColors.textSecondary)
            }
            
            // Gender picker
            VStack(alignment: .leading, spacing: 8) {
                Text("Gender")
                    .entativaLabelLarge()
                    .foregroundColor(EntativaColors.textPrimary)
                
                Picker("Gender", selection: $viewModel.signUpGender) {
                    ForEach(viewModel.genderOptions, id: \.value) { option in
                        Text(option.label)
                            .tag(option.value)
                    }
                }
                .pickerStyle(.segmented)
                
                Text("You can always change this later")
                    .entativaCaptionMedium()
                    .foregroundColor(EntativaColors.textSecondary)
            }
        }
    }
    
    // MARK: - Navigation Buttons
    
    private var navigationButtons: some View {
        HStack(spacing: 12) {
            if currentStep > 1 {
                Button {
                    withAnimation {
                        currentStep -= 1
                    }
                } label: {
                    Text("Back")
                        .entativaButtonMedium()
                        .foregroundColor(EntativaColors.buttonSecondaryText)
                        .frame(maxWidth: .infinity)
                        .frame(height: 48)
                        .background(EntativaColors.buttonSecondary)
                        .cornerRadius(8)
                }
            }
            
            Button {
                handleNextButton()
            } label: {
                Text(currentStep == 3 ? "Sign Up" : "Next")
                    .entativaButtonLarge()
                    .foregroundColor(EntativaColors.textOnPrimary)
                    .frame(maxWidth: .infinity)
                    .frame(height: 48)
                    .background(
                        canProceed() ?
                        EntativaColors.buttonPrimary :
                        EntativaColors.buttonPrimaryDisabled
                    )
                    .cornerRadius(8)
            }
            .disabled(!canProceed() || viewModel.isLoading)
        }
    }
    
    // MARK: - Terms Section
    
    private var termsSection: some View {
        VStack(spacing: 8) {
            Text("By signing up, you agree to our")
                .entativaCaptionMedium()
                .foregroundColor(EntativaColors.textSecondary)
            +
            Text(" Terms")
                .entativaCaptionMedium()
                .foregroundColor(EntativaColors.textLink)
            +
            Text(", ")
                .entativaCaptionMedium()
                .foregroundColor(EntativaColors.textSecondary)
            +
            Text("Privacy Policy")
                .entativaCaptionMedium()
                .foregroundColor(EntativaColors.textLink)
            +
            Text(" and ")
                .entativaCaptionMedium()
                .foregroundColor(EntativaColors.textSecondary)
            +
            Text("Cookies Policy")
                .entativaCaptionMedium()
                .foregroundColor(EntativaColors.textLink)
        }
        .multilineTextAlignment(.center)
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
                
                Text("Creating your account...")
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
    
    // MARK: - Helper Methods
    
    private func canProceedToStep2() -> Bool {
        let trimmedFirstName = viewModel.signUpFirstName.trimmingCharacters(in: .whitespaces)
        let trimmedLastName = viewModel.signUpLastName.trimmingCharacters(in: .whitespaces)
        return !trimmedFirstName.isEmpty && !trimmedLastName.isEmpty
    }
    
    private func canProceedToStep3() -> Bool {
        let trimmedEmail = viewModel.signUpEmail.trimmingCharacters(in: .whitespaces)
        return !trimmedEmail.isEmpty && !viewModel.signUpPassword.isEmpty
    }
    
    private func canProceed() -> Bool {
        switch currentStep {
        case 1:
            return canProceedToStep2()
        case 2:
            return canProceedToStep3()
        case 3:
            return true
        default:
            return false
        }
    }
    
    private func handleNextButton() {
        if currentStep < 3 {
            withAnimation {
                currentStep += 1
            }
            
            // Set focus on next step's first field
            DispatchQueue.main.asyncAfter(deadline: .now() + 0.5) {
                if currentStep == 2 {
                    focusedField = .email
                }
            }
        } else {
            // Final step - sign up
            Task {
                await viewModel.signUp()
                if viewModel.isAuthenticated {
                    dismiss()
                }
            }
        }
    }
}

// MARK: - Password Requirement View

struct PasswordRequirement: View {
    let text: String
    let isMet: Bool
    
    var body: some View {
        HStack(spacing: 8) {
            Image(systemName: isMet ? "checkmark.circle.fill" : "circle")
                .font(.system(size: 14))
                .foregroundColor(isMet ? EntativaColors.success : EntativaColors.textSecondary)
            
            Text(text)
                .entativaCaptionMedium()
                .foregroundColor(isMet ? EntativaColors.textPrimary : EntativaColors.textSecondary)
        }
    }
}

#Preview {
    EntativaSignUpView()
}
