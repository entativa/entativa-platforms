import SwiftUI

struct VignetteHomeView: View {
    @StateObject private var viewModel = VignetteHomeViewModel()
    @State private var showingCamera = false
    @State private var showingSearch = false
    @State private var selectedTab: VignetteTab = .home
    
    var body: some View {
        ZStack(alignment: .bottom) {
            // Main content
            VStack(spacing: 0) {
                // Top bar
                VignetteTopBar(
                    onPlusAction: {
                        showingCamera = true
                    },
                    onSearchAction: {
                        showingSearch = true
                    }
                )
                
                // Content based on selected tab
                TabView(selection: $selectedTab) {
                    // Home Feed
                    VignetteFeedView(viewModel: viewModel)
                        .tag(VignetteTab.home)
                    
                    // Takes (Reels)
                    VignetteTakesView()
                        .tag(VignetteTab.takes)
                    
                    // Messages (Direct)
                    VignetteDirectView()
                        .tag(VignetteTab.messages)
                    
                    // Activity (Notifications)
                    VignetteActivityView()
                        .tag(VignetteTab.activity)
                    
                    // Profile
                    VignetteProfileView()
                        .tag(VignetteTab.profile)
                }
                .tabViewStyle(.page(indexDisplayMode: .never))
            }
            
            // Floating liquid glass bottom navigation
            VignetteBottomNavBar(selectedTab: $selectedTab)
                .padding(.bottom, 8)
                .padding(.horizontal, 16)
        }
        .ignoresSafeArea(.keyboard)
        .fullScreenCover(isPresented: $showingCamera) {
            VignetteCameraView()
        }
        .sheet(isPresented: $showingSearch) {
            VignetteSearchView()
        }
    }
}

// MARK: - Top Bar Component
struct VignetteTopBar: View {
    let onPlusAction: () -> Void
    let onSearchAction: () -> Void
    
    var body: some View {
        HStack(spacing: 16) {
            // Plus button (left)
            Button(action: onPlusAction) {
                Image(systemName: "plus.app")
                    .font(.system(size: 26, weight: .regular))
                    .foregroundColor(VignetteColors.textPrimary)
            }
            
            Spacer()
            
            // Vignette logo (center)
            Text("Vignette")
                .font(.custom("Snell Roundhand", size: 32))
                .italic()
                .foregroundColor(VignetteColors.textPrimary)
            
            Spacer()
            
            // Search button (right)
            Button(action: onSearchAction) {
                Image(systemName: "magnifyingglass")
                    .font(.system(size: 22, weight: .regular))
                    .foregroundColor(VignetteColors.textPrimary)
            }
        }
        .padding(.horizontal, 20)
        .padding(.vertical, 12)
        .background(
            Color(UIColor.systemBackground)
                .shadow(color: .black.opacity(0.03), radius: 1, y: 1)
        )
    }
}

// MARK: - Bottom Navigation Bar (Liquid Glass Effect)
struct VignetteBottomNavBar: View {
    @Binding var selectedTab: VignetteTab
    
    var body: some View {
        HStack(spacing: 0) {
            ForEach(VignetteTab.allCases, id: \.self) { tab in
                Button(action: {
                    withAnimation(.spring(response: 0.3, dampingFraction: 0.7)) {
                        selectedTab = tab
                    }
                }) {
                    VStack(spacing: 0) {
                        Image(systemName: tab.icon)
                            .font(.system(size: 25, weight: selectedTab == tab ? .semibold : .regular))
                            .foregroundColor(selectedTab == tab ? VignetteColors.textPrimary : VignetteColors.textSecondary)
                            .frame(height: 32)
                    }
                    .frame(maxWidth: .infinity)
                    .padding(.vertical, 10)
                }
            }
        }
        .padding(.horizontal, 8)
        .background(
            // Liquid glass effect
            ZStack {
                // Blur background
                Color(UIColor.systemBackground)
                    .opacity(0.85)
                
                // Glass effect
                Rectangle()
                    .fill(.ultraThinMaterial)
            }
            .cornerRadius(24)
            .shadow(color: .black.opacity(0.08), radius: 16, y: 8)
        )
    }
}

// MARK: - Feed View (Instagram-style)
struct VignetteFeedView: View {
    @ObservedObject var viewModel: VignetteHomeViewModel
    
