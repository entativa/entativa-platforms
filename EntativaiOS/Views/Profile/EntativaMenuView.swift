import SwiftUI

// MARK: - Entativa Menu View (Facebook-Style Profile)
struct EntativaMenuView: View {
    @StateObject private var viewModel = MenuViewModel()
    @State private var showSettings = false
    
    var body: some View {
        ScrollView {
            VStack(spacing: 0) {
                // Profile Header Section
                VStack(spacing: 16) {
                    // Cover photo + profile pic
                    ZStack(alignment: .bottomLeading) {
                        // Cover photo
                        if let coverImage = viewModel.profile.coverImageUrl {
                            AsyncImage(url: URL(string: coverImage)) { image in
                                image
                                    .resizable()
                                    .aspectRatio(contentMode: .fill)
                            } placeholder: {
                                LinearGradient(
                                    colors: [
                                        Color(hex: "007CFC"),
                                        Color(hex: "6F3EFB"),
                                        Color(hex: "FC30E1")
                                    ],
                                    startPoint: .topLeading,
                                    endPoint: .bottomTrailing
                                )
                            }
                            .frame(height: 180)
                            .clipped()
                        } else {
                            LinearGradient(
                                colors: [
                                    Color(hex: "007CFC"),
                                    Color(hex: "6F3EFB"),
                                    Color(hex: "FC30E1")
                                ],
                                startPoint: .topLeading,
                                endPoint: .bottomTrailing
                            )
                            .frame(height: 180)
                        }
                        
                        // Profile picture
                        HStack {
                            ZStack {
                                Circle()
                                    .fill(Color.white)
                                    .frame(width: 120, height: 120)
                                
                                if let avatarUrl = viewModel.profile.avatarUrl {
                                    AsyncImage(url: URL(string: avatarUrl)) { image in
                                        image
                                            .resizable()
                                            .aspectRatio(contentMode: .fill)
                                    } placeholder: {
                                        Image(systemName: "person.circle.fill")
                                            .resizable()
                                            .foregroundColor(.gray)
                                    }
                                    .frame(width: 112, height: 112)
                                    .clipShape(Circle())
                                } else {
                                    Image(systemName: "person.circle.fill")
                                        .resizable()
                                        .foregroundColor(.gray)
                                        .frame(width: 112, height: 112)
                                }
                                
                                // Camera button
                                Circle()
                                    .fill(Color.gray.opacity(0.9))
                                    .frame(width: 32, height: 32)
                                    .overlay(
                                        Image(systemName: "camera.fill")
                                            .font(.system(size: 14))
                                            .foregroundColor(.white)
                                    )
                                    .offset(x: 40, y: 40)
                            }
                            .padding(.leading, 16)
                            .offset(y: 30)
                            
                            Spacer()
                        }
                    }
                    
                    // Profile info
                    VStack(alignment: .leading, spacing: 8) {
                        HStack(alignment: .top) {
                            VStack(alignment: .leading, spacing: 4) {
                                Text(viewModel.profile.fullName)
                                    .font(.system(size: 24, weight: .bold))
                                
                                Text("@\(viewModel.profile.username)")
                                    .font(.system(size: 15))
                                    .foregroundColor(.gray)
                                
                                if let bio = viewModel.profile.bio {
                                    Text(bio)
                                        .font(.system(size: 15))
                                        .foregroundColor(.primary)
                                        .padding(.top, 8)
                                }
                            }
                            
                            Spacer()
                            
                            Button(action: { showSettings = true }) {
                                Image(systemName: "gearshape.fill")
                                    .font(.system(size: 22))
                                    .foregroundColor(.gray)
                            }
                        }
                        .padding(.horizontal, 16)
                        .padding(.top, 40)
                        
                        // Stats row
                        HStack(spacing: 20) {
                            StatView(count: viewModel.profile.friendsCount, label: "Friends")
                            StatView(count: viewModel.profile.followersCount, label: "Followers")
                            StatView(count: viewModel.profile.followingCount, label: "Following")
                        }
                        .padding(.horizontal, 16)
                        .padding(.top, 12)
                        
                        // Action buttons
                        HStack(spacing: 12) {
                            Button(action: {}) {
                                HStack {
                                    Image(systemName: "plus")
                                        .font(.system(size: 14, weight: .semibold))
                                    Text("Add Story")
                                        .font(.system(size: 15, weight: .semibold))
                                }
                                .frame(maxWidth: .infinity)
                                .frame(height: 36)
                                .background(Color(hex: "007CFC"))
                                .foregroundColor(.white)
                                .cornerRadius(8)
                            }
                            
                            Button(action: {}) {
                                HStack {
                                    Image(systemName: "pencil")
                                        .font(.system(size: 14, weight: .semibold))
                                    Text("Edit Profile")
                                        .font(.system(size: 15, weight: .semibold))
                                }
                                .frame(maxWidth: .infinity)
                                .frame(height: 36)
                                .background(Color.gray.opacity(0.2))
                                .foregroundColor(.primary)
                                .cornerRadius(8)
                            }
                            
                            Button(action: {}) {
                                Image(systemName: "ellipsis")
                                    .font(.system(size: 14, weight: .semibold))
                                    .frame(width: 36, height: 36)
                                    .background(Color.gray.opacity(0.2))
                                    .foregroundColor(.primary)
                                    .cornerRadius(8)
                            }
                        }
                        .padding(.horizontal, 16)
                        .padding(.top, 12)
                    }
                }
                .background(Color(UIColor.systemBackground))
                
                Divider()
                    .padding(.top, 16)
                
                // Menu sections
                VStack(spacing: 0) {
                    // Your Shortcuts
                    MenuSection(title: "Your Shortcuts") {
                        MenuItemRow(
                            icon: "person.2.fill",
                            iconColor: Color(hex: "007CFC"),
                            title: "Friends",
                            action: {}
                        )
                        
                        MenuItemRow(
                            icon: "clock.fill",
                            iconColor: Color(hex: "6F3EFB"),
                            title: "Memories",
                            action: {}
                        )
                        
                        MenuItemRow(
                            icon: "bookmark.fill",
                            iconColor: Color(hex: "FC30E1"),
                            title: "Saved",
                            action: {}
                        )
                        
                        MenuItemRow(
                            icon: "flag.fill",
                            iconColor: .orange,
                            title: "Pages",
                            action: {}
                        )
                        
                        MenuItemRow(
                            icon: "play.rectangle.fill",
                            iconColor: .blue,
                            title: "Video",
                            action: {}
                        )
                        
                        MenuItemRow(
                            icon: "storefront.fill",
                            iconColor: .cyan,
                            title: "Marketplace",
                            action: {}
                        )
                        
                        MenuItemRow(
                            icon: "calendar",
                            iconColor: .red,
                            title: "Events",
                            action: {}
                        )
                        
                        Button(action: {}) {
                            HStack {
                                Image(systemName: "chevron.down")
                                    .font(.system(size: 14, weight: .semibold))
                                    .foregroundColor(.gray)
                                
                                Text("See More")
                                    .font(.system(size: 15, weight: .semibold))
                                    .foregroundColor(.gray)
                                
                                Spacer()
                            }
                            .padding(.horizontal, 16)
                            .padding(.vertical, 12)
                        }
                    }
                    
                    Divider()
                    
                    // Settings & Privacy
                    MenuSection(title: "Settings & Privacy") {
                        MenuItemRow(
                            icon: "gearshape.fill",
                            iconColor: .gray,
                            title: "Settings",
                            action: { showSettings = true }
                        )
                        
                        MenuItemRow(
                            icon: "lock.shield.fill",
                            iconColor: .gray,
                            title: "Privacy Checkup",
                            action: {}
                        )
                        
                        MenuItemRow(
                            icon: "questionmark.circle.fill",
                            iconColor: .gray,
                            title: "Help & Support",
                            action: {}
                        )
                    }
                    
                    Divider()
                    
                    // Account Actions
                    MenuSection(title: "") {
                        MenuItemRow(
                            icon: "moon.fill",
                            iconColor: .gray,
                            title: "Dark Mode",
                            action: {}
                        )
                        
                        MenuItemRow(
                            icon: "bell.fill",
                            iconColor: .gray,
                            title: "Notification Settings",
                            action: {}
                        )
                        
                        Button(action: {}) {
                            HStack {
                                Image(systemName: "rectangle.portrait.and.arrow.right")
                                    .font(.system(size: 20))
                                    .foregroundColor(.red)
                                    .frame(width: 36)
                                
                                Text("Log Out")
                                    .font(.system(size: 15))
                                    .foregroundColor(.red)
                                
                                Spacer()
                            }
                            .padding(.horizontal, 16)
                            .padding(.vertical, 12)
                        }
                    }
                }
                .padding(.top, 8)
            }
        }
        .background(Color(UIColor.systemGroupedBackground))
        .sheet(isPresented: $showSettings) {
            SettingsMenuView()
        }
    }
}

