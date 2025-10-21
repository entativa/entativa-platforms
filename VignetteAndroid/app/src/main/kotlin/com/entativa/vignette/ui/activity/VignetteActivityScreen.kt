package com.entativa.vignette.ui.activity

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
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.SpanStyle
import androidx.compose.ui.text.buildAnnotatedString
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.withStyle
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import coil.compose.AsyncImage
import com.entativa.vignette.R

// MARK: - Vignette Activity Screen (Instagram-Style)
@Composable
fun VignetteActivityScreen(
    viewModel: VignetteActivityViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
) {
    val todayActivities by viewModel.todayActivities.collectAsState()
    val weekActivities by viewModel.weekActivities.collectAsState()
    val monthActivities by viewModel.monthActivities.collectAsState()
    val earlierActivities by viewModel.earlierActivities.collectAsState()
    val followingActivities by viewModel.followingActivities.collectAsState()
    
    var selectedTab by remember { mutableStateOf(ActivityTab.YOU) }
    
    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(Color.White)
    ) {
        // Top bar
        Surface(
            modifier = Modifier.fillMaxWidth(),
            tonalElevation = 2.dp
        ) {
            Column {
                // Title
                Text(
                    text = "Notifications",
                    fontSize = 20.sp,
                    fontWeight = FontWeight.Bold,
                    modifier = Modifier.padding(16.dp)
                )
                
                // Tab selector
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(44.dp)
                ) {
                    ActivityTabButton(
                        title = "Following",
                        isSelected = selectedTab == ActivityTab.FOLLOWING,
                        onClick = { selectedTab = ActivityTab.FOLLOWING },
                        modifier = Modifier.weight(1f)
                    )
                    
                    ActivityTabButton(
                        title = "You",
                        isSelected = selectedTab == ActivityTab.YOU,
                        onClick = { selectedTab = ActivityTab.YOU },
                        modifier = Modifier.weight(1f)
                    )
                }
            }
        }
        
        // Activity list
        LazyColumn(
            modifier = Modifier
                .fillMaxSize()
                .background(Color.White)
        ) {
            if (selectedTab == ActivityTab.YOU) {
                // Today
                if (todayActivities.isNotEmpty()) {
                    item {
                        ActivitySectionHeader(title = "Today")
                    }
                    items(todayActivities) { activity ->
                        VignetteActivityRow(activity = activity)
                    }
                }
                
                // This Week
                if (weekActivities.isNotEmpty()) {
                    item {
                        ActivitySectionHeader(title = "This Week")
                    }
                    items(weekActivities) { activity ->
                        VignetteActivityRow(activity = activity)
                    }
                }
                
                // This Month
                if (monthActivities.isNotEmpty()) {
                    item {
                        ActivitySectionHeader(title = "This Month")
                    }
                    items(monthActivities) { activity ->
                        VignetteActivityRow(activity = activity)
                    }
                }
                
                // Earlier
                if (earlierActivities.isNotEmpty()) {
                    item {
                        ActivitySectionHeader(title = "Earlier")
                    }
                    items(earlierActivities) { activity ->
                        VignetteActivityRow(activity = activity)
                    }
                }
            } else {
                // Following tab
                items(followingActivities) { activity ->
                    VignetteActivityRow(activity = activity)
                }
            }
        }
    }
}

@Composable
fun ActivityTabButton(
    title: String,
    isSelected: Boolean,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .fillMaxHeight()
            .clickable(onClick = onClick),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Text(
            text = title,
            fontSize = 16.sp,
            fontWeight = if (isSelected) FontWeight.SemiBold else FontWeight.Normal,
            color = if (isSelected) Color.Black else Color.Gray
        )
        
        Spacer(modifier = Modifier.height(8.dp))
        
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .height(1.dp)
                .background(if (isSelected) Color.Black else Color.Transparent)
        )
    }
}

@Composable
fun ActivitySectionHeader(title: String) {
    Text(
        text = title,
        fontSize = 15.sp,
        fontWeight = FontWeight.SemiBold,
        modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp)
    )
}

