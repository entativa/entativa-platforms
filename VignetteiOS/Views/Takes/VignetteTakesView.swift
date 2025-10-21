import SwiftUI
import AVFoundation

struct VignetteTakesView: View {
    @StateObject private var viewModel = VignetteTakesViewModel()
    @State private var currentIndex: Int = 0
    
    var body: some View {
        GeometryReader { geometry in
            ZStack {
                // Background
                Color.black.ignoresSafeArea()
                
                // Reels Feed
                TabView(selection: $currentIndex) {
                    ForEach(Array(viewModel.takes.enumerated()), id: \.offset) { index, take in
                        VignetteTakePlayer(
                            take: take,
                            isCurrentlyPlaying: index == currentIndex,
                            geometry: geometry
                        )
                        .tag(index)
                    }
                }
                .tabViewStyle(.page(indexDisplayMode: .never))
                .ignoresSafeArea()
                
                // Top Bar (Instagram Reels style)
                VStack {
                    HStack(spacing: 20) {
                        Text("Takes")
                            .font(.system(size: 24, weight: .bold))
                            .foregroundColor(.white)
                        
                        Spacer()
                        
                        Button(action: {}) {
                            Image(systemName: "camera")
                                .font(.system(size: 24, weight: .regular))
                                .foregroundColor(.white)
                        }
                    }
                    .padding(.horizontal, 16)
                    .padding(.top, 50)
                    
                    Spacer()
                }
                .ignoresSafeArea()
            }
        }
    }
}

// MARK: - Take Player Card
struct VignetteTakePlayer: View {
    let take: VignetteTake
    let isCurrentlyPlaying: Bool
    let geometry: GeometryProxy
    
    @State private var isLiked = false
    @State private var showComments = false
    @State private var showShare = false
    @State private var isFollowing = false
    