// MARK: - Stat View
struct StatView: View {
    let count: Int
    let label: String
    
    var body: some View {
        VStack(spacing: 4) {
            Text("\(count)")
                .font(.system(size: 18, weight: .bold))
            
            Text(label)
                .font(.system(size: 13))
                .foregroundColor(.gray)
        }
    }
}

// MARK: - Menu Section
struct MenuSection<Content: View>: View {
    let title: String
    let content: Content
    
    init(title: String, @ViewBuilder content: () -> Content) {
        self.title = title
        self.content = content()
    }
    
    var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            if !title.isEmpty {
                Text(title)
                    .font(.system(size: 17, weight: .semibold))
                    .padding(.horizontal, 16)
                    .padding(.vertical, 12)
            }
            
            content
        }
        .background(Color(UIColor.systemBackground))
    }
}

// MARK: - Menu Item Row
struct MenuItemRow: View {
    let icon: String
    let iconColor: Color
    let title: String
    let action: () -> Void
    
    var body: some View {
        Button(action: action) {
            HStack(spacing: 12) {
                Image(systemName: icon)
                    .font(.system(size: 20))
                    .foregroundColor(iconColor)
                    .frame(width: 36)
                
                Text(title)
                    .font(.system(size: 15))
                    .foregroundColor(.primary)
                
                Spacer()
                
                Image(systemName: "chevron.right")
                    .font(.system(size: 14, weight: .semibold))
                    .foregroundColor(.gray)
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)
        }
    }
}

