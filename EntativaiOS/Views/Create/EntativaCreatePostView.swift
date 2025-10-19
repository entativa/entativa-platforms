import SwiftUI
import PhotosUI

// MARK: - Entativa Create Post View (Facebook-Style)
struct EntativaCreatePostView: View {
    @Environment(\.dismiss) var dismiss
    @StateObject private var viewModel = EntativaCreatePostViewModel()
    @State private var showMediaPicker = false
    @State private var showBackgroundPicker = false
    
    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 0) {
                    // User info
                    HStack(spacing: 12) {
                        Circle()
                            .fill(Color.gray.opacity(0.2))
                            .frame(width: 40, height: 40)
                            .overlay(
                                Image(systemName: "person.fill")
                                    .foregroundColor(.gray)
                            )
                        
                        VStack(alignment: .leading, spacing: 4) {
                            Text("Your Name")
                                .font(.system(size: 15, weight: .semibold))
                            
                            // Audience selector
                            Button(action: {}) {
                                HStack(spacing: 4) {
                                    Image(systemName: viewModel.audienceIcon)
                                        .font(.system(size: 12))
                                    Text(viewModel.audience.rawValue)
                                        .font(.system(size: 13))
                                    Image(systemName: "chevron.down")
                                        .font(.system(size: 10))
                                }
                                .foregroundColor(.gray)
                                .padding(.horizontal, 8)
                                .padding(.vertical, 4)
                                .background(Color.gray.opacity(0.1))
                                .cornerRadius(4)
                            }
                        }
                        
                        Spacer()
                    }
                    .padding(16)
                    
                    // Text input
                    if viewModel.selectedBackground == nil {
                        TextEditor(text: $viewModel.postText)
                            .font(.system(size: 18))
                            .frame(minHeight: 150)
                            .padding(.horizontal, 16)
                            .placeholder(when: viewModel.postText.isEmpty) {
                                Text("What's on your mind?")
                                    .font(.system(size: 18))
                                    .foregroundColor(.gray)
                                    .padding(.horizontal, 20)
                                    .padding(.top, 8)
                            }
                    } else {
                        // Background text post
                        ZStack {
                            viewModel.selectedBackground
                                .resizable()
                                .aspectRatio(contentMode: .fill)
                                .frame(height: 300)
                                .clipped()
                            
                            TextEditor(text: $viewModel.postText)
                                .font(.system(size: 24, weight: .bold))
                                .foregroundColor(.white)
                                .multilineTextAlignment(.center)
                                .padding(32)
                                .background(Color.clear)
                                .frame(height: 300)
                        }
                    }
                    
                    // Selected media preview
                    if !viewModel.selectedImages.isEmpty {
                        ScrollView(.horizontal, showsIndicators: false) {
                            HStack(spacing: 8) {
                                ForEach(viewModel.selectedImages.indices, id: \.self) { index in
                                    ZStack(alignment: .topTrailing) {
                                        Image(uiImage: viewModel.selectedImages[index])
                                            .resizable()
                                            .aspectRatio(contentMode: .fill)
                                            .frame(width: 120, height: 120)
                                            .clipped()
                                            .cornerRadius(8)
                                        
                                        Button(action: {
                                            viewModel.selectedImages.remove(at: index)
                                        }) {
                                            Image(systemName: "xmark.circle.fill")
                                                .font(.system(size: 20))
                                                .foregroundColor(.white)
                                                .background(
                                                    Circle()
                                                        .fill(Color.black.opacity(0.5))
                                                        .frame(width: 24, height: 24)
                                                )
                                        }
                                        .padding(8)
                                    }
                                }
                            }
                            .padding(.horizontal, 16)
                        }
                        .padding(.vertical, 12)
                    }
                    
                    Divider()
                        .padding(.horizontal, 16)
                    
                    // Action buttons
                    VStack(spacing: 0) {
                        PostActionButton(
                            icon: "photo.on.rectangle",
                            iconColor: Color.green,
                            title: "Photo/Video",
                            action: { showMediaPicker = true }
                        )
                        
                        PostActionButton(
                            icon: "person.2",
                            iconColor: Color(hex: "007CFC"),
                            title: "Tag people",
                            action: {}
                        )
                        
                        PostActionButton(
                            icon: "face.smiling",
                            iconColor: Color.yellow,
                            title: "Feeling/Activity",
                            action: {}
                        )
                        
                        PostActionButton(
                            icon: "location",
                            iconColor: Color.red,
                            title: "Check in",
                            action: {}
                        )
                        
                        PostActionButton(
                            icon: "camera.fill",
                            iconColor: Color(hex: "6F3EFB"),
                            title: "Live video",
                            action: {}
                        )
                        
                        PostActionButton(
                            icon: "paintpalette",
                            iconColor: Color.orange,
                            title: "Background",
                            action: { showBackgroundPicker = true }
                        )
                        
                        PostActionButton(
                            icon: "gift",
                            iconColor: Color.pink,
                            title: "Celebration",
                            action: {}
                        )
                    }
                    .padding(.vertical, 8)
                }
            }
            .navigationTitle("Create Post")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Cancel") {
                        dismiss()
                    }
                }
                
                ToolbarItem(placement: .confirmationAction) {
                    Button("Post") {
                        viewModel.createPost()
                        dismiss()
                    }
                    .fontWeight(.semibold)
                    .disabled(viewModel.postText.isEmpty && viewModel.selectedImages.isEmpty)
                }
            }
            .sheet(isPresented: $showMediaPicker) {
                MediaPickerSheet(selectedImages: $viewModel.selectedImages)
            }
            .sheet(isPresented: $showBackgroundPicker) {
                BackgroundPickerSheet(selectedBackground: $viewModel.selectedBackground)
            }
        }
    }
}

