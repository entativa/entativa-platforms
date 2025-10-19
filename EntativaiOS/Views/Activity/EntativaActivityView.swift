import SwiftUI

// MARK: - Entativa Activity View (Facebook-Style)
struct EntativaActivityView: View {
    @StateObject private var viewModel = EntativaActivityViewModel()
    
    var body: some View {
        NavigationView {
            ScrollView {
                LazyVStack(spacing: 0) {
                    // New notifications
                    if !viewModel.newActivities.isEmpty {
                        SectionHeaderView(title: "New")
                        
                        ForEach(viewModel.newActivities) { activity in
                            EntativaActivityRow(activity: activity)
                        }
                        
                        Divider()
                            .padding(.vertical, 8)
                    }
                    
                    // Earlier notifications
                    if !viewModel.earlierActivities.isEmpty {
                        SectionHeaderView(title: "Earlier")
                        
                        ForEach(viewModel.earlierActivities) { activity in
                            EntativaActivityRow(activity: activity)
                        }
                    }
                }
            }
            .background(Color(UIColor.systemGroupedBackground))
            .navigationTitle("Notifications")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button(action: {}) {
                        Image(systemName: "ellipsis")
                            .foregroundColor(.primary)
                    }
                }
            }
        }
    }
}

// MARK: - Section Header View
struct SectionHeaderView: View {
    let title: String
    
    var body: some View {
        HStack {
            Text(title)
                .font(.system(size: 17, weight: .semibold))
                .foregroundColor(.primary)
            
            Spacer()
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 12)
        .background(Color(UIColor.systemGroupedBackground))
    }
}

// MARK: - Entativa Activity Row
struct EntativaActivityRow: View {
    let activity: EntativaActivity
    @State private var showMenu = false
    
    var body: some View {
        Button(action: {
            // Navigate to relevant content
        }) {
            HStack(alignment: .top, spacing: 12) {
                // Icon with colored background
                ZStack {
                    Circle()
                        .fill(activity.iconBackgroundColor)
                        .frame(width: 56, height: 56)
                    
                    if let avatarUrl = activity.userAvatar {
                        AsyncImage(url: URL(string: avatarUrl)) { image in
                            image
                                .resizable()
                                .aspectRatio(contentMode: .fill)
                        } placeholder: {
                            Image(systemName: activity.iconName)
                                .font(.system(size: 24))
                                .foregroundColor(.white)
                        }
                        .frame(width: 52, height: 52)
                        .clipShape(Circle())
                    } else {
                        Image(systemName: activity.iconName)
                            .font(.system(size: 24))
                            .foregroundColor(.white)
                    }
                    
                    // Badge for specific types
                    if activity.showBadge {
                        Circle()
                            .fill(Color.red)
                            .frame(width: 16, height: 16)
                            .overlay(
                                Image(systemName: "exclamationmark")
                                    .font(.system(size: 10, weight: .bold))
                                    .foregroundColor(.white)
                            )
                            .offset(x: 20, y: -20)
                    }
                }
                
                // Content
                VStack(alignment: .leading, spacing: 4) {
                    Text(activity.text)
                        .font(.system(size: 15))
                        .foregroundColor(.primary)
                        .lineLimit(3)
                    
                    Text(activity.timeAgo)
                        .font(.system(size: 13))
                        .foregroundColor(.gray)
                    
                    // Action buttons for specific types
                    if activity.type == .friendRequest {
                        HStack(spacing: 12) {
                            Button(action: {}) {
                                Text("Confirm")
                                    .font(.system(size: 15, weight: .semibold))
                                    .foregroundColor(.white)
                                    .frame(maxWidth: .infinity)
                                    .frame(height: 36)
                                    .background(Color(hex: "007CFC"))
                                    .cornerRadius(8)
                            }
                            
                            Button(action: {}) {
                                Text("Delete")
                                    .font(.system(size: 15, weight: .semibold))
                                    .foregroundColor(.primary)
                                    .frame(maxWidth: .infinity)
                                    .frame(height: 36)
                                    .background(Color.gray.opacity(0.2))
                                    .cornerRadius(8)
                            }
                        }
                        .padding(.top, 8)
                    }
                }
                
                Spacer()
                
                // Post thumbnail if applicable
                if let thumbnail = activity.postThumbnail {
                    AsyncImage(url: URL(string: thumbnail)) { image in
                        image
                            .resizable()
                            .aspectRatio(contentMode: .fill)
                    } placeholder: {
                        Rectangle()
                            .fill(Color.gray.opacity(0.2))
                    }
                    .frame(width: 64, height: 64)
                    .cornerRadius(8)
                }
                
                // Menu button
                Button(action: { showMenu = true }) {
                    Image(systemName: "ellipsis")
                        .font(.system(size: 16))
                        .foregroundColor(.gray)
                        .frame(width: 32, height: 32)
                }
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)
            .background(activity.isRead ? Color(UIColor.systemBackground) : Color(hex: "007CFC").opacity(0.08))
        }
        .buttonStyle(PlainButtonStyle())
    }
}