    var body: some View {
        ZStack(alignment: .bottomTrailing) {
            // Video Player (placeholder)
            Rectangle()
                .fill(
                    LinearGradient(
                        colors: [
                            VignetteColors.moonstone.opacity(0.5),
                            VignetteColors.gunmetal.opacity(0.7),
                            VignetteColors.saffron.opacity(0.3)
                        ],
                        startPoint: .topLeading,
                        endPoint: .bottomTrailing
                    )
                )
                .overlay(
                    VStack {
                        Image(systemName: "play.circle")
                            .font(.system(size: 80, weight: .ultraLight))
                            .foregroundColor(.white.opacity(0.9))
                        
                        Text("Video Player")
                            .font(.system(size: 16, weight: .light))
                            .foregroundColor(.white.opacity(0.7))
                    }
                )
            
            // Right Side Actions (Instagram style)
            VStack(spacing: 20) {
                // Profile Avatar
                Button(action: {}) {
                    ZStack(alignment: .bottom) {
                        Circle()
                            .strokeBorder(Color.white, lineWidth: 2)
                            .background(
                                Circle()
                                    .fill(Color.gray.opacity(0.5))
                            )
                            .frame(width: 44, height: 44)
                        
                        if !isFollowing {
                            Circle()
                                .fill(EntativaColors.primaryBlue)
                                .frame(width: 22, height: 22)
                                .overlay(
                                    Image(systemName: "plus")
                                        .font(.system(size: 12, weight: .bold))
                                        .foregroundColor(.white)
                                )
                                .offset(y: 10)
                        }
                    }
                }
                
                // Like
                VStack(spacing: 6) {
                    Button(action: {
                        withAnimation(.spring(response: 0.3)) {
                            isLiked.toggle()
                        }
                    }) {
                        Image(systemName: isLiked ? "heart.fill" : "heart")
                            .font(.system(size: 28, weight: .regular))
                            .foregroundColor(isLiked ? .red : .white)
                            .scaleEffect(isLiked ? 1.15 : 1.0)
                    }
                    
                    Text(formatCount(take.likesCount))
                        .font(.system(size: 12, weight: .medium))
                        .foregroundColor(.white)
                }
                
                // Comments
                VStack(spacing: 6) {
                    Button(action: { showComments = true }) {
                        Image(systemName: "bubble.right")
                            .font(.system(size: 28, weight: .regular))
                            .foregroundColor(.white)
                    }
                    
                    Text(formatCount(take.commentsCount))
                        .font(.system(size: 12, weight: .medium))
                        .foregroundColor(.white)
                }
                
                // Share
                VStack(spacing: 6) {
                    Button(action: { showShare = true }) {
                        Image(systemName: "paperplane")
                            .font(.system(size: 26, weight: .regular))
                            .foregroundColor(.white)
                    }
                    
                    Text(formatCount(take.sharesCount))
                        .font(.system(size: 12, weight: .medium))
                        .foregroundColor(.white)
                }
                
                // Save
                Button(action: {}) {
                    Image(systemName: "bookmark")
                        .font(.system(size: 26, weight: .regular))
                        .foregroundColor(.white)
                }
                
                // More Options
                Button(action: {}) {
                    Image(systemName: "ellipsis")
                        .font(.system(size: 22, weight: .regular))
                        .foregroundColor(.white)
                        .rotationEffect(.degrees(90))
                }
                
                // Audio Icon (spinning record)
                Button(action: {}) {
                    ZStack {
                        Circle()
                            .fill(
                                LinearGradient(
                                    colors: [VignetteColors.moonstone, VignetteColors.saffron],
                                    startPoint: .topLeading,
                                    endPoint: .bottomTrailing
                                )
                            )
                            .frame(width: 36, height: 36)
                        
                        Image(systemName: "music.note")
                            .font(.system(size: 14))
                            .foregroundColor(.white)
                    }
                }
            }
            .padding(.trailing, 12)
            .padding(.bottom, 100)
            
            // Bottom Info
            VStack(alignment: .leading, spacing: 10) {
                // Username
                HStack(spacing: 8) {
                    Text("@\(take.username)")
                        .font(.system(size: 15, weight: .semibold))
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
                                .padding(.horizontal, 12)
                                .padding(.vertical, 4)
                                .background(
                                    Capsule()
                                        .stroke(Color.white, lineWidth: 1)
                                )
                        }
                    }
                }
                
                // Caption
                Text(take.caption)
                    .font(.system(size: 13))
                    .foregroundColor(.white)
                    .lineLimit(2)
                
                // Audio Info
                HStack(spacing: 6) {
                    Image(systemName: "music.note")
                        .font(.system(size: 11))
                    
                    Text(take.audioName)
                        .font(.system(size: 12))
                        .lineLimit(1)
                }
                .foregroundColor(.white)
            }
            .frame(maxWidth: .infinity, alignment: .leading)
            .padding(.horizontal, 16)
            .padding(.bottom, 100)
            .padding(.trailing, 70)
        }
        .sheet(isPresented: $showComments) {
            VignetteTakeCommentsSheet(take: take)
        }
        .sheet(isPresented: $showShare) {
            VignetteTakeShareSheet(take: take)
        }
    }
    
    private func formatCount(_ count: Int) -> String {
        if count >= 1_000_000 {
            return String(format: "%.1fM", Double(count) / 1_000_000)
        } else if count >= 1_000 {
            return String(format: "%.1fK", Double(count) / 1_000)
        } else {
            return "\(count)"
        }
    }
}

// MARK: - Comments Sheet
struct VignetteTakeCommentsSheet: View {
    let take: VignetteTake
    @Environment(\.dismiss) var dismiss
    @State private var commentText = ""
    
    var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Header
                HStack {
                    Text("Comments")
                        .font(.system(size: 16, weight: .semibold))
                    
                    Spacer()
                    
                    Button(action: { dismiss() }) {
                        Image(systemName: "xmark")
                            .font(.system(size: 16))
                            .foregroundColor(VignetteColors.textPrimary)
                    }
                }
                .padding()
                
