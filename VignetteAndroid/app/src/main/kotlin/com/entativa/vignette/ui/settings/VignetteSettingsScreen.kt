package com.entativa.vignette.ui.settings

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewmodel.compose.viewModel
import com.entativa.vignette.R
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow

// MARK: - Vignette Settings Screen (Instagram-Style)
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VignetteSettingsScreen(
    onDismiss: () -> Unit,
    viewModel: VignetteSettingsViewModel = viewModel()
) {
    var showLogoutDialog by remember { mutableStateOf(false) }
    val darkMode by viewModel.darkModeEnabled.collectAsState()

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Settings") },
                navigationIcon = {
                    TextButton(onClick = onDismiss) {
                        Text("Done", color = Color(0xFF007CFC))
                    }
                }
            )
        }
    ) { paddingValues ->
        LazyColumn(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            // Account Section
            item {
                SettingsSectionHeader("Account")
            }

            item {
                SettingsItem("Edit Profile", R.drawable.ic_person) {}
                SettingsItem("Change Password", R.drawable.ic_key) {}
                SettingsItem("Account Privacy", R.drawable.ic_lock) {}
                Divider()
            }

            // Content & Activity
            item {
                SettingsSectionHeader("Content & Activity")
            }

            item {
                SettingsItem("Notifications", R.drawable.ic_bell) {}
                SettingsItem("Posts You've Liked", R.drawable.ic_heart) {}
                SettingsItem("Saved", R.drawable.ic_bookmark) {}
                SettingsItem("Archive", R.drawable.ic_archive) {}
                Divider()
            }

            // Security
            item {
                SettingsSectionHeader("Security")
            }

            item {
                SettingsItem("Security", R.drawable.ic_shield) {}
                SettingsItem("Two-Factor Authentication", R.drawable.ic_lock_shield) {}
                SettingsItem("Login Activity", R.drawable.ic_clock_history) {}
                Divider()
            }

            // Privacy
            item {
                SettingsSectionHeader("Privacy")
            }

            item {
                SettingsItem("Privacy Controls", R.drawable.ic_privacy) {}
                SettingsItem("Blocked Accounts", R.drawable.ic_block) {}
                SettingsItem("Muted Accounts", R.drawable.ic_volume_off) {}
                SettingsItem("Restricted Accounts", R.drawable.ic_restricted) {}
                Divider()
            }

            // Preferences
            item {
                SettingsSectionHeader("Preferences")
            }

            item {
                LanguageItem()
                
                // Dark Mode Toggle
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 16.dp, vertical = 12.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        painter = painterResource(R.drawable.ic_moon),
                        contentDescription = null,
                        modifier = Modifier.size(24.dp)
                    )
                    Spacer(modifier = Modifier.width(12.dp))
                    Text("Dark Mode", fontSize = 16.sp, modifier = Modifier.weight(1f))
                    Switch(
                        checked = darkMode,
                        onCheckedChange = { viewModel.toggleDarkMode() }
                    )
                }
                
                SettingsItem("Data Usage", R.drawable.ic_data) {}
                Divider()
            }

            // Help & Support
            item {
                SettingsSectionHeader("Help & Support")
            }

            item {
                SettingsItem("Help", R.drawable.ic_help) {}
                SettingsItem("About", R.drawable.ic_info) {}
                SettingsItem("Account Status", R.drawable.ic_check_shield) {}
                Divider()
            }

            // Switch Account
            item {
                Divider()
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .clickable {}
                        .padding(horizontal = 16.dp, vertical = 12.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        painter = painterResource(R.drawable.ic_switch),
                        contentDescription = null,
                        tint = Color(0xFF007CFC),
                        modifier = Modifier.size(24.dp)
                    )
                    Spacer(modifier = Modifier.width(12.dp))
                    Text("Switch to Entativa", fontSize = 16.sp)
                }
                Divider()
            }

            // Logout
            item {
                Button(
                    onClick = { showLogoutDialog = true },
                    colors = ButtonDefaults.buttonColors(
                        containerColor = Color.Transparent,
                        contentColor = Color.Red
                    ),
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 16.dp, vertical = 12.dp)
                ) {
                    Text("Log Out", fontWeight = FontWeight.SemiBold)
                }
                Divider()
            }

            // Legal
            item {
                SettingsSectionHeader("Legal")
            }

            item {
                SettingsLinkItem("Terms of Service") {}
                SettingsLinkItem("Privacy Policy") {}
                SettingsLinkItem("Community Guidelines") {}
                Divider()
            }

            // App Info
            item {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 16.dp, vertical = 12.dp),
                    horizontalArrangement = Arrangement.SpaceBetween
                ) {
                    Text("Version", fontSize = 14.sp, color = Color.Gray)
                    Text("1.0.0", fontSize = 14.sp, color = Color.Gray)
                }
            }

            // Footer
            item {
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(16.dp),
                    contentAlignment = Alignment.Center
                ) {
                    Text(
                        "© 2025 Vignette, Inc.",
                        fontSize = 12.sp,
                        color = Color.Gray
                    )
                }
            }
        }
    }

    // Logout Dialog
    if (showLogoutDialog) {
        AlertDialog(
            onDismissRequest = { showLogoutDialog = false },
            title = { Text("Log Out") },
            text = { Text("Are you sure you want to log out?") },
            confirmButton = {
                TextButton(onClick = {
                    viewModel.logout()
                    showLogoutDialog = false
                    onDismiss()
                }) {
                    Text("Log Out", color = Color.Red)
                }
            },
            dismissButton = {
                TextButton(onClick = { showLogoutDialog = false }) {
                    Text("Cancel")
                }
            }
        )
    }
}