// MARK: - Settings Menu View
struct SettingsMenuView: View {
    @Environment(\.dismiss) var dismiss
    
    var body: some View {
        NavigationView {
            List {
                Section("Account") {
                    NavigationLink(destination: Text("Personal Information")) {
                        Label("Personal Information", systemImage: "person.text.rectangle")
                    }
                    
                    NavigationLink(destination: Text("Password and Security")) {
                        Label("Password and Security", systemImage: "lock.shield")
                    }
                    
                    NavigationLink(destination: Text("Your Activity")) {
                        Label("Your Activity", systemImage: "clock")
                    }
                }
                
                Section("Preferences") {
                    NavigationLink(destination: Text("Notifications")) {
                        Label("Notifications", systemImage: "bell")
                    }
                    
                    NavigationLink(destination: Text("Privacy")) {
                        Label("Privacy", systemImage: "hand.raised")
                    }
                    
                    NavigationLink(destination: Text("Language")) {
                        Label("Language", systemImage: "globe")
                    }
                }
                
                Section("Support") {
                    NavigationLink(destination: Text("Help Center")) {
                        Label("Help Center", systemImage: "questionmark.circle")
                    }
                    
                    NavigationLink(destination: Text("Report a Problem")) {
                        Label("Report a Problem", systemImage: "exclamationmark.triangle")
                    }
                    
                    NavigationLink(destination: Text("About")) {
                        Label("About", systemImage: "info.circle")
                    }
                }
            }
            .navigationTitle("Settings & Privacy")
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

// MARK: - Menu Profile Model
struct MenuProfile {
    let id: String
    let username: String
    let fullName: String
    let bio: String?
    let avatarUrl: String?
    let coverImageUrl: String?
    let friendsCount: Int
    let followersCount: Int
    let followingCount: Int
}

// MARK: - View Model
class MenuViewModel: ObservableObject {
    @Published var profile: MenuProfile
    @Published var isLoading = false
    
    init() {
        // Mock data
        self.profile = MenuProfile(
            id: "user123",
            username: "yourname",
            fullName: "Your Name",
            bio: "Welcome to my profile! 👋",
            avatarUrl: nil,
            coverImageUrl: nil,
            friendsCount: 342,
            followersCount: 1520,
            followingCount: 487
        )
    }
    
    func loadProfile() async {
        // TODO: Load from API
    }
}