    var body: some View {
        ScrollView {
            LazyVStack(spacing: 0) {
                // Stories Row (Instagram circular style)
                VignetteStoriesRow()
                    .padding(.bottom, 8)
                
                Divider()
                    .padding(.bottom, 8)
                
                // Posts (Instagram-style single posts)
                ForEach(viewModel.posts) { post in
                    VignettePostCard(post: post)
                        .padding(.bottom, 12)
                    
                    Divider()
                        .padding(.horizontal, 16)
                }
            }
        }
        .refreshable {
            await viewModel.refreshFeed()
        }
    }
}

// MARK: - Stories Row (Circular Instagram-style)
struct VignetteStoriesRow: View {
    @State private var stories: [VignetteStory] = VignetteStory.mockStories
    
    var body: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: 16) {
                // Your Story
                VStack(spacing: 6) {
                    ZStack(alignment: .bottomTrailing) {
                        Circle()
                            .fill(
                                LinearGradient(
                                    colors: [
                                        VignetteColors.moonstone.opacity(0.3),
                                        VignetteColors.saffron.opacity(0.3)
                                    ],
                                    startPoint: .topLeading,
                                    endPoint: .bottomTrailing
                                )
                            )
                            .frame(width: 72, height: 72)
                        
                        Circle()
                            .fill(EntativaColors.primaryBlue)
                            .frame(width: 24, height: 24)
                            .overlay(
                                Image(systemName: "plus")
                                    .font(.system(size: 14, weight: .semibold))
                                    .foregroundColor(.white)
                            )
                            .overlay(
                                Circle()
                                    .stroke(Color.white, lineWidth: 2)
                            )
                    }
                    
                    Text("Your story")
                        .vignetteCaptionSmall()
                        .foregroundColor(VignetteColors.textPrimary)
                }
                
                // Friend Stories
                ForEach(stories) { story in
                    VignetteStoryCircle(story: story)
                }
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 8)
        }
    }
}

struct VignetteStoryCircle: View {
    let story: VignetteStory
    
    var body: some View {
        VStack(spacing: 6) {
            Circle()
                .strokeBorder(
                    LinearGradient(
                        colors: [VignetteColors.moonstone, VignetteColors.saffron],
                        startPoint: .topLeading,
                        endPoint: .bottomTrailing
                    ),
                    lineWidth: 2.5
                )
                .background(
                    Circle()
                        .fill(Color.gray.opacity(0.3))
                )
                .frame(width: 72, height: 72)
            
            Text(story.userName)
                .vignetteCaptionSmall()
                .foregroundColor(VignetteColors.textPrimary)
                .lineLimit(1)
                .frame(width: 72)
        }
    }
}

// MARK: - Post Card (Instagram-style)
struct VignettePostCard: View {
    let post: VignettePost
    @State private var isLiked = false
    @State private var isSaved = false
    
    var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            // Post header
            HStack(spacing: 12) {
                Circle()
                    .strokeBorder(
                        LinearGradient(
                            colors: [VignetteColors.moonstone, VignetteColors.saffron],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        ),
                        lineWidth: 2
                    )
                    .background(
                        Circle()
                            .fill(Color.gray.opacity(0.3))
                    )
                    .frame(width: 32, height: 32)
                
                VStack(alignment: .leading, spacing: 2) {
                    Text(post.userName)
                        .vignetteBodyMedium()
                        .foregroundColor(VignetteColors.textPrimary)
                    
                    if let location = post.location {
                        Text(location)
                            .vignetteCaptionSmall()
                            .foregroundColor(VignetteColors.textSecondary)
                    }
                }
                
                Spacer()
                
                Button(action: {}) {
                    Image(systemName: "ellipsis")
                        .font(.system(size: 16))
                        .foregroundColor(VignetteColors.textPrimary)
                }
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)
            
            // Post image
            Rectangle()
                .fill(Color.gray.opacity(0.2))
                .aspectRatio(1, contentMode: .fill)
                .overlay(
                    Image(systemName: "photo")
                        .font(.system(size: 48))
                        .foregroundColor(.gray.opacity(0.5))
                )
            
