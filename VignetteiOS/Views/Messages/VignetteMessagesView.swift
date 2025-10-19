import SwiftUI

// MARK: - Vignette Messages View (Instagram Direct-Style)
struct VignetteMessagesView: View {
    @StateObject private var viewModel = MessagesViewModel()
    @State private var showNewMessage = false
    
    var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Header
                HStack {
                    Text("yourusername")
                        .font(.system(size: 20, weight: .bold))
                    
                    Image(systemName: "chevron.down")
                        .font(.system(size: 14, weight: .semibold))
                    
                    Spacer()
                    
                    Button(action: { showNewMessage = true }) {
                        Image(systemName: "square.and.pencil")
                            .font(.system(size: 22))
                            .foregroundColor(.primary)
                    }
                }
                .padding(.horizontal, 16)
                .padding(.vertical, 12)
                
                // Search bar
                HStack(spacing: 8) {
                    Image(systemName: "magnifyingglass")
                        .foregroundColor(.gray)
                    
                    TextField("Search", text: $viewModel.searchText)
                        .autocapitalization(.none)
                }
                .padding(.horizontal, 12)
                .padding(.vertical, 8)
                .background(Color.gray.opacity(0.1))
                .cornerRadius(10)
                .padding(.horizontal, 16)
                .padding(.bottom, 8)
                
                Divider()
                
                // Messages list
                if viewModel.conversations.isEmpty {
                    EmptyMessagesView()
                } else {
                    ScrollView {
                        LazyVStack(spacing: 0) {
                            ForEach(viewModel.conversations) { conversation in
                                NavigationLink(destination: ChatView(conversation: conversation)) {
                                    ConversationRow(conversation: conversation)
                                }
                                .buttonStyle(PlainButtonStyle())
                            }
                        }
                    }
                }
            }
            .sheet(isPresented: $showNewMessage) {
                NewMessageView()
            }
        }
    }
}

// MARK: - Conversation Row
struct ConversationRow: View {
    let conversation: Conversation
    
    var body: some View {
        HStack(spacing: 12) {
            // Avatar
            ZStack(alignment: .bottomTrailing) {
                Circle()
                    .fill(Color.gray.opacity(0.2))
                    .frame(width: 56, height: 56)
                    .overlay(
                        Image(systemName: "person.fill")
                            .foregroundColor(.gray)
                    )
                
                // Online indicator
                if conversation.isOnline {
                    Circle()
                        .fill(Color.green)
                        .frame(width: 14, height: 14)
                        .overlay(
                            Circle()
                                .stroke(Color.white, lineWidth: 2)
                        )
                }
            }
            
            // Content
            VStack(alignment: .leading, spacing: 4) {
                HStack {
                    Text(conversation.name)
                        .font(.system(size: 15, weight: conversation.hasUnread ? .semibold : .regular))
                    
                    Spacer()
                    
                    Text(conversation.timeAgo)
                        .font(.system(size: 13))
                        .foregroundColor(conversation.hasUnread ? .primary : .gray)
                }
                
                HStack(spacing: 4) {
                    if conversation.lastMessageIsYours {
                        Text("You:")
                            .font(.system(size: 14))
                            .foregroundColor(.gray)
                    }
                    
                    Text(conversation.lastMessage)
                        .font(.system(size: 14))
                        .foregroundColor(conversation.hasUnread ? .primary : .gray)
                        .lineLimit(2)
                    
                    Spacer()
                    
                    // Unread badge
                    if conversation.unreadCount > 0 {
                        Circle()
                            .fill(Color(hex: "007CFC"))
                            .frame(width: 8, height: 8)
                    }
                }
            }
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 8)
        .background(conversation.hasUnread ? Color(hex: "007CFC").opacity(0.05) : Color.clear)
        .contentShape(Rectangle())
    }
}

// MARK: - Chat View
struct ChatView: View {
    let conversation: Conversation
    @StateObject private var viewModel = ChatViewModel()
    @State private var messageText = ""
    @State private var showMediaPicker = false
    
