package com.entativa.ui.menu

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
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
import com.entativa.R
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow

// MARK: - Entativa Menu Screen (Facebook-Style)
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun EntativaMenuScreen(
    onDismiss: () -> Unit,
    viewModel: EntativaMenuViewModel = viewModel()
) {
    val userName by viewModel.userName.collectAsState()
    val username by viewModel.username.collectAsState()
    var showLogoutDialog by remember { mutableStateOf(false) }
    var showSettings by remember { mutableStateOf(false) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Menu", fontWeight = FontWeight.Bold) },
                actions = {
                    IconButton(onClick = onDismiss) {
                        Icon(
                            painter = painterResource(R.drawable.ic_close),
                            contentDescription = "Close"
                        )
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
            // User Profile Section
            item {
                UserProfileSection(
                    userName = userName,
                    username = username,
                    onProfileClick = {}
                )
                Divider()
            }

            // Your Shortcuts
            item {
                MenuSectionHeader("Your shortcuts")
            }

            items(viewModel.getShortcuts()) { shortcut ->
                MenuShortcutItem(
                    icon = shortcut.icon,
                    title = shortcut.title,
                    color = shortcut.color,
                    badge = shortcut.badge,
                    onClick = {}
                )
            }

            item {
                Divider(modifier = Modifier.padding(vertical = 8.dp))
            }

            // Settings & Privacy
            item {
                MenuSectionHeader("Settings & privacy")
            }

            item {
                MenuListItem(
                    icon = R.drawable.ic_settings,
                    title = "Settings",
                    color = Color.Gray,
                    onClick = { showSettings = true }
                )
                MenuListItem(
                    icon = R.drawable.ic_lock,
                    title = "Privacy Center",
                    color = Color(0xFF007CFC),
                    onClick = {}
                )
                MenuListItem(
                    icon = R.drawable.ic_clock,
                    title = "Activity log",
                    color = Color(0xFFFF9800),
                    onClick = {}
                )
            }

            item {
                Divider(modifier = Modifier.padding(vertical = 8.dp))
            }

            // Help & Support
            item {
                MenuSectionHeader("Help & support")
            }

            item {
                MenuListItem(
                    icon = R.drawable.ic_help,
                    title = "Help Center",
                    color = Color(0xFF9C27B0),
                    onClick = {}
                )
                MenuListItem(
                    icon = R.drawable.ic_flag,
                    title = "Report a problem",
                    color = Color.Red,
                    onClick = {}
                )
            }

            item {
                Divider(modifier = Modifier.padding(vertical = 8.dp))
            }

            // Also From Entativa
            item {
                MenuSectionHeader("Also from Entativa")
            }

            item {
                MenuListItem(
                    icon = R.drawable.ic_camera,
                    title = "Vignette",
                    subtitle = "Photo & video sharing",
                    color = Color(0xFF519CAB),
                    onClick = {}
                )
            }

            item {
                Divider(modifier = Modifier.padding(vertical = 8.dp))
            }

            // Logout
            item {
                Button(
                    onClick = { showLogoutDialog = true },
                    colors = ButtonDefaults.buttonColors(
                        containerColor = Color.Red.copy(alpha = 0.1f),
                        contentColor = Color.Red
                    ),
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 16.dp, vertical = 8.dp)
                ) {
                    Icon(
                        painter = painterResource(R.drawable.ic_logout),
                        contentDescription = null,
                        modifier = Modifier.size(20.dp)
                    )
                    Spacer(modifier = Modifier.width(8.dp))
                    Text("Log out", fontWeight = FontWeight.SemiBold)
                }
            }

            // Legal Footer
            item {
                Column(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(16.dp),
                    horizontalAlignment = Alignment.CenterHorizontally
                ) {
                    Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                        Text("Terms", fontSize = 12.sp, color = Color.Gray)
                        Text("•", fontSize = 12.sp, color = Color.Gray)
                        Text("Privacy Policy", fontSize = 12.sp, color = Color.Gray)
                        Text("•", fontSize = 12.sp, color = Color.Gray)
                        Text("Cookies", fontSize = 12.sp, color = Color.Gray)
                    }
                    Spacer(modifier = Modifier.height(8.dp))
                    Text("Entativa © 2025", fontSize = 12.sp, color = Color.Gray)
                }
            }
        }
    }

    // Logout Confirmation Dialog
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

    // Settings Screen
    if (showSettings) {
        EntativaSettingsScreen(onDismiss = { showSettings = false })
    }
}

