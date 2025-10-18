package com.entativa.viewmodel

import android.app.Application
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.viewModelScope
import com.entativa.network.EntativaAuthAPIClient
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import java.text.SimpleDateFormat
import java.util.*
import java.util.regex.Pattern

/**
 * Authentication ViewModel for Entativa
 * Enterprise-grade implementation with proper state management
 */
class EntativaAuthViewModel(application: Application) : AndroidViewModel(application) {
    
    private val apiClient = EntativaAuthAPIClient.getInstance(application)
    
    // MARK: - UI State
    
    data class AuthUiState(
        val isAuthenticated: Boolean = false,
        val currentUser: EntativaAuthAPIClient.User? = null,
        val isLoading: Boolean = false,
        val errorMessage: String? = null,
        val showError: Boolean = false
    )
    
    private val _uiState = MutableStateFlow(AuthUiState())
    val uiState: StateFlow<AuthUiState> = _uiState.asStateFlow()
    
    // MARK: - Sign Up Form
    
    data class SignUpFormState(
        val firstName: String = "",
        val lastName: String = "",
        val email: String = "",
        val password: String = "",
        val birthday: Calendar = Calendar.getInstance().apply {
            add(Calendar.YEAR, -18)
        },
        val gender: String = "prefer_not_to_say",
        
        val firstNameError: String? = null,
        val lastNameError: String? = null,
        val emailError: String? = null,
        val passwordError: String? = null,
        val birthdayError: String? = null
    )
    
    private val _signUpFormState = MutableStateFlow(SignUpFormState())
    val signUpFormState: StateFlow<SignUpFormState> = _signUpFormState.asStateFlow()
    
    // MARK: - Login Form
    
    data class LoginFormState(
        val emailOrUsername: String = "",
        val password: String = ""
    )
    
    private val _loginFormState = MutableStateFlow(LoginFormState())
    val loginFormState: StateFlow<LoginFormState> = _loginFormState.asStateFlow()
    
    // MARK: - Initialization
    
    init {
        checkAuthenticationStatus()
    }
    
    // MARK: - Authentication Status
    
    fun checkAuthenticationStatus() {
        viewModelScope.launch {
            if (apiClient.hasToken()) {
                apiClient.getCurrentUser().fold(
                    onSuccess = { user ->
                        _uiState.value = _uiState.value.copy(
                            isAuthenticated = true,
                            currentUser = user
                        )
                    },
                    onFailure = {
                        _uiState.value = _uiState.value.copy(
                            isAuthenticated = false,
                            currentUser = null
                        )
                    }
                )
            }
        }
    }
    
    // MARK: - Sign Up
    
    fun updateSignUpField(field: SignUpField, value: Any) {
        _signUpFormState.value = when (field) {
            SignUpField.FIRST_NAME -> _signUpFormState.value.copy(
                firstName = value as String,
                firstNameError = null
            )
            SignUpField.LAST_NAME -> _signUpFormState.value.copy(
                lastName = value as String,
                lastNameError = null
            )
            SignUpField.EMAIL -> _signUpFormState.value.copy(
                email = value as String,
                emailError = null
            )
            SignUpField.PASSWORD -> _signUpFormState.value.copy(
                password = value as String,
                passwordError = null
            )
            SignUpField.BIRTHDAY -> _signUpFormState.value.copy(
                birthday = value as Calendar,
                birthdayError = null
            )
            SignUpField.GENDER -> _signUpFormState.value.copy(
                gender = value as String
            )
        }
    }
    
    fun signUp() {
        if (!validateSignUpForm()) return
        
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, errorMessage = null, showError = false)
            
            val form = _signUpFormState.value
            val dateFormat = SimpleDateFormat("yyyy-MM-dd", Locale.US)
            val birthdayStr = dateFormat.format(form.birthday.time)
            
