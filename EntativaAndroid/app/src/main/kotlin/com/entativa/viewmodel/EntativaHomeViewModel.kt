package com.entativa.ui.home

import androidx.lifecycle.ViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow

class EntativaHomeViewModel : ViewModel() {
    private val _posts = MutableStateFlow(mockPosts)
    val posts: StateFlow<List<Post>> = _posts.asStateFlow()
    
    private val _isLoading = MutableStateFlow(false)
    val isLoading: StateFlow<Boolean> = _isLoading.asStateFlow()
    
    suspend fun refreshFeed() {
        _isLoading.value = true
        // Simulate network call
        kotlinx.coroutines.delay(1000)
        _isLoading.value = false
    }
    
    companion object {
        private val mockPosts = listOf(
            Post(
                id = "1",
                userName = "John Doe",
                timestamp = "2h ago",
                text = "Just finished an amazing workout session! 💪 Feeling great and ready to tackle the day. Who else is staying active today?",
                mediaUrls = listOf(""),
                likesCount = 124,
                commentsCount = 23,
                sharesCount = 5
            ),
            Post(
                id = "2",
                userName = "Jane Smith",
                timestamp = "5h ago",
                text = "Beautiful sunset at the beach today 🌅",
                mediaUrls = listOf("", "", ""),
                likesCount = 489,
                commentsCount = 67,
                sharesCount = 34
            ),
            Post(
                id = "3",
                userName = "Mike Johnson",
                timestamp = "1d ago",
                text = "New project launch! So excited to share what we've been working on. Stay tuned for more updates!",
                mediaUrls = emptyList(),
                likesCount = 256,
                commentsCount = 45,
                sharesCount = 12
            )
        )
    }
}
