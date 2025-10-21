package com.entativa.ui.auth

import androidx.compose.animation.AnimatedContent
import androidx.compose.animation.slideInHorizontally
import androidx.compose.animation.slideOutHorizontally
import androidx.compose.animation.togetherWith
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
import com.entativa.R
import com.entativa.viewmodel.EntativaAuthViewModel
import java.util.Calendar

/**
 * Entativa Sign Up Screen - Multi-step Facebook-inspired registration
 * Enterprise-grade implementation with Jetpack Compose
 */
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun EntativaSignUpScreen(
    viewModel: EntativaAuthViewModel = viewModel(),
    onNavigateBack: () -> Unit = {},
    onSignUpSuccess: () -> Unit = {}
) {
    val uiState by viewModel.uiState.collectAsState()
    val signUpForm by viewModel.signUpFormState.collectAsState()
    val focusManager = LocalFocusManager.current
    
    var currentStep by remember { mutableStateOf(1) }
    var showPassword by remember { mutableStateOf(false) }
    var showDatePicker by remember { mutableStateOf(false) }
    
    LaunchedEffect(uiState.isAuthenticated) {
        if (uiState.isAuthenticated) {
            onSignUpSuccess()
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
                .verticalScroll(rememberScrollState()),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            // Top bar with close button
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                IconButton(onClick = onNavigateBack) {
                    Icon(
                        painter = painterResource(R.drawable.ic_close),
                        contentDescription = "Close",
                        tint = colorResource(R.color.entativa_text_primary)
                    )
                }
            }
            
            Spacer(modifier = Modifier.height(8.dp))
            
            // Header
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                modifier = Modifier.padding(horizontal = 24.dp)
            ) {
                Text(
                    "Create Account",
                    fontSize = 28.sp,
                    fontWeight = FontWeight.SemiBold,
                    color = colorResource(R.color.entativa_text_primary)
                )
                
                Spacer(modifier = Modifier.height(8.dp))
                
                Text(
                    when (currentStep) {
                        1 -> "What's your name?"
                        2 -> "Enter your email and password"
                        3 -> "Tell us about yourself"
                        else -> ""
                    },
                    fontSize = 15.sp,
                    color = colorResource(R.color.entativa_text_secondary),
                    textAlign = TextAlign.Center
                )
            }
            
            Spacer(modifier = Modifier.height(24.dp))
            
            // Progress indicator
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = 24.dp),
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                repeat(3) { step ->
                    Box(
                        modifier = Modifier
                            .weight(1f)
                            .height(4.dp)
                            .background(
                                color = if (step < currentStep) {
                                    colorResource(R.color.entativa_button_primary)
                                } else {
                                    colorResource(R.color.entativa_border_default)
                                },
                                shape = RoundedCornerShape(2.dp)
                            )
                    )
                }
            }
            
            Spacer(modifier = Modifier.height(24.dp))
            
            // Form content with animation
            AnimatedContent(
                targetState = currentStep,
                transitionSpec = {
                    slideInHorizontally { it } togetherWith slideOutHorizontally { -it }
                },
                label = "step_transition"
            ) { step ->
                Column(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 24.dp)
                ) {
                    when (step) {
                        1 -> NameStep(
                            viewModel = viewModel,
                            signUpForm = signUpForm,
                            focusManager = focusManager
                        )
                        2 -> EmailPasswordStep(
                            viewModel = viewModel,
                            signUpForm = signUpForm,
                            showPassword = showPassword,
                            onTogglePassword = { showPassword = !showPassword },
                            focusManager = focusManager
                        )
                        3 -> BirthdayGenderStep(
                            viewModel = viewModel,
                            signUpForm = signUpForm,
                            onShowDatePicker = { showDatePicker = true }
                        )
                    }
                }
            }
            
            Spacer(modifier = Modifier.height(24.dp))
            
            // Navigation buttons
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = 24.dp),
                horizontalArrangement = Arrangement.spacedBy(12.dp)
            ) {
                if (currentStep > 1) {
                    Button(
                        onClick = { currentStep-- },
                        modifier = Modifier
                            .weight(1f)
                            .height(48.dp),
                        colors = ButtonDefaults.buttonColors(
                            containerColor = colorResource(R.color.entativa_button_secondary)
                        ),
                        shape = RoundedCornerShape(8.dp)
                    ) {
                        Text(
                            "Back",
                            color = colorResource(R.color.entativa_button_secondary_text),
                            fontSize = 15.sp,
                            fontWeight = FontWeight.SemiBold
                        )
                    }
                }
                
                Button(
                    onClick = {
                        if (currentStep < 3) {
                            currentStep++
                        } else {
                            viewModel.signUp()
                        }
                    },
                    modifier = Modifier
                        .weight(1f)
                        .height(48.dp),
                    enabled = !uiState.isLoading && canProceed(currentStep, signUpForm),
                    colors = ButtonDefaults.buttonColors(
                        containerColor = colorResource(R.color.entativa_button_primary),
                        disabledContainerColor = colorResource(R.color.entativa_button_primary_disabled)
                    ),
                    shape = RoundedCornerShape(8.dp)
                ) {
                    Text(
                        if (currentStep == 3) "Sign Up" else "Next",
                        color = colorResource(R.color.entativa_text_on_primary),
                        fontSize = 17.sp,
                        fontWeight = FontWeight.SemiBold
                    )
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
                color = colorResource(R.color.entativa_text_secondary),
                textAlign = TextAlign.Center,
                modifier = Modifier.padding(horizontal = 40.dp)
            )
            
            Spacer(modifier = Modifier.height(40.dp))
        }
        
        // Date picker dialog
        if (showDatePicker) {
            DatePickerModal(
                onDismiss = { showDatePicker = false },
                onDateSelected = { calendar ->
                    viewModel.updateSignUpField(
                        EntativaAuthViewModel.SignUpField.BIRTHDAY,
                        calendar
                    )
                    showDatePicker = false
                },
                initialDate = signUpForm.birthday
            )
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
                            "Creating your account...",
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
private fun NameStep(
    viewModel: EntativaAuthViewModel,
    signUpForm: EntativaAuthViewModel.SignUpFormState,
    focusManager: androidx.compose.ui.focus.FocusManager
) {
    Column(verticalArrangement = Arrangement.spacedBy(16.dp)) {
        OutlinedTextField(
            value = signUpForm.firstName,
            onValueChange = {
                viewModel.updateSignUpField(
                    EntativaAuthViewModel.SignUpField.FIRST_NAME,
                    it
                )
            },
            modifier = Modifier.fillMaxWidth(),
            label = { Text("First name") },
            isError = signUpForm.firstNameError != null,
            supportingText = signUpForm.firstNameError?.let { { Text(it) } },
            keyboardOptions = KeyboardOptions(
                capitalization = KeyboardCapitalization.Words,
                imeAction = ImeAction.Next
            ),
            keyboardActions = KeyboardActions(
                onNext = { focusManager.moveFocus(FocusDirection.Down) }
            ),
            singleLine = true,
            colors = OutlinedTextFieldDefaults.colors(
                focusedBorderColor = colorResource(R.color.entativa_border_focus),
                unfocusedBorderColor = colorResource(R.color.entativa_border_default),
                errorBorderColor = colorResource(R.color.entativa_border_error),
                focusedContainerColor = colorResource(R.color.entativa_background_secondary),
                unfocusedContainerColor = colorResource(R.color.entativa_background_secondary)
            ),
            shape = RoundedCornerShape(8.dp)
        )
        
        OutlinedTextField(
            value = signUpForm.lastName,
            onValueChange = {
                viewModel.updateSignUpField(
                    EntativaAuthViewModel.SignUpField.LAST_NAME,
                    it
                )
            },
            modifier = Modifier.fillMaxWidth(),
            label = { Text("Last name") },
            isError = signUpForm.lastNameError != null,
            supportingText = signUpForm.lastNameError?.let { { Text(it) } },
            keyboardOptions = KeyboardOptions(
                capitalization = KeyboardCapitalization.Words,
                imeAction = ImeAction.Done
            ),
            keyboardActions = KeyboardActions(
                onDone = { focusManager.clearFocus() }
            ),
            singleLine = true,
            colors = OutlinedTextFieldDefaults.colors(
                focusedBorderColor = colorResource(R.color.entativa_border_focus),
                unfocusedBorderColor = colorResource(R.color.entativa_border_default),
                errorBorderColor = colorResource(R.color.entativa_border_error),
                focusedContainerColor = colorResource(R.color.entativa_background_secondary),
                unfocusedContainerColor = colorResource(R.color.entativa_background_secondary)
            ),
            shape = RoundedCornerShape(8.dp)
        )
    }
}

@Composable
private fun EmailPasswordStep(
    viewModel: EntativaAuthViewModel,
    signUpForm: EntativaAuthViewModel.SignUpFormState,
    showPassword: Boolean,
    onTogglePassword: () -> Unit,
    focusManager: androidx.compose.ui.focus.FocusManager
) {
    Column(verticalArrangement = Arrangement.spacedBy(16.dp)) {
        OutlinedTextField(
            value = signUpForm.email,
            onValueChange = {
                viewModel.updateSignUpField(
                    EntativaAuthViewModel.SignUpField.EMAIL,
                    it
                )
            },
            modifier = Modifier.fillMaxWidth(),
            label = { Text("Email address") },
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
                focusedBorderColor = colorResource(R.color.entativa_border_focus),
                unfocusedBorderColor = colorResource(R.color.entativa_border_default),
                errorBorderColor = colorResource(R.color.entativa_border_error),
                focusedContainerColor = colorResource(R.color.entativa_background_secondary),
                unfocusedContainerColor = colorResource(R.color.entativa_background_secondary)
            ),
            shape = RoundedCornerShape(8.dp)
        )
        
        OutlinedTextField(
            value = signUpForm.password,
            onValueChange = {
                viewModel.updateSignUpField(
                    EntativaAuthViewModel.SignUpField.PASSWORD,
                    it
                )
            },
            modifier = Modifier.fillMaxWidth(),
            label = { Text("Password") },
            isError = signUpForm.passwordError != null,
            supportingText = signUpForm.passwordError?.let { { Text(it) } },
            visualTransformation = if (showPassword) VisualTransformation.None
            else PasswordVisualTransformation(),
            trailingIcon = {
                IconButton(onClick = onTogglePassword) {
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
                onDone = { focusManager.clearFocus() }
            ),
            singleLine = true,
            colors = OutlinedTextFieldDefaults.colors(
                focusedBorderColor = colorResource(R.color.entativa_border_focus),
                unfocusedBorderColor = colorResource(R.color.entativa_border_default),
                errorBorderColor = colorResource(R.color.entativa_border_error),
                focusedContainerColor = colorResource(R.color.entativa_background_secondary),
                unfocusedContainerColor = colorResource(R.color.entativa_background_secondary)
            ),
            shape = RoundedCornerShape(8.dp)
        )
        
        // Password requirements
        if (signUpForm.password.isNotEmpty()) {
            Column(
                verticalArrangement = Arrangement.spacedBy(6.dp),
                modifier = Modifier.padding(horizontal = 4.dp)
            ) {
                PasswordRequirement(
                    text = "At least 8 characters",
                    isMet = signUpForm.password.length >= 8
                )
                PasswordRequirement(
                    text = "Contains uppercase letter",
                    isMet = signUpForm.password.any { it.isUpperCase() }
                )
                PasswordRequirement(
                    text = "Contains lowercase letter",
                    isMet = signUpForm.password.any { it.isLowerCase() }
                )
                PasswordRequirement(
                    text = "Contains number",
                    isMet = signUpForm.password.any { it.isDigit() }
                )
            }
        }
    }
}

@Composable
private fun BirthdayGenderStep(
    viewModel: EntativaAuthViewModel,
    signUpForm: EntativaAuthViewModel.SignUpFormState,
    onShowDatePicker: () -> Unit
) {
    Column(verticalArrangement = Arrangement.spacedBy(24.dp)) {
        // Birthday
        Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
            Text(
                "Birthday",
                fontSize = 15.sp,
                fontWeight = FontWeight.Medium,
                color = colorResource(R.color.entativa_text_primary)
            )
            
            OutlinedButton(
                onClick = onShowDatePicker,
                modifier = Modifier
                    .fillMaxWidth()
                    .height(56.dp),
                shape = RoundedCornerShape(8.dp),
                colors = ButtonDefaults.outlinedButtonColors(
                    containerColor = colorResource(R.color.entativa_background_secondary)
                ),
                border = ButtonDefaults.outlinedButtonBorder.copy(
                    brush = androidx.compose.ui.graphics.SolidColor(
                        if (signUpForm.birthdayError != null) {
                            colorResource(R.color.entativa_border_error)
                        } else {
                            colorResource(R.color.entativa_border_default)
                        }
                    )
                )
            ) {
                Text(
                    formatDate(signUpForm.birthday),
                    fontSize = 16.sp,
                    color = colorResource(R.color.entativa_text_primary)
                )
            }
            
            signUpForm.birthdayError?.let { error ->
                Text(
                    error,
                    fontSize = 12.sp,
                    color = colorResource(R.color.entativa_error),
                    modifier = Modifier.padding(horizontal = 16.dp)
                )
            }
            
            Text(
                "You must be at least 13 years old",
                fontSize = 12.sp,
                color = colorResource(R.color.entativa_text_secondary),
                modifier = Modifier.padding(horizontal = 4.dp)
            )
        }
        
        // Gender
        Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
            Text(
                "Gender",
                fontSize = 15.sp,
                fontWeight = FontWeight.Medium,
                color = colorResource(R.color.entativa_text_primary)
            )
            
            Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                viewModel.genderOptions.chunked(2).forEach { row ->
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.spacedBy(8.dp)
                    ) {
                        row.forEach { (value, label) ->
                            FilterChip(
                                selected = signUpForm.gender == value,
                                onClick = {
                                    viewModel.updateSignUpField(
                                        EntativaAuthViewModel.SignUpField.GENDER,
                                        value
                                    )
                                },
                                label = { Text(label) },
                                modifier = Modifier.weight(1f),
                                colors = FilterChipDefaults.filterChipColors(
                                    selectedContainerColor = colorResource(R.color.entativa_button_primary),
                                    selectedLabelColor = colorResource(R.color.entativa_text_on_primary)
                                )
                            )
                        }
                        // Fill empty space if odd number
                        if (row.size < 2) {
                            Spacer(modifier = Modifier.weight(1f))
                        }
                    }
                }
            }
            
            Text(
                "You can always change this later",
                fontSize = 12.sp,
                color = colorResource(R.color.entativa_text_secondary),
                modifier = Modifier.padding(horizontal = 4.dp)
            )
        }
    }
}

