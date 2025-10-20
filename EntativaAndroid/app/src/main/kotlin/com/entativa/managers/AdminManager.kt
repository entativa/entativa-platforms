package com.entativa.managers

import android.content.Context
import android.util.Log
import androidx.biometric.BiometricManager
import androidx.biometric.BiometricPrompt
import androidx.core.content.ContextCompat
import androidx.fragment.app.FragmentActivity
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody
import org.json.JSONObject
import java.util.concurrent.TimeUnit

/**
 * AdminManager - Founder access control for @neoqiss
 * 
 * Provides secure admin panel access with:
 * - Founder account verification
 * - Triple-tap gesture detection
 * - Biometric authentication
 * - Device whitelisting
 * - Session timeout (15 minutes)
 * - Full audit logging
 */
object AdminManager {
    private const val TAG = "AdminManager"
    private const val ADMIN_BASE_URL = "http://10.0.2.2:8005/api/admin" // Admin service
    private const val SESSION_TIMEOUT_MS = 15 * 60 * 1000L // 15 minutes
    
    private var isAdminMode = false
    private var sessionExpiryTime: Long = 0
    private var sessionTimer: android.os.CountDownTimer? = null
    
    private val client = OkHttpClient.Builder()
        .connectTimeout(30, TimeUnit.SECONDS)
        .readTimeout(30, TimeUnit.SECONDS)
        .writeTimeout(30, TimeUnit.SECONDS)
        .build()
    
    // MARK: - Founder Verification
    
    /**
     * Checks if current user is the founder (@neoqiss)
     */
    fun isFounderAccount(context: Context): Boolean {
        val prefs = context.getSharedPreferences("auth_prefs", Context.MODE_PRIVATE)
        val token = prefs.getString("jwt_token", null) ?: return false
        
        // Decode JWT to check username and is_founder flag
        return checkFounderInToken(token)
    }
    
    private fun checkFounderInToken(token: String): Boolean {
        try {
            val parts = token.split(".")
            if (parts.size != 3) return false
            
            // Decode payload (base64url)
            val payload = parts[1]
            val decodedBytes = android.util.Base64.decode(
                payload.replace('-', '+').replace('_', '/'),
                android.util.Base64.URL_SAFE or android.util.Base64.NO_WRAP
            )
            
            val json = JSONObject(String(decodedBytes))
            val username = json.optString("username", "")
            val isFounder = json.optBoolean("is_founder", false)
            
            return username == "neoqiss" && isFounder
        } catch (e: Exception) {
            Log.e(TAG, "Failed to decode token", e)
            return false
        }
    }
    
    // MARK: - Admin Panel Access
    
    /**
     * Shows admin panel with biometric authentication
     */
    fun showAdminPanel(activity: FragmentActivity) {
        if (!isFounderAccount(activity)) {
            Log.w(TAG, "Not authorized for admin access")
            return
        }
        
        // Require biometric authentication
        authenticateWithBiometric(
            activity = activity,
            title = "Admin Panel Access",
            subtitle = "Authenticate as Founder"
        ) { success ->
            if (success) {
                launchAdminPanel(activity)
            } else {
                // Show error
                Log.e(TAG, "Biometric authentication failed")
            }
        }
    }
    
    private fun launchAdminPanel(activity: FragmentActivity) {
        isAdminMode = true
        startAdminSession()
        
        // Launch admin panel activity
        // TODO: Create AdminPanelActivity
        Log.i(TAG, "Admin panel launched - Founder mode active")
    }
    
    // MARK: - Session Management
    
    private fun startAdminSession() {
        sessionExpiryTime = System.currentTimeMillis() + SESSION_TIMEOUT_MS
        
        // Cancel existing timer
        sessionTimer?.cancel()
        
        // Start new session timer
        sessionTimer = object : android.os.CountDownTimer(SESSION_TIMEOUT_MS, 1000) {
            override fun onTick(millisUntilFinished: Long) {
                // Update session time
            }
            
            override fun onFinish() {
                endAdminSession()
            }
        }.start()
    }
    
    fun endAdminSession() {
        isAdminMode = false
        sessionExpiryTime = 0
        sessionTimer?.cancel()
        sessionTimer = null
        Log.i(TAG, "Admin session ended")
    }
    
    fun extendAdminSession(activity: FragmentActivity) {
        if (!isAdminMode) return
        
        authenticateWithBiometric(
            activity = activity,
            title = "Extend Admin Session",
            subtitle = "Re-authenticate to continue"
        ) { success ->
            if (success) {
                startAdminSession()
            } else {
                endAdminSession()
            }
        }
    }
    
