package com.entativa.ui.messages

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewmodel.compose.viewModel
import com.entativa.R
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import java.util.Date

// MARK: - E2EE Message Backup Settings
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun BackupSettingsScreen(
    onDismiss: () -> Unit,
    viewModel: BackupSettingsViewModel = viewModel()
) {
    val backupEnabled by viewModel.backupEnabled.collectAsState()
    val backupLocation by viewModel.backupLocation.collectAsState()
    val autoBackupEnabled by viewModel.autoBackupEnabled.collectAsState()
    val hasBackupKey by viewModel.hasBackupKey.collectAsState()
    var showPINSetup by remember { mutableStateOf(false) }
    var showBackupNow by remember { mutableStateOf(false) }
    var showThirdPartyWarning by remember { mutableStateOf(false) }
    
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Message Backups") },
                navigationIcon = {
                    IconButton(onClick = onDismiss) {
                        Icon(painterResource(R.drawable.ic_back), "Back")
                    }
                },
                actions = {
                    TextButton(onClick = {
                        viewModel.saveSettings()
                        onDismiss()
                    }) {
                        Text("Done")
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
            // Backup Status
            item {
                Surface(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(16.dp),
                    shape = RoundedCornerShape(16.dp),
                    color = if (backupEnabled) Color(0xFF4CAF50).copy(alpha = 0.1f) else Color.Gray.copy(alpha = 0.1f)
                ) {
                    Row(
                        modifier = Modifier.padding(20.dp),
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Icon(
                            painter = painterResource(
                                if (backupEnabled) R.drawable.ic_shield_check else R.drawable.ic_shield_x
                            ),
                            contentDescription = null,
                            tint = if (backupEnabled) Color(0xFF4CAF50) else Color.Gray,
                            modifier = Modifier.size(48.dp)
                        )
                        
                        Spacer(modifier = Modifier.width(16.dp))
                        
                        Column {
                            Text(
                                if (backupEnabled) "Backups Enabled" else "Backups Disabled",
                                fontSize = 17.sp,
                                fontWeight = FontWeight.SemiBold
                            )
                            
                            Text(
                                "Last backup: Never", // TODO: real data
                                fontSize = 14.sp,
                                color = Color.Gray
                            )
                        }
                    }
                }
            }
            
            // Backup Location Section
            item {
                SectionHeader("Backup Location")
            }
            
            item {
                Column(modifier = Modifier.padding(horizontal = 16.dp)) {
                    // Enable/Disable Toggle
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(vertical = 12.dp),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Text("Enable Backups", fontSize = 16.sp)
                        Switch(
                            checked = backupEnabled,
                            onCheckedChange = { viewModel.toggleBackup() }
                        )
                    }
                    
                    if (backupEnabled) {
                        Spacer(modifier = Modifier.height(16.dp))
                        
                        // Our Servers (Recommended)
                        BackupLocationCard(
                            icon = R.drawable.ic_server,
                            title = "Our Servers",
                            subtitle = "🔒 Most Secure - End-to-end encrypted with your PIN",
                            badge = "RECOMMENDED",
                            badgeColor = Color(0xFF4CAF50),
                            isSelected = backupLocation == BackupLocation.OUR_SERVERS,
                            onClick = { viewModel.selectBackupLocation(BackupLocation.OUR_SERVERS) }
                        )
                        
                        Spacer(modifier = Modifier.height(12.dp))
                        
                        // Google Drive
                        BackupLocationCard(
                            icon = R.drawable.ic_google_drive,
                            title = "Google Drive",
                            subtitle = "⚠️ Google can decrypt if pressured by authorities",
                            badge = null,
                            badgeColor = Color(0xFFFF9800),
                            isSelected = backupLocation == BackupLocation.GOOGLE_DRIVE,
                            onClick = { showThirdPartyWarning = true }
                        )
                    }
                }
            }
            
            // Warning Footer
            if (backupEnabled) {
                item {
                    Text(
                        if (backupLocation == BackupLocation.OUR_SERVERS)
                            "Your messages are encrypted with your PIN/passphrase. Only you can decrypt them."
                        else
                            "⚠️ Third-party providers (Google, Apple) may decrypt your backups if pressured by authorities or at their own discretion.",
                        fontSize = 13.sp,
                        color = if (backupLocation == BackupLocation.OUR_SERVERS) Color.Gray else Color(0xFFFF9800),
                        modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp)
                    )
                }
            }
            
            // Auto-Backup Section
            if (backupEnabled) {
                item {
                    SectionHeader("Automatic Backups")
                }
                
                item {
                    Column(modifier = Modifier.padding(horizontal = 16.dp)) {
                        // Auto-Backup Toggle
                        Row(
                            modifier = Modifier
                                .fillMaxWidth()
                                .padding(vertical = 12.dp),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text("Auto-Backup", fontSize = 16.sp)
                            Switch(
                                checked = autoBackupEnabled,
                                onCheckedChange = { viewModel.toggleAutoBackup() }
                            )
                        }
                        
                        // Frequency (if auto-backup enabled)
                        // Wi-Fi Only toggle
                    }
                }
            }
            
            // Manual Backup Section
            if (backupEnabled) {
                item {
                    SectionHeader("Manual Backup")
                }
                
                item {
                    Column(modifier = Modifier.padding(horizontal = 16.dp)) {
                        if (!hasBackupKey) {
                            // Setup PIN button
                            Surface(
                                onClick = { showPINSetup = true },
                                modifier = Modifier.fillMaxWidth(),
                                shape = RoundedCornerShape(12.dp),
                                color = MaterialTheme.colorScheme.secondaryContainer
                            ) {
                                Row(
                                    modifier = Modifier.padding(16.dp),
                                    verticalAlignment = Alignment.CenterVertically
                                ) {
                                    Icon(
                                        painterResource(R.drawable.ic_key),
                                        contentDescription = null,
                                        tint = Color(0xFF007CFC)
                                    )
                                    Spacer(modifier = Modifier.width(12.dp))
                                    Text("Set Up Backup PIN")
                                    Spacer(modifier = Modifier.weight(1f))
                                    Icon(
                                        painterResource(R.drawable.ic_chevron_right),
                                        contentDescription = null,
                                        tint = Color.Gray
                                    )
                                }
                            }
                        } else {
                            // Backup Now button
                            Surface(
                                onClick = { showBackupNow = true },
                                modifier = Modifier.fillMaxWidth(),
                                shape = RoundedCornerShape(12.dp),
                                color = MaterialTheme.colorScheme.secondaryContainer
                            ) {
                                Row(
                                    modifier = Modifier.padding(16.dp),
                                    verticalAlignment = Alignment.CenterVertically
                                ) {
                                    Icon(
                                        painterResource(R.drawable.ic_upload),
                                        contentDescription = null,
                                        tint = Color(0xFF007CFC)
                                    )
                                    Spacer(modifier = Modifier.width(12.dp))
                                    Text("Backup Now")
                                    Spacer(modifier = Modifier.weight(1f))
                                    Icon(
                                        painterResource(R.drawable.ic_chevron_right),
                                        contentDescription = null,
                                        tint = Color.Gray
                                    )
                                }
                            }
                            
                            Spacer(modifier = Modifier.height(12.dp))
                            
                            // Backup History button
                            Surface(
                                onClick = { /* Navigate to history */ },
                                modifier = Modifier.fillMaxWidth(),
                                shape = RoundedCornerShape(12.dp),
                                color = MaterialTheme.colorScheme.secondaryContainer
                            ) {
                                Row(
                                    modifier = Modifier.padding(16.dp),
                                    verticalAlignment = Alignment.CenterVertically
                                ) {
                                    Icon(
                                        painterResource(R.drawable.ic_history),
                                        contentDescription = null,
                                        tint = Color(0xFF007CFC)
                                    )
                                    Spacer(modifier = Modifier.width(12.dp))
                                    Text("Backup History")
                                    Spacer(modifier = Modifier.weight(1f))
                                    Icon(
                                        painterResource(R.drawable.ic_chevron_right),
                                        contentDescription = null,
                                        tint = Color.Gray
                                    )
                                }
                            }
                        }
                    }
                }
            }
            
            // Danger Zone
            if (backupEnabled && hasBackupKey) {
                item {
                    SectionHeader("Danger Zone")
                }
                
                item {
                    Column(modifier = Modifier.padding(horizontal = 16.dp)) {
                        Surface(
                            onClick = { viewModel.deleteAllBackups() },
                            modifier = Modifier.fillMaxWidth(),
                            shape = RoundedCornerShape(12.dp),
                            color = Color.Red.copy(alpha = 0.1f)
                        ) {
                            Row(
                                modifier = Modifier.padding(16.dp),
                                verticalAlignment = Alignment.CenterVertically
                            ) {
                                Icon(
                                    painterResource(R.drawable.ic_trash),
                                    contentDescription = null,
                                    tint = Color.Red
                                )
                                Spacer(modifier = Modifier.width(12.dp))
                                Text("Delete All Backups", color = Color.Red)
                            }
                        }
                        
                        Spacer(modifier = Modifier.height(8.dp))
                        
                        Text(
                            "This will permanently delete all your message backups. This cannot be undone.",
                            fontSize = 13.sp,
                            color = Color.Gray
                        )
                    }
                }
            }
            
            item {
                Spacer(modifier = Modifier.height(32.dp))
            }
        }
    }
    
    // Third-Party Warning Dialog
    if (showThirdPartyWarning) {
        AlertDialog(
            onDismissRequest = { showThirdPartyWarning = false },
            icon = {
                Icon(
                    painterResource(R.drawable.ic_warning),
                    contentDescription = null,
                    tint = Color(0xFFFF9800),
                    modifier = Modifier.size(48.dp)
                )
            },
            title = {
                Text("Third-Party Backup Warning")
            },
            text = {
                Text(
                    "⚠️ WARNING: Google (Google Drive) can decrypt your message backups if pressured by authorities or at their own discretion.\n\n" +
                    "For maximum security, we recommend using our servers where your messages are encrypted with your PIN and only you can decrypt them."
                )
            },
            confirmButton = {
                Button(
                    onClick = {
                        viewModel.acknowledgeThirdPartyWarning()
                        viewModel.selectBackupLocation(BackupLocation.GOOGLE_DRIVE)
                        showThirdPartyWarning = false
                    },
                    colors = ButtonDefaults.buttonColors(containerColor = Color(0xFFFF9800))
                ) {
                    Text("I Understand")
                }
            },
            dismissButton = {
                TextButton(onClick = { showThirdPartyWarning = false }) {
                    Text("Cancel")
                }
            }
        )
    }
    
    // PIN Setup Sheet
    if (showPINSetup) {
        BackupPINSetupSheet(
            onDismiss = { showPINSetup = false },
            viewModel = viewModel
        )
    }
    
    // Backup Now Sheet
    if (showBackupNow) {
        BackupNowSheet(
            onDismiss = { showBackupNow = false },
            viewModel = viewModel
        )
    }
}

