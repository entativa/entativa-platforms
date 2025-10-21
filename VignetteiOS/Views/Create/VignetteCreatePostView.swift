import SwiftUI
import PhotosUI

// MARK: - Vignette Create Post View (Instagram-Style)
struct VignetteCreatePostView: View {
    @Environment(\.dismiss) var dismiss
    @StateObject private var viewModel = CreatePostViewModel()
    @State private var selectedMedia: [PhotosPickerItem] = []
    @State private var selectedTab: MediaTab = .camera
    
    var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Media selector/preview
                if viewModel.selectedImages.isEmpty {
                    MediaSelectorView(
                        selectedTab: $selectedTab,
                        selectedMedia: $selectedMedia,
                        onMediaSelected: { items in
                            viewModel.loadImages(from: items)
                        }
                    )
                } else {
                    // Preview selected media
                    TabView {
                        ForEach(viewModel.selectedImages.indices, id: \.self) { index in
                            Image(uiImage: viewModel.selectedImages[index])
                                .resizable()
                                .aspectRatio(contentMode: .fill)
                                .frame(maxWidth: .infinity)
                                .clipped()
                        }
                    }
                    .tabViewStyle(.page(indexDisplayMode: .always))
                    .frame(height: 400)
                    
                    // Edit tools
                    ScrollView(.horizontal, showsIndicators: false) {
                        HStack(spacing: 12) {
                            EditToolButton(icon: "wand.and.stars", title: "Filter")
                            EditToolButton(icon: "crop", title: "Crop")
                            EditToolButton(icon: "slider.horizontal.3", title: "Adjust")
                            EditToolButton(icon: "textformat", title: "Text")
                            EditToolButton(icon: "scribble", title: "Draw")
                        }
                        .padding(.horizontal, 16)
                    }
                    .padding(.vertical, 12)
                    .background(Color(UIColor.systemBackground))
                }
                
                Divider()
                
                // Caption and details
                ScrollView {
                    VStack(spacing: 16) {
                        // Caption
                        TextField("Write a caption...", text: $viewModel.caption, axis: .vertical)
                            .font(.system(size: 16))
                            .lineLimit(5...10)
                            .padding(.horizontal, 16)
                            .padding(.vertical, 12)
                        
                        Divider()
                        
                        // Tag people
                        NavigationLink(destination: Text("Tag People")) {
                            HStack {
                                Text("Tag people")
                                    .font(.system(size: 16))
                                Spacer()
                                Image(systemName: "chevron.right")
                                    .font(.system(size: 14))
                                    .foregroundColor(.gray)
                            }
                            .padding(.horizontal, 16)
                            .padding(.vertical, 12)
                        }
                        .foregroundColor(.primary)
                        
                        Divider()
                        
                        // Add location
                        NavigationLink(destination: Text("Add Location")) {
                            HStack {
                                Text("Add location")
                                    .font(.system(size: 16))
                                Spacer()
                                Image(systemName: "chevron.right")
                                    .font(.system(size: 14))
                                    .foregroundColor(.gray)
                            }
                            .padding(.horizontal, 16)
                            .padding(.vertical, 12)
                        }
                        .foregroundColor(.primary)
                        
                        Divider()
                        
                        // Also post to
                        VStack(alignment: .leading, spacing: 12) {
                            Text("Also post to")
                                .font(.system(size: 13, weight: .semibold))
                                .foregroundColor(.gray)
                                .padding(.horizontal, 16)
                            
                            Toggle("Entativa", isOn: $viewModel.postToEntativa)
                                .padding(.horizontal, 16)
                            
                            Toggle("Twitter", isOn: $viewModel.postToTwitter)
                                .padding(.horizontal, 16)
                            
                            Toggle("Tumblr", isOn: $viewModel.postToTumblr)
                                .padding(.horizontal, 16)
                        }
                        .padding(.vertical, 8)
                        
                        Divider()
                        
                        // Advanced settings
                        Button(action: {}) {
                            HStack {
                                Text("Advanced settings")
                                    .font(.system(size: 16))
                                Spacer()
                                Image(systemName: "chevron.right")
                                    .font(.system(size: 14))
                                    .foregroundColor(.gray)
                            }
                            .padding(.horizontal, 16)
                            .padding(.vertical, 12)
                        }
                        .foregroundColor(.primary)
                    }
                }
            }
            .navigationTitle("New Post")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Cancel") {
                        dismiss()
                    }
                }
                
                ToolbarItem(placement: .confirmationAction) {
                    Button("Share") {
                        viewModel.createPost()
                        dismiss()
                    }
                    .fontWeight(.semibold)
                    .disabled(viewModel.selectedImages.isEmpty)
                }
            }
        }
    }
}

// MARK: - Media Selector View
struct MediaSelectorView: View {
    @Binding var selectedTab: MediaTab
    @Binding var selectedMedia: [PhotosPickerItem]
    let onMediaSelected: ([PhotosPickerItem]) -> Void
    
