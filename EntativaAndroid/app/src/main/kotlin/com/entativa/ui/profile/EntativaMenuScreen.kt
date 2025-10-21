package com.entativa.ui.profile

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import coil.compose.AsyncImage
import com.entativa.R

// MARK: - Entativa Menu Screen (Facebook-Style Profile)
@Composable
fun EntativaMenuScreen(
    viewModel: EntativaMenuViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
) {
    val profile by viewModel.profile.collectAsState()
    
    LazyColumn(
        modifier = Modifier
            .fillMaxSize()
            .background(Color(0xFFF0F2F5))
    ) {
        // Profile header
        item {
            Surface(
                modifier = Modifier.fillMaxWidth(),
                color = Color.White
            ) {
                Column {
                    // Cover photo + profile pic
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(180.dp)
                    ) {
                        // Cover photo
                        if (profile.coverImageUrl != null) {
                            AsyncImage(
                                model = profile.coverImageUrl,
                                contentDescription = "Cover photo",
                                modifier = Modifier.fillMaxSize(),
                                contentScale = ContentScale.Crop
                            )
                        } else {
                            Box(
                                modifier = Modifier
                                    .fillMaxSize()
                                    .background(
                                        Brush.horizontalGradient(
                                            colors = listOf(
                                                Color(0xFF007CFC),
                                                Color(0xFF6F3EFB),
                                                Color(0xFFFC30E1)
                                            )
                                        )
                                    )
                            )
                        }
                        
                        // Profile picture
                        Box(
                            modifier = Modifier
                                .align(Alignment.BottomStart)
                                .padding(start = 16.dp)
                                .offset(y = 30.dp)
                        ) {
                            Surface(
                                modifier = Modifier.size(120.dp),
                                shape = CircleShape,
                                color = Color.White
                            ) {
                                Box(
                                    modifier = Modifier.fillMaxSize(),
                                    contentAlignment = Alignment.Center
                                ) {
                                    if (profile.avatarUrl != null) {
                                        AsyncImage(
                                            model = profile.avatarUrl,
                                            contentDescription = "Profile",
                                            modifier = Modifier
                                                .size(112.dp)
                                                .clip(CircleShape),
                                            contentScale = ContentScale.Crop
                                        )
                                    } else {
                                        Icon(
                                            painter = painterResource(R.drawable.ic_person),
                                            contentDescription = "Profile",
                                            tint = Color.Gray,
                                            modifier = Modifier.size(112.dp)
                                        )
                                    }
                                }
                            }
                            
                            // Camera button
                            Surface(
                                modifier = Modifier
                                    .size(32.dp)
                                    .align(Alignment.BottomEnd),
                                shape = CircleShape,
                                color = Color.Gray
                            ) {
                                Box(
                                    modifier = Modifier.fillMaxSize(),
                                    contentAlignment = Alignment.Center
                                ) {
                                    Icon(
                                        painter = painterResource(R.drawable.ic_camera),
                                        contentDescription = "Change photo",
                                        tint = Color.White,
                                        modifier = Modifier.size(14.dp)
                                    )
                                }
                            }
                        }
                    }
                    
                    // Profile info
                    Column(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(top = 40.dp, start = 16.dp, end = 16.dp, bottom = 16.dp)
                    ) {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.Top
                        ) {
                            Column(modifier = Modifier.weight(1f)) {
                                Text(
                                    text = profile.fullName,
                                    fontSize = 24.sp,
                                    fontWeight = FontWeight.Bold
                                )
                                
                                Text(
                                    text = "@${profile.username}",
                                    fontSize = 15.sp,
                                    color = Color.Gray
                                )
                                
                                if (profile.bio != null) {
                                    Spacer(modifier = Modifier.height(8.dp))
                                    Text(
                                        text = profile.bio,
                                        fontSize = 15.sp
                                    )
                                }
                            }
                            
                            IconButton(onClick = {}) {
                                Icon(
                                    painter = painterResource(R.drawable.ic_settings),
                                    contentDescription = "Settings",
                                    tint = Color.Gray,
                                    modifier = Modifier.size(22.dp)
                                )
                            }
                        }
                        
                        Spacer(modifier = Modifier.height(12.dp))
                        
                        // Stats
                        Row(
                            horizontalArrangement = Arrangement.spacedBy(20.dp)
                        ) {
                            StatItem(count = profile.friendsCount, label = "Friends")
                            StatItem(count = profile.followersCount, label = "Followers")
                            StatItem(count = profile.followingCount, label = "Following")
                        }
                        
                        Spacer(modifier = Modifier.height(12.dp))
                        
                        // Action buttons
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.spacedBy(12.dp)
                        ) {
                            Button(
                                onClick = {},
                                modifier = Modifier.weight(1f),
                                colors = ButtonDefaults.buttonColors(
                                    containerColor = Color(0xFF007CFC)
                                ),
                                shape = RoundedCornerShape(8.dp)
                            ) {
                                Icon(
                                    painter = painterResource(R.drawable.ic_plus),
                                    contentDescription = null,
                                    modifier = Modifier.size(14.dp)
                                )
                                Spacer(modifier = Modifier.width(8.dp))
                                Text("Add Story", fontSize = 15.sp)
                            }
                            
                            OutlinedButton(
                                onClick = {},
                                modifier = Modifier.weight(1f),
                                shape = RoundedCornerShape(8.dp),
                                colors = ButtonDefaults.outlinedButtonColors(
                                    contentColor = Color.Black
                                )
                            ) {
                                Icon(
                                    painter = painterResource(R.drawable.ic_edit),
                                    contentDescription = null,
                                    modifier = Modifier.size(14.dp)
                                )
                                Spacer(modifier = Modifier.width(8.dp))
                                Text("Edit Profile", fontSize = 15.sp)
                            }
                            
                            OutlinedButton(
                                onClick = {},
                                modifier = Modifier.width(44.dp).height(36.dp),
                                shape = RoundedCornerShape(8.dp),
                                contentPadding = PaddingValues(0.dp)
                            ) {
                                Icon(
                                    painter = painterResource(R.drawable.ic_more),
                                    contentDescription = "More",
                                    modifier = Modifier.size(14.dp)
                                )
                            }
                        }
                    }
                }
            }
        }
        
        item { Spacer(modifier = Modifier.height(8.dp)) }
        
        // Your Shortcuts section
        item {
            MenuSection(title = "Your Shortcuts") {
                MenuItemRow(
                    icon = R.drawable.ic_people,
                    iconColor = Color(0xFF007CFC),
                    title = "Friends",
                    onClick = {}
                )
                
                MenuItemRow(
                    icon = R.drawable.ic_clock,
                    iconColor = Color(0xFF6F3EFB),
                    title = "Memories",
                    onClick = {}
                )
                
                MenuItemRow(
                    icon = R.drawable.ic_bookmark,
                    iconColor = Color(0xFFFC30E1),
                    title = "Saved",
                    onClick = {}
                )
                
                MenuItemRow(
                    icon = R.drawable.ic_flag,
                    iconColor = Color(0xFFFF9800),
                    title = "Pages",
                    onClick = {}
                )
                
                MenuItemRow(
                    icon = R.drawable.ic_play,
                    iconColor = Color(0xFF2196F3),
                    title = "Video",
                    onClick = {}
                )
                
                MenuItemRow(
                    icon = R.drawable.ic_shop,
                    iconColor = Color(0xFF00BCD4),
                    title = "Marketplace",
                    onClick = {}
                )
                
                MenuItemRow(
                    icon = R.drawable.ic_calendar,
                    iconColor = Color(0xFFF44336),
                    title = "Events",
                    onClick = {}
                )
                
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .clickable {}
                        .padding(16.dp)
                ) {
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Icon(
                            painter = painterResource(R.drawable.ic_chevron_down),
                            contentDescription = null,
                            tint = Color.Gray,
                            modifier = Modifier.size(14.dp)
                        )
                        Spacer(modifier = Modifier.width(8.dp))
                        Text(
                            text = "See More",
                            fontSize = 15.sp,
                            fontWeight = FontWeight.SemiBold,
                            color = Color.Gray
                        )
                    }
                }
            }
        }
        
        item { Spacer(modifier = Modifier.height(8.dp)) }
        
        // Settings & Privacy section
        item {
            MenuSection(title = "Settings & Privacy") {
                MenuItemRow(
                    icon = R.drawable.ic_settings,
                    iconColor = Color.Gray,
                    title = "Settings",
                    onClick = {}
                )
                
                MenuItemRow(
                    icon = R.drawable.ic_shield,
                    iconColor = Color.Gray,
                    title = "Privacy Checkup",
                    onClick = {}
                )
                
                MenuItemRow(
                    icon = R.drawable.ic_info,
                    iconColor = Color.Gray,
                    title = "Help & Support",
                    onClick = {}
                )
            }
        }
        
        item { Spacer(modifier = Modifier.height(8.dp)) }
        
        // Account section
        item {
            MenuSection(title = "") {
                MenuItemRow(
                    icon = R.drawable.ic_moon,
                    iconColor = Color.Gray,
                    title = "Dark Mode",
                    onClick = {}
                )
                
                MenuItemRow(
                    icon = R.drawable.ic_bell,
                    iconColor = Color.Gray,
                    title = "Notification Settings",
                    onClick = {}
                )
                
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .clickable {}
                        .padding(16.dp)
                ) {
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Icon(
                            painter = painterResource(R.drawable.ic_logout),
                            contentDescription = null,
                            tint = Color.Red,
                            modifier = Modifier.size(20.dp).padding(end = 12.dp)
                        )
                        Spacer(modifier = Modifier.width(12.dp))
                        Text(
                            text = "Log Out",
                            fontSize = 15.sp,
                            color = Color.Red
                        )
                    }
                }
            }
        }
    }
}

