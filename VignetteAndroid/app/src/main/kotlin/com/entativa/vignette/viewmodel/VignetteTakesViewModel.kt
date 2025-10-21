package com.entativa.vignette.ui.takes

import androidx.lifecycle.ViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow

class VignetteTakesViewModel : ViewModel() {
    private val _takes = MutableStateFlow(mockTakes)
    val takes: StateFlow<List<VignetteTake>> = _takes.asStateFlow()
    
    private val _isLoading = MutableStateFlow(false)
    val isLoading: StateFlow<Boolean> = _isLoading.asStateFlow()
    
    suspend fun loadMore() {
        _isLoading.value = true
        kotlinx.coroutines.delay(1000)
        _isLoading.value = false
    }
    
    companion object {
        private val mockTakes = listOf(
            VignetteTake(
                id = "1",
                username = "photo.vibes",
                caption = "Golden hour in the city 🌆✨ #photography",
                audioName = "Chill Vibes - Lofi Beats",
                likesCount = 234500,
                commentsCount = 3421,
                sharesCount = 8934,
                viewsCount = 1234500
            ),
            VignetteTake(
                id = "2",
                username = "fit.journey",
                caption = "Morning routine that changed my life! 💪",
                audioName = "Workout Mix 2024",
                likesCount = 89300,
                commentsCount = 1234,
                sharesCount = 3456,
                viewsCount = 567800
            ),
            VignetteTake(
                id = "3",
                username = "chef.athome",
                caption = "60-second pasta that tastes like heaven! 🍝",
                audioName = "Cooking Time - Kitchen Beats",
                likesCount = 456700,
                commentsCount = 12345,
                sharesCount = 23456,
                viewsCount = 2345600
            ),
            VignetteTake(
                id = "4",
                username = "style.daily",
                caption = "Transforming thrift finds into designer looks ✨👗",
                audioName = "Fashion Week Runway",
                likesCount = 678900,
                commentsCount = 8765,
                sharesCount = 34567,
                viewsCount = 3456700
            ),
            VignetteTake(
                id = "5",
                username = "pet.moments",
                caption = "When your dog understands the assignment 😂🐕",
                audioName = "Funny Pet Sounds",
                likesCount = 890100,
                commentsCount = 23456,
                sharesCount = 45678,
                viewsCount = 4567800
            )
        )
    }
}
