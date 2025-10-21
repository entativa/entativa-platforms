package com.entativa.vignette.ui.home

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
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
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.entativa.vignette.R
import com.entativa.vignette.ui.theme.*

@Composable
fun VignetteHomeScreen(
    viewModel: VignetteHomeViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
) {
    var selectedTab by remember { mutableStateOf(VignetteTab.HOME) }
    var showingCamera by remember { mutableStateOf(false) }
    var showingSearch by remember { mutableStateOf(false) }
    
    Scaffold(
        topBar = {
            VignetteTopBar(
                onPlusClick = { showingCamera = true },
                onSearchClick = { showingSearch = true }
            )
        },
        bottomBar = {
            VignetteBottomNavBar(
                selectedTab = selectedTab,
                onTabSelected = { selectedTab = it }
            )
        },
        containerColor = Color.White
    ) { paddingValues ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            when (selectedTab) {
                VignetteTab.HOME -> VignetteFeedScreen(viewModel)
                VignetteTab.TAKES -> VignetteTakesScreen()
                VignetteTab.MESSAGES -> VignetteDirectScreen()
                VignetteTab.ACTIVITY -> VignetteActivityScreen()
                VignetteTab.PROFILE -> VignetteProfileScreen()
            }
        }
    }
}

// MARK: - Top Bar
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VignetteTopBar(
    onPlusClick: () -> Unit,
    onSearchClick: () -> Unit
) {
    TopAppBar(
        title = {
            Box(
                modifier = Modifier.fillMaxWidth(),
                contentAlignment = Alignment.Center
            ) {
                Text(
                    text = "Vignette",
                    fontSize = 32.sp,
                    fontStyle = FontStyle.Italic,
                    fontWeight = FontWeight.Normal,
                    color = vignette_text_primary
                )
            }
        },
        navigationIcon = {
            IconButton(onClick = onPlusClick) {
                Icon(
                    painter = painterResource(R.drawable.ic_plus_circle),
                    contentDescription = "Camera",
                    tint = vignette_text_primary,
                    modifier = Modifier.size(26.dp)
                )
            }
        },
        actions = {
            IconButton(onClick = onSearchClick) {
                Icon(
                    painter = painterResource(R.drawable.ic_search),
                    contentDescription = "Search",
                    tint = vignette_text_primary,
                    modifier = Modifier.size(24.dp)
                )
            }
        },
        colors = TopAppBarDefaults.topAppBarColors(
            containerColor = Color.White
        )
    )
}

// MARK: - Bottom Navigation Bar (Semi-translucent)
@Composable
fun VignetteBottomNavBar(
    selectedTab: VignetteTab,
    onTabSelected: (VignetteTab) -> Unit
) {
    Surface(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp, vertical = 8.dp),
        shape = RoundedCornerShape(24.dp),
        color = Color.White.copy(alpha = 0.92f),
        shadowElevation = 8.dp,
        tonalElevation = 0.dp
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .height(58.dp)
                .padding(horizontal = 8.dp),
            horizontalArrangement = Arrangement.SpaceEvenly,
            verticalAlignment = Alignment.CenterVertically
        ) {
            VignetteTab.values().forEach { tab ->
                VignetteBottomNavItem(
                    tab = tab,
                    isSelected = selectedTab == tab,
                    onClick = { onTabSelected(tab) }
                )
            }
        }
    }
}

@Composable
fun VignetteBottomNavItem(
    tab: VignetteTab,
    isSelected: Boolean,
    onClick: () -> Unit
) {
    IconButton(
        onClick = onClick,
        modifier = Modifier
            .weight(1f)
            .fillMaxHeight()
    ) {
        Icon(
            painter = painterResource(tab.icon),
            contentDescription = tab.label,
            tint = if (isSelected) vignette_text_primary else vignette_text_secondary,
            modifier = Modifier.size(26.dp)
        )
    }
}

