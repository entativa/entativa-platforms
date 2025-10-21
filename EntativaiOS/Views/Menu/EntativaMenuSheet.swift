import SwiftUI

// MARK: - Entativa Menu Sheet (Facebook-Style Left Menu)
struct EntativaMenuSheet: View {
    @Environment(\.dismiss) var dismiss
    @StateObject private var viewModel = MenuViewModel()
    @State private var showSettings = false
    @State private var showLogoutConfirm = false
    
    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 0) {
                    // User Profile Section
                    UserProfileSection(viewModel: viewModel)
                        .padding(.vertical, 16)
                    
                    Divider()
                    
                    // Your Shortcuts
                    MenuSection(title: "Your shortcuts") {
                        MenuShortcutItem(
                            icon: "person.2.fill",
                            title: "Friends",
                            color: Color(hex: "007CFC"),
                            badge: viewModel.friendRequestsCount
                        )
                        
                        MenuShortcutItem(
                            icon: "clock.fill",
                            title: "Memories",
                            color: Color.blue
                        )
                        
                        MenuShortcutItem(
                            icon: "bookmark.fill",
                            title: "Saved",
                            color: Color.purple
                        )
                        
                        MenuShortcutItem(
                            icon: "person.3.fill",
                            title: "Groups",
                            color: Color(hex: "007CFC"),
                            badge: viewModel.groupNotifications
                        )
                        
                        MenuShortcutItem(
                            icon: "video.fill",
                            title: "Video",
                            color: Color.cyan
                        )
                        
                        MenuShortcutItem(
                            icon: "bag.fill",
                            title: "Marketplace",
                            color: Color.green
                        )
                        
                        MenuShortcutItem(
                            icon: "calendar",
                            title: "Events",
                            color: Color.red
                        )
                    }
                    
                    Divider()
                    
                    // Settings & Privacy
                    MenuSection(title: "Settings & privacy") {
                        NavigationLink(destination: SettingsMainView()) {
                            MenuListItem(
                                icon: "gearshape.fill",
                                title: "Settings",
                                color: Color.gray
                            )
                        }
                        
                        NavigationLink(destination: PrivacySettingsView()) {
                            MenuListItem(
                                icon: "lock.shield.fill",
                                title: "Privacy Center",
                                color: Color.blue
                            )
                        }
                        
                        NavigationLink(destination: ActivityLogView()) {
                            MenuListItem(
                                icon: "clock.arrow.circlepath",
                                title: "Activity log",
                                color: Color.orange
                            )
                        }
                    }
                    
                    Divider()
                    
                    // Help & Support
                    MenuSection(title: "Help & support") {
                        NavigationLink(destination: HelpCenterView()) {
                            MenuListItem(
                                icon: "questionmark.circle.fill",
                                title: "Help Center",
                                color: Color.purple
                            )
                        }
                        
                        NavigationLink(destination: ReportProblemView()) {
                            MenuListItem(
                                icon: "exclamationmark.triangle.fill",
                                title: "Report a problem",
                                color: Color.red
                            )
                        }
                    }
                    
                    Divider()
                    
                    // Also From Entativa
                    MenuSection(title: "Also from Entativa") {
                        MenuListItem(
                            icon: "photo.on.rectangle.angled",
                            title: "Vignette",
                            color: Color(hex: "519CAB"),
                            subtitle: "Photo & video sharing"
                        )
                    }
                    
                    Divider()
                    
                    // Logout
                    Button(action: { showLogoutConfirm = true }) {
                        HStack(spacing: 12) {
                            Image(systemName: "rectangle.portrait.and.arrow.right")
                                .font(.system(size: 20))
                                .foregroundColor(.red)
                                .frame(width: 36, height: 36)
                                .background(Color.red.opacity(0.1))
                                .clipShape(Circle())
                            
                            Text("Log out")
                                .font(.system(size: 16, weight: .semibold))
                                .foregroundColor(.red)
                            
                            Spacer()
                        }
                        .padding(.horizontal, 16)
                        .padding(.vertical, 12)
                    }
                    
                    // Legal
                    VStack(spacing: 8) {
                        HStack(spacing: 8) {
                            Text("Terms")
                                .foregroundColor(.gray)
                            Text("•")
                                .foregroundColor(.gray)
                            Text("Privacy Policy")
                                .foregroundColor(.gray)
                            Text("•")
                                .foregroundColor(.gray)
                            Text("Cookies")
                                .foregroundColor(.gray)
                        }
                        .font(.system(size: 12))
                        
                        Text("Entativa © 2025")
                            .font(.system(size: 12))
                            .foregroundColor(.gray)
                    }
                    .padding(.vertical, 16)
                }
            }
            .navigationTitle("Menu")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button(action: { dismiss() }) {
                        Image(systemName: "xmark.circle.fill")
                            .foregroundColor(.gray)
                            .font(.system(size: 24))
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

// MARK: - User Profile Section
struct UserProfileSection: View {
    @ObservedObject var viewModel: MenuViewModel
    
    var body: some View {
        HStack(spacing: 12) {
            // Profile picture
            Circle()
                .fill(Color.gray.opacity(0.2))
                .frame(width: 60, height: 60)
                .overlay(
                    Image(systemName: "person.fill")
                        .font(.system(size: 28))
                        .foregroundColor(.gray)
                )
            
            VStack(alignment: .leading, spacing: 4) {
                Text(viewModel.userName)
                    .font(.system(size: 20, weight: .bold))
                
                Text("@\(viewModel.username)")
                    .font(.system(size: 14))
                    .foregroundColor(.gray)
                
                NavigationLink(destination: Text("View Profile")) {
                    Text("See your profile")
                        .font(.system(size: 14, weight: .medium))
                        .foregroundColor(Color(hex: "007CFC"))
                }
            }
            
            Spacer()
        }
        .padding(.horizontal, 16)
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
            Text(title)
                .font(.system(size: 16, weight: .bold))
                .foregroundColor(.gray)
                .padding(.horizontal, 16)
                .padding(.vertical, 12)
            
            content
        }
    }
}

