package com.entativa.ui.messages

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
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.entativa.R

// MARK: - Entativa Messages Screen (Facebook Messenger-Style)
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun EntativaMessagesScreen(
    viewModel: EntativaMessagesViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
) {
    val conversations by viewModel.conversations.collectAsState()
    var selectedTab by remember { mutableStateOf(MessengerTab.CHATS) }
    var showNewMessage by remember { mutableStateOf(false) }
    
    Scaffold(
        topBar = {
            Column {
                TopAppBar(
                    title = { Text("Messages", fontWeight = FontWeight.Bold) },
                    navigationIcon = {
                        IconButton(onClick = {}) {
                            Icon(painter = painterResource(R.drawable.ic_settings), contentDescription = "Settings")
                        }
                    },
                    actions = {
                        IconButton(onClick = { showNewMessage = true }) {
                            Icon(painter = painterResource(R.drawable.ic_edit), contentDescription = "New")
                        }
                    }
                )
                
                // Tabs
                Row(modifier = Modifier.fillMaxWidth().height(44.dp)) {
                    MessengerTabButton("Chats", MessengerTab.CHATS, selectedTab) { selectedTab = it }
                    MessengerTabButton("Calls", MessengerTab.CALLS, selectedTab) { selectedTab = it }
                    MessengerTabButton("People", MessengerTab.PEOPLE, selectedTab) { selectedTab = it }
                }
                
                Divider()
            }
        }
    ) { paddingValues ->
        Column(modifier = Modifier.padding(paddingValues)) {
            when (selectedTab) {
                MessengerTab.CHATS -> EntativaChatsListView(conversations)
                MessengerTab.CALLS -> EntativaCallsListView()
                MessengerTab.PEOPLE -> EntativaPeopleListView()
            }
        }
    }
}

@Composable
fun MessengerTabButton(
    title: String,
    tab: MessengerTab,
    selectedTab: MessengerTab,
    onSelect: (MessengerTab) -> Unit
) {
    Column(
        modifier = Modifier
            .weight(1f)
            .fillMaxHeight()
            .clickable { onSelect(tab) },
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Text(
            title,
            fontSize = 16.sp,
            fontWeight = if (selectedTab == tab) FontWeight.SemiBold else FontWeight.Normal,
            color = if (selectedTab == tab) Color.Black else Color.Gray
        )
        Spacer(modifier = Modifier.height(8.dp))
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .height(2.dp)
                .background(if (selectedTab == tab) Color(0xFF007CFC) else Color.Transparent)
        )
    }
}

@Composable
fun EntativaChatsListView(conversations: List<EntativaConv>) {
    var searchText by remember { mutableStateOf("") }
    
    LazyColumn {
        // Search
        item {
            OutlinedTextField(
                value = searchText,
                onValueChange = { searchText = it },
                placeholder = { Text("Search messages") },
                leadingIcon = {
                    Icon(painter = painterResource(R.drawable.ic_search), contentDescription = null)
                },
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp),
                shape = RoundedCornerShape(10.dp)
            )
        }
        
        // Conversations
        items(conversations) { conversation ->
            EntativaConversationRow(conversation)
        }
    }
}