@Composable
fun VignetteActivityRow(activity: VignetteActivity) {
    var isFollowing by remember { mutableStateOf(activity.isFollowing) }
    
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable {}
            .background(if (!activity.isRead) Color(0xFFE3F2FD) else Color.White)
            .padding(horizontal = 16.dp, vertical = 8.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        // Profile picture
        Box(
            modifier = Modifier.size(44.dp)
        ) {
            if (activity.userAvatar != null) {
                AsyncImage(
                    model = activity.userAvatar,
                    contentDescription = null,
                    modifier = Modifier
                        .size(44.dp)
                        .clip(CircleShape),
                    contentScale = ContentScale.Crop
                )
            } else {
                Surface(
                    modifier = Modifier.size(44.dp),
                    shape = CircleShape,
                    color = Color.Gray.copy(alpha = 0.2f)
                ) {
                    Box(
                        modifier = Modifier.fillMaxSize(),
                        contentAlignment = Alignment.Center
                    ) {
                        Icon(
                            painter = painterResource(R.drawable.ic_person),
                            contentDescription = null,
                            tint = Color.Gray.copy(alpha = 0.5f),
                            modifier = Modifier.size(24.dp)
                        )
                    }
                }
            }
        }
        
        Spacer(modifier = Modifier.width(12.dp))
        
        // Activity text
        Column(
            modifier = Modifier.weight(1f)
        ) {
            Text(
                text = buildAnnotatedString {
                    withStyle(style = SpanStyle(fontWeight = FontWeight.SemiBold)) {
                        append(activity.username)
                    }
                    append(" ${activity.action}")
                },
                fontSize = 14.sp,
                lineHeight = 18.sp,
                maxLines = 3
            )
            
            Spacer(modifier = Modifier.height(4.dp))
            
            Text(
                text = activity.timeAgo,
                fontSize = 12.sp,
                color = Color.Gray
            )
        }
        
        Spacer(modifier = Modifier.width(12.dp))
        
        // Right side (thumbnail or follow button)
        when {
            activity.postThumbnail != null -> {
                AsyncImage(
                    model = activity.postThumbnail,
                    contentDescription = null,
                    modifier = Modifier
                        .size(44.dp)
                        .clip(RoundedCornerShape(4.dp)),
                    contentScale = ContentScale.Crop
                )
            }
            activity.type == VignetteActivityType.FOLLOW || 
            activity.type == VignetteActivityType.FOLLOW_REQUEST -> {
                Button(
                    onClick = { isFollowing = !isFollowing },
                    colors = ButtonDefaults.buttonColors(
                        containerColor = if (isFollowing) Color.Gray.copy(alpha = 0.2f) else Color(0xFF007CFC)
                    ),
                    shape = RoundedCornerShape(8.dp),
                    modifier = Modifier
                        .width(100.dp)
                        .height(32.dp),
                    contentPadding = PaddingValues(0.dp)
                ) {
                    Text(
                        text = if (isFollowing) "Following" else "Follow",
                        fontSize = 14.sp,
                        fontWeight = FontWeight.SemiBold,
                        color = if (isFollowing) Color.Black else Color.White
                    )
                }
            }
        }
    }
}

// Models
enum ActivityTab {
    FOLLOWING, YOU
}

enum VignetteActivityType {
    LIKE, COMMENT, FOLLOW, FOLLOW_REQUEST, MENTION, TAG, REPLY
}

data class VignetteActivity(
    val id: String,
    val type: VignetteActivityType,
    val username: String,
    val userAvatar: String?,
    val action: String,
    val postThumbnail: String?,
    val timeAgo: String,
    val isRead: Boolean,
    val isFollowing: Boolean,
    val timestamp: Long
)

// ViewModel
class VignetteActivityViewModel : androidx.lifecycle.ViewModel() {
    private val _todayActivities = kotlinx.coroutines.flow.MutableStateFlow<List<VignetteActivity>>(emptyList())
    val todayActivities: kotlinx.coroutines.flow.StateFlow<List<VignetteActivity>> = _todayActivities
    
