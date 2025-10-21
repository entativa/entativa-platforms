package com.entativa.ui.admin

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
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

// MARK: - Admin Panel (Founder Only)
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun AdminPanelScreen(
    onDismiss: () -> Unit,
    viewModel: AdminPanelViewModel = viewModel()
) {
    var selectedSection by remember { mutableStateOf(AdminSection.DASHBOARD) }
    
    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Icon(
                            painter = painterResource(R.drawable.ic_crown),
                            contentDescription = "Founder",
                            tint = Color(0xFFFFD700),
                            modifier = Modifier.size(24.dp)
                        )
                        Spacer(modifier = Modifier.width(8.dp))
                        Text("Admin Panel", fontWeight = FontWeight.Bold)
                    }
                },
                actions = {
                    // Session indicator
                    Surface(
                        shape = RoundedCornerShape(12.dp),
                        color = Color.Green.copy(alpha = 0.1f)
                    ) {
                        Row(
                            modifier = Modifier.padding(horizontal = 12.dp, vertical = 6.dp),
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Box(
                                modifier = Modifier
                                    .size(8.dp)
                                    .background(Color.Green, CircleShape)
                            )
                            Spacer(modifier = Modifier.width(6.dp))
                            Text(
                                "Admin Active",
                                fontSize = 12.sp,
                                fontWeight = FontWeight.Medium,
                                color = Color.Green
                            )
                        }
                    }
                    
                    IconButton(onClick = onDismiss) {
                        Icon(painter = painterResource(R.drawable.ic_close), contentDescription = "Close")
                    }
                }
            )
        }
    ) { paddingValues ->
        Row(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            // Side Navigation
            AdminSideNavigation(
                selectedSection = selectedSection,
                onSectionSelected = { selectedSection = it }
            )
            
            // Main Content
            Box(
                modifier = Modifier
                    .weight(1f)
                    .fillMaxHeight()
            ) {
                when (selectedSection) {
                    AdminSection.DASHBOARD -> AdminDashboardContent(viewModel)
                    AdminSection.USERS -> AdminUsersContent(viewModel)
                    AdminSection.CONTENT -> AdminContentContent(viewModel)
                    AdminSection.PLATFORM -> AdminPlatformContent(viewModel)
                    AdminSection.ANALYTICS -> AdminAnalyticsContent(viewModel)
                    AdminSection.DEVELOPER -> AdminDeveloperContent(viewModel)
                    AdminSection.SECURITY -> AdminSecurityContent(viewModel)
                    AdminSection.AUDIT -> AdminAuditContent(viewModel)
                }
            }
        }
    }
}

// MARK: - Admin Sections
enum class AdminSection(val title: String, val icon: Int, val color: Color) {
    DASHBOARD("Dashboard", R.drawable.ic_dashboard, Color(0xFF007CFC)),
    USERS("Users", R.drawable.ic_people, Color(0xFF6F3EFB)),
    CONTENT("Content", R.drawable.ic_document, Color.Green),
    PLATFORM("Platform", R.drawable.ic_settings, Color(0xFFFF9800)),
    ANALYTICS("Analytics", R.drawable.ic_chart, Color.Cyan),
    DEVELOPER("Developer", R.drawable.ic_code, Color(0xFFFC30E1)),
    SECURITY("Security", R.drawable.ic_shield, Color.Red),
    AUDIT("Audit", R.drawable.ic_list, Color(0xFF9C27B0))
}

