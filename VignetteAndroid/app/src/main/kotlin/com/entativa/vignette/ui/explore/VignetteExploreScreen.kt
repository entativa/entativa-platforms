package com.entativa.vignette.ui.explore

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
import com.entativa.vignette.R
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow

// MARK: - Vignette Explore Screen (Instagram-Style)
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun VignetteExploreScreen(
    viewModel: VignetteExploreViewModel = viewModel()
) {
    var searchText by remember { mutableStateOf("") }
    val isSearching = searchText.isNotEmpty()
    val selectedTab by viewModel.selectedTab.collectAsState()
    val recentSearches by viewModel.recentSearches.collectAsState()

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    OutlinedTextField(
                        value = searchText,
                        onValueChange = { searchText = it },
                        placeholder = { Text("Search") },
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
                        shape = RoundedCornerShape(12.dp),
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
                // Search Tabs
                LazyRow(
                    modifier = Modifier.fillMaxWidth(),
                    contentPadding = PaddingValues(horizontal = 16.dp, vertical = 8.dp),
                    horizontalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    items(viewModel.getTabs()) { tab ->
                        SearchTabButton(
                            text = tab,
                            selected = selectedTab == tab,
                            onClick = { viewModel.selectTab(tab) }
                        )
                    }
                }

                Divider()

                // Search Results
                SearchResultsContent(selectedTab, searchText)
            } else {
                // Explore Grid
                ExploreGridView()
            }
        }
    }
}

// MARK: - Search Tab Button
@Composable
fun SearchTabButton(text: String, selected: Boolean, onClick: () -> Unit) {
    Surface(
        onClick = onClick,
        shape = RoundedCornerShape(8.dp),
        color = if (selected) Color.Black else Color.Transparent,
        border = if (!selected) androidx.compose.foundation.BorderStroke(1.dp, Color.Gray.copy(alpha = 0.3f)) else null
    ) {
        Text(
            text,
            fontSize = 14.sp,
            fontWeight = FontWeight.Medium,
            color = if (selected) Color.White else Color.Black,
            modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp)
        )
    }
}

// MARK: - Explore Grid View
@Composable
fun ExploreGridView() {
    LazyVerticalGrid(
        columns = GridCells.Fixed(3),
        modifier = Modifier.fillMaxSize(),
        contentPadding = PaddingValues(1.dp)
    ) {
        items(30) { index ->
            Box(
                modifier = Modifier
                    .aspectRatio(1f)
                    .padding(1.dp)
            ) {
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = Color.Gray.copy(alpha = 0.3f)
                ) {
                    // Photo placeholder
                }

                // Video indicator for some items
                if (index % 5 == 0) {
                    Icon(
                        painter = painterResource(R.drawable.ic_play),
                        contentDescription = "Video",
                        tint = Color.White,
                        modifier = Modifier
                            .align(Alignment.TopEnd)
                            .padding(8.dp)
                            .size(20.dp)
                    )
                }
            }
        }
    }
}

// MARK: - Search Results Content
@Composable
fun SearchResultsContent(tab: String, query: String) {
    when (tab) {
        "Top" -> TopResultsView(query)
        "Accounts" -> AccountsResultsView(query)
        "Audio" -> AudioResultsView(query)
        "Tags" -> TagsResultsView(query)
        "Places" -> PlacesResultsView(query)
    }
}

// MARK: - Top Results View
@Composable
fun TopResultsView(query: String) {
    LazyColumn(modifier = Modifier.fillMaxSize()) {
        // Top Accounts
        item {
            Text(
                "Accounts",
                fontSize = 16.sp,
                fontWeight = FontWeight.Bold,
                modifier = Modifier.padding(16.dp)
            )
        }

        items(3) { index ->
            AccountResultRow(
                username = "user$index",
                name = "User Name $index",
                isVerified = index == 0,
                onClick = {}
            )
        }

        // Top Posts Grid
        item {
            Text(
                "Posts",
                fontSize = 16.sp,
                fontWeight = FontWeight.Bold,
                modifier = Modifier.padding(16.dp)
            )
        }

        item {
            LazyVerticalGrid(
                columns = GridCells.Fixed(3),
                modifier = Modifier
                    .fillMaxWidth()
                    .height(400.dp),
                contentPadding = PaddingValues(1.dp)
            ) {
                items(9) {
                    Surface(
                        modifier = Modifier
                            .aspectRatio(1f)
                            .padding(1.dp),
                        color = Color.Gray.copy(alpha = 0.3f)
                    ) {}
                }
            }
        }
    }
}

// MARK: - Accounts Results View
@Composable
fun AccountsResultsView(query: String) {
    LazyColumn(modifier = Modifier.fillMaxSize()) {
        items(20) { index ->
            AccountResultRow(
                username = "user$index",
                name = "User Name $index",
                isVerified = index % 5 == 0,
                onClick = {}
            )
        }
    }
}

