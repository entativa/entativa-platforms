import SwiftUI

// MARK: - Vignette Profile View (Full-Bleed Immersive Design)
struct VignetteProfileView: View {
    @StateObject private var viewModel = ProfileViewModel()
    @State private var showEditProfile = false
    @State private var showSettings = false
    @State private var selectedTab: ProfileTab = .posts
    
    var body: some View {
        ZStack {
            // Full-bleed background image
            if let profileImage = viewModel.profile.profileImageUrl {
                AsyncImage(url: URL(string: profileImage)) { image in
                    image
                        .resizable()
                        .aspectRatio(contentMode: .fill)
                        .frame(maxWidth: .infinity, maxHeight: .infinity)
                        .clipped()
                } placeholder: {
                    LinearGradient(
                        colors: [
                            Color(hex: "C3E7F1"),
                            Color(hex: "519CAB")
                        ],
                        startPoint: .topLeading,
                        endPoint: .bottomTrailing
                    )
                }
            } else {
                // Default gradient background
                LinearGradient(
                    colors: [
                        Color(hex: "C3E7F1"),
                        Color(hex: "519CAB"),
                        Color(hex: "20373B")
                    ],
                    startPoint: .top,
                    endPoint: .bottom
                )
            }
            
            // Dark gradient overlay for readability
            LinearGradient(
                colors: [
                    Color.black.opacity(0.6),
                    Color.black.opacity(0.3),
                    Color.black.opacity(0.6)
                ],
                startPoint: .top,
                endPoint: .bottom
            )
            
            // Content layer
            ScrollView {
                VStack(spacing: 0) {
                    // Top section with profile info
                    VStack(spacing: 16) {
                        // Header buttons
                        HStack {
                            Button(action: {}) {
                                Image(systemName: "lock.fill")
                                    .font(.system(size: 14, weight: .semibold))
                                    .foregroundColor(.white)
                                Text(viewModel.profile.username)
                                    .font(.system(size: 18, weight: .bold))
                                    .foregroundColor(.white)
                                Image(systemName: "chevron.down")
                                    .font(.system(size: 14, weight: .semibold))
                                    .foregroundColor(.white)
                            }
                            
                            Spacer()
                            
                            HStack(spacing: 16) {
                                Button(action: {}) {
                                    Image(systemName: "plus.app")
                                        .font(.system(size: 24))
                                        .foregroundColor(.white)
                                }
                                
                                Button(action: { showSettings = true }) {
                                    Image(systemName: "line.3.horizontal")
                                        .font(.system(size: 24))
                                        .foregroundColor(.white)
                                }
                            }
                        }
                        .padding(.horizontal, 16)
                        .padding(.top, 50)
                        
                        // Profile picture (circular with gradient border)
                        ZStack {
                            Circle()
                                .fill(
                                    LinearGradient(
                                        colors: [
                                            Color(hex: "FFC64F"),
                                            Color(hex: "FC30E1"),
                                            Color(hex: "6F3EFB")
                                        ],
                                        startPoint: .topLeading,
                                        endPoint: .bottomTrailing
                                    )
                                )
                                .frame(width: 90, height: 90)
                            
                            Circle()
                                .fill(Color.white.opacity(0.2))
                                .frame(width: 84, height: 84)
                            
                            if let avatarUrl = viewModel.profile.avatarUrl {
                                AsyncImage(url: URL(string: avatarUrl)) { image in
                                    image
                                        .resizable()
                                        .aspectRatio(contentMode: .fill)
                                } placeholder: {
                                    Image(systemName: "person.circle.fill")
                                        .resizable()
                                        .foregroundColor(.white.opacity(0.5))
                                }
                                .frame(width: 80, height: 80)
                                .clipShape(Circle())
                            } else {
                                Image(systemName: "person.circle.fill")
                                    .resizable()
                                    .foregroundColor(.white.opacity(0.5))
                                    .frame(width: 80, height: 80)
                            }
                        }
                        .padding(.top, 20)
                        
                        // Name and bio
                        VStack(spacing: 8) {
                            Text(viewModel.profile.fullName)
                                .font(.system(size: 16, weight: .semibold))
                                .foregroundColor(.white)
                            
                            if let bio = viewModel.profile.bio {
                                Text(bio)
                                    .font(.system(size: 14))
                                    .foregroundColor(.white.opacity(0.9))
                                    .multilineTextAlignment(.center)
                                    .padding(.horizontal, 32)
                            }
                            
                            if let link = viewModel.profile.link {
                                Button(action: {}) {
                                    Text(link)
                                        .font(.system(size: 14, weight: .medium))
                                        .foregroundColor(Color(hex: "FFC64F"))
                                }
                            }
                        }
                        .padding(.horizontal, 16)
                        
                        // Stats row (frosted glass)
                        HStack(spacing: 0) {
                            StatButton(
                                count: viewModel.profile.postsCount,
                                label: "Posts"
                            )
                            
                            StatButton(
                                count: viewModel.profile.followersCount,
                                label: "Followers"
                            )
                            
                            StatButton(
                                count: viewModel.profile.followingCount,
                                label: "Following"
                            )
                        }
                        .padding(.horizontal, 16)
                        .padding(.top, 20)
                        
                        // Action buttons (frosted glass)
                        HStack(spacing: 12) {
                            Button(action: { showEditProfile = true }) {
                                Text("Edit Profile")
                                    .font(.system(size: 14, weight: .semibold))
                                    .foregroundColor(.white)
                                    .frame(maxWidth: .infinity)
                                    .frame(height: 32)
                                    .background(
                                        RoundedRectangle(cornerRadius: 8)
                                            .fill(.ultraThinMaterial)
                                    )
                            }
                            
                            Button(action: {}) {
                                Text("Share Profile")
                                    .font(.system(size: 14, weight: .semibold))
                                    .foregroundColor(.white)
                                    .frame(maxWidth: .infinity)
                                    .frame(height: 32)
                                    .background(
                                        RoundedRectangle(cornerRadius: 8)
                                            .fill(.ultraThinMaterial)
                                    )
                            }
                        }
                        .padding(.horizontal, 16)
                        .padding(.top, 12)
                        
                        // Story highlights (frosted glass circles)
                        ScrollView(.horizontal, showsIndicators: false) {
                            HStack(spacing: 16) {
                                // Add highlight button
                                VStack(spacing: 8) {
                                    Circle()
                                        .fill(.ultraThinMaterial)
                                        .frame(width: 64, height: 64)
                                        .overlay(
                                            Image(systemName: "plus")
                                                .font(.system(size: 24, weight: .medium))
                                                .foregroundColor(.white)
                                        )
                                    
                                    Text("New")
                                        .font(.system(size: 12))
                                        .foregroundColor(.white.opacity(0.9))
                                }
                                
                                // Existing highlights
                                ForEach(viewModel.highlights) { highlight in
                                    HighlightView(highlight: highlight)
                                }
                            }
                            .padding(.horizontal, 16)
                        }
                        .padding(.top, 20)
                    }
                    .padding(.bottom, 20)
                    
                    // Tab selector (frosted glass)
                    HStack(spacing: 0) {
                        ProfileTabButton(
                            icon: "square.grid.3x3",
                            isSelected: selectedTab == .posts,
                            action: { selectedTab = .posts }
                        )
                        
                        ProfileTabButton(
                            icon: "play.rectangle",
                            isSelected: selectedTab == .reels,
                            action: { selectedTab = .reels }
                        )
                        
                        ProfileTabButton(
                            icon: "person.crop.square",
                            isSelected: selectedTab == .tagged,
                            action: { selectedTab = .tagged }
                        )
                    }
                    .frame(height: 44)
                    .background(.ultraThinMaterial)
                    
                    // Content grid
                    LazyVGrid(
                        columns: [
                            GridItem(.flexible(), spacing: 2),
                            GridItem(.flexible(), spacing: 2),
                            GridItem(.flexible(), spacing: 2)
                        ],
                        spacing: 2
                    ) {
                        ForEach(viewModel.posts) { post in
                            PostGridItem(post: post)
                        }
                    }
                    .background(.ultraThinMaterial)
                }
            }
            .ignoresSafeArea()
        }
        .sheet(isPresented: $showEditProfile) {
            EditProfileView(profile: viewModel.profile)
        }
        .sheet(isPresented: $showSettings) {
            SettingsView()
        }
    }
}

