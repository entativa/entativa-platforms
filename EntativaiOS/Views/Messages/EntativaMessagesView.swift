import SwiftUI

// MARK: - Entativa Messages View (Facebook Messenger-Style)
struct EntativaMessagesView: View {
    @StateObject private var viewModel = EntativaMessagesViewModel()
    @State private var showNewMessage = false
    @State private var selectedTab: MessengerTab = .chats
    
    var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Tab selector
                HStack(spacing: 0) {
                    MessengerTabButton(
                        title: "Chats",
                        tab: .chats,
                        selectedTab: $selectedTab,
                        count: viewModel.unreadChatsCount
                    )
                    
                    MessengerTabButton(
                        title: "Calls",
                        tab: .calls,
                        selectedTab: $selectedTab,
                        count: 0
                    )
                    
                    MessengerTabButton(
                        title: "People",
                        tab: .people,
                        selectedTab: $selectedTab,
                        count: 0
                    )
                }
                .frame(height: 44)
                
                Divider()
                
                // Content
                if selectedTab == .chats {
                    ChatsListView(viewModel: viewModel)
                } else if selectedTab == .calls {
                    CallsListView()
                } else {
                    PeopleListView()
                }
            }
            .navigationTitle("Messages")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button(action: {}) {
                        Image(systemName: "gearshape.fill")
                            .foregroundColor(.gray)
                    }
                }
                
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button(action: { showNewMessage = true }) {
                        Image(systemName: "square.and.pencil")
                            .foregroundColor(.primary)
                    }
                }
            }
            .sheet(isPresented: $showNewMessage) {
                EntativaNewMessageView()
            }
        }
    }
}

// MARK: - Messenger Tab Button
struct MessengerTabButton: View {
    let title: String
    let tab: MessengerTab
    @Binding var selectedTab: MessengerTab
    let count: Int
    
    var body: some View {
        Button(action: {
            selectedTab = tab
        }) {
            VStack(spacing: 8) {
                HStack(spacing: 4) {
                    Text(title)
                        .font(.system(size: 16, weight: selectedTab == tab ? .semibold : .regular))
                        .foregroundColor(selectedTab == tab ? .primary : .gray)
                    
                    if count > 0 {
                        Circle()
                            .fill(Color.red)
                            .frame(width: 8, height: 8)
                    }
                }
                .frame(maxWidth: .infinity)
                
                Rectangle()
                    .fill(selectedTab == tab ? Color(hex: "007CFC") : Color.clear)
                    .frame(height: 2)
            }
        }
    }
}

// MARK: - Chats List View
struct ChatsListView: View {
    @ObservedObject var viewModel: EntativaMessagesViewModel
    
    var body: some View {
        ScrollView {
            LazyVStack(spacing: 0) {
                // Search bar
                HStack(spacing: 8) {
                    Image(systemName: "magnifyingglass")
                        .foregroundColor(.gray)
                    
                    TextField("Search messages", text: $viewModel.searchText)
                        .autocapitalization(.none)
                }
                .padding(.horizontal, 12)
                .padding(.vertical, 8)
                .background(Color.gray.opacity(0.1))
                .cornerRadius(10)
                .padding(16)
                
                // Conversations
                if viewModel.conversations.isEmpty {
                    EntativaEmptyChatsView()
                } else {
                    ForEach(viewModel.conversations) { conversation in
                        NavigationLink(destination: EntativaChatView(conversation: conversation)) {
                            EntativaConversationRow(conversation: conversation)
                        }
                        .buttonStyle(PlainButtonStyle())
                    }
                }
            }
        }
    }
}

// MARK: - Conversation Row (Facebook-Style)
struct EntativaConversationRow: View {
    let conversation: EntativaConversation
    
    var body: some View {
        HStack(spacing: 12) {
            // Avatar with online indicator
            ZStack(alignment: .bottomTrailing) {
                Circle()
                    .fill(Color.gray.opacity(0.2))
                    .frame(width: 56, height: 56)
                    .overlay(
                        Image(systemName: "person.fill")
                            .foregroundColor(.gray)
                    )
                
                if conversation.isOnline {
                    Circle()
                        .fill(Color.green)
                        .frame(width: 16, height: 16)
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
                        .font(.system(size: 15, weight: conversation.hasUnread ? .bold : .semibold))
                    
                    // E2EE lock icon
                    Image(systemName: "lock.fill")
                        .font(.system(size: 10))
                        .foregroundColor(.gray)
                    
                    Spacer()
                    
                    HStack(spacing: 4) {
                        Text(conversation.timeAgo)
                            .font(.system(size: 13))
                            .foregroundColor(.gray)
                        
                        if conversation.unreadCount > 0 {
                            Circle()
                                .fill(Color(hex: "007CFC"))
                                .frame(width: 8, height: 8)
                        }
                    }
                }
                
                HStack {
                    Text(conversation.lastMessage)
                        .font(.system(size: 14))
                        .foregroundColor(conversation.hasUnread ? .primary : .gray)
                        .lineLimit(2)
                    
                    Spacer()
                }
            }
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 8)
        .background(conversation.hasUnread ? Color(hex: "007CFC").opacity(0.05) : Color.clear)
    }
}

