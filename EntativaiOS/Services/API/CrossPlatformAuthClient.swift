import Foundation

/// Cross-Platform Authentication Client
/// Enables seamless sign-in between Entativa and Vignette using existing credentials
class CrossPlatformAuthClient {
    static let shared = CrossPlatformAuthClient()
    
    private let session: URLSession
    
    private init() {
        let configuration = URLSessionConfiguration.default
        configuration.timeoutIntervalForRequest = 30
        configuration.timeoutIntervalForResource = 60
        self.session = URLSession(configuration: configuration)
    }
    
    // MARK: - Models
    
    struct CrossPlatformSignInRequest: Codable {
        let platform: String // "vignette" or "entativa"
        let accessToken: String
    }
    
    struct CrossPlatformSignInResponse: Codable {
        let success: Bool
        let message: String
        let data: CrossPlatformData?
        
        struct CrossPlatformData: Codable {
            let user: EntativaUser
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
        
        struct EntativaUser: Codable {
            let id: String
            let firstName: String
            let lastName: String
            let email: String
            let username: String
            let profilePictureUrl: String?
            
            enum CodingKeys: String, CodingKey {
                case id
                case firstName = "first_name"
                case lastName = "last_name"
                case email
                case username
                case profilePictureUrl = "profile_picture_url"
            }
        }
    }
    
    /// Sign in to Entativa using Vignette credentials
    func signInWithVignette(vignetteToken: String) async throws -> CrossPlatformSignInResponse {
        #if DEBUG
        let baseURL = "http://localhost:8001/api/v1"
        #else
        let baseURL = "https://api.entativa.com/api/v1"
        #endif
        
        let endpoint = "\(baseURL)/auth/cross-platform/signin"
        
        guard let url = URL(string: endpoint) else {
            throw CrossPlatformAuthError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        let requestBody = CrossPlatformSignInRequest(
            platform: "vignette",
            accessToken: vignetteToken
        )
        
        request.httpBody = try JSONEncoder().encode(requestBody)
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw CrossPlatformAuthError.invalidResponse
        }
        
        if httpResponse.statusCode == 200 || httpResponse.statusCode == 201 {
            let authResponse = try JSONDecoder().decode(CrossPlatformSignInResponse.self, from: data)
            
            // Store token
            if let token = authResponse.data?.accessToken {
                try KeychainManager.shared.save(token: token)
            }
            
            return authResponse
        } else {
            let errorResponse = try? JSONDecoder().decode(ErrorResponse.self, from: data)
            throw CrossPlatformAuthError.serverError(errorResponse?.error ?? "Sign in failed")
        }
    }
    
    /// Verify if user has an account on the other platform
    func checkCrossPlatformAccount(email: String, targetPlatform: String) async throws -> Bool {
        #if DEBUG
        let baseURL = targetPlatform == "vignette" ? 
            "http://localhost:8002/api/v1" : "http://localhost:8001/api/v1"
        #else
        let baseURL = targetPlatform == "vignette" ? 
            "https://api.vignette.app/api/v1" : "https://api.entativa.com/api/v1"
        #endif
        
        let endpoint = "\(baseURL)/auth/cross-platform/check"
        
        guard var components = URLComponents(string: endpoint) else {
            throw CrossPlatformAuthError.invalidURL
        }
        
        components.queryItems = [URLQueryItem(name: "email", value: email)]
        
        guard let url = components.url else {
            throw CrossPlatformAuthError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw CrossPlatformAuthError.invalidResponse
        }
        
        if httpResponse.statusCode == 200 {
            struct CheckResponse: Codable {
                let exists: Bool
            }
            let checkResponse = try JSONDecoder().decode(CheckResponse.self, from: data)
            return checkResponse.exists
        } else {
            return false
        }
    }
    
    struct ErrorResponse: Codable {
        let error: String
    }
}

enum CrossPlatformAuthError: LocalizedError {
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