// MARK: - Menu Shortcut Item
struct MenuShortcutItem: View {
    let icon: String
    let title: String
    let color: Color
    var badge: Int = 0
    
    var body: some View {
        Button(action: {}) {
            HStack(spacing: 12) {
                Image(systemName: icon)
                    .font(.system(size: 20))
                    .foregroundColor(.white)
                    .frame(width: 36, height: 36)
                    .background(color)
                    .clipShape(Circle())
                
                Text(title)
                    .font(.system(size: 16, weight: .medium))
                    .foregroundColor(.primary)
                
                Spacer()
                
                if badge > 0 {
                    Text("\(badge)")
                        .font(.system(size: 12, weight: .bold))
                        .foregroundColor(.white)
                        .padding(.horizontal, 8)
                        .padding(.vertical, 4)
                        .background(Color.red)
                        .clipShape(Capsule())
                }
                
                Image(systemName: "chevron.right")
                    .font(.system(size: 14))
                    .foregroundColor(.gray)
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 8)
        }
    }
}

// MARK: - Menu List Item
struct MenuListItem: View {
    let icon: String
    let title: String
    let color: Color
    var subtitle: String? = nil
    
    var body: some View {
        HStack(spacing: 12) {
            Image(systemName: icon)
                .font(.system(size: 20))
                .foregroundColor(.white)
                .frame(width: 36, height: 36)
                .background(color)
                .clipShape(Circle())
            
            VStack(alignment: .leading, spacing: 2) {
                Text(title)
                    .font(.system(size: 16, weight: .medium))
                    .foregroundColor(.primary)
                
                if let subtitle = subtitle {
                    Text(subtitle)
                        .font(.system(size: 13))
                        .foregroundColor(.gray)
                }
            }
            
            Spacer()
            
            Image(systemName: "chevron.right")
                .font(.system(size: 14))
                .foregroundColor(.gray)
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 8)
    }
}

// MARK: - Settings Main View
struct SettingsMainView: View {
    @StateObject private var viewModel = SettingsViewModel()
    
