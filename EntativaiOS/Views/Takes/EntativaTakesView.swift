import SwiftUI
import AVFoundation

struct EntativaTakesView: View {
    @StateObject private var viewModel = TakesViewModel()
    @State private var currentIndex: Int = 0
    @GestureState private var dragOffset: CGFloat = 0
    
    var body: some View {
        GeometryReader { geometry in
            ZStack {
                // Background
                Color.black.ignoresSafeArea()
                
                // Video Feed
                TabView(selection: $currentIndex) {
                    ForEach(Array(viewModel.takes.enumerated()), id: \.offset) { index, take in
                        TakeVideoPlayer(
                            take: take,
                            isCurrentlyPlaying: index == currentIndex,
                            geometry: geometry,
                            viewModel: viewModel
                        )
                        .tag(index)
                        .onAppear {
                            // Load more when near end
                            if index == viewModel.takes.count - 2 {
                                Task {
                                    await viewModel.loadMore()
                                }
                            }
                        }
                    }
                }
                .tabViewStyle(.page(indexDisplayMode: .never))
                .ignoresSafeArea()
                
                // Top Bar (minimal)
                VStack {
                    HStack {
                        Text("Takes")
                            .font(.system(size: 20, weight: .semibold))
                            .foregroundColor(.white)
                        
                        Spacer()
                        
                        Button(action: {}) {
                            Image(systemName: "camera")
                                .font(.system(size: 22))
                                .foregroundColor(.white)
                        }
                    }
                    .padding(.horizontal, 20)
                    .padding(.top, 50)
                    
                    Spacer()
                }
                .ignoresSafeArea()
            }
        }
    }
}

// MARK: - Video Player Card
struct TakeVideoPlayer: View {
    let take: TakeModel
    let isCurrentlyPlaying: Bool
    let geometry: GeometryProxy
    let viewModel: TakesViewModel
    
    @State private var showComments = false
    @State private var showShare = false
    @State private var isFollowing = false
    