// MARK: - Chat View (Messenger-Style)
struct EntativaChatView: View {
    let conversation: EntativaConversation
    @StateObject private var viewModel = EntativaChatViewModel()
    @State private var messageText = ""
    @State private var showOptions = false
    
    var body: some View {
        VStack(spacing: 0) {
            // Messages
            ScrollViewReader { proxy in
                ScrollView {
                    LazyVStack(spacing: 8) {
                        ForEach(viewModel.messages) { message in
                            EntativaMessageBubble(message: message)
                                .id(message.id)
                        }
                    }
                    .padding(16)
                    .rotationEffect(.degrees(180))
                }
                .rotationEffect(.degrees(180))
            }
            
            // E2EE indicator bar
            HStack(spacing: 4) {
                Image(systemName: "lock.shield.fill")
                    .font(.system(size: 11))
                    .foregroundColor(.gray)
                
                Text("End-to-end encrypted • Messages are secure")
                    .font(.system(size: 11))
                    .foregroundColor(.gray)
            }
            .padding(.vertical, 6)
            .background(Color.gray.opacity(0.05))
            
            Divider()
            
            // Input bar
            HStack(spacing: 8) {
                Button(action: {}) {
                    Image(systemName: "plus.circle.fill")
                        .font(.system(size: 28))
                        .foregroundColor(Color(hex: "007CFC"))
                }
                
                Button(action: {}) {
                    Image(systemName: "camera.fill")
                        .font(.system(size: 20))
                        .foregroundColor(Color(hex: "007CFC"))
                }
                
                Button(action: {}) {
                    Image(systemName: "photo.fill")
                        .font(.system(size: 20))
                        .foregroundColor(Color(hex: "007CFC"))
                }
                
                Button(action: {}) {
                    Image(systemName: "mic.fill")
                        .font(.system(size: 20))
                        .foregroundColor(Color(hex: "007CFC"))
                }
                
                // Text field
                HStack {
                    TextField("Message...", text: $messageText)
                    
                    if !messageText.isEmpty {
                        Button(action: {
                            viewModel.sendMessage(messageText)
                            messageText = ""
                        }) {
                            Text("Send")
                                .font(.system(size: 15, weight: .semibold))
                                .foregroundColor(Color(hex: "007CFC"))
                        }
                    }
                }
                .padding(.horizontal, 12)
                .padding(.vertical, 8)
                .background(Color.gray.opacity(0.1))
                .cornerRadius(20)
            }
            .padding(.horizontal, 12)
            .padding(.vertical, 8)
        }
        .navigationTitle(conversation.name)
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
            ToolbarItem(placement: .navigationBarTrailing) {
                HStack(spacing: 12) {
                    Button(action: {}) {
                        Image(systemName: "phone.fill")
                            .foregroundColor(Color(hex: "007CFC"))
                    }
                    
                    Button(action: {}) {
                        Image(systemName: "video.fill")
                            .foregroundColor(Color(hex: "007CFC"))
                    }
                    
                    Button(action: { showOptions = true }) {
                        Image(systemName: "info.circle")
                            .foregroundColor(.gray)
                    }
                }
            }
        }
        .sheet(isPresented: $showOptions) {
            ChatOptionsView(conversation: conversation)
        }
    }
}

// MARK: - Message Bubble (Messenger-Style)
struct EntativaMessageBubble: View {
    let message: EntativaMessage
    
