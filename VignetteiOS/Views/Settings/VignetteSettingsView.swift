import SwiftUI

// MARK: - Vignette Settings View (Instagram-Style from Profile)
struct VignetteSettingsView: View {
    @Environment(\.dismiss) var dismiss
    @StateObject private var viewModel = VignetteSettingsViewModel()
    @State private var showLogoutConfirm = false
    
    var body: some View {
        NavigationView {
            List {
                // Account
                Section {
                    NavigationLink(destination: EditProfileView(viewModel: viewModel)) {
                        Label("Edit Profile", systemImage: "person.circle")
                    }
                    
                    NavigationLink(destination: ChangePasswordView()) {
                        Label("Change Password", systemImage: "key.fill")
                    }
                    
                    NavigationLink(destination: AccountPrivacyView(viewModel: viewModel)) {
                        Label("Account Privacy", systemImage: "lock.fill")
                    }
                } header: {
                    Text("Account")
                }
                
                // Content & Activity
                Section {
                    NavigationLink(destination: NotificationPreferencesView(viewModel: viewModel)) {
                        Label("Notifications", systemImage: "bell.fill")
                    }
                    
                    NavigationLink(destination: PostsYouLikedView()) {
                        Label("Posts You've Liked", systemImage: "heart.fill")
                    }
                    
                    NavigationLink(destination: SavedPostsView()) {
                        Label("Saved", systemImage: "bookmark.fill")
                    }
                    
                    NavigationLink(destination: ArchiveView()) {
                        Label("Archive", systemImage: "archivebox.fill")
                    }
                } header: {
                    Text("Content & Activity")
                }
                
                // Security
                Section {
                    NavigationLink(destination: SecurityView()) {
                        Label("Security", systemImage: "shield.fill")
                    }
                    
                    NavigationLink(destination: TwoFactorAuthView()) {
                        Label("Two-Factor Authentication", systemImage: "lock.shield.fill")
                    }
                    
                    NavigationLink(destination: LoginActivityView()) {
                        Label("Login Activity", systemImage: "clock.arrow.circlepath")
                    }
                } header: {
                    Text("Security")
                }
                
                // Privacy
                Section {
                    NavigationLink(destination: PrivacyControlsView(viewModel: viewModel)) {
                        Label("Privacy", systemImage: "hand.raised.fill")
                    }
                    
                    NavigationLink(destination: BlockedAccountsView()) {
                        Label("Blocked Accounts", systemImage: "person.fill.xmark")
                    }
                    
                    NavigationLink(destination: MutedAccountsView()) {
                        Label("Muted Accounts", systemImage: "speaker.slash.fill")
                    }
                    
                    NavigationLink(destination: RestrictedAccountsView()) {
                        Label("Restricted Accounts", systemImage: "exclamationmark.shield.fill")
                    }
                } header: {
                    Text("Privacy")
                }
                
                // Preferences
                Section {
                    NavigationLink(destination: LanguageView()) {
                        HStack {
                            Label("Language", systemImage: "globe")
                            Spacer()
                            Text("English")
                                .foregroundColor(.gray)
                        }
                    }
                    
                    Toggle(isOn: $viewModel.darkModeEnabled) {
                        Label("Dark Mode", systemImage: "moon.fill")
                    }
                    
                    NavigationLink(destination: DataUsageView(viewModel: viewModel)) {
                        Label("Data Usage", systemImage: "arrow.up.arrow.down")
                    }
                } header: {
                    Text("Preferences")
                }
                
                // Help
                Section {
                    NavigationLink(destination: HelpCenterView()) {
                        Label("Help", systemImage: "questionmark.circle")
                    }
                    
                    NavigationLink(destination: AboutView()) {
                        Label("About", systemImage: "info.circle")
                    }
                    
                    NavigationLink(destination: AccountStatusView()) {
                        Label("Account Status", systemImage: "checkmark.shield")
                    }
                } header: {
                    Text("Help & Support")
                }
                
                // Switch Account
                Section {
                    NavigationLink(destination: Text("Switch to Entativa")) {
                        HStack {
                            Image(systemName: "arrow.left.arrow.right.circle.fill")
                                .foregroundColor(Color(hex: "007CFC"))
                            Text("Switch to Entativa")
                        }
                    }
                }
                
                // Logout
                Section {
                    Button(action: { showLogoutConfirm = true }) {
                        HStack {
                            Spacer()
                            Text("Log Out")
                                .foregroundColor(.red)
                                .fontWeight(.semibold)
                            Spacer()
                        }
                    }
                }
                
                // Legal
                Section {
                    Link("Terms of Service", destination: URL(string: "https://vignette.app/terms")!)
                    Link("Privacy Policy", destination: URL(string: "https://vignette.app/privacy")!)
                    Link("Community Guidelines", destination: URL(string: "https://vignette.app/guidelines")!)
                } header: {
                    Text("Legal")
                }
                
                // App Info
                Section {
                    HStack {
                        Text("Version")
                        Spacer()
                        Text("1.0.0")
                            .foregroundColor(.gray)
                    }
                    .font(.system(size: 14))
                } footer: {
                    Text("© 2025 Vignette, Inc.")
                        .frame(maxWidth: .infinity)
                        .font(.system(size: 12))
                        .foregroundColor(.gray)
                        .padding(.top, 8)
                }
            }
            .navigationTitle("Settings")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Done") {
                        dismiss()
                    }
                }
            }
        }
        .alert("Log Out", isPresented: $showLogoutConfirm) {
            Button("Cancel", role: .cancel) {}
            Button("Log Out", role: .destructive) {
                viewModel.logout()
            }
        } message: {
            Text("Are you sure you want to log out?")
        }
    }
}

