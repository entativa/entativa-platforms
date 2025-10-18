import Foundation
import Combine

/// Authentication API client for Vignette backend
class VignetteAuthAPIClient {
    static let shared = VignetteAuthAPIClient()
    
    private let baseURL: String
    private let session: URLSession
    
    private init() {
        #if DEBUG
        self.baseURL = "http://localhost:8002/api/v1"
        #else
        self.baseURL = "https://api.vignette.app/api/v1"
        #endif
        
        let configuration = URLSessionConfiguration.default
        configuration.timeoutIntervalForRequest = 30
        configuration.timeoutIntervalForResource = 60
        self.session = URLSession(configuration: configuration)
    }
    
    // MARK: - Models
    
    struct SignUpRequest: Codable {
        let username: String
        let email: String
        let fullName: String
        let password: String
        
        enum CodingKeys: String, CodingKey {
            case username
            case email
            case fullName = "full_name"
            case password
        }
    }
    
    struct LoginRequest: Codable {
        let usernameOrEmail: String
        let password: String
        
        enum CodingKeys: String, CodingKey {
            case usernameOrEmail = "username_or_email"
            case password
        }
    }
    
    struct AuthResponse: Codable {
        let success: Bool
        let message: String
        let data: AuthData?
        
        struct AuthData: Codable {
            let user: User
            let accessToken: String
            let tokenType: String
            let expiresIn: Int
            
            enum CodingKeys: String, CodingKey {
                case user
                case accessToken = "access_token"
                case tokenType = "token_type"
                case expiresIn = "expires_in"
            }
        }
        
        struct User: Codable {
            let id: String
            let username: String
            let email: String
            let fullName: String
            let bio: String?
            let website: String?
            let profilePictureUrl: String?
            let isPrivate: Bool
            let isVerified: Bool
            let isActive: Bool
            let followersCount: Int
            let followingCount: Int
            let postsCount: Int
            let createdAt: String
            
            enum CodingKeys: String, CodingKey {
                case id
                case username
                case email
                case fullName = "full_name"
                case bio
                case website
                case profilePictureUrl = "profile_picture_url"
                case isPrivate = "is_private"
                case isVerified = "is_verified"
                case isActive = "is_active"
                case followersCount = "followers_count"
                case followingCount = "following_count"
                case postsCount = "posts_count"
                case createdAt = "created_at"
            }
        }
    }
    
    struct ErrorResponse: Codable {
        let success: Bool
        let error: String
        let details: [String: [String]]?
    }
    
    // MARK: - API Methods
    
    /// Sign up a new user
    func signUp(
        username: String,
        email: String,
        fullName: String,
        password: String
    ) async throws -> AuthResponse {
        let endpoint = "\(baseURL)/auth/signup"
        
        guard let url = URL(string: endpoint) else {
            throw VignetteAuthError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        let signUpRequest = SignUpRequest(
            username: username,
            email: email,
            fullName: fullName,
            password: password
        )
        
        request.httpBody = try JSONEncoder().encode(signUpRequest)
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw VignetteAuthError.invalidResponse
        }
        
        if httpResponse.statusCode == 200 || httpResponse.statusCode == 201 {
            let authResponse = try JSONDecoder().decode(AuthResponse.self, from: data)
            
            // Store token
            if let token = authResponse.data?.accessToken {
                try VignetteKeychainManager.shared.save(token: token)
            }
            
            return authResponse
        } else {
            let errorResponse = try? JSONDecoder().decode(ErrorResponse.self, from: data)
            throw VignetteAuthError.serverError(errorResponse?.error ?? "Sign up failed")
        }
    }
    
