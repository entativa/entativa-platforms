import SwiftUI

// MARK: - Entativa Search View (Facebook-Style)
struct EntativaSearchView: View {
    @StateObject private var viewModel = EntativaSearchViewModel()
    @State private var searchText = ""
    @State private var isSearching = false
    
    var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Search bar
                HStack(spacing: 12) {
                    HStack(spacing: 8) {
                        Image(systemName: "magnifyingglass")
                            .foregroundColor(.gray)
                        
                        TextField("Search Entativa", text: $searchText, onEditingChanged: { editing in
                            isSearching = editing || !text.isEmpty
                        })
                        .autocapitalization(.none)
                        
                        if !searchText.isEmpty {
                            Button(action: {
                                searchText = ""
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
                            searchText = ""
                            isSearching = false
                            UIApplication.shared.sendAction(#selector(UIResponder.resignFirstResponder), to: nil, from: nil, for: nil)
                        }
                        .foregroundColor(.primary)
                    }
                }
                .padding(.horizontal, 16)
                .padding(.vertical, 8)
                
                Divider()
                
                if isSearching && !searchText.isEmpty {
                    // Search results with filters
                    SearchResultsWithFiltersView(searchText: searchText, viewModel: viewModel)
                } else {
                    // Recent searches and suggestions
                    RecentAndSuggestionsView(viewModel: viewModel)
                }
            }
            .navigationTitle("Search")
            .navigationBarTitleDisplayMode(.inline)
        }
    }
}

// MARK: - Recent and Suggestions View
struct RecentAndSuggestionsView: View {
    @ObservedObject var viewModel: EntativaSearchViewModel
    
    var body: some View {
        ScrollView {
            VStack(spacing: 0) {
                // Recent searches
                if !viewModel.recentSearches.isEmpty {
                    SearchSectionHeader(title: "Recent", showClear: true, onClear: {
                        viewModel.clearRecentSearches()
                    })
                    
                    ForEach(viewModel.recentSearches) { search in
                        EntativaRecentSearchRow(search: search, onDelete: {
                            viewModel.deleteRecentSearch(search)
                        })
                    }
                    
                    Divider()
                        .padding(.vertical, 8)
                }
                
                // Suggested searches
                SearchSectionHeader(title: "Suggested Searches", showClear: false, onClear: {})
                
                ForEach(viewModel.suggestedSearches) { suggestion in
                    SuggestedSearchRow(suggestion: suggestion)
                }
            }
        }
        .background(Color(UIColor.systemBackground))
    }
}

// MARK: - Search Results with Filters View
struct SearchResultsWithFiltersView: View {
    let searchText: String
    @ObservedObject var viewModel: EntativaSearchViewModel
    @State private var selectedFilter: SearchFilter = .all
    
    var body: some View {
        VStack(spacing: 0) {
            // Filter tabs
            ScrollView(.horizontal, showsIndicators: false) {
                HStack(spacing: 12) {
                    FilterChip(title: "All", filter: .all, selectedFilter: $selectedFilter)
                    FilterChip(title: "People", filter: .people, selectedFilter: $selectedFilter)
                    FilterChip(title: "Posts", filter: .posts, selectedFilter: $selectedFilter)
                    FilterChip(title: "Photos", filter: .photos, selectedFilter: $selectedFilter)
                    FilterChip(title: "Videos", filter: .videos, selectedFilter: $selectedFilter)
                    FilterChip(title: "Pages", filter: .pages, selectedFilter: $selectedFilter)
                    FilterChip(title: "Groups", filter: .groups, selectedFilter: $selectedFilter)
                    FilterChip(title: "Events", filter: .events, selectedFilter: $selectedFilter)
                }
                .padding(.horizontal, 16)
            }
            .padding(.vertical, 12)
            
            Divider()
            
            // Results
            ScrollView {
                LazyVStack(spacing: 0) {
                    switch selectedFilter {
                    case .all:
                        AllSearchResults(searchText: searchText, viewModel: viewModel)
                    case .people:
                        PeopleSearchResults(searchText: searchText, viewModel: viewModel)
                    case .posts:
                        PostsSearchResults(searchText: searchText)
                    case .photos:
                        PhotosSearchResults(searchText: searchText)
                    case .videos:
                        VideosSearchResults(searchText: searchText)
                    case .pages:
                        PagesSearchResults(searchText: searchText)
                    case .groups:
                        GroupsSearchResults(searchText: searchText)
                    case .events:
                        EventsSearchResults(searchText: searchText)
                    }
                }
            }
        }
    }
}

