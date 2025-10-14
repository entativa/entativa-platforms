import os

def create_vignette_backend():
    root_files = [
        "docker-compose.yml",
        "docker-compose.dev.yml",
        "docker-compose.prod.yml",
        ".gitignore",
        "README.md",
        "ARCHITECTURE.md",
        "Makefile",
        ".env.example"
    ]
    
    structure = {
        "services/user-service": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            ".dockerignore",
            "README.md"
        ],
        "services/user-service/cmd/api": [
            "main.go"
        ],
        "services/user-service/internal/handler": [
            "auth_handler.go",
            "user_handler.go",
            "profile_handler.go",
            "follow_handler.go",
            "block_handler.go",
            "mute_handler.go",
            "restrict_handler.go",
            "close_friends_handler.go",
            "verification_handler.go"
        ],
        "services/user-service/internal/service": [
            "auth_service.go",
            "user_service.go",
            "profile_service.go",
            "follow_service.go",
            "session_service.go",
            "verification_service.go",
            "privacy_service.go",
            "two_factor_service.go"
        ],
        "services/user-service/internal/repository": [
            "user_repository.go",
            "profile_repository.go",
            "follow_repository.go",
            "session_repository.go",
            "privacy_repository.go"
        ],
        "services/user-service/internal/model": [
            "user.go",
            "profile.go",
            "follow.go",
            "session.go",
            "privacy.go",
            "verification.go"
        ],
        "services/user-service/internal/middleware": [
            "auth_middleware.go",
            "cors_middleware.go",
            "rate_limit_middleware.go",
            "logger_middleware.go"
        ],
        "services/user-service/internal/config": [
            "config.go"
        ],
        "services/user-service/internal/util": [
            "jwt.go",
            "password.go",
            "validation.go",
            "email.go"
        ],
        "services/user-service/pkg/database": [
            "postgres.go",
            "redis.go",
            "migrations.go"
        ],
        "services/user-service/migrations": [
            "001_create_users_table.up.sql",
            "001_create_users_table.down.sql",
            "002_create_profiles_table.up.sql",
            "002_create_profiles_table.down.sql",
            "003_create_follows_table.up.sql",
            "003_create_follows_table.down.sql"
        ],
        "services/user-service/test": [
            "auth_test.go",
            "user_test.go",
            "follow_test.go",
            "integration_test.go"
        ],
        
        "services/feed-service": [
            "requirements.txt",
            "Dockerfile",
            ".dockerignore",
            "README.md",
            "setup.py",
            "pyproject.toml"
        ],
        "services/feed-service/app": [
            "__init__.py",
            "main.py",
            "config.py",
            "dependencies.py"
        ],
        "services/feed-service/app/api/v1": [
            "__init__.py",
            "feed.py",
            "posts.py",
            "comments.py",
            "likes.py",
            "shares.py",
            "saves.py",
            "archives.py"
        ],
        "services/feed-service/app/core": [
            "__init__.py",
            "security.py",
            "exceptions.py",
            "middleware.py"
        ],
        "services/feed-service/app/models": [
            "__init__.py",
            "post.py",
            "comment.py",
            "like.py",
            "save.py"
        ],
        "services/feed-service/app/schemas": [
            "__init__.py",
            "post.py",
            "comment.py",
            "feed.py"
        ],
        "services/feed-service/app/services": [
            "__init__.py",
            "feed_service.py",
            "post_service.py",
            "ranking_service.py",
            "algorithm_service.py",
            "personalization_service.py",
            "explore_service.py"
        ],
        "services/feed-service/app/ml": [
            "__init__.py",
            "content_ranker.py",
            "engagement_predictor.py",
            "recommendation_engine.py",
            "embeddings.py",
            "trending_detector.py"
        ],
        "services/feed-service/app/db": [
            "__init__.py",
            "mongodb.py",
            "redis_client.py",
            "repositories.py"
        ],
        "services/feed-service/tests": [
            "__init__.py",
            "test_feed.py",
            "test_posts.py",
            "test_algorithm.py"
        ],
        
        "services/takes-service": [
            "requirements.txt",
            "Dockerfile",
            ".dockerignore",
            "README.md"
        ],
        "services/takes-service/app": [
            "__init__.py",
            "main.py",
            "config.py"
        ],
        "services/takes-service/app/api": [
            "__init__.py",
            "takes.py",
            "behind_the_takes.py",
            "challenges.py",
            "effects.py",
            "audio.py"
        ],
        "services/takes-service/app/models": [
            "__init__.py",
            "take.py",
            "behind_the_take.py",
            "challenge.py",
            "effect.py"
        ],
        "services/takes-service/app/services": [
            "__init__.py",
            "take_service.py",
            "btt_service.py",
            "challenge_service.py",
            "discovery_service.py",
            "trending_service.py",
            "effect_service.py",
            "audio_library_service.py"
        ],
        "services/takes-service/app/ml": [
            "__init__.py",
            "viral_predictor.py",
            "content_classifier.py",
            "recommendation_engine.py",
            "hashtag_extractor.py"
        ],
        "services/takes-service/app/db": [
            "__init__.py",
            "mongodb.py",
            "redis_client.py"
        ],
        "services/takes-service/tests": [
            "__init__.py",
            "test_takes.py",
            "test_challenges.py",
            "test_btt.py"
        ],
        
        "services/story-service": [
            "requirements.txt",
            "Dockerfile",
            ".dockerignore",
            "README.md"
        ],
        "services/story-service/app": [
            "__init__.py",
            "main.py",
            "config.py"
        ],
        "services/story-service/app/api": [
            "__init__.py",
            "stories.py",
            "highlights.py",
            "viewers.py"
        ],
        "services/story-service/app/services": [
            "__init__.py",
            "story_service.py",
            "expiration_service.py",
            "viewer_service.py",
            "highlight_service.py",
            "close_friends_service.py"
        ],
        "services/story-service/app/models": [
            "__init__.py",
            "story.py",
            "highlight.py",
            "viewer.py"
        ],
        "services/story-service/app/db": [
            "__init__.py",
            "mongodb.py",
            "redis_client.py"
        ],
        "services/story-service/tests": [
            "__init__.py",
            "test_stories.py",
            "test_highlights.py"
        ],
        
        "services/messaging-service": [
            "build.sbt",
            "Dockerfile",
            ".dockerignore",
            "README.md"
        ],
        "services/messaging-service/src/main/scala/com/vignette/messaging": [
            "Main.scala",
            "Config.scala"
        ],
        "services/messaging-service/src/main/scala/com/vignette/messaging/actor": [
            "ChatActor.scala",
            "UserActor.scala",
            "GroupChatActor.scala",
            "MessageRouter.scala",
            "PresenceActor.scala",
            "TypingIndicatorActor.scala"
        ],
        "services/messaging-service/src/main/scala/com/vignette/messaging/api": [
            "ChatRoutes.scala",
            "MessageRoutes.scala",
            "WebSocketHandler.scala",
            "GroupChatRoutes.scala"
        ],
        "services/messaging-service/src/main/scala/com/vignette/messaging/model": [
            "Message.scala",
            "Chat.scala",
            "User.scala",
            "Presence.scala",
            "Reaction.scala"
        ],
        "services/messaging-service/src/main/scala/com/vignette/messaging/service": [
            "ChatService.scala",
            "MessageService.scala",
            "EncryptionService.scala",
            "PresenceService.scala",
            "VoiceMessageService.scala",
            "MediaMessageService.scala"
        ],
        "services/messaging-service/src/main/scala/com/vignette/messaging/repository": [
            "MessageRepository.scala",
            "ChatRepository.scala"
        ],
        "services/messaging-service/src/test/scala/com/vignette/messaging": [
            "ChatActorSpec.scala",
            "MessageServiceSpec.scala"
        ],
        
        "services/media-service": [
            "Cargo.toml",
            "Cargo.lock",
            "Dockerfile",
            ".dockerignore",
            "README.md"
        ],
        "services/media-service/src": [
            "main.rs",
            "config.rs"
        ],
        "services/media-service/src/handlers": [
            "mod.rs",
            "upload.rs",
            "download.rs",
            "streaming.rs",
            "processing.rs"
        ],
        "services/media-service/src/services": [
            "mod.rs",
            "image_processor.rs",
            "video_processor.rs",
            "audio_processor.rs",
            "thumbnail_generator.rs",
            "transcoding_service.rs",
            "compression_service.rs",
            "filter_service.rs",
            "ar_filter_service.rs"
        ],
        "services/media-service/src/storage": [
            "mod.rs",
            "s3_client.rs",
            "minio_client.rs",
            "local_storage.rs",
            "cdn_manager.rs"
        ],
        "services/media-service/src/models": [
            "mod.rs",
            "media.rs",
            "upload.rs",
            "metadata.rs"
        ],
        "services/media-service/src/utils": [
            "mod.rs",
            "validation.rs",
            "mime_types.rs",
            "crypto.rs"
        ],
        "services/media-service/tests": [
            "integration_test.rs",
            "processing_test.rs"
        ],
        
        "services/notification-service": [
            "build.sbt",
            "Dockerfile",
            ".dockerignore",
            "README.md"
        ],
        "services/notification-service/src/main/scala/com/vignette/notification": [
            "Main.scala",
            "Config.scala"
        ],
        "services/notification-service/src/main/scala/com/vignette/notification/actor": [
            "NotificationActor.scala",
            "PushNotificationActor.scala",
            "EmailActor.scala",
            "DeviceRegistry.scala",
            "ActivityActor.scala"
        ],
        "services/notification-service/src/main/scala/com/vignette/notification/api": [
            "NotificationRoutes.scala",
            "SubscriptionRoutes.scala",
            "ActivityRoutes.scala"
        ],
        "services/notification-service/src/main/scala/com/vignette/notification/model": [
            "Notification.scala",
            "Device.scala",
            "Template.scala",
            "Activity.scala"
        ],
        "services/notification-service/src/main/scala/com/vignette/notification/service": [
            "NotificationService.scala",
            "FCMService.scala",
            "APNService.scala",
            "EmailService.scala",
            "TemplateService.scala",
            "ActivityService.scala"
        ],
        "services/notification-service/src/main/scala/com/vignette/notification/repository": [
            "NotificationRepository.scala",
            "DeviceRepository.scala",
            "ActivityRepository.scala"
        ],
        
        "services/recommendation-service": [
            "requirements.txt",
            "Dockerfile",
            ".dockerignore",
            "README.md"
        ],
        "services/recommendation-service/app": [
            "__init__.py",
            "main.py",
            "config.py"
        ],
        "services/recommendation-service/app/api": [
            "__init__.py",
            "recommendations.py",
            "personalization.py",
            "explore.py"
        ],
        "services/recommendation-service/app/models": [
            "__init__.py",
            "user_embedding.py",
            "content_embedding.py",
            "interaction.py"
        ],
        "services/recommendation-service/app/ml": [
            "__init__.py",
            "collaborative_filtering.py",
            "content_based.py",
            "hybrid_model.py",
            "neural_cf.py",
            "matrix_factorization.py",
            "feature_engineering.py",
            "takes_recommender.py",
            "challenge_recommender.py"
        ],
        "services/recommendation-service/app/services": [
            "__init__.py",
            "recommendation_service.py",
            "training_service.py",
            "inference_service.py",
            "ab_testing_service.py",
            "explore_service.py"
        ],
        "services/recommendation-service/app/db": [
            "__init__.py",
            "vector_db.py",
            "cache.py"
        ],
        "services/recommendation-service/notebooks": [
            "model_training.ipynb",
            "data_exploration.ipynb",
            "takes_analysis.ipynb"
        ],
        
        "services/search-service": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            ".dockerignore",
            "README.md"
        ],
        "services/search-service/cmd/api": [
            "main.go"
        ],
        "services/search-service/internal/handler": [
            "search_handler.go",
            "indexing_handler.go",
            "autocomplete_handler.go",
            "hashtag_handler.go",
            "location_handler.go"
        ],
        "services/search-service/internal/service": [
            "search_service.go",
            "indexing_service.go",
            "ranking_service.go",
            "autocomplete_service.go",
            "hashtag_service.go",
            "location_service.go"
        ],
        "services/search-service/internal/elasticsearch": [
            "client.go",
            "indices.go",
            "queries.go"
        ],
        "services/search-service/internal/model": [
            "search.go",
            "document.go",
            "hashtag.go"
        ],
        
        "services/api-gateway": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            ".dockerignore",
            "README.md"
        ],
        "services/api-gateway/cmd/gateway": [
            "main.go"
        ],
        "services/api-gateway/internal/handler": [
            "proxy_handler.go",
            "health_handler.go"
        ],
        "services/api-gateway/internal/middleware": [
            "auth_middleware.go",
            "rate_limit_middleware.go",
            "cors_middleware.go",
            "logging_middleware.go",
            "circuit_breaker_middleware.go",
            "compression_middleware.go"
        ],
        "services/api-gateway/internal/proxy": [
            "router.go",
            "load_balancer.go",
            "service_discovery.go"
        ],
        "services/api-gateway/internal/config": [
            "config.go",
            "routes.yaml"
        ],
        
        "services/live-streaming-service": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            "README.md"
        ],
        "services/live-streaming-service/cmd/api": [
            "main.go"
        ],
        "services/live-streaming-service/internal/handler": [
            "stream_handler.go",
            "viewer_handler.go",
            "chat_handler.go",
            "gift_handler.go"
        ],
        "services/live-streaming-service/internal/service": [
            "streaming_service.go",
            "rtmp_service.go",
            "hls_service.go",
            "moderation_service.go",
            "viewer_service.go",
            "chat_service.go"
        ],
        "services/live-streaming-service/internal/webrtc": [
            "signaling.go",
            "peer_connection.go"
        ],
        
        "services/shopping-service": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            "README.md"
        ],
        "services/shopping-service/cmd/api": [
            "main.go"
        ],
        "services/shopping-service/internal/handler": [
            "product_handler.go",
            "shop_handler.go",
            "cart_handler.go",
            "checkout_handler.go",
            "order_handler.go",
            "wishlist_handler.go"
        ],
        "services/shopping-service/internal/service": [
            "product_service.go",
            "shop_service.go",
            "cart_service.go",
            "payment_service.go",
            "order_service.go",
            "wishlist_service.go"
        ],
        "services/shopping-service/internal/repository": [
            "product_repository.go",
            "order_repository.go",
            "shop_repository.go"
        ],
        
        "services/analytics-service": [
            "requirements.txt",
            "Dockerfile",
            "README.md"
        ],
        "services/analytics-service/app": [
            "__init__.py",
            "main.py",
            "config.py"
        ],
        "services/analytics-service/app/api": [
            "__init__.py",
            "metrics.py",
            "events.py",
            "insights.py"
        ],
        "services/analytics-service/app/services": [
            "__init__.py",
            "tracking_service.py",
            "aggregation_service.py",
            "reporting_service.py",
            "insights_service.py",
            "creator_analytics_service.py"
        ],
        "services/analytics-service/app/processors": [
            "__init__.py",
            "event_processor.py",
            "batch_processor.py",
            "real_time_processor.py"
        ],
        
        "services/moderation-service": [
            "requirements.txt",
            "Dockerfile",
            "README.md"
        ],
        "services/moderation-service/app": [
            "__init__.py",
            "main.py",
            "config.py"
        ],
        "services/moderation-service/app/api": [
            "__init__.py",
            "moderation.py",
            "reports.py"
        ],
        "services/moderation-service/app/services": [
            "__init__.py",
            "content_moderation_service.py",
            "report_service.py",
            "spam_detection_service.py",
            "auto_moderation_service.py"
        ],
        "services/moderation-service/app/ml": [
            "__init__.py",
            "toxicity_detector.py",
            "nsfw_detector.py",
            "spam_classifier.py",
            "hate_speech_detector.py",
            "violence_detector.py"
        ],
        
        "services/guide-service": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            "README.md"
        ],
        "services/guide-service/cmd/api": [
            "main.go"
        ],
        "services/guide-service/internal/handler": [
            "guide_handler.go",
            "collection_handler.go"
        ],
        "services/guide-service/internal/service": [
            "guide_service.go",
            "collection_service.go",
            "curation_service.go"
        ],
        "services/guide-service/internal/repository": [
            "guide_repository.go",
            "collection_repository.go"
        ],
        
        "services/creator-service": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            "README.md"
        ],
        "services/creator-service/cmd/api": [
            "main.go"
        ],
        "services/creator-service/internal/handler": [
            "creator_handler.go",
            "badge_handler.go",
            "subscription_handler.go",
            "monetization_handler.go"
        ],
        "services/creator-service/internal/service": [
            "creator_service.go",
            "badge_service.go",
            "subscription_service.go",
            "monetization_service.go",
            "payout_service.go"
        ],
        "services/creator-service/internal/repository": [
            "creator_repository.go",
            "subscription_repository.go",
            "payout_repository.go"
        ],
        
        "services/hashtag-service": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            "README.md"
        ],
        "services/hashtag-service/cmd/api": [
            "main.go"
        ],
        "services/hashtag-service/internal/handler": [
            "hashtag_handler.go",
            "trending_handler.go"
        ],
        "services/hashtag-service/internal/service": [
            "hashtag_service.go",
            "trending_service.go",
            "analytics_service.go"
        ],
        "services/hashtag-service/internal/repository": [
            "hashtag_repository.go"
        ],
        
        "services/location-service": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            "README.md"
        ],
        "services/location-service/cmd/api": [
            "main.go"
        ],
        "services/location-service/internal/handler": [
            "location_handler.go",
            "place_handler.go"
        ],
        "services/location-service/internal/service": [
            "location_service.go",
            "geocoding_service.go",
            "place_service.go"
        ],
        "services/location-service/internal/repository": [
            "location_repository.go"
        ],
        
        "infrastructure/kafka": [
            "docker-compose.yml",
            "topics.json",
            "README.md"
        ],
        "infrastructure/kafka/config": [
            "server.properties",
            "producer.properties",
            "consumer.properties"
        ],
        
        "infrastructure/postgres": [
            "docker-compose.yml",
            "init.sql",
            "README.md"
        ],
        "infrastructure/postgres/scripts": [
            "backup.sh",
            "restore.sh",
            "migrate.sh"
        ],
        
        "infrastructure/mongodb": [
            "docker-compose.yml",
            "init.js",
            "README.md"
        ],
        "infrastructure/mongodb/scripts": [
            "backup.sh",
            "restore.sh"
        ],
        
        "infrastructure/redis": [
            "docker-compose.yml",
            "redis.conf",
            "README.md"
        ],
        "infrastructure/redis/clusters": [
            "cluster.conf"
        ],
        
        "infrastructure/elasticsearch": [
            "docker-compose.yml",
            "elasticsearch.yml",
            "README.md"
        ],
        
        "infrastructure/minio": [
            "docker-compose.yml",
            "README.md"
        ],
        "infrastructure/minio/policies": [
            "public-read.json",
            "private.json"
        ],
        
        "infrastructure/nginx": [
            "nginx.conf",
            "Dockerfile",
            "README.md"
        ],
        "infrastructure/nginx/conf.d": [
            "gateway.conf",
            "ssl.conf",
            "cache.conf"
        ],
        
        "infrastructure/prometheus": [
            "prometheus.yml",
            "docker-compose.yml",
            "README.md"
        ],
        "infrastructure/prometheus/rules": [
            "alerts.yml",
            "recording_rules.yml"
        ],
        
        "infrastructure/grafana": [
            "docker-compose.yml",
            "grafana.ini",
            "README.md"
        ],
        "infrastructure/grafana/dashboards": [
            "api-gateway.json",
            "services-overview.json",
            "system-metrics.json",
            "takes-metrics.json",
            "user-engagement.json"
        ],
        
        "infrastructure/kubernetes": [
            "README.md"
        ],
        "infrastructure/kubernetes/base": [
            "namespace.yaml",
            "configmap.yaml",
            "secrets.yaml"
        ],
        "infrastructure/kubernetes/deployments": [
            "user-service.yaml",
            "feed-service.yaml",
            "takes-service.yaml",
            "media-service.yaml",
            "messaging-service.yaml",
            "notification-service.yaml",
            "api-gateway.yaml",
            "story-service.yaml",
            "shopping-service.yaml"
        ],
        "infrastructure/kubernetes/services": [
            "user-service.yaml",
            "feed-service.yaml",
            "takes-service.yaml",
            "api-gateway.yaml"
        ],
        "infrastructure/kubernetes/ingress": [
            "ingress.yaml",
            "tls-secret.yaml"
        ],
        "infrastructure/kubernetes/hpa": [
            "takes-service-hpa.yaml",
            "media-service-hpa.yaml",
            "api-gateway-hpa.yaml"
        ],
        
        "infrastructure/terraform": [
            "main.tf",
            "variables.tf",
            "outputs.tf",
            "provider.tf",
            "README.md"
        ],
        "infrastructure/terraform/modules/aws": [
            "eks.tf",
            "rds.tf",
            "s3.tf",
            "elasticache.tf",
            "cloudfront.tf"
        ],
        "infrastructure/terraform/modules/gcp": [
            "gke.tf",
            "cloud_sql.tf",
            "cloud_storage.tf",
            "memorystore.tf"
        ],
        
        "shared/proto": [
            "user.proto",
            "post.proto",
            "take.proto",
            "message.proto",
            "notification.proto",
            "media.proto",
            "challenge.proto"
        ],
        
        "shared/events": [
            "user_events.json",
            "post_events.json",
            "take_events.json",
            "message_events.json",
            "challenge_events.json",
            "engagement_events.json"
        ],
        
        "scripts": [
            "setup.sh",
            "start-all.sh",
            "stop-all.sh",
            "deploy.sh",
            "migrate.sh",
            "seed-data.sh",
            "health-check.sh",
            "backup.sh",
            "restore.sh"
        ],
        
        "scripts/kafka": [
            "create-topics.sh",
            "consume-events.sh",
            "produce-test-event.sh"
        ],
        
        "scripts/monitoring": [
            "setup-prometheus.sh",
            "setup-grafana.sh",
            "setup-alerts.sh"
        ],
        
        "scripts/deployment": [
            "deploy-production.sh",
            "deploy-staging.sh",
            "rollback.sh",
            "scale-services.sh"
        ],
        
        "docs": [
            "README.md",
            "ARCHITECTURE.md",
            "API_DOCUMENTATION.md",
            "DEPLOYMENT.md",
            "CONTRIBUTING.md",
            "MICROSERVICES.md",
            "EVENT_DRIVEN.md",
            "SCALING.md",
            "SECURITY.md",
            "TAKES_SYSTEM.md",
            "CHALLENGE_SYSTEM.md"
        ],
        
        "docs/api": [
            "user-service.md",
            "feed-service.md",
            "takes-service.md",
            "messaging-service.md",
            "media-service.md",
            "shopping-service.md",
            "creator-service.md"
        ],
        
        "tests/integration": [
            "user_flow_test.py",
            "post_creation_flow_test.py",
            "take_creation_flow_test.py",
            "messaging_flow_test.py",
            "challenge_flow_test.py",
            "shopping_flow_test.py"
        ],
        
        "tests/load": [
            "locustfile.py",
            "k6-script.js",
            "takes-load-test.js",
            "README.md"
        ],
        
        "tests/e2e": [
            "feed_test.py",
            "takes_test.py",
            "stories_test.py",
            "direct_test.py"
        ],
        
        ".github/workflows": [
            "ci.yml",
            "cd.yml",
            "user-service.yml",
            "feed-service.yml",
            "takes-service.yml",
            "media-service.yml",
            "security-scan.yml"
        ],
        
        "monitoring/alerts": [
            "service-down.yml",
            "high-latency.yml",
            "error-rate.yml",
            "takes-processing-delay.yml",
            "storage-capacity.yml"
        ],
        
        "config": [
            "development.env",
            "staging.
