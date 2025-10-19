package com.entativa.vignette.ui.takes

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.rotate
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.entativa.vignette.R
import com.entativa.vignette.ui.theme.*
import com.google.accompanist.pager.ExperimentalPagerApi
import com.google.accompanist.pager.VerticalPager
import com.google.accompanist.pager.rememberPagerState

@OptIn(ExperimentalPagerApi::class)
@Composable
fun VignetteTakesScreen(
    viewModel: VignetteTakesViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
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
            VignetteTakeCard(
                take = takes[page],
                isCurrentlyPlaying = page == pagerState.currentPage
            )
        }
        
        // Top Bar (Instagram style)
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
                    fontSize = 24.sp,
                    fontWeight = FontWeight.Bold,
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
fun VignetteTakeCard(
    take: VignetteTake,
    isCurrentlyPlaying: Boolean
) {
    var isLiked by remember { mutableStateOf(false) }
    var showComments by remember { mutableStateOf(false) }
    var showShare by remember { mutableStateOf(false) }
    var isFollowing by remember { mutableStateOf(false) }
    
    Box(modifier = Modifier.fillMaxSize()) {
        // Video Player (Placeholder)
        Box(
            modifier = Modifier
                .fillMaxSize()
                .background(
                    Brush.linearGradient(
                        colors = listOf(
                            vignette_accent_moonstone.copy(alpha = 0.5f),
                            vignette_gunmetal.copy(alpha = 0.7f),
                            vignette_accent_saffron.copy(alpha = 0.3f)
                        )
                    )
                ),
            contentAlignment = Alignment.Center
        ) {
            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                Icon(
                    painter = painterResource(R.drawable.ic_play),
                    contentDescription = "Play",
                    tint = Color.White.copy(alpha = 0.9f),
                    modifier = Modifier.size(80.dp)
                )
                Spacer(modifier = Modifier.height(8.dp))
                Text(
                    "Video Player",
                    color = Color.White.copy(alpha = 0.7f),
                    fontSize = 16.sp,
                    fontWeight = FontWeight.Light
                )
            }
        }
        
        // Right Side Actions (Instagram Reels style)
        Column(
            modifier = Modifier
                .align(Alignment.CenterEnd)
                .padding(end = 12.dp)
                .padding(bottom = 100.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(20.dp)
        ) {
            // Profile Avatar
            Box(contentAlignment = Alignment.BottomCenter) {
                Box(
                    modifier = Modifier
                        .size(44.dp)
                        .clip(CircleShape)
                        .background(Color.Gray.copy(alpha = 0.5f))
                        .padding(2.dp)
                        .background(Color.White, CircleShape)
                )
                
                if (!isFollowing) {
                    Box(
                        modifier = Modifier
                            .offset(y = 10.dp)
                            .size(22.dp)
                            .clip(CircleShape)
                            .background(entativa_blue),
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
                verticalArrangement = Arrangement.spacedBy(6.dp)
            ) {
                IconButton(
                    onClick = { isLiked = !isLiked },
                    modifier = Modifier.size(36.dp)
                ) {
                    Icon(
                        painter = painterResource(if (isLiked) R.drawable.ic_heart_filled else R.drawable.ic_heart),
                        contentDescription = "Like",
                        tint = if (isLiked) Color.Red else Color.White,
                        modifier = Modifier.size(28.dp)
                    )
                }
                Text(
                    text = formatCount(take.likesCount),
                    fontSize = 12.sp,
                    fontWeight = FontWeight.Medium,
                    color = Color.White
                )
            }
            
            // Comments
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(6.dp)
            ) {
                IconButton(
                    onClick = { showComments = true },
                    modifier = Modifier.size(36.dp)
                ) {
                    Icon(
                        painter = painterResource(R.drawable.ic_comment),
                        contentDescription = "Comments",
                        tint = Color.White,
                        modifier = Modifier.size(28.dp)
                    )
                }
                Text(
                    text = formatCount(take.commentsCount),
                    fontSize = 12.sp,
                    fontWeight = FontWeight.Medium,
                    color = Color.White
                )
            }
            
            // Share
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(6.dp)
            ) {
                IconButton(
                    onClick = { showShare = true },
                    modifier = Modifier.size(36.dp)
                ) {
                    Icon(
                        painter = painterResource(R.drawable.ic_share),
                        contentDescription = "Share",
                        tint = Color.White,
                        modifier = Modifier.size(26.dp)
                    )
                }
                Text(
                    text = formatCount(take.sharesCount),
                    fontSize = 12.sp,
                    fontWeight = FontWeight.Medium,
                    color = Color.White
                )
            }
            
            // Save
            IconButton(
                onClick = {},
                modifier = Modifier.size(36.dp)
            ) {
                Icon(
                    painter = painterResource(R.drawable.ic_bookmark),
                    contentDescription = "Save",
                    tint = Color.White,
                    modifier = Modifier.size(26.dp)
                )
            }
            
            // More
            IconButton(
                onClick = {},
                modifier = Modifier.size(36.dp)
            ) {
                Icon(
                    painter = painterResource(R.drawable.ic_more),
                    contentDescription = "More",
                    tint = Color.White,
                    modifier = Modifier
                        .size(22.dp)
                        .rotate(90f)
                )
            }
            
            // Audio (spinning record)
            Box(
                modifier = Modifier
                    .size(36.dp)
                    .clip(CircleShape)
                    .background(
                        Brush.linearGradient(
                            colors = listOf(
                                vignette_accent_moonstone,
                                vignette_accent_saffron
                            )
                        )
                    ),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    painter = painterResource(R.drawable.ic_music),
                    contentDescription = "Audio",
                    tint = Color.White,
                    modifier = Modifier.size(14.dp)
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
                .padding(end = 70.dp),
            verticalArrangement = Arrangement.spacedBy(10.dp)
        ) {
            // Username and Follow
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                Text(
                    text = "@${take.username}",
                    fontSize = 15.sp,
                    fontWeight = FontWeight.SemiBold,
                    color = Color.White
                )
                
                if (!isFollowing) {
                    OutlinedButton(
                        onClick = { isFollowing = true },
                        colors = ButtonDefaults.outlinedButtonColors(
                            contentColor = Color.White
                        ),
                        border = ButtonDefaults.outlinedButtonBorder.copy(
                            brush = Brush.linearGradient(listOf(Color.White, Color.White))
                        ),
                        contentPadding = PaddingValues(horizontal = 12.dp, vertical = 4.dp),
                        shape = RoundedCornerShape(4.dp),
                        modifier = Modifier.height(28.dp)
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
                fontSize = 13.sp,
                color = Color.White,
                maxLines = 2
            )
            
            // Audio
            Row(
                horizontalArrangement = Arrangement.spacedBy(6.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Icon(
                    painter = painterResource(R.drawable.ic_music),
                    contentDescription = "Audio",
                    tint = Color.White,
                    modifier = Modifier.size(11.dp)
                )
                Text(
                    text = take.audioName,
                    fontSize = 12.sp,
                    color = Color.White,
                    maxLines = 1
                )
            }
        }
    }
    
    // Comments Sheet
    if (showComments) {
        VignetteTakeCommentsSheet(
            take = take,
            onDismiss = { showComments = false }
        )
    }
    
    // Share Sheet
    if (showShare) {
        VignetteTakeShareSheet(
            take = take,
            onDismiss = { showShare = false }
        )
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VignetteTakeCommentsSheet(
    take: VignetteTake,
    onDismiss: () -> Unit
) {
    var commentText by remember { mutableStateOf("") }
    
    ModalBottomSheet(
        onDismissRequest = onDismiss,
        containerColor = Color.White
    ) {
        Column(modifier = Modifier.fillMaxWidth()) {
            // Header
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(
                    text = "Comments",
                    fontSize = 16.sp,
                    fontWeight = FontWeight.SemiBold
                )
                
                IconButton(onClick = onDismiss) {
                    Icon(
                        painter = painterResource(R.drawable.ic_close),
                        contentDescription = "Close",
                        tint = vignette_text_primary
                    )
                }
            }
            
            Divider()
            
            // Comments List
            LazyColumn(
                modifier = Modifier
                    .weight(1f)
                    .fillMaxWidth(),
                contentPadding = PaddingValues(16.dp)
            ) {
                items(15) { index ->
                    VignetteTakeCommentRow(
                        username = "user$index",
                        comment = "This is awesome! 🔥",
                        timestamp = "${index + 1}h",
                        likesCount = (5..500).random()
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
                        .size(28.dp)
                        .clip(CircleShape)
                        .background(Color.Gray.copy(alpha = 0.2f))
                )
                
                TextField(
                    value = commentText,
                    onValueChange = { commentText = it },
                    placeholder = { Text("Add a comment...", fontSize = 14.sp) },
                    modifier = Modifier.weight(1f),
                    colors = TextFieldDefaults.colors(
                        focusedContainerColor = Color.Transparent,
                        unfocusedContainerColor = Color.Transparent,
                        focusedIndicatorColor = Color.Transparent,
                        unfocusedIndicatorColor = Color.Transparent
                    )
                )
                
                if (commentText.isNotEmpty()) {
                    TextButton(onClick = { commentText = "" }) {
                        Text("Post", color = entativa_blue, fontSize = 14.sp, fontWeight = FontWeight.SemiBold)
                    }
                }
            }
        }
    }
}

@Composable
fun VignetteTakeCommentRow(
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
                .background(Color.Gray.copy(alpha = 0.2f))
        )
        
        Column(modifier = Modifier.weight(1f)) {
            Text(
                username,
                fontSize = 13.sp,
                fontWeight = FontWeight.SemiBold
            )
            
            Text(
                comment,
                fontSize = 13.sp,
                modifier = Modifier.padding(top = 4.dp)
            )
            
            Row(
                modifier = Modifier.padding(top = 8.dp),
                horizontalArrangement = Arrangement.spacedBy(12.dp)
            ) {
                Text(
                    timestamp,
                    fontSize = 12.sp,
                    color = vignette_text_secondary
                )
                
                if (likesCount > 0) {
                    Text(
                        "$likesCount likes",
                        fontSize = 12.sp,
                        color = vignette_text_secondary
                    )
                }
                
                Text(
                    "Reply",
                    fontSize = 12.sp,
                    color = vignette_text_secondary
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
                tint = if (isLiked) Color.Red else vignette_text_secondary,
                modifier = Modifier.size(12.dp)
            )
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VignetteTakeShareSheet(
    take: VignetteTake,
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
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(
                    "Share",
                    fontSize = 16.sp,
                    fontWeight = FontWeight.SemiBold
                )
                
                IconButton(onClick = onDismiss) {
                    Icon(
                        painter = painterResource(R.drawable.ic_close),
                        contentDescription = "Close"
                    )
                }
            }
            
            Spacer(modifier = Modifier.height(24.dp))
            
            // Share Options
            Column(verticalArrangement = Arrangement.spacedBy(20.dp)) {
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceEvenly
                ) {
                    VignetteShareOptionItem("Copy Link", R.drawable.ic_link)
                    VignetteShareOptionItem("Share", R.drawable.ic_share)
                    VignetteShareOptionItem("Save", R.drawable.ic_bookmark)
                }
                
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceEvenly
                ) {
                    VignetteShareOptionItem("Not Interested", R.drawable.ic_eye_slash)
                    VignetteShareOptionItem("Report", R.drawable.ic_flag)
                    VignetteShareOptionItem("Hide", R.drawable.ic_eye_slash)
                }
            }
            
            Spacer(modifier = Modifier.height(32.dp))
        }
    }
}

@Composable
fun VignetteShareOptionItem(title: String, icon: Int) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(8.dp),
        modifier = Modifier.width(80.dp)
    ) {
        Box(
            modifier = Modifier
                .size(56.dp)
                .clip(CircleShape)
                .background(Color.Transparent)
                .padding(1.dp)
                .background(Color.Transparent),
            contentAlignment = Alignment.Center
        ) {
            Icon(
                painter = painterResource(icon),
                contentDescription = title,
                tint = vignette_text_primary,
                modifier = Modifier.size(22.dp)
            )
        }
        Text(
            title,
            fontSize = 12.sp,
            color = vignette_text_primary,
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
data class VignetteTake(
    val id: String,
    val username: String,
    val caption: String,
    val audioName: String,
    val likesCount: Int,
    val commentsCount: Int,
    val sharesCount: Int,
    val viewsCount: Int
)
