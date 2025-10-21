package com.entativa.ui.activity

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
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import coil.compose.AsyncImage
import com.entativa.R

// MARK: - Entativa Activity Screen (Facebook-Style)
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun EntativaActivityScreen(
    viewModel: EntativaActivityViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
) {
    val newActivities by viewModel.newActivities.collectAsState()
    val earlierActivities by viewModel.earlierActivities.collectAsState()
    
    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        "Notifications",
                        fontSize = 20.sp,
                        fontWeight = FontWeight.Bold
                    )
                },
                actions = {
                    IconButton(onClick = {}) {
                        Icon(
                            painter = painterResource(R.drawable.ic_more),
                            contentDescription = "More"
                        )
                    }
                }
            )
        }
    ) { paddingValues ->
        LazyColumn(
            modifier = Modifier
                .fillMaxSize()
                .background(Color(0xFFF0F2F5))
                .padding(paddingValues)
        ) {
            // New notifications
            if (newActivities.isNotEmpty()) {
                item {
                    EntativaActivitySectionHeader(title = "New")
                }
                items(newActivities) { activity ->
                    EntativaActivityRow(activity = activity)
                }
                item {
                    Spacer(modifier = Modifier.height(8.dp))
                }
            }
            
            // Earlier notifications
            if (earlierActivities.isNotEmpty()) {
                item {
                    EntativaActivitySectionHeader(title = "Earlier")
                }
                items(earlierActivities) { activity ->
                    EntativaActivityRow(activity = activity)
                }
            }
        }
    }
}

@Composable
fun EntativaActivitySectionHeader(title: String) {
    Surface(
        modifier = Modifier.fillMaxWidth(),
        color = Color(0xFFF0F2F5)
    ) {
        Text(
            text = title,
            fontSize = 17.sp,
            fontWeight = FontWeight.SemiBold,
            modifier = Modifier.padding(horizontal = 16.dp, vertical = 12.dp)
        )
    }
}

@Composable
fun EntativaActivityRow(activity: EntativaActivityData) {
    Surface(
        modifier = Modifier.fillMaxWidth(),
        color = if (!activity.isRead) Color(0xFFE3F2FD) else Color.White
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .clickable {}
                .padding(16.dp),
            verticalAlignment = Alignment.Top
        ) {
            // Icon with colored background
            Box(
                modifier = Modifier.size(56.dp),
                contentAlignment = Alignment.Center
            ) {
                Surface(
                    modifier = Modifier.size(56.dp),
                    shape = CircleShape,
                    color = activity.iconBackgroundColor
                ) {
                    Box(
                        modifier = Modifier.fillMaxSize(),
                        contentAlignment = Alignment.Center
                    ) {
                        if (activity.userAvatar != null) {
                            AsyncImage(
                                model = activity.userAvatar,
                                contentDescription = null,
                                modifier = Modifier
                                    .size(52.dp)
                                    .clip(CircleShape),
                                contentScale = ContentScale.Crop
                            )
                        } else {
                            Icon(
                                painter = painterResource(activity.iconResId),
                                contentDescription = null,
                                tint = Color.White,
                                modifier = Modifier.size(24.dp)
                            )
                        }
                    }
                }
                
                // Badge for specific types
                if (activity.showBadge) {
                    Surface(
                        modifier = Modifier
                            .size(16.dp)
                            .align(Alignment.TopEnd),
                        shape = CircleShape,
                        color = Color.Red
                    ) {
                        Box(
                            modifier = Modifier.fillMaxSize(),
                            contentAlignment = Alignment.Center
                        ) {
                            Icon(
                                painter = painterResource(R.drawable.ic_info),
                                contentDescription = null,
                                tint = Color.White,
                                modifier = Modifier.size(10.dp)
                            )
                        }
                    }
                }
            }
            
            Spacer(modifier = Modifier.width(12.dp))
            
            // Content
            Column(
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = activity.text,
                    fontSize = 15.sp,
                    lineHeight = 20.sp,
                    maxLines = 3
                )
                
                Spacer(modifier = Modifier.height(4.dp))
                
                Text(
                    text = activity.timeAgo,
                    fontSize = 13.sp,
                    color = Color.Gray
                )
                
                // Action buttons for specific types
                if (activity.type == EntativaActivityType.FRIEND_REQUEST) {
                    Spacer(modifier = Modifier.height(8.dp))
                    
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
                            Text("Confirm", fontSize = 15.sp, fontWeight = FontWeight.SemiBold)
                        }
                        
                        OutlinedButton(
                            onClick = {},
                            modifier = Modifier.weight(1f),
                            shape = RoundedCornerShape(8.dp),
                            colors = ButtonDefaults.outlinedButtonColors(
                                contentColor = Color.Black
                            )
                        ) {
                            Text("Delete", fontSize = 15.sp, fontWeight = FontWeight.SemiBold)
                        }
                    }
                }
            }
            
            Spacer(modifier = Modifier.width(12.dp))
            
            // Post thumbnail if applicable
            if (activity.postThumbnail != null) {
                AsyncImage(
                    model = activity.postThumbnail,
                    contentDescription = null,
                    modifier = Modifier
                        .size(64.dp)
                        .clip(RoundedCornerShape(8.dp)),
                    contentScale = ContentScale.Crop
                )
            }
            
            // Menu button
            IconButton(onClick = {}) {
                Icon(
                    painter = painterResource(R.drawable.ic_more),
                    contentDescription = "More",
                    tint = Color.Gray,
                    modifier = Modifier.size(16.dp)
                )
            }
        }
    }
}