// MARK: - Search Section Header
struct SearchSectionHeader: View {
    let title: String
    let showClear: Bool
    let onClear: () -> Void
    
    var body: some View {
        HStack {
            Text(title)
                .font(.system(size: 17, weight: .semibold))
                .foregroundColor(.primary)
            
            Spacer()
            
            if showClear {
                Button("Clear all") {
                    onClear()
                }
                .font(.system(size: 15))
                .foregroundColor(Color(hex: "007CFC"))
            }
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 12)
    }
}

// MARK: - Recent Search Row
struct EntativaRecentSearchRow: View {
    let search: EntativaRecentSearch
    let onDelete: () -> Void
    
    var body: some View {
        HStack(spacing: 12) {
            // Icon
            Circle()
                .fill(Color.gray.opacity(0.2))
                .frame(width: 40, height: 40)
                .overlay(
                    Image(systemName: search.icon)
                        .font(.system(size: 18))
                        .foregroundColor(.gray)
                )
            
            // Text
            VStack(alignment: .leading, spacing: 4) {
                Text(search.text)
                    .font(.system(size: 15))
                
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

// MARK: - Suggested Search Row
struct SuggestedSearchRow: View {
    let suggestion: SuggestedSearch
    
    var body: some View {
        HStack(spacing: 12) {
            // Icon
            Circle()
                .fill(Color(hex: "007CFC").opacity(0.1))
                .frame(width: 40, height: 40)
                .overlay(
                    Image(systemName: suggestion.icon)
                        .font(.system(size: 18))
                        .foregroundColor(Color(hex: "007CFC"))
                )
            
            // Text
            Text(suggestion.text)
                .font(.system(size: 15))
            
            Spacer()
            
            Image(systemName: "chevron.right")
                .font(.system(size: 14))
                .foregroundColor(.gray)
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 8)
    }
}

// MARK: - Filter Chip
struct FilterChip: View {
    let title: String
    let filter: SearchFilter
    @Binding var selectedFilter: SearchFilter
    
    var body: some View {
        Button(action: {
            selectedFilter = filter
        }) {
            Text(title)
                .font(.system(size: 14, weight: selectedFilter == filter ? .semibold : .regular))
                .foregroundColor(selectedFilter == filter ? .white : .primary)
                .padding(.horizontal, 16)
                .padding(.vertical, 8)
                .background(selectedFilter == filter ? Color(hex: "007CFC") : Color.gray.opacity(0.1))
                .cornerRadius(20)
        }
    }
}

// MARK: - All Search Results
struct AllSearchResults: View {
    let searchText: String
    @ObservedObject var viewModel: EntativaSearchViewModel
    
    var body: some View {
        VStack(spacing: 0) {
            // People section
            if !viewModel.searchPeople.isEmpty {
                SearchSectionHeader(title: "People", showClear: false, onClear: {})
                
                ForEach(viewModel.searchPeople.prefix(3)) { person in
                    PersonResultRow(person: person)
                }
                
                Divider().padding(.vertical, 8)
            }
            
            // Posts section
            SearchSectionHeader(title: "Posts", showClear: false, onClear: {})
            ForEach(0..<3) { _ in
                PostResultRow()
            }
            
            Divider().padding(.vertical, 8)
            
            // Pages section
            SearchSectionHeader(title: "Pages", showClear: false, onClear: {})
            ForEach(0..<2) { _ in
                PageResultRow()
            }
        }
    }
}

// MARK: - People Search Results
struct PeopleSearchResults: View {
    let searchText: String
    @ObservedObject var viewModel: EntativaSearchViewModel
    
    var body: some View {
        ForEach(viewModel.searchPeople) { person in
            PersonResultRow(person: person)
        }
    }
}

// MARK: - Person Result Row
struct PersonResultRow: View {
    let person: SearchPerson
    @State private var isFriend = false
    
    var body: some View {
        HStack(spacing: 12) {
            // Profile picture
            Circle()
                .fill(Color.gray.opacity(0.2))
                .frame(width: 56, height: 56)
                .overlay(
                    Image(systemName: "person.fill")
                        .foregroundColor(.gray)
                )
            
            // Info
            VStack(alignment: .leading, spacing: 4) {
                Text(person.name)
                    .font(.system(size: 15, weight: .semibold))
                
                if let subtitle = person.subtitle {
                    Text(subtitle)
                        .font(.system(size: 13))
                        .foregroundColor(.gray)
                }
            }
            
            Spacer()
            
            // Action button
            Button(action: {
                isFriend.toggle()
            }) {
                Image(systemName: isFriend ? "checkmark" : "person.badge.plus")
                    .font(.system(size: 16))
                    .foregroundColor(Color(hex: "007CFC"))
                    .frame(width: 36, height: 36)
                    .background(Color(hex: "007CFC").opacity(0.1))
                    .clipShape(Circle())
            }
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 8)
    }
}

// MARK: - Other Search Results
struct PostsSearchResults: View {
    let searchText: String
    
    var body: some View {
        ForEach(0..<10) { _ in
            PostResultRow()
        }
    }
}

struct PostResultRow: View {
    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            HStack(spacing: 8) {
                Circle()
                    .fill(Color.gray.opacity(0.2))
                    .frame(width: 32, height: 32)
                
                VStack(alignment: .leading, spacing: 2) {
                    Text("User Name")
                        .font(.system(size: 13, weight: .semibold))
                    Text("2 hours ago")
                        .font(.system(size: 11))
                        .foregroundColor(.gray)
                }
            }
            
            Text("This is a sample post that matches the search query...")
                .font(.system(size: 15))
                .lineLimit(3)
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 12)
        .background(Color(UIColor.systemBackground))
    }
}

struct PhotosSearchResults: View {
    let searchText: String
    
    var body: some View {
        LazyVGrid(
            columns: [
                GridItem(.flexible(), spacing: 2),
                GridItem(.flexible(), spacing: 2),
                GridItem(.flexible(), spacing: 2)
            ],
            spacing: 2
        ) {
            ForEach(0..<20) { _ in
                Rectangle()
                    .fill(Color.gray.opacity(0.2))
                    .aspectRatio(1, contentMode: .fit)
            }
        }
    }
}

struct VideosSearchResults: View {
    let searchText: String
    
    var body: some View {
        LazyVGrid(
            columns: [
                GridItem(.flexible(), spacing: 2),
                GridItem(.flexible(), spacing: 2)
            ],
            spacing: 2
        ) {
            ForEach(0..<10) { _ in
                Rectangle()
                    .fill(Color.gray.opacity(0.2))
                    .aspectRatio(16/9, contentMode: .fit)
                    .overlay(
                        Image(systemName: "play.fill")
                            .foregroundColor(.white)
                    )
            }
        }
    }
}

struct PagesSearchResults: View {
    let searchText: String
    
    var body: some View {
        ForEach(0..<10) { _ in
            PageResultRow()
        }
    }
}

struct PageResultRow: View {
    var body: some View {
        HStack(spacing: 12) {
            RoundedRectangle(cornerRadius: 8)
                .fill(Color.gray.opacity(0.2))
                .frame(width: 56, height: 56)
            
            VStack(alignment: .leading, spacing: 4) {
                Text("Page Name")
                    .font(.system(size: 15, weight: .semibold))
                
                Text("Category • 1.2K likes")
                    .font(.system(size: 13))
                    .foregroundColor(.gray)
            }
            
            Spacer()
            
            Button(action: {}) {
                Image(systemName: "heart")
                    .font(.system(size: 16))
                    .foregroundColor(Color(hex: "007CFC"))
                    .frame(width: 36, height: 36)
                    .background(Color(hex: "007CFC").opacity(0.1))
                    .clipShape(Circle())
            }
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 8)
    }
}

struct GroupsSearchResults: View {
    let searchText: String
    
    var body: some View {
        ForEach(0..<10) { _ in
            GroupResultRow()
        }
    }
}

struct GroupResultRow: View {
    var body: some View {
        HStack(spacing: 12) {
            RoundedRectangle(cornerRadius: 8)
                .fill(Color.gray.opacity(0.2))
                .frame(width: 56, height: 56)
            
            VStack(alignment: .leading, spacing: 4) {
                Text("Group Name")
                    .font(.system(size: 15, weight: .semibold))
                
                Text("Public group • 5.6K members")
                    .font(.system(size: 13))
                    .foregroundColor(.gray)
            }
            
            Spacer()
            
            Button(action: {}) {
                Image(systemName: "plus")
                    .font(.system(size: 16))
                    .foregroundColor(Color(hex: "007CFC"))
                    .frame(width: 36, height: 36)
                    .background(Color(hex: "007CFC").opacity(0.1))
                    .clipShape(Circle())
            }
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 8)
    }
}

struct EventsSearchResults: View {
    let searchText: String
    
    var body: some View {
        ForEach(0..<10) { _ in
            EventResultRow()
        }
    }
}

struct EventResultRow: View {
    var body: some View {
        HStack(spacing: 12) {
            RoundedRectangle(cornerRadius: 8)
                .fill(Color.gray.opacity(0.2))
                .frame(width: 56, height: 56)
            
            VStack(alignment: .leading, spacing: 4) {
                Text("Event Name")
                    .font(.system(size: 15, weight: .semibold))
                
                Text("Tomorrow at 7:00 PM")
                    .font(.system(size: 13))
                    .foregroundColor(.gray)
            }
            
            Spacer()
            
            Button(action: {}) {
                Image(systemName: "star")
                    .font(.system(size: 16))
                    .foregroundColor(Color(hex: "007CFC"))
                    .frame(width: 36, height: 36)
                    .background(Color(hex: "007CFC").opacity(0.1))
                    .clipShape(Circle())
            }
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 8)
    }
}

// MARK: - Search Filter Enum
enum SearchFilter {
    case all
    case people
    case posts
    case photos
    case videos
    case pages
    case groups
    case events
}

// MARK: - Models
struct EntativaRecentSearch: Identifiable {
    let id = UUID()
    let text: String
    let subtitle: String?
    let icon: String
}

struct SuggestedSearch: Identifiable {
    let id = UUID()
    let text: String
    let icon: String
}

struct SearchPerson: Identifiable {
    let id = UUID()
    let name: String
    let subtitle: String?
}

// MARK: - View Model
class EntativaSearchViewModel: ObservableObject {
    @Published var recentSearches: [EntativaRecentSearch] = []
    @Published var suggestedSearches: [SuggestedSearch] = []
    @Published var searchPeople: [SearchPerson] = []
    
    init() {
        loadMockData()
    }
    
    func loadMockData() {
        // Recent searches
        recentSearches = [
            EntativaRecentSearch(text: "Sarah Johnson", subtitle: "Friend", icon: "person.fill"),
            EntativaRecentSearch(text: "Photography Group", subtitle: "Group • 2.3K members", icon: "person.2.fill"),
            EntativaRecentSearch(text: "Summer Festival", subtitle: "Event", icon: "calendar"),
            EntativaRecentSearch(text: "Tech News", subtitle: "Page • 45K likes", icon: "flag.fill")
        ]
        
        // Suggested searches
        suggestedSearches = [
            SuggestedSearch(text: "Friends", icon: "person.2.fill"),
            SuggestedSearch(text: "Pages you may like", icon: "flag.fill"),
            SuggestedSearch(text: "Groups", icon: "person.3.fill"),
            SuggestedSearch(text: "Events near you", icon: "mappin.circle.fill"),
            SuggestedSearch(text: "Marketplace", icon: "cart.fill")
        ]
        
        // Search people
        searchPeople = [
            SearchPerson(name: "Alex Creative", subtitle: "3 mutual friends"),
            SearchPerson(name: "Mike Wilson", subtitle: "Photographer • Followed by Sarah"),
            SearchPerson(name: "Emma Davis", subtitle: "Designer")
        ]
    }
    
    func clearRecentSearches() {
        recentSearches.removeAll()
    }
    
    func deleteRecentSearch(_ search: EntativaRecentSearch) {
        recentSearches.removeAll { $0.id == search.id }
    }
}
