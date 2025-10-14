import os

def create_socialink_project():
    root_files = [
        "Socialink.xcodeproj",
        "SocialinkApp.swift",
        "SceneDelegate.swift",
        "AppDelegate.swift"
    ]
    
    structure = {
        "Models": [
            "User.swift",
            "Post.swift",
            "Comment.swift",
            "Story.swift",
            "Pod.swift",
            "Group.swift",
            "Event.swift",
            "MarketplaceItem.swift",
            "Notification.swift",
            "Message.swift",
            "Ad.swift",
            "LiveStream.swift",
            "Reaction.swift",
            "Poll.swift",
            "Page.swift",
            "SavedCollection.swift",
            "Memory.swift",
            "Fundraiser.swift",
            "Watch.swift",
            "Gaming.swift",
            "Dating.swift",
            "Job.swift",
            "Sticker.swift",
            "GIF.swift",
            "Badge.swift",
            "Subscription.swift",
            "Payment.swift",
            "Report.swift",
            "Block.swift",
            "FriendRequest.swift"
        ],
        "Views/Feed": [
            "FeedView.swift",
            "PostDetailView.swift"
        ],
        "Views/Profile": [
            "ProfileView.swift",
            "FriendListView.swift",
            "EditProfileView.swift"
        ],
        "Views/Stories": [
            "StoriesView.swift"
        ],
        "Views/Pods": [
            "PodsFeedView.swift",
            "PodPlayerView.swift",
            "CreatePodView.swift"
        ],
        "Views/Marketplace": [
            "MarketplaceView.swift",
            "ItemDetailView.swift",
            "CreateListingView.swift"
        ],
        "Views/Messaging": [
            "ChatListView.swift",
            "ChatView.swift",
            "GroupChatView.swift"
        ],
        "Views/Live": [
            "LiveStreamView.swift"
        ],
        "Views/PostCreation": [
            "CreatePostView.swift"
        ],
        "Views/Search": [
            "SearchView.swift"
        ],
        "Views/Notifications": [
            "NotificationsView.swift"
        ],
        "Views/Settings": [
            "SettingsView.swift",
            "PrivacySettingsView.swift"
        ],
        "Views/Login": [
            "LoginView.swift"
        ],
        "Views/Onboarding": [
            "OnboardingView.swift"
        ],
        "Views/Pages": [
            "PagesView.swift",
            "PageDetailView.swift",
            "CreatePageView.swift"
        ],
        "Views/Groups": [
            "GroupsView.swift",
            "GroupDetailView.swift",
            "CreateGroupView.swift",
            "GroupSettingsView.swift"
        ],
        "Views/Events": [
            "EventsView.swift",
            "EventDetailView.swift",
            "CreateEventView.swift"
        ],
        "Views/Watch": [
            "WatchView.swift",
            "WatchPlayerView.swift"
        ],
        "Views/Gaming": [
            "GamingView.swift",
            "PlayGameView.swift"
        ],
        "Views/Dating": [
            "DatingView.swift",
            "ProfileMatchView.swift",
            "MatchesView.swift"
        ],
        "Views/Jobs": [
            "JobsView.swift",
            "JobDetailView.swift",
            "ApplyJobView.swift"
        ],
        "Views/Fundraisers": [
            "FundraisersView.swift",
            "FundraiserDetailView.swift",
            "CreateFundraiserView.swift"
        ],
        "Views/Memories": [
            "MemoriesView.swift"
        ],
        "Views/Saved": [
            "SavedView.swift",
            "CollectionsView.swift"
        ],
        "Views/Menu": [
            "MenuView.swift",
            "ShortcutsView.swift"
        ],
        "Views/Moderation": [
            "ReportView.swift",
            "BlockedUsersView.swift",
            "ContentModerationView.swift"
        ],
        "Views": [
            "Main.storyboard"
        ],
        "ViewModels": [
            "FeedViewModel.swift",
            "ProfileViewModel.swift",
            "StoriesViewModel.swift",
            "PodsViewModel.swift",
            "MarketplaceViewModel.swift",
            "ChatViewModel.swift",
            "LiveViewModel.swift",
            "SearchViewModel.swift",
            "NotificationsViewModel.swift",
            "AuthViewModel.swift",
            "AdViewModel.swift",
            "GroupViewModel.swift",
            "EventViewModel.swift",
            "PageViewModel.swift",
            "WatchViewModel.swift",
            "GamingViewModel.swift",
            "DatingViewModel.swift",
            "JobsViewModel.swift",
            "FundraiserViewModel.swift",
            "MemoriesViewModel.swift",
            "SavedViewModel.swift",
            "ModerationViewModel.swift"
        ],
        "Services/API": [
            "APIService.swift",
            "AuthService.swift",
            "FeedAPIService.swift",
            "PodsAPIService.swift",
            "MarketplaceAPIService.swift",
            "GroupAPIService.swift",
            "EventAPIService.swift",
            "PageAPIService.swift",
            "WatchAPIService.swift",
            "GamingAPIService.swift",
            "DatingAPIService.swift",
            "JobsAPIService.swift",
            "FundraiserAPIService.swift",
            "ModerationAPIService.swift"
        ],
        "Services/Storage": [
            "DatabaseService.swift",
            "ImageCache.swift",
            "VideoCache.swift"
        ],
        "Services/Notifications": [
            "NotificationService.swift"
        ],
        "Services/Analytics": [
            "AnalyticsService.swift"
        ],
        "Services/Location": [
            "LocationService.swift"
        ],
        "Services/Media": [
            "MediaService.swift"
        ],
        "Services/Ads": [
            "AdService.swift"
        ],
        "Services/Live": [
            "LiveStreamingService.swift"
        ],
        "Services/Payment": [
            "PaymentService.swift",
            "SubscriptionService.swift",
            "DonationService.swift"
        ],
        "Services/Moderation": [
            "ContentModerationService.swift",
            "SpamDetectionService.swift",
            "ReportingService.swift"
        ],
        "Services/Recommendations": [
            "RecommendationEngine.swift",
            "FeedAlgorithmService.swift",
            "PersonalizationService.swift"
        ],
        "Services/Security": [
            "EncryptionService.swift",
            "BiometricAuthService.swift",
            "TwoFactorAuthService.swift"
        ],
        "Services/Social": [
            "FriendSuggestionService.swift",
            "ShareService.swift",
            "InviteService.swift"
        ],
        "Services/Sync": [
            "SyncService.swift",
            "ConflictResolutionService.swift",
            "OfflineModeService.swift"
        ],
        "Utilities/Extensions": [
            "Date+Formatting.swift",
            "String+Validation.swift",
            "UIImage+Effects.swift",
            "URL+Helpers.swift",
            "Array+SafeAccess.swift",
            "UIView+Animations.swift",
            "Color+Palette.swift",
            "Data+Compression.swift"
        ],
        "Utilities": [
            "Constants.swift",
            "ErrorHandler.swift",
            "Logger.swift",
            "ThemeManager.swift",
            "PermissionsManager.swift",
            "NetworkMonitor.swift",
            "KeychainManager.swift",
            "DeepLinkHandler.swift",
            "ImageCompressor.swift",
            "VideoCompressor.swift",
            "ValidationHelper.swift",
            "DateHelper.swift",
            "CurrencyFormatter.swift",
            "LinkPreviewGenerator.swift"
        ],
        "Utilities/Managers": [
            "SessionManager.swift",
            "CacheManager.swift",
            "DownloadManager.swift",
            "UploadManager.swift",
            "AccessibilityManager.swift"
        ],
        "Resources": [
            "Assets.xcassets",
            "Localizable.strings",
            "PrivacyInfo.xcprivacy",
            "Info.plist"
        ],
        "Resources/Fonts": [],
        "Resources/Sounds": [],
        "Resources/Videos": [],
        "Resources/Animations": [],
        "Config": [
            "Development.xcconfig",
            "Staging.xcconfig",
            "Production.xcconfig"
        ],
        "Components": [
            "CustomButton.swift",
            "CustomTextField.swift",
            "LoadingIndicator.swift",
            "EmptyStateView.swift",
            "ErrorView.swift",
            "MediaPicker.swift",
            "ImageViewer.swift",
            "VideoPlayer.swift",
            "AudioPlayer.swift",
            "PullToRefresh.swift",
            "InfiniteScroll.swift",
            "TabBar.swift",
            "NavigationBar.swift",
            "SearchBar.swift",
            "FilterChips.swift",
            "ActionSheet.swift",
            "Toast.swift",
            "Badge.swift",
            "Avatar.swift",
            "CardView.swift",
            "ListRow.swift",
            "Separator.swift"
        ],
        "Middleware": [
            "AuthMiddleware.swift",
            "LoggingMiddleware.swift",
            "ErrorMiddleware.swift",
            "RateLimitMiddleware.swift"
        ],
        "Coordinators": [
            "AppCoordinator.swift",
            "AuthCoordinator.swift",
            "MainCoordinator.swift",
            "FeedCoordinator.swift",
            "ProfileCoordinator.swift"
        ],
        "Tests/SocialinkTests": [
            "FeedViewModelTests.swift",
            "PodsViewModelTests.swift",
            "APIServiceTests.swift",
            "AuthServiceTests.swift",
            "MarketplaceViewModelTests.swift",
            "ChatViewModelTests.swift"
        ],
        "Tests/SocialinkUITests": [
            "FeedUITests.swift",
            "PodsUITests.swift",
            "LoginUITests.swift",
            "MarketplaceUITests.swift",
            "MessagingUITests.swift",
            "GroupsUITests.swift",
            "EventsUITests.swift"
        ],
        "Tests/SocialinkPerformanceTests": [
            "FeedPerformanceTests.swift",
            "ImageLoadingPerformanceTests.swift",
            "VideoStreamingPerformanceTests.swift"
        ],
        "Tests/SocialinkIntegrationTests": [
            "AuthFlowTests.swift",
            "PaymentFlowTests.swift",
            "NotificationFlowTests.swift"
        ],
        "Documentation": [
            "README.md",
            "ARCHITECTURE.md",
            "API_DOCS.md",
            "CONTRIBUTING.md",
            "CHANGELOG.md"
        ],
        "Scripts": [
            "setup.sh",
            "build.sh",
            "deploy.sh",
            "test.sh",
            "lint.sh"
        ]
    }
    
    base_dir = "Socialink"
    
    if not os.path.exists(base_dir):
        os.makedirs(base_dir)
        print(f"Created root directory: {base_dir}")
    
    for file in root_files:
        file_path = os.path.join(base_dir, file)
        if not file.endswith('.xcodeproj'):
            open(file_path, 'a').close()
            print(f"Created: {file_path}")
        else:
            if not os.path.exists(file_path):
                os.makedirs(file_path)
                print(f"Created: {file_path}")
    
    for folder, files in structure.items():
        folder_path = os.path.join(base_dir, folder)
        if not os.path.exists(folder_path):
            os.makedirs(folder_path)
            print(f"Created directory: {folder_path}")
        
        for file in files:
            file_path = os.path.join(folder_path, file)
            if not file.endswith('.xcassets') and not file.endswith('.storyboard'):
                open(file_path, 'a').close()
                print(f"Created: {file_path}")
            else:
                if not os.path.exists(file_path):
                    os.makedirs(file_path)
                    print(f"Created: {file_path}")
    
    print(f"\n✅ Socialink project structure created successfully in '{base_dir}/' directory!")
    print(f"📱 Total files/folders created: {sum(len(files) for files in structure.values()) + len(root_files)} items")

if __name__ == "__main__":
    create_socialink_project()