// MARK: - Edit Profile View
struct EditProfileView: View {
    @ObservedObject var viewModel: VignetteSettingsViewModel
    @State private var name = ""
    @State private var username = ""
    @State private var bio = ""
    @State private var website = ""
    @State private var isSaving = false
    
    var body: some View {
        List {
            Section {
                HStack {
                    Spacer()
                    VStack(spacing: 12) {
                        Circle()
                            .fill(Color.gray.opacity(0.2))
                            .frame(width: 80, height: 80)
                            .overlay(
                                Image(systemName: "person.fill")
                                    .font(.system(size: 36))
                                    .foregroundColor(.gray)
                            )
                        
                        Button("Change Photo") {
                            // Change photo
                        }
                        .font(.system(size: 14, weight: .semibold))
                        .foregroundColor(Color(hex: "007CFC"))
                    }
                    Spacer()
                }
                .padding(.vertical, 8)
            }
            
            Section {
                VStack(alignment: .leading, spacing: 4) {
                    Text("Name")
                        .font(.system(size: 12, weight: .medium))
                        .foregroundColor(.gray)
                    TextField("Name", text: $name)
                }
                
                VStack(alignment: .leading, spacing: 4) {
                    Text("Username")
                        .font(.system(size: 12, weight: .medium))
                        .foregroundColor(.gray)
                    TextField("Username", text: $username)
                        .autocapitalization(.none)
                }
                
                VStack(alignment: .leading, spacing: 4) {
                    Text("Bio")
                        .font(.system(size: 12, weight: .medium))
                        .foregroundColor(.gray)
                    TextEditor(text: $bio)
                        .frame(height: 80)
                }
                
                VStack(alignment: .leading, spacing: 4) {
                    Text("Website")
                        .font(.system(size: 12, weight: .medium))
                        .foregroundColor(.gray)
                    TextField("Website", text: $website)
                        .autocapitalization(.none)
                        .keyboardType(.URL)
                }
            }
            
            Section {
                Button(action: {
                    isSaving = true
                    viewModel.updateProfile(name: name, username: username, bio: bio, website: website) {
                        isSaving = false
                    }
                }) {
                    HStack {
                        Spacer()
                        if isSaving {
                            ProgressView()
                        } else {
                            Text("Save")
                                .fontWeight(.semibold)
                        }
                        Spacer()
                    }
                }
                .disabled(isSaving)
            }
        }
        .navigationTitle("Edit Profile")
        .onAppear {
            name = viewModel.userName
            username = viewModel.username
            bio = viewModel.userBio
            website = viewModel.userWebsite
        }
    }
}

