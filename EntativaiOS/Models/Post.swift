import Foundation

struct Post: Identifiable, Codable {
    let id: String
    let userId: String
    let userName: String
    let userAvatar: String?
    let text: String?
    let media: [MediaItem]
    let likesCount: Int
    let commentsCount: Int
    let sharesCount: Int
    let timestamp: String
    let createdAt: Date
    
    // Mock data
    static let mockPosts: [Post] = [
        Post(
            id: "1",
            userId: "user1",
            userName: "John Doe",
            userAvatar: nil,
            text: "Just finished an amazing workout session! 💪 Feeling great and ready to tackle the day. Who else is staying active today?",
            media: [MediaItem(type: .image, url: "")],
            likesCount: 124,
            commentsCount: 23,
            sharesCount: 5,
            timestamp: "2h ago",
            createdAt: Date()
        ),
        Post(
            id: "2",
            userId: "user2",
            userName: "Jane Smith",
            userAvatar: nil,
            text: "Beautiful sunset at the beach today 🌅",
            media: [
                MediaItem(type: .image, url: ""),
                MediaItem(type: .image, url: ""),
                MediaItem(type: .image, url: "")
            ],
            likesCount: 489,
            commentsCount: 67,
            sharesCount: 34,
            timestamp: "5h ago",
            createdAt: Date()
        ),
        Post(
            id: "3",
            userId: "user3",
            userName: "Mike Johnson",
            userAvatar: nil,
            text: "New project launch! So excited to share what we've been working on. Stay tuned for more updates!",
            media: [],
            likesCount: 256,
            commentsCount: 45,
            sharesCount: 12,
            timestamp: "1d ago",
            createdAt: Date()
        )
    ]
}

struct MediaItem: Codable {
    let type: MediaType
    let url: String
    let thumbnail: String?
    let duration: Int? // For videos
    
    init(type: MediaType, url: String, thumbnail: String? = nil, duration: Int? = nil) {
        self.type = type
        self.url = url
        self.thumbnail = thumbnail
        self.duration = duration
    }
}

enum MediaType: String, Codable {
    case image
    case video
    case gif
}
