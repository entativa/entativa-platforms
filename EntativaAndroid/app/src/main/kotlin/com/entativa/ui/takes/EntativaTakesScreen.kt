package com.entativa.ui.takes

import androidx.compose.foundation.background
import androidx.compose.foundation.gestures.detectVerticalDragGestures
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.input.pointer.pointerInput
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.entativa.R
import com.entativa.ui.theme.*
import com.google.accompanist.pager.ExperimentalPagerApi
import com.google.accompanist.pager.VerticalPager
import com.google.accompanist.pager.rememberPagerState

@OptIn(ExperimentalPagerApi::class)
@Composable
fun EntativaTakesScreen(
    viewModel: EntativaTakesViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
) {
    val takes by viewModel.takes.collectAsState()
    val pagerState = rememberPagerState()
    
    Box(modifier = Modifier.fillMaxSize().background(Color.Black)) {
        // Vertical Pager for Takes
        VerticalPager(
            count = takes.size,
            state = pagerState,
            modifier = Modifier.fillMaxSize()
        ) { page ->
            TakeVideoCard(
                take = takes[page],
                isCurrentlyPlaying = page == pagerState.currentPage,
                viewModel = viewModel
            )
            
            // Load more when near end
            LaunchedEffect(page) {
                if (page >= takes.size - 2) {
                    viewModel.loadMore()
                }
            }
        }
        
        // Top Bar
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .statusBarsPadding()
                .padding(16.dp)
        ) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(
                    text = "Takes",
                    fontSize = 20.sp,
                    fontWeight = FontWeight.SemiBold,
                    color = Color.White
                )
                
                IconButton(onClick = {}) {
                    Icon(
                        painter = painterResource(R.drawable.ic_camera),
                        contentDescription = "Camera",
                        tint = Color.White,
                        modifier = Modifier.size(24.dp)
                    )
                }
            }
        }
    }
}

