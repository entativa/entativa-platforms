package com.entativa.ui.home

import androidx.compose.animation.core.Spring
import androidx.compose.animation.core.spring
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
import androidx.compose.ui.draw.blur
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.entativa.R
import com.entativa.ui.theme.*

@Composable
fun EntativaHomeScreen(
    viewModel: EntativaHomeViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
) {
    var selectedTab by remember { mutableStateOf(HomeTab.HOME) }
    var showingCreatePost by remember { mutableStateOf(false) }
    var showingSearch by remember { mutableStateOf(false) }
    
    Scaffold(
        topBar = {
            EntativaTopBar(
                onPlusClick = { showingCreatePost = true },
                onSearchClick = { showingSearch = true }
            )
        },
        bottomBar = {
            EntativaBottomNavBar(
                selectedTab = selectedTab,
                onTabSelected = { selectedTab = it }
            )
        },
        containerColor = Color(0xFFF5F5F5)
    ) { paddingValues ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            when (selectedTab) {
                HomeTab.HOME -> EntativaFeedScreen(viewModel)
                HomeTab.TAKES -> EntativaTakesScreen()
                HomeTab.MESSAGES -> EntativaMessagesScreen()
                HomeTab.ACTIVITY -> EntativaActivityScreen()
                HomeTab.MENU -> EntativaMenuScreen()
            }
        }
    }
}

// MARK: - Top Bar
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun EntativaTopBar(
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
                    text = "entativa",
                    fontSize = 28.sp,
                    fontWeight = FontWeight.Bold,
                    fontStyle = FontStyle.Italic,
                    style = MaterialTheme.typography.displaySmall.copy(
                        brush = Brush.horizontalGradient(
                            colors = listOf(
                                entativa_primary_blue,
                                entativa_primary_purple,
                                entativa_primary_pink
                            )
                        )
                    )
                )
            }
        },
        navigationIcon = {
            IconButton(onClick = onPlusClick) {
                Icon(
                    painter = painterResource(R.drawable.ic_plus_circle),
                    contentDescription = "Create Post",
                    tint = entativa_primary_blue,
                    modifier = Modifier.size(28.dp)
                )
            }
        },
        actions = {
            IconButton(onClick = onSearchClick) {
                Icon(
                    painter = painterResource(R.drawable.ic_search),
                    contentDescription = "Search",
                    tint = entativa_text_primary,
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
fun EntativaBottomNavBar(
    selectedTab: HomeTab,
    onTabSelected: (HomeTab) -> Unit
) {
    Surface(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp, vertical = 8.dp),
        shape = RoundedCornerShape(24.dp),
        color = Color.White.copy(alpha = 0.95f),
        shadowElevation = 8.dp,
        tonalElevation = 0.dp
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .height(60.dp)
                .padding(horizontal = 8.dp),
            horizontalArrangement = Arrangement.SpaceEvenly,
            verticalAlignment = Alignment.CenterVertically
        ) {
            HomeTab.values().forEach { tab ->
                BottomNavItem(
                    tab = tab,
                    isSelected = selectedTab == tab,
                    onClick = { onTabSelected(tab) }
                )
            }
        }
    }
}

@Composable
fun BottomNavItem(
    tab: HomeTab,
    isSelected: Boolean,
    onClick: () -> Unit
) {
    Column(
        modifier = Modifier
            .weight(1f)
            .fillMaxHeight()
            .padding(vertical = 8.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        IconButton(
            onClick = onClick,
            modifier = Modifier.size(36.dp)
        ) {
            Icon(
                painter = painterResource(tab.icon),
                contentDescription = tab.label,
                tint = if (isSelected) entativa_primary_blue else entativa_text_secondary,
                modifier = Modifier.size(26.dp)
            )
        }
        
        if (isSelected) {
            Box(
                modifier = Modifier
                    .size(4.dp)
                    .clip(CircleShape)
                    .background(entativa_primary_blue)
            )
        }
    }
}

// MARK: - Feed Screen
@Composable
fun EntativaFeedScreen(viewModel: EntativaHomeViewModel) {
    val posts by viewModel.posts.collectAsState()
    
    LazyColumn(
        modifier = Modifier.fillMaxSize(),
        contentPadding = PaddingValues(bottom = 16.dp)
    ) {
        // Stories Row (Card Style)
        item {
            EntativaStoriesRow()
            Spacer(modifier = Modifier.height(12.dp))
        }
        
        // Posts (Carousel Style)
        items(posts) { post ->
            EntativaCarouselPostCard(post = post)
            Spacer(modifier = Modifier.height(12.dp))
        }
    }
}

// MARK: - Stories Row (Card Style)
@Composable
fun EntativaStoriesRow() {
    LazyRow(
        contentPadding = PaddingValues(horizontal = 16.dp, vertical = 8.dp),
        horizontalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        // Create Story Card
        item {
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                Box(
                    modifier = Modifier
                        .width(110.dp)
                        .height(160.dp)
                        .clip(RoundedCornerShape(16.dp))
                        .background(
                            Brush.linearGradient(
                                colors = listOf(
                                    entativa_primary_blue.copy(alpha = 0.3f),
                                    entativa_primary_purple.copy(alpha = 0.3f)
                                )
                            )
                        ),
                    contentAlignment = Alignment.BottomEnd
                ) {
                    Icon(
                        painter = painterResource(R.drawable.ic_plus_circle),
                        contentDescription = "Create Story",
                        tint = entativa_primary_blue,
                        modifier = Modifier
                            .padding(8.dp)
                            .size(32.dp)
                            .background(Color.White, CircleShape)
                            .padding(4.dp)
                    )
                }
                
                Text(
                    text = "Create",
                    style = MaterialTheme.typography.labelSmall,
                    color = entativa_text_primary
                )
            }
        }
        
        // Story Cards
        items(5) { index ->
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                Box(
                    modifier = Modifier
                        .width(110.dp)
                        .height(160.dp)
                        .clip(RoundedCornerShape(16.dp))
                        .background(
                            Brush.linearGradient(
                                colors = listOf(
                                    Color.Blue.copy(alpha = 0.5f),
                                    Color.Purple.copy(alpha = 0.5f)
                                )
                            )
                        ),
                    contentAlignment = Alignment.TopStart
                ) {
                    Box(
                        modifier = Modifier
                            .padding(8.dp)
                            .size(40.dp)
                            .clip(CircleShape)
                            .background(
                                Brush.linearGradient(
                                    colors = listOf(
                                        entativa_primary_blue,
                                        entativa_primary_purple
                                    )
                                )
                            )
                            .padding(3.dp)
                            .background(Color.White, CircleShape)
                            .padding(1.dp)
                            .background(Color.Gray, CircleShape)
                    )
                }
                
                Text(
                    text = "User $index",
                    style = MaterialTheme.typography.labelSmall,
                    color = entativa_text_primary
                )
            }
        }
    }
}

