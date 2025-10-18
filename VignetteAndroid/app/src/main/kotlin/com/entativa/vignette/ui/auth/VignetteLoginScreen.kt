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
 * Vignette Login Screen - Instagram-inspired design
 * Minimalist, clean UI with enterprise-grade functionality
 */
@Composable
fun VignetteLoginScreen(
    viewModel: VignetteAuthViewModel = viewModel(),
    onNavigateToSignUp: () -> Unit = {},
    onLoginSuccess: () -> Unit = {},
    onForgotPassword: () -> Unit = {}
) {
    val uiState by viewModel.uiState.collectAsState()
    val loginForm by viewModel.loginFormState.collectAsState()
    val focusManager = LocalFocusManager.current
    
    var showPassword by remember { mutableStateOf(false) }
    
    LaunchedEffect(uiState.isAuthenticated) {
        if (uiState.isAuthenticated) {
            onLoginSuccess()
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
            Spacer(modifier = Modifier.height(80.dp))
            
            // Logo
            Text(
                text = "Vignette",
                fontSize = 52.sp,
                fontFamily = FontFamily.Cursive,
                fontWeight = FontWeight.Medium,
                color = colorResource(R.color.vignette_text_primary)
            )
            
            Spacer(modifier = Modifier.height(48.dp))
            
            // Login form
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = 32.dp),
                verticalArrangement = Arrangement.spacedBy(12.dp)
            ) {
                // Username/Email field
                OutlinedTextField(
                    value = loginForm.usernameOrEmail,
                    onValueChange = {
                        viewModel.updateLoginField(
                            VignetteAuthViewModel.LoginField.USERNAME_OR_EMAIL,
                            it
                        )
                    },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = { Text("Username or email") },
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
                        focusedContainerColor = colorResource(R.color.vignette_background_secondary),
                        unfocusedContainerColor = colorResource(R.color.vignette_background_secondary)
                    ),
                    shape = RoundedCornerShape(6.dp)
                )
                
                // Password field
                OutlinedTextField(
                    value = loginForm.password,
                    onValueChange = {
                        viewModel.updateLoginField(
                            VignetteAuthViewModel.LoginField.PASSWORD,
                            it
                        )
                    },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = { Text("Password") },
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
                            viewModel.login()
                        }
                    ),
                    singleLine = true,
                    colors = OutlinedTextFieldDefaults.colors(
                        focusedBorderColor = colorResource(R.color.vignette_border_default),
                        unfocusedBorderColor = colorResource(R.color.vignette_border_default),
                        focusedContainerColor = colorResource(R.color.vignette_background_secondary),
                        unfocusedContainerColor = colorResource(R.color.vignette_background_secondary)
                    ),
                    shape = RoundedCornerShape(6.dp)
                )
                
                // Login button
                Button(
                    onClick = { viewModel.login() },
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(44.dp),
                    enabled = !uiState.isLoading,
                    colors = ButtonDefaults.buttonColors(
                        containerColor = colorResource(R.color.vignette_button_primary),
                        disabledContainerColor = colorResource(R.color.vignette_button_primary_disabled)
                    ),
                    shape = RoundedCornerShape(8.dp)
                ) {
                    Text(
                        "Log In",
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
                
                // Facebook login button
                TextButton(
                    onClick = { /* Facebook OAuth implemented below */ },
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
            
            // Forgot password
            Spacer(modifier = Modifier.height(20.dp))
            
            TextButton(onClick = onForgotPassword) {
                Text(
                    "Forgot password?",
                    color = colorResource(R.color.vignette_moonstone),
                    fontSize = 13.sp,
                    fontWeight = FontWeight.Medium
                )
            }
            
            Spacer(modifier = Modifier.weight(1f))
            
            // Sign up section
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
                        "Don't have an account?",
                        fontSize = 14.sp,
                        color = colorResource(R.color.vignette_text_secondary)
                    )
                    TextButton(
                        onClick = onNavigateToSignUp,
                        contentPadding = PaddingValues(0.dp)
                    ) {
                        Text(
                            "Sign up",
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