    var body: some View {
        VStack(spacing: 0) {
            // Messages list
            ScrollViewReader { proxy in
                ScrollView {
                    LazyVStack(spacing: 12) {
                        ForEach(viewModel.messages) { message in
                            MessageBubble(message: message)
                                .id(message.id)
                        }
                    }
                    .padding(16)
                    .rotationEffect(.degrees(180))
                }
                .rotationEffect(.degrees(180))
                .onChange(of: viewModel.messages.count) { _, _ in
                    if let lastMessage = viewModel.messages.last {
                        proxy.scrollTo(lastMessage.id, anchor: .bottom)
                    }
                }
            }
            
            Divider()
            
            // Input bar
            HStack(spacing: 12) {
                Button(action: { showMediaPicker = true }) {
                    Image(systemName: "camera.fill")
                        .font(.system(size: 22))
                        .foregroundColor(Color(hex: "007CFC"))
                }
                
                Button(action: {}) {
                    Image(systemName: "photo")
                        .font(.system(size: 22))
                        .foregroundColor(Color(hex: "007CFC"))
                }
                
                // Text field
                HStack(spacing: 8) {
                    TextField("Message...", text: $messageText)
                    
                    if !messageText.isEmpty {
                        Button(action: {
                            viewModel.sendMessage(messageText)
                            messageText = ""
                        }) {
                            Image(systemName: "arrow.up.circle.fill")
                                .font(.system(size: 28))
                                .foregroundColor(Color(hex: "007CFC"))
                        }
                    } else {
                        Button(action: {}) {
                            Image(systemName: "mic.fill")
                                .font(.system(size: 20))
                                .foregroundColor(Color(hex: "007CFC"))
                        }
                        
                        Button(action: {}) {
                            Image(systemName: "face.smiling")
                                .font(.system(size: 20))
                                .foregroundColor(Color(hex: "007CFC"))
                        }
                    }
                }
                .padding(.horizontal, 12)
                .padding(.vertical, 8)
                .background(Color.gray.opacity(0.1))
                .cornerRadius(20)
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 8)
        }
        .navigationTitle(conversation.name)
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
            ToolbarItem(placement: .navigationBarTrailing) {
                HStack(spacing: 16) {
                    Button(action: {}) {
                        Image(systemName: "phone.fill")
                    }
                    
                    Button(action: {}) {
                        Image(systemName: "video.fill")
                    }
                    
                    Button(action: {}) {
                        Image(systemName: "info.circle")
                    }
                }
            }
        }
    }
}

// MARK: - Message Bubble
struct MessageBubble: View {
    let message: ChatMessage
    
    var body: some View {
        HStack {
            if message.isSender {
                Spacer()
            }
            
            VStack(alignment: message.isSender ? .trailing : .leading, spacing: 4) {
                // Message bubble
                Text(message.content)
                    .font(.system(size: 15))
                    .foregroundColor(message.isSender ? .white : .primary)
                    .padding(.horizontal, 12)
                    .padding(.vertical, 8)
                    .background(
                        message.isSender ?
                            Color(hex: "007CFC") :
                            Color.gray.opacity(0.15)
                    )
                    .cornerRadius(18)
                
                // Status
                HStack(spacing: 4) {
                    Text(message.timeAgo)
                        .font(.system(size: 11))
                        .foregroundColor(.gray)
                    
                    if message.isSender {
                        if message.isRead {
                            Text("Read")
                                .font(.system(size: 11))
                                .foregroundColor(.gray)
                        } else if message.isDelivered {
                            Text("Delivered")
                                .font(.system(size: 11))
                                .foregroundColor(.gray)
                        }
                        
                        // E2EE indicator
                        Image(systemName: "lock.fill")
                            .font(.system(size: 9))
                            .foregroundColor(.gray)
                    }
                }
            }
            
            if !message.isSender {
                Spacer()
            }
        }
    }
}

// MARK: - Empty Messages View
struct EmptyMessagesView: View {
    var body: some View {
        VStack(spacing: 16) {
            Image(systemName: "message")
                .font(.system(size: 64))
                .foregroundColor(.gray.opacity(0.5))
            
            Text("Your Messages")
                .font(.system(size: 22, weight: .bold))
            
            Text("Send private photos and messages to a friend or group")
                .font(.system(size: 14))
                .foregroundColor(.gray)
                .multilineTextAlignment(.center)
                .padding(.horizontal, 40)
            
            Button(action: {}) {
                Text("Send Message")
                    .font(.system(size: 15, weight: .semibold))
                    .foregroundColor(.white)
                    .frame(width: 200, height: 44)
                    .background(Color(hex: "007CFC"))
                    .cornerRadius(8)
            }
            .padding(.top, 8)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }
}

// MARK: - New Message View
struct NewMessageView: View {
    @Environment(\.dismiss) var dismiss
    @State private var searchText = ""
    @State private var selectedUsers: [String] = []
    
    var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Search bar
                HStack(spacing: 8) {
                    Image(systemName: "magnifyingglass")
                        .foregroundColor(.gray)
                    
                    TextField("Search...", text: $searchText)
                        .autocapitalization(.none)
                }
                .padding(.horizontal, 12)
                .padding(.vertical, 8)
                .background(Color.gray.opacity(0.1))
                .cornerRadius(10)
                .padding(16)
                
                // Suggested users
                ScrollView {
                    LazyVStack(spacing: 0) {
                        ForEach(0..<20) { index in
                            UserSelectRow(
                                username: "user\(index)",
                                fullName: "User Name \(index)",
                                isSelected: false,
                                onToggle: {}
                            )
                        }
                    }
                }
            }
            .navigationTitle("New Message")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Cancel") {
                        dismiss()
                    }
                }
                
                ToolbarItem(placement: .confirmationAction) {
                    Button("Next") {
                        // Create conversation
                        dismiss()
                    }
                    .fontWeight(.semibold)
                    .disabled(selectedUsers.isEmpty)
                }
            }
        }
    }
}