            // Action buttons
            HStack(spacing: 16) {
                Button(action: {
                    withAnimation(.spring(response: 0.3)) {
                        isLiked.toggle()
                    }
                }) {
                    Image(systemName: isLiked ? "heart.fill" : "heart")
                        .font(.system(size: 26, weight: .regular))
                        .foregroundColor(isLiked ? .red : VignetteColors.textPrimary)
                }
                
                Button(action: {}) {
                    Image(systemName: "bubble.right")
                        .font(.system(size: 25, weight: .regular))
                        .foregroundColor(VignetteColors.textPrimary)
                }
                
                Button(action: {}) {
                    Image(systemName: "paperplane")
                        .font(.system(size: 24, weight: .regular))
                        .foregroundColor(VignetteColors.textPrimary)
                }
                
                Spacer()
                
                Button(action: {
                    withAnimation(.spring(response: 0.3)) {
                        isSaved.toggle()
                    }
                }) {
                    Image(systemName: isSaved ? "bookmark.fill" : "bookmark")
                        .font(.system(size: 24, weight: .regular))
                        .foregroundColor(VignetteColors.textPrimary)
                }
            }
            .padding(.horizontal, 16)
            .padding(.top, 12)
            
            // Likes count
            if post.likesCount > 0 {
                Text("\(post.likesCount) likes")
                    .vignetteBodySmall()
                    .foregroundColor(VignetteColors.textPrimary)
                    .padding(.horizontal, 16)
                    .padding(.top, 8)
            }
            
            // Caption
            if let caption = post.caption {
                HStack(alignment: .top, spacing: 4) {
                    Text(post.userName)
                        .vignetteBodyMedium()
                        .foregroundColor(VignetteColors.textPrimary)
                    
                    Text(caption)
                        .vignetteBodyRegular()
                        .foregroundColor(VignetteColors.textPrimary)
                        .lineLimit(2)
                }
                .padding(.horizontal, 16)
                .padding(.top, 4)
            }
            
            // View comments
            if post.commentsCount > 0 {
                Text("View all \(post.commentsCount) comments")
                    .vignetteCaptionMedium()
                    .foregroundColor(VignetteColors.textSecondary)
                    .padding(.horizontal, 16)
                    .padding(.top, 4)
            }
            
            // Timestamp
            Text(post.timestamp)
                .vignetteCaptionSmall()
                .foregroundColor(VignetteColors.textSecondary)
                .padding(.horizontal, 16)
                .padding(.top, 4)
                .padding(.bottom, 12)
        }
    }
}

// MARK: - Placeholder Views
struct VignetteTakesView: View {
    var body: some View {
        ZStack {
            Color.black
                .ignoresSafeArea()
            
            VStack {
                Text("Takes")
                    .vignetteHeadlineLarge()
                    .foregroundColor(.white)
                Text("Coming Soon")
                    .vignetteBodyRegular()
                    .foregroundColor(.white.opacity(0.7))
            }
        }
    }
}

struct VignetteDirectView: View {
    var body: some View {
        VStack {
            Text("Direct")
                .vignetteHeadlineLarge()
            Text("Coming Soon")
                .vignetteBodyRegular()
                .foregroundColor(VignetteColors.textSecondary)
        }
    }
}

struct VignetteActivityView: View {
    var body: some View {
        VStack {
            Text("Activity")
                .vignetteHeadlineLarge()
            Text("Coming Soon")
                .vignetteBodyRegular()
                .foregroundColor(VignetteColors.textSecondary)
        }
    }
}

struct VignetteProfileView: View {
    var body: some View {
        ScrollView {
            VStack(spacing: 20) {
                // Profile header
                VStack(spacing: 16) {
                    Circle()
                        .fill(Color.gray.opacity(0.3))
                        .frame(width: 86, height: 86)
                    
                    Text("@username")
                        .vignetteBodyLarge()
                        .foregroundColor(VignetteColors.textPrimary)
                }
                .padding(.top, 20)
                
                // Stats
                HStack(spacing: 40) {
                    StatView(count: "245", label: "Posts")
                    StatView(count: "1.2K", label: "Followers")
                    StatView(count: "892", label: "Following")
                }
                
                // Edit Profile button
                Button(action: {}) {
                    Text("Edit Profile")
                        .vignetteButtonMedium()
                        .foregroundColor(EntativaColors.primaryBlue)
                        .frame(maxWidth: .infinity)
                        .frame(height: 44)
                        .background(
                            RoundedRectangle(cornerRadius: 8)
                                .stroke(VignetteColors.separator, lineWidth: 1)
                        )
                }
                .padding(.horizontal, 16)
                
                Divider()
                    .padding(.top, 8)
                
                Text("Coming Soon")
                    .vignetteBodyRegular()
                    .foregroundColor(VignetteColors.textSecondary)
            }
        }
    }
}

struct StatView: View {
    let count: String
    let label: String
    