@Composable
fun AccountResultRow(
    username: String,
    name: String,
    isVerified: Boolean,
    onClick: () -> Unit
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(horizontal = 16.dp, vertical = 8.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Surface(
            modifier = Modifier.size(48.dp),
            shape = CircleShape,
            color = Color.Gray.copy(alpha = 0.2f)
        ) {}

        Spacer(modifier = Modifier.width(12.dp))

        Column(modifier = Modifier.weight(1f)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Text(username, fontSize = 15.sp, fontWeight = FontWeight.SemiBold)
                if (isVerified) {
                    Spacer(modifier = Modifier.width(4.dp))
                    Icon(
                        painter = painterResource(R.drawable.ic_verified),
                        contentDescription = "Verified",
                        tint = Color(0xFF007CFC),
                        modifier = Modifier.size(14.dp)
                    )
                }
            }
            Text(name, fontSize = 14.sp, color = Color.Gray)
        }

        Button(
            onClick = {},
            colors = ButtonDefaults.buttonColors(
                containerColor = Color(0xFF007CFC)
            ),
            contentPadding = PaddingValues(horizontal = 16.dp, vertical = 6.dp)
        ) {
            Text("Follow", fontSize = 14.sp)
        }
    }
}

// MARK: - Audio Results View
@Composable
fun AudioResultsView(query: String) {
    LazyColumn(modifier = Modifier.fillMaxSize()) {
        items(15) { index ->
            AudioResultRow(
                trackName = "Audio Track $index",
                artist = "Artist Name $index",
                usedCount = "${(index + 1) * 10}K",
                onClick = {}
            )
        }
    }
}

@Composable
fun AudioResultRow(
    trackName: String,
    artist: String,
    usedCount: String,
    onClick: () -> Unit
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(horizontal = 16.dp, vertical = 12.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Surface(
            modifier = Modifier.size(48.dp),
            shape = RoundedCornerShape(4.dp),
            color = Color.Gray.copy(alpha = 0.2f)
        ) {
            Box(contentAlignment = Alignment.Center) {
                Icon(
                    painter = painterResource(R.drawable.ic_music),
                    contentDescription = null,
                    tint = Color.Gray
                )
            }
        }

        Spacer(modifier = Modifier.width(12.dp))

        Column(modifier = Modifier.weight(1f)) {
            Text(trackName, fontSize = 15.sp, fontWeight = FontWeight.SemiBold)
            Text("$artist • $usedCount posts", fontSize = 13.sp, color = Color.Gray)
        }

        Icon(
            painter = painterResource(R.drawable.ic_chevron_right),
            contentDescription = null,
            tint = Color.Gray
        )
    }
}

// MARK: - Tags Results View
@Composable
fun TagsResultsView(query: String) {
    LazyColumn(modifier = Modifier.fillMaxSize()) {
        items(15) { index ->
            TagResultRow(
                tag = "tag$index",
                postCount = "${(index + 1) * 100}K",
                onClick = {}
            )
        }
    }
}

@Composable
fun TagResultRow(
    tag: String,
    postCount: String,
    onClick: () -> Unit
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(horizontal = 16.dp, vertical = 12.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Surface(
            modifier = Modifier.size(48.dp),
            shape = RoundedCornerShape(4.dp),
            color = Color.Gray.copy(alpha = 0.2f)
        ) {
            Box(contentAlignment = Alignment.Center) {
                Text(
                    "#",
                    fontSize = 24.sp,
                    fontWeight = FontWeight.Bold,
                    color = Color.Gray
                )
            }
        }

        Spacer(modifier = Modifier.width(12.dp))

        Column(modifier = Modifier.weight(1f)) {
            Text("#$tag", fontSize = 15.sp, fontWeight = FontWeight.SemiBold)
            Text("$postCount posts", fontSize = 13.sp, color = Color.Gray)
        }

        Icon(
            painter = painterResource(R.drawable.ic_chevron_right),
            contentDescription = null,
            tint = Color.Gray
        )
    }
}

// MARK: - Places Results View
@Composable
fun PlacesResultsView(query: String) {
    LazyColumn(modifier = Modifier.fillMaxSize()) {
        items(15) { index ->
            PlaceResultRow(
                place = "Place Name $index",
                location = "City, Country",
                postCount = "${(index + 1) * 50}K",
                onClick = {}
            )
        }
    }
}

@Composable
fun PlaceResultRow(
    place: String,
    location: String,
    postCount: String,
    onClick: () -> Unit
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(horizontal = 16.dp, vertical = 12.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Surface(
            modifier = Modifier.size(48.dp),
            shape = RoundedCornerShape(4.dp),
            color = Color.Gray.copy(alpha = 0.2f)
        ) {
            Box(contentAlignment = Alignment.Center) {
                Icon(
                    painter = painterResource(R.drawable.ic_location),
                    contentDescription = null,
                    tint = Color.Gray
                )
            }
        }

        Spacer(modifier = Modifier.width(12.dp))

        Column(modifier = Modifier.weight(1f)) {
            Text(place, fontSize = 15.sp, fontWeight = FontWeight.SemiBold)
            Text("$location • $postCount posts", fontSize = 13.sp, color = Color.Gray)
        }

        Icon(
            painter = painterResource(R.drawable.ic_chevron_right),
            contentDescription = null,
            tint = Color.Gray
        )
    }
}

// MARK: - ViewModel
class VignetteExploreViewModel : ViewModel() {
    private val _selectedTab = MutableStateFlow("Top")
    val selectedTab: StateFlow<String> = _selectedTab

    private val _recentSearches = MutableStateFlow(
        listOf("travel", "food", "art")
    )
    val recentSearches: StateFlow<List<String>> = _recentSearches

    fun getTabs(): List<String> {
        return listOf("Top", "Accounts", "Audio", "Tags", "Places")
    }

    fun selectTab(tab: String) {
        _selectedTab.value = tab
    }

    fun clearRecentSearches() {
        _recentSearches.value = emptyList()
    }
}