    // MARK: - Quick Admin Actions
    
    /**
     * Ban a user with biometric confirmation
     */
    suspend fun quickBanUser(
        context: Context,
        activity: FragmentActivity,
        userId: String,
        reason: String
    ): Result<Unit> {
        return withContext(Dispatchers.Main) {
            var result: Result<Unit>? = null
            
            authenticateWithBiometric(
                activity = activity,
                title = "Ban User",
                subtitle = "Confirm action"
            ) { success ->
                if (success) {
                    kotlinx.coroutines.GlobalScope.launch(Dispatchers.IO) {
                        result = performBanUser(context, userId, reason)
                    }
                } else {
                    result = Result.failure(Exception("Authentication failed"))
                }
            }
            
            // Wait for result
            while (result == null) {
                kotlinx.coroutines.delay(100)
            }
            
            result!!
        }
    }
    
    private suspend fun performBanUser(
        context: Context,
        userId: String,
        reason: String
    ): Result<Unit> = withContext(Dispatchers.IO) {
        try {
            val prefs = context.getSharedPreferences("auth_prefs", Context.MODE_PRIVATE)
            val token = prefs.getString("jwt_token", null) ?: throw Exception("Not authenticated")
            
            val json = JSONObject().apply {
                put("reason", reason)
                put("duration", 0) // 0 = permanent
            }
            
            val body = json.toString().toRequestBody("application/json".toMediaType())
            
            val request = Request.Builder()
                .url("$ADMIN_BASE_URL/users/$userId/ban")
                .post(body)
                .addHeader("Authorization", "Bearer $token")
                .addHeader("X-Device-ID", getDeviceID(context))
                .build()
            
            val response = client.newCall(request).execute()
            
            if (response.isSuccessful) {
                Result.success(Unit)
            } else {
                Result.failure(Exception("Ban failed: ${response.code}"))
            }
        } catch (e: Exception) {
            Log.e(TAG, "Failed to ban user", e)
            Result.failure(e)
        }
    }
    
    /**
     * Shadowban a user
     */
    suspend fun quickShadowbanUser(
        context: Context,
        activity: FragmentActivity,
        userId: String,
        reason: String
    ): Result<Unit> = withContext(Dispatchers.IO) {
        try {
            val prefs = context.getSharedPreferences("auth_prefs", Context.MODE_PRIVATE)
            val token = prefs.getString("jwt_token", null) ?: throw Exception("Not authenticated")
            
            val json = JSONObject().apply {
                put("reason", reason)
            }
            
            val body = json.toString().toRequestBody("application/json".toMediaType())
            
            val request = Request.Builder()
                .url("$ADMIN_BASE_URL/users/$userId/shadowban")
                .post(body)
                .addHeader("Authorization", "Bearer $token")
                .addHeader("X-Device-ID", getDeviceID(context))
                .build()
            
            val response = client.newCall(request).execute()
            
            if (response.isSuccessful) {
                Result.success(Unit)
            } else {
                Result.failure(Exception("Shadowban failed: ${response.code}"))
            }
        } catch (e: Exception) {
            Log.e(TAG, "Failed to shadowban user", e)
            Result.failure(e)
        }
    }
    
    /**
     * Impersonate a user (step-up authentication with password)
     */
    suspend fun impersonateUser(
        context: Context,
        activity: FragmentActivity,
        userId: String,
        reason: String,
        password: String
    ): Result<String> {
        // Require detailed reason (min 20 chars)
        if (reason.length < 20) {
            return Result.failure(Exception("Reason must be at least 20 characters"))
        }
        
        return withContext(Dispatchers.Main) {
            var result: Result<String>? = null
            
            // First biometric check
            authenticateWithBiometric(
                activity = activity,
                title = "Impersonate User - Step 1",
                subtitle = "First authentication"
            ) { success ->
                if (success) {
                    // Second authentication with password
                    kotlinx.coroutines.GlobalScope.launch(Dispatchers.IO) {
                        result = performImpersonation(context, userId, reason, password)
                    }
                } else {
                    result = Result.failure(Exception("Biometric authentication failed"))
                }
            }
            
            // Wait for result
            while (result == null) {
                kotlinx.coroutines.delay(100)
            }
            
            result!!
        }
    }
    
