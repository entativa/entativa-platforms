import os

def create_socialink_android_project():
    root_files = [
        "build.gradle.kts",
        "settings.gradle.kts",
        "gradle.properties",
        "gradlew",
        "gradlew.bat",
        ".gitignore",
        "local.properties"
    ]
    
    structure = {
        "app": [
            "build.gradle.kts",
            "proguard-rules.pro"
        ],
        "app/src/main": [
            "AndroidManifest.xml"
        ],
        "app/src/main/kotlin/com/entativa/socialink": [
            "MainActivity.kt",
            "SocialinkApplication.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/di": [
            "AppModule.kt",
            "NetworkModule.kt",
            "RepositoryModule.kt",
            "DatabaseModule.kt",
            "ViewModelModule.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/data/local/entities": [
            "UserEntity.kt",
            "PostEntity.kt",
            "PodEntity.kt",
            "StoryEntity.kt",
            "CommentEntity.kt",
            "MessageEntity.kt",
            "GroupEntity.kt",
            "EventEntity.kt",
            "MarketplaceItemEntity.kt",
            "NotificationEntity.kt",
            "PageEntity.kt",
            "ReactionEntity.kt",
            "PollEntity.kt",
            "SavedItemEntity.kt",
            "FriendRequestEntity.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/data/local/dao": [
            "UserDao.kt",
            "PostDao.kt",
            "PodDao.kt",
            "StoryDao.kt",
            "CommentDao.kt",
            "MessageDao.kt",
            "GroupDao.kt",
            "EventDao.kt",
            "MarketplaceItemDao.kt",
            "NotificationDao.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/data/local": [
            "AppDatabase.kt",
            "Converters.kt",
            "PreferencesManager.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/data/remote/dto": [
            "UserDto.kt",
            "PostDto.kt",
            "PodDto.kt",
            "StoryDto.kt",
            "CommentDto.kt",
            "MessageDto.kt",
            "GroupDto.kt",
            "EventDto.kt",
            "MarketplaceItemDto.kt",
            "NotificationDto.kt",
            "AdDto.kt",
            "LiveStreamDto.kt",
            "PageDto.kt",
            "ReactionDto.kt",
            "PollDto.kt",
            "JobDto.kt",
            "FundraiserDto.kt",
            "DatingProfileDto.kt",
            "PaymentDto.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/data/remote/api": [
            "FeedApi.kt",
            "PodsApi.kt",
            "AuthApi.kt",
            "ProfileApi.kt",
            "ChatApi.kt",
            "GroupApi.kt",
            "EventApi.kt",
            "MarketplaceApi.kt",
            "StoryApi.kt",
            "LiveApi.kt",
            "PageApi.kt",
            "NotificationApi.kt",
            "SearchApi.kt",
            "WatchApi.kt",
            "GamingApi.kt",
            "DatingApi.kt",
            "JobsApi.kt",
            "FundraiserApi.kt",
            "PaymentApi.kt",
            "ModerationApi.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/data/remote": [
            "ApiClient.kt",
            "AuthInterceptor.kt",
            "NetworkErrorHandler.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/data/repository": [
            "UserRepository.kt",
            "FeedRepository.kt",
            "PodsRepository.kt",
            "StoryRepository.kt",
            "ChatRepository.kt",
            "GroupRepository.kt",
            "EventRepository.kt",
            "MarketplaceRepository.kt",
            "NotificationRepository.kt",
            "LiveRepository.kt",
            "PageRepository.kt",
            "SearchRepository.kt",
            "WatchRepository.kt",
            "GamingRepository.kt",
            "DatingRepository.kt",
            "JobsRepository.kt",
            "FundraiserRepository.kt",
            "PaymentRepository.kt",
            "ModerationRepository.kt",
            "AuthRepository.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/domain/model": [
            "User.kt",
            "Post.kt",
            "Pod.kt",
            "Story.kt",
            "Comment.kt",
            "Message.kt",
            "Group.kt",
            "Event.kt",
            "MarketplaceItem.kt",
            "Notification.kt",
            "Ad.kt",
            "LiveStream.kt",
            "Page.kt",
            "Reaction.kt",
            "Poll.kt",
            "SavedCollection.kt",
            "Memory.kt",
            "Fundraiser.kt",
            "Job.kt",
            "DatingProfile.kt",
            "Payment.kt",
            "Subscription.kt",
            "Badge.kt",
            "Sticker.kt",
            "GIF.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/domain/usecase": [
            "GetFeedUseCase.kt",
            "CreatePodUseCase.kt",
            "LoginUseCase.kt",
            "LogoutUseCase.kt",
            "RegisterUseCase.kt",
            "GetProfileUseCase.kt",
            "UpdateProfileUseCase.kt",
            "CreatePostUseCase.kt",
            "DeletePostUseCase.kt",
            "LikePostUseCase.kt",
            "CommentPostUseCase.kt",
            "SharePostUseCase.kt",
            "GetStoriesUseCase.kt",
            "CreateStoryUseCase.kt",
            "GetMessagesUseCase.kt",
            "SendMessageUseCase.kt",
            "GetNotificationsUseCase.kt",
            "MarkNotificationReadUseCase.kt",
            "GetGroupsUseCase.kt",
            "JoinGroupUseCase.kt",
            "GetEventsUseCase.kt",
            "RSVPEventUseCase.kt",
            "GetMarketplaceUseCase.kt",
            "PurchaseItemUseCase.kt",
            "StartLiveStreamUseCase.kt",
            "SearchUseCase.kt",
            "ReportContentUseCase.kt",
            "BlockUserUseCase.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/theme": [
            "Color.kt",
            "Theme.kt",
            "Type.kt",
            "Shape.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/navigation": [
            "NavGraph.kt",
            "Screen.kt",
            "NavigationActions.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/components": [
            "PostCard.kt",
            "PodPlayer.kt",
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
            "TabBar.kt",
            "TopAppBar.kt",
            "FloatingActionButton.kt",
            "SwipeableCard.kt",
            "CarouselView.kt",
            "StoryRing.kt",
            "ReactionPicker.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/feed": [
            "FeedScreen.kt",
            "FeedViewModel.kt",
            "PostDetailScreen.kt",
            "PostDetailViewModel.kt",
            "CreatePostScreen.kt",
            "CreatePostViewModel.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/profile": [
            "ProfileScreen.kt",
            "ProfileViewModel.kt",
            "EditProfileScreen.kt",
            "EditProfileViewModel.kt",
            "FriendListScreen.kt",
            "FriendListViewModel.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/pods": [
            "PodsFeedScreen.kt",
            "PodsFeedViewModel.kt",
            "PodPlayerScreen.kt",
            "PodPlayerViewModel.kt",
            "CreatePodScreen.kt",
            "CreatePodViewModel.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/stories": [
            "StoriesScreen.kt",
            "StoriesViewModel.kt",
            "StoryViewerScreen.kt",
            "CreateStoryScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/marketplace": [
            "MarketplaceScreen.kt",
            "MarketplaceViewModel.kt",
            "ItemDetailScreen.kt",
            "ItemDetailViewModel.kt",
            "CreateListingScreen.kt",
            "CreateListingViewModel.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/messaging": [
            "ChatListScreen.kt",
            "ChatListViewModel.kt",
            "ChatScreen.kt",
            "ChatViewModel.kt",
            "GroupChatScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/live": [
            "LiveStreamScreen.kt",
            "LiveStreamViewModel.kt",
            "LiveViewerScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/search": [
            "SearchScreen.kt",
            "SearchViewModel.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/notifications": [
            "NotificationsScreen.kt",
            "NotificationsViewModel.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/settings": [
            "SettingsScreen.kt",
            "SettingsViewModel.kt",
            "PrivacySettingsScreen.kt",
            "NotificationSettingsScreen.kt",
            "AccountSettingsScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/auth": [
            "LoginScreen.kt",
            "LoginViewModel.kt",
            "RegisterScreen.kt",
            "RegisterViewModel.kt",
            "ForgotPasswordScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/onboarding": [
            "OnboardingScreen.kt",
            "OnboardingViewModel.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/groups": [
            "GroupsScreen.kt",
            "GroupsViewModel.kt",
            "GroupDetailScreen.kt",
            "GroupDetailViewModel.kt",
            "CreateGroupScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/events": [
            "EventsScreen.kt",
            "EventsViewModel.kt",
            "EventDetailScreen.kt",
            "EventDetailViewModel.kt",
            "CreateEventScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/pages": [
            "PagesScreen.kt",
            "PagesViewModel.kt",
            "PageDetailScreen.kt",
            "CreatePageScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/watch": [
            "WatchScreen.kt",
            "WatchViewModel.kt",
            "WatchPlayerScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/gaming": [
            "GamingScreen.kt",
            "GamingViewModel.kt",
            "PlayGameScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/dating": [
            "DatingScreen.kt",
            "DatingViewModel.kt",
            "ProfileMatchScreen.kt",
            "MatchesScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/jobs": [
            "JobsScreen.kt",
            "JobsViewModel.kt",
            "JobDetailScreen.kt",
            "ApplyJobScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/fundraisers": [
            "FundraisersScreen.kt",
            "FundraisersViewModel.kt",
            "FundraiserDetailScreen.kt",
            "CreateFundraiserScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/memories": [
            "MemoriesScreen.kt",
            "MemoriesViewModel.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/saved": [
            "SavedScreen.kt",
            "SavedViewModel.kt",
            "CollectionsScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/menu": [
            "MenuScreen.kt",
            "MenuViewModel.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/ui/moderation": [
            "ReportScreen.kt",
            "BlockedUsersScreen.kt",
            "ContentModerationScreen.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/service": [
            "NotificationService.kt",
            "MessagingService.kt",
            "UploadService.kt",
            "DownloadService.kt",
            "SyncService.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/receiver": [
            "NetworkReceiver.kt",
            "NotificationReceiver.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/worker": [
            "SyncWorker.kt",
            "UploadWorker.kt",
            "NotificationWorker.kt",
            "CleanupWorker.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/util": [
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
            "CurrencyFormatter.kt",
            "LinkPreviewGenerator.kt",
            "BiometricHelper.kt",
            "EncryptionHelper.kt"
        ],
        "app/src/main/kotlin/com/entativa/socialink/util/manager": [
            "SessionManager.kt",
            "CacheManager.kt",
            "DownloadManager.kt",
            "UploadManager.kt",
            "PermissionManager.kt",
            "NotificationManager.kt"
        ],
        "app/src/main/res/drawable": [],
        "app/src/main/res/font": [],
        "app/src/main/res/raw": [],
        "app/src/main/res/values": [
            "colors.xml",
            "strings.xml",
            "themes.xml",
            "dimens.xml",
            "styles.xml"
        ],
        "app/src/main/res/values-night": [
            "colors.xml",
            "themes.xml"
        ],
        "app/src/main/res/xml": [
            "network_security_config.xml",
            "file_paths.xml",
            "backup_rules.xml"
        ],
        "app/src/test/kotlin/com/entativa/socialink/repository": [
            "FeedRepositoryTest.kt",
            "PodsRepositoryTest.kt",
            "AuthRepositoryTest.kt",
            "ChatRepositoryTest.kt",
            "GroupRepositoryTest.kt"
        ],
        "app/src/test/kotlin/com/entativa/socialink/usecase": [
            "GetFeedUseCaseTest.kt",
            "CreatePodUseCaseTest.kt",
            "LoginUseCaseTest.kt",
            "SendMessageUseCaseTest.kt"
        ],
        "app/src/test/kotlin/com/entativa/socialink/viewmodel": [
            "FeedViewModelTest.kt",
            "PodsViewModelTest.kt",
            "AuthViewModelTest.kt"
        ],
        "app/src/androidTest/kotlin/com/entativa/socialink/ui": [
            "FeedScreenTest.kt",
            "PodsScreenTest.kt",
            "LoginScreenTest.kt",
            "ChatScreenTest.kt",
            "GroupsScreenTest.kt"
        ],
        "app/src/androidTest/kotlin/com/entativa/socialink/database": [
            "UserDaoTest.kt",
            "PostDaoTest.kt",
            "MessageDaoTest.kt"
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
            "SETUP.md"
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
            "pr_checks.yml"
        ]
    }
    
    base_dir = "SocialinkAndroid"
    
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
    
    print(f"\n✅ Socialink Android project structure created successfully in '{base_dir}/' directory!")
    print(f"📱 Total files/folders created: {sum(len(files) for files in structure.values()) + len(root_files) + 2} items")
    print(f"\n💡 Next steps:")
    print(f"   1. Open project in Android Studio")
    print(f"   2. Sync Gradle dependencies")
    print(f"   3. Configure API keys in local.properties")
    print(f"   4. Update AndroidManifest.xml with required permissions")

if __name__ == "__main__":
    create_socialink_android_project()
