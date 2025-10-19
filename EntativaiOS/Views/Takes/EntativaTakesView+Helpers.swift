import Foundation

// MARK: - Number Formatting Helper
func formatCount(_ count: Int) -> String {
    if count >= 1_000_000 {
        return String(format: "%.1fM", Double(count) / 1_000_000)
    } else if count >= 1_000 {
        return String(format: "%.1fK", Double(count) / 1_000)
    } else {
        return "\(count)"
    }
}

// MARK: - Time Formatting Helper
func formatTimestamp(_ dateString: String) -> String {
    let formatter = ISO8601DateFormatter()
    guard let date = formatter.date(from: dateString) else {
        return "Just now"
    }
    
    let now = Date()
    let components = Calendar.current.dateComponents([.second, .minute, .hour, .day, .weekOfYear], from: date, to: now)
    
    if let weeks = components.weekOfYear, weeks > 0 {
        return "\(weeks)w ago"
    } else if let days = components.day, days > 0 {
        return "\(days)d ago"
    } else if let hours = components.hour, hours > 0 {
        return "\(hours)h ago"
    } else if let minutes = components.minute, minutes > 0 {
        return "\(minutes)m ago"
    } else {
        return "Just now"
    }
}
