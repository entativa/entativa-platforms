package com.entativa.ui.search

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewmodel.compose.viewModel
import com.entativa.R
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow

// MARK: - Entativa Search Screen (Facebook-Style)
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun EntativaSearchScreen(
    viewModel: EntativaSearchViewModel = viewModel()
) {
    var searchText by remember { mutableStateOf("") }
    val isSearching = searchText.isNotEmpty()
    val selectedFilter by viewModel.selectedFilter.collectAsState()
    val recentSearches by viewModel.recentSearches.collectAsState()

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    OutlinedTextField(
                        value = searchText,
                        onValueChange = { searchText = it },
                        placeholder = { Text("Search Entativa") },
                        leadingIcon = {
                            Icon(painterResource(R.drawable.ic_search), null)
                        },
                        trailingIcon = {
                            if (searchText.isNotEmpty()) {
                                IconButton(onClick = { searchText = "" }) {
                                    Icon(painterResource(R.drawable.ic_close), null)
                                }
                            }
                        },
                        modifier = Modifier.fillMaxWidth(),
                        singleLine = true,
                        shape = RoundedCornerShape(24.dp),
                        colors = TextFieldDefaults.colors(
                            focusedContainerColor = Color.Gray.copy(alpha = 0.1f),
                            unfocusedContainerColor = Color.Gray.copy(alpha = 0.1f),
                            focusedIndicatorColor = Color.Transparent,
                            unfocusedIndicatorColor = Color.Transparent
                        )
                    )
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            if (isSearching) {
                // Filter Chips
                LazyRow(
                    modifier = Modifier.fillMaxWidth(),
                    contentPadding = PaddingValues(horizontal = 16.dp, vertical = 8.dp),
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    items(viewModel.getFilters()) { filter ->
                        FilterChip(
                            selected = selectedFilter == filter,
                            onClick = { viewModel.selectFilter(filter) },
                            label = { Text(filter) }
                        )
                    }
                }

                Divider()

                // Search Results
                SearchResultsView(selectedFilter, searchText)
            } else {
                // Recent & Suggested Searches
                RecentAndSuggestedView(
                    recentSearches = recentSearches,
                    onClearAll = { viewModel.clearRecentSearches() },
                    onSearchClick = { searchText = it }
                )
            }
        }
    }
}

// MARK: - Recent & Suggested View
@Composable
fun RecentAndSuggestedView(
    recentSearches: List<String>,
    onClearAll: () -> Unit,
    onSearchClick: (String) -> Unit
) {
    LazyColumn(modifier = Modifier.fillMaxSize()) {
        // Recent Searches
        if (recentSearches.isNotEmpty()) {
            item {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 16.dp, vertical = 12.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(
                        "Recent",
                        fontSize = 18.sp,
                        fontWeight = FontWeight.Bold
                    )
                    TextButton(onClick = onClearAll) {
                        Text("Clear all", color = Color(0xFF007CFC))
                    }
                }
            }

            items(recentSearches.size) { index ->
                RecentSearchRow(
                    text = recentSearches[index],
                    onClick = { onSearchClick(recentSearches[index]) }
                )
            }
        }

        // Suggested Searches
        item {
            Text(
                "Suggested",
                fontSize = 18.sp,
                fontWeight = FontWeight.Bold,
                modifier = Modifier.padding(horizontal = 16.dp, vertical = 12.dp)
            )
        }

        items(10) { index ->
            SuggestedSearchRow(
                icon = R.drawable.ic_trending,
                text = "Trending topic $index",
                onClick = { onSearchClick("Trending topic $index") }
            )
        }
    }
}