// MARK: - Side Navigation
@Composable
fun AdminSideNavigation(
    selectedSection: AdminSection,
    onSectionSelected: (AdminSection) -> Unit
) {
    Surface(
        modifier = Modifier
            .width(80.dp)
            .fillMaxHeight(),
        color = Color.Black
    ) {
        Column {
            // Crown icon
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(vertical = 20.dp),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Icon(
                    painter = painterResource(R.drawable.ic_crown),
                    contentDescription = "Founder",
                    tint = Color(0xFFFFD700),
                    modifier = Modifier.size(24.dp)
                )
                Spacer(modifier = Modifier.height(4.dp))
                Text(
                    "FOUNDER",
                    fontSize = 8.sp,
                    fontWeight = FontWeight.Bold,
                    color = Color(0xFFFFD700)
                )
            }
            
            Divider(color = Color.White.copy(alpha = 0.2f))
            
            // Navigation buttons
            AdminSection.values().forEach { section ->
                AdminNavButton(
                    section = section,
                    isSelected = selectedSection == section,
                    onClick = { onSectionSelected(section) }
                )
            }
        }
    }
}

@Composable
fun AdminNavButton(
    section: AdminSection,
    isSelected: Boolean,
    onClick: () -> Void
) {
    Column(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .background(if (isSelected) Color.White.copy(alpha = 0.15f) else Color.Transparent)
            .padding(vertical = 12.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Icon(
            painter = painterResource(section.icon),
            contentDescription = section.title,
            tint = if (isSelected) section.color else Color.White.copy(alpha = 0.6f),
            modifier = Modifier.size(22.dp)
        )
        Spacer(modifier = Modifier.height(6.dp))
        Text(
            section.title,
            fontSize = 9.sp,
            fontWeight = if (isSelected) FontWeight.SemiBold else FontWeight.Normal,
            color = if (isSelected) section.color else Color.White.copy(alpha = 0.6f),
            maxLines = 1
        )
    }
}

// MARK: - Dashboard Content
@Composable
fun AdminDashboardContent(viewModel: AdminPanelViewModel) {
    val liveMetrics by viewModel.liveMetrics.collectAsState()
    
    LazyColumn(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp)
    ) {
        item {
            Text(
                "Live Platform Metrics",
                fontSize = 20.sp,
                fontWeight = FontWeight.Bold
            )
            Spacer(modifier = Modifier.height(16.dp))
        }
        
        item {
            LazyVerticalGrid(
                columns = GridCells.Fixed(2),
                modifier = Modifier.height(400.dp),
                horizontalArrangement = Arrangement.spacedBy(16.dp),
                verticalArrangement = Arrangement.spacedBy(16.dp)
            ) {
                item {
                    MetricCard("Active Users", "12,453", R.drawable.ic_people, Color(0xFF007CFC))
                }
                item {
                    MetricCard("Posts/Second", "47", R.drawable.ic_document, Color.Green)
                }
                item {
                    MetricCard("Server Load", "23%", R.drawable.ic_cpu, Color.Orange)
                }
                item {
                    MetricCard("Error Rate", "0.02%", R.drawable.ic_warning, Color.Red)
                }
            }
        }
        
        item {
            Spacer(modifier = Modifier.height(24.dp))
            Divider()
            Spacer(modifier = Modifier.height(24.dp))
        }
        
        item {
            Text(
                "Quick Actions",
                fontSize = 20.sp,
                fontWeight = FontWeight.Bold
            )
            Spacer(modifier = Modifier.height(16.dp))
        }
        
        item {
            QuickActionButton(
                icon = R.drawable.ic_search,
                title = "Search Users",
                color = Color(0xFF007CFC),
                onClick = {}
            )
            Spacer(modifier = Modifier.height(12.dp))
            QuickActionButton(
                icon = R.drawable.ic_trash,
                title = "Moderation Queue",
                color = Color.Red,
                badge = 5,
                onClick = {}
            )
            Spacer(modifier = Modifier.height(12.dp))
            QuickActionButton(
                icon = R.drawable.ic_bell,
                title = "Broadcast Notification",
                color = Color(0xFF9C27B0),
                onClick = {}
            )
        }
    }
}

