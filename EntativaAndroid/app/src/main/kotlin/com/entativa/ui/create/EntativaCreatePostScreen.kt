package com.entativa.ui.create

import android.net.Uri
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
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

// MARK: - Entativa Create Post Screen (Facebook-Style)
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun EntativaCreatePostScreen(
    onDismiss: () -> Unit,
    viewModel: EntativaCreatePostViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
) {
    val postText by viewModel.postText.collectAsState()
    val selectedImages by viewModel.selectedImages.collectAsState()
    val audience by viewModel.audience.collectAsState()
    
    val imagePickerLauncher = rememberLauncherForActivityResult(
        contract = ActivityResultContracts.GetMultipleContents()
    ) { uris ->
        viewModel.setSelectedImages(uris)
    }
    
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Create Post", fontWeight = FontWeight.Bold) },
                navigationIcon = {
                    IconButton(onClick = onDismiss) {
                        Icon(painter = painterResource(R.drawable.ic_close), contentDescription = "Close")
                    }
                },
                actions = {
                    TextButton(
                        onClick = {
                            viewModel.createPost()
                            onDismiss()
                        },
                        enabled = postText.isNotEmpty() || selectedImages.isNotEmpty()
                    ) {
                        Text(
                            "Post",
                            fontWeight = FontWeight.SemiBold,
                            color = if (postText.isNotEmpty() || selectedImages.isNotEmpty()) 
                                Color(0xFF007CFC) else Color.Gray
                        )
                    }
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .verticalScroll(rememberScrollState())
        ) {
            // User info and audience
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                verticalAlignment = Alignment.Top
            ) {
                // Profile picture
                Surface(
                    modifier = Modifier.size(40.dp),
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
                            tint = Color.Gray
                        )
                    }
                }
                
                Spacer(modifier = Modifier.width(12.dp))
                
                Column {
                    Text(
                        "Your Name",
                        fontSize = 15.sp,
                        fontWeight = FontWeight.SemiBold
                    )
                    
                    Spacer(modifier = Modifier.height(4.dp))
                    
                    // Audience selector
                    Surface(
                        modifier = Modifier.clickable {},
                        shape = RoundedCornerShape(4.dp),
                        color = Color.Gray.copy(alpha = 0.1f)
                    ) {
                        Row(
                            modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp),
                            verticalAlignment = Alignment.CenterVertically,
                            horizontalArrangement = Arrangement.spacedBy(4.dp)
                        ) {
                            Icon(
                                painter = painterResource(
                                    when (audience) {
                                        PostAudience.PUBLIC -> R.drawable.ic_globe
                                        PostAudience.FRIENDS -> R.drawable.ic_people
                                        PostAudience.ONLY_ME -> R.drawable.ic_lock
                                    }
                                ),
                                contentDescription = null,
                                tint = Color.Gray,
                                modifier = Modifier.size(12.dp)
                            )
                            
                            Text(
                                audience.displayName,
                                fontSize = 13.sp,
                                color = Color.Gray
                            )
                            
                            Icon(
                                painter = painterResource(R.drawable.ic_chevron_down),
                                contentDescription = null,
                                tint = Color.Gray,
                                modifier = Modifier.size(10.dp)
                            )
                        }
                    }
                }
            }
            
            // Text input
            OutlinedTextField(
                value = postText,
                onValueChange = { viewModel.setPostText(it) },
                placeholder = { Text("What's on your mind?", fontSize = 18.sp) },
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = 16.dp),
                minLines = 5,
                colors = OutlinedTextFieldDefaults.colors(
                    focusedBorderColor = Color.Transparent,
                    unfocusedBorderColor = Color.Transparent
                )
            )
            
            // Selected media preview
            if (selectedImages.isNotEmpty()) {
                LazyRow(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 16.dp, vertical = 12.dp),
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    items(selectedImages) { uri ->
                        Box {
                            AsyncImage(
                                model = uri,
                                contentDescription = null,
                                modifier = Modifier
                                    .size(120.dp)
                                    .clip(RoundedCornerShape(8.dp)),
                                contentScale = ContentScale.Crop
                            )
                            
                            // Remove button
                            IconButton(
                                onClick = { viewModel.removeImage(uri) },
                                modifier = Modifier
                                    .align(Alignment.TopEnd)
                                    .padding(8.dp)
                            ) {
                                Surface(
                                    shape = CircleShape,
                                    color = Color.Black.copy(alpha = 0.5f)
                                ) {
                                    Icon(
                                        painter = painterResource(R.drawable.ic_close),
                                        contentDescription = "Remove",
                                        tint = Color.White,
                                        modifier = Modifier
                                            .size(24.dp)
                                            .padding(4.dp)
                                    )
                                }
                            }
                        }
                    }
                }
            }
            
            Divider()
            
            // Action buttons
            Column(
                modifier = Modifier.padding(vertical = 8.dp)
            ) {
                PostActionButton(
                    icon = R.drawable.ic_photo,
                    iconColor = Color.Green,
                    title = "Photo/Video",
                    onClick = { imagePickerLauncher.launch("image/*") }
                )
                
                PostActionButton(
                    icon = R.drawable.ic_people,
                    iconColor = Color(0xFF007CFC),
                    title = "Tag people",
                    onClick = {}
                )
                
                PostActionButton(
                    icon = R.drawable.ic_emoji,
                    iconColor = Color.Yellow,
                    title = "Feeling/Activity",
                    onClick = {}
                )
                
                PostActionButton(
                    icon = R.drawable.ic_location,
                    iconColor = Color.Red,
                    title = "Check in",
                    onClick = {}
                )
                
                PostActionButton(
                    icon = R.drawable.ic_camera,
                    iconColor = Color(0xFF6F3EFB),
                    title = "Live video",
                    onClick = {}
                )
                
                PostActionButton(
                    icon = R.drawable.ic_background,
                    iconColor = Color(0xFFFF9800),
                    title = "Background",
                    onClick = {}
                )
                
                PostActionButton(
                    icon = R.drawable.ic_gift,
                    iconColor = Color(0xFFE91E63),
                    title = "Celebration",
                    onClick = {}
                )
            }
        }
    }
}