@Composable
fun EntativaConversationRow(conversation: EntativaConv) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable {}
            .background(if (conversation.hasUnread) Color(0xFFE3F2FD) else Color.White)
            .padding(16.dp)
    ) {
        // Avatar
        Box {
            Surface(
                modifier = Modifier.size(56.dp),
                shape = CircleShape,
                color = Color.Gray.copy(alpha = 0.2f)
            ) {}
            
            if (conversation.isOnline) {
                Surface(
                    modifier = Modifier
                        .size(16.dp)
                        .align(Alignment.BottomEnd),
                    shape = CircleShape,
                    color = Color.Green,
                    border = androidx.compose.foundation.BorderStroke(2.dp, Color.White)
                ) {}
            }
        }
        
        Spacer(modifier = Modifier.width(12.dp))
        
        Column(modifier = Modifier.weight(1f)) {
            Row {
                Text(
                    conversation.name,
                    fontSize = 15.sp,
                    fontWeight = if (conversation.hasUnread) FontWeight.Bold else FontWeight.SemiBold
                )
                Spacer(modifier = Modifier.width(4.dp))
                Icon(
                    painter = painterResource(R.drawable.ic_lock),
                    contentDescription = "Encrypted",
                    tint = Color.Gray,
                    modifier = Modifier.size(10.dp)
                )
                Spacer()
                Row(horizontalArrangement = Arrangement.spacedBy(4.dp)) {
                    Text(conversation.timeAgo, fontSize = 13.sp, color = Color.Gray)
                    if (conversation.unreadCount > 0) {
                        Surface(
                            modifier = Modifier.size(8.dp),
                            shape = CircleShape,
                            color = Color(0xFF007CFC)
                        ) {}
                    }
                }
            }
            
            Spacer(modifier = Modifier.height(4.dp))
            
            Text(
                conversation.lastMessage,
                fontSize = 14.sp,
                color = if (conversation.hasUnread) Color.Black else Color.Gray,
                maxLines = 2
            )
        }
    }
}

@Composable
fun EntativaCallsListView() {
    Text("Calls list - Coming soon", modifier = Modifier.padding(16.dp))
}

@Composable
fun EntativaPeopleListView() {
    Text("People list - Coming soon", modifier = Modifier.padding(16.dp))
}

// MARK: - Chat Screen
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun EntativaChatScreen(
    conversation: EntativaConv,
    viewModel: EntativaChatViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
) {
    val messages by viewModel.messages.collectAsState()
    var messageText by remember { mutableStateOf("") }
    
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text(conversation.name, fontWeight = FontWeight.SemiBold) },
                actions = {
                    IconButton(onClick = {}) {
                        Icon(painter = painterResource(R.drawable.ic_phone), contentDescription = "Call", tint = Color(0xFF007CFC))
                    }
                    IconButton(onClick = {}) {
                        Icon(painter = painterResource(R.drawable.ic_video), contentDescription = "Video", tint = Color(0xFF007CFC))
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
                reverseLayout = true,
                contentPadding = PaddingValues(16.dp)
            ) {
                items(messages) { message ->
                    EntativaMessageBubble(message)
                }
            }
            
            // E2EE indicator
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .background(Color.Gray.copy(alpha = 0.05f))
                    .padding(vertical = 6.dp),
                horizontalArrangement = Arrangement.Center
            ) {
                Icon(
                    painter = painterResource(R.drawable.ic_lock),
                    contentDescription = null,
                    tint = Color.Gray,
                    modifier = Modifier.size(11.dp)
                )
                Spacer(modifier = Modifier.width(4.dp))
                Text(
                    "End-to-end encrypted • Messages are secure",
                    fontSize = 11.sp,
                    color = Color.Gray
                )
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
                    Icon(painter = painterResource(R.drawable.ic_plus_circle), contentDescription = "More", tint = Color(0xFF007CFC))
                }
                IconButton(onClick = {}) {
                    Icon(painter = painterResource(R.drawable.ic_camera), contentDescription = "Camera", tint = Color(0xFF007CFC))
                }
                IconButton(onClick = {}) {
                    Icon(painter = painterResource(R.drawable.ic_photo), contentDescription = "Photo", tint = Color(0xFF007CFC))
                }
                IconButton(onClick = {}) {
                    Icon(painter = painterResource(R.drawable.ic_mic), contentDescription = "Voice", tint = Color(0xFF007CFC))
                }
                
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
                        TextButton(onClick = {
                            viewModel.sendMessage(messageText)
                            messageText = ""
                        }) {
                            Text("Send", fontWeight = FontWeight.SemiBold, color = Color(0xFF007CFC))
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun EntativaMessageBubble(message: EntativaChatMessage) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 4.dp),
        horizontalArrangement = if (message.isSender) Arrangement.End else Arrangement.Start
    ) {
        Column(
            horizontalAlignment = if (message.isSender) Alignment.End else Alignment.Start
        ) {
            Surface(
                shape = RoundedCornerShape(18.dp),
                color = if (message.isSender) Color.Transparent else Color.Gray.copy(alpha = 0.15f)
            ) {
                Box(
                    modifier = Modifier
                        .then(
                            if (message.isSender) {
                                Modifier.background(
                                    Brush.linearGradient(
                                        colors = listOf(
                                            Color(0xFF007CFC),
                                            Color(0xFF6F3EFB)
                                        )
                                    )
                                )
                            } else {
                                Modifier
                            }
                        )
                        .padding(horizontal = 14.dp, vertical = 10.dp)
                ) {
                    Text(
                        message.content,
                        fontSize = 15.sp,
                        color = if (message.isSender) Color.White else Color.Black
                    )
                }
            }
            
            Spacer(modifier = Modifier.height(2.dp))
            
            Row(horizontalArrangement = Arrangement.spacedBy(4.dp)) {
                if (message.isSender) {
                    Icon(
                        painter = painterResource(R.drawable.ic_lock),
                        contentDescription = "Encrypted",
                        tint = Color.Gray,
                        modifier = Modifier.size(8.dp)
                    )
                    
                    if (message.isRead) {
                        Icon(
                            painter = painterResource(R.drawable.ic_check_circle_filled),
                            contentDescription = "Read",
                            tint = Color(0xFF007CFC),
                            modifier = Modifier.size(12.dp)
                        )
                    } else if (message.isDelivered) {
                        Icon(
                            painter = painterResource(R.drawable.ic_circle),
                            contentDescription = "Delivered",
                            tint = Color.Gray,
                            modifier = Modifier.size(12.dp)
                        )
                    }
                }
                
                Text(message.timeAgo, fontSize = 11.sp, color = Color.Gray)
            }
        }
    }
}

