package com.entativa.vignette.viewmodel

import android.app.Application
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.viewModelScope
import com.entativa.vignette.network.VignetteAuthAPIClient
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import java.util.regex.Pattern

/**
 * Authentication ViewModel for Vignette
 * Instagram-style authentication with enterprise-grade state management
 */
class VignetteAuthViewModel(application: Application) : AndroidViewModel(application) {
    
    private val apiClient = VignetteAuthAPIClient.getInstance(application)
    
    // MARK: - UI State
    
    data class AuthUiState(
        val isAuthenticated: Boolean = false,
        val currentUser: VignetteAuthAPIClient.User? = null,
        val isLoading: Boolean = false,
        val errorMessage: String? = null,
        val showError: Boolean = false
    )
    
    private val _uiState = MutableStateFlow(AuthUiState())
    val uiState: StateFlow<AuthUiState> = _uiState.asStateFlow()
    
    // MARK: - Sign Up Form
    
    data class SignUpFormState(
        val username: String = "",
        val email: String = "",
        val fullName: String = "",
        val password: String = "",
        
        val usernameError: String? = null,
        val emailError: String? = null,
        val fullNameError: String? = null,
        val passwordError: String? = null
    )
    
    private val _signUpFormState = MutableStateFlow(SignUpFormState())
    val signUpFormState: StateFlow<SignUpFormState> = _signUpFormState.asStateFlow()
    
    // MARK: - Login Form
    
    data class LoginFormState(
        val usernameOrEmail: String = "",
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
    
    fun updateSignUpField(field: SignUpField, value: String) {
        _signUpFormState.value = when (field) {
            SignUpField.USERNAME -> _signUpFormState.value.copy(
                username = value.lowercase(),
                usernameError = null
            )
            SignUpField.EMAIL -> _signUpFormState.value.copy(
                email = value,
                emailError = null
            )
            SignUpField.FULL_NAME -> _signUpFormState.value.copy(
                fullName = value,
                fullNameError = null
            )
            SignUpField.PASSWORD -> _signUpFormState.value.copy(
                password = value,
                passwordError = null
            )
        }
    }
    
    fun signUp() {
        if (!validateSignUpForm()) return
        
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, errorMessage = null, showError = false)
            
            val form = _signUpFormState.value
            
            apiClient.signUp(
                username = form.username.trim().lowercase(),
                email = form.email.trim().lowercase(),
                fullName = form.fullName.trim(),
                password = form.password
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
        
        // Validate username (Instagram-style rules)
        val trimmedUsername = form.username.trim().lowercase()
        when {
            trimmedUsername.isEmpty() -> {
                errors = errors.copy(usernameError = "Username is required")
                isValid = false
            }
            trimmedUsername.length < 3 -> {
                errors = errors.copy(usernameError = "Username must be at least 3 characters")
                isValid = false
            }
            trimmedUsername.length > 30 -> {
                errors = errors.copy(usernameError = "Username must be 30 characters or less")
                isValid = false
            }
            !isValidUsername(trimmedUsername) -> {
                errors = errors.copy(usernameError = "Username can only contain letters, numbers, periods, and underscores")
                isValid = false
            }
            trimmedUsername.startsWith(".") || trimmedUsername.endsWith(".") -> {
                errors = errors.copy(usernameError = "Username cannot start or end with a period")
                isValid = false
            }
            trimmedUsername.contains("..") -> {
                errors = errors.copy(usernameError = "Username cannot have consecutive periods")
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
        
        // Validate full name
        val trimmedFullName = form.fullName.trim()
        when {
            trimmedFullName.isEmpty() -> {
                errors = errors.copy(fullNameError = "Full name is required")
                isValid = false
            }
            trimmedFullName.length < 2 -> {
                errors = errors.copy(fullNameError = "Full name must be at least 2 characters")
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
        
        if (!isValid) {
            _signUpFormState.value = _signUpFormState.value.copy(
                usernameError = errors.usernameError,
                emailError = errors.emailError,
                fullNameError = errors.fullNameError,
                passwordError = errors.passwordError
            )
        }
        
        return isValid
    }
    
    // MARK: - Login
    
    fun updateLoginField(field: LoginField, value: String) {
        _loginFormState.value = when (field) {
            LoginField.USERNAME_OR_EMAIL -> _loginFormState.value.copy(usernameOrEmail = value)
            LoginField.PASSWORD -> _loginFormState.value.copy(password = value)
        }
    }
    
    fun login() {
        val form = _loginFormState.value
        
        if (form.usernameOrEmail.trim().isEmpty()) {
            _uiState.value = _uiState.value.copy(
                errorMessage = "Please enter your username or email",
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
                usernameOrEmail = form.usernameOrEmail.trim(),
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
    
    private fun isValidUsername(username: String): Boolean {
        // Instagram-style username validation: only letters, numbers, periods, and underscores
        val usernamePattern = Pattern.compile("^[a-zA-Z0-9._]+$")
        return usernamePattern.matcher(username).matches()
    }
    
    private fun isValidEmail(email: String): Boolean {
        val emailPattern = Pattern.compile(
            "[a-zA-Z0-9._-]+@[a-z]+\\.+[a-z]+"
        )
        return emailPattern.matcher(email).matches()
    }
    
    // MARK: - Enums
    
    enum class SignUpField {
        USERNAME, EMAIL, FULL_NAME, PASSWORD
    }
    
    enum class LoginField {
        USERNAME_OR_EMAIL, PASSWORD
    }
}