@Composable
fun RecentSearchRow(text: String, onClick: () -> Unit) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(horizontal = 16.dp, vertical = 12.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Icon(
            painter = painterResource(R.drawable.ic_clock),
            contentDescription = null,
            tint = Color.Gray,
            modifier = Modifier.size(24.dp)
        )
        Spacer(modifier = Modifier.width(12.dp))
        Text(text, fontSize = 16.sp, modifier = Modifier.weight(1f))
        IconButton(onClick = { /* Remove */ }) {
            Icon(
                painter = painterResource(R.drawable.ic_close),
                contentDescription = "Remove",
                tint = Color.Gray
            )
        }
    }
}

@Composable
fun SuggestedSearchRow(icon: Int, text: String, onClick: () -> Unit) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(horizontal = 16.dp, vertical = 12.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Surface(
            modifier = Modifier.size(40.dp),
            shape = CircleShape,
            color = Color.Gray.copy(alpha = 0.2f)
        ) {
            Box(contentAlignment = Alignment.Center) {
                Icon(
                    painter = painterResource(icon),
                    contentDescription = null,
                    tint = Color(0xFF007CFC)
                )
            }
        }
        Spacer(modifier = Modifier.width(12.dp))
        Text(text, fontSize = 16.sp)
    }
}

// MARK: - Search Results View
@Composable
fun SearchResultsView(filter: String, query: String) {
    LazyColumn(modifier = Modifier.fillMaxSize()) {
        when (filter) {
            "All", "People" -> {
                items(10) { index ->
                    PersonResultRow(
                        name = "User $index",
                        username = "user$index",
                        onClick = {}
                    )
                }
            }
            "Posts" -> {
                items(10) { index ->
                    PostResultRow(
                        author = "User $index",
                        content = "This is a post about $query...",
                        timeAgo = "2h",
                        onClick = {}
                    )
                }
            }
            "Photos" -> {
                item {
                    LazyVerticalGrid(
                        columns = GridCells.Fixed(3),
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(600.dp),
                        contentPadding = PaddingValues(2.dp)
                    ) {
                        items(15) {
                            Surface(
                                modifier = Modifier
                                    .aspectRatio(1f)
                                    .padding(2.dp),
                                color = Color.Gray.copy(alpha = 0.3f)
                            ) {}
                        }
                    }
                }
            }
            "Videos" -> {
                item {
                    LazyVerticalGrid(
                        columns = GridCells.Fixed(2),
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(600.dp),
                        contentPadding = PaddingValues(4.dp)
                    ) {
                        items(10) {
                            Surface(
                                modifier = Modifier
                                    .aspectRatio(9f / 16f)
                                    .padding(4.dp),
                                color = Color.Gray.copy(alpha = 0.3f)
                            ) {}
                        }
                    }
                }
            }
            "Pages" -> {
                items(10) { index ->
                    PageResultRow(
                        name = "Page $index",
                        category = "Category",
                        likes = "1.2K likes",
                        onClick = {}
                    )
                }
            }
            "Groups" -> {
                items(10) { index ->
                    GroupResultRow(
                        name = "Group $index",
                        members = "15K members",
                        onClick = {}
                    )
                }
            }
            "Events" -> {
                items(10) { index ->
                    EventResultRow(
                        name = "Event $index",
                        date = "Tomorrow at 7:00 PM",
                        onClick = {}
                    )
                }
            }
        }
    }
}

// Result Row Composables
@Composable
fun PersonResultRow(name: String, username: String, onClick: () -> Unit) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(16.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Surface(
            modifier = Modifier.size(48.dp),
            shape = CircleShape,
            color = Color.Gray.copy(alpha = 0.2f)
        ) {}
        Spacer(modifier = Modifier.width(12.dp))
        Column(modifier = Modifier.weight(1f)) {
            Text(name, fontSize = 16.sp, fontWeight = FontWeight.SemiBold)
            Text("@$username", fontSize = 14.sp, color = Color.Gray)
        }
        Button(
            onClick = {},
            colors = ButtonDefaults.buttonColors(containerColor = Color(0xFF007CFC))
        ) {
            Text("Add Friend")
        }
    }
}

