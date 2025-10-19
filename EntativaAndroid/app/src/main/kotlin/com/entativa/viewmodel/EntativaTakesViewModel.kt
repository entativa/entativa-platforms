package com.entativa.ui.takes

import androidx.lifecycle.ViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow

class EntativaTakesViewModel : ViewModel() {
    private val _takes = MutableStateFlow(mockTakes)
    val takes: StateFlow<List<Take>> = _takes.asStateFlow()
    
    private val _isLoading = MutableStateFlow(false)
    val isLoading: StateFlow<Boolean> = _isLoading.asStateFlow()
    
    suspend fun loadMore() {
        _isLoading.value = true
        kotlinx.coroutines.delay(1000)
        _isLoading.value = false
    }
    
    companion object {
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