    var body: some View {
        HStack {
            if message.isSender {
                Spacer()
            }
            
            VStack(alignment: message.isSender ? .trailing : .leading, spacing: 2) {
                // Message bubble
                Text(message.content)
                    .font(.system(size: 15))
                    .foregroundColor(message.isSender ? .white : .primary)
                    .padding(.horizontal, 14)
                    .padding(.vertical, 10)
                    .background(
                        message.isSender ?
                            LinearGradient(
                                colors: [
                                    Color(hex: "007CFC"),
                                    Color(hex: "6F3EFB")
                                ],
                                startPoint: .topLeading,
                                endPoint: .bottomTrailing
                            ) :
                            AnyShapeStyle(Color.gray.opacity(0.15))
                    )
                    .clipShape(RoundedRectangle(cornerRadius: 18))
                
                // Status row
                HStack(spacing: 4) {
                    if message.isSender {
                        // E2EE indicator
                        Image(systemName: "lock.fill")
                            .font(.system(size: 8))
                            .foregroundColor(.gray)
                        
                        if message.isRead {
                            Image(systemName: "checkmark.circle.fill")
                                .font(.system(size: 12))
                                .foregroundColor(Color(hex: "007CFC"))
                        } else if message.isDelivered {
                            Image(systemName: "checkmark.circle")
                                .font(.system(size: 12))
                                .foregroundColor(.gray)
                        }
                    }
                    
                    Text(message.timeAgo)
                        .font(.system(size: 11))
                        .foregroundColor(.gray)
                }
            }
            
            if !message.isSender {
                Spacer()
            }
        }
    }
}

// MARK: - Empty Chats View
struct EntativaEmptyChatsView: View {
    var body: some View {
        VStack(spacing: 16) {
            Image(systemName: "bubble.left.and.bubble.right")
                .font(.system(size: 64))
                .foregroundColor(.gray.opacity(0.5))
            
            Text("No messages yet")
                .font(.system(size: 20, weight: .semibold))
            
            Text("Start a conversation with your friends")
                .font(.system(size: 14))
                .foregroundColor(.gray)
                .multilineTextAlignment(.center)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .padding(.top, 100)
    }
}

// MARK: - Calls List View
struct CallsListView: View {
    var body: some View {
        ScrollView {
            LazyVStack(spacing: 0) {
                ForEach(0..<10) { index in
                    CallHistoryRow(
                        name: "User \(index)",
                        isVideo: index % 2 == 0,
                        isMissed: index % 3 == 0,
                        timeAgo: "2 hours ago"
                    )
                }
            }
        }
    }
}

// MARK: - Call History Row
struct CallHistoryRow: View {
    let name: String
    let isVideo: Bool
    let isMissed: Bool
    let timeAgo: String
    
    var body: some View {
        HStack(spacing: 12) {
            Circle()
                .fill(Color.gray.opacity(0.2))
                .frame(width: 48, height: 48)
            
            VStack(alignment: .leading, spacing: 4) {
                Text(name)
                    .font(.system(size: 15, weight: .semibold))
                
                HStack(spacing: 4) {
                    Image(systemName: isMissed ? "phone.down.fill" : "phone.fill")
                        .font(.system(size: 12))
                        .foregroundColor(isMissed ? .red : .gray)
                    
                    Text(timeAgo)
                        .font(.system(size: 13))
                        .foregroundColor(.gray)
                }
            }
            
            Spacer()
            
            Button(action: {}) {
                Image(systemName: isVideo ? "video.fill" : "phone.fill")
                    .font(.system(size: 20))
                    .foregroundColor(Color(hex: "007CFC"))
            }
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 8)
    }
}

// MARK: - People List View
struct PeopleListView: View {
    var body: some View {
        ScrollView {
            LazyVStack(spacing: 0) {
                ForEach(0..<15) { index in
                    PersonRow(name: "Friend \(index)")
                }
            }
        }
    }
}

// MARK: - Person Row
struct PersonRow: View {
    let name: String
    
    var body: some View {
        HStack(spacing: 12) {
            Circle()
                .fill(Color.gray.opacity(0.2))
                .frame(width: 48, height: 48)
            
            Text(name)
                .font(.system(size: 15, weight: .semibold))
            
            Spacer()
            
            Button(action: {}) {
                Image(systemName: "message.fill")
                    .font(.system(size: 18))
                    .foregroundColor(Color(hex: "007CFC"))
            }
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 8)
    }
}

// MARK: - Chat Options View
struct ChatOptionsView: View {
    let conversation: EntativaConversation
    @Environment(\.dismiss) var dismiss
    
    var body: some View {
        NavigationView {
            List {
                Section {
                    HStack {
                        Spacer()
                        
                        VStack(spacing: 12) {
                            Circle()
                                .fill(Color.gray.opacity(0.2))
                                .frame(width: 80, height: 80)
                            
                            Text(conversation.name)
                                .font(.system(size: 18, weight: .bold))
                        }
                        
                        Spacer()
                    }
                    .padding(.vertical, 16)
                }
                
                Section {
                    NavigationLink(destination: Text("Search in Conversation")) {
                        Label("Search in Conversation", systemImage: "magnifyingglass")
                    }
                    
                    NavigationLink(destination: Text("Media, Links & Docs")) {
                        Label("Media, Links & Docs", systemImage: "photo.on.rectangle")
                    }
                }
                
                Section {
                    Toggle(isOn: .constant(true)) {
                        Label("Notifications", systemImage: "bell.fill")
                    }
                    
                    NavigationLink(destination: Text("Disappearing Messages")) {
                        Label("Disappearing Messages", systemImage: "timer")
                    }
                }
                
                Section {
                    Button(action: {}) {
                        Label("Block", systemImage: "hand.raised.fill")
                            .foregroundColor(.red)
                    }
                    
                    Button(action: {}) {
                        Label("Delete Chat", systemImage: "trash")
                            .foregroundColor(.red)
                    }
                }
            }
            .navigationTitle("Chat Settings")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Done") {
                        dismiss()
                    }
                }
            }
        }
    }
}