// MARK: - User Profile Section
@Composable
fun UserProfileSection(
    userName: String,
    username: String,
    onProfileClick: () -> Unit
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onProfileClick)
            .padding(16.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Surface(
            modifier = Modifier.size(60.dp),
            shape = CircleShape,
            color = Color.Gray.copy(alpha = 0.2f)
        ) {
            Box(contentAlignment = Alignment.Center) {
                Icon(
                    painter = painterResource(R.drawable.ic_person),
                    contentDescription = null,
                    tint = Color.Gray,
                    modifier = Modifier.size(32.dp)
                )
            }
        }

        Spacer(modifier = Modifier.width(12.dp))

        Column {
            Text(userName, fontSize = 20.sp, fontWeight = FontWeight.Bold)
            Text("@$username", fontSize = 14.sp, color = Color.Gray)
            Spacer(modifier = Modifier.height(4.dp))
            Text(
                "See your profile",
                fontSize = 14.sp,
                fontWeight = FontWeight.Medium,
                color = Color(0xFF007CFC)
            )
        }
    }
}

// MARK: - Menu Section Header
@Composable
fun MenuSectionHeader(title: String) {
    Text(
        text = title,
        fontSize = 16.sp,
        fontWeight = FontWeight.Bold,
        color = Color.Gray,
        modifier = Modifier.padding(horizontal = 16.dp, vertical = 12.dp)
    )
}

// MARK: - Menu Shortcut Item
@Composable
fun MenuShortcutItem(
    icon: Int,
    title: String,
    color: Color,
    badge: Int = 0,
    onClick: () -> Unit
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(horizontal = 16.dp, vertical = 8.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Surface(
            modifier = Modifier.size(36.dp),
            shape = CircleShape,
            color = color
        ) {
            Box(contentAlignment = Alignment.Center) {
                Icon(
                    painter = painterResource(icon),
                    contentDescription = null,
                    tint = Color.White,
                    modifier = Modifier.size(20.dp)
                )
            }
        }

        Spacer(modifier = Modifier.width(12.dp))

        Text(
            text = title,
            fontSize = 16.sp,
            fontWeight = FontWeight.Medium,
            modifier = Modifier.weight(1f)
        )

        if (badge > 0) {
            Surface(
                shape = CircleShape,
                color = Color.Red
            ) {
                Text(
                    text = badge.toString(),
                    fontSize = 12.sp,
                    fontWeight = FontWeight.Bold,
                    color = Color.White,
                    modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp)
                )
            }
        }

        Icon(
            painter = painterResource(R.drawable.ic_chevron_right),
            contentDescription = null,
            tint = Color.Gray,
            modifier = Modifier.size(20.dp)
        )
    }
}

// MARK: - Menu List Item
@Composable
fun MenuListItem(
    icon: Int,
    title: String,
    color: Color,
    subtitle: String? = null,
    onClick: () -> Unit
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(horizontal = 16.dp, vertical = 8.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Surface(
            modifier = Modifier.size(36.dp),
            shape = CircleShape,
            color = color
        ) {
            Box(contentAlignment = Alignment.Center) {
                Icon(
                    painter = painterResource(icon),
                    contentDescription = null,
                    tint = Color.White,
                    modifier = Modifier.size(20.dp)
                )
            }
        }

        Spacer(modifier = Modifier.width(12.dp))

        Column(modifier = Modifier.weight(1f)) {
            Text(
                text = title,
                fontSize = 16.sp,
                fontWeight = FontWeight.Medium
            )
            if (subtitle != null) {
                Text(
                    text = subtitle,
                    fontSize = 13.sp,
                    color = Color.Gray
                )
            }
        }

        Icon(
            painter = painterResource(R.drawable.ic_chevron_right),
            contentDescription = null,
            tint = Color.Gray,
            modifier = Modifier.size(20.dp)
        )
    }
}

