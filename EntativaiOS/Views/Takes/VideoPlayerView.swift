import SwiftUI
import AVKit
import AVFoundation

struct VideoPlayerView: View {
    let videoURL: URL
    let isPlaying: Bool
    @State private var player: AVPlayer?
    @State private var isMuted: Bool = false
    @State private var showControls: Bool = false
    
    var body: some View {
        ZStack {
            if let player = player {
                VideoPlayer(player: player)
                    .onAppear {
                        setupPlayer()
                    }
                    .onChange(of: isPlaying) { oldValue, newValue in
                        if newValue {
                            player.play()
                        } else {
                            player.pause()
                        }
                    }
                    .onDisappear {
                        player.pause()
                    }
                    .onTapGesture {
                        togglePlayPause()
                    }
                
                // Mute/Unmute Button
                VStack {
                    Spacer()
                    
                    HStack {
                        Spacer()
                        
                        Button(action: {
                            isMuted.toggle()
                            player.isMuted = isMuted
                        }) {
                            Image(systemName: isMuted ? "speaker.slash.fill" : "speaker.wave.2.fill")
                                .font(.system(size: 20))
                                .foregroundColor(.white)
                                .padding(12)
                                .background(Color.black.opacity(0.5))
                                .clipShape(Circle())
                        }
                        .padding(.trailing, 16)
                        .padding(.bottom, 200)
                    }
                }
            }
        }
        .ignoresSafeArea()
    }
    
    private func setupPlayer() {
        let playerItem = AVPlayerItem(url: videoURL)
        player = AVPlayer(playerItem: playerItem)
        player?.isMuted = isMuted
        
        // Loop video
        NotificationCenter.default.addObserver(
            forName: .AVPlayerItemDidPlayToEndTime,
            object: playerItem,
            queue: .main
        ) { _ in
            player?.seek(to: .zero)
            player?.play()
        }
        
        if isPlaying {
            player?.play()
        }
    }
    
    private func togglePlayPause() {
        guard let player = player else { return }
        
        if player.timeControlStatus == .playing {
            player.pause()
        } else {
            player.play()
        }
        
        // Show controls temporarily
        withAnimation {
            showControls = true
        }
        
        DispatchQueue.main.asyncAfter(deadline: .now() + 2) {
            withAnimation {
                showControls = false
            }
        }
    }
}

// Custom Video Player with more control
struct CustomVideoPlayer: UIViewRepresentable {
    let player: AVPlayer
    
    func makeUIView(context: Context) -> UIView {
        let view = UIView(frame: .zero)
        
        let playerLayer = AVPlayerLayer(player: player)
        playerLayer.frame = view.bounds
        playerLayer.videoGravity = .resizeAspectFill
        view.layer.addSublayer(playerLayer)
        
        context.coordinator.playerLayer = playerLayer
        
        return view
    }
    
    func updateUIView(_ uiView: UIView, context: Context) {
        context.coordinator.playerLayer?.frame = uiView.bounds
    }
    
    func makeCoordinator() -> Coordinator {
        Coordinator()
    }
    
    class Coordinator {
        var playerLayer: AVPlayerLayer?
    }
}

// MARK: - Video Cache Manager
class VideoCache {
    static let shared = VideoCache()
    
    private var cache: [URL: AVPlayerItem] = [:]
    private let queue = DispatchQueue(label: "com.entativa.videocache")
    
    func preload(url: URL) {
        queue.async {
            if self.cache[url] == nil {
                let playerItem = AVPlayerItem(url: url)
                self.cache[url] = playerItem
            }
        }
    }
    
    func getPlayerItem(for url: URL) -> AVPlayerItem {
        return queue.sync {
            if let cached = cache[url] {
                return cached
            }
            
            let playerItem = AVPlayerItem(url: url)
            cache[url] = playerItem
            return playerItem
        }
    }
    
    func clearCache() {
        queue.async {
            self.cache.removeAll()
        }
    }
}