// MARK: - Components

@Composable
fun SectionHeader(title: String) {
    Text(
        title,
        fontSize = 13.sp,
        fontWeight = FontWeight.SemiBold,
        color = Color.Gray,
        modifier = Modifier.padding(horizontal = 16.dp, vertical = 16.dp)
    )
}

@Composable
fun BackupLocationCard(
    icon: Int,
    title: String,
    subtitle: String,
    badge: String?,
    badgeColor: Color,
    isSelected: Boolean,
    onClick: () -> Unit
) {
    Surface(
        onClick = onClick,
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(16.dp),
        color = if (isSelected) Color(0xFF007CFC).copy(alpha = 0.1f) else MaterialTheme.colorScheme.secondaryContainer,
        border = if (isSelected) androidx.compose.foundation.BorderStroke(2.dp, Color(0xFF007CFC)) else null
    ) {
        Row(
            modifier = Modifier.padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Icon(
                painter = painterResource(icon),
                contentDescription = null,
                tint = if (isSelected) Color(0xFF007CFC) else Color.Gray,
                modifier = Modifier.size(36.dp)
            )
            
            Spacer(modifier = Modifier.width(12.dp))
            
            Column(modifier = Modifier.weight(1f)) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Text(
                        title,
                        fontSize = 16.sp,
                        fontWeight = FontWeight.SemiBold
                    )
                    
                    if (badge != null) {
                        Spacer(modifier = Modifier.width(8.dp))
                        Surface(
                            shape = RoundedCornerShape(4.dp),
                            color = badgeColor
                        ) {
                            Text(
                                badge,
                                fontSize = 10.sp,
                                fontWeight = FontWeight.Bold,
                                color = Color.White,
                                modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp)
                            )
                        }
                    }
                }
                
                Spacer(modifier = Modifier.height(4.dp))
                
                Text(
                    subtitle,
                    fontSize = 13.sp,
                    color = Color.Gray
                )
            }
            
            if (isSelected) {
                Icon(
                    painter = painterResource(R.drawable.ic_check_circle),
                    contentDescription = null,
                    tint = Color(0xFF007CFC),
                    modifier = Modifier.size(22.dp)
                )
            }
        }
    }
}

