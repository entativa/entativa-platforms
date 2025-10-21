package com.entativa.vignette.ui.messages

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
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.entativa.vignette.R

// MARK: - Vignette Messages Screen (Instagram Direct-Style)
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VignetteMessagesScreen(
    viewModel: VignetteMessagesViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
) {
    val conversations by viewModel.conversations.collectAsState()
    var searchText by remember { mutableStateOf("") }
    var showNewMessage by remember { mutableStateOf(false) }
    
    Scaffold(
        topBar = {
            Column(modifier = Modifier.background(Color.White)) {
                // Header
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(16.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Text(
                            "yourusername",
                            fontSize = 20.sp,
                            fontWeight = FontWeight.Bold
                        )
                        Spacer(modifier = Modifier.width(8.dp))
                        Icon(
                            painter = painterResource(R.drawable.ic_chevron_down),
                            contentDescription = null,
                            modifier = Modifier.size(14.dp)
                        )
                    }
                    
                    IconButton(onClick = { showNewMessage = true }) {
                        Icon(
                            painter = painterResource(R.drawable.ic_edit),
                            contentDescription = "New message",
                            modifier = Modifier.size(22.dp)
                        )
                    }
                }
                
                // Search bar
                OutlinedTextField(
                    value = searchText,
                    onValueChange = { searchText = it },
                    placeholder = { Text("Search") },
                    leadingIcon = {
                        Icon(
                            painter = painterResource(R.drawable.ic_search),
                            contentDescription = null,
                            tint = Color.Gray
                        )
                    },
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 16.dp)
                        .padding(bottom = 8.dp),
                    shape = RoundedCornerShape(10.dp),
                    colors = OutlinedTextFieldDefaults.colors(
                        unfocusedBorderColor = Color.Gray.copy(alpha = 0.2f),
                        focusedBorderColor = Color.Gray.copy(alpha = 0.3f)
                    )
                )
                
                Divider()
            }
        }
    ) { paddingValues ->
        LazyColumn(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            items(conversations) { conversation ->
                VignetteConversationRow(
                    conversation = conversation,
                    onClick = { /* Navigate to chat */ }
                )
            }
        }
    }
}