    var body: some View {
        List {
            // Account
            Section {
                NavigationLink(destination: AccountSettingsView(viewModel: viewModel)) {
                    SettingsRow(icon: "person.circle.fill", title: "Account", color: .blue)
                }
                
                NavigationLink(destination: PrivacySettingsView()) {
                    SettingsRow(icon: "lock.fill", title: "Privacy & Security", color: .red)
                }
                
                NavigationLink(destination: NotificationSettingsView(viewModel: viewModel)) {
                    SettingsRow(icon: "bell.fill", title: "Notifications", color: .purple)
                }
            }
            
            // Content
            Section {
                NavigationLink(destination: Text("Language")) {
                    SettingsRow(icon: "globe", title: "Language", color: .green)
                }
                
                NavigationLink(destination: Text("Accessibility")) {
                    SettingsRow(icon: "accessibility", title: "Accessibility", color: .orange)
                }
                
                Toggle(isOn: $viewModel.darkModeEnabled) {
                    SettingsRow(icon: "moon.fill", title: "Dark Mode", color: .indigo)
                }
            }
            
            // Data
            Section {
                NavigationLink(destination: DataUsageView(viewModel: viewModel)) {
                    SettingsRow(icon: "externaldrive.fill", title: "Data Usage", color: .cyan)
                }
                
                NavigationLink(destination: Text("Storage")) {
                    SettingsRow(icon: "internaldrive.fill", title: "Storage", color: .gray)
                }
            }
            
            // About
            Section {
                NavigationLink(destination: Text("About")) {
                    SettingsRow(icon: "info.circle.fill", title: "About", color: .blue)
                }
                
                NavigationLink(destination: Text("Help")) {
                    SettingsRow(icon: "questionmark.circle.fill", title: "Help Center", color: .green)
                }
            }
            
            // Danger Zone
            Section {
                Button(action: { viewModel.showDeleteAccountAlert = true }) {
                    SettingsRow(icon: "trash.fill", title: "Delete Account", color: .red)
                }
            }
        }
        .navigationTitle("Settings")
        .navigationBarTitleDisplayMode(.large)
        .alert("Delete Account", isPresented: $viewModel.showDeleteAccountAlert) {
            Button("Cancel", role: .cancel) {}
            Button("Delete", role: .destructive) {
                viewModel.deleteAccount()
            }
        } message: {
            Text("Are you sure? This action cannot be undone.")
        }
    }
}

// MARK: - Settings Row
struct SettingsRow: View {
    let icon: String
    let title: String
    let color: Color
    
    var body: some View {
        HStack(spacing: 12) {
            Image(systemName: icon)
                .foregroundColor(color)
                .frame(width: 28)
            
            Text(title)
                .font(.system(size: 16))
        }
    }
}

// MARK: - Account Settings View
struct AccountSettingsView: View {
    @ObservedObject var viewModel: SettingsViewModel
    @State private var email = ""
    @State private var phone = ""
    @State private var bio = ""
    @State private var isSaving = false
    
    var body: some View {
        List {
            Section(header: Text("Personal Information")) {
                VStack(alignment: .leading, spacing: 8) {
                    Text("Email")
                        .font(.system(size: 14, weight: .medium))
                        .foregroundColor(.gray)
                    TextField("Email", text: $email)
                        .textContentType(.emailAddress)
                        .autocapitalization(.none)
                }
                .padding(.vertical, 4)
                
                VStack(alignment: .leading, spacing: 8) {
                    Text("Phone")
                        .font(.system(size: 14, weight: .medium))
                        .foregroundColor(.gray)
                    TextField("Phone number", text: $phone)
                        .textContentType(.telephoneNumber)
                        .keyboardType(.phonePad)
                }
                .padding(.vertical, 4)
                
                VStack(alignment: .leading, spacing: 8) {
                    Text("Bio")
                        .font(.system(size: 14, weight: .medium))
                        .foregroundColor(.gray)
                    TextEditor(text: $bio)
                        .frame(height: 100)
                }
                .padding(.vertical, 4)
            }
            
            Section(header: Text("Account Actions")) {
                NavigationLink(destination: ChangePasswordView()) {
                    Text("Change Password")
                }
                
                NavigationLink(destination: BlockedUsersView()) {
                    Text("Blocked Users")
                }
                
                NavigationLink(destination: Text("Download Your Data")) {
                    Text("Download Your Data")
                }
            }
            
            Section {
                Button(action: {
                    isSaving = true
                    viewModel.updateAccount(email: email, phone: phone, bio: bio) {
                        isSaving = false
                    }
                }) {
                    HStack {
                        Spacer()
                        if isSaving {
                            ProgressView()
                        } else {
                            Text("Save Changes")
                                .fontWeight(.semibold)
                        }
                        Spacer()
                    }
                }
                .disabled(isSaving)
            }
        }
        .navigationTitle("Account")
        .onAppear {
            email = viewModel.userEmail
            phone = viewModel.userPhone
            bio = viewModel.userBio
        }
    }
}