@Composable
fun MenuSection(
    title: String,
    content: @Composable () -> Unit
) {
    Surface(
        modifier = Modifier.fillMaxWidth(),
        color = Color.White
    ) {
        Column {
            if (title.isNotEmpty()) {
                Text(
                    text = title,
                    fontSize = 17.sp,
                    fontWeight = FontWeight.SemiBold,
                    modifier = Modifier.padding(16.dp)
                )
            }
            
            content()
        }
    }
}

@Composable
fun StatItem(count: Int, label: String) {
    Column(
        verticalArrangement = Arrangement.spacedBy(4.dp)
    ) {
        Text(
            text = count.toString(),
            fontSize = 18.sp,
            fontWeight = FontWeight.Bold
        )
        
        Text(
            text = label,
            fontSize = 13.sp,
            color = Color.Gray
        )
    }
}

@Composable
fun MenuItemRow(
    icon: Int,
    iconColor: Color,
    title: String,
    onClick: () -> Void
) {
    Box(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(16.dp)
    ) {
        Row(
            verticalAlignment = Alignment.CenterVertically
        ) {
            Icon(
                painter = painterResource(icon),
                contentDescription = null,
                tint = iconColor,
                modifier = Modifier.size(20.dp)
            )
            
            Spacer(modifier = Modifier.width(12.dp))
            
            Text(
                text = title,
                fontSize = 15.sp,
                modifier = Modifier.weight(1f)
            )
            
            Icon(
                painter = painterResource(R.drawable.ic_chevron_right),
                contentDescription = null,
                tint = Color.Gray,
                modifier = Modifier.size(14.dp)
            )
        }
    }
}

// Models
data class EntativaMenuProfile(
    val id: String,
    val username: String,
    val fullName: String,
    val bio: String?,
    val avatarUrl: String?,
    val coverImageUrl: String?,
    val friendsCount: Int,
    val followersCount: Int,
    val followingCount: Int
)

// ViewModel
class EntativaMenuViewModel : androidx.lifecycle.ViewModel() {
    private val _profile = kotlinx.coroutines.flow.MutableStateFlow(
        EntativaMenuProfile(
            id = "user123",
            username = "yourname",
            fullName = "Your Name",
            bio = "Welcome to my profile! 👋",
            avatarUrl = null,
            coverImageUrl = null,
            friendsCount = 342,
            followersCount = 1520,
            followingCount = 487
        )
    )
    val profile: kotlinx.coroutines.flow.StateFlow<EntativaMenuProfile> = _profile
}
