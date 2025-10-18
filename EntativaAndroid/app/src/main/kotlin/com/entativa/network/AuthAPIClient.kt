package com.entativa.network

import android.content.Context
import android.content.SharedPreferences
import android.security.keystore.KeyGenParameterSpec
import android.security.keystore.KeyProperties
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
import java.io.IOException
import java.util.concurrent.TimeUnit

/**
 * Authentication API Client for Entativa Backend
 * Enterprise-grade implementation with secure token storage
 */
class EntativaAuthAPIClient(private val context: Context) {
    
    companion object {
        private const val BASE_URL_DEBUG = "http://10.0.2.2:8001/api/v1"
        private const val BASE_URL_PRODUCTION = "https://api.entativa.com/api/v1"
        private const val TOKEN_KEY = "auth_token"
        
        @Volatile
        private var INSTANCE: EntativaAuthAPIClient? = null
        
        fun getInstance(context: Context): EntativaAuthAPIClient {
            return INSTANCE ?: synchronized(this) {
                INSTANCE ?: EntativaAuthAPIClient(context.applicationContext).also {
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
            "entativa_secure_prefs",
            masterKey,
            EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
            EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
        )
    }
    
    // MARK: - Data Models
    
    data class SignUpRequest(
        @SerializedName("first_name") val firstName: String,
        @SerializedName("last_name") val lastName: String,
        val email: String,
        val password: String,
        val birthday: String,
        val gender: String
    )
    
    data class LoginRequest(
        @SerializedName("email_or_username") val emailOrUsername: String,
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
        @SerializedName("first_name") val firstName: String,
        @SerializedName("last_name") val lastName: String,
        val email: String,
        val username: String,
        val birthday: String?,
        val gender: String?,
        @SerializedName("profile_picture_url") val profilePictureUrl: String?,
        @SerializedName("cover_photo_url") val coverPhotoUrl: String?,
        @SerializedName("is_active") val isActive: Boolean,
        @SerializedName("created_at") val createdAt: String
    )
    
    data class ErrorResponse(
        val success: Boolean,
        val error: String,
        val details: Map<String, List<String>>?
    )
    
    // MARK: - API Methods
    
    /**
     * Sign up a new user
     */
    suspend fun signUp(
        firstName: String,
        lastName: String,
        email: String,
        password: String,
        birthday: String,
        gender: String
    ): Result<AuthResponse> = withContext(Dispatchers.IO) {
        try {
            val requestBody = SignUpRequest(
                firstName = firstName,
                lastName = lastName,
                email = email,
                password = password,
                birthday = birthday,
                gender = gender
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
                Result.failure(AuthException(errorResponse.error))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    
    /**
     * Log in an existing user
     */
    suspend fun login(
        emailOrUsername: String,
        password: String
    ): Result<AuthResponse> = withContext(Dispatchers.IO) {
        try {
            val requestBody = LoginRequest(
                emailOrUsername = emailOrUsername,
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
                Result.failure(AuthException(errorResponse.error))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    
    /**
     * Get current authenticated user
     */
    suspend fun getCurrentUser(): Result<User> = withContext(Dispatchers.IO) {
        try {
            val token = getToken() ?: return@withContext Result.failure(AuthException("Not authenticated"))
            
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
                Result.failure(AuthException("Failed to get user"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    
    /**
     * Log out current user
     */
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
            deleteToken() // Clear token even if request fails
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
    
    // Helper class for user response
    data class UserResponse(
        val success: Boolean,
        val data: User
    )
}

/**
 * Custom exception for authentication errors
 */
class AuthException(message: String) : Exception(message)

// BuildConfig placeholder - should be generated by build system
object BuildConfig {
    const val DEBUG = true
}