    var body: some View {
        ZStack(alignment: .bottomTrailing) {
            // Video Player (Real AVPlayer)
            if let videoURL = URL(string: take.videoUrl) {
                VideoPlayerView(
                    videoURL: videoURL,
                    isPlaying: isCurrentlyPlaying
                )
            } else {
                // Fallback placeholder
                Rectangle()
                    .fill(
                        LinearGradient(
                            colors: [
                                Color.blue.opacity(0.4),
                                Color.purple.opacity(0.6),
                                Color.pink.opacity(0.4)
                            ],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                    )
                    .overlay(
                        VStack {
                            Image(systemName: "play.circle.fill")
                                .font(.system(size: 80))
                                .foregroundColor(.white.opacity(0.8))
                            
                            Text("Video Unavailable")
                                .font(.system(size: 16))
                                .foregroundColor(.white.opacity(0.6))
                        }
                    )
            }
            
            // Right Side Actions
            VStack(spacing: 24) {
                // Profile Avatar
                Button(action: {}) {
                    ZStack(alignment: .bottom) {
                        Circle()
                            .fill(
                                LinearGradient(
                                    colors: [
                                        EntativaColors.primaryBlue,
                                        EntativaColors.primaryPurple
                                    ],
                                    startPoint: .topLeading,
                                    endPoint: .bottomTrailing
                                )
                            )
                            .frame(width: 48, height: 48)
                        
                        if !isFollowing {
                            Circle()
                                .fill(EntativaColors.primaryBlue)
                                .frame(width: 20, height: 20)
                                .overlay(
                                    Image(systemName: "plus")
                                        .font(.system(size: 12, weight: .bold))
                                        .foregroundColor(.white)
                                )
                                .offset(y: 8)
                        }
                    }
                }
                
                // Like
                VStack(spacing: 4) {
                    Button(action: {
                        Task {
                            if take.isLiked {
                                await viewModel.unlikeTake(takeID: take.id)
                            } else {
                                await viewModel.likeTake(takeID: take.id)
                            }
                        }
                    }) {
                        Image(systemName: take.isLiked ? "heart.fill" : "heart")
                            .font(.system(size: 32, weight: .medium))
                            .foregroundColor(take.isLiked ? .red : .white)
                            .scaleEffect(take.isLiked ? 1.1 : 1.0)
                    }
                    
                    Text(formatCount(take.likesCount))
                        .font(.system(size: 12, weight: .semibold))
                        .foregroundColor(.white)
                }
                
                // Comments
                VStack(spacing: 4) {
                    Button(action: { showComments = true }) {
                        Image(systemName: "bubble.right.fill")
                            .font(.system(size: 30, weight: .medium))
                            .foregroundColor(.white)
                    }
                    
                    Text(formatCount(take.commentsCount))
                        .font(.system(size: 12, weight: .semibold))
                        .foregroundColor(.white)
                }
                
                // Share
                VStack(spacing: 4) {
                    Button(action: { showShare = true }) {
                        Image(systemName: "paperplane.fill")
                            .font(.system(size: 28, weight: .medium))
                            .foregroundColor(.white)
                    }
                    
                    Text(formatCount(take.sharesCount))
                        .font(.system(size: 12, weight: .semibold))
                        .foregroundColor(.white)
                }
                
                // More Options
                Button(action: {}) {
                    Image(systemName: "ellipsis")
                        .font(.system(size: 24, weight: .medium))
                        .foregroundColor(.white)
                }
            }
            .padding(.trailing, 16)
            .padding(.bottom, 100)
            
            // Bottom Info
            VStack(alignment: .leading, spacing: 12) {
                // Username
                HStack(spacing: 8) {
                    Text("@\(take.username)")
                        .font(.system(size: 16, weight: .semibold))
                        .foregroundColor(.white)
                    
                    if !isFollowing {
                        Button(action: {
                            withAnimation {
                                isFollowing.toggle()
                            }
                        }) {
                            Text("Follow")
                                .font(.system(size: 14, weight: .semibold))
                                .foregroundColor(.white)
                                .padding(.horizontal, 16)
                                .padding(.vertical, 6)
                                .background(
                                    LinearGradient(
                                        colors: [
                                            EntativaColors.primaryBlue,
                                            EntativaColors.primaryPurple
                                        ],
                                        startPoint: .leading,
                                        endPoint: .trailing
                                    )
                                )
                                .cornerRadius(4)
                        }
                    }
                }
                
                // Caption
                Text(take.caption)
                    .font(.system(size: 14))
                    .foregroundColor(.white)
                    .lineLimit(2)
                
                // Audio Info
                HStack(spacing: 8) {
                    Image(systemName: "music.note")
                        .font(.system(size: 12))
                    
                    Text(take.audioName)
                        .font(.system(size: 13))
                    
                    Spacer()
                }
                .foregroundColor(.white.opacity(0.9))
            }
            .frame(maxWidth: .infinity, alignment: .leading)
            .padding(.horizontal, 16)
            .padding(.bottom, 100)
            .padding(.trailing, 80) // Space for right sidebar
        }
        .sheet(isPresented: $showComments) {
            TakeCommentsSheet(take: take)
        }
        .sheet(isPresented: $showShare) {
            TakeShareSheet(take: take)
        }
    }
}

// MARK: - Comments Sheet
struct TakeCommentsSheet: View {
    let take: Take
    @Environment(\.dismiss) var dismiss
    @State private var commentText = ""
    
    var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Comments List
                ScrollView {
                    LazyVStack(spacing: 16) {
                        ForEach(0..<10) { index in
                            CommentRow(
                                username: "user\(index)",
                                comment: "This is an amazing take! Love it! 🔥",
                                timestamp: "\(index + 1)h ago",
                                likesCount: Int.random(in: 10...500)
                            )
                        }
                    }
                    .padding()
                }
                
                Divider()
                
                // Comment Input
                HStack(spacing: 12) {
                    Circle()
                        .fill(Color.gray.opacity(0.3))
                        .frame(width: 32, height: 32)
                    
                    TextField("Add comment...", text: $commentText)
                        .padding(.horizontal, 12)
                        .padding(.vertical, 8)
                        .background(Color.gray.opacity(0.1))
                        .cornerRadius(20)
                    
                    if !commentText.isEmpty {
                        Button(action: {
                            // Post comment
                            commentText = ""
                        }) {
                            Text("Post")
                                .font(.system(size: 15, weight: .semibold))
                                .foregroundColor(EntativaColors.primaryBlue)
                        }
                    }
                }
                .padding()
            }
            .navigationTitle("\(take.commentsCount) comments")
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

struct CommentRow: View {
    let username: String
    let comment: String
    let timestamp: String
    let likesCount: Int
    @State private var isLiked = false
    