@Composable
fun BackupPINSetupSheet(onDismiss: () -> Unit, viewModel: BackupSettingsViewModel) {
    // TODO: Implement PIN setup dialog
    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("Set Up Backup PIN") },
        text = { Text("PIN setup coming...") },
        confirmButton = {
            Button(onClick = onDismiss) {
                Text("OK")
            }
        }
    )
}

@Composable
fun BackupNowSheet(onDismiss: () -> Unit, viewModel: BackupSettingsViewModel) {
    // TODO: Implement backup now dialog
    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("Backup Messages") },
        text = { Text("Backup functionality coming...") },
        confirmButton = {
            Button(onClick = onDismiss) {
                Text("OK")
            }
        }
    )
}

// MARK: - Models

enum class BackupLocation {
    OUR_SERVERS,
    GOOGLE_DRIVE
}

// MARK: - ViewModel

class BackupSettingsViewModel : ViewModel() {
    private val _backupEnabled = MutableStateFlow(true)
    val backupEnabled: StateFlow<Boolean> = _backupEnabled
    
    private val _backupLocation = MutableStateFlow(BackupLocation.OUR_SERVERS)
    val backupLocation: StateFlow<BackupLocation> = _backupLocation
    
    private val _autoBackupEnabled = MutableStateFlow(true)
    val autoBackupEnabled: StateFlow<Boolean> = _autoBackupEnabled
    
    private val _hasBackupKey = MutableStateFlow(false)
    val hasBackupKey: StateFlow<Boolean> = _hasBackupKey
    
    fun toggleBackup() {
        _backupEnabled.value = !_backupEnabled.value
        // TODO: API call
    }
    
    fun selectBackupLocation(location: BackupLocation) {
        _backupLocation.value = location
        // TODO: API call
    }
    
    fun toggleAutoBackup() {
        _autoBackupEnabled.value = !_autoBackupEnabled.value
        // TODO: API call
    }
    
    fun acknowledgeThirdPartyWarning() {
        // TODO: API call
    }
    
    fun saveSettings() {
        // TODO: API call
    }
    
    fun deleteAllBackups() {
        // TODO: API call with confirmation
    }
}
