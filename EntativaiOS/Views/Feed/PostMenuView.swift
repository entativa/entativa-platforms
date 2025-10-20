import SwiftUI

// MARK: - Post Menu (3-dot button on posts)
struct PostMenuButton: View {
    let post: Post
    let isOwnPost: Bool
    @State private var showMenu = false
    
    var body: some View {
        Button(action: { showMenu = true }) {
            Image(systemName: "ellipsis")
                .font(.system(size: 18, weight: .semibold))
                .foregroundColor(.primary)
                .frame(width: 32, height: 32)
                .background(Color.gray.opacity(0.1))
                .clipShape(Circle())
        }
        .confirmationDialog("Post Options", isPresented: $showMenu, titleVisibility: .hidden) {
            if isOwnPost {
                // Own post options
                Button("Edit Post") {
                    editPost()
                }
                
                Button("Delete Post", role: .destructive) {
                    deletePost()
                }
                
                Button("Pin to Profile") {
                    pinPost()
                }
                
                Button("Archive Post") {
                    archivePost()
                }
                
                Button(post.commentsDisabled ? "Enable Comments" : "Turn Off Comments") {
                    toggleComments()
                }
                
                Button("Share") {
                    sharePost()
                }
                
                Button("Copy Link") {
                    copyLink()
                }
            } else {
                // Others' post options
                Button("Report Post") {
                    reportPost()
                }
                
                Button("Block \(post.author.username)") {
                    blockUser()
                }
                
                Button("Hide Post") {
                    hidePost()
                }
                
                Button("Not Interested") {
                    notInterested()
                }
                
                Button("Share") {
                    sharePost()
                }
                
                Button("Copy Link") {
                    copyLink()
                }
            }
            
            Button("Cancel", role: .cancel) {}
        }
    }
    
    // MARK: - Own Post Actions
    
    private func editPost() {
        // Navigate to edit post screen
        print("Edit post: \(post.id)")
    }
    
    private func deletePost() {
        // Show confirmation alert, then delete
        print("Delete post: \(post.id)")
    }
    
    private func pinPost() {
        // Pin post to profile
        print("Pin post: \(post.id)")
    }
    
    private func archivePost() {
        // Archive post
        print("Archive post: \(post.id)")
    }
    
    private func toggleComments() {
        // Toggle comments on/off
        print("Toggle comments: \(post.id)")
    }
    
    // MARK: - Others' Post Actions
    
    private func reportPost() {
        // Show report options
        print("Report post: \(post.id)")
    }
    
    private func blockUser() {
        // Show block confirmation
        print("Block user: \(post.author.id)")
    }
    
    private func hidePost() {
        // Hide post from feed
        print("Hide post: \(post.id)")
    }
    
    private func notInterested() {
        // Mark as not interested (affects algorithm)
        print("Not interested: \(post.id)")
    }
    
    // MARK: - Common Actions
    
    private func sharePost() {
        // Show share sheet
        guard let url = URL(string: "https://entativa.com/p/\(post.id)") else { return }
        let activityVC = UIActivityViewController(activityItems: [url], applicationActivities: nil)
        
        if let windowScene = UIApplication.shared.connectedScenes.first as? UIWindowScene,
           let window = windowScene.windows.first,
           let rootVC = window.rootViewController {
            rootVC.present(activityVC, animated: true)
        }
    }
    
    private func copyLink() {
        // Copy post link to clipboard
        let link = "https://entativa.com/p/\(post.id)"
        UIPasteboard.general.string = link
        
        // Show toast
        print("Link copied: \(link)")
    }
}

// MARK: - Post Model Extension
extension Post {
    var commentsDisabled: Bool {
        // This would come from the API
        return false
    }
}

// MARK: - Report Post Sheet
struct ReportPostSheet: View {
    @Environment(\.dismiss) var dismiss
    let post: Post
    @State private var selectedReason: ReportReason?
    @State private var additionalInfo = ""
    
    var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Header
                VStack(spacing: 8) {
                    Text("Report Post")
                        .font(.system(size: 20, weight: .bold))
                    
                    Text("Help us understand what's happening")
                        .font(.system(size: 14))
                        .foregroundColor(.gray)
                }
                .padding(.top, 24)
                .padding(.bottom, 16)
                
                Divider()
                
                // Reasons
                ScrollView {
                    VStack(spacing: 0) {
                        ForEach(ReportReason.allCases, id: \.self) { reason in
                            Button(action: { selectedReason = reason }) {
                                HStack {
                                    VStack(alignment: .leading, spacing: 4) {
                                        Text(reason.title)
                                            .font(.system(size: 16, weight: .medium))
                                            .foregroundColor(.primary)
                                        
                                        Text(reason.subtitle)
                                            .font(.system(size: 13))
                                            .foregroundColor(.gray)
                                    }
                                    
                                    Spacer()
                                    
                                    if selectedReason == reason {
                                        Image(systemName: "checkmark.circle.fill")
                                            .foregroundColor(Color(hex: "007CFC"))
                                    } else {
                                        Image(systemName: "circle")
                                            .foregroundColor(.gray)
                                    }
                                }
                                .padding(.horizontal, 20)
                                .padding(.vertical, 16)
                            }
                            
                            if reason != ReportReason.allCases.last {
                                Divider()
                                    .padding(.leading, 20)
                            }
                        }
                    }
                }
                
                // Additional Info
                if selectedReason != nil {
                    Divider()
                    
                    VStack(alignment: .leading, spacing: 8) {
                        Text("Additional Information (Optional)")
                            .font(.system(size: 13))
                            .foregroundColor(.gray)
                        
                        TextEditor(text: $additionalInfo)
                            .frame(height: 100)
                            .padding(8)
                            .background(Color.gray.opacity(0.1))
                            .cornerRadius(8)
                    }
                    .padding(20)
                }
                
                // Submit Button
                Button(action: submitReport) {
                    Text("Submit Report")
                        .font(.system(size: 17, weight: .semibold))
                        .foregroundColor(.white)
                        .frame(maxWidth: .infinity)
                        .frame(height: 50)
                        .background(selectedReason != nil ? Color.red : Color.gray)
                        .cornerRadius(12)
                }
                .disabled(selectedReason == nil)
                .padding(.horizontal, 20)
                .padding(.bottom, 32)
            }
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Cancel") {
                        dismiss()
                    }
                }
            }
        }
    }
    
    private func submitReport() {
        // TODO: Submit report to API
        dismiss()
    }
}

// MARK: - Report Reasons
enum ReportReason: CaseIterable {
    case spam
    case inappropriate
    case harassment
    case hateSpeech
    case violence
    case falseInfo
    case scam
    case other
    
    var title: String {
        switch self {
        case .spam: return "Spam"
        case .inappropriate: return "Inappropriate Content"
        case .harassment: return "Harassment or Bullying"
        case .hateSpeech: return "Hate Speech"
        case .violence: return "Violence or Dangerous Organizations"
        case .falseInfo: return "False Information"
        case .scam: return "Scam or Fraud"
        case .other: return "Something Else"
        }
    }
    
    var subtitle: String {
        switch self {
        case .spam: return "Unwanted commercial content or spam"
        case .inappropriate: return "Nudity, sexual content, or graphic violence"
        case .harassment: return "Bullying or harassment of others"
        case .hateSpeech: return "Hateful or discriminatory content"
        case .violence: return "Content promoting violence or terrorism"
        case .falseInfo: return "Misleading or false information"
        case .scam: return "Fraudulent content or scams"
        case .other: return "Something not listed here"
        }
    }
}