                Divider()
                
                // Comments List
                ScrollView {
                    LazyVStack(spacing: 12) {
                        ForEach(0..<15) { index in
                            VignetteTakeCommentRow(
                                username: "user\(index)",
                                comment: "This is awesome! 🔥",
                                timestamp: "\(index + 1)h",
                                likesCount: Int.random(in: 5...500)
                            )
                        }
                    }
                    .padding()
                }
                
                Divider()
                
                // Comment Input
                HStack(spacing: 12) {
                    Circle()
                        .fill(Color.gray.opacity(0.2))
                        .frame(width: 28, height: 28)
                    
                    TextField("Add a comment...", text: $commentText)
                        .font(.system(size: 14))
                    
                    if !commentText.isEmpty {
                        Button(action: {
                            commentText = ""
                        }) {
                            Text("Post")
                                .font(.system(size: 14, weight: .semibold))
                                .foregroundColor(EntativaColors.primaryBlue)
                        }
                    }
                }
                .padding()
            }
        }
    }
}

struct VignetteTakeCommentRow: View {
    let username: String
    let comment: String
    let timestamp: String
    let likesCount: Int
    @State private var isLiked = false
    
    var body: some View {
        HStack(alignment: .top, spacing: 12) {
            Circle()
                .fill(Color.gray.opacity(0.2))
                .frame(width: 32, height: 32)
            
            VStack(alignment: .leading, spacing: 4) {
                Text(username)
                    .font(.system(size: 13, weight: .semibold))
                    .foregroundColor(VignetteColors.textPrimary)
                
                Text(comment)
                    .font(.system(size: 13))
                    .foregroundColor(VignetteColors.textPrimary)
                
                HStack(spacing: 12) {
                    Text(timestamp)
                        .font(.system(size: 12))
                        .foregroundColor(VignetteColors.textSecondary)
                    
                    if likesCount > 0 {
                        Text("\(likesCount) likes")
                            .font(.system(size: 12))
                            .foregroundColor(VignetteColors.textSecondary)
                    }
                    
                    Button(action: {}) {
                        Text("Reply")
                            .font(.system(size: 12))
                            .foregroundColor(VignetteColors.textSecondary)
                    }
                }
            }
            
            Spacer()
            
            Button(action: {
                withAnimation {
                    isLiked.toggle()
                }
            }) {
                Image(systemName: isLiked ? "heart.fill" : "heart")
                    .font(.system(size: 12))
                    .foregroundColor(isLiked ? .red : VignetteColors.textSecondary)
            }
        }
    }
}

// MARK: - Share Sheet
struct VignetteTakeShareSheet: View {
    let take: VignetteTake
    @Environment(\.dismiss) var dismiss
    
    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 24) {
                    // Share to Stories
                    VStack(spacing: 16) {
                        Button(action: {}) {
                            HStack {
                                Image(systemName: "plus.circle.fill")
                                    .font(.system(size: 24))
                                    .foregroundColor(EntativaColors.primaryBlue)
                                
                                VStack(alignment: .leading) {
                                    Text("Share to Story")
                                        .font(.system(size: 15, weight: .medium))
                                    Text("Share this take with your followers")
                                        .font(.system(size: 13))
                                        .foregroundColor(VignetteColors.textSecondary)
                                }
                                
                                Spacer()
                            }
                            .padding()
                            .background(Color(UIColor.systemGray6))
                            .cornerRadius(12)
                        }
                    }
                    .padding(.horizontal)
                    
                    Divider()
                    
                    // Share Options
                    LazyVGrid(columns: [
                        GridItem(.flexible()),
                        GridItem(.flexible()),
                        GridItem(.flexible())
                    ], spacing: 24) {
                        VignetteShareOption(icon: "link", title: "Copy Link")
                        VignetteShareOption(icon: "square.and.arrow.up", title: "Share")
                        VignetteShareOption(icon: "bookmark", title: "Save")
                        VignetteShareOption(icon: "person.badge.minus", title: "Not Interested")
                        VignetteShareOption(icon: "exclamationmark.bubble", title: "Report")
                        VignetteShareOption(icon: "eye.slash", title: "Hide")
                    }
                    .padding()
                }
            }
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .principal) {
                    Text("Share")
                        .font(.system(size: 16, weight: .semibold))
                }
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button(action: { dismiss() }) {
                        Image(systemName: "xmark")
                            .foregroundColor(VignetteColors.textPrimary)
                    }
                }
            }
        }
    }
}

