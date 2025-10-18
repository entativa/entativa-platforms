package com.entativa.ui.auth

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
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalFocusManager
import androidx.compose.ui.res.colorResource
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.lifecycle.viewmodel.compose.viewModel
import com.entativa.R
import com.entativa.viewmodel.EntativaAuthViewModel

/**
 * Entativa Login Screen - Facebook-inspired design
 * Enterprise-grade implementation with Jetpack Compose
 */
@Composable
fun EntativaLoginScreen(
    viewModel: EntativaAuthViewModel = viewModel(),
    onNavigateToSignUp: () -> Unit = {},
    onLoginSuccess: () -> Unit = {}
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
            .background(colorResource(R.color.entativa_background_primary))
    ) {
        Column(
            modifier = Modifier
                .fillMaxSize()
                .verticalScroll(rememberScrollState())
                .padding(horizontal = 24.dp),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Spacer(modifier = Modifier.height(60.dp))
            
            // Logo and branding
            EntativaLogo()
            
            Spacer(modifier = Modifier.height(40.dp))
            
            // Login form
            OutlinedTextField(
                value = loginForm.emailOrUsername,
                onValueChange = {
                    viewModel.updateLoginField(
                        EntativaAuthViewModel.LoginField.EMAIL_OR_USERNAME,
                        it
                    )
                },
                modifier = Modifier.fillMaxWidth(),
                label = { Text("Email address or username") },
                keyboardOptions = KeyboardOptions(
                    keyboardType = KeyboardType.Email,
                    imeAction = ImeAction.Next
                ),
                keyboardActions = KeyboardActions(
                    onNext = { focusManager.moveFocus(FocusDirection.Down) }
                ),
                singleLine = true,
                colors = OutlinedTextFieldDefaults.colors(
                    focusedBorderColor = colorResource(R.color.entativa_border_focus),
                    unfocusedBorderColor = colorResource(R.color.entativa_border_default),
                    focusedContainerColor = colorResource(R.color.entativa_background_secondary),
                    unfocusedContainerColor = colorResource(R.color.entativa_background_secondary)
                ),
                shape = RoundedCornerShape(8.dp)
            )
            
            Spacer(modifier = Modifier.height(12.dp))
            
            OutlinedTextField(
                value = loginForm.password,
                onValueChange = {
                    viewModel.updateLoginField(
                        EntativaAuthViewModel.LoginField.PASSWORD,
                        it
                    )
                },
                modifier = Modifier.fillMaxWidth(),
                label = { Text("Password") },
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
                            tint = colorResource(R.color.entativa_text_secondary)
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
                    focusedBorderColor = colorResource(R.color.entativa_border_focus),
                    unfocusedBorderColor = colorResource(R.color.entativa_border_default),
                    focusedContainerColor = colorResource(R.color.entativa_background_secondary),
                    unfocusedContainerColor = colorResource(R.color.entativa_background_secondary)
                ),
                shape = RoundedCornerShape(8.dp)
            )
            
            Spacer(modifier = Modifier.height(4.dp))
            
            // Forgot password link
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.End
            ) {
                TextButton(onClick = { /* TODO: Forgot password */ }) {
                    Text(
                        "Forgotten password?",
                        color = colorResource(R.color.entativa_text_link),
                        fontSize = 13.sp,
                        fontWeight = FontWeight.Medium
                    )
                }
            }
            
            Spacer(modifier = Modifier.height(12.dp))
            
            // Login button
            Button(
                onClick = { viewModel.login() },
                modifier = Modifier
                    .fillMaxWidth()
                    .height(48.dp),
                enabled = !uiState.isLoading,
                colors = ButtonDefaults.buttonColors(
                    containerColor = colorResource(R.color.entativa_button_primary),
                    disabledContainerColor = colorResource(R.color.entativa_button_primary_disabled)
                ),
                shape = RoundedCornerShape(8.dp)
            ) {
                Text(
                    "Log In",
                    color = colorResource(R.color.entativa_text_on_primary),
                    fontSize = 17.sp,
                    fontWeight = FontWeight.SemiBold
                )
            }
            
            Spacer(modifier = Modifier.height(28.dp))
            
            // OR divider
            Row(
                modifier = Modifier.fillMaxWidth(),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Divider(
                    modifier = Modifier.weight(1f),
                    color = colorResource(R.color.entativa_border_default)
                )
                Text(
                    "OR",
                    modifier = Modifier.padding(horizontal = 12.dp),
                    color = colorResource(R.color.entativa_text_secondary),
                    fontSize = 13.sp,
                    fontWeight = FontWeight.Medium
                )
                Divider(
                    modifier = Modifier.weight(1f),
                    color = colorResource(R.color.entativa_border_default)
                )
            }
            
            Spacer(modifier = Modifier.height(28.dp))
            
            // Sign up prompt
            Divider(color = colorResource(R.color.entativa_border_default))
            
            Spacer(modifier = Modifier.height(16.dp))
            
            Text(
                "Don't have an account?",
                color = colorResource(R.color.entativa_text_secondary),
                fontSize = 15.sp
            )
            
            Spacer(modifier = Modifier.height(16.dp))
            
            // Sign up button
            Button(
                onClick = onNavigateToSignUp,
                modifier = Modifier
                    .fillMaxWidth()
                    .height(48.dp),
                colors = ButtonDefaults.buttonColors(
                    containerColor = colorResource(R.color.entativa_button_primary_deemph)
                ),
                shape = RoundedCornerShape(8.dp)
            ) {
                Text(
                    "Create New Account",
                    color = colorResource(R.color.entativa_button_primary_deemph_text),
                    fontSize = 15.sp,
                    fontWeight = FontWeight.SemiBold
                )
            }
            
            Spacer(modifier = Modifier.height(40.dp))
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
                    modifier = Modifier.padding(32.dp),
                    shape = RoundedCornerShape(16.dp),
                    colors = CardDefaults.cardColors(
                        containerColor = Color.White.copy(alpha = 0.95f)
                    )
                ) {
                    Column(
                        modifier = Modifier.padding(32.dp),
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        CircularProgressIndicator(
                            color = colorResource(R.color.entativa_button_primary)
                        )
                        Spacer(modifier = Modifier.height(16.dp))
                        Text(
                            "Logging in...",
                            fontSize = 15.sp,
                            fontWeight = FontWeight.Medium
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
private fun EntativaLogo() {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            text = "entativa",
            fontSize = 48.sp,
            fontWeight = FontWeight.Bold,
            fontStyle = FontStyle.Italic,
            style = LocalTextStyle.current.copy(
                brush = Brush.linearGradient(
                    colors = listOf(
                        colorResource(R.color.entativa_primary_blue),
                        colorResource(R.color.entativa_primary_purple),
                        colorResource(R.color.entativa_primary_pink)
                    )
                )
            )
        )
        
        Spacer(modifier = Modifier.height(12.dp))
        
        Text(
            "Connect with friends and the world around you",
            color = colorResource(R.color.entativa_text_secondary),
            fontSize = 15.sp,
            textAlign = TextAlign.Center,
            modifier = Modifier.padding(horizontal = 40.dp)
        )
    }
}