    var body: some View {
        HStack(alignment: .top, spacing: 12) {
            Circle()
                .fill(Color.gray.opacity(0.3))
                .frame(width: 32, height: 32)
            
            VStack(alignment: .leading, spacing: 4) {
                HStack {
                    Text(username)
                        .font(.system(size: 13, weight: .semibold))
                        .foregroundColor(EntativaColors.textPrimary)
                    
                    Text(timestamp)
                        .font(.system(size: 12))
                        .foregroundColor(EntativaColors.textSecondary)
                }
                
                Text(comment)
                    .font(.system(size: 14))
                    .foregroundColor(EntativaColors.textPrimary)
                
                HStack(spacing: 16) {
                    Text("\(likesCount) likes")
                        .font(.system(size: 12, weight: .medium))
                        .foregroundColor(EntativaColors.textSecondary)
                    
                    Button(action: {}) {
                        Text("Reply")
                            .font(.system(size: 12, weight: .medium))
                            .foregroundColor(EntativaColors.textSecondary)
                    }
                }
                .padding(.top, 2)
            }
            
            Spacer()
            
            Button(action: {
                withAnimation {
                    isLiked.toggle()
                }
            }) {
                Image(systemName: isLiked ? "heart.fill" : "heart")
                    .font(.system(size: 14))
                    .foregroundColor(isLiked ? .red : EntativaColors.textSecondary)
            }
        }
    }
}

// MARK: - Share Sheet
struct TakeShareSheet: View {
    let take: Take
    @Environment(\.dismiss) var dismiss
    
    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    // Share Options
                    LazyVGrid(columns: [
                        GridItem(.flexible()),
                        GridItem(.flexible()),
                        GridItem(.flexible()),
                        GridItem(.flexible())
                    ], spacing: 20) {
                        ShareOption(icon: "person.2.fill", title: "Friends")
                        ShareOption(icon: "link", title: "Copy Link")
                        ShareOption(icon: "square.and.arrow.up", title: "Share to...")
                        ShareOption(icon: "bookmark", title: "Save")
                        ShareOption(icon: "flag", title: "Report")
                        ShareOption(icon: "eye.slash", title: "Not Interested")
                    }
                    .padding()
                }
            }
            .navigationTitle("Share")
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

struct ShareOption: View {
    let icon: String
    let title: String
    
    var body: some View {
        VStack(spacing: 8) {
            Circle()
                .fill(Color.gray.opacity(0.1))
                .frame(width: 60, height: 60)
                .overlay(
                    Image(systemName: icon)
                        .font(.system(size: 24))
                        .foregroundColor(EntativaColors.textPrimary)
                )
            
            Text(title)
                .font(.system(size: 12))
                .foregroundColor(EntativaColors.textPrimary)
                .multilineTextAlignment(.center)
        }
    }
}

// MARK: - View Model
class TakesViewModel: ObservableObject {
    @Published var takes: [TakeModel] = []
    @Published var isLoading = false
    @Published var errorMessage: String?
    
    private var currentPage = 1
    private var hasMore = true
    private let apiClient = TakesAPIClient.shared
    
    init() {
        Task {
            await loadFeed()
        }
    }
    
    @MainActor
    func loadFeed() async {
        guard !isLoading else { return }
        
        isLoading = true
        errorMessage = nil
        
        do {
            let response = try await apiClient.getFeed(page: currentPage, limit: 10)
            self.takes = response.takes
            self.hasMore = response.hasMore
            
            // Preload next videos
            if takes.count > 1 {
                for i in 0..<min(3, takes.count) {
                    if let url = URL(string: takes[i].videoUrl) {
                        VideoCache.shared.preload(url: url)
                    }
                }
            }
        } catch {
            self.errorMessage = error.localizedDescription
            // Fallback to mock data if API fails
            self.takes = convertMockTakes()
        }
        
        isLoading = false
    }
    
    @MainActor
    func loadMore() async {
        guard !isLoading, hasMore else { return }
        
        isLoading = true
        currentPage += 1
        
        do {
            let response = try await apiClient.getFeed(page: currentPage, limit: 10)
            self.takes.append(contentsOf: response.takes)
            self.hasMore = response.hasMore
        } catch {
            self.errorMessage = error.localizedDescription
        }
        
        isLoading = false
    }
    