// MARK: - Feed Screen (Instagram-style)
@Composable
fun VignetteFeedScreen(viewModel: VignetteHomeViewModel) {
    val posts by viewModel.posts.collectAsState()
    
    LazyColumn(
        modifier = Modifier.fillMaxSize()
    ) {
        // Stories Row (Circular Instagram-style)
        item {
            VignetteStoriesRow()
            Divider(
                modifier = Modifier.padding(vertical = 8.dp),
                color = vignette_separator
            )
        }
        
        // Posts (Instagram-style)
        items(posts) { post ->
            VignettePostCard(post = post)
            Divider(
                modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp),
                color = vignette_separator
            )
        }
    }
}

// MARK: - Stories Row (Circular)
@Composable
fun VignetteStoriesRow() {
    LazyRow(
        contentPadding = PaddingValues(horizontal = 16.dp, vertical = 8.dp),
        horizontalArrangement = Arrangement.spacedBy(16.dp)
    ) {
        // Your Story
        item {
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(6.dp)
            ) {
                Box(
                    modifier = Modifier.size(72.dp),
                    contentAlignment = Alignment.BottomEnd
                ) {
                    Box(
                        modifier = Modifier
                            .size(72.dp)
                            .clip(CircleShape)
                            .background(
                                Brush.linearGradient(
                                    colors = listOf(
                                        vignette_accent_moonstone.copy(alpha = 0.3f),
                                        vignette_accent_saffron.copy(alpha = 0.3f)
                                    )
                                )
                            )
                    )
                    
                    Box(
                        modifier = Modifier
                            .size(24.dp)
                            .clip(CircleShape)
                            .background(entativa_blue)
                            .padding(4.dp),
                        contentAlignment = Alignment.Center
                    ) {
                        Icon(
                            painter = painterResource(R.drawable.ic_plus),
                            contentDescription = "Add Story",
                            tint = Color.White,
                            modifier = Modifier.size(14.dp)
                        )
                    }
                }
                
                Text(
                    text = "Your story",
                    style = MaterialTheme.typography.labelSmall,
                    color = vignette_text_primary
                )
            }
        }
        
        // Friend Stories
        items(5) { index ->
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(6.dp)
            ) {
                Box(
                    modifier = Modifier
                        .size(72.dp)
                        .clip(CircleShape)
                        .background(
                            Brush.linearGradient(
                                colors = listOf(
                                    vignette_accent_moonstone,
                                    vignette_accent_saffron
                                ),
                                start = androidx.compose.ui.geometry.Offset(0f, 0f),
                                end = androidx.compose.ui.geometry.Offset(100f, 100f)
                            )
                        )
                        .padding(2.5.dp)
                        .background(Color.White, CircleShape)
                        .padding(2.dp)
                        .background(Color.Gray.copy(alpha = 0.3f), CircleShape)
                )
                
                Text(
                    text = "user_$index",
                    style = MaterialTheme.typography.labelSmall,
                    color = vignette_text_primary
                )
            }
        }
    }
}