@Composable
fun MetricCard(title: String, value: String, icon: Int, color: Color) {
    Surface(
        modifier = Modifier
            .fillMaxWidth()
            .height(120.dp),
        shape = RoundedCornerShape(12.dp),
        color = MaterialTheme.colorScheme.secondaryContainer
    ) {
        Column(
            modifier = Modifier.padding(16.dp),
            verticalArrangement = Arrangement.SpaceBetween
        ) {
            Icon(
                painter = painterResource(icon),
                contentDescription = null,
                tint = color,
                modifier = Modifier.size(24.dp)
            )
            
            Text(
                value,
                fontSize = 32.sp,
                fontWeight = FontWeight.Bold
            )
            
            Text(
                title,
                fontSize = 14.sp,
                color = Color.Gray
            )
        }
    }
}

@Composable
fun QuickActionButton(
    icon: Int,
    title: String,
    color: Color,
    badge: Int = 0,
    onClick: () -> Unit
) {
    Surface(
        onClick = onClick,
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp),
        color = MaterialTheme.colorScheme.secondaryContainer
    ) {
        Row(
            modifier = Modifier.padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Box {
                Surface(
                    modifier = Modifier.size(48.dp),
                    shape = RoundedCornerShape(10.dp),
                    color = color
                ) {
                    Box(contentAlignment = Alignment.Center) {
                        Icon(
                            painter = painterResource(icon),
                            contentDescription = null,
                            tint = Color.White,
                            modifier = Modifier.size(22.dp)
                        )
                    }
                }
                
                if (badge > 0) {
                    Surface(
                        modifier = Modifier
                            .align(Alignment.TopEnd)
                            .offset(x = 8.dp, y = (-8).dp),
                        shape = CircleShape,
                        color = Color.Red
                    ) {
                        Text(
                            badge.toString(),
                            fontSize = 10.sp,
                            fontWeight = FontWeight.Bold,
                            color = Color.White,
                            modifier = Modifier.padding(horizontal = 6.dp, vertical = 3.dp)
                        )
                    }
                }
            }
            
            Spacer(modifier = Modifier.width(16.dp))
            
            Text(
                title,
                fontSize = 16.sp,
                fontWeight = FontWeight.SemiBold
            )
            
            Spacer(modifier = Modifier.weight(1f))
            
            Icon(
                painter = painterResource(R.drawable.ic_chevron_right),
                contentDescription = null,
                tint = Color.Gray
            )
        }
    }
}

// MARK: - Content Screens (Placeholders)
@Composable
fun AdminUsersContent(viewModel: AdminPanelViewModel) {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Text("Users Management", fontSize = 24.sp)
    }
}

@Composable
fun AdminContentContent(viewModel: AdminPanelViewModel) {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Text("Content Moderation", fontSize = 24.sp)
    }
}

@Composable
fun AdminPlatformContent(viewModel: AdminPanelViewModel) {
    val features by viewModel.featureFlags.collectAsState()
    
    LazyColumn(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp)
    ) {
        item {
            Text(
                "Feature Flags",
                fontSize = 20.sp,
                fontWeight = FontWeight.Bold
            )
            Spacer(modifier = Modifier.height(16.dp))
        }
        
        items(features.size) { index ->
            FeatureFlagRow(
                feature = features[index],
                onToggle = { viewModel.toggleFeature(features[index].name) }
            )
            Spacer(modifier = Modifier.height(12.dp))
        }
        
        item {
            Spacer(modifier = Modifier.height(24.dp))
            Text(
                "Emergency Kill Switches",
                fontSize = 20.sp,
                fontWeight = FontWeight.Bold
            )
            Spacer(modifier = Modifier.height(16.dp))
        }
        
        item {
            KillSwitchCard(
                title = "Disable All Posting",
                isActive = false,
                onActivate = { viewModel.activateKillSwitch("posting") }
            )
            Spacer(modifier = Modifier.height(12.dp))
            KillSwitchCard(
                title = "Disable All Commenting",
                isActive = false,
                onActivate = { viewModel.activateKillSwitch("commenting") }
            )
            Spacer(modifier = Modifier.height(12.dp))
            KillSwitchCard(
                title = "Emergency Lockdown",
                isActive = false,
                onActivate = { viewModel.activateKillSwitch("lockdown") }
            )
        }
    }
}