// MARK: - Account Privacy View
struct AccountPrivacyView: View {
    @ObservedObject var viewModel: VignetteSettingsViewModel
    
    var body: some View {
        List {
            Section {
                Toggle("Private Account", isOn: $viewModel.isPrivateAccount)
                
                Text("When your account is private, only people you approve can see your photos and videos. Your existing followers won't be affected.")
                    .font(.system(size: 13))
                    .foregroundColor(.gray)
            }
            
            Section(header: Text("Interactions")) {
                Picker("Comments", selection: $viewModel.commentSettings) {
                    Text("Everyone").tag("everyone")
                    Text("People You Follow").tag("following")
                    Text("Your Followers").tag("followers")
                    Text("Off").tag("off")
                }
                
                Picker("Mentions", selection: $viewModel.mentionSettings) {
                    Text("Everyone").tag("everyone")
                    Text("People You Follow").tag("following")
                    Text("Off").tag("off")
                }
                
                Picker("Story Sharing", selection: $viewModel.storySharing) {
                    Text("Everyone").tag("everyone")
                    Text("People You Follow").tag("following")
                    Text("Off").tag("off")
                }
            }
            
            Section(header: Text("Messages")) {
                Toggle("Allow Message Requests", isOn: $viewModel.allowMessageRequests)
                Toggle("Show Activity Status", isOn: $viewModel.showActivityStatus)
                Toggle("Read Receipts", isOn: $viewModel.readReceipts)
            }
        }
        .navigationTitle("Account Privacy")
        .onChange(of: viewModel.isPrivateAccount) { _, _ in
            viewModel.savePrivacySettings()
        }
    }
}

// MARK: - Notification Preferences View
struct NotificationPreferencesView: View {
    @ObservedObject var viewModel: VignetteSettingsViewModel
    
    var body: some View {
        List {
            Section(header: Text("Push Notifications")) {
                Toggle("Likes", isOn: $viewModel.notifyLikes)
                Toggle("Comments", isOn: $viewModel.notifyComments)
                Toggle("New Followers", isOn: $viewModel.notifyFollowers)
                Toggle("Direct Messages", isOn: $viewModel.notifyMessages)
                Toggle("Video Views", isOn: $viewModel.notifyVideoViews)
            }
            
            Section(header: Text("Email Notifications")) {
                Toggle("Feedback Emails", isOn: $viewModel.emailFeedback)
                Toggle("Reminder Emails", isOn: $viewModel.emailReminders)
                Toggle("Product Emails", isOn: $viewModel.emailProduct)
            }
            
            Section(header: Text("Live & IGTV")) {
                Toggle("Live Videos", isOn: $viewModel.notifyLiveVideos)
                Toggle("IGTV Video Uploads", isOn: $viewModel.notifyIGTV)
            }
        }
        .navigationTitle("Notifications")
        .onChange(of: viewModel.notifyLikes) { _, _ in
            viewModel.saveNotificationSettings()
        }
    }
}

// MARK: - Privacy Controls View
struct PrivacyControlsView: View {
    @ObservedObject var viewModel: VignetteSettingsViewModel
    
    var body: some View {
        List {
            Section(header: Text("Discoverability")) {
                Toggle("Similar Account Suggestions", isOn: $viewModel.similarAccountSuggestions)
                Toggle("Include in Recommendations", isOn: $viewModel.includeInRecommendations)
            }
            
            Section(header: Text("Data")) {
                NavigationLink(destination: Text("Download Data")) {
                    Text("Download Data")
                }
                
                NavigationLink(destination: Text("Search History")) {
                    Text("Search History")
                }
            }
            
            Section(header: Text("Sensitive Content Control")) {
                Picker("Sensitive Content", selection: $viewModel.sensitiveContentControl) {
                    Text("Allow").tag("allow")
                    Text("Limit").tag("limit")
                    Text("Limit More").tag("limit_more")
                }
            }
        }
        .navigationTitle("Privacy Controls")
    }
}

