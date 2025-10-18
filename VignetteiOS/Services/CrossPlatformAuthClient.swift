import Foundation

/// Cross-Platform Authentication Client for Vignette
/// Enables seamless sign-in using Entativa credentials
class CrossPlatformAuthClientVignette {
    static let shared = CrossPlatformAuthClientVignette()
    
    private let session: URLSession
    
    private init() {
        let configuration = URLSessionConfiguration.default
        configuration.timeoutIntervalForRequest = 30
        configuration.timeoutIntervalForResource = 60
        self.session = URLSession(configuration: configuration)
    }
    
    // MARK: - Models
    
    struct CrossPlatformSignInRequest: Codable {
        let platform: String // "entativa" or "vignette"
        let accessToken: String
    }
    
    struct CrossPlatformSignInResponse: Codable {
        let success: Bool
        let message: String
        let data: CrossPlatformData?
        
        struct CrossPlatformData: Codable {
            let user: VignetteUser
            let accessToken: String
            let tokenType: String
            let expiresIn: Int
            let isNewAccount: Bool
            
            enum CodingKeys: String, CodingKey {
                case user
                case accessToken = "access_token"
                case tokenType = "token_type"
                case expiresIn = "expires_in"
                case isNewAccount = "is_new_account"
            }
        }
        
        struct VignetteUser: Codable {
            let id: String
            let username: String
            let email: String
            let fullName: String
            let profilePictureUrl: String?
            
            enum CodingKeys: String, CodingKey {
                case id
                case username
                case email
                case fullName = "full_name"
                case profilePictureUrl = "profile_picture_url"
            }
        }
    }
    
    /// Sign in to Vignette using Entativa credentials
    func signInWithEntativa(entativaToken: String) async throws -> CrossPlatformSignInResponse {
        #if DEBUG
        let baseURL = "http://localhost:8002/api/v1"
        #else
        let baseURL = "https://api.vignette.app/api/v1"
        #endif
        
        let endpoint = "\(baseURL)/auth/cross-platform/signin"
        
        guard let url = URL(string: endpoint) else {
            throw CrossPlatformAuthErrorVignette.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        let requestBody = CrossPlatformSignInRequest(
            platform: "entativa",
            accessToken: entativaToken
        )
        
        request.httpBody = try JSONEncoder().encode(requestBody)
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw CrossPlatformAuthErrorVignette.invalidResponse
        }
        
        if httpResponse.statusCode == 200 || httpResponse.statusCode == 201 {
            let authResponse = try JSONDecoder().decode(CrossPlatformSignInResponse.self, from: data)
            
            // Store token
            if let token = authResponse.data?.accessToken {
                try VignetteKeychainManager.shared.save(token: token)
            }
            
            return authResponse
        } else {
            let errorResponse = try? JSONDecoder().decode(ErrorResponse.self, from: data)
            throw CrossPlatformAuthErrorVignette.serverError(errorResponse?.error ?? "Sign in failed")
        }
    }
    
    struct ErrorResponse: Codable {
        let error: String
    }
}

enum CrossPlatformAuthErrorVignette: LocalizedError {
    case invalidURL
    case invalidResponse
    case serverError(String)
    case tokenExpired
    case accountNotFound
    
    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Invalid API URL"
        case .invalidResponse:
            return "Invalid server response"
        case .serverError(let message):
            return message
        case .tokenExpired:
            return "Your session has expired. Please sign in again."
        case .accountNotFound:
            return "No account found on the other platform."
        }
    }
}