// MARK: - Stat Button (Frosted Glass)
struct StatButton: View {
    let count: Int
    let label: String
    
    var body: some View {
        Button(action: {}) {
            VStack(spacing: 4) {
                Text(formatCount(count))
                    .font(.system(size: 18, weight: .bold))
                    .foregroundColor(.white)
                
                Text(label)
                    .font(.system(size: 12))
                    .foregroundColor(.white.opacity(0.8))
            }
            .frame(maxWidth: .infinity)
            .padding(.vertical, 8)
        }
    }
}

// MARK: - Profile Tab Button
struct ProfileTabButton: View {
    let icon: String
    let isSelected: Bool
    let action: () -> Void
    
    var body: some View {
        Button(action: action) {
            Image(systemName: icon)
                .font(.system(size: 20))
                .foregroundColor(isSelected ? .white : .white.opacity(0.5))
                .frame(maxWidth: .infinity)
                .frame(height: 44)
                .background(
                    isSelected ?
                        AnyView(
                            Rectangle()
                                .fill(Color.white.opacity(0.2))
                        ) :
                        AnyView(Color.clear)
                )
        }
    }
}

// MARK: - Highlight View
struct HighlightView: View {
    let highlight: Highlight
    
    var body: some View {
        VStack(spacing: 8) {
            Circle()
                .fill(.ultraThinMaterial)
                .frame(width: 64, height: 64)
                .overlay(
                    AsyncImage(url: URL(string: highlight.coverImage)) { image in
                        image
                            .resizable()
                            .aspectRatio(contentMode: .fill)
                    } placeholder: {
                        Image(systemName: "photo")
                            .foregroundColor(.white.opacity(0.5))
                    }
                    .frame(width: 60, height: 60)
                    .clipShape(Circle())
                )
            
            Text(highlight.title)
                .font(.system(size: 12))
                .foregroundColor(.white.opacity(0.9))
                .lineLimit(1)
        }
    }
}