// MARK: - Carousel Post Card
@Composable
fun EntativaCarouselPostCard(post: Post) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp),
        shape = RoundedCornerShape(20.dp),
        colors = CardDefaults.cardColors(
            containerColor = Color.White
        ),
        elevation = CardDefaults.cardElevation(
            defaultElevation = 2.dp
        )
    ) {
        Column(
            modifier = Modifier.padding(16.dp)
        ) {
            // Post Header
            Row(
                modifier = Modifier.fillMaxWidth(),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Box(
                    modifier = Modifier
                        .size(44.dp)
                        .clip(CircleShape)
                        .background(
                            Brush.linearGradient(
                                colors = listOf(
                                    entativa_primary_blue,
                                    entativa_primary_purple
                                )
                            )
                        )
                )
                
                Spacer(modifier = Modifier.width(12.dp))
                
                Column(modifier = Modifier.weight(1f)) {
                    Text(
                        text = post.userName,
                        style = MaterialTheme.typography.bodyMedium,
                        fontWeight = FontWeight.SemiBold,
                        color = entativa_text_primary
                    )
                    Text(
                        text = post.timestamp,
                        style = MaterialTheme.typography.bodySmall,
                        color = entativa_text_secondary
                    )
                }
                
                IconButton(onClick = {}) {
                    Icon(
                        painter = painterResource(R.drawable.ic_more),
                        contentDescription = "More",
                        tint = entativa_text_secondary
                    )
                }
            }
            
            // Post Text
            if (post.text.isNotEmpty()) {
                Spacer(modifier = Modifier.height(12.dp))
                Text(
                    text = post.text,
                    style = MaterialTheme.typography.bodyMedium,
                    color = entativa_text_primary
                )
            }
            
            // Media (if any)
            if (post.mediaUrls.isNotEmpty()) {
                Spacer(modifier = Modifier.height(12.dp))
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(300.dp)
                        .clip(RoundedCornerShape(16.dp))
                        .background(Color.Gray.copy(alpha = 0.2f)),
                    contentAlignment = Alignment.Center
                ) {
                    Icon(
                        painter = painterResource(R.drawable.ic_photo),
                        contentDescription = "Post Media",
                        tint = Color.Gray.copy(alpha = 0.5f),
                        modifier = Modifier.size(48.dp)
                    )
                }
            }
            
            // Action Buttons
            Spacer(modifier = Modifier.height(12.dp))
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.spacedBy(20.dp)
            ) {
                PostActionButton(
                    icon = R.drawable.ic_heart,
                    count = post.likesCount,
                    tint = entativa_text_secondary
                )
                PostActionButton(
                    icon = R.drawable.ic_comment,
                    count = post.commentsCount,
                    tint = entativa_text_secondary
                )
                PostActionButton(
                    icon = R.drawable.ic_share,
                    count = post.sharesCount,
                    tint = entativa_text_secondary
                )
                
                Spacer(modifier = Modifier.weight(1f))
                
                IconButton(onClick = {}) {
                    Icon(
                        painter = painterResource(R.drawable.ic_bookmark),
                        contentDescription = "Save",
                        tint = entativa_text_secondary
                    )
                }
            }
        }
    }
}