    var body: some View {
        VStack(spacing: 4) {
            Text(count)
                .vignetteBodyLarge()
                .foregroundColor(VignetteColors.textPrimary)
            Text(label)
                .vignetteCaptionMedium()
                .foregroundColor(VignetteColors.textSecondary)
        }
    }
}

struct VignetteCameraView: View {
    @Environment(\.dismiss) var dismiss
    
    var body: some View {
        ZStack {
            Color.black
                .ignoresSafeArea()
            
            VStack {
                HStack {
                    Button(action: { dismiss() }) {
                        Image(systemName: "xmark")
                            .font(.system(size: 24))
                            .foregroundColor(.white)
                    }
                    .padding()
                    
                    Spacer()
                }
                
                Spacer()
                
                Text("Camera")
                    .vignetteHeadlineLarge()
                    .foregroundColor(.white)
                Text("Coming Soon")
                    .vignetteBodyRegular()
                    .foregroundColor(.white.opacity(0.7))
                
                Spacer()
            }
        }
    }
}

struct VignetteSearchView: View {
    var body: some View {
        Text("Search")
            .vignetteHeadlineLarge()
    }
}

// MARK: - Supporting Types
enum VignetteTab: CaseIterable {
    case home, takes, messages, activity, profile
    
    var icon: String {
        switch self {
        case .home: return "house"
        case .takes: return "play.rectangle"
        case .messages: return "paperplane"
        case .activity: return "heart"
        case .profile: return "person.circle"
        }
    }
}

// MARK: - View Model
class VignetteHomeViewModel: ObservableObject {
    @Published var posts: [VignettePost] = VignettePost.mockPosts
    @Published var isLoading = false
    
    @MainActor
    func refreshFeed() async {
        isLoading = true
        try? await Task.sleep(nanoseconds: 1_000_000_000)
        // Load new posts
        isLoading = false
    }
}

// MARK: - Models
struct VignettePost: Identifiable {
    let id: String
    let userId: String
    let userName: String
    let userAvatar: String?
    let caption: String?
    let mediaUrl: String
    let location: String?
    let likesCount: Int
    let commentsCount: Int
    let timestamp: String
    let createdAt: Date
    
    static let mockPosts: [VignettePost] = [
        VignettePost(
            id: "1",
            userId: "user1",
            userName: "alexjohnson",
            userAvatar: nil,
            caption: "Golden hour magic ✨ Perfect end to a perfect day",
            mediaUrl: "",
            location: "San Francisco, CA",
            likesCount: 1247,
            commentsCount: 89,
            timestamp: "2 HOURS AGO",
            createdAt: Date()
        ),
        VignettePost(
            id: "2",
            userId: "user2",
            userName: "sarahcreative",
            userAvatar: nil,
            caption: "New artwork incoming! Can't wait to share the full collection 🎨",
            mediaUrl: "",
            location: "Brooklyn, NY",
            likesCount: 2134,
            commentsCount: 156,
            timestamp: "5 HOURS AGO",
            createdAt: Date()
        ),
        VignettePost(
            id: "3",
            userId: "user3",
            userName: "mikefitness",
            userAvatar: nil,
            caption: "Day 100 of my fitness journey! 💪",
            mediaUrl: "",
            location: nil,
            likesCount: 892,
            commentsCount: 67,
            timestamp: "1 DAY AGO",
            createdAt: Date()
        )
    ]
}

struct VignetteStory: Identifiable {
    let id: String
    let userId: String
    let userName: String
    let userAvatar: String?
    let mediaUrl: String
    let createdAt: Date
    let expiresAt: Date
    let isViewed: Bool
    
    static let mockStories: [VignetteStory] = [
        VignetteStory(
            id: "story1",
            userId: "user1",
            userName: "emma.lee",
            userAvatar: nil,
            mediaUrl: "",
            createdAt: Date(),
            expiresAt: Date().addingTimeInterval(86400),
            isViewed: false
        ),
        VignetteStory(
            id: "story2",
            userId: "user2",
            userName: "travel_tom",
            userAvatar: nil,
            mediaUrl: "",
            createdAt: Date(),
            expiresAt: Date().addingTimeInterval(86400),
            isViewed: false
        ),
        VignetteStory(
            id: "story3",
            userId: "user3",
            userName: "foodie.jane",
            userAvatar: nil,
            mediaUrl: "",
            createdAt: Date(),
            expiresAt: Date().addingTimeInterval(86400),
            isViewed: true
        )
    ]
}

// MARK: - Preview
struct VignetteHomeView_Previews: PreviewProvider {
    static var previews: some View {
        VignetteHomeView()
    }
}
