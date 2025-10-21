package com.entativa.ui.takes

import android.app.Application
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.viewModelScope
import com.entativa.network.TakesAPIClient
import com.entativa.network.TakeData
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

class EntativaTakesViewModel(application: Application) : AndroidViewModel(application) {
    private val apiClient = TakesAPIClient(application)
    
    private val _takes = MutableStateFlow<List<TakeData>>(emptyList())
    val takes: StateFlow<List<TakeData>> = _takes.asStateFlow()
    
    private val _isLoading = MutableStateFlow(false)
    val isLoading: StateFlow<Boolean> = _isLoading.asStateFlow()
    
    private var currentPage = 1
    private var hasMore = true
    
    init {
        loadFeed()
    }
    
    fun loadFeed() {
        viewModelScope.launch {
            _isLoading.value = true
            
            apiClient.getFeed(page = 1, limit = 10).fold(
                onSuccess = { response ->
                    _takes.value = response.takes
                    hasMore = response.hasMore
                    currentPage = 1
                },
                onFailure = { error ->
                    // Fallback to mock data
                    _takes.value = mockTakesData
                }
            )
            
            _isLoading.value = false
        }
    }
    
    fun loadMore() {
        if (!hasMore || _isLoading.value) return
        
        viewModelScope.launch {
            _isLoading.value = true
            currentPage++
            
            apiClient.getFeed(page = currentPage, limit = 10).fold(
                onSuccess = { response ->
                    _takes.value = _takes.value + response.takes
                    hasMore = response.hasMore
                },
                onFailure = { /* handle error */ }
            )
            
            _isLoading.value = false
        }
    }
    
    fun likeTake(takeId: String) {
        viewModelScope.launch {
            apiClient.likeTake(takeId).fold(
                onSuccess = { updatedTake ->
                    _takes.value = _takes.value.map { 
                        if (it.id == takeId) updatedTake else it 
                    }
                },
                onFailure = { /* handle error */ }
            )
        }
    }
    
    fun unlikeTake(takeId: String) {
        viewModelScope.launch {
            apiClient.unlikeTake(takeId).fold(
                onSuccess = { updatedTake ->
                    _takes.value = _takes.value.map { 
                        if (it.id == takeId) updatedTake else it 
                    }
                },
                onFailure = { /* handle error */ }
            )
        }
    }
    
    companion object {
        private val mockTakesData = listOf(
            TakeData(
                id = "1",
                userId = "user1",
                username = "alexcreator",
                userAvatar = null,
                videoUrl = "https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_1mb.mp4",
                thumbnailUrl = null,
                caption = "Check out this amazing transformation! 💪 #fitness #motivation",
                audioName = "Original Audio - alexcreator",
                audioUrl = null,
                duration = 30,
                likesCount = 45200,
                commentsCount = 892,
                sharesCount = 1234,
                viewsCount = 234500,
                isLiked = false,
                isSaved = false,
                hashtags = listOf("fitness", "motivation"),
                createdAt = "2025-10-18T12:00:00Z"
            ),
            TakeData(
                id = "2",
                userId = "user2",
                username = "foodie.life",
                userAvatar = null,
                videoUrl = "https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_2mb.mp4",
                thumbnailUrl = null,
                caption = "Best pasta recipe ever! 🍝 Try it and let me know what you think!",
                audioName = "Cooking Vibes - Sound Library",
                audioUrl = null,
                duration = 45,
                likesCount = 78300,
                commentsCount = 1456,
                sharesCount = 2890,
                viewsCount = 456700,
                isLiked = false,
                isSaved = false,
                hashtags = listOf("cooking", "food"),
                createdAt = "2025-10-18T10:00:00Z"
            ),
            TakeData(
                id = "3",
                userId = "user3",
                username = "travel.with.me",
                userAvatar = null,
                videoUrl = "https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_5mb.mp4",
                thumbnailUrl = null,
                caption = "Hidden gems in Bali you NEED to visit! 🌴✨ #travel #bali",
                audioName = "Tropical Summer - Music Mix",
                audioUrl = null,
                duration = 60,
                likesCount = 123400,
                commentsCount = 3421,
                sharesCount = 5678,
                viewsCount = 890200,
                isLiked = false,
                isSaved = false,
                hashtags = listOf("travel", "bali"),
                createdAt = "2025-10-17T18:00:00Z"
            )
        )
        
        private val mockTakes = listOf(
            Take(
                id = "1",
                username = "alexcreator",
                caption = "Check out this amazing transformation! 💪 #fitness #motivation",
                audioName = "Original Audio - alexcreator",
                likesCount = 45200,
                commentsCount = 892,
                sharesCount = 1234,
                viewsCount = 234500
            ),
            Take(
                id = "2",
                username = "foodie.life",
                caption = "Best pasta recipe ever! 🍝 Try it and let me know what you think!",
                audioName = "Cooking Vibes - Sound Library",
                likesCount = 78300,
                commentsCount = 1456,
                sharesCount = 2890,
                viewsCount = 456700
            ),
            Take(
                id = "3",
                username = "travel.with.me",
                caption = "Hidden gems in Bali you NEED to visit! 🌴✨ #travel #bali",
                audioName = "Tropical Summer - Music Mix",
                likesCount = 123400,
                commentsCount = 3421,
                sharesCount = 5678,
                viewsCount = 890200
            ),
            Take(
                id = "4",
                username = "tech.reviews",
                caption = "iPhone 16 Pro vs Samsung S24 Ultra - The TRUTH! 📱",
                audioName = "Tech Beat - Audio Track",
                likesCount = 56700,
                commentsCount = 2134,
                sharesCount = 3456,
                viewsCount = 567800
            ),
            Take(
                id = "5",
                username = "comedy.gold",
                caption = "When your code finally works 😂💻 #coding #funny",
                audioName = "Funny Moments - Sound Effect",
                likesCount = 234500,
                commentsCount = 8976,
                sharesCount = 12345,
                viewsCount = 1234500
            )
        )
    }
}