// MARK: - Post Grid Item
struct PostGridItem: View {
    let post: ProfilePost
    
    var body: some View {
        GeometryReader { geometry in
            AsyncImage(url: URL(string: post.thumbnailUrl)) { image in
                image
                    .resizable()
                    .aspectRatio(contentMode: .fill)
            } placeholder: {
                Rectangle()
                    .fill(Color.white.opacity(0.1))
                    .overlay(
                        Image(systemName: "photo")
                            .foregroundColor(.white.opacity(0.3))
                    )
            }
            .frame(width: geometry.size.width, height: geometry.size.width)
            .clipped()
            .overlay(
                // Video indicator
                post.isVideo ?
                    AnyView(
                        Image(systemName: "play.fill")
                            .font(.system(size: 14))
                            .foregroundColor(.white)
                            .padding(4)
                            .background(Color.black.opacity(0.4))
                            .clipShape(Circle())
                            .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .topTrailing)
                            .padding(8)
                    ) :
                    AnyView(EmptyView())
            )
        }
        .aspectRatio(1, contentMode: .fit)
    }
}

// MARK: - Edit Profile View
struct EditProfileView: View {
    let profile: UserProfile
    @Environment(\.dismiss) var dismiss
    
    @State private var name: String
    @State private var bio: String
    @State private var link: String
    
    init(profile: UserProfile) {
        self.profile = profile
        _name = State(initialValue: profile.fullName)
        _bio = State(initialValue: profile.bio ?? "")
        _link = State(initialValue: profile.link ?? "")
    }
    
    var body: some View {
        NavigationView {
            Form {
                Section {
                    HStack {
                        Spacer()
                        
                        ZStack {
                            Circle()
                                .fill(Color.gray.opacity(0.2))
                                .frame(width: 80, height: 80)
                            
                            Image(systemName: "person.circle.fill")
                                .resizable()
                                .frame(width: 80, height: 80)
                                .foregroundColor(.gray)
                        }
                        
                        Spacer()
                    }
                    .padding(.vertical, 8)
                    
                    Button("Change Profile Photo") {}
                        .frame(maxWidth: .infinity)
                }
                
                Section(header: Text("Profile Information")) {
                    TextField("Name", text: $name)
                    TextField("Bio", text: $bio, axis: .vertical)
                        .lineLimit(3...6)
                    TextField("Link", text: $link)
                        .autocapitalization(.none)
                        .keyboardType(.URL)
                }
            }
            .navigationTitle("Edit Profile")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Cancel") {
                        dismiss()
                    }
                }
                