// MARK: - New Message View
struct EntativaNewMessageView: View {
    @Environment(\.dismiss) var dismiss
    @State private var searchText = ""
    
    var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Search
                HStack(spacing: 8) {
                    Image(systemName: "magnifyingglass")
                        .foregroundColor(.gray)
                    
                    TextField("Search for people...", text: $searchText)
                        .autocapitalization(.none)
                }
                .padding(.horizontal, 12)
                .padding(.vertical, 8)
                .background(Color.gray.opacity(0.1))
                .cornerRadius(10)
                .padding(16)
                
                // Suggested people
                ScrollView {
                    LazyVStack(spacing: 0) {
                        ForEach(0..<15) { index in
                            Button(action: {}) {
                                HStack(spacing: 12) {
                                    Circle()
                                        .fill(Color.gray.opacity(0.2))
                                        .frame(width: 48, height: 48)
                                    
                                    VStack(alignment: .leading, spacing: 4) {
                                        Text("User \(index)")
                                            .font(.system(size: 15, weight: .semibold))
                                            .foregroundColor(.primary)
                                        
                                        Text("@username\(index)")
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
            }
        }
    }
}

// MARK: - Messenger Tab Enum
enum MessengerTab {
    case chats
    case calls
    case people
}

// MARK: - Models
struct EntativaConversation: Identifiable {
    let id: String
    let name: String
    let avatarUrl: String?
    let lastMessage: String
    let timeAgo: String
    let hasUnread: Bool
    let unreadCount: Int
    let isOnline: Bool
    let isEncrypted: Bool
}

struct EntativaMessage: Identifiable {
    let id: String
    let content: String
    let isSender: Bool
    let timeAgo: String
    let isRead: Bool
    let isDelivered: Bool
    let timestamp: Date
}

// MARK: - View Models
class EntativaMessagesViewModel: ObservableObject {
    @Published var conversations: [EntativaConversation] = []
    @Published var searchText: String = ""
    @Published var unreadChatsCount = 0
    
    init() {
        loadMockData()
    }
    
    func loadMockData() {
        conversations = [
            EntativaConversation(
                id: "1",
                name: "Sarah Johnson",
                avatarUrl: nil,
                lastMessage: "You: Sounds good! See you then 👍",
                timeAgo: "5m",
                hasUnread: true,
                unreadCount: 3,
                isOnline: true,
                isEncrypted: true
            ),
            EntativaConversation(
                id: "2",
                name: "Mike Wilson",
                avatarUrl: nil,
                lastMessage: "That's perfect, thanks!",
                timeAgo: "2h",
                hasUnread: false,
                unreadCount: 0,
                isOnline: false,
                isEncrypted: true
            )
        ]
        
        unreadChatsCount = conversations.filter { $0.hasUnread }.count
    }
    
    func clearRecentSearches() {
        // TODO: Clear searches
    }
}

class EntativaChatViewModel: ObservableObject {
    @Published var messages: [EntativaMessage] = []
    @Published var isTyping = false
    
    init() {
        loadMockMessages()
    }
    
    func loadMockMessages() {
        let now = Date()
        messages = [
            EntativaMessage(
                id: "1",
                content: "Hey! How's it going?",
                isSender: false,
                timeAgo: "10:30 AM",
                isRead: true,
                isDelivered: true,
                timestamp: now.addingTimeInterval(-3600)
            ),
            EntativaMessage(
                id: "2",
                content: "Going great! Just wrapped up that project we discussed",
                isSender: true,
                timeAgo: "10:32 AM",
                isRead: true,
                isDelivered: true,
                timestamp: now.addingTimeInterval(-3480)
            ),
            EntativaMessage(
                id: "3",
                content: "Amazing! Would love to see it",
                isSender: false,
                timeAgo: "10:35 AM",
                isRead: true,
                isDelivered: true,
                timestamp: now.addingTimeInterval(-3300)
            )
        ]
    }
    
    func sendMessage(_ text: String) {
        let newMessage = EntativaMessage(
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