@Composable
fun VignetteConversationRow(
    conversation: VignetteConversation,
    onClick: () -> Unit
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .background(if (conversation.hasUnread) Color(0xFFE3F2FD) else Color.White)
            .padding(horizontal = 16.dp, vertical = 8.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        // Avatar with online indicator
        Box {
            Surface(
                modifier = Modifier.size(56.dp),
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
            
            if (conversation.isOnline) {
                Surface(
                    modifier = Modifier
                        .size(14.dp)
                        .align(Alignment.BottomEnd),
                    shape = CircleShape,
                    color = Color.Green,
                    border = androidx.compose.foundation.BorderStroke(2.dp, Color.White)
                ) {}
            }
        }
        
        Spacer(modifier = Modifier.width(12.dp))
        
        // Content
        Column(modifier = Modifier.weight(1f)) {
            Row {
                Text(
                    conversation.name,
                    fontSize = 15.sp,
                    fontWeight = if (conversation.hasUnread) FontWeight.SemiBold else FontWeight.Normal
                )
                Spacer()
                Text(
                    conversation.timeAgo,
                    fontSize = 13.sp,
                    color = if (conversation.hasUnread) Color.Black else Color.Gray
                )
            }
            
            Spacer(modifier = Modifier.height(4.dp))
            
            Row {
                Text(
                    conversation.lastMessage,
                    fontSize = 14.sp,
                    color = if (conversation.hasUnread) Color.Black else Color.Gray,
                    maxLines = 2
                )
                Spacer()
                if (conversation.unreadCount > 0) {
                    Box(
                        modifier = Modifier.size(8.dp),
                        contentAlignment = Alignment.Center
                    ) {
                        Surface(
                            shape = CircleShape,
                            color = Color(0xFF007CFC)
                        ) {}
                    }
                }
            }
        }
    }
}

// MARK: - Chat Screen
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VignetteChatScreen(
    conversation: VignetteConversation,
    viewModel: VignetteChatViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
) {
    val messages by viewModel.messages.collectAsState()
    var messageText by remember { mutableStateOf("") }
    
    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Row(verticalAlignment = Alignment.CenterVertically) {
                        Surface(
                            modifier = Modifier.size(32.dp),
                            shape = CircleShape,
                            color = Color.Gray.copy(alpha = 0.2f)
                        ) {}
                        Spacer(modifier = Modifier.width(8.dp))
                        Column {
                            Text(conversation.name, fontSize = 16.sp, fontWeight = FontWeight.SemiBold)
                            if (conversation.isOnline) {
                                Text("Active now", fontSize = 12.sp, color = Color.Gray)
                            }
                        }
                    }
                },
                actions = {
                    IconButton(onClick = {}) {
                        Icon(painter = painterResource(R.drawable.ic_phone), contentDescription = "Call")
                    }
                    IconButton(onClick = {}) {
                        Icon(painter = painterResource(R.drawable.ic_video), contentDescription = "Video")
                    }
                    IconButton(onClick = {}) {
                        Icon(painter = painterResource(R.drawable.ic_info), contentDescription = "Info")
                    }
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            // Messages
            LazyColumn(
                modifier = Modifier.weight(1f),
                reverseLayout = true
            ) {
                items(messages) { message ->
                    VignetteMessageBubble(message = message)
                }
            }
            
            Divider()
            
            // Input bar
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(8.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                IconButton(onClick = {}) {
                    Icon(
                        painter = painterResource(R.drawable.ic_camera),
                        contentDescription = "Camera",
                        tint = Color(0xFF007CFC)
                    )
                }
                
                IconButton(onClick = {}) {
                    Icon(
                        painter = painterResource(R.drawable.ic_photo),
                        contentDescription = "Photo",
                        tint = Color(0xFF007CFC)
                    )
                }
                
                // Text field
                Row(
                    modifier = Modifier
                        .weight(1f)
                        .clip(RoundedCornerShape(20.dp))
                        .background(Color.Gray.copy(alpha = 0.1f))
                        .padding(horizontal = 12.dp, vertical = 8.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    TextField(
                        value = messageText,
                        onValueChange = { messageText = it },
                        placeholder = { Text("Message...") },
                        modifier = Modifier.weight(1f),
                        colors = TextFieldDefaults.colors(
                            unfocusedContainerColor = Color.Transparent,
                            focusedContainerColor = Color.Transparent,
                            unfocusedIndicatorColor = Color.Transparent,
                            focusedIndicatorColor = Color.Transparent
                        )
                    )
                    
                    if (messageText.isNotEmpty()) {
                        IconButton(onClick = {
                            viewModel.sendMessage(messageText)
                            messageText = ""
                        }) {
                            Icon(
                                painter = painterResource(R.drawable.ic_send),
                                contentDescription = "Send",
                                tint = Color(0xFF007CFC)
                            )
                        }
                    } else {
                        Row {
                            IconButton(onClick = {}) {
                                Icon(
                                    painter = painterResource(R.drawable.ic_mic),
                                    contentDescription = "Voice",
                                    tint = Color(0xFF007CFC),
                                    modifier = Modifier.size(20.dp)
                                )
                            }
                            IconButton(onClick = {}) {
                                Icon(
                                    painter = painterResource(R.drawable.ic_emoji),
                                    contentDescription = "Emoji",
                                    tint = Color(0xFF007CFC),
                                    modifier = Modifier.size(20.dp)
                                )
                            }
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun VignetteMessageBubble(message: VignetteChatMessage) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp, vertical = 4.dp),
        horizontalArrangement = if (message.isSender) Arrangement.End else Arrangement.Start
    ) {
        Column(
            horizontalAlignment = if (message.isSender) Alignment.End else Alignment.Start
        ) {
            // Message bubble
            Text(
                message.content,
                fontSize = 15.sp,
                color = if (message.isSender) Color.White else Color.Black,
                modifier = Modifier
                    .clip(RoundedCornerShape(18.dp))
                    .background(if (message.isSender) Color(0xFF007CFC) else Color.Gray.copy(alpha = 0.15f))
                    .padding(horizontal = 12.dp, vertical = 8.dp)
            )
            
            // Status
            Row(
                horizontalArrangement = Arrangement.spacedBy(4.dp)
            ) {
                Text(
                    message.timeAgo,
                    fontSize = 11.sp,
                    color = Color.Gray
                )
                
                if (message.isSender) {
                    if (message.isRead) {
                        Text("Read", fontSize = 11.sp, color = Color.Gray)
                    } else if (message.isDelivered) {
                        Text("Delivered", fontSize = 11.sp, color = Color.Gray)
                    }
                    
                    Icon(
                        painter = painterResource(R.drawable.ic_lock),
                        contentDescription = "Encrypted",
                        tint = Color.Gray,
                        modifier = Modifier.size(9.dp)
                    )
                }
            }
        }
    }
}

// Models
data class VignetteConversation(
    val id: String,
    val name: String,
    val avatarUrl: String?,
    val lastMessage: String,
    val lastMessageIsYours: Boolean,
    val timeAgo: String,
    val hasUnread: Boolean,
    val unreadCount: Int,
    val isOnline: Boolean,
    val isEncrypted: Boolean
)

data class VignetteChatMessage(
    val id: String,
    val content: String,
    val isSender: Boolean,
    val timeAgo: String,
    val isRead: Boolean,
    val isDelivered: Boolean,
    val timestamp: Long
)

// ViewModels
class VignetteMessagesViewModel : androidx.lifecycle.ViewModel() {
    private val _conversations = kotlinx.coroutines.flow.MutableStateFlow<List<VignetteConversation>>(emptyList())
    val conversations: kotlinx.coroutines.flow.StateFlow<List<VignetteConversation>> = _conversations
    
    init {
        loadMockData()
    }
    
    private fun loadMockData() {
        _conversations.value = listOf(
            VignetteConversation(
                id = "1",
                name = "sarah_jones",
                avatarUrl = null,
                lastMessage = "That sounds great! When are you free?",
                lastMessageIsYours = false,
                timeAgo = "2m",
                hasUnread = true,
                unreadCount = 2,
                isOnline = true,
                isEncrypted = true
            ),
            VignetteConversation(
                id = "2",
                name = "mike_wilson",
                avatarUrl = null,
                lastMessage = "Thanks for sharing! 🔥",
                lastMessageIsYours = true,
                timeAgo = "1h",
                hasUnread = false,
                unreadCount = 0,
                isOnline = false,
                isEncrypted = true
            )
        )
    }
}

class VignetteChatViewModel : androidx.lifecycle.ViewModel() {
    private val _messages = kotlinx.coroutines.flow.MutableStateFlow<List<VignetteChatMessage>>(emptyList())
    val messages: kotlinx.coroutines.flow.StateFlow<List<VignetteChatMessage>> = _messages
    
    init {
        loadMockMessages()
    }
    
    private fun loadMockMessages() {
        val now = System.currentTimeMillis()
        _messages.value = listOf(
            VignetteChatMessage(
                id = "1",
                content = "Hey! How are you?",
                isSender = false,
                timeAgo = "10:30 AM",
                isRead = true,
                isDelivered = true,
                timestamp = now - 3600000
            ),
            VignetteChatMessage(
                id = "2",
                content = "I'm great! Just finished a new project 🎉",
                isSender = true,
                timeAgo = "10:32 AM",
                isRead = true,
                isDelivered = true,
                timestamp = now - 3480000
            )
        )
    }
    
    fun sendMessage(text: String) {
        // TODO: Send via WebSocket with E2EE
    }
}