// MARK: - Data Usage View (Vignette)
struct DataUsageView: View {
    @ObservedObject var viewModel: VignetteSettingsViewModel
    
    var body: some View {
        List {
            Section(header: Text("Cellular Data Use")) {
                Toggle("Use Less Data", isOn: $viewModel.useLessData)
                
                Text("Loads lower resolution photos and videos. Best for saving cellular data.")
                    .font(.system(size: 13))
                    .foregroundColor(.gray)
            }
            
            Section(header: Text("Media Upload Quality")) {
                Picker("Upload Quality", selection: $viewModel.uploadQuality) {
                    Text("High").tag("high")
                    Text("Normal").tag("normal")
                    Text("Basic").tag("basic")
                }
            }
            
            Section(header: Text("Storage")) {
                HStack {
                    Text("Cache Size")
                    Spacer()
                    Text("256 MB")
                        .foregroundColor(.gray)
                }
                
                Button("Clear Cache") {
                    viewModel.clearCache()
                }
                .foregroundColor(.red)
            }
        }
        .navigationTitle("Data Usage")
        .onChange(of: viewModel.useLessData) { _, _ in
            viewModel.saveDataSettings()
        }
    }
}

// MARK: - Helper Views
struct PostsYouLikedView: View {
    var body: some View {
        ScrollView {
            LazyVGrid(columns: [GridItem(.flexible()), GridItem(.flexible()), GridItem(.flexible())], spacing: 2) {
                ForEach(0..<15) { _ in
                    Color.gray.opacity(0.3)
                        .aspectRatio(1, contentMode: .fit)
                }
            }
        }
        .navigationTitle("Posts You've Liked")
    }
}

struct SavedPostsView: View {
    var body: some View {
        ScrollView {
            LazyVGrid(columns: [GridItem(.flexible()), GridItem(.flexible()), GridItem(.flexible())], spacing: 2) {
                ForEach(0..<15) { _ in
                    Color.gray.opacity(0.3)
                        .aspectRatio(1, contentMode: .fit)
                }
            }
        }
        .navigationTitle("Saved")
    }
}

struct ArchiveView: View {
    var body: some View {
        Text("No archived posts")
            .foregroundColor(.gray)
            .navigationTitle("Archive")
    }
}

struct SecurityView: View {
    var body: some View {
        List {
            NavigationLink(destination: Text("Password")) {
                Text("Password")
            }
            NavigationLink(destination: Text("Login Activity")) {
                Text("Login Activity")
            }
            NavigationLink(destination: Text("Saved Login Info")) {
                Text("Saved Login Info")
            }
            NavigationLink(destination: TwoFactorAuthView()) {
                Text("Two-Factor Authentication")
            }
        }
        .navigationTitle("Security")
    }
}

struct TwoFactorAuthView: View {
    @State private var isEnabled = false
    
    var body: some View {
        List {
            Section {
                Toggle("Two-Factor Authentication", isOn: $isEnabled)
                
                Text("Add an extra layer of security to your account. You'll be asked for a security code when you log in from an unrecognized device.")
                    .font(.system(size: 13))
                    .foregroundColor(.gray)
            }
        }
        .navigationTitle("Two-Factor Authentication")
    }
}

struct LoginActivityView: View {
    var body: some View {
        List {
            Section {
                HStack {
                    VStack(alignment: .leading) {
                        Text("iPhone 15 Pro")
                            .font(.system(size: 15, weight: .semibold))
                        Text("San Francisco, CA • Active now")
                            .font(.system(size: 13))
                            .foregroundColor(.gray)
                    }
                    Spacer()
                    Text("✓")
                        .foregroundColor(.green)
                }
            }
        }
        .navigationTitle("Login Activity")
    }
}