    /// Log in an existing user
    func login(usernameOrEmail: String, password: String) async throws -> AuthResponse {
        let endpoint = "\(baseURL)/auth/login"
        
        guard let url = URL(string: endpoint) else {
            throw VignetteAuthError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        let loginRequest = LoginRequest(
            usernameOrEmail: usernameOrEmail,
            password: password
        )
        
        request.httpBody = try JSONEncoder().encode(loginRequest)
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw VignetteAuthError.invalidResponse
        }
        
        if httpResponse.statusCode == 200 {
            let authResponse = try JSONDecoder().decode(AuthResponse.self, from: data)
            
            // Store token
            if let token = authResponse.data?.accessToken {
                try VignetteKeychainManager.shared.save(token: token)
            }
            
            return authResponse
        } else {
            let errorResponse = try? JSONDecoder().decode(ErrorResponse.self, from: data)
            throw VignetteAuthError.serverError(errorResponse?.error ?? "Login failed")
        }
    }
    
    /// Get current authenticated user
    func getCurrentUser() async throws -> AuthResponse.User {
        let endpoint = "\(baseURL)/auth/me"
        
        guard let url = URL(string: endpoint) else {
            throw VignetteAuthError.invalidURL
        }
        
        guard let token = try? VignetteKeychainManager.shared.getToken() else {
            throw VignetteAuthError.notAuthenticated
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw VignetteAuthError.invalidResponse
        }
        
        if httpResponse.statusCode == 200 {
            struct UserResponse: Codable {
                let success: Bool
                let data: AuthResponse.User
            }
            let userResponse = try JSONDecoder().decode(UserResponse.self, from: data)
            return userResponse.data
        } else {
            throw VignetteAuthError.notAuthenticated
        }
    }
    
    /// Log out current user
    func logout() async throws {
        let endpoint = "\(baseURL)/auth/logout"
        
        guard let url = URL(string: endpoint) else {
            throw VignetteAuthError.invalidURL
        }
        
        guard let token = try? VignetteKeychainManager.shared.getToken() else {
            // Already logged out
            return
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        
        let (_, response) = try await session.data(for: request)
        
        // Clear token regardless of server response
        try VignetteKeychainManager.shared.deleteToken()
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw VignetteAuthError.invalidResponse
        }
        
        if httpResponse.statusCode != 200 {
            print("Warning: Logout returned status \(httpResponse.statusCode)")
        }
    }
}

// MARK: - Auth Errors

enum VignetteAuthError: LocalizedError {
    case invalidURL
    case invalidResponse
    case notAuthenticated
    case serverError(String)
    case networkError(Error)
    
    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Invalid API URL"
        case .invalidResponse:
            return "Invalid server response"
        case .notAuthenticated:
            return "Not authenticated. Please log in."
        case .serverError(let message):
            return message
        case .networkError(let error):
            return "Network error: \(error.localizedDescription)"
        }
    }
}

// MARK: - Keychain Manager

class VignetteKeychainManager {
    static let shared = VignetteKeychainManager()
    private let service = "com.vignette.app"
    private let account = "authToken"
    
    private init() {}
    
    func save(token: String) throws {
        let data = token.data(using: .utf8)!
        
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: account,
            kSecValueData as String: data
        ]
        
        // Delete any existing item
        SecItemDelete(query as CFDictionary)
        
        // Add new item
        let status = SecItemAdd(query as CFDictionary, nil)
        
        guard status == errSecSuccess else {
            throw VignetteKeychainError.saveFailed
        }
    }
    
    func getToken() throws -> String {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: account,
            kSecReturnData as String: true
        ]
        
        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)
        
        guard status == errSecSuccess,
              let data = result as? Data,
              let token = String(data: data, encoding: .utf8) else {
            throw VignetteKeychainError.notFound
        }
        
        return token
    }
    
    func deleteToken() throws {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: account
        ]
        
        let status = SecItemDelete(query as CFDictionary)
        
        guard status == errSecSuccess || status == errSecItemNotFound else {
            throw VignetteKeychainError.deleteFailed
        }
    }
}

enum VignetteKeychainError: Error {
    case saveFailed
    case notFound
    case deleteFailed
}
