import SwiftUI

struct EntativaHomeView: View {
    @StateObject private var viewModel = HomeViewModel()
    @State private var showingCreatePost = false
    @State private var showingSearch = false
    @State private var selectedTab: HomeTab = .home
    
    var body: some View {
        ZStack(alignment: .bottom) {
            // Main content
            VStack(spacing: 0) {
                // Top bar
                EntativaTopBar(
                    onPlusAction: {
                        showingCreatePost = true
                    },
                    onSearchAction: {
                        showingSearch = true
                    }
                )
                
                // Content based on selected tab
                TabView(selection: $selectedTab) {
                    // Home Feed
                    EntativaFeedView(viewModel: viewModel)
                        .tag(HomeTab.home)
                    
                    // Takes
                    EntativaTakesView()
                        .tag(HomeTab.takes)
                    
                    // Messages
                    EntativaMessagesView()
                        .tag(HomeTab.messages)
                    
                    // Activity (Notifications)
                    EntativaActivityView()
                        .tag(HomeTab.activity)
                    
                    // Menu
                    EntativaMenuView()
                        .tag(HomeTab.menu)
                }
                .tabViewStyle(.page(indexDisplayMode: .never))
            }
            
            // Floating liquid glass bottom navigation
            EntativaBottomNavBar(selectedTab: $selectedTab)
                .padding(.bottom, 8)
                .padding(.horizontal, 16)
        }
        .ignoresSafeArea(.keyboard)
        .sheet(isPresented: $showingCreatePost) {
            EntativaCreatePostView()
        }
        .sheet(isPresented: $showingSearch) {
            EntativaSearchView()
        }
    }
}

// MARK: - Top Bar Component
struct EntativaTopBar: View {
    let onPlusAction: () -> Void
    let onSearchAction: () -> Void
    
    var body: some View {
        HStack(spacing: 16) {
            // Plus button (left)
            Button(action: onPlusAction) {
                Image(systemName: "plus.circle.fill")
                    .font(.system(size: 28, weight: .medium))
                    .foregroundColor(EntativaColors.primaryBlue)
            }
            
            Spacer()
            
            // Entativa logo (center)
            HStack(spacing: 2) {
                Text("entativa")
                    .font(.custom("SFProRounded-Bold", size: 28))
                    .italic()
                    .foregroundStyle(
                        LinearGradient(
                            colors: [
                                EntativaColors.primaryBlue,
                                EntativaColors.primaryPurple,
                                EntativaColors.primaryPink
                            ],
                            startPoint: .leading,
                            endPoint: .trailing
                        )
                    )
            }
            
            Spacer()
            
            // Search button (right)
            Button(action: onSearchAction) {
                Image(systemName: "magnifyingglass")
                    .font(.system(size: 22, weight: .semibold))
                    .foregroundColor(EntativaColors.textPrimary)
            }
        }
        .padding(.horizontal, 20)
        .padding(.vertical, 12)
        .background(
            Color(UIColor.systemBackground)
                .shadow(color: .black.opacity(0.05), radius: 1, y: 1)
        )
    }
}

// MARK: - Bottom Navigation Bar (Liquid Glass Effect)
struct EntativaBottomNavBar: View {
    @Binding var selectedTab: HomeTab
    
    var body: some View {
        HStack(spacing: 0) {
            ForEach(HomeTab.allCases, id: \.self) { tab in
                Button(action: {
                    withAnimation(.spring(response: 0.3, dampingFraction: 0.7)) {
                        selectedTab = tab
                    }
                }) {
                    VStack(spacing: 4) {
                        Image(systemName: tab.icon)
                            .font(.system(size: 24, weight: selectedTab == tab ? .semibold : .regular))
                            .foregroundColor(selectedTab == tab ? EntativaColors.primaryBlue : EntativaColors.textSecondary)
                        
                        if selectedTab == tab {
                            Circle()
                                .fill(EntativaColors.primaryBlue)
                                .frame(width: 4, height: 4)
                                .transition(.scale.combined(with: .opacity))
                        }
                    }
                    .frame(maxWidth: .infinity)
                    .padding(.vertical, 12)
                }
            }
        }
        .padding(.horizontal, 8)
        .background(
            // Liquid glass effect
            ZStack {
                // Blur background
                Color(UIColor.systemBackground)
                    .opacity(0.8)
                
                // Glass effect
                Rectangle()
                    .fill(.ultraThinMaterial)
            }
            .cornerRadius(24)
            .shadow(color: .black.opacity(0.1), radius: 20, y: 10)
        )
    }
}

// MARK: - Feed View (Carousel Posts)
struct EntativaFeedView: View {
    @ObservedObject var viewModel: HomeViewModel
    @State private var combinedPosts: [String: [Post]] = [:]
    