                ToolbarItem(placement: .confirmationAction) {
                    Button("Done") {
                        // Save profile
                        dismiss()
                    }
                    .fontWeight(.semibold)
                }
            }
        }
    }
}

// MARK: - Settings View
struct SettingsView: View {
    @Environment(\.dismiss) var dismiss
    
    var body: some View {
        NavigationView {
            List {
                Section {
                    NavigationLink(destination: Text("Account Settings")) {
                        Label("Account", systemImage: "person.circle")
                    }
                    
                    NavigationLink(destination: Text("Privacy Settings")) {
                        Label("Privacy", systemImage: "lock.circle")
                    }
                    
                    NavigationLink(destination: Text("Notifications")) {
                        Label("Notifications", systemImage: "bell.circle")
                    }
                }
                
                Section {
                    NavigationLink(destination: Text("Saved")) {
                        Label("Saved", systemImage: "bookmark")
                    }
                    
                    NavigationLink(destination: Text("Archive")) {
                        Label("Archive", systemImage: "archivebox")
                    }
                    
                    NavigationLink(destination: Text("Your Activity")) {
                        Label("Your Activity", systemImage: "clock")
                    }
                }
                
                Section {
                    Button(action: {}) {
                        Label("Log Out", systemImage: "rectangle.portrait.and.arrow.right")
                            .foregroundColor(.red)
                    }
                }
            }
            .navigationTitle("Settings")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Close") {
                        dismiss()
                    }
                }
            }
        }
    }
}

// MARK: - Profile Tab Enum
enum ProfileTab {
    case posts
    case reels
    case tagged
}

// MARK: - Models
struct UserProfile {
    let id: String
    let username: String
    let fullName: String
    let bio: String?
    let link: String?
    let avatarUrl: String?
    let profileImageUrl: String?
    let postsCount: Int
    let followersCount: Int
    let followingCount: Int
    let isVerified: Bool
    let isPrivate: Bool
}

struct Highlight: Identifiable {
    let id: String
    let title: String
    let coverImage: String
}

struct ProfilePost: Identifiable {
    let id: String
    let thumbnailUrl: String
    let isVideo: Bool
    let likesCount: Int
}

// MARK: - View Model
class ProfileViewModel: ObservableObject {
    @Published var profile: UserProfile
    @Published var highlights: [Highlight] = []
    @Published var posts: [ProfilePost] = []
    @Published var isLoading = false
    
    init() {
        // Mock data
        self.profile = UserProfile(
            id: "user123",
            username: "yourusername",
            fullName: "Your Name",
            bio: "✨ Living life to the fullest\n📍 San Francisco, CA\n💼 Creator & Entrepreneur",
            link: "yourwebsite.com",
            avatarUrl: nil,
            profileImageUrl: nil,
            postsCount: 142,
            followersCount: 12500,
            followingCount: 890,
            isVerified: false,
            isPrivate: false
        )
        
        // Mock highlights
        self.highlights = [
            Highlight(id: "1", title: "Travel", coverImage: ""),
            Highlight(id: "2", title: "Food", coverImage: ""),
            Highlight(id: "3", title: "Work", coverImage: "")
        ]
        
        // Mock posts
        self.posts = (1...24).map { index in
            ProfilePost(
                id: "post\(index)",
                thumbnailUrl: "",
                isVideo: index % 5 == 0,
                likesCount: Int.random(in: 100...5000)
            )
        }
    }
    
    func loadProfile() async {
        // TODO: Load from API
    }
}

// Helper function
private func formatCount(_ count: Int) -> String {
    if count >= 1_000_000 {
        return String(format: "%.1fM", Double(count) / 1_000_000)
    } else if count >= 1_000 {
        return String(format: "%.1fK", Double(count) / 1_000)
    } else {
        return "\(count)"
    }
}