// MARK: - User Select Row
struct UserSelectRow: View {
    let username: String
    let fullName: String
    let isSelected: Bool
    let onToggle: () -> Void
    
    var body: some View {
        Button(action: onToggle) {
            HStack(spacing: 12) {
                ZStack(alignment: .bottomTrailing) {
                    Circle()
                        .fill(Color.gray.opacity(0.2))
                        .frame(width: 44, height: 44)
                        .overlay(
                            Image(systemName: "person.fill")
                                .foregroundColor(.gray)
                        )
                    
                    if isSelected {
                        Circle()
                            .fill(Color(hex: "007CFC"))
                            .frame(width: 20, height: 20)
                            .overlay(
                                Image(systemName: "checkmark")
                                    .font(.system(size: 12, weight: .bold))
                                    .foregroundColor(.white)
                            )
                    }
                }
                
                VStack(alignment: .leading, spacing: 4) {
                    Text(username)
                        .font(.system(size: 15, weight: .semibold))
                        .foregroundColor(.primary)
                    
                    Text(fullName)
                        .font(.system(size: 13))
                        .foregroundColor(.gray)
                }
                
                Spacer()
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 8)
        }
    }
}

// MARK: - Models
struct Conversation: Identifiable {
    let id: String
    let name: String
    let avatarUrl: String?
    let lastMessage: String
    let lastMessageIsYours: Bool
    let timeAgo: String
    let hasUnread: Bool
    let unreadCount: Int
    let isOnline: Bool
    let isEncrypted: Bool
}

struct ChatMessage: Identifiable {
    let id: String
    let content: String
    let isSender: Bool
    let timeAgo: String
    let isRead: Bool
    let isDelivered: Bool
    let timestamp: Date
}

// MARK: - View Models
class MessagesViewModel: ObservableObject {
    @Published var conversations: [Conversation] = []
    @Published var searchText: String = ""
    @Published var isLoading = false
    
    init() {
        loadMockData()
    }
    
    func loadMockData() {
        conversations = [
            Conversation(
                id: "1",
                name: "sarah_jones",
                avatarUrl: nil,
                lastMessage: "That sounds great! When are you free?",
                lastMessageIsYours: false,
                timeAgo: "2m",
                hasUnread: true,
                unreadCount: 2,
                isOnline: true,
                isEncrypted: true
            ),
            Conversation(
                id: "2",
                name: "mike_wilson",
                avatarUrl: nil,
                lastMessage: "Thanks for sharing! 🔥",
                lastMessageIsYours: true,
                timeAgo: "1h",
                hasUnread: false,
                unreadCount: 0,
                isOnline: false,
                isEncrypted: true
            ),
            Conversation(
                id: "3",
                name: "emma_davis",
                avatarUrl: nil,
                lastMessage: "See you tomorrow!",
                lastMessageIsYours: false,
                timeAgo: "3h",
                hasUnread: false,
                unreadCount: 0,
                isOnline: true,
                isEncrypted: true
            )
        ]
    }
}

class ChatViewModel: ObservableObject {
    @Published var messages: [ChatMessage] = []
    @Published var isTyping = false
    
    init() {
        loadMockMessages()
    }
    
    func loadMockMessages() {
        let now = Date()
        messages = [
            ChatMessage(
                id: "1",
                content: "Hey! How are you?",
                isSender: false,
                timeAgo: "10:30 AM",
                isRead: true,
                isDelivered: true,
                timestamp: now.addingTimeInterval(-3600)
            ),
            ChatMessage(
                id: "2",
                content: "I'm great! Just finished a new project 🎉",
                isSender: true,
                timeAgo: "10:32 AM",
                isRead: true,
                isDelivered: true,
                timestamp: now.addingTimeInterval(-3480)
            ),
            ChatMessage(
                id: "3",
                content: "That's awesome! Can you share some details?",
                isSender: false,
                timeAgo: "10:35 AM",
                isRead: true,
                isDelivered: true,
                timestamp: now.addingTimeInterval(-3300)
            ),
            ChatMessage(
                id: "4",
                content: "Sure! It's a photo editing app with lots of cool filters",
                isSender: true,
                timeAgo: "10:37 AM",
                isRead: false,
                isDelivered: true,
                timestamp: now.addingTimeInterval(-3180)
            )
        ]
    }
    
    func sendMessage(_ text: String) {
        let newMessage = ChatMessage(
            id: UUID().uuidString,
            content: text,
            isSender: true,
            timeAgo: "Now",
            isRead: false,
            isDelivered: false,
            timestamp: Date()
        )
        
        messages.append(newMessage)
        
        // TODO: Send via WebSocket with E2EE
    }
}
