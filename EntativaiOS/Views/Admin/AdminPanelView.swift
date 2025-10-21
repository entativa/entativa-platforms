import SwiftUI

// MARK: - Admin Panel View (Founder Only - @neoqiss)
struct AdminPanelView: View {
    @Environment(\.dismiss) var dismiss
    @StateObject private var viewModel = AdminPanelViewModel()
    @State private var selectedSection: AdminSection = .dashboard
    
    var body: some View {
        NavigationView {
            HStack(spacing: 0) {
                // Side Navigation
                AdminSideNavigation(selectedSection: $selectedSection)
                    .frame(width: 80)
                    .background(Color.black)
                
                // Main Content
                VStack(spacing: 0) {
                    // Header
                    AdminHeaderView(
                        section: selectedSection,
                        onDismiss: { dismiss() }
                    )
                    
                    Divider()
                    
                    // Content Area
                    switch selectedSection {
                    case .dashboard:
                        AdminDashboardView(viewModel: viewModel)
                    case .users:
                        AdminUsersView(viewModel: viewModel)
                    case .content:
                        AdminContentView(viewModel: viewModel)
                    case .platform:
                        AdminPlatformView(viewModel: viewModel)
                    case .analytics:
                        AdminAnalyticsView(viewModel: viewModel)
                    case .developer:
                        AdminDeveloperView(viewModel: viewModel)
                    case .security:
                        AdminSecurityView(viewModel: viewModel)
                    case .audit:
                        AdminAuditView(viewModel: viewModel)
                    }
                }
            }
            .navigationBarHidden(true)
        }
    }
}

// MARK: - Admin Section Enum
enum AdminSection: String, CaseIterable {
    case dashboard = "Dashboard"
    case users = "Users"
    case content = "Content"
    case platform = "Platform"
    case analytics = "Analytics"
    case developer = "Developer"
    case security = "Security"
    case audit = "Audit"
    
    var icon: String {
        switch self {
        case .dashboard: return "chart.bar.fill"
        case .users: return "person.3.fill"
        case .content: return "doc.text.fill"
        case .platform: return "gearshape.2.fill"
        case .analytics: return "chart.line.uptrend.xyaxis"
        case .developer: return "hammer.fill"
        case .security: return "lock.shield.fill"
        case .audit: return "list.clipboard.fill"
        }
    }
    
    var color: Color {
        switch self {
        case .dashboard: return Color(hex: "007CFC")
        case .users: return Color(hex: "6F3EFB")
        case .content: return Color.green
        case .platform: return Color.orange
        case .analytics: return Color.cyan
        case .developer: return Color.pink
        case .security: return Color.red
        case .audit: return Color.purple
        }
    }
}

// MARK: - Side Navigation
struct AdminSideNavigation: View {
    @Binding var selectedSection: AdminSection
    
    var body: some View {
        VStack(spacing: 0) {
            // Crown icon for Founder
            VStack(spacing: 4) {
                Image(systemName: "crown.fill")
                    .font(.system(size: 24))
                    .foregroundColor(.yellow)
                
                Text("FOUNDER")
                    .font(.system(size: 8, weight: .bold))
                    .foregroundColor(.yellow)
            }
            .padding(.vertical, 20)
            
            Divider()
                .background(Color.white.opacity(0.2))
            
            // Navigation Items
            ScrollView {
                VStack(spacing: 8) {
                    ForEach(AdminSection.allCases, id: \.self) { section in
                        AdminNavButton(
                            section: section,
                            isSelected: selectedSection == section,
                            action: { selectedSection = section }
                        )
                    }
                }
                .padding(.vertical, 12)
            }
            
            Spacer()
        }
    }
}

// MARK: - Nav Button
struct AdminNavButton: View {
    let section: AdminSection
    let isSelected: Bool
    let action: () -> Void
    