            apiClient.signUp(
                firstName = form.firstName.trim(),
                lastName = form.lastName.trim(),
                email = form.email.trim().lowercase(),
                password = form.password,
                birthday = birthdayStr,
                gender = form.gender
            ).fold(
                onSuccess = { response ->
                    _uiState.value = _uiState.value.copy(
                        isAuthenticated = true,
                        currentUser = response.data?.user,
                        isLoading = false
                    )
                    clearSignUpForm()
                },
                onFailure = { error ->
                    _uiState.value = _uiState.value.copy(
                        errorMessage = error.message ?: "Sign up failed",
                        showError = true,
                        isLoading = false
                    )
                }
            )
        }
    }
    
    private fun validateSignUpForm(): Boolean {
        val form = _signUpFormState.value
        var isValid = true
        var errors = SignUpFormState()
        
        // Validate first name
        val trimmedFirstName = form.firstName.trim()
        when {
            trimmedFirstName.isEmpty() -> {
                errors = errors.copy(firstNameError = "First name is required")
                isValid = false
            }
            trimmedFirstName.length < 2 -> {
                errors = errors.copy(firstNameError = "First name must be at least 2 characters")
                isValid = false
            }
            !trimmedFirstName.all { it.isLetter() || it.isWhitespace() || it == '-' || it == '\'' } -> {
                errors = errors.copy(firstNameError = "First name can only contain letters")
                isValid = false
            }
        }
        
        // Validate last name
        val trimmedLastName = form.lastName.trim()
        when {
            trimmedLastName.isEmpty() -> {
                errors = errors.copy(lastNameError = "Last name is required")
                isValid = false
            }
            trimmedLastName.length < 2 -> {
                errors = errors.copy(lastNameError = "Last name must be at least 2 characters")
                isValid = false
            }
            !trimmedLastName.all { it.isLetter() || it.isWhitespace() || it == '-' || it == '\'' } -> {
                errors = errors.copy(lastNameError = "Last name can only contain letters")
                isValid = false
            }
        }
        
        // Validate email
        val trimmedEmail = form.email.trim()
        when {
            trimmedEmail.isEmpty() -> {
                errors = errors.copy(emailError = "Email is required")
                isValid = false
            }
            !isValidEmail(trimmedEmail) -> {
                errors = errors.copy(emailError = "Please enter a valid email address")
                isValid = false
            }
        }
        
        // Validate password
        when {
            form.password.isEmpty() -> {
                errors = errors.copy(passwordError = "Password is required")
                isValid = false
            }
            form.password.length < 8 -> {
                errors = errors.copy(passwordError = "Password must be at least 8 characters")
                isValid = false
            }
            !form.password.any { it.isUpperCase() } -> {
                errors = errors.copy(passwordError = "Password must contain at least one uppercase letter")
                isValid = false
            }
            !form.password.any { it.isLowerCase() } -> {
                errors = errors.copy(passwordError = "Password must contain at least one lowercase letter")
                isValid = false
            }
            !form.password.any { it.isDigit() } -> {
                errors = errors.copy(passwordError = "Password must contain at least one number")
                isValid = false
            }
        }
        
        // Validate age (must be 13+)
        val age = Calendar.getInstance().get(Calendar.YEAR) - form.birthday.get(Calendar.YEAR)
        when {
            age < 13 -> {
                errors = errors.copy(birthdayError = "You must be at least 13 years old to sign up")
                isValid = false
            }
            age > 120 -> {
                errors = errors.copy(birthdayError = "Please enter a valid birthday")
                isValid = false
            }
        }
        
        if (!isValid) {
            _signUpFormState.value = _signUpFormState.value.copy(
                firstNameError = errors.firstNameError,
                lastNameError = errors.lastNameError,
                emailError = errors.emailError,
                passwordError = errors.passwordError,
                birthdayError = errors.birthdayError
            )
        }
        
        return isValid
    }
    
    // MARK: - Login
    
    fun updateLoginField(field: LoginField, value: String) {
        _loginFormState.value = when (field) {
            LoginField.EMAIL_OR_USERNAME -> _loginFormState.value.copy(emailOrUsername = value)
            LoginField.PASSWORD -> _loginFormState.value.copy(password = value)
        }
    }
    
    fun login() {
        val form = _loginFormState.value
        
        if (form.emailOrUsername.trim().isEmpty()) {
            _uiState.value = _uiState.value.copy(
                errorMessage = "Please enter your email or username",
                showError = true
            )
            return
        }
        
        if (form.password.isEmpty()) {
            _uiState.value = _uiState.value.copy(
                errorMessage = "Please enter your password",
                showError = true
            )
            return
        }
        
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, errorMessage = null, showError = false)
            
            apiClient.login(
                emailOrUsername = form.emailOrUsername.trim(),
                password = form.password
            ).fold(
                onSuccess = { response ->
                    _uiState.value = _uiState.value.copy(
                        isAuthenticated = true,
                        currentUser = response.data?.user,
                        isLoading = false
                    )
                    clearLoginForm()
                },
                onFailure = { error ->
                    _uiState.value = _uiState.value.copy(
                        errorMessage = error.message ?: "Login failed",
                        showError = true,
                        isLoading = false
                    )
                }
            )
        }
    }
    
    // MARK: - Logout
    
    fun logout() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true)
            
            apiClient.logout()
            
            _uiState.value = _uiState.value.copy(
                isAuthenticated = false,
                currentUser = null,
                isLoading = false
            )
        }
    }
    
    // MARK: - Form Management
    
    private fun clearSignUpForm() {
        _signUpFormState.value = SignUpFormState()
    }
    
    private fun clearLoginForm() {
        _loginFormState.value = LoginFormState()
    }
    
    fun clearError() {
        _uiState.value = _uiState.value.copy(errorMessage = null, showError = false)
    }
    
    // MARK: - Validation Helpers
    
    private fun isValidEmail(email: String): Boolean {
        val emailPattern = Pattern.compile(
            "[a-zA-Z0-9._-]+@[a-z]+\\.+[a-z]+"
        )
        return emailPattern.matcher(email).matches()
    }
    
    // MARK: - Enums
    
    enum class SignUpField {
        FIRST_NAME, LAST_NAME, EMAIL, PASSWORD, BIRTHDAY, GENDER
    }
    
    enum class LoginField {
        EMAIL_OR_USERNAME, PASSWORD
    }
    
    // MARK: - Gender Options
    
    val genderOptions = listOf(
        "male" to "Male",
        "female" to "Female",
        "non_binary" to "Non-binary",
        "prefer_not_to_say" to "Prefer not to say",
        "custom" to "Custom"
    )
}
