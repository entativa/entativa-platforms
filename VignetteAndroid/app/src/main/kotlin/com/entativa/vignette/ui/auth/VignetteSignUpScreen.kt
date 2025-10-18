package com.entativa.vignette.ui.auth

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.focus.FocusDirection
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalFocusManager
import androidx.compose.ui.res.colorResource
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardCapitalization
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.lifecycle.viewmodel.compose.viewModel
import com.entativa.vignette.R
import com.entativa.vignette.viewmodel.VignetteAuthViewModel

/**
 * Vignette Sign Up Screen - Instagram-inspired single-page registration
 * Clean, minimalist design with comprehensive validation
 */
@Composable
fun VignetteSignUpScreen(
    viewModel: VignetteAuthViewModel = viewModel(),
    onNavigateBack: () -> Unit = {},
    onSignUpSuccess: () -> Unit = {}
) {
    val uiState by viewModel.uiState.collectAsState()
    val signUpForm by viewModel.signUpFormState.collectAsState()
    val focusManager = LocalFocusManager.current
    
    var showPassword by remember { mutableStateOf(false) }
    
    LaunchedEffect(uiState.isAuthenticated) {
        if (uiState.isAuthenticated) {
            onSignUpSuccess()
        }
    }
    
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(colorResource(R.color.vignette_background_primary))
    ) {
        Column(
            modifier = Modifier
                .fillMaxSize()
                .verticalScroll(rememberScrollState()),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            // Top bar
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                horizontalArrangement = Arrangement.Start
            ) {
                IconButton(onClick = onNavigateBack) {
                    Icon(
                        painter = painterResource(R.drawable.ic_close),
                        contentDescription = "Close",
                        tint = colorResource(R.color.vignette_text_primary)
                    )
                }
            }
            
            Spacer(modifier = Modifier.height(24.dp))
            
            // Logo
            Text(
                text = "Vignette",
                fontSize = 48.sp,
                fontFamily = FontFamily.Cursive,
                fontWeight = FontWeight.Medium,
                color = colorResource(R.color.vignette_text_primary)
            )
            
            Spacer(modifier = Modifier.height(16.dp))
            
            Text(
                "Sign up to see photos and videos\nfrom your friends.",
                fontSize = 14.sp,
                color = colorResource(R.color.vignette_text_secondary),
                textAlign = TextAlign.Center,
                modifier = Modifier.padding(horizontal = 40.dp)
            )
            
            Spacer(modifier = Modifier.height(32.dp))
            
            // Sign up form
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = 32.dp),
                verticalArrangement = Arrangement.spacedBy(12.dp)
            ) {
                // Email field
                OutlinedTextField(
                    value = signUpForm.email,
                    onValueChange = {
                        viewModel.updateSignUpField(
                            VignetteAuthViewModel.SignUpField.EMAIL,
                            it
                        )
                    },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = { Text("Email") },
                    isError = signUpForm.emailError != null,
                    supportingText = signUpForm.emailError?.let { { Text(it) } },
                    keyboardOptions = KeyboardOptions(
                        keyboardType = KeyboardType.Email,
                        imeAction = ImeAction.Next
                    ),
                    keyboardActions = KeyboardActions(
                        onNext = { focusManager.moveFocus(FocusDirection.Down) }
                    ),
                    singleLine = true,
                    colors = OutlinedTextFieldDefaults.colors(
                        focusedBorderColor = colorResource(R.color.vignette_border_default),
                        unfocusedBorderColor = colorResource(R.color.vignette_border_default),
                        errorBorderColor = colorResource(R.color.vignette_border_error),
                        focusedContainerColor = colorResource(R.color.vignette_background_secondary),
                        unfocusedContainerColor = colorResource(R.color.vignette_background_secondary)
                    ),
                    shape = RoundedCornerShape(6.dp)
                )
                
                // Full name field
                OutlinedTextField(
                    value = signUpForm.fullName,
                    onValueChange = {
                        viewModel.updateSignUpField(
                            VignetteAuthViewModel.SignUpField.FULL_NAME,
                            it
                        )
                    },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = { Text("Full Name") },
                    isError = signUpForm.fullNameError != null,
                    supportingText = signUpForm.fullNameError?.let { { Text(it) } },
                    keyboardOptions = KeyboardOptions(
                        capitalization = KeyboardCapitalization.Words,
                        imeAction = ImeAction.Next
                    ),
                    keyboardActions = KeyboardActions(
                        onNext = { focusManager.moveFocus(FocusDirection.Down) }
                    ),
                    singleLine = true,
                    colors = OutlinedTextFieldDefaults.colors(
                        focusedBorderColor = colorResource(R.color.vignette_border_default),
                        unfocusedBorderColor = colorResource(R.color.vignette_border_default),
                        errorBorderColor = colorResource(R.color.vignette_border_error),
                        focusedContainerColor = colorResource(R.color.vignette_background_secondary),
                        unfocusedContainerColor = colorResource(R.color.vignette_background_secondary)
                    ),
                    shape = RoundedCornerShape(6.dp)
                )
                
                // Username field
                OutlinedTextField(
                    value = signUpForm.username,
                    onValueChange = {
                        viewModel.updateSignUpField(
                            VignetteAuthViewModel.SignUpField.USERNAME,
                            it.lowercase()
                        )
                    },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = { Text("Username") },
                    isError = signUpForm.usernameError != null,
                    supportingText = signUpForm.usernameError?.let { { Text(it) } },
                    keyboardOptions = KeyboardOptions(
                        keyboardType = KeyboardType.Ascii,
                        imeAction = ImeAction.Next
                    ),
                    keyboardActions = KeyboardActions(
                        onNext = { focusManager.moveFocus(FocusDirection.Down) }
                    ),
                    singleLine = true,
                    colors = OutlinedTextFieldDefaults.colors(
                        focusedBorderColor = colorResource(R.color.vignette_border_default),
                        unfocusedBorderColor = colorResource(R.color.vignette_border_default),
                        errorBorderColor = colorResource(R.color.vignette_border_error),
                        focusedContainerColor = colorResource(R.color.vignette_background_secondary),
                        unfocusedContainerColor = colorResource(R.color.vignette_background_secondary)
                    ),
                    shape = RoundedCornerShape(6.dp)
                )
                
                // Username hint
                if (signUpForm.usernameError == null && signUpForm.username.isNotEmpty()) {
                    Text(
                        "Can contain letters, numbers, periods, and underscores",
                        fontSize = 11.sp,
                        color = colorResource(R.color.vignette_text_secondary),
                        modifier = Modifier.padding(horizontal = 4.dp)
                    )
                }
                
                // Password field
                OutlinedTextField(
                    value = signUpForm.password,
                    onValueChange = {
                        viewModel.updateSignUpField(
                            VignetteAuthViewModel.SignUpField.PASSWORD,
                            it
                        )
                    },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = { Text("Password") },
                    isError = signUpForm.passwordError != null,
                    supportingText = signUpForm.passwordError?.let { { Text(it) } },
                    visualTransformation = if (showPassword) VisualTransformation.None
                    else PasswordVisualTransformation(),
                    trailingIcon = {
                        IconButton(onClick = { showPassword = !showPassword }) {
                            Icon(
                                painter = painterResource(
                                    if (showPassword) R.drawable.ic_eye_slash
                                    else R.drawable.ic_eye
                                ),
                                contentDescription = if (showPassword) "Hide password" else "Show password",
                                tint = colorResource(R.color.vignette_text_secondary)
                            )
                        }
                    },
                    keyboardOptions = KeyboardOptions(
                        keyboardType = KeyboardType.Password,
                        imeAction = ImeAction.Done
                    ),
                    keyboardActions = KeyboardActions(
                        onDone = {
                            focusManager.clearFocus()
                            if (canSignUp(signUpForm)) {
                                viewModel.signUp()
                            }
                        }
                    ),
                    singleLine = true,
                    colors = OutlinedTextFieldDefaults.colors(
                        focusedBorderColor = colorResource(R.color.vignette_border_default),
                        unfocusedBorderColor = colorResource(R.color.vignette_border_default),
                        errorBorderColor = colorResource(R.color.vignette_border_error),
                        focusedContainerColor = colorResource(R.color.vignette_background_secondary),
                        unfocusedContainerColor = colorResource(R.color.vignette_background_secondary)
                    ),
                    shape = RoundedCornerShape(6.dp)
                )
                
                // Password requirements
                if (signUpForm.password.isNotEmpty()) {
                    Column(
                        verticalArrangement = Arrangement.spacedBy(4.dp),
                        modifier = Modifier.padding(horizontal = 4.dp)
                    ) {
                        PasswordRequirement(
                            text = "8+ characters",
                            isMet = signUpForm.password.length >= 8
                        )
                        PasswordRequirement(
                            text = "Uppercase & lowercase",
                            isMet = signUpForm.password.any { it.isUpperCase() } &&
                                    signUpForm.password.any { it.isLowerCase() }
                        )
                        PasswordRequirement(
                            text = "Contains a number",
                            isMet = signUpForm.password.any { it.isDigit() }
                        )
                    }
                }
                
                // Sign up button
                Button(
                    onClick = { viewModel.signUp() },
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(44.dp),
                    enabled = !uiState.isLoading && canSignUp(signUpForm),
                    colors = ButtonDefaults.buttonColors(
                        containerColor = colorResource(R.color.vignette_button_primary),
                        disabledContainerColor = colorResource(R.color.vignette_button_primary_disabled)
                    ),
                    shape = RoundedCornerShape(8.dp)
                ) {
                    Text(
                        "Sign Up",
                        color = colorResource(R.color.vignette_text_on_primary),
                        fontSize = 16.sp,
                        fontWeight = FontWeight.SemiBold
                    )
                }
                
                // OR divider
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(vertical = 24.dp),
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(16.dp)
                ) {
                    Divider(
                        modifier = Modifier.weight(1f),
                        color = colorResource(R.color.vignette_border_default)
                    )
                    Text(
                        "OR",
                        fontSize = 13.sp,
                        fontWeight = FontWeight.SemiBold,
                        color = colorResource(R.color.vignette_text_secondary)
                    )
                    Divider(
                        modifier = Modifier.weight(1f),
                        color = colorResource(R.color.vignette_border_default)
                    )
                }
                
                // Facebook sign up
                TextButton(
                    onClick = { /* Facebook OAuth */ },
                    modifier = Modifier.fillMaxWidth()
                ) {
                    Row(
                        horizontalArrangement = Arrangement.spacedBy(8.dp),
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Icon(
                            painter = painterResource(R.drawable.ic_facebook),
                            contentDescription = "Facebook",
                            tint = colorResource(R.color.vignette_moonstone),
                            modifier = Modifier.size(20.dp)
                        )
                        Text(
                            "Log in with Facebook",
                            color = colorResource(R.color.vignette_moonstone),
                            fontSize = 14.sp,
                            fontWeight = FontWeight.SemiBold
                        )
                    }
                }
            }
            
            Spacer(modifier = Modifier.height(32.dp))
            
            // Terms
            Text(
                buildString {
                    append("By signing up, you agree to our ")
                    append("Terms, Privacy Policy ")
                    append("and Cookies Policy.")
                },
                fontSize = 12.sp,
                color = colorResource(R.color.vignette_text_secondary),
                textAlign = TextAlign.Center,
                modifier = Modifier.padding(horizontal = 40.dp)
            )
            
            Spacer(modifier = Modifier.weight(1f))
            
            // Login link
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(16.dp)
            ) {
                Divider(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 32.dp),
                    color = colorResource(R.color.vignette_border_default)
                )
                
                Row(
                    horizontalArrangement = Arrangement.spacedBy(4.dp)
                ) {
                    Text(
                        "Have an account?",
                        fontSize = 14.sp,
                        color = colorResource(R.color.vignette_text_secondary)
                    )
                    TextButton(
                        onClick = onNavigateBack,
                        contentPadding = PaddingValues(0.dp)
                    ) {
                        Text(
                            "Log in",
                            fontSize = 14.sp,
                            fontWeight = FontWeight.SemiBold,
                            color = colorResource(R.color.vignette_moonstone)
                        )
                    }
                }
            }
            
            Spacer(modifier = Modifier.height(20.dp))
        }
        
        // Loading overlay
        if (uiState.isLoading) {
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .background(Color.Black.copy(alpha = 0.3f)),
                contentAlignment = Alignment.Center
            ) {
                Card(
                    shape = RoundedCornerShape(12.dp),
                    colors = CardDefaults.cardColors(
                        containerColor = Color.White
                    ),
                    elevation = CardDefaults.cardElevation(defaultElevation = 10.dp)
                ) {
                    Box(
                        modifier = Modifier.padding(40.dp),
                        contentAlignment = Alignment.Center
                    ) {
                        CircularProgressIndicator(
                            color = colorResource(R.color.vignette_button_primary),
                            modifier = Modifier.size(40.dp)
                        )
                    }
                }
            }
        }
        
        // Error dialog
        if (uiState.showError) {
            AlertDialog(
                onDismissRequest = { viewModel.clearError() },
                title = { Text("Error") },
                text = { Text(uiState.errorMessage ?: "An error occurred") },
                confirmButton = {
                    TextButton(onClick = { viewModel.clearError() }) {
                        Text("OK")
                    }
                }
            )
        }
    }
}

@Composable
private fun PasswordRequirement(text: String, isMet: Boolean) {
    Row(
        horizontalArrangement = Arrangement.spacedBy(6.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Icon(
            painter = painterResource(
                if (isMet) R.drawable.ic_check_circle_filled
                else R.drawable.ic_circle
            ),
            contentDescription = null,
            tint = if (isMet) {
                colorResource(R.color.vignette_success)
            } else {
                colorResource(R.color.vignette_text_tertiary)
            },
            modifier = Modifier.size(12.dp)
        )
        
        Text(
            text,
            fontSize = 11.sp,
            color = if (isMet) {
                colorResource(R.color.vignette_text_primary)
            } else {
                colorResource(R.color.vignette_text_secondary)
            }
        )
    }
}

private fun canSignUp(form: VignetteAuthViewModel.SignUpFormState): Boolean {
    return form.email.trim().isNotEmpty() &&
            form.fullName.trim().isNotEmpty() &&
            form.username.trim().isNotEmpty() &&
            form.password.isNotEmpty() &&
            form.password.length >= 8
}