    var body: some View {
        Button(action: action) {
            VStack(spacing: 6) {
                Image(systemName: section.icon)
                    .font(.system(size: 22))
                    .foregroundColor(isSelected ? section.color : .white.opacity(0.6))
                
                Text(section.rawValue)
                    .font(.system(size: 9, weight: isSelected ? .semibold : .regular))
                    .foregroundColor(isSelected ? section.color : .white.opacity(0.6))
                    .lineLimit(1)
                    .minimumScaleFactor(0.7)
            }
            .frame(width: 70, height: 70)
            .background(isSelected ? Color.white.opacity(0.15) : Color.clear)
            .cornerRadius(12)
        }
    }
}

// MARK: - Header View
struct AdminHeaderView: View {
    let section: AdminSection
    let onDismiss: () -> Void
    
    var body: some View {
        HStack {
            Image(systemName: section.icon)
                .font(.system(size: 20))
                .foregroundColor(section.color)
            
            Text(section.rawValue)
                .font(.system(size: 24, weight: .bold))
            
            Spacer()
            
            // Session indicator
            HStack(spacing: 8) {
                Circle()
                    .fill(Color.green)
                    .frame(width: 8, height: 8)
                
                Text("Admin Active")
                    .font(.system(size: 12, weight: .medium))
                    .foregroundColor(.green)
            }
            .padding(.horizontal, 12)
            .padding(.vertical, 6)
            .background(Color.green.opacity(0.1))
            .cornerRadius(12)
            
            Button(action: onDismiss) {
                Image(systemName: "xmark.circle.fill")
                    .font(.system(size: 28))
                    .foregroundColor(.gray)
            }
        }
        .padding(.horizontal, 20)
        .padding(.vertical, 16)
        .background(Color(UIColor.systemBackground))
    }
}

// MARK: - Dashboard View
struct AdminDashboardView: View {
    @ObservedObject var viewModel: AdminPanelViewModel
    
    var body: some View {
        ScrollView {
            VStack(spacing: 20) {
                // Live Metrics
                Text("Live Platform Metrics")
                    .font(.system(size: 20, weight: .bold))
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .padding(.horizontal)
                
                LazyVGrid(columns: [GridItem(.flexible()), GridItem(.flexible())], spacing: 16) {
                    MetricCard(title: "Active Users", value: "12,453", icon: "person.3.fill", color: .blue)
                    MetricCard(title: "Posts/Second", value: "47", icon: "doc.text.fill", color: .green)
                    MetricCard(title: "Server Load", value: "23%", icon: "cpu", color: .orange)
                    MetricCard(title: "Error Rate", value: "0.02%", icon: "exclamationmark.triangle", color: .red)
                }
                .padding(.horizontal)
                
                Divider()
                
                // Quick Actions
                Text("Quick Actions")
                    .font(.system(size: 20, weight: .bold))
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .padding(.horizontal)
                
                VStack(spacing: 12) {
                    QuickActionButton(
                        icon: "magnifyingglass",
                        title: "Search Users",
                        color: Color(hex: "007CFC"),
                        action: {}
                    )
                    
                    QuickActionButton(
                        icon: "trash.fill",
                        title: "Moderation Queue",
                        color: .red,
                        badge: 5,
                        action: {}
                    )
                    
                    QuickActionButton(
                        icon: "bell.badge.fill",
                        title: "Broadcast Notification",
                        color: .purple,
                        action: {}
                    )
                }
                .padding(.horizontal)
            }
            .padding(.vertical, 20)
        }
    }
}

// MARK: - Metric Card
struct MetricCard: View {
    let title: String
    let value: String
    let icon: String
    let color: Color
    
    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack {
                Image(systemName: icon)
                    .font(.system(size: 20))
                    .foregroundColor(color)
                
                Spacer()
            }
            
            Text(value)
                .font(.system(size: 32, weight: .bold))
            
            Text(title)
                .font(.system(size: 14))
                .foregroundColor(.gray)
        }
        .padding(16)
        .frame(maxWidth: .infinity, alignment: .leading)
        .background(Color(UIColor.secondarySystemBackground))
        .cornerRadius(12)
    }
}

// MARK: - Quick Action Button
struct QuickActionButton: View {
    let icon: String
    let title: String
    let color: Color
    var badge: Int = 0
    let action: () -> Void
    
