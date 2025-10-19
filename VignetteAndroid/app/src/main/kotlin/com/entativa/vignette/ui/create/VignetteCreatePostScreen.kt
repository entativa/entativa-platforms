package com.entativa.vignette.ui.create

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
import com.entativa.vignette.R
import com.google.accompanist.pager.ExperimentalPagerApi
import com.google.accompanist.pager.HorizontalPager
import com.google.accompanist.pager.HorizontalPagerIndicator
import com.google.accompanist.pager.rememberPagerState

// MARK: - Vignette Create Post Screen (Instagram-Style)
@OptIn(ExperimentalMaterial3Api::class, ExperimentalPagerApi::class)
@Composable
fun VignetteCreatePostScreen(
    onDismiss: () -> Unit,
    viewModel: VignetteCreatePostViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
) {
    val selectedImages by viewModel.selectedImages.collectAsState()
    val caption by viewModel.caption.collectAsState()
    
    val imagePickerLauncher = rememberLauncherForActivityResult(
        contract = ActivityResultContracts.GetMultipleContents()
    ) { uris ->
        viewModel.setSelectedImages(uris)
    }
    
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("New Post", fontWeight = FontWeight.Bold) },
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
                        enabled = selectedImages.isNotEmpty()
                    ) {
                        Text(
                            "Share",
                            fontWeight = FontWeight.SemiBold,
                            color = if (selectedImages.isNotEmpty()) Color(0xFF007CFC) else Color.Gray
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
            // Media preview or selector
            if (selectedImages.isEmpty()) {
                // Media selector
                Column(
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(400.dp)
                        .clickable { imagePickerLauncher.launch("image/*") },
                    horizontalAlignment = Alignment.CenterHorizontally,
                    verticalArrangement = Arrangement.Center
                ) {
                    Icon(
                        painter = painterResource(R.drawable.ic_photo),
                        contentDescription = "Select photos",
                        tint = Color(0xFF007CFC),
                        modifier = Modifier.size(64.dp)
                    )
                    
                    Spacer(modifier = Modifier.height(16.dp))
                    
                    Text(
                        "Select photos",
                        fontSize = 18.sp,
                        fontWeight = FontWeight.SemiBold
                    )
                    
                    Spacer(modifier = Modifier.height(8.dp))
                    
                    Text(
                        "You can select up to 10 photos",
                        fontSize = 14.sp,
                        color = Color.Gray
                    )
                }
            } else {
                // Media preview with pager
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(400.dp)
                ) {
                    val pagerState = rememberPagerState()
                    
                    HorizontalPager(
                        count = selectedImages.size,
                        state = pagerState,
                        modifier = Modifier.fillMaxSize()
                    ) { page ->
                        AsyncImage(
                            model = selectedImages[page],
                            contentDescription = null,
                            modifier = Modifier.fillMaxSize(),
                            contentScale = ContentScale.Crop
                        )
                    }
                    
                    // Page indicator
                    if (selectedImages.size > 1) {
                        HorizontalPagerIndicator(
                            pagerState = pagerState,
                            modifier = Modifier
                                .align(Alignment.BottomCenter)
                                .padding(16.dp),
                            activeColor = Color.White,
                            inactiveColor = Color.White.copy(alpha = 0.5f)
                        )
                    }
                }
                
                // Edit tools
                LazyRow(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(vertical = 12.dp),
                    horizontalArrangement = Arrangement.spacedBy(12.dp),
                    contentPadding = PaddingValues(horizontal = 16.dp)
                ) {
                    item { EditToolButton(icon = R.drawable.ic_filter, title = "Filter") }
                    item { EditToolButton(icon = R.drawable.ic_crop, title = "Crop") }
                    item { EditToolButton(icon = R.drawable.ic_adjust, title = "Adjust") }
                    item { EditToolButton(icon = R.drawable.ic_text, title = "Text") }
                    item { EditToolButton(icon = R.drawable.ic_draw, title = "Draw") }
                }
            }
            
            Divider()
            
            // Caption
            OutlinedTextField(
                value = caption,
                onValueChange = { viewModel.setCaption(it) },
                placeholder = { Text("Write a caption...") },
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                minLines = 3,
                maxLines = 10,
                colors = OutlinedTextFieldDefaults.colors(
                    focusedBorderColor = Color.Transparent,
                    unfocusedBorderColor = Color.Transparent
                )
            )
            
            Divider()
            
            // Tag people
            CreatePostOption(
                title = "Tag people",
                onClick = {}
            )
            
            Divider()
            
            // Add location
            CreatePostOption(
                title = "Add location",
                onClick = {}
            )
            
            Divider()
            
            // Cross-post options
            Column(
                modifier = Modifier.padding(16.dp)
            ) {
                Text(
                    "Also post to",
                    fontSize = 13.sp,
                    fontWeight = FontWeight.SemiBold,
                    color = Color.Gray,
                    modifier = Modifier.padding(bottom = 12.dp)
                )
                
                var postToEntativa by remember { mutableStateOf(false) }
                var postToTwitter by remember { mutableStateOf(false) }
                
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text("Entativa", fontSize = 16.sp)
                    Switch(
                        checked = postToEntativa,
                        onCheckedChange = { postToEntativa = it }
                    )
                }
                
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text("Twitter", fontSize = 16.sp)
                    Switch(
                        checked = postToTwitter,
                        onCheckedChange = { postToTwitter = it }
                    )
                }
            }
            
            Divider()
            
            // Advanced settings
            CreatePostOption(
                title = "Advanced settings",
                onClick = {}
            )
        }
    }
}

@Composable
fun EditToolButton(icon: Int, title: String) {
    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        Surface(
            modifier = Modifier.size(50.dp),
            shape = CircleShape,
            color = Color.Gray.copy(alpha = 0.1f)
        ) {
            Box(
                modifier = Modifier.fillMaxSize(),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    painter = painterResource(icon),
                    contentDescription = title,
                    modifier = Modifier.size(24.dp)
                )
            }
        }
        
        Text(
            title,
            fontSize = 12.sp
        )
    }
}

@Composable
fun CreatePostOption(
    title: String,
    onClick: () -> Unit
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(horizontal = 16.dp, vertical = 12.dp),
        horizontalArrangement = Arrangement.SpaceBetween,
        verticalAlignment = Alignment.CenterVertically
    ) {
        Text(
            title,
            fontSize = 16.sp
        )
        
        Icon(
            painter = painterResource(R.drawable.ic_chevron_right),
            contentDescription = null,
            tint = Color.Gray,
            modifier = Modifier.size(14.dp)
        )
    }
}

// ViewModel
class VignetteCreatePostViewModel : androidx.lifecycle.ViewModel() {
    private val _selectedImages = kotlinx.coroutines.flow.MutableStateFlow<List<Uri>>(emptyList())
    val selectedImages: kotlinx.coroutines.flow.StateFlow<List<Uri>> = _selectedImages
    
    private val _caption = kotlinx.coroutines.flow.MutableStateFlow("")
    val caption: kotlinx.coroutines.flow.StateFlow<String> = _caption
    
    fun setSelectedImages(images: List<Uri>) {
        _selectedImages.value = images
    }
    
    fun setCaption(text: String) {
        _caption.value = text
    }
    
    fun createPost() {
        // TODO: Upload images and create post via API
    }
}
