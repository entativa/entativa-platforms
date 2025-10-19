package com.entativa.vignette.ui.profile

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.blur
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import coil.compose.AsyncImage
import com.entativa.vignette.R

// MARK: - Vignette Profile Screen (Full-Bleed Immersive Design)
@Composable
fun VignetteProfileScreen(
    viewModel: VignetteProfileViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
) {
    val profile by viewModel.profile.collectAsState()
    val posts by viewModel.posts.collectAsState()
    var selectedTab by remember { mutableStateOf(ProfileTab.POSTS) }
    var showSettings by remember { mutableStateOf(false) }
    
    Box(modifier = Modifier.fillMaxSize()) {
        // Full-bleed background image
        if (profile.profileImageUrl != null) {
            AsyncImage(
                model = profile.profileImageUrl,
                contentDescription = null,
                modifier = Modifier
                    .fillMaxSize()
                    .blur(radius = 0.dp),
                contentScale = ContentScale.Crop
            )
        } else {
            // Default gradient background
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .background(
                        Brush.verticalGradient(
                            colors = listOf(
                                Color(0xFFC3E7F1),
                                Color(0xFF519CAB),
                                Color(0xFF20373B)
                            )
                        )
                    )
            )
        }
        
        // Dark gradient overlay for readability
        Box(
            modifier = Modifier
                .fillMaxSize()
                .background(
                    Brush.verticalGradient(
                        colors = listOf(
                            Color.Black.copy(alpha = 0.6f),
                            Color.Black.copy(alpha = 0.3f),
                            Color.Black.copy(alpha = 0.6f)
                        )
                    )
                )
        )
        
        // Content layer
        LazyColumn(
            modifier = Modifier.fillMaxSize()
        ) {
            // Header section
            item {
                Column(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(top = 50.dp, start = 16.dp, end = 16.dp)
                ) {
                    // Top bar
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            Icon(
                                painter = painterResource(R.drawable.ic_lock),
                                contentDescription = "Private",
                                tint = Color.White,
                                modifier = Modifier.size(14.dp)
                            )
                            
                            Spacer(modifier = Modifier.width(8.dp))
                            
                            Text(
                                text = profile.username,
                                fontSize = 18.sp,
                                fontWeight = FontWeight.Bold,
                                color = Color.White
                            )
                            
                            Spacer(modifier = Modifier.width(8.dp))
                            
                            Icon(
                                painter = painterResource(R.drawable.ic_chevron_down),
                                contentDescription = "Dropdown",
                                tint = Color.White,
                                modifier = Modifier.size(14.dp)
                            )
                        }
                        
                        Row {
                            IconButton(onClick = {}) {
                                Icon(
                                    painter = painterResource(R.drawable.ic_plus_circle),
                                    contentDescription = "Add",
                                    tint = Color.White,
                                    modifier = Modifier.size(24.dp)
                                )
                            }
                            
                            IconButton(onClick = { showSettings = true }) {
                                Icon(
                                    painter = painterResource(R.drawable.ic_menu),
                                    contentDescription = "Menu",
                                    tint = Color.White,
                                    modifier = Modifier.size(24.dp)
                                )
                            }
                        }
                    }
                    
                    Spacer(modifier = Modifier.height(20.dp))
                    
                    // Profile picture with gradient border
                    Box(
                        modifier = Modifier.align(Alignment.CenterHorizontally)
                    ) {
                        // Gradient border
                        Box(
                            modifier = Modifier
                                .size(90.dp)
                                .clip(CircleShape)
                                .background(
                                    Brush.linearGradient(
                                        colors = listOf(
                                            Color(0xFFFFC64F),
                                            Color(0xFFFC30E1),
                                            Color(0xFF6F3EFB)
                                        )
                                    )
                                ),
                            contentAlignment = Alignment.Center
                        ) {
                            // Inner circle
                            Box(
                                modifier = Modifier
                                    .size(84.dp)
                                    .clip(CircleShape)
                                    .background(Color.White.copy(alpha = 0.2f)),
                                contentAlignment = Alignment.Center
                            ) {
                                if (profile.avatarUrl != null) {
                                    AsyncImage(
                                        model = profile.avatarUrl,
                                        contentDescription = "Profile",
                                        modifier = Modifier
                                            .size(80.dp)
                                            .clip(CircleShape),
                                        contentScale = ContentScale.Crop
                                    )
                                } else {
                                    Icon(
                                        painter = painterResource(R.drawable.ic_person),
                                        contentDescription = "Profile",
                                        tint = Color.White.copy(alpha = 0.5f),
                                        modifier = Modifier.size(80.dp)
                                    )
                                }
                            }
                        }
                    }
                    
                    Spacer(modifier = Modifier.height(16.dp))
                    
                    // Name and bio
                    Column(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalAlignment = Alignment.CenterHorizontally
                    ) {
                        Text(
                            text = profile.fullName,
                            fontSize = 16.sp,
                            fontWeight = FontWeight.SemiBold,
                            color = Color.White
                        )
                        
                        if (profile.bio != null) {
                            Spacer(modifier = Modifier.height(8.dp))
                            Text(
                                text = profile.bio,
                                fontSize = 14.sp,
                                color = Color.White.copy(alpha = 0.9f),
                                modifier = Modifier.padding(horizontal = 32.dp)
                            )
                        }
                        
                        if (profile.link != null) {
                            Spacer(modifier = Modifier.height(4.dp))
                            Text(
                                text = profile.link,
                                fontSize = 14.sp,
                                fontWeight = FontWeight.Medium,
                                color = Color(0xFFFFC64F),
                                modifier = Modifier.clickable {}
                            )
                        }
                    }
                    
                    Spacer(modifier = Modifier.height(20.dp))
                    
                    // Stats row (frosted glass)
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.SpaceEvenly
                    ) {
                        StatButton(count = profile.postsCount, label = "Posts")
                        StatButton(count = profile.followersCount, label = "Followers")
                        StatButton(count = profile.followingCount, label = "Following")
                    }
                    
                    Spacer(modifier = Modifier.height(12.dp))
                    
                    // Action buttons (frosted glass)
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.spacedBy(12.dp)
                    ) {
                        Surface(
                            modifier = Modifier
                                .weight(1f)
                                .height(32.dp),
                            shape = RoundedCornerShape(8.dp),
                            color = Color.White.copy(alpha = 0.2f)
                        ) {
                            Box(
                                modifier = Modifier
                                    .fillMaxSize()
                                    .clickable {},
                                contentAlignment = Alignment.Center
                            ) {
                                Text(
                                    text = "Edit Profile",
                                    fontSize = 14.sp,
                                    fontWeight = FontWeight.SemiBold,
                                    color = Color.White
                                )
                            }
                        }
                        
                        Surface(
                            modifier = Modifier
                                .weight(1f)
                                .height(32.dp),
                            shape = RoundedCornerShape(8.dp),
                            color = Color.White.copy(alpha = 0.2f)
                        ) {
                            Box(
                                modifier = Modifier
                                    .fillMaxSize()
                                    .clickable {},
                                contentAlignment = Alignment.Center
                            ) {
                                Text(
                                    text = "Share Profile",
                                    fontSize = 14.sp,
                                    fontWeight = FontWeight.SemiBold,
                                    color = Color.White
                                )
                            }
                        }
                    }
                    
                    Spacer(modifier = Modifier.height(20.dp))
                    
                    // Story highlights
                    LazyRow(
                        horizontalArrangement = Arrangement.spacedBy(16.dp)
                    ) {
                        // Add new highlight
                        item {
                            Column(
                                horizontalAlignment = Alignment.CenterHorizontally,
                                verticalArrangement = Arrangement.spacedBy(8.dp)
                            ) {
                                Surface(
                                    modifier = Modifier.size(64.dp),
                                    shape = CircleShape,
                                    color = Color.White.copy(alpha = 0.2f)
                                ) {
                                    Box(
                                        modifier = Modifier.fillMaxSize(),
                                        contentAlignment = Alignment.Center
                                    ) {
                                        Icon(
                                            painter = painterResource(R.drawable.ic_plus),
                                            contentDescription = "Add",
                                            tint = Color.White,
                                            modifier = Modifier.size(24.dp)
                                        )
                                    }
                                }
                                
                                Text(
                                    text = "New",
                                    fontSize = 12.sp,
                                    color = Color.White.copy(alpha = 0.9f)
                                )
                            }
                        }
                        
                        // Existing highlights
                        items(3) { index ->
                            HighlightItem(title = listOf("Travel", "Food", "Work")[index])
                        }
                    }
                    
                    Spacer(modifier = Modifier.height(20.dp))
                }
            }
            
            // Tab selector (frosted glass)
            item {
                Surface(
                    modifier = Modifier.fillMaxWidth(),
                    color = Color.White.copy(alpha = 0.2f)
                ) {
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(44.dp)
                    ) {
                        TabButton(
                            icon = R.drawable.ic_photo,
                            isSelected = selectedTab == ProfileTab.POSTS,
                            onClick = { selectedTab = ProfileTab.POSTS },
                            modifier = Modifier.weight(1f)
                        )
                        
                        TabButton(
                            icon = R.drawable.ic_play_rect,
                            isSelected = selectedTab == ProfileTab.REELS,
                            onClick = { selectedTab = ProfileTab.REELS },
                            modifier = Modifier.weight(1f)
                        )
                        
                        TabButton(
                            icon = R.drawable.ic_person,
                            isSelected = selectedTab == ProfileTab.TAGGED,
                            onClick = { selectedTab = ProfileTab.TAGGED },
                            modifier = Modifier.weight(1f)
                        )
                    }
                }
            }
            
            // Posts grid
            item {
                Surface(
                    modifier = Modifier.fillMaxWidth(),
                    color = Color.White.copy(alpha = 0.2f)
                ) {
                    LazyVerticalGrid(
                        columns = GridCells.Fixed(3),
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(600.dp),
                        horizontalArrangement = Arrangement.spacedBy(2.dp),
                        verticalArrangement = Arrangement.spacedBy(2.dp)
                    ) {
                        items(posts) { post ->
                            PostGridItem(post = post)
                        }
                    }
                }
            }
        }
    }
    
    // Settings sheet
    if (showSettings) {
        // Settings modal would go here
    }
}