// MARK: - Privacy Settings View
struct PrivacySettingsView: View {
    @StateObject private var viewModel = PrivacyViewModel()
    
    var body: some View {
        List {
            Section(header: Text("Profile Privacy")) {
                Toggle("Private Account", isOn: $viewModel.isPrivateAccount)
                Toggle("Activity Status", isOn: $viewModel.showActivityStatus)
                Toggle("Read Receipts", isOn: $viewModel.readReceipts)
            }
            
            Section(header: Text("Content Privacy")) {
                Picker("Who can see your posts", selection: $viewModel.postsVisibility) {
                    Text("Everyone").tag("everyone")
                    Text("Friends").tag("friends")
                    Text("Only Me").tag("only_me")
                }
                
                Picker("Who can comment", selection: $viewModel.commentsAllowed) {
                    Text("Everyone").tag("everyone")
                    Text("Friends").tag("friends")
                    Text("No One").tag("no_one")
                }
            }
            
            Section(header: Text("Messaging")) {
                Toggle("Message Requests", isOn: $viewModel.allowMessageRequests)
                Toggle("Group Invitations", isOn: $viewModel.allowGroupInvites)
            }
            
            Section(header: Text("Data & History")) {
                NavigationLink(destination: Text("Clear Search History")) {
                    Text("Clear Search History")
                }
                
                NavigationLink(destination: Text("Clear Watch History")) {
                    Text("Clear Watch History")
                }
            }
        }
        .navigationTitle("Privacy & Security")
        .onChange(of: viewModel.isPrivateAccount) { _, _ in
            viewModel.savePrivacySettings()
        }
        .onChange(of: viewModel.showActivityStatus) { _, _ in
            viewModel.savePrivacySettings()
        }
    }
}

// MARK: - Notification Settings View
struct NotificationSettingsView: View {
    @ObservedObject var viewModel: SettingsViewModel
    
    var body: some View {
        List {
            Section(header: Text("Push Notifications")) {
                Toggle("Likes", isOn: $viewModel.notifyLikes)
                Toggle("Comments", isOn: $viewModel.notifyComments)
                Toggle("New Followers", isOn: $viewModel.notifyFollowers)
                Toggle("Messages", isOn: $viewModel.notifyMessages)
                Toggle("Friend Requests", isOn: $viewModel.notifyFriendRequests)
            }
            
            Section(header: Text("Email Notifications")) {
                Toggle("Weekly Summary", isOn: $viewModel.emailWeeklySummary)
                Toggle("Product Updates", isOn: $viewModel.emailProductUpdates)
                Toggle("Tips & Recommendations", isOn: $viewModel.emailTips)
            }
            
            Section(header: Text("In-App Notifications")) {
                Toggle("Sound", isOn: $viewModel.notificationSound)
                Toggle("Vibration", isOn: $viewModel.notificationVibration)
                Toggle("Badge Count", isOn: $viewModel.showBadgeCount)
            }
        }
        .navigationTitle("Notifications")
        .onChange(of: viewModel.notifyLikes) { _, _ in
            viewModel.saveNotificationSettings()
        }
    }
}

// MARK: - Data Usage View
struct DataUsageView: View {
    @ObservedObject var viewModel: SettingsViewModel
    
    var body: some View {
        List {
            Section(header: Text("Media Quality")) {
                Picker("Upload Quality", selection: $viewModel.uploadQuality) {
                    Text("High").tag("high")
                    Text("Medium").tag("medium")
                    Text("Low").tag("low")
                }
                
                Picker("Autoplay", selection: $viewModel.autoplaySettings) {
                    Text("Always").tag("always")
                    Text("Wi-Fi Only").tag("wifi")
                    Text("Never").tag("never")
                }
            }
            
            Section(header: Text("Data Saver")) {
                Toggle("Data Saver Mode", isOn: $viewModel.dataSaverMode)
                
                if viewModel.dataSaverMode {
                    Text("Reduces data usage by loading lower quality media")
                        .font(.system(size: 13))
                        .foregroundColor(.gray)
                }
            }
            
            Section(header: Text("Cache")) {
                HStack {
                    Text("Cache Size")
                    Spacer()
                    Text("142 MB")
                        .foregroundColor(.gray)
                }
                
                Button("Clear Cache") {
                    viewModel.clearCache()
                }
            }
        }
        .navigationTitle("Data Usage")
        .onChange(of: viewModel.dataSaverMode) { _, _ in
            viewModel.saveDataSettings()
        }
    }
}