@Composable
fun PostActionButton(
    icon: Int,
    iconColor: Color,
    title: String,
    onClick: () -> Unit
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(horizontal = 16.dp, vertical = 12.dp),
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        Icon(
            painter = painterResource(icon),
            contentDescription = null,
            tint = iconColor,
            modifier = Modifier.size(24.dp)
        )
        
        Text(
            title,
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

// Models
enum class PostAudience(val displayName: String) {
    PUBLIC("Public"),
    FRIENDS("Friends"),
    ONLY_ME("Only Me")
}

// ViewModel
class EntativaCreatePostViewModel : androidx.lifecycle.ViewModel() {
    private val _postText = kotlinx.coroutines.flow.MutableStateFlow("")
    val postText: kotlinx.coroutines.flow.StateFlow<String> = _postText
    
    private val _selectedImages = kotlinx.coroutines.flow.MutableStateFlow<List<Uri>>(emptyList())
    val selectedImages: kotlinx.coroutines.flow.StateFlow<List<Uri>> = _selectedImages
    
    private val _audience = kotlinx.coroutines.flow.MutableStateFlow(PostAudience.PUBLIC)
    val audience: kotlinx.coroutines.flow.StateFlow<PostAudience> = _audience
    
    fun setPostText(text: String) {
        _postText.value = text
    }
    
    fun setSelectedImages(images: List<Uri>) {
        _selectedImages.value = images
    }
    
    fun removeImage(uri: Uri) {
        _selectedImages.value = _selectedImages.value.filter { it != uri }
    }
    
    fun setAudience(audience: PostAudience) {
        _audience.value = audience
    }
    
    fun createPost() {
        // TODO: Upload media and create post via API
    }
}
