package com.entativa.ui.takes

import android.content.Context
import android.net.Uri
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material3.Icon
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.unit.dp
import androidx.compose.ui.viewinterop.AndroidView
import androidx.media3.common.MediaItem
import androidx.media3.common.Player
import androidx.media3.exoplayer.ExoPlayer
import androidx.media3.ui.PlayerView
import com.entativa.R

@Composable
fun VideoPlayer(
    videoUrl: String,
    isPlaying: Boolean,
    modifier: Modifier = Modifier,
    onPlayerReady: (ExoPlayer) -> Unit = {}
) {
    val context = LocalContext.current
    var isMuted by remember { mutableStateOf(false) }
    var showControls by remember { mutableStateOf(false) }
    
    val exoPlayer = remember {
        ExoPlayer.Builder(context).build().apply {
            val mediaItem = MediaItem.fromUri(Uri.parse(videoUrl))
            setMediaItem(mediaItem)
            prepare()
            repeatMode = Player.REPEAT_MODE_ONE
            volume = if (isMuted) 0f else 1f
            onPlayerReady(this)
        }
    }
    
    DisposableEffect(isPlaying) {
        if (isPlaying) {
            exoPlayer.play()
        } else {
            exoPlayer.pause()
        }
        
        onDispose {
            exoPlayer.pause()
        }
    }
    
    DisposableEffect(Unit) {
        onDispose {
            exoPlayer.release()
        }
    }
    
    Box(modifier = modifier.fillMaxSize()) {
        AndroidView(
            factory = { ctx ->
                PlayerView(ctx).apply {
                    player = exoPlayer
                    useController = false
                    controllerAutoShow = false
                }
            },
            modifier = Modifier
                .fillMaxSize()
                .clickable {
                    // Toggle play/pause
                    if (exoPlayer.isPlaying) {
                        exoPlayer.pause()
                    } else {
                        exoPlayer.play()
                    }
                    showControls = !showControls
                }
        )
        
        // Mute/Unmute Button
        Box(
            modifier = Modifier
                .align(Alignment.BottomEnd)
                .padding(end = 16.dp)
                .padding(bottom = 200.dp)
        ) {
            IconButton(
                onClick = {
                    isMuted = !isMuted
                    exoPlayer.volume = if (isMuted) 0f else 1f
                },
                modifier = Modifier
                    .size(40.dp)
                    .clip(CircleShape)
                    .background(Color.Black.copy(alpha = 0.5f))
            ) {
                Icon(
                    painter = painterResource(if (isMuted) R.drawable.ic_volume_off else R.drawable.ic_volume_up),
                    contentDescription = if (isMuted) "Unmute" else "Mute",
                    tint = Color.White,
                    modifier = Modifier.size(20.dp)
                )
            }
        }
        
        // Play/Pause Indicator (temporary)
        if (showControls) {
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .background(Color.Black.copy(alpha = 0.3f)),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    painter = painterResource(if (exoPlayer.isPlaying) R.drawable.ic_pause else R.drawable.ic_play),
                    contentDescription = null,
                    tint = Color.White.copy(alpha = 0.8f),
                    modifier = Modifier.size(64.dp)
                )
            }
            
            LaunchedEffect(Unit) {
                kotlinx.coroutines.delay(1000)
                showControls = false
            }
        }
    }
}

@Composable
fun IconButton(
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    content: @Composable () -> Unit
) {
    Box(
        modifier = modifier.clickable(onClick = onClick),
        contentAlignment = Alignment.Center
    ) {
        content()
    }
}

// MARK: - Video Cache Manager
object VideoCache {
    private val cache = mutableMapOf<String, ExoPlayer>()
    
    fun preload(context: Context, videoUrl: String) {
        if (!cache.containsKey(videoUrl)) {
            val player = ExoPlayer.Builder(context).build().apply {
                val mediaItem = MediaItem.fromUri(Uri.parse(videoUrl))
                setMediaItem(mediaItem)
                prepare()
            }
            cache[videoUrl] = player
        }
    }
    
    fun getPlayer(context: Context, videoUrl: String): ExoPlayer {
        return cache[videoUrl] ?: ExoPlayer.Builder(context).build().apply {
            val mediaItem = MediaItem.fromUri(Uri.parse(videoUrl))
            setMediaItem(mediaItem)
            prepare()
            cache[videoUrl] = this
        }
    }
    
    fun release(videoUrl: String) {
        cache[videoUrl]?.release()
        cache.remove(videoUrl)
    }
    
    fun releaseAll() {
        cache.values.forEach { it.release() }
        cache.clear()
    }
}