@Composable
fun FeatureFlagRow(feature: FeatureFlag, onToggle: () -> Void) {
    Surface(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(10.dp),
        color = MaterialTheme.colorScheme.secondaryContainer
    ) {
        Row(
            modifier = Modifier.padding(12.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                feature.name,
                fontSize = 16.sp,
                fontWeight = FontWeight.Medium,
                modifier = Modifier.weight(1f)
            )
            Switch(
                checked = feature.isEnabled,
                onCheckedChange = { onToggle() }
            )
        }
    }
}

@Composable
fun KillSwitchCard(title: String, isActive: Boolean, onActivate: () -> Void) {
    Surface(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp),
        color = MaterialTheme.colorScheme.secondaryContainer,
        border = if (isActive) androidx.compose.foundation.BorderStroke(2.dp, Color.Red) else null
    ) {
        Row(
            modifier = Modifier
                .clickable(onClick = onActivate)
                .padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Box(
                modifier = Modifier
                    .size(12.dp)
                    .background(if (isActive) Color.Red else Color.Gray, CircleShape)
            )
            
            Spacer(modifier = Modifier.width(12.dp))
            
            Text(
                title,
                fontSize = 16.sp,
                fontWeight = FontWeight.Medium,
                color = if (isActive) Color.Red else Color.Black,
                modifier = Modifier.weight(1f)
            )
            
            Surface(
                shape = RoundedCornerShape(12.dp),
                color = if (isActive) Color.Red else Color.Gray
            ) {
                Text(
                    if (isActive) "ACTIVE" else "INACTIVE",
                    fontSize = 12.sp,
                    fontWeight = FontWeight.Bold,
                    color = Color.White,
                    modifier = Modifier.padding(horizontal = 12.dp, vertical = 6.dp)
                )
            }
        }
    }
}

@Composable
fun AdminAnalyticsContent(viewModel: AdminPanelViewModel) {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Text("Analytics Dashboard", fontSize = 24.sp)
    }
}

@Composable
fun AdminDeveloperContent(viewModel: AdminPanelViewModel) {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Text("Developer Tools", fontSize = 24.sp)
    }
}

@Composable
fun AdminSecurityContent(viewModel: AdminPanelViewModel) {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Text("Security & Audit", fontSize = 24.sp)
    }
}

@Composable
fun AdminAuditContent(viewModel: AdminPanelViewModel) {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Text("Audit Logs", fontSize = 24.sp)
    }
}

// MARK: - Models
data class FeatureFlag(
    val name: String,
    val isEnabled: Boolean,
    val rolloutPercentage: Int = 100
)

data class LiveMetrics(
    val activeUsers: Int,
    val postsPerSecond: Int,
    val serverLoad: Double,
    val errorRate: Double
)

// MARK: - ViewModel
class AdminPanelViewModel : ViewModel() {
    private val _featureFlags = MutableStateFlow(listOf(
        FeatureFlag("Posting", true),
        FeatureFlag("Commenting", true),
        FeatureFlag("Messaging", true),
        FeatureFlag("Takes", true),
        FeatureFlag("Stories", true),
        FeatureFlag("Live Streaming", false)
    ))
    val featureFlags: StateFlow<List<FeatureFlag>> = _featureFlags
    
    private val _liveMetrics = MutableStateFlow(
        LiveMetrics(
            activeUsers = 12453,
            postsPerSecond = 47,
            serverLoad = 23.0,
            errorRate = 0.02
        )
    )
    val liveMetrics: StateFlow<LiveMetrics> = _liveMetrics
    
    fun toggleFeature(featureName: String) {
        // TODO: API call to toggle feature
    }
    
    fun activateKillSwitch(switchName: String) {
        // TODO: API call to activate kill switch
    }
}
