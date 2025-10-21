import SwiftUI

// MARK: - Vignette Activity View (Instagram-Style)
struct VignetteActivityView: View {
    @StateObject private var viewModel = ActivityViewModel()
    @State private var selectedTab: ActivityTab = .you
    
    var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Tab selector
                HStack(spacing: 0) {
                    TabButton(
                        title: "Following",
                        isSelected: selectedTab == .following,
                        action: { selectedTab = .following }
                    )
                    
                    TabButton(
                        title: "You",
                        isSelected: selectedTab == .you,
                        action: { selectedTab = .you }
                    )
                }
                .frame(height: 44)
                .background(Color(UIColor.systemBackground))
                
                Divider()
                
                // Activity list
                ScrollView {
                    LazyVStack(spacing: 0) {
                        if selectedTab == .you {
                            // Today section
                            if !viewModel.todayActivities.isEmpty {
                                SectionHeader(title: "Today")
                                
                                ForEach(viewModel.todayActivities) { activity in
                                    ActivityRow(activity: activity)
                                }
                            }
                            
                            // This Week section
                            if !viewModel.weekActivities.isEmpty {
                                SectionHeader(title: "This Week")
                                
                                ForEach(viewModel.weekActivities) { activity in
                                    ActivityRow(activity: activity)
                                }
                            }
                            
                            // This Month section
                            if !viewModel.monthActivities.isEmpty {
                                SectionHeader(title: "This Month")
                                
                                ForEach(viewModel.monthActivities) { activity in
                                    ActivityRow(activity: activity)
                                }
                            }
                            
                            // Earlier section
                            if !viewModel.earlierActivities.isEmpty {
                                SectionHeader(title: "Earlier")
                                
                                ForEach(viewModel.earlierActivities) { activity in
                                    ActivityRow(activity: activity)
                                }
                            }
                        } else {
                            // Following tab
                            ForEach(viewModel.followingActivities) { activity in
                                ActivityRow(activity: activity)
                            }
                        }
                    }
                    .padding(.vertical, 8)
                }
                .background(Color(UIColor.systemBackground))
            }
            .navigationTitle("Notifications")
            .navigationBarTitleDisplayMode(.inline)
        }
    }
}

// MARK: - Tab Button
struct TabButton: View {
    let title: String
    let isSelected: Bool
    let action: () -> Void
    
    var body: some View {
        Button(action: action) {
            VStack(spacing: 8) {
                Text(title)
                    .font(.system(size: 16, weight: isSelected ? .semibold : .regular))
                    .foregroundColor(isSelected ? .primary : .gray)
                    .frame(maxWidth: .infinity)
                
                Rectangle()
                    .fill(isSelected ? Color.primary : Color.clear)
                    .frame(height: 1)
            }
        }
    }
}

// MARK: - Section Header
struct SectionHeader: View {
    let title: String
    
    var body: some View {
        HStack {
            Text(title)
                .font(.system(size: 15, weight: .semibold))
                .foregroundColor(.primary)
            
            Spacer()
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 8)
        .background(Color(UIColor.systemBackground))
    }
}

// MARK: - Activity Row
struct ActivityRow: View {
    let activity: Activity
    @State private var isFollowing: Bool
    
    init(activity: Activity) {
        self.activity = activity
        _isFollowing = State(initialValue: activity.isFollowing)
    }
    
    var body: some View {
        HStack(alignment: .center, spacing: 12) {
            // Profile picture
            AsyncImage(url: URL(string: activity.userAvatar ?? "")) { image in
                image
                    .resizable()
                    .aspectRatio(contentMode: .fill)
            } placeholder: {
                Circle()
                    .fill(Color.gray.opacity(0.2))
                    .overlay(
                        Image(systemName: "person.fill")
                            .foregroundColor(.gray.opacity(0.5))
                    )
            }
            .frame(width: 44, height: 44)
            .clipShape(Circle())
            
            // Activity text and time
            VStack(alignment: .leading, spacing: 4) {
                Text(attributedActivityText)
                    .font(.system(size: 14))
                    .lineLimit(3)
                
                Text(activity.timeAgo)
                    .font(.system(size: 12))
                    .foregroundColor(.gray)
            }
            
            Spacer()
            
            // Right side (post thumbnail or follow button)
            if let postThumbnail = activity.postThumbnail {
                AsyncImage(url: URL(string: postThumbnail)) { image in
                    image
                        .resizable()
                        .aspectRatio(contentMode: .fill)
                } placeholder: {
                    Rectangle()
                        .fill(Color.gray.opacity(0.2))
                }
                .frame(width: 44, height: 44)
                .cornerRadius(4)
            } else if activity.type == .follow || activity.type == .followRequest {
                Button(action: {
                    isFollowing.toggle()
                }) {
                    Text(isFollowing ? "Following" : "Follow")
                        .font(.system(size: 14, weight: .semibold))
                        .foregroundColor(isFollowing ? .primary : .white)
                        .frame(width: 100, height: 32)
                        .background(isFollowing ? Color.gray.opacity(0.2) : Color(hex: "007CFC"))
                        .cornerRadius(8)
                }
            }
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 8)
        .background(activity.isRead ? Color.clear : Color(hex: "007CFC").opacity(0.05))
    }
    