    var body: some View {
        Button(action: action) {
            HStack(spacing: 16) {
                ZStack(alignment: .topTrailing) {
                    Image(systemName: icon)
                        .font(.system(size: 22))
                        .foregroundColor(.white)
                        .frame(width: 48, height: 48)
                        .background(color)
                        .cornerRadius(10)
                    
                    if badge > 0 {
                        Text("\(badge)")
                            .font(.system(size: 10, weight: .bold))
                            .foregroundColor(.white)
                            .padding(.horizontal, 6)
                            .padding(.vertical, 3)
                            .background(Color.red)
                            .cornerRadius(8)
                            .offset(x: 8, y: -8)
                    }
                }
                
                Text(title)
                    .font(.system(size: 16, weight: .semibold))
                    .foregroundColor(.primary)
                
                Spacer()
                
                Image(systemName: "chevron.right")
                    .foregroundColor(.gray)
            }
            .padding(16)
            .background(Color(UIColor.secondarySystemBackground))
            .cornerRadius(12)
        }
    }
}

// MARK: - Users View
struct AdminUsersView: View {
    @ObservedObject var viewModel: AdminPanelViewModel
    @State private var searchText = ""
    
    var body: some View {
        VStack(spacing: 0) {
            // Search bar
            HStack {
                Image(systemName: "magnifyingglass")
                    .foregroundColor(.gray)
                
                TextField("Search users by username, email, or ID", text: $searchText)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
                
                if !searchText.isEmpty {
                    Button("Search") {
                        viewModel.searchUsers(query: searchText)
                    }
                    .buttonStyle(.borderedProminent)
                }
            }
            .padding()
            
            Divider()
            
            // Results
            if viewModel.searchResults.isEmpty {
                VStack(spacing: 16) {
                    Image(systemName: "person.3")
                        .font(.system(size: 48))
                        .foregroundColor(.gray)
                    
                    Text("Search for users to manage")
                        .font(.system(size: 16))
                        .foregroundColor(.gray)
                }
                .frame(maxHeight: .infinity)
            } else {
                ScrollView {
                    LazyVStack(spacing: 0) {
                        ForEach(viewModel.searchResults) { user in
                            AdminUserRow(user: user)
                        }
                    }
                }
            }
        }
    }
}

// MARK: - Admin User Row
struct AdminUserRow: View {
    let user: AdminUserResult
    @State private var showActions = false
    
    var body: some View {
        HStack(spacing: 12) {
            // Avatar
            Circle()
                .fill(Color.gray.opacity(0.2))
                .frame(width: 48, height: 48)
            
            // User info
            VStack(alignment: .leading, spacing: 4) {
                HStack(spacing: 6) {
                    Text(user.username)
                        .font(.system(size: 16, weight: .semibold))
                    
                    // Founder badge
                    if user.isFounder {
                        Image("neoqissCheck")
                            .resizable()
                            .frame(width: 16, height: 16)
                    }
                }
                
                Text(user.email)
                    .font(.system(size: 14))
                    .foregroundColor(.gray)
                
                HStack(spacing: 8) {
                    if user.isBanned {
                        StatusPill(text: "BANNED", color: .red)
                    }
                    if user.isShadowbanned {
                        StatusPill(text: "SHADOWBANNED", color: .orange)
                    }
                    if user.isSuspended {
                        StatusPill(text: "SUSPENDED", color: .yellow)
                    }
                }
            }
            
            Spacer()
            
            // Actions button
            Button(action: { showActions = true }) {
                Image(systemName: "ellipsis.circle.fill")
                    .font(.system(size: 24))
                    .foregroundColor(Color(hex: "007CFC"))
            }
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 12)
        .confirmationDialog("User Actions", isPresented: $showActions) {
            Button("View Profile") {}
            Button("Ban User", role: .destructive) {}
            Button("Shadowban") {}
            Button("Force Logout") {}
            Button("Reset Password") {}
            if !user.isFounder {
                Button("Impersonate User") {}
            }
            Button("Cancel", role: .cancel) {}
        }
    }
}

// MARK: - Status Pill
struct StatusPill: View {
    let text: String
    let color: Color
    
    var body: some View {
        Text(text)
            .font(.system(size: 10, weight: .bold))
            .foregroundColor(.white)
            .padding(.horizontal, 8)
            .padding(.vertical, 4)
            .background(color)
            .cornerRadius(8)
    }
}