// MARK: - Models
enum EntativaActivityType {
    case like
    case comment
    case share
    case friendRequest
    case friendAccepted
    case tag
    case mention
    case event
    case birthday
    case memory
}

struct EntativaActivity: Identifiable {
    let id: String
    let type: EntativaActivityType
    let text: String
    let timeAgo: String
    let userAvatar: String?
    let postThumbnail: String?
    let isRead: Bool
    let showBadge: Bool
    let timestamp: Date
    
    var iconName: String {
        switch type {
        case .like:
            return "heart.fill"
        case .comment:
            return "bubble.left.fill"
        case .share:
            return "arrow.turn.up.right"
        case .friendRequest:
            return "person.fill"
        case .friendAccepted:
            return "person.2.fill"
        case .tag:
            return "tag.fill"
        case .mention:
            return "at"
        case .event:
            return "calendar"
        case .birthday:
            return "gift.fill"
        case .memory:
            return "clock.fill"
        }
    }
    
    var iconBackgroundColor: Color {
        switch type {
        case .like:
            return Color.red
        case .comment:
            return Color(hex: "007CFC")
        case .share:
            return Color.green
        case .friendRequest, .friendAccepted:
            return Color(hex: "007CFC")
        case .tag:
            return Color.orange
        case .mention:
            return Color(hex: "6F3EFB")
        case .event:
            return Color.red
        case .birthday:
            return Color.pink
        case .memory:
            return Color(hex: "6F3EFB")
        }
    }
}

// MARK: - View Model
class EntativaActivityViewModel: ObservableObject {
    @Published var newActivities: [EntativaActivity] = []
    @Published var earlierActivities: [EntativaActivity] = []
    @Published var isLoading = false
    
    init() {
        loadMockData()
    }
    
    func loadMockData() {
        let now = Date()
        
        // New notifications
        newActivities = [
            EntativaActivity(
                id: "1",
                type: .friendRequest,
                text: "Sarah Johnson sent you a friend request.",
                timeAgo: "2 hours ago",
                userAvatar: nil,
                postThumbnail: nil,
                isRead: false,
                showBadge: true,
                timestamp: now.addingTimeInterval(-2 * 3600)
            ),
            EntativaActivity(
                id: "2",
                type: .like,
                text: "Mike Wilson and 12 others reacted to your post.",
                timeAgo: "4 hours ago",
                userAvatar: nil,
                postThumbnail: "",
                isRead: false,
                showBadge: false,
                timestamp: now.addingTimeInterval(-4 * 3600)
            ),
            EntativaActivity(
                id: "3",
                type: .comment,
                text: "Emma Davis commented on your photo: \"This is amazing! 🔥\"",
                timeAgo: "6 hours ago",
                userAvatar: nil,
                postThumbnail: "",
                isRead: false,
                showBadge: false,
                timestamp: now.addingTimeInterval(-6 * 3600)
            ),
            EntativaActivity(
                id: "4",
                type: .birthday,
                text: "It's Alex Chen's birthday today! Write on their timeline.",
                timeAgo: "8 hours ago",
                userAvatar: nil,
                postThumbnail: nil,
                isRead: true,
                showBadge: false,
                timestamp: now.addingTimeInterval(-8 * 3600)
            )
        ]
        
        // Earlier notifications
        earlierActivities = [
            EntativaActivity(
                id: "5",
                type: .friendAccepted,
                text: "Chris Taylor accepted your friend request.",
                timeAgo: "Yesterday",
                userAvatar: nil,
                postThumbnail: nil,
                isRead: true,
                showBadge: false,
                timestamp: now.addingTimeInterval(-24 * 3600)
            ),
            EntativaActivity(
                id: "6",
                type: .tag,
                text: "You were tagged in a photo by Jessica Brown.",
                timeAgo: "2 days ago",
                userAvatar: nil,
                postThumbnail: "",
                isRead: true,
                showBadge: false,
                timestamp: now.addingTimeInterval(-2 * 86400)
            ),
            EntativaActivity(
                id: "7",
                type: .memory,
                text: "We found a memory from 3 years ago that you might like.",
                timeAgo: "3 days ago",
                userAvatar: nil,
                postThumbnail: "",
                isRead: true,
                showBadge: false,
                timestamp: now.addingTimeInterval(-3 * 86400)
            ),
            EntativaActivity(
                id: "8",
                type: .event,
                text: "Summer Festival is happening tomorrow. Are you going?",
                timeAgo: "4 days ago",
                userAvatar: nil,
                postThumbnail: nil,
                isRead: true,
                showBadge: false,
                timestamp: now.addingTimeInterval(-4 * 86400)
            )
        ]
    }
    
    func loadActivities() async {
        // TODO: Load from API
    }
    
    func markAsRead(activityId: String) {
        // TODO: Mark as read via API
    }
}
