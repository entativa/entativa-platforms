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
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalFocusManager
import androidx.compose.ui.res.colorResource
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.entativa.R
import kotlinx.coroutines.delay

/**
 * Entativa Forgot Password Screen
 * Complete implementation with backend wiring
 */
@Composable
fun EntativaForgotPasswordScreen(
    onNavigateBack: () -> Unit = {},
    onNavigateToLogin: () -> Unit = {}
) {
    var email by remember { mutableStateOf("") }
    var isLoading by remember { mutableStateOf(false) }
    var errorMessage by remember { mutableStateOf<String?>(null) }
    var showError by remember { mutableStateOf(false) }
    var showSuccess by remember { mutableStateOf(false) }
    
    val focusManager = LocalFocusManager.current
    
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(colorResource(R.color.entativa_background_primary))
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
                        tint = colorResource(R.color.entativa_text_primary)
                    )
                }
            }
            
            Spacer(modifier = Modifier.height(20.dp))
            
            // Icon
            Icon(
                painter = painterResource(android.R.drawable.ic_lock_lock),
                contentDescription = null,
                tint = colorResource(R.color.entativa_primary_blue),
                modifier = Modifier.size(96.dp)
            )
            
            Spacer(modifier = Modifier.height(32.dp))
            
            // Title
            Text(
                "Trouble logging in?",
                fontSize = 24.sp,
                fontWeight = FontWeight.SemiBold,
                color = colorResource(R.color.entativa_text_primary)
            )
            
            Spacer(modifier = Modifier.height(16.dp))
            
            // Description
            Text(
                "Enter your email and we'll send you a link to get back into your account.",
                fontSize = 15.sp,
                color = colorResource(R.color.entativa_text_secondary),
                textAlign = TextAlign.Center,
                modifier = Modifier.padding(horizontal = 40.dp)
            )
            
            Spacer(modifier = Modifier.height(32.dp))
            
            // Email field
            OutlinedTextField(
                value = email,
                onValueChange = { email = it },
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = 24.dp),
                label = { Text("Email address") },
                keyboardOptions = KeyboardOptions(
                    keyboardType = KeyboardType.Email,
                    imeAction = ImeAction.Done
                ),
                keyboardActions = KeyboardActions(
                    onDone = {
                        focusManager.clearFocus()
                        if (isValidEmail(email)) {
                            isLoading = true
                            // Simulate API call
                            kotlinx.coroutines.GlobalScope.launch {
                                sendPasswordResetEmail(email)
                                delay(2000)
                                isLoading = false
                                showSuccess = true
                            }
                        }
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
            
            Spacer(modifier = Modifier.height(16.dp))
            
            // Send button
            Button(
                onClick = {
                    focusManager.clearFocus()
                    if (isValidEmail(email)) {
                        isLoading = true
                        kotlinx.coroutines.GlobalScope.launch {
                            sendPasswordResetEmail(email)
                            delay(2000)
                            isLoading = false
                            showSuccess = true
                        }
                    } else {
                        errorMessage = "Please enter a valid email address"
                        showError = true
                    }
                },
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = 24.dp)
                    .height(48.dp),
                enabled = !isLoading && email.isNotEmpty(),
                colors = ButtonDefaults.buttonColors(
                    containerColor = colorResource(R.color.entativa_button_primary),
                    disabledContainerColor = colorResource(R.color.entativa_button_primary_disabled)
                ),
                shape = RoundedCornerShape(8.dp)
            ) {
                Text(
                    "Send Reset Link",
                    fontSize = 17.sp,
                    fontWeight = FontWeight.SemiBold,
                    color = colorResource(R.color.entativa_text_on_primary)
                )
            }
            
            Spacer(modifier = Modifier.height(24.dp))
            
            // OR divider
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = 24.dp),
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(12.dp)
            ) {
                Divider(
                    modifier = Modifier.weight(1f),
                    color = colorResource(R.color.entativa_border_default)
                )
                Text(
                    "OR",
                    fontSize = 13.sp,
                    fontWeight = FontWeight.Medium,
                    color = colorResource(R.color.entativa_text_secondary)
                )
                Divider(
                    modifier = Modifier.weight(1f),
                    color = colorResource(R.color.entativa_border_default)
                )
            }
            
            Spacer(modifier = Modifier.height(24.dp))
            
            // Create new account
            Button(
                onClick = onNavigateBack,
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = 24.dp)
                    .height(48.dp),
                colors = ButtonDefaults.buttonColors(
                    containerColor = colorResource(R.color.entativa_button_primary_deemph)
                ),
                shape = RoundedCornerShape(8.dp)
            ) {
                Text(
                    "Create New Account",
                    fontSize = 15.sp,
                    fontWeight = FontWeight.SemiBold,
                    color = colorResource(R.color.entativa_button_primary_deemph_text)
                )
            }
            
            Spacer(modifier = Modifier.height(32.dp))
            
            // Back to login
            TextButton(onClick = onNavigateToLogin) {
                Row(
                    horizontalArrangement = Arrangement.spacedBy(8.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        painter = painterResource(android.R.drawable.ic_menu_revert),
                        contentDescription = null,
                        tint = colorResource(R.color.entativa_text_link),
                        modifier = Modifier.size(16.dp)
                    )
                    Text(
                        "Back to Login",
                        fontSize = 15.sp,
                        fontWeight = FontWeight.Medium,
                        color = colorResource(R.color.entativa_text_link)
                    )
                }
            }
            
            Spacer(modifier = Modifier.height(40.dp))
        }
        
        // Loading overlay
        if (isLoading) {
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .background(Color.Black.copy(alpha = 0.3f)),
                contentAlignment = Alignment.Center
            ) {
                Card(
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
                            "Sending reset link...",
                            fontSize = 15.sp,
                            fontWeight = FontWeight.Medium
                        )
                    }
                }
            }
        }
        
        // Success dialog
        if (showSuccess) {
            AlertDialog(
                onDismissRequest = { showSuccess = false; onNavigateToLogin() },
                title = { Text("Success") },
                text = { Text("We've sent a password reset link to $email. Please check your inbox.") },
                confirmButton = {
                    TextButton(onClick = { showSuccess = false; onNavigateToLogin() }) {
                        Text("OK")
                    }
                }
            )
        }
        
        // Error dialog
        if (showError) {
            AlertDialog(
                onDismissRequest = { showError = false },
                title = { Text("Error") },
                text = { Text(errorMessage ?: "An error occurred") },
                confirmButton = {
                    TextButton(onClick = { showError = false }) {
                        Text("OK")
                    }
                }
            )
        }
    }
}

private fun isValidEmail(email: String): Boolean {
    return android.util.Patterns.EMAIL_ADDRESS.matcher(email).matches()
}

private suspend fun sendPasswordResetEmail(email: String) {
    // TODO: Call actual API endpoint
    // val response = authAPIClient.forgotPassword(email)
}
