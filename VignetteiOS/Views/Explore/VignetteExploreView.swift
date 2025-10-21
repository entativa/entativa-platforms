import SwiftUI

// MARK: - Vignette Explore View (Instagram-Style)
struct VignetteExploreView: View {
    @StateObject private var viewModel = ExploreViewModel()
    @State private var searchText = ""
    @State private var isSearching = false
    
    var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Search bar
                SearchBar(
                    text: $searchText,
                    isSearching: $isSearching,
                    placeholder: "Search"
                )
                .padding(.horizontal, 16)
                .padding(.vertical, 8)
                
                if isSearching {
                    // Search results
                    SearchResultsView(searchText: searchText, viewModel: viewModel)
                } else {
                    // Explore grid
                    ScrollView {
                        LazyVGrid(
                            columns: [
                                GridItem(.flexible(), spacing: 2),
                                GridItem(.flexible(), spacing: 2),
                                GridItem(.flexible(), spacing: 2)
                            ],
                            spacing: 2
                        ) {
                            ForEach(viewModel.explorePosts) { post in
                                ExplorePostCell(post: post)
                            }
                        }
                    }
                }
            }
            .navigationBarHidden(true)
        }
    }
}

// MARK: - Search Bar
struct SearchBar: View {
    @Binding var text: String
    @Binding var isSearching: Bool
    let placeholder: String
    
    var body: some View {
        HStack(spacing: 12) {
            HStack(spacing: 8) {
                Image(systemName: "magnifyingglass")
                    .foregroundColor(.gray)
                
                TextField(placeholder, text: $text, onEditingChanged: { editing in
                    isSearching = editing || !text.isEmpty
                })
                .autocapitalization(.none)
                
                if !text.isEmpty {
                    Button(action: {
                        text = ""
                    }) {
                        Image(systemName: "xmark.circle.fill")
                            .foregroundColor(.gray)
                    }
                }
            }
            .padding(.horizontal, 12)
            .padding(.vertical, 8)
            .background(Color.gray.opacity(0.1))
            .cornerRadius(10)
            
            if isSearching {
                Button("Cancel") {
                    text = ""
                    isSearching = false
                    UIApplication.shared.sendAction(#selector(UIResponder.resignFirstResponder), to: nil, from: nil, for: nil)
                }
                .foregroundColor(.primary)
            }
        }
    }
}

// MARK: - Search Results View
struct SearchResultsView: View {
    let searchText: String
    @ObservedObject var viewModel: ExploreViewModel
    @State private var selectedTab: SearchTab = .top
    
    var body: some View {
        VStack(spacing: 0) {
            // Tab selector
            ScrollView(.horizontal, showsIndicators: false) {
                HStack(spacing: 20) {
                    SearchTabButton(title: "Top", tab: .top, selectedTab: $selectedTab)
                    SearchTabButton(title: "Accounts", tab: .accounts, selectedTab: $selectedTab)
                    SearchTabButton(title: "Audio", tab: .audio, selectedTab: $selectedTab)
                    SearchTabButton(title: "Tags", tab: .tags, selectedTab: $selectedTab)
                    SearchTabButton(title: "Places", tab: .places, selectedTab: $selectedTab)
                }
                .padding(.horizontal, 16)
            }
            .padding(.vertical, 12)
            
            Divider()
            
            // Results
            ScrollView {
                LazyVStack(spacing: 0) {
                    if searchText.isEmpty {
                        // Recent searches
                        if !viewModel.recentSearches.isEmpty {
                            SectionHeader(title: "Recent")
                            
                            ForEach(viewModel.recentSearches) { search in
                                RecentSearchRow(search: search, onDelete: {
                                    viewModel.deleteRecentSearch(search)
                                })
                            }
                        }
                    } else {
                        // Search results by tab
                        switch selectedTab {
                        case .top:
                            TopSearchResults(searchText: searchText, viewModel: viewModel)
                        case .accounts:
                            AccountSearchResults(searchText: searchText, viewModel: viewModel)
                        case .audio:
                            AudioSearchResults(searchText: searchText)
                        case .tags:
                            TagSearchResults(searchText: searchText)
                        case .places:
                            PlaceSearchResults(searchText: searchText)
                        }
                    }
                }
            }
        }
    }
}

// MARK: - Search Tab Button
struct SearchTabButton: View {
    let title: String
    let tab: SearchTab
    @Binding var selectedTab: SearchTab
    