// MARK: - Helper Views
struct HelpCenterView: View {
    var body: some View {
        List {
            Section("Getting Started") {
                NavigationLink(destination: Text("How to create a post")) {
                    Text("How to create a post")
                }
                NavigationLink(destination: Text("How to add friends")) {
                    Text("How to add friends")
                }
            }
            
            Section("Privacy & Safety") {
                NavigationLink(destination: Text("Managing your privacy")) {
                    Text("Managing your privacy")
                }
                NavigationLink(destination: Text("Blocking users")) {
                    Text("Blocking users")
                }
            }
        }
        .navigationTitle("Help Center")
    }
}

struct ReportProblemView: View {
    @State private var problemDescription = ""
    
    var body: some View {
        VStack {
            Text("Describe the problem")
                .font(.headline)
                .padding()
            
            TextEditor(text: $problemDescription)
                .frame(height: 200)
                .padding()
                .border(Color.gray.opacity(0.3))
            
            Button("Submit Report") {
                // Submit report
            }
            .buttonStyle(.borderedProminent)
            .padding()
            
            Spacer()
        }
        .navigationTitle("Report Problem")
    }
}

struct ActivityLogView: View {
    var body: some View {
        Text("Activity Log")
            .navigationTitle("Activity Log")
    }
}

struct ChangePasswordView: View {
    @State private var currentPassword = ""
    @State private var newPassword = ""
    @State private var confirmPassword = ""
    
    var body: some View {
        List {
            Section {
                SecureField("Current Password", text: $currentPassword)
                SecureField("New Password", text: $newPassword)
                SecureField("Confirm Password", text: $confirmPassword)
            }
            
            Section {
                Button("Change Password") {
                    // Change password
                }
            }
        }
        .navigationTitle("Change Password")
    }
}

struct BlockedUsersView: View {
    var body: some View {
        List {
            Text("No blocked users")
                .foregroundColor(.gray)
        }
        .navigationTitle("Blocked Users")
    }
}

// MARK: - View Models
class MenuViewModel: ObservableObject {
    @Published var userName = "John Doe"
    @Published var username = "johndoe"
    @Published var friendRequestsCount = 3
    @Published var groupNotifications = 5
    
    func logout() {
        // Clear tokens
        KeychainManager.shared.deleteToken()
        // Navigate to login
    }
}

class SettingsViewModel: ObservableObject {
    // Account
    @Published var userEmail = "user@example.com"
    @Published var userPhone = ""
    @Published var userBio = ""
    
    // Appearance
    @Published var darkModeEnabled = false
    
    // Notifications
    @Published var notifyLikes = true
    @Published var notifyComments = true
    @Published var notifyFollowers = true
    @Published var notifyMessages = true
    @Published var notifyFriendRequests = true
    @Published var emailWeeklySummary = true
    @Published var emailProductUpdates = false
    @Published var emailTips = false
    @Published var notificationSound = true
    @Published var notificationVibration = true
    @Published var showBadgeCount = true
    
    // Data
    @Published var uploadQuality = "high"
    @Published var autoplaySettings = "wifi"
    @Published var dataSaverMode = false
    
    // Alerts
    @Published var showDeleteAccountAlert = false
    
    func updateAccount(email: String, phone: String, bio: String, completion: @escaping () -> Void) {
        // TODO: API call
        DispatchQueue.main.asyncAfter(deadline: .now() + 1) {
            self.userEmail = email
            self.userPhone = phone
            self.userBio = bio
            completion()
        }
    }
    
    func saveNotificationSettings() {
        // TODO: API call
    }
    
    func saveDataSettings() {
        // TODO: API call
    }
    
    func clearCache() {
        // TODO: Clear cache
    }
    
    func deleteAccount() {
        // TODO: API call to delete account
    }
}

class PrivacyViewModel: ObservableObject {
    @Published var isPrivateAccount = false
    @Published var showActivityStatus = true
    @Published var readReceipts = true
    @Published var postsVisibility = "everyone"
    @Published var commentsAllowed = "everyone"
    @Published var allowMessageRequests = true
    @Published var allowGroupInvites = true
    
    func savePrivacySettings() {
        // TODO: API call
    }
}