// MARK: - Platform View
struct AdminPlatformView: View {
    @ObservedObject var viewModel: AdminPanelViewModel
    
    var body: some View {
        ScrollView {
            VStack(spacing: 24) {
                // Feature Flags
                FeatureFlagsSection(viewModel: viewModel)
                
                Divider()
                
                // Kill Switches
                KillSwitchesSection(viewModel: viewModel)
                
                Divider()
                
                // Maintenance Mode
                MaintenanceModeSection(viewModel: viewModel)
            }
            .padding(20)
        }
    }
}

// MARK: - Feature Flags Section
struct FeatureFlagsSection: View {
    @ObservedObject var viewModel: AdminPanelViewModel
    
    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text("Feature Flags")
                .font(.system(size: 20, weight: .bold))
            
            VStack(spacing: 12) {
                FeatureFlagRow(name: "Posting", isEnabled: viewModel.featurePosting)
                FeatureFlagRow(name: "Commenting", isEnabled: viewModel.featureCommenting)
                FeatureFlagRow(name: "Messaging", isEnabled: viewModel.featureMessaging)
                FeatureFlagRow(name: "Takes", isEnabled: viewModel.featureTakes)
                FeatureFlagRow(name: "Stories", isEnabled: viewModel.featureStories)
                FeatureFlagRow(name: "Live Streaming", isEnabled: viewModel.featureLiveStreaming)
            }
        }
    }
}

// MARK: - Feature Flag Row
struct FeatureFlagRow: View {
    let name: String
    let isEnabled: Bool
    
    var body: some View {
        HStack {
            Text(name)
                .font(.system(size: 16, weight: .medium))
            
            Spacer()
            
            Toggle("", isOn: .constant(isEnabled))
                .labelsHidden()
        }
        .padding(12)
        .background(Color(UIColor.secondarySystemBackground))
        .cornerRadius(10)
    }
}

// MARK: - Kill Switches Section
struct KillSwitchesSection: View {
    @ObservedObject var viewModel: AdminPanelViewModel
    @State private var showKillSwitchAlert = false
    @State private var selectedSwitch = ""
    
    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            HStack {
                Image(systemName: "exclamationmark.triangle.fill")
                    .foregroundColor(.red)
                Text("Emergency Kill Switches")
                    .font(.system(size: 20, weight: .bold))
            }
            
            VStack(spacing: 12) {
                KillSwitchButton(
                    title: "Disable All Posting",
                    isActive: false,
                    action: {
                        selectedSwitch = "posting"
                        showKillSwitchAlert = true
                    }
                )
                
                KillSwitchButton(
                    title: "Disable All Commenting",
                    isActive: false,
                    action: {
                        selectedSwitch = "commenting"
                        showKillSwitchAlert = true
                    }
                )
                
                KillSwitchButton(
                    title: "Disable Messaging",
                    isActive: false,
                    action: {
                        selectedSwitch = "messaging"
                        showKillSwitchAlert = true
                    }
                )
                
                KillSwitchButton(
                    title: "Emergency Lockdown",
                    isActive: false,
                    action: {
                        selectedSwitch = "lockdown"
                        showKillSwitchAlert = true
                    }
                )
            }
        }
        .alert("Activate Kill Switch?", isPresented: $showKillSwitchAlert) {
            Button("Cancel", role: .cancel) {}
            Button("Activate", role: .destructive) {
                viewModel.activateKillSwitch(selectedSwitch)
            }
        } message: {
            Text("This will immediately disable \(selectedSwitch) for all users. This action is logged.")
        }
    }
}

// MARK: - Kill Switch Button
struct KillSwitchButton: View {
    let title: String
    let isActive: Bool
    let action: () -> Void
    