    private suspend fun performImpersonation(
        context: Context,
        userId: String,
        reason: String,
        password: String
    ): Result<String> = withContext(Dispatchers.IO) {
        try {
            val prefs = context.getSharedPreferences("auth_prefs", Context.MODE_PRIVATE)
            val token = prefs.getString("jwt_token", null) ?: throw Exception("Not authenticated")
            
            val json = JSONObject().apply {
                put("reason", reason)
                put("password", password)
            }
            
            val body = json.toString().toRequestBody("application/json".toMediaType())
            
            val request = Request.Builder()
                .url("$ADMIN_BASE_URL/users/$userId/impersonate")
                .post(body)
                .addHeader("Authorization", "Bearer $token")
                .addHeader("X-Device-ID", getDeviceID(context))
                .build()
            
            val response = client.newCall(request).execute()
            val responseBody = response.body?.string()
            
            if (response.isSuccessful && responseBody != null) {
                val responseJson = JSONObject(responseBody)
                val impersonationToken = responseJson.optString("impersonation_token", "")
                
                if (impersonationToken.isNotEmpty()) {
                    // Start 10-minute auto-termination
                    startImpersonationTimer(context, userId, impersonationToken)
                    Result.success(impersonationToken)
                } else {
                    Result.failure(Exception("No impersonation token returned"))
                }
            } else {
                Result.failure(Exception("Impersonation failed: ${response.code}"))
            }
        } catch (e: Exception) {
            Log.e(TAG, "Failed to impersonate user", e)
            Result.failure(e)
        }
    }
    
    private fun startImpersonationTimer(context: Context, userId: String, token: String) {
        // Auto-terminate after 10 minutes
        object : android.os.CountDownTimer(10 * 60 * 1000L, 1000) {
            override fun onTick(millisUntilFinished: Long) {}
            
            override fun onFinish() {
                kotlinx.coroutines.GlobalScope.launch(Dispatchers.IO) {
                    endImpersonation(context, userId)
                }
            }
        }.start()
    }
    
    private suspend fun endImpersonation(context: Context, userId: String) = withContext(Dispatchers.IO) {
        try {
            val prefs = context.getSharedPreferences("auth_prefs", Context.MODE_PRIVATE)
            val token = prefs.getString("jwt_token", null) ?: return@withContext
            
            val request = Request.Builder()
                .url("$ADMIN_BASE_URL/users/$userId/end-impersonation")
                .post("{}".toRequestBody("application/json".toMediaType()))
                .addHeader("Authorization", "Bearer $token")
                .addHeader("X-Device-ID", getDeviceID(context))
                .build()
            
            client.newCall(request).execute()
            Log.i(TAG, "Impersonation session ended")
        } catch (e: Exception) {
            Log.e(TAG, "Failed to end impersonation", e)
        }
    }
    
    // MARK: - Biometric Authentication
    
    private fun authenticateWithBiometric(
        activity: FragmentActivity,
        title: String,
        subtitle: String,
        callback: (Boolean) -> Unit
    ) {
        val executor = ContextCompat.getMainExecutor(activity)
        
        val biometricPrompt = BiometricPrompt(
            activity,
            executor,
            object : BiometricPrompt.AuthenticationCallback() {
                override fun onAuthenticationSucceeded(result: BiometricPrompt.AuthenticationResult) {
                    super.onAuthenticationSucceeded(result)
                    callback(true)
                }
                
                override fun onAuthenticationFailed() {
                    super.onAuthenticationFailed()
                    callback(false)
                }
                
                override fun onAuthenticationError(errorCode: Int, errString: CharSequence) {
                    super.onAuthenticationError(errorCode, errString)
                    callback(false)
                }
            }
        )
        
        val promptInfo = BiometricPrompt.PromptInfo.Builder()
            .setTitle(title)
            .setSubtitle(subtitle)
            .setNegativeButtonText("Cancel")
            .build()
        
        biometricPrompt.authenticate(promptInfo)
    }
    
    // MARK: - Helpers
    
    private fun getDeviceID(context: Context): String {
        return android.provider.Settings.Secure.getString(
            context.contentResolver,
            android.provider.Settings.Secure.ANDROID_ID
        )
    }
}

/**
 * Multi-tap listener for gesture detection
 */
abstract class OnMultiTapListener(private val requiredTaps: Int) : android.view.View.OnClickListener {
    private var tapCount = 0
    private var lastTapTime = 0L
    private val tapTimeout = 500L // 500ms between taps
    
    override fun onClick(v: android.view.View?) {
        val currentTime = System.currentTimeMillis()
        
        if (currentTime - lastTapTime < tapTimeout) {
            tapCount++
            if (tapCount >= requiredTaps) {
                onMultiTap()
                tapCount = 0
            }
        } else {
            tapCount = 1
        }
        
        lastTapTime = currentTime
    }
    
    abstract fun onMultiTap()
}