@Composable
fun PostActionButton(
    icon: Int,
    count: Int,
    tint: Color
) {
    Row(
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(6.dp)
    ) {
        IconButton(onClick = {}, modifier = Modifier.size(36.dp)) {
            Icon(
                painter = painterResource(icon),
                contentDescription = null,
                tint = tint,
                modifier = Modifier.size(20.dp)
            )
        }
        if (count > 0) {
            Text(
                text = count.toString(),
                style = MaterialTheme.typography.bodySmall,
                color = tint
            )
        }
    }
}

// MARK: - Placeholder Screens
@Composable
fun EntativaTakesScreen() {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Column(horizontalAlignment = Alignment.CenterHorizontally) {
            Text("Takes", style = MaterialTheme.typography.headlineLarge)
            Text("Coming Soon", style = MaterialTheme.typography.bodyMedium)
        }
    }
}

@Composable
fun EntativaMessagesScreen() {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Column(horizontalAlignment = Alignment.CenterHorizontally) {
            Text("Messages", style = MaterialTheme.typography.headlineLarge)
            Text("Coming Soon", style = MaterialTheme.typography.bodyMedium)
        }
    }
}

@Composable
fun EntativaActivityScreen() {
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
fun EntativaMenuScreen() {
    LazyColumn(
        modifier = Modifier.fillMaxSize(),
        contentPadding = PaddingValues(20.dp)
    ) {
        item {
            Text(
                "Menu",
                style = MaterialTheme.typography.headlineLarge,
                modifier = Modifier.padding(bottom = 20.dp)
            )
        }
        
        item {
            MenuSection(
                title = "Your Shortcuts",
                items = listOf(
                    MenuItem("Friends", R.drawable.ic_people),
                    MenuItem("Groups", R.drawable.ic_group),
                    MenuItem("Watch", R.drawable.ic_play),
                    MenuItem("Marketplace", R.drawable.ic_shop)
                )
            )
            
            Spacer(modifier = Modifier.height(20.dp))
        }
        
        item {
            MenuSection(
                title = "Settings & Privacy",
                items = listOf(
                    MenuItem("Settings", R.drawable.ic_settings),
                    MenuItem("Privacy Center", R.drawable.ic_shield),
                    MenuItem("About", R.drawable.ic_info)
                )
            )
        }
    }
}

@Composable
fun MenuSection(
    title: String,
    items: List<MenuItem>
) {
    Column {
        Text(
            text = title,
            style = MaterialTheme.typography.labelLarge,
            color = entativa_text_secondary,
            modifier = Modifier.padding(bottom = 12.dp)
        )
        
        items.forEach { item ->
            Surface(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(vertical = 4.dp),
                color = Color(0xFFF5F5F5),
                shape = RoundedCornerShape(12.dp)
            ) {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(16.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        painter = painterResource(item.icon),
                        contentDescription = item.title,
                        tint = entativa_primary_blue,
                        modifier = Modifier.size(24.dp)
                    )
                    
                    Spacer(modifier = Modifier.width(16.dp))
                    
                    Text(
                        text = item.title,
                        style = MaterialTheme.typography.bodyMedium,
                        color = entativa_text_primary,
                        modifier = Modifier.weight(1f)
                    )
                    
                    Icon(
                        painter = painterResource(R.drawable.ic_chevron_right),
                        contentDescription = null,
                        tint = entativa_text_secondary,
                        modifier = Modifier.size(16.dp)
                    )
                }
            }
        }
    }
}

// MARK: - Supporting Types
enum class HomeTab(val icon: Int, val label: String) {
    HOME(R.drawable.ic_home, "Home"),
    TAKES(R.drawable.ic_play_rect, "Takes"),
    MESSAGES(R.drawable.ic_message, "Messages"),
    ACTIVITY(R.drawable.ic_bell, "Activity"),
    MENU(R.drawable.ic_menu, "Menu")
}

data class MenuItem(
    val title: String,
    val icon: Int
)

data class Post(
    val id: String,
    val userName: String,
    val timestamp: String,
    val text: String,
    val mediaUrls: List<String>,
    val likesCount: Int,
    val commentsCount: Int,
    val sharesCount: Int
)