    var body: some View {
        Button(action: {
            selectedTab = tab
        }) {
            VStack(spacing: 8) {
                Text(title)
                    .font(.system(size: 16, weight: selectedTab == tab ? .semibold : .regular))
                    .foregroundColor(selectedTab == tab ? .primary : .gray)
                
                Rectangle()
                    .fill(selectedTab == tab ? Color.primary : Color.clear)
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
        .padding(.vertical, 12)
    }
}

// MARK: - Recent Search Row
struct RecentSearchRow: View {
    let search: RecentSearch
    let onDelete: () -> Void
    
    var body: some View {
        HStack(spacing: 12) {
            // Icon
            Circle()
                .fill(Color.gray.opacity(0.2))
                .frame(width: 44, height: 44)
                .overlay(
                    Image(systemName: search.type.icon)
                        .foregroundColor(.primary)
                )
            
            // Text
            VStack(alignment: .leading, spacing: 4) {
                Text(search.text)
                    .font(.system(size: 15, weight: .semibold))
                
                if let subtitle = search.subtitle {
                    Text(subtitle)
                        .font(.system(size: 13))
                        .foregroundColor(.gray)
                }
            }
            
            Spacer()
            
            // Delete button
            Button(action: onDelete) {
                Image(systemName: "xmark")
                    .font(.system(size: 14))
                    .foregroundColor(.gray)
            }
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 8)
    }
}

// MARK: - Top Search Results
struct TopSearchResults: View {
    let searchText: String
    @ObservedObject var viewModel: ExploreViewModel
    
    var body: some View {
        VStack(spacing: 0) {
            // Accounts section
            if !viewModel.searchAccounts.isEmpty {
                SectionHeader(title: "Accounts")
                
                ForEach(viewModel.searchAccounts.prefix(3)) { account in
                    AccountResultRow(account: account)
                }
            }
            
            // Tags section
            if !viewModel.searchTags.isEmpty {
                SectionHeader(title: "Tags")
                
                ForEach(viewModel.searchTags.prefix(3)) { tag in
                    TagResultRow(tag: tag)
                }
            }
            
            // Places section
            SectionHeader(title: "Places")
            ForEach(0..<3) { _ in
                PlaceResultRow()
            }
        }
    }
}

// MARK: - Account Search Results
struct AccountSearchResults: View {
    let searchText: String
    @ObservedObject var viewModel: ExploreViewModel
    
    var body: some View {
        ForEach(viewModel.searchAccounts) { account in
            AccountResultRow(account: account)
        }
    }
}

// MARK: - Account Result Row
struct AccountResultRow: View {
    let account: SearchAccount
    @State private var isFollowing = false
    
    var body: some View {
        HStack(spacing: 12) {
            // Profile picture
            Circle()
                .fill(Color.gray.opacity(0.2))
                .frame(width: 44, height: 44)
                .overlay(
                    Image(systemName: "person.fill")
                        .foregroundColor(.gray)
                )
            
            // Info
            VStack(alignment: .leading, spacing: 4) {
                Text(account.username)
                    .font(.system(size: 15, weight: .semibold))
                
                Text(account.fullName)
                    .font(.system(size: 13))
                    .foregroundColor(.gray)
                
                if let subtitle = account.subtitle {
                    Text(subtitle)
                        .font(.system(size: 13))
                        .foregroundColor(.gray)
                }
            }
            
            Spacer()
            
            // Follow button
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
        .padding(.horizontal, 16)
        .padding(.vertical, 8)
    }
}

// MARK: - Tag Result Row
struct TagResultRow: View {
    let tag: SearchTag
    
    var body: some View {
        HStack(spacing: 12) {
            // Icon
            RoundedRectangle(cornerRadius: 8)
                .fill(Color.gray.opacity(0.2))
                .frame(width: 44, height: 44)
                .overlay(
                    Text("#")
                        .font(.system(size: 24, weight: .bold))
                        .foregroundColor(.primary)
                )
            
            // Info
            VStack(alignment: .leading, spacing: 4) {
                Text("#\(tag.name)")
                    .font(.system(size: 15, weight: .semibold))
                
                Text("\(formatCount(tag.postsCount)) posts")
                    .font(.system(size: 13))
                    .foregroundColor(.gray)
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

// MARK: - Audio/Place Search Results
struct AudioSearchResults: View {
    let searchText: String
    
    var body: some View {
        ForEach(0..<10) { _ in
            AudioResultRow()
        }
    }
}

struct PlaceSearchResults: View {
    let searchText: String
    
    var body: some View {
        ForEach(0..<10) { _ in
            PlaceResultRow()
        }
    }
}

struct AudioResultRow: View {
    var body: some View {
        HStack(spacing: 12) {
            RoundedRectangle(cornerRadius: 8)
                .fill(Color.gray.opacity(0.2))
                .frame(width: 44, height: 44)
                .overlay(
                    Image(systemName: "music.note")
                        .foregroundColor(.primary)
                )
            
            VStack(alignment: .leading, spacing: 4) {
                Text("Audio Name")
                    .font(.system(size: 15, weight: .semibold))
                
                Text("Artist Name • 1.2M posts")
                    .font(.system(size: 13))
                    .foregroundColor(.gray)
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

struct PlaceResultRow: View {
    var body: some View {
        HStack(spacing: 12) {
            RoundedRectangle(cornerRadius: 8)
                .fill(Color.gray.opacity(0.2))
                .frame(width: 44, height: 44)
                .overlay(
                    Image(systemName: "mappin.circle.fill")
                        .foregroundColor(.red)
                )
            
            VStack(alignment: .leading, spacing: 4) {
                Text("Location Name")
                    .font(.system(size: 15, weight: .semibold))
                
                Text("City, Country")
                    .font(.system(size: 13))
                    .foregroundColor(.gray)
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

// MARK: - Explore Post Cell
struct ExplorePostCell: View {
    let post: ExplorePost
    
    var body: some View {
        GeometryReader { geometry in
            ZStack(alignment: .topTrailing) {
                Rectangle()
                    .fill(Color.gray.opacity(0.2))
                    .aspectRatio(1, contentMode: .fill)
                
                // Video indicator
                if post.isVideo {
                    Image(systemName: "play.fill")
                        .font(.system(size: 14))
                        .foregroundColor(.white)
                        .padding(4)
                        .background(Color.black.opacity(0.4))
                        .clipShape(Circle())
                        .padding(8)
                }
            }
        }
        .aspectRatio(1, contentMode: .fit)
    }
}

// MARK: - Search Tab Enum
enum SearchTab {
    case top
    case accounts
    case audio
    case tags
    case places
}

// MARK: - Models
struct RecentSearch: Identifiable {
    let id = UUID()
    let type: SearchType
    let text: String
    let subtitle: String?
    
    enum SearchType {
        case account
        case tag
        case place
        case audio
        
        var icon: String {
            switch self {
            case .account: return "person.fill"
            case .tag: return "number"
            case .place: return "mappin.circle.fill"
            case .audio: return "music.note"
            }
        }
    }
}

struct SearchAccount: Identifiable {
    let id = UUID()
    let username: String
    let fullName: String
    let subtitle: String?
    let isVerified: Bool
}

struct SearchTag: Identifiable {
    let id = UUID()
    let name: String
    let postsCount: Int
}

struct ExplorePost: Identifiable {
    let id = UUID()
    let imageUrl: String?
    let isVideo: Bool
    let likesCount: Int
}

// MARK: - View Model
class ExploreViewModel: ObservableObject {
    @Published var explorePosts: [ExplorePost] = []
    @Published var recentSearches: [RecentSearch] = []
    @Published var searchAccounts: [SearchAccount] = []
    @Published var searchTags: [SearchTag] = []
    
    init() {
        loadMockData()
    }
    
    func loadMockData() {
        // Explore posts
        explorePosts = (0..<30).map { index in
            ExplorePost(
                imageUrl: nil,
                isVideo: index % 5 == 0,
                likesCount: Int.random(in: 100...50000)
            )
        }
        
        // Recent searches
        recentSearches = [
            RecentSearch(type: .account, text: "sarah_jones", subtitle: "Sarah Jones • Following"),
            RecentSearch(type: .tag, text: "photography", subtitle: "1.2M posts"),
            RecentSearch(type: .place, text: "San Francisco", subtitle: "California, USA"),
            RecentSearch(type: .audio, text: "Trending Audio", subtitle: "456K posts")
        ]
        
        // Search accounts
        searchAccounts = [
            SearchAccount(username: "alex_creative", fullName: "Alex Creative", subtitle: "Followed by sarah + 3 others", isVerified: false),
            SearchAccount(username: "mike_photographer", fullName: "Mike Wilson", subtitle: "Photographer", isVerified: true),
            SearchAccount(username: "emma_designs", fullName: "Emma Davis", subtitle: "Designer • Followed by john", isVerified: false)
        ]
        
        // Search tags
        searchTags = [
            SearchTag(name: "photography", postsCount: 1200000),
            SearchTag(name: "travel", postsCount: 890000),
            SearchTag(name: "food", postsCount: 750000)
        ]
    }
    
    func deleteRecentSearch(_ search: RecentSearch) {
        recentSearches.removeAll { $0.id == search.id }
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