// MARK: - Post Action Button
struct PostActionButton: View {
    let icon: String
    let iconColor: Color
    let title: String
    let action: () -> Void
    
    var body: some View {
        Button(action: action) {
            HStack(spacing: 12) {
                Image(systemName: icon)
                    .font(.system(size: 20))
                    .foregroundColor(iconColor)
                    .frame(width: 24)
                
                Text(title)
                    .font(.system(size: 15))
                    .foregroundColor(.primary)
                
                Spacer()
                
                Image(systemName: "chevron.right")
                    .font(.system(size: 14))
                    .foregroundColor(.gray)
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)
        }
    }
}

// MARK: - Media Picker Sheet
struct MediaPickerSheet: View {
    @Environment(\.dismiss) var dismiss
    @Binding var selectedImages: [UIImage]
    @State private var selectedItems: [PhotosPickerItem] = []
    
    var body: some View {
        NavigationView {
            PhotosPicker(
                selection: $selectedItems,
                maxSelectionCount: 10,
                matching: .images
            ) {
                VStack(spacing: 16) {
                    Image(systemName: "photo.on.rectangle.angled")
                        .font(.system(size: 64))
                        .foregroundColor(Color(hex: "007CFC"))
                    
                    Text("Select photos or videos")
                        .font(.system(size: 18, weight: .semibold))
                    
                    Text("You can select up to 10 items")
                        .font(.system(size: 14))
                        .foregroundColor(.gray)
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
            }
            .onChange(of: selectedItems) { _, newValue in
                loadImages(from: newValue)
            }
            .navigationTitle("Select Media")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Cancel") {
                        dismiss()
                    }
                }
                
                ToolbarItem(placement: .confirmationAction) {
                    Button("Done") {
                        dismiss()
                    }
                }
            }
        }
    }
    
    func loadImages(from items: [PhotosPickerItem]) {
        Task {
            var loadedImages: [UIImage] = []
            
            for item in items {
                if let data = try? await item.loadTransferable(type: Data.self),
                   let image = UIImage(data: data) {
                    loadedImages.append(image)
                }
            }
            
            await MainActor.run {
                self.selectedImages = loadedImages
                dismiss()
            }
        }
    }
}

// MARK: - Background Picker Sheet
struct BackgroundPickerSheet: View {
    @Environment(\.dismiss) var dismiss
    @Binding var selectedBackground: Image?
    
    let backgrounds: [Color] = [
        Color(hex: "007CFC"),
        Color(hex: "6F3EFB"),
        Color(hex: "FC30E1"),
        Color.red,
        Color.orange,
        Color.green,
        Color.purple,
        Color.pink
    ]
    
    var body: some View {
        NavigationView {
            ScrollView {
                LazyVGrid(columns: [
                    GridItem(.flexible()),
                    GridItem(.flexible())
                ], spacing: 16) {
                    ForEach(backgrounds, id: \.self) { color in
                        Button(action: {
                            selectedBackground = Image(systemName: "square.fill")
                            dismiss()
                        }) {
                            RoundedRectangle(cornerRadius: 12)
                                .fill(color)
                                .frame(height: 200)
                                .overlay(
                                    Text("Aa")
                                        .font(.system(size: 48, weight: .bold))
                                        .foregroundColor(.white)
                                )
                        }
                    }
                }
                .padding(16)
            }
            .navigationTitle("Choose Background")
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
}

// MARK: - Text Editor Placeholder Extension
extension View {
    func placeholder<Content: View>(
        when shouldShow: Bool,
        alignment: Alignment = .leading,
        @ViewBuilder placeholder: () -> Content
    ) -> some View {
        ZStack(alignment: alignment) {
            placeholder().opacity(shouldShow ? 1 : 0)
            self
        }
    }
}

// MARK: - Audience Enum
enum PostAudience: String {
    case public = "Public"
    case friends = "Friends"
    case onlyMe = "Only Me"
    case custom = "Custom"
}

// MARK: - View Model
class EntativaCreatePostViewModel: ObservableObject {
    @Published var postText: String = ""
    @Published var selectedImages: [UIImage] = []
    @Published var selectedBackground: Image?
    @Published var audience: PostAudience = .public
    @Published var taggedPeople: [String] = []
    @Published var feeling: String = ""
    @Published var location: String = ""
    @Published var isLoading = false
    
    var audienceIcon: String {
        switch audience {
        case .public:
            return "globe"
        case .friends:
            return "person.2"
        case .onlyMe:
            return "lock"
        case .custom:
            return "gearshape"
        }
    }
    
    func createPost() {
        // TODO: Upload media and create post via API
        isLoading = true
        
        Task {
            // Simulate upload
            try? await Task.sleep(nanoseconds: 2_000_000_000)
            
            await MainActor.run {
                isLoading = false
            }
        }
    }
}