// Models
enum class EntativaActivityType {
    LIKE, COMMENT, SHARE, FRIEND_REQUEST, FRIEND_ACCEPTED,
    TAG, MENTION, EVENT, BIRTHDAY, MEMORY
}

data class EntativaActivityData(
    val id: String,
    val type: EntativaActivityType,
    val text: String,
    val timeAgo: String,
    val userAvatar: String?,
    val postThumbnail: String?,
    val isRead: Boolean,
    val showBadge: Boolean,
    val timestamp: Long
) {
    val iconResId: Int
        get() = when (type) {
            EntativaActivityType.LIKE -> R.drawable.ic_heart
            EntativaActivityType.COMMENT -> R.drawable.ic_comment
            EntativaActivityType.SHARE -> R.drawable.ic_share
            EntativaActivityType.FRIEND_REQUEST -> R.drawable.ic_person
            EntativaActivityType.FRIEND_ACCEPTED -> R.drawable.ic_people
            EntativaActivityType.TAG -> R.drawable.ic_flag
            EntativaActivityType.MENTION -> R.drawable.ic_at
            EntativaActivityType.EVENT -> R.drawable.ic_calendar
            EntativaActivityType.BIRTHDAY -> R.drawable.ic_gift
            EntativaActivityType.MEMORY -> R.drawable.ic_clock
        }
    
    val iconBackgroundColor: Color
        get() = when (type) {
            EntativaActivityType.LIKE -> Color.Red
            EntativaActivityType.COMMENT -> Color(0xFF007CFC)
            EntativaActivityType.SHARE -> Color.Green
            EntativaActivityType.FRIEND_REQUEST -> Color(0xFF007CFC)
            EntativaActivityType.FRIEND_ACCEPTED -> Color(0xFF007CFC)
            EntativaActivityType.TAG -> Color(0xFFFF9800)
            EntativaActivityType.MENTION -> Color(0xFF6F3EFB)
            EntativaActivityType.EVENT -> Color.Red
            EntativaActivityType.BIRTHDAY -> Color(0xFFE91E63)
            EntativaActivityType.MEMORY -> Color(0xFF6F3EFB)
        }
}

// ViewModel
class EntativaActivityViewModel : androidx.lifecycle.ViewModel() {
    private val _newActivities = kotlinx.coroutines.flow.MutableStateFlow<List<EntativaActivityData>>(emptyList())
    val newActivities: kotlinx.coroutines.flow.StateFlow<List<EntativaActivityData>> = _newActivities
    
    private val _earlierActivities = kotlinx.coroutines.flow.MutableStateFlow<List<EntativaActivityData>>(emptyList())
    val earlierActivities: kotlinx.coroutines.flow.StateFlow<List<EntativaActivityData>> = _earlierActivities
    
    init {
        loadMockData()
    }
    
    private fun loadMockData() {
        val now = System.currentTimeMillis()
        
        _newActivities.value = listOf(
            EntativaActivityData(
                id = "1",
                type = EntativaActivityType.FRIEND_REQUEST,
                text = "Sarah Johnson sent you a friend request.",
                timeAgo = "2 hours ago",
                userAvatar = null,
                postThumbnail = null,
                isRead = false,
                showBadge = true,
                timestamp = now - 2 * 3600 * 1000
            ),
            EntativaActivityData(
                id = "2",
                type = EntativaActivityType.LIKE,
                text = "Mike Wilson and 12 others reacted to your post.",
                timeAgo = "4 hours ago",
                userAvatar = null,
                postThumbnail = null,
                isRead = false,
                showBadge = false,
                timestamp = now - 4 * 3600 * 1000
            ),
            EntativaActivityData(
                id = "3",
                type = EntativaActivityType.COMMENT,
                text = "Emma Davis commented on your photo: \"This is amazing! 🔥\"",
                timeAgo = "6 hours ago",
                userAvatar = null,
                postThumbnail = null,
                isRead = false,
                showBadge = false,
                timestamp = now - 6 * 3600 * 1000
            ),
            EntativaActivityData(
                id = "4",
                type = EntativaActivityType.BIRTHDAY,
                text = "It's Alex Chen's birthday today! Write on their timeline.",
                timeAgo = "8 hours ago",
                userAvatar = null,
                postThumbnail = null,
                isRead = true,
                showBadge = false,
                timestamp = now - 8 * 3600 * 1000
            )
        )
        
        _earlierActivities.value = listOf(
            EntativaActivityData(
                id = "5",
                type = EntativaActivityType.FRIEND_ACCEPTED,
                text = "Chris Taylor accepted your friend request.",
                timeAgo = "Yesterday",
                userAvatar = null,
                postThumbnail = null,
                isRead = true,
                showBadge = false,
                timestamp = now - 24 * 3600 * 1000
            ),
            EntativaActivityData(
                id = "6",
                type = EntativaActivityType.TAG,
                text = "You were tagged in a photo by Jessica Brown.",
                timeAgo = "2 days ago",
                userAvatar = null,
                postThumbnail = null,
                isRead = true,
                showBadge = false,
                timestamp = now - 2 * 86400 * 1000
            ),
            EntativaActivityData(
                id = "7",
                type = EntativaActivityType.MEMORY,
                text = "We found a memory from 3 years ago that you might like.",
                timeAgo = "3 days ago",
                userAvatar = null,
                postThumbnail = null,
                isRead = true,
                showBadge = false,
                timestamp = now - 3 * 86400 * 1000
            )
        )
    }
}