@Composable
fun StatButton(count: Int, label: String) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(4.dp),
        modifier = Modifier.clickable {}
    ) {
        Text(
            text = formatCount(count),
            fontSize = 18.sp,
            fontWeight = FontWeight.Bold,
            color = Color.White
        )
        
        Text(
            text = label,
            fontSize = 12.sp,
            color = Color.White.copy(alpha = 0.8f)
        )
    }
}

@Composable
fun HighlightItem(title: String) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        Surface(
            modifier = Modifier.size(64.dp),
            shape = CircleShape,
            color = Color.White.copy(alpha = 0.2f)
        ) {
            Box(
                modifier = Modifier.fillMaxSize(),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    painter = painterResource(R.drawable.ic_photo),
                    contentDescription = title,
                    tint = Color.White.copy(alpha = 0.5f),
                    modifier = Modifier.size(32.dp)
                )
            }
        }
        
        Text(
            text = title,
            fontSize = 12.sp,
            color = Color.White.copy(alpha = 0.9f)
        )
    }
}

@Composable
fun TabButton(
    icon: Int,
    isSelected: Boolean,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Box(
        modifier = modifier
            .fillMaxHeight()
            .clickable(onClick = onClick)
            .background(if (isSelected) Color.White.copy(alpha = 0.2f) else Color.Transparent),
        contentAlignment = Alignment.Center
    ) {
        Icon(
            painter = painterResource(icon),
            contentDescription = null,
            tint = if (isSelected) Color.White else Color.White.copy(alpha = 0.5f),
            modifier = Modifier.size(20.dp)
        )
    }
}