@Composable
private fun PasswordRequirement(text: String, isMet: Boolean) {
    Row(
        horizontalArrangement = Arrangement.spacedBy(8.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Icon(
            painter = painterResource(
                if (isMet) R.drawable.ic_check_circle_filled
                else R.drawable.ic_circle
            ),
            contentDescription = null,
            tint = if (isMet) {
                colorResource(R.color.entativa_success)
            } else {
                colorResource(R.color.entativa_text_secondary)
            },
            modifier = Modifier.size(14.dp)
        )
        
        Text(
            text,
            fontSize = 12.sp,
            color = if (isMet) {
                colorResource(R.color.entativa_text_primary)
            } else {
                colorResource(R.color.entativa_text_secondary)
            }
        )
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun DatePickerModal(
    onDismiss: () -> Unit,
    onDateSelected: (Calendar) -> Unit,
    initialDate: Calendar
) {
    val datePickerState = rememberDatePickerState(
        initialSelectedDateMillis = initialDate.timeInMillis
    )
    
    DatePickerDialog(
        onDismissRequest = onDismiss,
        confirmButton = {
            TextButton(
                onClick = {
                    datePickerState.selectedDateMillis?.let { millis ->
                        val calendar = Calendar.getInstance().apply {
                            timeInMillis = millis
                        }
                        onDateSelected(calendar)
                    }
                }
            ) {
                Text("OK")
            }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) {
                Text("Cancel")
            }
        }
    ) {
        DatePicker(state = datePickerState)
    }
}

private fun formatDate(calendar: Calendar): String {
    val month = calendar.get(Calendar.MONTH) + 1
    val day = calendar.get(Calendar.DAY_OF_MONTH)
    val year = calendar.get(Calendar.YEAR)
    return "$month/$day/$year"
}

private fun canProceed(
    step: Int,
    form: EntativaAuthViewModel.SignUpFormState
): Boolean {
    return when (step) {
        1 -> form.firstName.trim().isNotEmpty() && form.lastName.trim().isNotEmpty()
        2 -> form.email.trim().isNotEmpty() && form.password.isNotEmpty()
        3 -> true
        else -> false
    }
}
