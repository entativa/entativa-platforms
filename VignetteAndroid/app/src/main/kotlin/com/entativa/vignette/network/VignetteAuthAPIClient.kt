package com.entativa.vignette.network

import android.content.Context
import android.content.SharedPreferences
import androidx.security.crypto.EncryptedSharedPreferences
import androidx.security.crypto.MasterKey
import com.google.gson.Gson
import com.google.gson.annotations.SerializedName
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody
import java.util.concurrent.TimeUnit

/**
 * Authentication API Client for Vignette Backend
 * Instagram-style authentication with enterprise-grade security
 */
class VignetteAuthAPIClient(private val context: Context) {
    
    companion object {
        private const val BASE_URL_DEBUG = "http://10.0.2.2:8002/api/v1"
        private const val BASE_URL_PRODUCTION = "https://api.vignette.app/api/v1"
        private const val TOKEN_KEY = "auth_token"
        
        @Volatile
        private var INSTANCE: VignetteAuthAPIClient? = null
        
        fun getInstance(context: Context): VignetteAuthAPIClient {
            return INSTANCE ?: synchronized(this) {
                INSTANCE ?: VignetteAuthAPIClient(context.applicationContext).also {
                    INSTANCE = it
                }
            }
        }
    }
    
    private val baseUrl = if (BuildConfig.DEBUG) BASE_URL_DEBUG else BASE_URL_PRODUCTION
    
    private val client = OkHttpClient.Builder()
        .connectTimeout(30, TimeUnit.SECONDS)
        .readTimeout(30, TimeUnit.SECONDS)
        .writeTimeout(30, TimeUnit.SECONDS)
        .build()
    
    private val gson = Gson()
    
    private val securePrefs: SharedPreferences by lazy {
        val masterKey = MasterKey.Builder(context)
            .setKeyScheme(MasterKey.KeyScheme.AES256_GCM)
            .build()
        
        EncryptedSharedPreferences.create(
            context,
            "vignette_secure_prefs",
            masterKey,
            EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
            EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
        )
    }
    
    // MARK: - Data Models
    
    data class SignUpRequest(
        val username: String,
        val email: String,
        @SerializedName("full_name") val fullName: String,
        val password: String
    )
    
    data class LoginRequest(
        @SerializedName("username_or_email") val usernameOrEmail: String,
        val password: String
    )
    
    data class AuthResponse(
        val success: Boolean,
        val message: String,
        val data: AuthData?
    )
    
    data class AuthData(
        val user: User,
        @SerializedName("access_token") val accessToken: String,
        @SerializedName("token_type") val tokenType: String,
        @SerializedName("expires_in") val expiresIn: Int
    )
    
    data class User(
        val id: String,
        val username: String,
        val email: String,
        @SerializedName("full_name") val fullName: String,
        val bio: String?,
        val website: String?,
        @SerializedName("profile_picture_url") val profilePictureUrl: String?,
        @SerializedName("is_private") val isPrivate: Boolean,
        @SerializedName("is_verified") val isVerified: Boolean,
        @SerializedName("is_active") val isActive: Boolean,
        @SerializedName("followers_count") val followersCount: Int,
        @SerializedName("following_count") val followingCount: Int,
        @SerializedName("posts_count") val postsCount: Int,
        @SerializedName("created_at") val createdAt: String
    )
    
    data class ErrorResponse(
        val success: Boolean,
        val error: String,
        val details: Map<String, List<String>>?
    )
    
    // MARK: - API Methods
    
    suspend fun signUp(
        username: String,
        email: String,
        fullName: String,
        password: String
    ): Result<AuthResponse> = withContext(Dispatchers.IO) {
        try {
            val requestBody = SignUpRequest(
                username = username,
                email = email,
                fullName = fullName,
                password = password
            )
            
            val json = gson.toJson(requestBody)
            val body = json.toRequestBody("application/json".toMediaType())
            
            val request = Request.Builder()
                .url("$baseUrl/auth/signup")
                .post(body)
                .build()
            
            val response = client.newCall(request).execute()
            
            if (response.isSuccessful) {
                val authResponse = gson.fromJson(response.body?.string(), AuthResponse::class.java)
                authResponse.data?.accessToken?.let { saveToken(it) }
                Result.success(authResponse)
            } else {
                val errorResponse = gson.fromJson(response.body?.string(), ErrorResponse::class.java)
                Result.failure(VignetteAuthException(errorResponse.error))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    
    suspend fun login(
        usernameOrEmail: String,
        password: String
    ): Result<AuthResponse> = withContext(Dispatchers.IO) {
        try {
            val requestBody = LoginRequest(
                usernameOrEmail = usernameOrEmail,
                password = password
            )
            
            val json = gson.toJson(requestBody)
            val body = json.toRequestBody("application/json".toMediaType())
            
            val request = Request.Builder()
                .url("$baseUrl/auth/login")
                .post(body)
                .build()
            
            val response = client.newCall(request).execute()
            
            if (response.isSuccessful) {
                val authResponse = gson.fromJson(response.body?.string(), AuthResponse::class.java)
                authResponse.data?.accessToken?.let { saveToken(it) }
                Result.success(authResponse)
            } else {
                val errorResponse = gson.fromJson(response.body?.string(), ErrorResponse::class.java)
                Result.failure(VignetteAuthException(errorResponse.error))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    
    suspend fun getCurrentUser(): Result<User> = withContext(Dispatchers.IO) {
        try {
            val token = getToken() ?: return@withContext Result.failure(
                VignetteAuthException("Not authenticated")
            )
            
            val request = Request.Builder()
                .url("$baseUrl/auth/me")
                .header("Authorization", "Bearer $token")
                .get()
                .build()
            
            val response = client.newCall(request).execute()
            
            if (response.isSuccessful) {
                val userResponse = gson.fromJson(response.body?.string(), UserResponse::class.java)
                Result.success(userResponse.data)
            } else {
                Result.failure(VignetteAuthException("Failed to get user"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    
    suspend fun logout(): Result<Unit> = withContext(Dispatchers.IO) {
        try {
            val token = getToken()
            
            if (token != null) {
                val request = Request.Builder()
                    .url("$baseUrl/auth/logout")
                    .header("Authorization", "Bearer $token")
                    .post("".toRequestBody())
                    .build()
                
                client.newCall(request).execute()
            }
            
            deleteToken()
            Result.success(Unit)
        } catch (e: Exception) {
            deleteToken()
            Result.success(Unit)
        }
    }
    
    // MARK: - Token Management
    
    fun saveToken(token: String) {
        securePrefs.edit().putString(TOKEN_KEY, token).apply()
    }
    
    fun getToken(): String? {
        return securePrefs.getString(TOKEN_KEY, null)
    }
    
    fun deleteToken() {
        securePrefs.edit().remove(TOKEN_KEY).apply()
    }
    
    fun hasToken(): Boolean {
        return getToken() != null
    }
    
    data class UserResponse(
        val success: Boolean,
        val data: User
    )
}

class VignetteAuthException(message: String) : Exception(message)

object BuildConfig {
    const val DEBUG = true
}
