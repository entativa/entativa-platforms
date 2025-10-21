package com.entativa.vignette.ui.home

import androidx.lifecycle.ViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow

class VignetteHomeViewModel : ViewModel() {
    private val _posts = MutableStateFlow(mockPosts)
    val posts: StateFlow<List<VignettePost>> = _posts.asStateFlow()
    
    private val _isLoading = MutableStateFlow(false)
    val isLoading: StateFlow<Boolean> = _isLoading.asStateFlow()
    
    suspend fun refreshFeed() {
        _isLoading.value = true
        kotlinx.coroutines.delay(1000)
        _isLoading.value = false
    }
    
    companion object {
        private val mockPosts = listOf(
            VignettePost(
                id = "1",
                userName = "alexjohnson",
                location = "San Francisco, CA",
                caption = "Golden hour magic ✨ Perfect end to a perfect day",
                mediaUrl = "",
                likesCount = 1247,
                commentsCount = 89,
                timestamp = "2 HOURS AGO"
            ),
            VignettePost(
                id = "2",
                userName = "sarahcreative",
                location = "Brooklyn, NY",
                caption = "New artwork incoming! Can't wait to share the full collection 🎨",
                mediaUrl = "",
                likesCount = 2134,
                commentsCount = 156,
                timestamp = "5 HOURS AGO"
            ),
            VignettePost(
                id = "3",
                userName = "mikefitness",
                location = null,
                caption = "Day 100 of my fitness journey! 💪",
                mediaUrl = "",
                likesCount = 892,
                commentsCount = 67,
                timestamp = "1 DAY AGO"
            )
        )
    }
}