    @MainActor
    func likeTake(takeID: String) async {
        do {
            let updatedTake = try await apiClient.likeTake(takeID: takeID)
            if let index = takes.firstIndex(where: { $0.id == takeID }) {
                takes[index] = updatedTake
            }
        } catch {
            self.errorMessage = error.localizedDescription
        }
    }
    
    @MainActor
    func unlikeTake(takeID: String) async {
        do {
            let updatedTake = try await apiClient.unlikeTake(takeID: takeID)
            if let index = takes.firstIndex(where: { $0.id == takeID }) {
                takes[index] = updatedTake
            }
        } catch {
            self.errorMessage = error.localizedDescription
        }
    }
    
    // Convert mock data to TakeModel format (fallback)
    private func convertMockTakes() -> [TakeModel] {
        return Take.mockTakes.map { mockTake in
            TakeModel(
                id: mockTake.id,
                userId: mockTake.userId,
                username: mockTake.username,
                userAvatar: mockTake.userAvatar,
                videoUrl: "https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_1mb.mp4",
                thumbnailUrl: mockTake.thumbnailUrl,
                caption: mockTake.caption,
                audioName: mockTake.audioName,
                audioUrl: nil,
                duration: 30,
                likesCount: mockTake.likesCount,
                commentsCount: mockTake.commentsCount,
                sharesCount: mockTake.sharesCount,
                viewsCount: mockTake.viewsCount,
                isLiked: false,
                isSaved: false,
                hashtags: [],
                createdAt: ISO8601DateFormatter().string(from: mockTake.createdAt)
            )
        }
    }
}

// MARK: - Take Model
struct Take: Identifiable {
    let id: String
    let userId: String
    let username: String
    let userAvatar: String?
    let videoUrl: String
    let thumbnailUrl: String?
    let caption: String
    let audioName: String
    let likesCount: Int
    let commentsCount: Int
    let sharesCount: Int
    let viewsCount: Int
    let createdAt: Date
    
    static let mockTakes: [Take] = [
        Take(
            id: "1",
            userId: "user1",
            username: "alexcreator",
            userAvatar: nil,
            videoUrl: "",
            thumbnailUrl: nil,
            caption: "Check out this amazing transformation! 💪 #fitness #motivation",
            audioName: "Original Audio - alexcreator",
            likesCount: 45_200,
            commentsCount: 892,
            sharesCount: 1_234,
            viewsCount: 234_500,
            createdAt: Date()
        ),
        Take(
            id: "2",
            userId: "user2",
            username: "foodie.life",
            userAvatar: nil,
            videoUrl: "",
            thumbnailUrl: nil,
            caption: "Best pasta recipe ever! 🍝 Try it and let me know what you think!",
            audioName: "Cooking Vibes - Sound Library",
            likesCount: 78_300,
            commentsCount: 1_456,
            sharesCount: 2_890,
            viewsCount: 456_700,
            createdAt: Date()
        ),
        Take(
            id: "3",
            userId: "user3",
            username: "travel.with.me",
            userAvatar: nil,
            videoUrl: "",
            thumbnailUrl: nil,
            caption: "Hidden gems in Bali you NEED to visit! 🌴✨ #travel #bali",
            audioName: "Tropical Summer - Music Mix",
            likesCount: 123_400,
            commentsCount: 3_421,
            sharesCount: 5_678,
            viewsCount: 890_200,
            createdAt: Date()
        ),
        Take(
            id: "4",
            userId: "user4",
            username: "tech.reviews",
            userAvatar: nil,
            videoUrl: "",
            thumbnailUrl: nil,
            caption: "iPhone 16 Pro vs Samsung S24 Ultra - The TRUTH! 📱",
            audioName: "Tech Beat - Audio Track",
            likesCount: 56_700,
            commentsCount: 2_134,
            sharesCount: 3_456,
            viewsCount: 567_800,
            createdAt: Date()
        ),
        Take(
            id: "5",
            userId: "user5",
            username: "comedy.gold",
            userAvatar: nil,
            videoUrl: "",
            thumbnailUrl: nil,
            caption: "When your code finally works 😂💻 #coding #funny",
            audioName: "Funny Moments - Sound Effect",
            likesCount: 234_500,
            commentsCount: 8_976,
            sharesCount: 12_345,
            viewsCount: 1_234_500,
            createdAt: Date()
        )
    ]
}

// MARK: - Preview
struct EntativaTakesView_Previews: PreviewProvider {
    static var previews: some View {
        EntativaTakesView()
    }
}