    var body: some View {
        Button(action: action) {
            HStack {
                Circle()
                    .fill(isActive ? Color.red : Color.gray)
                    .frame(width: 12, height: 12)
                
                Text(title)
                    .font(.system(size: 16, weight: .medium))
                    .foregroundColor(isActive ? .red : .primary)
                
                Spacer()
                
                Text(isActive ? "ACTIVE" : "INACTIVE")
                    .font(.system(size: 12, weight: .bold))
                    .foregroundColor(.white)
                    .padding(.horizontal, 12)
                    .padding(.vertical, 6)
                    .background(isActive ? Color.red : Color.gray)
                    .cornerRadius(12)
            }
            .padding(16)
            .background(Color(UIColor.secondarySystemBackground))
            .cornerRadius(12)
            .overlay(
                RoundedRectangle(cornerRadius: 12)
                    .stroke(isActive ? Color.red : Color.clear, lineWidth: 2)
            )
        }
    }
}

// MARK: - Maintenance Mode Section
struct MaintenanceModeSection: View {
    @ObservedObject var viewModel: AdminPanelViewModel
    @State private var showMaintenanceAlert = false
    
    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text("Maintenance Mode")
                .font(.system(size: 20, weight: .bold))
            
            Button(action: { showMaintenanceAlert = true }) {
                HStack {
                    Image(systemName: "wrench.and.screwdriver.fill")
                        .font(.system(size: 20))
                    
                    Text(viewModel.isMaintenanceMode ? "Disable Maintenance Mode" : "Enable Maintenance Mode")
                        .font(.system(size: 16, weight: .semibold))
                    
                    Spacer()
                }
                .foregroundColor(.white)
                .padding(16)
                .background(viewModel.isMaintenanceMode ? Color.green : Color.orange)
                .cornerRadius(12)
            }
        }
        .alert("Maintenance Mode", isPresented: $showMaintenanceAlert) {
            Button("Cancel", role: .cancel) {}
            Button(viewModel.isMaintenanceMode ? "Disable" : "Enable") {
                viewModel.toggleMaintenanceMode()
            }
        } message: {
            Text(viewModel.isMaintenanceMode ? "Users will be able to access the platform again." : "Users won't be able to access the platform during maintenance.")
        }
    }
}

// MARK: - Content, Analytics, Developer, Security, Audit Views (Placeholders)
struct AdminContentView: View {
    @ObservedObject var viewModel: AdminPanelViewModel
    var body: some View {
        Text("Content Moderation")
            .font(.largeTitle)
    }
}

struct AdminAnalyticsView: View {
    @ObservedObject var viewModel: AdminPanelViewModel
    var body: some View {
        Text("Analytics Dashboard")
            .font(.largeTitle)
    }
}

struct AdminDeveloperView: View {
    @ObservedObject var viewModel: AdminPanelViewModel
    var body: some View {
        Text("Developer Tools")
            .font(.largeTitle)
    }
}

struct AdminSecurityView: View {
    @ObservedObject var viewModel: AdminPanelViewModel
    var body: some View {
        Text("Security & Audit")
            .font(.largeTitle)
    }
}

struct AdminAuditView: View {
    @ObservedObject var viewModel: AdminPanelViewModel
    var body: some View {
        Text("Audit Logs")
            .font(.largeTitle)
    }
}

// MARK: - Models
struct AdminUserResult: Identifiable {
    let id: String
    let username: String
    let email: String
    let isFounder: Bool
    let isBanned: Bool
    let isShadowbanned: Bool
    let isSuspended: Bool
}

// MARK: - View Model
class AdminPanelViewModel: ObservableObject {
    @Published var searchResults: [AdminUserResult] = []
    
    // Feature flags
    @Published var featurePosting = true
    @Published var featureCommenting = true
    @Published var featureMessaging = true
    @Published var featureTakes = true
    @Published var featureStories = true
    @Published var featureLiveStreaming = false
    
    // Platform state
    @Published var isMaintenanceMode = false
    
    func searchUsers(query: String) {
        // TODO: API call
        searchResults = [
            AdminUserResult(
                id: "1",
                username: "testuser",
                email: "test@example.com",
                isFounder: false,
                isBanned: false,
                isShadowbanned: false,
                isSuspended: false
            )
        ]
    }
    
    func activateKillSwitch(_ switchName: String) {
        // TODO: API call
    }
    
    func toggleMaintenanceMode() {
        isMaintenanceMode.toggle()
        // TODO: API call
    }
}