    var body: some View {
        ScrollView {
            LazyVStack(spacing: 0) {
                // Card Stories (Facebook-style)
                EntativaStoriesRow()
                    .padding(.bottom, 12)
                
                // Posts in carousel format (like Threads)
                ForEach(viewModel.posts) { post in
                    EntativaCarouselPostCard(post: post)
                        .padding(.horizontal, 16)
                        .padding(.vertical, 8)
                }
            }
        }
        .refreshable {
            await viewModel.refreshFeed()
        }
    }
}

// MARK: - Stories Row (Card Style)
struct EntativaStoriesRow: View {
    @State private var stories: [Story] = Story.mockStories
    
    var body: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: 12) {
                // Create Story Card
                VStack(spacing: 8) {
                    ZStack(alignment: .bottomTrailing) {
                        RoundedRectangle(cornerRadius: 16)
                            .fill(
                                LinearGradient(
                                    colors: [
                                        EntativaColors.primaryBlue.opacity(0.3),
                                        EntativaColors.primaryPurple.opacity(0.3)
                                    ],
                                    startPoint: .topLeading,
                                    endPoint: .bottomTrailing
                                )
                            )
                            .frame(width: 110, height: 160)
                        
                        Image(systemName: "plus.circle.fill")
                            .font(.system(size: 32))
                            .foregroundStyle(EntativaColors.primaryBlue)
                            .background(
                                Circle()
                                    .fill(Color.white)
                                    .frame(width: 36, height: 36)
                            )
                            .offset(x: -8, y: -8)
                    }
                    
                    Text("Create")
                        .entativaLabelSmall()
                        .foregroundColor(EntativaColors.textPrimary)
                }
                
                // Story Cards
                ForEach(stories) { story in
                    StoryCard(story: story)
                }
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 8)
        }
    }
}

struct StoryCard: View {
    let story: Story
    
    var body: some View {
        VStack(spacing: 8) {
            ZStack(alignment: .topLeading) {
                // Story preview image
                RoundedRectangle(cornerRadius: 16)
                    .fill(
                        LinearGradient(
                            colors: [Color.blue.opacity(0.5), Color.purple.opacity(0.5)],
                            startPoint: .top,
                            endPoint: .bottom
                        )
                    )
                    .frame(width: 110, height: 160)
                
                // User avatar
                Circle()
                    .fill(
                        LinearGradient(
                            colors: [EntativaColors.primaryBlue, EntativaColors.primaryPurple],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                    )
                    .frame(width: 40, height: 40)
                    .overlay(
                        Circle()
                            .stroke(Color.white, lineWidth: 3)
                    )
                    .padding(8)
            }
            
            Text(story.userName)
                .entativaLabelSmall()
                .foregroundColor(EntativaColors.textPrimary)
                .lineLimit(1)
        }
    }
}

// MARK: - Carousel Post Card (Threads-style)
struct EntativaCarouselPostCard: View {
    let post: Post
    @State private var currentMediaIndex = 0
    @GestureState private var dragOffset: CGFloat = 0
    
    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            // Post header
            HStack(spacing: 12) {
                Circle()
                    .fill(
                        LinearGradient(
                            colors: [EntativaColors.primaryBlue, EntativaColors.primaryPurple],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                    )
                    .frame(width: 44, height: 44)
                
                VStack(alignment: .leading, spacing: 2) {
                    Text(post.userName)
                        .entativaBodyMedium()
                        .foregroundColor(EntativaColors.textPrimary)
                    
                    Text(post.timestamp)
                        .entativaCaptionSmall()
                        .foregroundColor(EntativaColors.textSecondary)
                }
                
                Spacer()
                
                Button(action: {}) {
                    Image(systemName: "ellipsis")
                        .font(.system(size: 20))
                        .foregroundColor(EntativaColors.textSecondary)
                }
            }
            
            // Post text
            if let text = post.text {
                Text(text)
                    .entativaBodyRegular()
                    .foregroundColor(EntativaColors.textPrimary)
            }
            
            // Media carousel
            if !post.media.isEmpty {
                TabView(selection: $currentMediaIndex) {
                    ForEach(Array(post.media.enumerated()), id: \.offset) { index, media in
                        RoundedRectangle(cornerRadius: 16)
                            .fill(Color.gray.opacity(0.2))
                            .frame(height: 300)
                            .overlay(
                                Image(systemName: "photo")
                                    .font(.system(size: 48))
                                    .foregroundColor(.gray.opacity(0.5))
                            )
                            .tag(index)
                    }
                }
                .frame(height: 300)
                .tabViewStyle(.page(indexDisplayMode: .automatic))
                .indexViewStyle(.page(backgroundDisplayMode: .always))
            }
            
            // Action buttons
            HStack(spacing: 20) {
                ActionButton(icon: "heart", count: post.likesCount, color: EntativaColors.textSecondary)
                ActionButton(icon: "bubble.right", count: post.commentsCount, color: EntativaColors.textSecondary)
                ActionButton(icon: "paperplane", count: post.sharesCount, color: EntativaColors.textSecondary)
                
                Spacer()
                
                Button(action: {}) {
                    Image(systemName: "bookmark")
                        .font(.system(size: 20))
                        .foregroundColor(EntativaColors.textSecondary)
                }
            }
        }
        .padding(16)
        .background(
            RoundedRectangle(cornerRadius: 20)
                .fill(Color(UIColor.secondarySystemBackground))
        )
    }
}