// Models
enum class MessengerTab { CHATS, CALLS, PEOPLE }

data class EntativaConv(
    val id: String,
    val name: String,
    val avatarUrl: String?,
    val lastMessage: String,
    val timeAgo: String,
    val hasUnread: Boolean,
    val unreadCount: Int,
    val isOnline: Boolean,
    val isEncrypted: Boolean
)

data class EntativaChatMessage(
    val id: String,
    val content: String,
    val isSender: Boolean,
    val timeAgo: String,
    val isRead: Boolean,
    val isDelivered: Boolean,
    val timestamp: Long
)

// ViewModels
class EntativaMessagesViewModel : androidx.lifecycle.ViewModel() {
    private val _conversations = kotlinx.coroutines.flow.MutableStateFlow<List<EntativaConv>>(emptyList())
    val conversations: kotlinx.coroutines.flow.StateFlow<List<EntativaConv>> = _conversations
    
    init {
        loadMockData()
    }
    
    private fun loadMockData() {
        _conversations.value = listOf(
            EntativaConv(
                id = "1",
                name = "Sarah Johnson",
                avatarUrl = null,
                lastMessage = "You: Sounds good! See you then 👍",
                timeAgo = "5m",
                hasUnread = true,
                unreadCount = 3,
                isOnline = true,
                isEncrypted = true
            ),
            EntativaConv(
                id = "2",
                name = "Mike Wilson",
                avatarUrl = null,
                lastMessage = "That's perfect, thanks!",
                timeAgo = "2h",
                hasUnread = false,
                unreadCount = 0,
                isOnline = false,
                isEncrypted = true
            )
        )
    }
}

class EntativaChatViewModel : androidx.lifecycle.ViewModel() {
    private val _messages = kotlinx.coroutines.flow.MutableStateFlow<List<EntativaChatMessage>>(emptyList())
    val messages: kotlinx.coroutines.flow.StateFlow<List<EntativaChatMessage>> = _messages
    
    init {
        loadMockMessages()
    }
    
    private fun loadMockMessages() {
        val now = System.currentTimeMillis()
        _messages.value = listOf(
            EntativaChatMessage(
                id = "1",
                content = "Hey! How's it going?",
                isSender = false,
                timeAgo = "10:30 AM",
                isRead = true,
                isDelivered = true,
                timestamp = now - 3600000
            ),
            EntativaChatMessage(
                id = "2",
                content = "Going great! Just wrapped up that project we discussed",
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