    private val _weekActivities = kotlinx.coroutines.flow.MutableStateFlow<List<VignetteActivity>>(emptyList())
    val weekActivities: kotlinx.coroutines.flow.StateFlow<List<VignetteActivity>> = _weekActivities
    
    private val _monthActivities = kotlinx.coroutines.flow.MutableStateFlow<List<VignetteActivity>>(emptyList())
    val monthActivities: kotlinx.coroutines.flow.StateFlow<List<VignetteActivity>> = _monthActivities
    
    private val _earlierActivities = kotlinx.coroutines.flow.MutableStateFlow<List<VignetteActivity>>(emptyList())
    val earlierActivities: kotlinx.coroutines.flow.StateFlow<List<VignetteActivity>> = _earlierActivities
    
    private val _followingActivities = kotlinx.coroutines.flow.MutableStateFlow<List<VignetteActivity>>(emptyList())
    val followingActivities: kotlinx.coroutines.flow.StateFlow<List<VignetteActivity>> = _followingActivities
    
    init {
        loadMockData()
    }
    
    private fun loadMockData() {
        val now = System.currentTimeMillis()
        
        _todayActivities.value = listOf(
            VignetteActivity(
                id = "1",
                type = VignetteActivityType.LIKE,
                username = "sarah_jones",
                userAvatar = null,
                action = "liked your photo.",
                postThumbnail = null,
                timeAgo = "2h",
                isRead = false,
                isFollowing = true,
                timestamp = now - 2 * 3600 * 1000
            ),
            VignetteActivity(
                id = "2",
                type = VignetteActivityType.COMMENT,
                username = "mike_wilson",
                userAvatar = null,
                action = "commented: \"Amazing shot! 🔥\"",
                postThumbnail = null,
                timeAgo = "4h",
                isRead = false,
                isFollowing = true,
                timestamp = now - 4 * 3600 * 1000
            ),
            VignetteActivity(
                id = "3",
                type = VignetteActivityType.FOLLOW,
                username = "alex_creative",
                userAvatar = null,
                action = "started following you.",
                postThumbnail = null,
                timeAgo = "5h",
                isRead = true,
                isFollowing = false,
                timestamp = now - 5 * 3600 * 1000
            )
        )
        
        _weekActivities.value = listOf(
            VignetteActivity(
                id = "4",
                type = VignetteActivityType.LIKE,
                username = "emma_davis",
                userAvatar = null,
                action = "liked your photo.",
                postThumbnail = null,
                timeAgo = "2d",
                isRead = true,
                isFollowing = true,
                timestamp = now - 2 * 86400 * 1000
            )
        )
        
        _monthActivities.value = listOf(
            VignetteActivity(
                id = "5",
                type = VignetteActivityType.FOLLOW,
                username = "photographer_pro",
                userAvatar = null,
                action = "started following you.",
                postThumbnail = null,
                timeAgo = "1w",
                isRead = true,
                isFollowing = false,
                timestamp = now - 7 * 86400 * 1000
            )
        )
        
        _earlierActivities.value = listOf(
            VignetteActivity(
                id = "6",
                type = VignetteActivityType.COMMENT,
                username = "design_lover",
                userAvatar = null,
                action = "commented: \"Love your style!\"",
                postThumbnail = null,
                timeAgo = "3w",
                isRead = true,
                isFollowing = true,
                timestamp = now - 21 * 86400 * 1000
            )
        )
        
        _followingActivities.value = listOf(
            VignetteActivity(
                id = "7",
                type = VignetteActivityType.LIKE,
                username = "sarah_jones",
                userAvatar = null,
                action = "liked a photo by mike_wilson.",
                postThumbnail = null,
                timeAgo = "1h",
                isRead = true,
                isFollowing = true,
                timestamp = now - 1 * 3600 * 1000
            )
        )
    }
}