@Composable
fun PostGridItem(post: ProfilePost) {
    Box(
        modifier = Modifier
            .aspectRatio(1f)
            .clickable {}
    ) {
        if (post.thumbnailUrl != null) {
            AsyncImage(
                model = post.thumbnailUrl,
                contentDescription = null,
                modifier = Modifier.fillMaxSize(),
                contentScale = ContentScale.Crop
            )
        } else {
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .background(Color.White.copy(alpha = 0.1f)),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    painter = painterResource(R.drawable.ic_photo),
                    contentDescription = null,
                    tint = Color.White.copy(alpha = 0.3f),
                    modifier = Modifier.size(32.dp)
                )
            }
        }
        
        if (post.isVideo) {
            Icon(
                painter = painterResource(R.drawable.ic_play),
                contentDescription = "Video",
                tint = Color.White,
                modifier = Modifier
                    .size(14.dp)
                    .align(Alignment.TopEnd)
                    .padding(8.dp)
            )
        }
    }
}

// Models
data class VignetteProfile(
    val id: String,
    val username: String,
    val fullName: String,
    val bio: String?,
    val link: String?,
    val avatarUrl: String?,
    val profileImageUrl: String?,
    val postsCount: Int,
    val followersCount: Int,
    val followingCount: Int
)

data class ProfilePost(
    val id: String,
    val thumbnailUrl: String?,
    val isVideo: Boolean
)

enum class ProfileTab {
    POSTS, REELS, TAGGED
}

// ViewModel
class VignetteProfileViewModel : androidx.lifecycle.ViewModel() {
    private val _profile = kotlinx.coroutines.flow.MutableStateFlow(
        VignetteProfile(
            id = "user123",
            username = "yourusername",
            fullName = "Your Name",
            bio = "✨ Living life to the fullest\n📍 San Francisco, CA\n💼 Creator & Entrepreneur",
            link = "yourwebsite.com",
            avatarUrl = null,
            profileImageUrl = null,
            postsCount = 142,
            followersCount = 12500,
            followingCount = 890
        )
    )
    val profile: kotlinx.coroutines.flow.StateFlow<VignetteProfile> = _profile
    
    private val _posts = kotlinx.coroutines.flow.MutableStateFlow(
        (1..24).map { index ->
            ProfilePost(
                id = "post$index",
                thumbnailUrl = null,
                isVideo = index % 5 == 0
            )
        }
    )
    val posts: kotlinx.coroutines.flow.StateFlow<List<ProfilePost>> = _posts
}

fun formatCount(count: Int): String {
    return when {
        count >= 1_000_000 -> String.format("%.1fM", count / 1_000_000.0)
        count >= 1_000 -> String.format("%.1fK", count / 1_000.0)
        else -> count.toString()
    }
}