@Composable
fun PostResultRow(author: String, content: String, timeAgo: String, onClick: () -> Unit) {
    Column(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(16.dp)
    ) {
        Row(verticalAlignment = Alignment.CenterVertically) {
            Surface(
                modifier = Modifier.size(40.dp),
                shape = CircleShape,
                color = Color.Gray.copy(alpha = 0.2f)
            ) {}
            Spacer(modifier = Modifier.width(8.dp))
            Column {
                Text(author, fontSize = 15.sp, fontWeight = FontWeight.SemiBold)
                Text(timeAgo, fontSize = 13.sp, color = Color.Gray)
            }
        }
        Spacer(modifier = Modifier.height(8.dp))
        Text(content, fontSize = 15.sp, maxLines = 3)
    }
    Divider()
}

@Composable
fun PageResultRow(name: String, category: String, likes: String, onClick: () -> Unit) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(16.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Surface(
            modifier = Modifier.size(48.dp),
            shape = RoundedCornerShape(8.dp),
            color = Color.Gray.copy(alpha = 0.2f)
        ) {}
        Spacer(modifier = Modifier.width(12.dp))
        Column(modifier = Modifier.weight(1f)) {
            Text(name, fontSize = 16.sp, fontWeight = FontWeight.SemiBold)
            Text(category, fontSize = 14.sp, color = Color.Gray)
            Text(likes, fontSize = 13.sp, color = Color.Gray)
        }
        Button(
            onClick = {},
            colors = ButtonDefaults.buttonColors(containerColor = Color(0xFF007CFC))
        ) {
            Text("Like")
        }
    }
}

@Composable
fun GroupResultRow(name: String, members: String, onClick: () -> Unit) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(16.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Surface(
            modifier = Modifier.size(48.dp),
            shape = RoundedCornerShape(8.dp),
            color = Color.Gray.copy(alpha = 0.2f)
        ) {}
        Spacer(modifier = Modifier.width(12.dp))
        Column(modifier = Modifier.weight(1f)) {
            Text(name, fontSize = 16.sp, fontWeight = FontWeight.SemiBold)
            Text(members, fontSize = 14.sp, color = Color.Gray)
        }
        Button(
            onClick = {},
            colors = ButtonDefaults.buttonColors(containerColor = Color(0xFF007CFC))
        ) {
            Text("Join")
        }
    }
}

@Composable
fun EventResultRow(name: String, date: String, onClick: () -> Unit) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(16.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Surface(
            modifier = Modifier.size(48.dp),
            shape = RoundedCornerShape(8.dp),
            color = Color.Gray.copy(alpha = 0.2f)
        ) {}
        Spacer(modifier = Modifier.width(12.dp))
        Column(modifier = Modifier.weight(1f)) {
            Text(name, fontSize = 16.sp, fontWeight = FontWeight.SemiBold)
            Text(date, fontSize = 14.sp, color = Color.Gray)
        }
        Button(
            onClick = {},
            colors = ButtonDefaults.buttonColors(containerColor = Color.Gray.copy(alpha = 0.2f)),
            contentPadding = PaddingValues(horizontal = 16.dp, vertical = 8.dp)
        ) {
            Text("Interested", color = Color.Black)
        }
    }
}

// MARK: - ViewModel
class EntativaSearchViewModel : ViewModel() {
    private val _selectedFilter = MutableStateFlow("All")
    val selectedFilter: StateFlow<String> = _selectedFilter

    private val _recentSearches = MutableStateFlow(
        listOf("John Doe", "Travel Photos", "Tech News")
    )
    val recentSearches: StateFlow<List<String>> = _recentSearches

    fun getFilters(): List<String> {
        return listOf("All", "People", "Posts", "Photos", "Videos", "Pages", "Groups", "Events")
    }

    fun selectFilter(filter: String) {
        _selectedFilter.value = filter
    }

    fun clearRecentSearches() {
        _recentSearches.value = emptyList()
    }
}