// MARK: - Helper Composables
@Composable
fun SettingsSectionHeader(title: String) {
    Text(
        text = title,
        fontSize = 14.sp,
        fontWeight = FontWeight.Bold,
        color = Color.Gray,
        modifier = Modifier.padding(horizontal = 16.dp, vertical = 12.dp)
    )
}

@Composable
fun SettingsItem(title: String, icon: Int, onClick: () -> Unit) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(horizontal = 16.dp, vertical = 12.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Icon(
            painter = painterResource(icon),
            contentDescription = null,
            modifier = Modifier.size(24.dp)
        )
        Spacer(modifier = Modifier.width(12.dp))
        Text(title, fontSize = 16.sp, modifier = Modifier.weight(1f))
        Icon(
            painter = painterResource(R.drawable.ic_chevron_right),
            contentDescription = null,
            tint = Color.Gray
        )
    }
}

@Composable
fun SettingsLinkItem(title: String, onClick: () -> Unit) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(horizontal = 16.dp, vertical = 12.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Text(
            title,
            fontSize = 16.sp,
            color = Color(0xFF007CFC),
            modifier = Modifier.weight(1f)
        )
    }
}

@Composable
fun LanguageItem() {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable {}
            .padding(horizontal = 16.dp, vertical = 12.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Icon(
            painter = painterResource(R.drawable.ic_globe),
            contentDescription = null,
            modifier = Modifier.size(24.dp)
        )
        Spacer(modifier = Modifier.width(12.dp))
        Text("Language", fontSize = 16.sp, modifier = Modifier.weight(1f))
        Text("English", fontSize = 16.sp, color = Color.Gray)
        Spacer(modifier = Modifier.width(8.dp))
        Icon(
            painter = painterResource(R.drawable.ic_chevron_right),
            contentDescription = null,
            tint = Color.Gray
        )
    }
}

// MARK: - Edit Profile Screen
@Composable
fun EditProfileScreen(onDismiss: () -> Unit) {
    var name by remember { mutableStateOf("") }
    var username by remember { mutableStateOf("") }
    var bio by remember { mutableStateOf("") }
    var website by remember { mutableStateOf("") }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Edit Profile") },
                navigationIcon = {
                    IconButton(onClick = onDismiss) {
                        Icon(painterResource(R.drawable.ic_arrow_back), null)
                    }
                },
                actions = {
                    TextButton(onClick = { /* Save */ }) {
                        Text("Save", color = Color(0xFF007CFC), fontWeight = FontWeight.SemiBold)
                    }
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(16.dp)
        ) {
            // Profile Photo
            Box(
                modifier = Modifier.fillMaxWidth(),
                contentAlignment = Alignment.Center
            ) {
                Column(horizontalAlignment = Alignment.CenterHorizontally) {
                    Surface(
                        modifier = Modifier.size(80.dp),
                        shape = MaterialTheme.shapes.large,
                        color = Color.Gray.copy(alpha = 0.2f)
                    ) {}
                    Spacer(modifier = Modifier.height(12.dp))
                    TextButton(onClick = {}) {
                        Text("Change Photo", color = Color(0xFF007CFC))
                    }
                }
            }

            Spacer(modifier = Modifier.height(24.dp))

            // Form Fields
            OutlinedTextField(
                value = name,
                onValueChange = { name = it },
                label = { Text("Name") },
                modifier = Modifier.fillMaxWidth()
            )
            Spacer(modifier = Modifier.height(16.dp))

            OutlinedTextField(
                value = username,
                onValueChange = { username = it },
                label = { Text("Username") },
                modifier = Modifier.fillMaxWidth()
            )
            Spacer(modifier = Modifier.height(16.dp))

            OutlinedTextField(
                value = bio,
                onValueChange = { bio = it },
                label = { Text("Bio") },
                modifier = Modifier.fillMaxWidth(),
                minLines = 3
            )
            Spacer(modifier = Modifier.height(16.dp))

            OutlinedTextField(
                value = website,
                onValueChange = { website = it },
                label = { Text("Website") },
                modifier = Modifier.fillMaxWidth()
            )
        }
    }
}

// MARK: - ViewModel
class VignetteSettingsViewModel : ViewModel() {
    private val _darkModeEnabled = MutableStateFlow(false)
    val darkModeEnabled: StateFlow<Boolean> = _darkModeEnabled

    fun toggleDarkMode() {
        _darkModeEnabled.value = !_darkModeEnabled.value
        // TODO: Save to backend
    }

    fun logout() {
        // TODO: Clear tokens and navigate to login
    }
}