@Composable
fun TakeVideoCard(
    take: TakeData,
    isCurrentlyPlaying: Boolean,
    viewModel: EntativaTakesViewModel
) {
    var showComments by remember { mutableStateOf(false) }
    var showShare by remember { mutableStateOf(false) }
    var isFollowing by remember { mutableStateOf(false) }
    
    Box(modifier = Modifier.fillMaxSize()) {
        // Video Player (Real ExoPlayer)
        VideoPlayer(
            videoUrl = take.videoUrl,
            isPlaying = isCurrentlyPlaying,
            modifier = Modifier.fillMaxSize()
        )
        
        // Right Side Actions
        Column(
            modifier = Modifier
                .align(Alignment.CenterEnd)
                .padding(end = 16.dp)
                .padding(bottom = 100.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(24.dp)
        ) {
            // Profile Avatar
            Box(contentAlignment = Alignment.BottomCenter) {
                Box(
                    modifier = Modifier
                        .size(48.dp)
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
                
                if (!isFollowing) {
                    Box(
                        modifier = Modifier
                            .offset(y = 8.dp)
                            .size(20.dp)
                            .clip(CircleShape)
                            .background(entativa_primary_blue),
                        contentAlignment = Alignment.Center
                    ) {
                        Icon(
                            painter = painterResource(R.drawable.ic_plus),
                            contentDescription = "Follow",
                            tint = Color.White,
                            modifier = Modifier.size(12.dp)
                        )
                    }
                }
            }
            
            // Like
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(4.dp)
            ) {
                IconButton(
                    onClick = {
                        if (take.isLiked) {
                            viewModel.unlikeTake(take.id)
                        } else {
                            viewModel.likeTake(take.id)
                        }
                    },
                    modifier = Modifier.size(40.dp)
                ) {
                    Icon(
                        painter = painterResource(if (take.isLiked) R.drawable.ic_heart_filled else R.drawable.ic_heart),
                        contentDescription = "Like",
                        tint = if (take.isLiked) Color.Red else Color.White,
                        modifier = Modifier.size(32.dp)
                    )
                }
                Text(
                    text = formatCount(take.likesCount),
                    fontSize = 12.sp,
                    fontWeight = FontWeight.SemiBold,
                    color = Color.White
                )
            }
            
            // Comments
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(4.dp)
            ) {
                IconButton(
                    onClick = { showComments = true },
                    modifier = Modifier.size(40.dp)
                ) {
                    Icon(
                        painter = painterResource(R.drawable.ic_comment),
                        contentDescription = "Comments",
                        tint = Color.White,
                        modifier = Modifier.size(30.dp)
                    )
                }
                Text(
                    text = formatCount(take.commentsCount),
                    fontSize = 12.sp,
                    fontWeight = FontWeight.SemiBold,
                    color = Color.White
                )
            }
            
            // Share
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(4.dp)
            ) {
                IconButton(
                    onClick = { showShare = true },
                    modifier = Modifier.size(40.dp)
                ) {
                    Icon(
                        painter = painterResource(R.drawable.ic_share),
                        contentDescription = "Share",
                        tint = Color.White,
                        modifier = Modifier.size(28.dp)
                    )
                }
                Text(
                    text = formatCount(take.sharesCount),
                    fontSize = 12.sp,
                    fontWeight = FontWeight.SemiBold,
                    color = Color.White
                )
            }
            
            // More
            IconButton(
                onClick = {},
                modifier = Modifier.size(40.dp)
            ) {
                Icon(
                    painter = painterResource(R.drawable.ic_more),
                    contentDescription = "More",
                    tint = Color.White,
                    modifier = Modifier.size(24.dp)
                )
            }
        }
        
        // Bottom Info
        Column(
            modifier = Modifier
                .align(Alignment.BottomStart)
                .fillMaxWidth()
                .padding(16.dp)
                .padding(bottom = 100.dp)
                .padding(end = 80.dp),
            verticalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            // Username and Follow
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                Text(
                    text = "@${take.username}",
                    fontSize = 16.sp,
                    fontWeight = FontWeight.SemiBold,
                    color = Color.White
                )
                
                if (!isFollowing) {
                    Button(
                        onClick = { isFollowing = true },
                        colors = ButtonDefaults.buttonColors(
                            containerColor = entativa_primary_blue
                        ),
                        contentPadding = PaddingValues(horizontal = 16.dp, vertical = 6.dp),
                        shape = RoundedCornerShape(4.dp)
                    ) {
                        Text(
                            "Follow",
                            fontSize = 14.sp,
                            fontWeight = FontWeight.SemiBold
                        )
                    }
                }
            }
            
            // Caption
            Text(
                text = take.caption,
                fontSize = 14.sp,
                color = Color.White,
                maxLines = 2
            )
            
            // Audio
            Row(
                horizontalArrangement = Arrangement.spacedBy(8.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Icon(
                    painter = painterResource(R.drawable.ic_music),
                    contentDescription = "Audio",
                    tint = Color.White.copy(alpha = 0.9f),
                    modifier = Modifier.size(12.dp)
                )
                Text(
                    text = take.audioName,
                    fontSize = 13.sp,
                    color = Color.White.copy(alpha = 0.9f),
                    maxLines = 1
                )
            }
        }
    }
    
    // Comments Sheet
    if (showComments) {
        TakeCommentsSheet(
            take = take,
            onDismiss = { showComments = false }
        )
    }
    
    // Share Sheet
    if (showShare) {
        TakeShareSheet(
            take = take,
            onDismiss = { showShare = false }
        )
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TakeCommentsSheet(
    take: Take,
    onDismiss: () -> Unit
) {
    var commentText by remember { mutableStateOf("") }
    
    ModalBottomSheet(
        onDismissRequest = onDismiss,
        containerColor = Color.White
    ) {
        Column(modifier = Modifier.fillMaxWidth()) {
            // Header
            Text(
                text = "${take.commentsCount} comments",
                fontSize = 16.sp,
                fontWeight = FontWeight.SemiBold,
                modifier = Modifier.padding(16.dp)
            )
            
            Divider()
            
            // Comments List
            LazyColumn(
                modifier = Modifier
                    .weight(1f)
                    .fillMaxWidth(),
                contentPadding = PaddingValues(16.dp)
            ) {
                items(15) { index ->
                    TakeCommentRow(
                        username = "user$index",
                        comment = "This is amazing! 🔥",
                        timestamp = "${index + 1}h ago",
                        likesCount = (10..500).random()
                    )
                    Spacer(modifier = Modifier.height(12.dp))
                }
            }
            
            Divider()
            
            // Comment Input
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(12.dp)
            ) {
                Box(
                    modifier = Modifier
                        .size(32.dp)
                        .clip(CircleShape)
                        .background(Color.Gray.copy(alpha = 0.3f))
                )
                
                OutlinedTextField(
                    value = commentText,
                    onValueChange = { commentText = it },
                    placeholder = { Text("Add comment...") },
                    modifier = Modifier.weight(1f),
                    shape = RoundedCornerShape(20.dp)
                )
                
                if (commentText.isNotEmpty()) {
                    TextButton(onClick = { commentText = "" }) {
                        Text("Post", color = entativa_primary_blue)
                    }
                }
            }
        }
    }
}