// MARK: - Post Card (Instagram-style)
@Composable
fun VignettePostCard(post: VignettePost) {
    var isLiked by remember { mutableStateOf(false) }
    var isSaved by remember { mutableStateOf(false) }
    
    Column(
        modifier = Modifier.fillMaxWidth()
    ) {
        // Post Header
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 16.dp, vertical = 12.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Box(
                modifier = Modifier
                    .size(32.dp)
                    .clip(CircleShape)
                    .background(
                        Brush.linearGradient(
                            colors = listOf(
                                vignette_accent_moonstone,
                                vignette_accent_saffron
                            )
                        )
                    )
                    .padding(2.dp)
                    .background(Color.White, CircleShape)
                    .padding(1.dp)
                    .background(Color.Gray.copy(alpha = 0.3f), CircleShape)
            )
            
            Spacer(modifier = Modifier.width(12.dp))
            
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = post.userName,
                    style = MaterialTheme.typography.bodyMedium,
                    fontWeight = FontWeight.SemiBold,
                    color = vignette_text_primary
                )
                
                post.location?.let { location ->
                    Text(
                        text = location,
                        style = MaterialTheme.typography.bodySmall,
                        color = vignette_text_secondary
                    )
                }
            }
            
            IconButton(onClick = {}) {
                Icon(
                    painter = painterResource(R.drawable.ic_more),
                    contentDescription = "More",
                    tint = vignette_text_primary,
                    modifier = Modifier.size(16.dp)
                )
            }
        }
        
        // Post Image
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .aspectRatio(1f)
                .background(Color.Gray.copy(alpha = 0.2f)),
            contentAlignment = Alignment.Center
        ) {
            Icon(
                painter = painterResource(R.drawable.ic_photo),
                contentDescription = "Post Image",
                tint = Color.Gray.copy(alpha = 0.5f),
                modifier = Modifier.size(48.dp)
            )
        }
        
        // Action Buttons
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 16.dp, vertical = 12.dp),
            horizontalArrangement = Arrangement.spacedBy(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            IconButton(
                onClick = { isLiked = !isLiked },
                modifier = Modifier.size(32.dp)
            ) {
                Icon(
                    painter = painterResource(if (isLiked) R.drawable.ic_heart_filled else R.drawable.ic_heart),
                    contentDescription = "Like",
                    tint = if (isLiked) Color.Red else vignette_text_primary,
                    modifier = Modifier.size(26.dp)
                )
            }
            
            IconButton(
                onClick = {},
                modifier = Modifier.size(32.dp)
            ) {
                Icon(
                    painter = painterResource(R.drawable.ic_comment),
                    contentDescription = "Comment",
                    tint = vignette_text_primary,
                    modifier = Modifier.size(25.dp)
                )
            }
            
            IconButton(
                onClick = {},
                modifier = Modifier.size(32.dp)
            ) {
                Icon(
                    painter = painterResource(R.drawable.ic_share),
                    contentDescription = "Share",
                    tint = vignette_text_primary,
                    modifier = Modifier.size(24.dp)
                )
            }
            
            Spacer(modifier = Modifier.weight(1f))
            
            IconButton(
                onClick = { isSaved = !isSaved },
                modifier = Modifier.size(32.dp)
            ) {
                Icon(
                    painter = painterResource(if (isSaved) R.drawable.ic_bookmark_filled else R.drawable.ic_bookmark),
                    contentDescription = "Save",
                    tint = vignette_text_primary,
                    modifier = Modifier.size(24.dp)
                )
            }
        }
        
        // Likes Count
        if (post.likesCount > 0) {
            Text(
                text = "${post.likesCount} likes",
                style = MaterialTheme.typography.bodySmall,
                fontWeight = FontWeight.SemiBold,
                color = vignette_text_primary,
                modifier = Modifier.padding(horizontal = 16.dp, vertical = 4.dp)
            )
        }
        
        // Caption
        post.caption?.let { caption ->
            Row(
                modifier = Modifier.padding(horizontal = 16.dp, vertical = 4.dp)
            ) {
                Text(
                    text = post.userName,
                    style = MaterialTheme.typography.bodySmall,
                    fontWeight = FontWeight.SemiBold,
                    color = vignette_text_primary
                )
                Spacer(modifier = Modifier.width(4.dp))
                Text(
                    text = caption,
                    style = MaterialTheme.typography.bodySmall,
                    color = vignette_text_primary,
                    maxLines = 2
                )
            }
        }
        
        // View Comments
        if (post.commentsCount > 0) {
            Text(
                text = "View all ${post.commentsCount} comments",
                style = MaterialTheme.typography.bodySmall,
                color = vignette_text_secondary,
                modifier = Modifier.padding(horizontal = 16.dp, vertical = 4.dp)
            )
        }
        
        // Timestamp
        Text(
            text = post.timestamp,
            style = MaterialTheme.typography.labelSmall,
            color = vignette_text_secondary,
            modifier = Modifier.padding(horizontal = 16.dp, vertical = 4.dp)
        )
    }
}