struct VignetteShareOption: View {
    let icon: String
    let title: String
    
    var body: some View {
        VStack(spacing: 8) {
            Circle()
                .stroke(VignetteColors.separator, lineWidth: 1)
                .frame(width: 56, height: 56)
                .overlay(
                    Image(systemName: icon)
                        .font(.system(size: 22))
                        .foregroundColor(VignetteColors.textPrimary)
                )
            
            Text(title)
                .font(.system(size: 12))
                .foregroundColor(VignetteColors.textPrimary)
                .multilineTextAlignment(.center)
                .lineLimit(2)
        }
    }
}

// MARK: - View Model
class VignetteTakesViewModel: ObservableObject {
    @Published var takes: [VignetteTake] = VignetteTake.mockTakes
    @Published var isLoading = false
    
    func loadMore() {
        // Load more takes
    }
}

// MARK: - Model
struct VignetteTake: Identifiable {
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
    
    static let mockTakes: [VignetteTake] = [
        VignetteTake(
            id: "1",
            userId: "user1",
            username: "photo.vibes",
            userAvatar: nil,
            videoUrl: "",
            thumbnailUrl: nil,
            caption: "Golden hour in the city 🌆✨ #photography",
            audioName: "Chill Vibes - Lofi Beats",
            likesCount: 234_500,
            commentsCount: 3_421,
            sharesCount: 8_934,
            viewsCount: 1_234_500,
            createdAt: Date()
        ),
        VignetteTake(
            id: "2",
            userId: "user2",
            username: "fit.journey",
            userAvatar: nil,
            videoUrl: "",
            thumbnailUrl: nil,
            caption: "Morning routine that changed my life! 💪",
            audioName: "Workout Mix 2024",
            likesCount: 89_300,
            commentsCount: 1_234,
            sharesCount: 3_456,
            viewsCount: 567_800,
            createdAt: Date()
        ),
        VignetteTake(
            id: "3",
            userId: "user3",
            username: "chef.athome",
            userAvatar: nil,
            videoUrl: "",
            thumbnailUrl: nil,
            caption: "60-second pasta that tastes like heaven! 🍝",
            audioName: "Cooking Time - Kitchen Beats",
            likesCount: 456_700,
            commentsCount: 12_345,
            sharesCount: 23_456,
            viewsCount: 2_345_600,
            createdAt: Date()
        ),
        VignetteTake(
            id: "4",
            userId: "user4",
            username: "style.daily",
            userAvatar: nil,
            videoUrl: "",
            thumbnailUrl: nil,
            caption: "Transforming thrift finds into designer looks ✨👗",
            audioName: "Fashion Week Runway",
            likesCount: 678_900,
            commentsCount: 8_765,
            sharesCount: 34_567,
            viewsCount: 3_456_700,
            createdAt: Date()
        ),
        VignetteTake(
            id: "5",
            userId: "user5",
            username: "pet.moments",
            userAvatar: nil,
            videoUrl: "",
            thumbnailUrl: nil,
            caption: "When your dog understands the assignment 😂🐕",
            audioName: "Funny Pet Sounds",
            likesCount: 890_100,
            commentsCount: 23_456,
            sharesCount: 45_678,
            viewsCount: 4_567_800,
            createdAt: Date()
        )
    ]
}

// MARK: - Preview
struct VignetteTakesView_Previews: PreviewProvider {
    static var previews: some View {
        VignetteTakesView()
    }
}