struct BlockedAccountsView: View {
    var body: some View {
        List {
            Text("No blocked accounts")
                .foregroundColor(.gray)
        }
        .navigationTitle("Blocked Accounts")
    }
}

struct MutedAccountsView: View {
    var body: some View {
        List {
            Text("No muted accounts")
                .foregroundColor(.gray)
        }
        .navigationTitle("Muted Accounts")
    }
}

struct RestrictedAccountsView: View {
    var body: some View {
        List {
            Text("No restricted accounts")
                .foregroundColor(.gray)
        }
        .navigationTitle("Restricted Accounts")
    }
}

struct LanguageView: View {
    var body: some View {
        List {
            Text("English")
            Text("Español")
            Text("Français")
            Text("Deutsch")
        }
        .navigationTitle("Language")
    }
}

struct AboutView: View {
    var body: some View {
        List {
            Section {
                HStack {
                    Text("Version")
                    Spacer()
                    Text("1.0.0")
                        .foregroundColor(.gray)
                }
                
                NavigationLink(destination: Text("Data Policy")) {
                    Text("Data Policy")
                }
                
                NavigationLink(destination: Text("Terms of Use")) {
                    Text("Terms of Use")
                }
                
                NavigationLink(destination: Text("Open Source Libraries")) {
                    Text("Open Source Libraries")
                }
            }
        }
        .navigationTitle("About")
    }
}

struct AccountStatusView: View {
    var body: some View {
        VStack(spacing: 24) {
            Image(systemName: "checkmark.shield.fill")
                .font(.system(size: 64))
                .foregroundColor(.green)
            
            Text("Your Account is in Good Standing")
                .font(.system(size: 20, weight: .bold))
            
            Text("You're following our Community Guidelines and Terms of Service.")
                .font(.system(size: 15))
                .foregroundColor(.gray)
                .multilineTextAlignment(.center)
                .padding(.horizontal, 32)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .navigationTitle("Account Status")
    }
}

// MARK: - View Model
class VignetteSettingsViewModel: ObservableObject {
    // Profile
    @Published var userName = "John Doe"
    @Published var username = "johndoe"
    @Published var userBio = "Living my best life! 🌟"
    @Published var userWebsite = "www.example.com"
    
    // Privacy
    @Published var isPrivateAccount = false
    @Published var commentSettings = "everyone"
    @Published var mentionSettings = "everyone"
    @Published var storySharing = "everyone"
    @Published var allowMessageRequests = true
    @Published var showActivityStatus = true
    @Published var readReceipts = true
    @Published var similarAccountSuggestions = true
    @Published var includeInRecommendations = true
    @Published var sensitiveContentControl = "limit"
    
    // Notifications
    @Published var notifyLikes = true
    @Published var notifyComments = true
    @Published var notifyFollowers = true
    @Published var notifyMessages = true
    @Published var notifyVideoViews = true
    @Published var notifyLiveVideos = true
    @Published var notifyIGTV = true
    @Published var emailFeedback = true
    @Published var emailReminders = false
    @Published var emailProduct = false
    
    // Data
    @Published var useLessData = false
    @Published var uploadQuality = "high"
    
    // Appearance
    @Published var darkModeEnabled = false
    
    func updateProfile(name: String, username: String, bio: String, website: String, completion: @escaping () -> Void) {
        // TODO: API call to update profile
        DispatchQueue.main.asyncAfter(deadline: .now() + 1) {
            self.userName = name
            self.username = username
            self.userBio = bio
            self.userWebsite = website
            completion()
        }
    }
    
    func savePrivacySettings() {
        // TODO: API call to save privacy settings
    }
    
    func saveNotificationSettings() {
        // TODO: API call to save notification settings
    }
    
    func saveDataSettings() {
        // TODO: API call to save data settings
    }
    
    func clearCache() {
        // TODO: Clear cache
    }
    
    func logout() {
        // Clear tokens
        KeychainManager.shared.deleteToken()
        // Navigate to login
    }
}