    var body: some View {
        VStack(spacing: 0) {
            // Tab selector
            HStack(spacing: 0) {
                MediaTabButton(
                    title: "Camera",
                    icon: "camera",
                    isSelected: selectedTab == .camera,
                    action: { selectedTab = .camera }
                )
                
                MediaTabButton(
                    title: "Gallery",
                    icon: "photo.on.rectangle",
                    isSelected: selectedTab == .gallery,
                    action: { selectedTab = .gallery }
                )
            }
            .frame(height: 44)
            .background(Color(UIColor.systemBackground))
            
            Divider()
            
            // Content
            if selectedTab == .camera {
                CameraPreviewView()
            } else {
                PhotosPickerView(
                    selectedMedia: $selectedMedia,
                    onMediaSelected: onMediaSelected
                )
            }
        }
    }
}

// MARK: - Camera Preview View
struct CameraPreviewView: View {
    var body: some View {
        ZStack {
            Color.black
            
            VStack {
                Spacer()
                
                HStack(spacing: 40) {
                    // Gallery button
                    Button(action: {}) {
                        RoundedRectangle(cornerRadius: 8)
                            .fill(Color.white.opacity(0.3))
                            .frame(width: 60, height: 60)
                            .overlay(
                                Image(systemName: "photo.on.rectangle")
                                    .font(.system(size: 24))
                                    .foregroundColor(.white)
                            )
                    }
                    
                    // Capture button
                    Button(action: {}) {
                        Circle()
                            .stroke(Color.white, lineWidth: 4)
                            .frame(width: 80, height: 80)
                            .overlay(
                                Circle()
                                    .fill(Color.white)
                                    .frame(width: 68, height: 68)
                            )
                    }
                    
                    // Flip camera button
                    Button(action: {}) {
                        Circle()
                            .fill(Color.white.opacity(0.3))
                            .frame(width: 60, height: 60)
                            .overlay(
                                Image(systemName: "arrow.triangle.2.circlepath.camera")
                                    .font(.system(size: 24))
                                    .foregroundColor(.white)
                            )
                    }
                }
                .padding(.bottom, 40)
            }
        }
        .frame(height: 400)
    }
}

// MARK: - Photos Picker View
struct PhotosPickerView: View {
    @Binding var selectedMedia: [PhotosPickerItem]
    let onMediaSelected: ([PhotosPickerItem]) -> Void
    
    var body: some View {
        PhotosPicker(
            selection: $selectedMedia,
            maxSelectionCount: 10,
            matching: .images
        ) {
            VStack(spacing: 16) {
                Image(systemName: "photo.on.rectangle.angled")
                    .font(.system(size: 64))
                    .foregroundColor(Color(hex: "007CFC"))
                
                Text("Select photos")
                    .font(.system(size: 18, weight: .semibold))
                
                Text("You can select up to 10 photos")
                    .font(.system(size: 14))
                    .foregroundColor(.gray)
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
        }
        .onChange(of: selectedMedia) { _, newValue in
            onMediaSelected(newValue)
        }
    }
}

// MARK: - Media Tab Button
struct MediaTabButton: View {
    let title: String
    let icon: String
    let isSelected: Bool
    let action: () -> Void
    
    var body: some View {
        Button(action: action) {
            VStack(spacing: 8) {
                HStack(spacing: 6) {
                    Image(systemName: icon)
                        .font(.system(size: 16))
                    Text(title)
                        .font(.system(size: 16, weight: isSelected ? .semibold : .regular))
                }
                .foregroundColor(isSelected ? .primary : .gray)
                .frame(maxWidth: .infinity)
                
                Rectangle()
                    .fill(isSelected ? Color.primary : Color.clear)
                    .frame(height: 2)
            }
        }
    }
}

// MARK: - Edit Tool Button
struct EditToolButton: View {
    let icon: String
    let title: String
    
    var body: some View {
        Button(action: {}) {
            VStack(spacing: 8) {
                Image(systemName: icon)
                    .font(.system(size: 24))
                    .foregroundColor(.primary)
                    .frame(width: 50, height: 50)
                    .background(Color.gray.opacity(0.1))
                    .clipShape(Circle())
                
                Text(title)
                    .font(.system(size: 12))
                    .foregroundColor(.primary)
            }
        }
    }
}

// MARK: - Media Tab Enum
enum MediaTab {
    case camera
    case gallery
}

// MARK: - View Model
class CreatePostViewModel: ObservableObject {
    @Published var selectedImages: [UIImage] = []
    @Published var caption: String = ""
    @Published var location: String = ""
    @Published var taggedPeople: [String] = []
    @Published var postToEntativa: Bool = false
    @Published var postToTwitter: Bool = false
    @Published var postToTumblr: Bool = false
    @Published var isLoading = false
    
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
            }
        }
    }
    
    func createPost() {
        // TODO: Upload images and create post via API
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
