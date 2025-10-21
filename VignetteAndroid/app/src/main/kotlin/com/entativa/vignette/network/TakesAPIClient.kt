package com.entativa.vignette.network

import android.content.Context
import com.google.gson.Gson
import com.google.gson.annotations.SerializedName
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import okhttp3.*
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.RequestBody.Companion.toRequestBody
import java.io.IOException

class TakesAPIClient(private val context: Context) {
    private val client = OkHttpClient.Builder()
        .connectTimeout(30, java.util.concurrent.TimeUnit.SECONDS)
        .readTimeout(30, java.util.concurrent.TimeUnit.SECONDS)
        .build()
    
    private val gson = Gson()
    
    private val baseUrl = if (BuildConfig.DEBUG) {
        "http://10.0.2.2:8002/api/v1"  // Android emulator localhost
    } else {
        "https://api.vignette.app/api/v1"
    }
    
    // MARK: - Get Feed
    suspend fun getFeed(page: Int = 1, limit: Int = 10): Result<TakesFeedResponse> = withContext(Dispatchers.IO) {
        try {
            val url = "$baseUrl/takes/feed?page=$page&limit=$limit"
            
            val requestBuilder = Request.Builder()
                .url(url)
                .get()
            
            // Add auth token if available
            getAuthToken()?.let { token ->
                requestBuilder.addHeader("Authorization", "Bearer $token")
            }
            
            val request = requestBuilder.build()
            val response = client.newCall(request).execute()
            
            val responseBody = response.body?.string() ?: throw IOException("Empty response")
            
            if (!response.isSuccessful) {
                return@withContext Result.failure(Exception("HTTP ${response.code}"))
            }
            
            val apiResponse = gson.fromJson(responseBody, TakesAPIResponse::class.java)
            
            if (apiResponse.success && apiResponse.data != null) {
                Result.success(apiResponse.data)
            } else {
                Result.failure(Exception(apiResponse.error ?: "Unknown error"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    
    // MARK: - Like Take
    suspend fun likeTake(takeId: String): Result<TakeData> = withContext(Dispatchers.IO) {
        try {
            val url = "$baseUrl/takes/$takeId/like"
            val token = getAuthToken() ?: return@withContext Result.failure(Exception("Not authenticated"))
            
            val request = Request.Builder()
                .url(url)
                .post("".toRequestBody())
                .addHeader("Authorization", "Bearer $token")
                .build()
            
            val response = client.newCall(request).execute()
            val responseBody = response.body?.string() ?: throw IOException("Empty response")
            
            if (!response.isSuccessful) {
                return@withContext Result.failure(Exception("HTTP ${response.code}"))
            }
            
            val apiResponse = gson.fromJson(responseBody, TakeSingleAPIResponse::class.java)
            
            if (apiResponse.success && apiResponse.data != null) {
                Result.success(apiResponse.data)
            } else {
                Result.failure(Exception(apiResponse.error ?: "Failed to like"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    
    // MARK: - Unlike Take
    suspend fun unlikeTake(takeId: String): Result<TakeData> = withContext(Dispatchers.IO) {
        try {
            val url = "$baseUrl/takes/$takeId/unlike"
            val token = getAuthToken() ?: return@withContext Result.failure(Exception("Not authenticated"))
            
            val request = Request.Builder()
                .url(url)
                .post("".toRequestBody())
                .addHeader("Authorization", "Bearer $token")
                .build()
            
            val response = client.newCall(request).execute()
            val responseBody = response.body?.string() ?: throw IOException("Empty response")
            
            if (!response.isSuccessful) {
                return@withContext Result.failure(Exception("HTTP ${response.code}"))
            }
            
            val apiResponse = gson.fromJson(responseBody, TakeSingleAPIResponse::class.java)
            
            if (apiResponse.success && apiResponse.data != null) {
                Result.success(apiResponse.data)
            } else {
                Result.failure(Exception(apiResponse.error ?: "Failed to unlike"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    
    // MARK: - Get Comments
    suspend fun getComments(takeId: String, page: Int = 1): Result<CommentsResponse> = withContext(Dispatchers.IO) {
        try {
            val url = "$baseUrl/takes/$takeId/comments?page=$page"
            
            val request = Request.Builder()
                .url(url)
                .get()
                .build()
            
            val response = client.newCall(request).execute()
            val responseBody = response.body?.string() ?: throw IOException("Empty response")
            
            if (!response.isSuccessful) {
                return@withContext Result.failure(Exception("HTTP ${response.code}"))
            }
            
            val apiResponse = gson.fromJson(responseBody, CommentsAPIResponse::class.java)
            
            if (apiResponse.success && apiResponse.data != null) {
                Result.success(apiResponse.data)
            } else {
                Result.failure(Exception(apiResponse.error ?: "Failed to fetch comments"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    
    // MARK: - Add Comment
    suspend fun addComment(takeId: String, text: String): Result<Unit> = withContext(Dispatchers.IO) {
        try {
            val url = "$baseUrl/takes/$takeId/comments"
            val token = getAuthToken() ?: return@withContext Result.failure(Exception("Not authenticated"))
            
            val body = mapOf("text" to text)
            val jsonBody = gson.toJson(body).toRequestBody("application/json".toMediaType())
            
            val request = Request.Builder()
                .url(url)
                .post(jsonBody)
                .addHeader("Authorization", "Bearer $token")
                .build()
            
            val response = client.newCall(request).execute()
            
            if (response.isSuccessful) {
                Result.success(Unit)
            } else {
                Result.failure(Exception("HTTP ${response.code}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    
    private fun getAuthToken(): String? {
        return try {
            val masterKey = androidx.security.crypto.MasterKey.Builder(context)
                .setKeyScheme(androidx.security.crypto.MasterKey.KeyScheme.AES256_GCM)
                .build()
            
            val securePrefs = androidx.security.crypto.EncryptedSharedPreferences.create(
                context,
                "auth_prefs",
                masterKey,
                androidx.security.crypto.EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
                androidx.security.crypto.EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
            )
            
            securePrefs.getString("access_token", null)
        } catch (e: Exception) {
            null
        }
    }
}

// MARK: - Models
data class TakeData(
    @SerializedName("id") val id: String,
    @SerializedName("user_id") val userId: String,
    @SerializedName("username") val username: String,
    @SerializedName("user_avatar") val userAvatar: String?,
    @SerializedName("video_url") val videoUrl: String,
    @SerializedName("thumbnail_url") val thumbnailUrl: String?,
    @SerializedName("caption") val caption: String,
    @SerializedName("audio_name") val audioName: String,
    @SerializedName("audio_url") val audioUrl: String?,
    @SerializedName("duration") val duration: Int,
    @SerializedName("likes_count") val likesCount: Int,
    @SerializedName("comments_count") val commentsCount: Int,
    @SerializedName("shares_count") val sharesCount: Int,
    @SerializedName("views_count") val viewsCount: Int,
    @SerializedName("is_liked") val isLiked: Boolean,
    @SerializedName("is_saved") val isSaved: Boolean,
    @SerializedName("hashtags") val hashtags: List<String>?,
    @SerializedName("created_at") val createdAt: String
)

data class TakesFeedResponse(
    @SerializedName("takes") val takes: List<TakeData>,
    @SerializedName("page") val page: Int,
    @SerializedName("limit") val limit: Int,
    @SerializedName("has_more") val hasMore: Boolean
)

data class CommentData(
    @SerializedName("id") val id: String,
    @SerializedName("take_id") val takeId: String,
    @SerializedName("user_id") val userId: String,
    @SerializedName("username") val username: String,
    @SerializedName("user_avatar") val userAvatar: String?,
    @SerializedName("text") val text: String,
    @SerializedName("likes_count") val likesCount: Int,
    @SerializedName("is_liked") val isLiked: Boolean,
    @SerializedName("created_at") val createdAt: String
)

data class CommentsResponse(
    @SerializedName("comments") val comments: List<CommentData>,
    @SerializedName("page") val page: Int,
    @SerializedName("limit") val limit: Int,
    @SerializedName("has_more") val hasMore: Boolean
)

// MARK: - API Response Wrappers
data class TakesAPIResponse(
    @SerializedName("success") val success: Boolean,
    @SerializedName("message") val message: String?,
    @SerializedName("data") val data: TakesFeedResponse?,
    @SerializedName("error") val error: String?
)

data class TakeSingleAPIResponse(
    @SerializedName("success") val success: Boolean,
    @SerializedName("message") val message: String?,
    @SerializedName("data") val data: TakeData?,
    @SerializedName("error") val error: String?
)

data class CommentsAPIResponse(
    @SerializedName("success") val success: Boolean,
    @SerializedName("message") val message: String?,
    @SerializedName("data") val data: CommentsResponse?,
    @SerializedName("error") val error: String?
)