struct ActionButton: View {
    let icon: String
    let count: Int
    let color: Color
    
    var body: some View {
        Button(action: {}) {
            HStack(spacing: 6) {
                Image(systemName: icon)
                    .font(.system(size: 20))
                
                if count > 0 {
                    Text("\(count)")
                        .entativaCaptionMedium()
                }
            }
            .foregroundColor(color)
        }
    }
}

// MARK: - Placeholder Views
struct EntativaTakesView: View {
    var body: some View {
        VStack {
            Text("Takes")
                .entativaHeadlineLarge()
            Text("Coming Soon")
                .entativaBodyRegular()
                .foregroundColor(EntativaColors.textSecondary)
        }
    }
}

struct EntativaMessagesView: View {
    var body: some View {
        VStack {
            Text("Messages")
                .entativaHeadlineLarge()
            Text("Coming Soon")
                .entativaBodyRegular()
                .foregroundColor(EntativaColors.textSecondary)
        }
    }
}

struct EntativaActivityView: View {
    var body: some View {
        VStack {
            Text("Activity")
                .entativaHeadlineLarge()
            Text("Coming Soon")
                .entativaBodyRegular()
                .foregroundColor(EntativaColors.textSecondary)
        }
    }
}

struct EntativaMenuView: View {
    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 20) {
                Text("Menu")
                    .entativaHeadlineLarge()
                    .padding(.horizontal, 20)
                
                MenuSection(title: "Your Shortcuts", items: [
                    MenuItem(icon: "person.2", title: "Friends"),
                    MenuItem(icon: "rectangle.3.group", title: "Groups"),
                    MenuItem(icon: "play.rectangle", title: "Watch"),
                    MenuItem(icon: "bag", title: "Marketplace"),
                ])
                
                MenuSection(title: "Settings & Privacy", items: [
                    MenuItem(icon: "gearshape", title: "Settings"),
                    MenuItem(icon: "lock.shield", title: "Privacy Center"),
                    MenuItem(icon: "info.circle", title: "About"),
                ])
            }
            .padding(.vertical, 20)
        }
    }
}

struct MenuSection: View {
    let title: String
    let items: [MenuItem]
    
    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text(title)
                .entativaLabelLarge()
                .foregroundColor(EntativaColors.textSecondary)
                .padding(.horizontal, 20)
            
            ForEach(items) { item in
                HStack(spacing: 16) {
                    Image(systemName: item.icon)
                        .font(.system(size: 24))
                        .foregroundColor(EntativaColors.primaryBlue)
                        .frame(width: 32)
                    
                    Text(item.title)
                        .entativaBodyMedium()
                        .foregroundColor(EntativaColors.textPrimary)
                    
                    Spacer()
                    
                    Image(systemName: "chevron.right")
                        .font(.system(size: 14))
                        .foregroundColor(EntativaColors.textSecondary)
                }
                .padding(.horizontal, 20)
                .padding(.vertical, 12)
                .background(Color(UIColor.secondarySystemBackground))
            }
        }
    }
}

struct MenuItem: Identifiable {
    let id = UUID()
    let icon: String
    let title: String
}

struct EntativaCreatePostView: View {
    var body: some View {
        Text("Create Post")
            .entativaHeadlineLarge()
    }
}

struct EntativaSearchView: View {
    var body: some View {
        Text("Search")
            .entativaHeadlineLarge()
    }
}

// MARK: - Supporting Types
enum HomeTab: CaseIterable {
    case home, takes, messages, activity, menu
    
    var icon: String {
        switch self {
        case .home: return "house.fill"
        case .takes: return "play.rectangle.fill"
        case .messages: return "message.fill"
        case .activity: return "bell.fill"
        case .menu: return "line.3.horizontal"
        }
    }
}

// MARK: - View Model
class HomeViewModel: ObservableObject {
    @Published var posts: [Post] = Post.mockPosts
    @Published var isLoading = false
    
    @MainActor
    func refreshFeed() async {
        isLoading = true
        try? await Task.sleep(nanoseconds: 1_000_000_000)
        // Load new posts
        isLoading = false
    }
}

// MARK: - Preview
struct EntativaHomeView_Previews: PreviewProvider {
    static var previews: some View {
        EntativaHomeView()
    }
}
