import Foundation

class TakesAPIClient {
    static let shared = TakesAPIClient()
    
    private let baseURL: String
    private let session: URLSession
    
    private init() {
        #if DEBUG
        self.baseURL = "http://localhost:8001/api/v1"
        #else
        self.baseURL = "https://api.entativa.com/api/v1"
        #endif
        
        let configuration = URLSessionConfiguration.default
        configuration.timeoutIntervalForRequest = 30
        self.session = URLSession(configuration: configuration)
    }
    
    // MARK: - Get Takes Feed
    func getFeed(page: Int = 1, limit: Int = 10) async throws -> TakesFeedResponse {
        guard let url = URL(string: "\(baseURL)/takes/feed?page=\(page)&limit=\(limit)") else {
            throw TakesError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        
        // Add auth token if available
        if let token = try? KeychainManager.shared.retrieve(key: "authToken") {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw TakesError.invalidResponse
        }
        
        guard httpResponse.statusCode == 200 else {
            throw TakesError.serverError(statusCode: httpResponse.statusCode)
        }
        
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        
        let apiResponse = try decoder.decode(TakesAPIResponse.self, from: data)
        
        guard apiResponse.success, let feedData = apiResponse.data else {
            throw TakesError.decodingError
        }
        
        return feedData
    }
    
    // MARK: - Like Take
    func likeTake(takeID: String) async throws -> TakeModel {
        guard let url = URL(string: "\(baseURL)/takes/\(takeID)/like") else {
            throw TakesError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        // Add auth token
        guard let token = try? KeychainManager.shared.retrieve(key: "authToken") else {
            throw TakesError.unauthorized
        }
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw TakesError.invalidResponse
        }
        
        guard httpResponse.statusCode == 200 else {
            throw TakesError.serverError(statusCode: httpResponse.statusCode)
        }
        
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        
        let apiResponse = try decoder.decode(TakeSingleAPIResponse.self, from: data)
        
        guard apiResponse.success, let take = apiResponse.data else {
            throw TakesError.decodingError
        }
        
        return take
    }
    
    // MARK: - Unlike Take
    func unlikeTake(takeID: String) async throws -> TakeModel {
        guard let url = URL(string: "\(baseURL)/takes/\(takeID)/unlike") else {
            throw TakesError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        guard let token = try? KeychainManager.shared.retrieve(key: "authToken") else {
            throw TakesError.unauthorized
        }
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw TakesError.invalidResponse
        }
        
        guard httpResponse.statusCode == 200 else {
            throw TakesError.serverError(statusCode: httpResponse.statusCode)
        }
        
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        
        let apiResponse = try decoder.decode(TakeSingleAPIResponse.self, from: data)
        
        guard apiResponse.success, let take = apiResponse.data else {
            throw TakesError.decodingError
        }
        
        return take
    }
    
    // MARK: - Get Comments
    func getComments(takeID: String, page: Int = 1) async throws -> CommentsResponse {
        guard let url = URL(string: "\(baseURL)/takes/\(takeID)/comments?page=\(page)") else {
            throw TakesError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw TakesError.invalidResponse
        }
        
        guard httpResponse.statusCode == 200 else {
            throw TakesError.serverError(statusCode: httpResponse.statusCode)
        }
        
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        
        let apiResponse = try decoder.decode(CommentsAPIResponse.self, from: data)
        
        guard apiResponse.success, let commentsData = apiResponse.data else {
            throw TakesError.decodingError
        }
        
        return commentsData
    }
    
    // MARK: - Add Comment
    func addComment(takeID: String, text: String) async throws {
        guard let url = URL(string: "\(baseURL)/takes/\(takeID)/comments") else {
            throw TakesError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        guard let token = try? KeychainManager.shared.retrieve(key: "authToken") else {
            throw TakesError.unauthorized
        }
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        
        let body = ["text": text]
        request.httpBody = try JSONEncoder().encode(body)
        
        let (_, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse else {
            throw TakesError.invalidResponse
        }
        
        guard httpResponse.statusCode == 201 else {
            throw TakesError.serverError(statusCode: httpResponse.statusCode)
        }
    }
}

// MARK: - Models
struct TakeModel: Codable, Identifiable {
    let id: String
    let userId: String
    let username: String
    let userAvatar: String?
    let videoUrl: String
    let thumbnailUrl: String?
    let caption: String
    let audioName: String
    let audioUrl: String?
    let duration: Int
    let likesCount: Int
    let commentsCount: Int
    let sharesCount: Int
    let viewsCount: Int
    let isLiked: Bool
    let isSaved: Bool
    let hashtags: [String]?
    let createdAt: String
}

struct TakesFeedResponse: Codable {
    let takes: [TakeModel]
    let page: Int
    let limit: Int
    let hasMore: Bool
}

struct CommentsResponse: Codable {
    let comments: [TakeComment]
    let page: Int
    let limit: Int
    let hasMore: Bool
}

struct TakeComment: Codable, Identifiable {
    let id: String
    let takeId: String
    let userId: String
    let username: String
    let userAvatar: String?
    let text: String
    let likesCount: Int
    let isLiked: Bool
    let createdAt: String
}

// MARK: - API Response Wrappers
struct TakesAPIResponse: Codable {
    let success: Bool
    let message: String?
    let data: TakesFeedResponse?
}

struct TakeSingleAPIResponse: Codable {
    let success: Bool
    let message: String?
    let data: TakeModel?
}

struct CommentsAPIResponse: Codable {
    let success: Bool
    let message: String?
    let data: CommentsResponse?
}

// MARK: - Errors
enum TakesError: LocalizedError {
    case invalidURL
    case invalidResponse
    case unauthorized
    case serverError(statusCode: Int)
    case decodingError
    case networkError(Error)
    
    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Invalid URL"
        case .invalidResponse:
            return "Invalid response from server"
        case .unauthorized:
            return "You must be logged in to perform this action"
        case .serverError(let statusCode):
            return "Server error: \(statusCode)"
        case .decodingError:
            return "Failed to decode response"
        case .networkError(let error):
            return "Network error: \(error.localizedDescription)"
        }
    }
}