// MARK: - Settings Screen
@Composable
fun EntativaSettingsScreen(onDismiss: () -> Unit) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Settings") },
                navigationIcon = {
                    IconButton(onClick = onDismiss) {
                        Icon(
                            painter = painterResource(R.drawable.ic_arrow_back),
                            contentDescription = "Back"
                        )
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
                SettingsSection("Account") {
                    SettingsItem("Account Settings", R.drawable.ic_person) {}
                    SettingsItem("Privacy & Security", R.drawable.ic_lock) {}
                    SettingsItem("Notifications", R.drawable.ic_bell) {}
                }
            }

            // Content Section
            item {
                SettingsSection("Content") {
                    SettingsItem("Language", R.drawable.ic_globe) {}
                    SettingsItem("Accessibility", R.drawable.ic_accessibility) {}
                }
            }

            // Data Section
            item {
                SettingsSection("Data") {
                    SettingsItem("Data Usage", R.drawable.ic_data) {}
                    SettingsItem("Storage", R.drawable.ic_storage) {}
                }
            }

            // About Section
            item {
                SettingsSection("About") {
                    SettingsItem("About", R.drawable.ic_info) {}
                    SettingsItem("Help Center", R.drawable.ic_help) {}
                }
            }
        }
    }
}

@Composable
fun SettingsSection(title: String, content: @Composable () -> Unit) {
    Column {
        Text(
            text = title,
            fontSize = 14.sp,
            fontWeight = FontWeight.Bold,
            color = Color.Gray,
            modifier = Modifier.padding(horizontal = 16.dp, vertical = 12.dp)
        )
        content()
        Divider(modifier = Modifier.padding(vertical = 8.dp))
    }
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
            tint = Color(0xFF007CFC),
            modifier = Modifier.size(24.dp)
        )
        Spacer(modifier = Modifier.width(12.dp))
        Text(title, fontSize = 16.sp)
        Spacer(modifier = Modifier.weight(1f))
        Icon(
            painter = painterResource(R.drawable.ic_chevron_right),
            contentDescription = null,
            tint = Color.Gray
        )
    }
}

// MARK: - Data Classes
data class MenuShortcut(
    val icon: Int,
    val title: String,
    val color: Color,
    val badge: Int = 0
)

// MARK: - ViewModel
class EntativaMenuViewModel : ViewModel() {
    private val _userName = MutableStateFlow("John Doe")
    val userName: StateFlow<String> = _userName

    private val _username = MutableStateFlow("johndoe")
    val username: StateFlow<String> = _username

    fun getShortcuts(): List<MenuShortcut> {
        return listOf(
            MenuShortcut(R.drawable.ic_people, "Friends", Color(0xFF007CFC), badge = 3),
            MenuShortcut(R.drawable.ic_clock, "Memories", Color(0xFF2196F3)),
            MenuShortcut(R.drawable.ic_bookmark, "Saved", Color(0xFF9C27B0)),
            MenuShortcut(R.drawable.ic_people, "Groups", Color(0xFF007CFC), badge = 5),
            MenuShortcut(R.drawable.ic_video, "Video", Color(0xFF00BCD4)),
            MenuShortcut(R.drawable.ic_shopping, "Marketplace", Color(0xFF4CAF50)),
            MenuShortcut(R.drawable.ic_calendar, "Events", Color(0xFFF44336))
        )
    }

    fun logout() {
        // TODO: Clear tokens and navigate to login
    }
}