    var attributedActivityText: AttributedString {
        var text = AttributedString(activity.username)
        text.font = .system(size: 14, weight: .semibold)
        
        var action = AttributedString(" \(activity.action)")
        action.font = .system(size: 14)
        
        return text + action
    }
}

// MARK: - Models
enum ActivityTab {
    case following
    case you
}

enum ActivityType {
    case like
    case comment
    case follow
    case followRequest
    case mention
    case tag
    case reply
}

struct Activity: Identifiable {
    let id: String
    let type: ActivityType
    let username: String
    let userAvatar: String?
    let action: String
    let postThumbnail: String?
    let timeAgo: String
    let isRead: Bool
    let isFollowing: Bool
    let timestamp: Date
}

// MARK: - View Model
class ActivityViewModel: ObservableObject {
    @Published var todayActivities: [Activity] = []
    @Published var weekActivities: [Activity] = []
    @Published var monthActivities: [Activity] = []
    @Published var earlierActivities: [Activity] = []
    @Published var followingActivities: [Activity] = []
    @Published var isLoading = false
    
    init() {
        loadMockData()
    }
    
    func loadMockData() {
        let now = Date()
        
        // Today
        todayActivities = [
            Activity(
                id: "1",
                type: .like,
                username: "sarah_jones",
                userAvatar: nil,
                action: "liked your photo.",
                postThumbnail: "",
                timeAgo: "2h",
                isRead: false,
                isFollowing: true,
                timestamp: now.addingTimeInterval(-2 * 3600)
            ),
            Activity(
                id: "2",
                type: .comment,
                username: "mike_wilson",
                userAvatar: nil,
                action: "commented: \"Amazing shot! 🔥\"",
                postThumbnail: "",
                timeAgo: "4h",
                isRead: false,
                isFollowing: true,
                timestamp: now.addingTimeInterval(-4 * 3600)
            ),
            Activity(
                id: "3",
                type: .follow,
                username: "alex_creative",
                userAvatar: nil,
                action: "started following you.",
                postThumbnail: nil,
                timeAgo: "5h",
                isRead: true,
                isFollowing: false,
                timestamp: now.addingTimeInterval(-5 * 3600)
            )
        ]
        
        // This Week
        weekActivities = [
            Activity(
                id: "4",
                type: .like,
                username: "emma_davis",
                userAvatar: nil,
                action: "liked your photo.",
                postThumbnail: "",
                timeAgo: "2d",
                isRead: true,
                isFollowing: true,
                timestamp: now.addingTimeInterval(-2 * 86400)
            ),
            Activity(
                id: "5",
                type: .mention,
                username: "chris_photo",
                userAvatar: nil,
                action: "mentioned you in a comment: @you Check this out!",
                postThumbnail: "",
                timeAgo: "3d",
                isRead: true,
                isFollowing: true,
                timestamp: now.addingTimeInterval(-3 * 86400)
            )
        ]
        
        // This Month
        monthActivities = [
            Activity(
                id: "6",
                type: .follow,
                username: "photographer_pro",
                userAvatar: nil,
                action: "started following you.",
                postThumbnail: nil,
                timeAgo: "1w",
                isRead: true,
                isFollowing: false,
                timestamp: now.addingTimeInterval(-7 * 86400)
            ),
            Activity(
                id: "7",
                type: .like,
                username: "travel_enthusiast",
                userAvatar: nil,
                action: "liked 3 of your photos.",
                postThumbnail: "",
                timeAgo: "2w",
                isRead: true,
                isFollowing: true,
                timestamp: now.addingTimeInterval(-14 * 86400)
            )
        ]
        
        // Earlier
        earlierActivities = [
            Activity(
                id: "8",
                type: .comment,
                username: "design_lover",
                userAvatar: nil,
                action: "commented: \"Love your style!\"",
                postThumbnail: "",
                timeAgo: "3w",
                isRead: true,
                isFollowing: true,
                timestamp: now.addingTimeInterval(-21 * 86400)
            )
        ]
        
        // Following tab (recent activity from people you follow)
        followingActivities = [
            Activity(
                id: "9",
                type: .like,
                username: "sarah_jones",
                userAvatar: nil,
                action: "liked a photo by mike_wilson.",
                postThumbnail: "",
                timeAgo: "1h",
                isRead: true,
                isFollowing: true,
                timestamp: now.addingTimeInterval(-1 * 3600)
            ),
            Activity(
                id: "10",
                type: .follow,
                username: "mike_wilson",
                userAvatar: nil,
                action: "started following alex_creative.",
                postThumbnail: nil,
                timeAgo: "3h",
                isRead: true,
                isFollowing: true,
                timestamp: now.addingTimeInterval(-3 * 3600)
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