// MARK: - Placeholder Screens
@Composable
fun VignetteTakesScreen() {
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(Color.Black),
        contentAlignment = Alignment.Center
    ) {
        Column(horizontalAlignment = Alignment.CenterHorizontally) {
            Text("Takes", style = MaterialTheme.typography.headlineLarge, color = Color.White)
            Text("Coming Soon", style = MaterialTheme.typography.bodyMedium, color = Color.White.copy(alpha = 0.7f))
        }
    }
}

@Composable
fun VignetteDirectScreen() {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Column(horizontalAlignment = Alignment.CenterHorizontally) {
            Text("Direct", style = MaterialTheme.typography.headlineLarge)
            Text("Coming Soon", style = MaterialTheme.typography.bodyMedium)
        }
    }
}

@Composable
fun VignetteActivityScreen() {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Column(horizontalAlignment = Alignment.CenterHorizontally) {
            Text("Activity", style = MaterialTheme.typography.headlineLarge)
            Text("Coming Soon", style = MaterialTheme.typography.bodyMedium)
        }
    }
}

@Composable
fun VignetteProfileScreen() {
    LazyColumn(
        modifier = Modifier.fillMaxSize(),
        contentPadding = PaddingValues(vertical = 20.dp)
    ) {
        item {
            Column(
                modifier = Modifier.fillMaxWidth(),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Box(
                    modifier = Modifier
                        .size(86.dp)
                        .clip(CircleShape)
                        .background(Color.Gray.copy(alpha = 0.3f))
                )
                
                Spacer(modifier = Modifier.height(16.dp))
                
                Text(
                    text = "@username",
                    style = MaterialTheme.typography.bodyLarge,
                    color = vignette_text_primary
                )
                
                Spacer(modifier = Modifier.height(20.dp))
                
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceEvenly
                ) {
                    ProfileStat("245", "Posts")
                    ProfileStat("1.2K", "Followers")
                    ProfileStat("892", "Following")
                }
                
                Spacer(modifier = Modifier.height(20.dp))
                
                OutlinedButton(
                    onClick = {},
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 16.dp)
                        .height(44.dp),
                    shape = RoundedCornerShape(8.dp),
                    colors = ButtonDefaults.outlinedButtonColors(
                        contentColor = entativa_blue
                    )
                ) {
                    Text("Edit Profile")
                }
                
                Spacer(modifier = Modifier.height(20.dp))
                
                Divider()
                
                Spacer(modifier = Modifier.height(20.dp))
                
                Text("Coming Soon", color = vignette_text_secondary)
            }
        }
    }
}

@Composable
fun ProfileStat(count: String, label: String) {
    Column(horizontalAlignment = Alignment.CenterHorizontally) {
        Text(
            text = count,
            style = MaterialTheme.typography.bodyLarge,
            fontWeight = FontWeight.SemiBold,
            color = vignette_text_primary
        )
        Text(
            text = label,
            style = MaterialTheme.typography.bodySmall,
            color = vignette_text_secondary
        )
    }
}

// MARK: - Supporting Types
enum class VignetteTab(val icon: Int, val label: String) {
    HOME(R.drawable.ic_home, "Home"),
    TAKES(R.drawable.ic_play_rect, "Takes"),
    MESSAGES(R.drawable.ic_message, "Direct"),
    ACTIVITY(R.drawable.ic_heart, "Activity"),
    PROFILE(R.drawable.ic_person, "Profile")
}

data class VignettePost(
    val id: String,
    val userName: String,
    val location: String?,
    val caption: String?,
    val mediaUrl: String,
    val likesCount: Int,
    val commentsCount: Int,
    val timestamp: String
)
