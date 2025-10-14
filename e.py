import os

def create_vignette_android_project():
    root_files = [
        "build.gradle.kts",
        "settings.gradle.kts",
        "gradle.properties",
        "gradlew",
        "gradlew.bat",
        ".gitignore",
        "local.properties",
        "README.md"
    ]
    
    structure = {
        "app": [
            "build.gradle.kts",
            "proguard-rules.pro"
        ],
        "app/src/main": [
            "AndroidManifest.xml"
        ],
        "app/src/main/kotlin/com/entativa/vignette": [
            "MainActivity.kt",
            "VignetteApplication.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/di": [
            "AppModule.kt",
            "NetworkModule.kt",
            "RepositoryModule.kt",
            "DatabaseModule.kt",
            "ViewModelModule.kt",
            "UseCaseModule.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/data/local/entities": [
            "UserEntity.kt",
            "PostEntity.kt",
            "TakeEntity.kt",
            "StoryEntity.kt",
            "CommentEntity.kt",
            "LikeEntity.kt",
            "FollowEntity.kt",
            "MessageEntity.kt",
            "CollectionEntity.kt",
            "HighlightEntity.kt",
            "TagEntity.kt",
            "LocationEntity.kt",
            "HashtagEntity.kt",
            "SavedPostEntity.kt",
            "ArchiveEntity.kt",
            "NotificationEntity.kt",
            "LiveStreamEntity.kt",
            "GuideEntity.kt",
            "ProductEntity.kt",
            "OrderEntity.kt",
            "ChallengeEntity.kt",
            "BehindTheTakeEntity.kt",
            "CreatorEntity.kt",
            "SubscriptionEntity.kt",
            "CloseFriendEntity.kt",
            "BlockedUserEntity.kt",
            "MutedUserEntity.kt",
            "RestrictedUserEntity.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/data/local/dao": [
            "UserDao.kt",
            "PostDao.kt",
            "TakeDao.kt",
            "StoryDao.kt",
            "CommentDao.kt",
            "MessageDao.kt",
            "CollectionDao.kt",
            "NotificationDao.kt",
            "ChallengeDao.kt",
            "BehindTheTakeDao.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/data/local": [
            "AppDatabase.kt",
            "Converters.kt",
            "PreferencesManager.kt",
            "SecurePreferences.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/data/remote/dto": [
            "UserDto.kt",
            "PostDto.kt",
            "TakeDto.kt",
            "StoryDto.kt",
            "CommentDto.kt",
            "LikeDto.kt",
            "FollowDto.kt",
            "MessageDto.kt",
            "CollectionDto.kt",
            "HighlightDto.kt",
            "TagDto.kt",
            "LocationDto.kt",
            "HashtagDto.kt",
            "NotificationDto.kt",
            "LiveStreamDto.kt",
            "GuideDto.kt",
            "ProductDto.kt",
            "OrderDto.kt",
            "ChallengeDto.kt",
            "BehindTheTakeDto.kt",
            "CreatorDto.kt",
            "InsightsDto.kt",
            "AnalyticsDto.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/data/remote/api": [
            "AuthApi.kt",
            "FeedApi.kt",
            "ProfileApi.kt",
            "TakeApi.kt",
            "StoryApi.kt",
            "ExploreApi.kt",
            "DirectApi.kt",
            "ActivityApi.kt",
            "ShoppingApi.kt",
            "LiveApi.kt",
            "GuideApi.kt",
            "SearchApi.kt",
            "HashtagApi.kt",
            "LocationApi.kt",
            "UploadApi.kt",
            "ChallengeApi.kt",
            "BehindTheTakeApi.kt",
            "CreatorApi.kt",
            "InsightsApi.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/data/remote": [
            "ApiClient.kt",
            "AuthInterceptor.kt",
            "NetworkErrorHandler.kt",
            "RetryInterceptor.kt",
            "LoggingInterceptor.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/data/repository": [
            "AuthRepository.kt",
            "UserRepository.kt",
            "FeedRepository.kt",
            "TakeRepository.kt",
            "StoryRepository.kt",
            "ExploreRepository.kt",
            "DirectRepository.kt",
            "ActivityRepository.kt",
            "ShoppingRepository.kt",
            "LiveRepository.kt",
            "GuideRepository.kt",
            "SearchRepository.kt",
            "HashtagRepository.kt",
            "LocationRepository.kt",
            "UploadRepository.kt",
            "ChallengeRepository.kt",
            "BehindTheTakeRepository.kt",
            "CreatorRepository.kt",
            "InsightsRepository.kt",
            "SocialRepository.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/domain/model": [
            "User.kt",
            "Post.kt",
            "Take.kt",
            "Story.kt",
            "Comment.kt",
            "Like.kt",
            "Follow.kt",
            "Message.kt",
            "Collection.kt",
            "Highlight.kt",
            "Tag.kt",
            "Location.kt",
            "Hashtag.kt",
            "Mention.kt",
            "SavedPost.kt",
            "Archive.kt",
            "Activity.kt",
            "Notification.kt",
            "LiveStream.kt",
            "Guide.kt",
            "Product.kt",
            "ShoppingBag.kt",
            "Order.kt",
            "Challenge.kt",
            "BehindTheTake.kt",
            "Creator.kt",
            "Badge.kt",
            "Subscription.kt",
            "CloseFriend.kt",
            "BlockedUser.kt",
            "MutedUser.kt",
            "RestrictedUser.kt",
            "Insights.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/domain/usecase": [
            "LoginUseCase.kt",
            "LogoutUseCase.kt",
            "RegisterUseCase.kt",
            "GetFeedUseCase.kt",
            "CreatePostUseCase.kt",
            "DeletePostUseCase.kt",
            "LikePostUseCase.kt",
            "UnlikePostUseCase.kt",
            "CommentPostUseCase.kt",
            "SharePostUseCase.kt",
            "SavePostUseCase.kt",
            "GetProfileUseCase.kt",
            "UpdateProfileUseCase.kt",
            "FollowUserUseCase.kt",
            "UnfollowUserUseCase.kt",
            "GetStoriesUseCase.kt",
            "CreateStoryUseCase.kt",
            "ViewStoryUseCase.kt",
            "GetTakesUseCase.kt",
            "CreateTakeUseCase.kt",
            "LikeTakeUseCase.kt",
            "CommentTakeUseCase.kt",
            "GetBehindTheTakeUseCase.kt",
            "CreateBehindTheTakeUseCase.kt",
            "GetChallengesUseCase.kt",
            "JoinChallengeUseCase.kt",
            "CreateChallengeUseCase.kt",
            "GetMessagesUseCase.kt",
            "SendMessageUseCase.kt",
            "GetNotificationsUseCase.kt",
            "MarkNotificationReadUseCase.kt",
            "SearchUseCase.kt",
            "ExploreContentUseCase.kt",
            "StartLiveStreamUseCase.kt",
            "BlockUserUseCase.kt",
            "MuteUserUseCase.kt",
            "RestrictUserUseCase.kt",
            "ReportContentUseCase.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/theme": [
            "Color.kt",
            "Theme.kt",
            "Type.kt",
            "Shape.kt",
            "Dimension.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/navigation": [
            "NavGraph.kt",
            "Screen.kt",
            "NavigationActions.kt",
            "BottomNavigation.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/components": [
            "PostCard.kt",
            "TakePlayer.kt",
            "StoryRing.kt",
            "LoadingIndicator.kt",
            "ErrorView.kt",
            "EmptyStateView.kt",
            "CustomButton.kt",
            "CustomTextField.kt",
            "ImagePicker.kt",
            "VideoPlayer.kt",
            "AudioPlayer.kt",
            "PullToRefresh.kt",
            "SearchBar.kt",
            "FilterChips.kt",
            "BottomSheet.kt",
            "Dialog.kt",
            "Toast.kt",
            "Avatar.kt",
            "Badge.kt",
            "Separator.kt",
            "TopAppBar.kt",
            "FloatingActionButton.kt",
            "SwipeableCard.kt",
            "CarouselView.kt",
            "ReactionPicker.kt",
            "LikeButton.kt",
            "CommentButton.kt",
            "ShareButton.kt",
            "BookmarkButton.kt",
            "FollowButton.kt",
            "HashtagText.kt",
            "MentionText.kt",
            "VerifiedBadge.kt",
            "DoubleTapLike.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/feed": [
            "FeedScreen.kt",
            "FeedViewModel.kt",
            "PostDetailScreen.kt",
            "PostDetailViewModel.kt",
            "CreatePostScreen.kt",
            "CreatePostViewModel.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/stories": [
            "StoriesBarScreen.kt",
            "StoriesBarViewModel.kt",
            "StoryScreen.kt",
            "StoryViewModel.kt",
            "StoryPlayerScreen.kt",
            "StoryCreatorScreen.kt",
            "StoryCreatorViewModel.kt",
            "StoryViewersScreen.kt",
            "HighlightsScreen.kt",
            "HighlightCreatorScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/takes": [
            "TakesScreen.kt",
            "TakesViewModel.kt",
            "TakePlayerScreen.kt",
            "TakeCreatorScreen.kt",
            "TakeCreatorViewModel.kt",
            "TakeEffectsScreen.kt",
            "TakeAudioPickerScreen.kt",
            "TakesCameraScreen.kt",
            "BehindTheTakeScreen.kt",
            "BehindTheTakeViewModel.kt",
            "BTTCreatorScreen.kt",
            "ChallengeScreen.kt",
            "ChallengeDetailScreen.kt",
            "CreateChallengeScreen.kt",
            "TrendingChallengesScreen.kt",
            "ChallengeFeedScreen.kt",
            "ChallengeViewModel.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/profile": [
            "ProfileScreen.kt",
            "ProfileViewModel.kt",
            "EditProfileScreen.kt",
            "ProfileGridScreen.kt",
            "ProfileHeaderScreen.kt",
            "FollowersScreen.kt",
            "FollowingScreen.kt",
            "TaggedPostsScreen.kt",
            "SavedPostsScreen.kt",
            "ArchiveScreen.kt",
            "CloseFriendsScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/explore": [
            "ExploreScreen.kt",
            "ExploreViewModel.kt",
            "ExploreGridScreen.kt",
            "TrendingScreen.kt",
            "SearchScreen.kt",
            "SearchViewModel.kt",
            "RecentSearchesScreen.kt",
            "SuggestedAccountsScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/camera": [
            "CameraScreen.kt",
            "CameraViewModel.kt",
            "CameraPreviewScreen.kt",
            "CameraControlsScreen.kt",
            "FilterPickerScreen.kt",
            "EffectsLibraryScreen.kt",
            "ARFiltersScreen.kt",
            "BoomerangScreen.kt",
            "LayoutScreen.kt",
            "SuperzoomScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/postcreation": [
            "MediaPickerScreen.kt",
            "MultiSelectScreen.kt",
            "FilterScreen.kt",
            "EditPhotoScreen.kt",
            "EditVideoScreen.kt",
            "CaptionScreen.kt",
            "TagPeopleScreen.kt",
            "AddLocationScreen.kt",
            "AdvancedSettingsScreen.kt",
            "AltTextScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/direct": [
            "DirectInboxScreen.kt",
            "DirectInboxViewModel.kt",
            "ChatScreen.kt",
            "ChatViewModel.kt",
            "ChatBubbleScreen.kt",
            "MediaMessageScreen.kt",
            "VoiceMessageScreen.kt",
            "GroupChatScreen.kt",
            "NewMessageScreen.kt",
            "ThreadSettingsScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/activity": [
            "ActivityScreen.kt",
            "ActivityViewModel.kt",
            "NotificationsScreen.kt",
            "FollowRequestsScreen.kt",
            "LikesScreen.kt",
            "CommentsScreen.kt",
            "MentionsScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/shopping": [
            "ShopScreen.kt",
            "ShopViewModel.kt",
            "ProductDetailScreen.kt",
            "ShoppingBagScreen.kt",
            "CheckoutScreen.kt",
            "OrderHistoryScreen.kt",
            "WishlistScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/live": [
            "LiveStreamScreen.kt",
            "LiveStreamViewModel.kt",
            "LiveViewerScreen.kt",
            "LiveCommentsScreen.kt",
            "LiveRequestsScreen.kt",
            "GoLiveScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/guides": [
            "GuidesScreen.kt",
            "GuidesViewModel.kt",
            "GuideDetailScreen.kt",
            "CreateGuideScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/settings": [
            "SettingsScreen.kt",
            "SettingsViewModel.kt",
            "AccountSettingsScreen.kt",
            "PrivacySettingsScreen.kt",
            "SecuritySettingsScreen.kt",
            "NotificationSettingsScreen.kt",
            "DataUsageScreen.kt",
            "LanguageScreen.kt",
            "ThemeScreen.kt",
            "AboutScreen.kt",
            "HelpScreen.kt",
            "BlockedAccountsScreen.kt",
            "MutedAccountsScreen.kt",
            "RestrictedAccountsScreen.kt",
            "CloseFriendsSettingsScreen.kt",
            "TwoFactorAuthScreen.kt",
            "LoginActivityScreen.kt",
            "DownloadDataScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/ui/auth": [
            "LoginScreen.kt",
            "LoginViewModel.kt",
            "SignupScreen.kt",
            "SignupViewModel.kt",
            "ForgotPasswordScreen.kt",
            "OnboardingScreen.kt",
            "ProfileSetupScreen.kt",
            "FindFriendsScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/service": [
            "NotificationService.kt",
            "MessagingService.kt",
            "UploadService.kt",
            "DownloadService.kt",
            "SyncService.kt",
            "LiveStreamService.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/receiver": [
            "NetworkReceiver.kt",
            "NotificationReceiver.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/worker": [
            "SyncWorker.kt",
            "UploadWorker.kt",
            "NotificationWorker.kt",
            "CleanupWorker.kt",
            "StoryExpirationWorker.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/util": [
            "Extensions.kt",
            "Constants.kt",
            "ErrorHandler.kt",
            "Permissions.kt",
            "NetworkMonitor.kt",
            "DateHelper.kt",
            "ImageCompressor.kt",
            "VideoCompressor.kt",
            "ValidationHelper.kt",
            "DeepLinkHandler.kt",
            "BiometricHelper.kt",
            "EncryptionHelper.kt",
            "HashtagParser.kt",
            "MentionParser.kt",
            "LinkPreviewGenerator.kt",
            "TextFormatter.kt",
            "NumberFormatter.kt",
            "ShareHelper.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/util/manager": [
            "SessionManager.kt",
            "CacheManager.kt",
            "DownloadManager.kt",
            "UploadManager.kt",
            "PermissionManager.kt",
            "NotificationManager.kt",
            "ThemeManager.kt",
            "LocalizationManager.kt",
            "HapticManager.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/camera": [
            "CameraManager.kt",
            "CameraPreview.kt",
            "CameraCapture.kt",
            "FilterEngine.kt",
            "ARFilterManager.kt",
            "EffectProcessor.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/media": [
            "ImageProcessor.kt",
            "VideoProcessor.kt",
            "AudioProcessor.kt",
            "ThumbnailGenerator.kt",
            "MediaCompressor.kt",
            "FilterApplier.kt",
            "EffectRenderer.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/player": [
            "VideoPlayerManager.kt",
            "ExoPlayerManager.kt",
            "AudioPlayerManager.kt",
            "PlaybackController.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/websocket": [
            "WebSocketClient.kt",
            "WebSocketManager.kt",
            "MessageHandler.kt",
            "PresenceManager.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/analytics": [
            "AnalyticsManager.kt",
            "EventTracker.kt",
            "InsightsCalculator.kt",
            "PerformanceMonitor.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/moderation": [
            "ContentModerationManager.kt",
            "SpamDetector.kt",
            "ReportHandler.kt",
            "FilterManager.kt"
        ],
        "app/src/main/kotlin/com/entativa/vignette/recommendation": [
            "RecommendationEngine.kt",
            "PersonalizationManager.kt",
            "AlgorithmManager.kt"
        ],
        "app/src/main/res/drawable": [],
        "app/src/main/res/font": [],
        "app/src/main/res/raw": [],
        "app/src/main/res/values": [
            "colors.xml",
            "strings.xml",
            "themes.xml",
            "dimens.xml",
            "styles.xml",
            "attrs.xml"
        ],
        "app/src/main/res/values-night": [
            "colors.xml",
            "themes.xml"
        ],
        "app/src/main/res/xml": [
            "network_security_config.xml",
            "file_paths.xml",
            "backup_rules.xml",
            "data_extraction_rules.xml"
        ],
        "app/src/main/res/anim": [],
        "app/src/main/res/animator": [],
        "app/src/main/res/layout": [
            "activity_main.xml"
        ],
        "app/src/test/kotlin/com/entativa/vignette/repository": [
            "FeedRepositoryTest.kt",
            "TakeRepositoryTest.kt",
            "StoryRepositoryTest.kt",
            "AuthRepositoryTest.kt",
            "DirectRepositoryTest.kt",
            "ChallengeRepositoryTest.kt"
        ],
        "app/src/test/kotlin/com/entativa/vignette/usecase": [
            "GetFeedUseCaseTest.kt",
            "CreateTakeUseCaseTest.kt",
            "CreateBehindTheTakeUseCaseTest.kt",
            "LoginUseCaseTest.kt",
            "SendMessageUseCaseTest.kt",
            "JoinChallengeUseCaseTest.kt"
        ],
        "app/src/test/kotlin/com/entativa/vignette/viewmodel": [
            "FeedViewModelTest.kt",
            "TakesViewModelTest.kt",
            "AuthViewModelTest.kt",
            "ChallengeViewModelTest.kt"
        ],
        "app/src/androidTest/kotlin/com/entativa/vignette/ui": [
            "FeedScreenTest.kt",
            "TakesScreenTest.kt",
            "StoryScreenTest.kt",
            "LoginScreenTest.kt",
            "ChatScreenTest.kt",
            "ChallengeScreenTest.kt",
            "BehindTheTakeScreenTest.kt"
        ],
        "app/src/androidTest/kotlin/com/entativa/vignette/database": [
            "UserDaoTest.kt",
            "PostDaoTest.kt",
            "TakeDaoTest.kt",
            "MessageDaoTest.kt",
            "ChallengeDaoTest.kt"
        ],
        "config": [
            "development.properties",
            "staging.properties",
            "production.properties"
        ],
        "buildSrc/src/main/kotlin": [
            "Dependencies.kt",
            "Versions.kt",
            "BuildConfig.kt"
        ],
        "docs": [
            "README.md",
            "ARCHITECTURE.md",
            "API_DOCS.md",
            "CONTRIBUTING.md",
            "CHANGELOG.md",
            "SETUP.md",
            "FEATURES.md"
        ],
        "scripts": [
            "setup.sh",
            "build.sh",
            "deploy.sh",
            "test.sh",
            "lint.sh",
            "clean.sh"
        ],
        ".github/workflows": [
            "ci.yml",
            "cd.yml",
            "pr_checks.yml",
            "release.yml"
        ]
    }
    
    base_dir = "VignetteAndroid"
    
    if not os.path.exists(base_dir):
        os.makedirs(base_dir)
        print(f"Created root directory: {base_dir}")
    
    for file in root_files:
        file_path = os.path.join(base_dir, file)
        open(file_path, 'a').close()
        print(f"Created: {file_path}")
    
    for folder, files in structure.items():
        folder_path = os.path.join(base_dir, folder)
        if not os.path.exists(folder_path):
            os.makedirs(folder_path)
            print(f"Created directory: {folder_path}")
        
        for file in files:
            file_path = os.path.join(folder_path, file)
            open(file_path, 'a').close()
            print(f"Created: {file_path}")
    
    gradle_wrapper_dir = os.path.join(base_dir, "gradle/wrapper")
    os.makedirs(gradle_wrapper_dir, exist_ok=True)
    open(os.path.join(gradle_wrapper_dir, "gradle-wrapper.properties"), 'a').close()
    open(os.path.join(gradle_wrapper_dir, "gradle-wrapper.jar"), 'a').close()
    print(f"Created gradle wrapper directory")
    
    print(f"\n✅ Vignette Android project structure created successfully in '{base_dir}/' directory!")
    print(f"\n📸 Vignette - Instagram Competitor Features:")
    print(f"   ✨ Core: Feed, Stories, Takes (short-form video), Live, Explore")
    print(f"   🎬 Takes: Viral short videos with vertical swipe discovery")
    print(f"   🎥 Behind-The-Takes (BTT): Creator transparency & making-of content")
    print(f"   🏆 Challenges: Trending challenges, participation, leaderboards")
    print(f"   💬 Social: Direct Messages (E2E encrypted), Comments, Likes, Follows")
    print(f"   🎨 Creative: Advanced Camera, AR Filters, Photo/Video Editing")
    print(f"   🛍️ Commerce: Shopping, Products, Checkout, Orders, Wishlist")
    print(f"   📚 Content: Guides, Collections, Saved Posts, Archive")
    print(f"   🔒 Privacy: Close Friends, Restrict, Mute, Block")
    print(f"   📊 Creator Tools: Insights, Analytics, Badges, Subscriptions")
    print(f"\n🏗️ Architecture:")
    print(f"   - Clean Architecture (MVVM + Use Cases)")
    print(f"   - Jetpack Compose for modern UI")
    print(f"   - Hilt for dependency injection")
    print(f"   - Room for offline-first storage")
    print(f"   - Retrofit + OkHttp for networking")
    print(f"   - ExoPlayer for video playback")
    print(f"   - WebSocket for real-time messaging")
    print(f"   - WorkManager for background tasks")
    print(f"   - CameraX for camera features")
    print(f"\n💡 Unique Features:")
    print(f"   - Behind-The-Takes system for viral content transparency")
    print(f"   - Challenge ecosystem with trending feeds")
    print(f"   - Creator-first analytics and insights")
    print(f"   - Advanced AR filters and effects")
    print(f"\n📁 Total files/folders created: {sum(len(files) for files in structure.values()) + len(root_files) + 2} items")

if __name__ == "__main__":
    create_vignette_android_project()