@Composable
fun TakeCommentRow(
    username: String,
    comment: String,
    timestamp: String,
    likesCount: Int
) {
    var isLiked by remember { mutableStateOf(false) }
    
    Row(
        modifier = Modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        Box(
            modifier = Modifier
                .size(32.dp)
                .clip(CircleShape)
                .background(Color.Gray.copy(alpha = 0.3f))
        )
        
        Column(modifier = Modifier.weight(1f)) {
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                Text(
                    username,
                    fontSize = 13.sp,
                    fontWeight = FontWeight.SemiBold
                )
                Text(
                    timestamp,
                    fontSize = 12.sp,
                    color = entativa_text_secondary
                )
            }
            
            Text(
                comment,
                fontSize = 14.sp,
                modifier = Modifier.padding(top = 4.dp)
            )
            
            Row(
                modifier = Modifier.padding(top = 8.dp),
                horizontalArrangement = Arrangement.spacedBy(16.dp)
            ) {
                Text(
                    "$likesCount likes",
                    fontSize = 12.sp,
                    fontWeight = FontWeight.Medium,
                    color = entativa_text_secondary
                )
                Text(
                    "Reply",
                    fontSize = 12.sp,
                    fontWeight = FontWeight.Medium,
                    color = entativa_text_secondary
                )
            }
        }
        
        IconButton(
            onClick = { isLiked = !isLiked },
            modifier = Modifier.size(24.dp)
        ) {
            Icon(
                painter = painterResource(if (isLiked) R.drawable.ic_heart_filled else R.drawable.ic_heart),
                contentDescription = "Like",
                tint = if (isLiked) Color.Red else entativa_text_secondary,
                modifier = Modifier.size(14.dp)
            )
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TakeShareSheet(
    take: Take,
    onDismiss: () -> Unit
) {
    ModalBottomSheet(
        onDismissRequest = onDismiss,
        containerColor = Color.White
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp)
        ) {
            Text(
                "Share",
                fontSize = 16.sp,
                fontWeight = FontWeight.SemiBold,
                modifier = Modifier.padding(bottom = 20.dp)
            )
            
            // Share Options Grid
            Column(verticalArrangement = Arrangement.spacedBy(16.dp)) {
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceEvenly
                ) {
                    ShareOptionItem("Friends", R.drawable.ic_people)
                    ShareOptionItem("Copy Link", R.drawable.ic_link)
                    ShareOptionItem("Share to...", R.drawable.ic_share)
                }
                
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceEvenly
                ) {
                    ShareOptionItem("Save", R.drawable.ic_bookmark)
                    ShareOptionItem("Report", R.drawable.ic_flag)
                    ShareOptionItem("Not Interested", R.drawable.ic_eye_slash)
                }
            }
            
            Spacer(modifier = Modifier.height(32.dp))
        }
    }
}

@Composable
fun ShareOptionItem(title: String, icon: Int) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(8.dp),
        modifier = Modifier.width(80.dp)
    ) {
        Box(
            modifier = Modifier
                .size(60.dp)
                .clip(CircleShape)
                .background(Color.Gray.copy(alpha = 0.1f)),
            contentAlignment = Alignment.Center
        ) {
            Icon(
                painter = painterResource(icon),
                contentDescription = title,
                tint = entativa_text_primary,
                modifier = Modifier.size(24.dp)
            )
        }
        Text(
            title,
            fontSize = 12.sp,
            color = entativa_text_primary,
            maxLines = 2
        )
    }
}

fun formatCount(count: Int): String {
    return when {
        count >= 1_000_000 -> String.format("%.1fM", count / 1_000_000.0)
        count >= 1_000 -> String.format("%.1fK", count / 1_000.0)
        else -> count.toString()
    }
}

// MARK: - Models
data class Take(
    val id: String,
    val username: String,
    val caption: String,
    val audioName: String,
    val likesCount: Int,
    val commentsCount: Int,
    val sharesCount: Int,
    val viewsCount: Int
)
