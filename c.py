import os

def create_socialink_backend():
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
            "friends_handler.go",
            "blocking_handler.go",
            "privacy_handler.go"
        ],
        "services/user-service/internal/service": [
            "auth_service.go",
            "user_service.go",
            "profile_service.go",
            "friends_service.go",
            "session_service.go",
            "verification_service.go"
        ],
        "services/user-service/internal/repository": [
            "user_repository.go",
            "profile_repository.go",
            "friends_repository.go",
            "session_repository.go"
        ],
        "services/user-service/internal/model": [
            "user.go",
            "profile.go",
            "friend.go",
            "session.go",
            "privacy.go"
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
            "002_create_profiles_table.down.sql"
        ],
        "services/user-service/test": [
            "auth_test.go",
            "user_test.go",
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
            "reactions.py",
            "shares.py",
            "polls.py"
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
            "reaction.py",
            "poll.py"
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
            "personalization_service.py"
        ],
        "services/feed-service/app/ml": [
            "__init__.py",
            "content_ranker.py",
            "engagement_predictor.py",
            "recommendation_engine.py",
            "embeddings.py"
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
        
        "services/messaging-service": [
            "build.sbt",
            "Dockerfile",
            ".dockerignore",
            "README.md"
        ],
        "services/messaging-service/src/main/scala/com/socialink/messaging": [
            "Main.scala",
            "Config.scala"
        ],
        "services/messaging-service/src/main/scala/com/socialink/messaging/actor": [
            "ChatActor.scala",
            "UserActor.scala",
            "GroupChatActor.scala",
            "MessageRouter.scala",
            "PresenceActor.scala"
        ],
        "services/messaging-service/src/main/scala/com/socialink/messaging/api": [
            "ChatRoutes.scala",
            "MessageRoutes.scala",
            "WebSocketHandler.scala"
        ],
        "services/messaging-service/src/main/scala/com/socialink/messaging/model": [
            "Message.scala",
            "Chat.scala",
            "User.scala",
            "Presence.scala"
        ],
        "services/messaging-service/src/main/scala/com/socialink/messaging/service": [
            "ChatService.scala",
            "MessageService.scala",
            "EncryptionService.scala",
            "PresenceService.scala"
        ],
        "services/messaging-service/src/main/scala/com/socialink/messaging/repository": [
            "MessageRepository.scala",
            "ChatRepository.scala"
        ],
        "services/messaging-service/src/test/scala/com/socialink/messaging": [
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
            "compression_service.rs"
        ],
        "services/media-service/src/storage": [
            "mod.rs",
            "s3_client.rs",
            "minio_client.rs",
            "local_storage.rs"
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
        "services/notification-service/src/main/scala/com/socialink/notification": [
            "Main.scala",
            "Config.scala"
        ],
        "services/notification-service/src/main/scala/com/socialink/notification/actor": [
            "NotificationActor.scala",
            "PushNotificationActor.scala",
            "EmailActor.scala",
            "SMSActor.scala",
            "DeviceRegistry.scala"
        ],
        "services/notification-service/src/main/scala/com/socialink/notification/api": [
            "NotificationRoutes.scala",
            "SubscriptionRoutes.scala"
        ],
        "services/notification-service/src/main/scala/com/socialink/notification/model": [
            "Notification.scala",
            "Device.scala",
            "Template.scala"
        ],
        "services/notification-service/src/main/scala/com/socialink/notification/service": [
            "NotificationService.scala",
            "FCMService.scala",
            "APNService.scala",
            "EmailService.scala",
            "SMSService.scala",
            "TemplateService.scala"
        ],
        "services/notification-service/src/main/scala/com/socialink/notification/repository": [
            "NotificationRepository.scala",
            "DeviceRepository.scala"
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
            "personalization.py"
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
            "feature_engineering.py"
        ],
        "services/recommendation-service/app/services": [
            "__init__.py",
            "recommendation_service.py",
            "training_service.py",
            "inference_service.py",
            "ab_testing_service.py"
        ],
        "services/recommendation-service/app/db": [
            "__init__.py",
            "vector_db.py",
            "cache.py"
        ],
        "services/recommendation-service/notebooks": [
            "model_training.ipynb",
            "data_exploration.ipynb"
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
            "autocomplete_handler.go"
        ],
        "services/search-service/internal/service": [
            "search_service.go",
            "indexing_service.go",
            "ranking_service.go",
            "autocomplete_service.go"
        ],
        "services/search-service/internal/elasticsearch": [
            "client.go",
            "indices.go",
            "queries.go"
        ],
        "services/search-service/internal/model": [
            "search.go",
            "document.go"
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
            "circuit_breaker_middleware.go"
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
        
        "services/group-service": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            "README.md"
        ],
        "services/group-service/cmd/api": [
            "main.go"
        ],
        "services/group-service/internal/handler": [
            "group_handler.go",
            "member_handler.go",
            "post_handler.go"
        ],
        "services/group-service/internal/service": [
            "group_service.go",
            "member_service.go",
            "moderation_service.go"
        ],
        "services/group-service/internal/repository": [
            "group_repository.go",
            "member_repository.go"
        ],
        "services/group-service/internal/model": [
            "group.go",
            "member.go"
        ],
        
        "services/event-service": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            "README.md"
        ],
        "services/event-service/cmd/api": [
            "main.go"
        ],
        "services/event-service/internal/handler": [
            "event_handler.go",
            "rsvp_handler.go",
            "reminder_handler.go"
        ],
        "services/event-service/internal/service": [
            "event_service.go",
            "rsvp_service.go",
            "reminder_service.go",
            "calendar_service.go"
        ],
        "services/event-service/internal/repository": [
            "event_repository.go",
            "rsvp_repository.go"
        ],
        
        "services/marketplace-service": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            "README.md"
        ],
        "services/marketplace-service/cmd/api": [
            "main.go"
        ],
        "services/marketplace-service/internal/handler": [
            "listing_handler.go",
            "transaction_handler.go",
            "review_handler.go"
        ],
        "services/marketplace-service/internal/service": [
            "listing_service.go",
            "transaction_service.go",
            "payment_service.go",
            "shipping_service.go"
        ],
        "services/marketplace-service/internal/repository": [
            "listing_repository.go",
            "transaction_repository.go"
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
            "chat_handler.go"
        ],
        "services/live-streaming-service/internal/service": [
            "streaming_service.go",
            "rtmp_service.go",
            "hls_service.go",
            "moderation_service.go"
        ],
        "services/live-streaming-service/internal/webrtc": [
            "signaling.go",
            "peer_connection.go"
        ],
        
        "services/story-service": [
            "requirements.txt",
            "Dockerfile",
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
            "highlights.py"
        ],
        "services/story-service/app/services": [
            "__init__.py",
            "story_service.py",
            "expiration_service.py",
            "viewer_service.py"
        ],
        "services/story-service/app/models": [
            "__init__.py",
            "story.py",
            "highlight.py"
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
            "events.py"
        ],
        "services/analytics-service/app/services": [
            "__init__.py",
            "tracking_service.py",
            "aggregation_service.py",
            "reporting_service.py"
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
            "spam_detection_service.py"
        ],
        "services/moderation-service/app/ml": [
            "__init__.py",
            "toxicity_detector.py",
            "nsfw_detector.py",
            "spam_classifier.py",
            "hate_speech_detector.py"
        ],
        
        "services/payment-service": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            "README.md"
        ],
        "services/payment-service/cmd/api": [
            "main.go"
        ],
        "services/payment-service/internal/handler": [
            "payment_handler.go",
            "subscription_handler.go",
            "webhook_handler.go"
        ],
        "services/payment-service/internal/service": [
            "payment_service.go",
            "subscription_service.go",
            "stripe_service.go",
            "paypal_service.go"
        ],
        "services/payment-service/internal/repository": [
            "transaction_repository.go",
            "subscription_repository.go"
        ],
        
        "services/ad-service": [
            "requirements.txt",
            "Dockerfile",
            "README.md"
        ],
        "services/ad-service/app": [
            "__init__.py",
            "main.py",
            "config.py"
        ],
        "services/ad-service/app/api": [
            "__init__.py",
            "ads.py",
            "campaigns.py"
        ],
        "services/ad-service/app/services": [
            "__init__.py",
            "ad_service.py",
            "targeting_service.py",
            "bidding_service.py",
            "impression_service.py"
        ],
        "services/ad-service/app/ml": [
            "__init__.py",
            "targeting_model.py",
            "ctr_predictor.py"
        ],
        
        "services/page-service": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            "README.md"
        ],
        "services/page-service/cmd/api": [
            "main.go"
        ],
        "services/page-service/internal/handler": [
            "page_handler.go",
            "follower_handler.go"
        ],
        "services/page-service/internal/service": [
            "page_service.go",
            "follower_service.go",
            "insights_service.go"
        ],
        
        "services/job-service": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            "README.md"
        ],
        "services/job-service/cmd/api": [
            "main.go"
        ],
        "services/job-service/internal/handler": [
            "job_handler.go",
            "application_handler.go"
        ],
        "services/job-service/internal/service": [
            "job_service.go",
            "application_service.go",
            "matching_service.go"
        ],
        
        "services/dating-service": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            "README.md"
        ],
        "services/dating-service/cmd/api": [
            "main.go"
        ],
        "services/dating-service/internal/handler": [
            "profile_handler.go",
            "match_handler.go",
            "swipe_handler.go"
        ],
        "services/dating-service/internal/service": [
            "matching_service.go",
            "profile_service.go",
            "algorithm_service.go"
        ],
        
        "services/gaming-service": [
            "go.mod",
            "go.sum",
            "Dockerfile",
            "README.md"
        ],
        "services/gaming-service/cmd/api": [
            "main.go"
        ],
        "services/gaming-service/internal/handler": [
            "game_handler.go",
            "leaderboard_handler.go",
            "achievement_handler.go"
        ],
        "services/gaming-service/internal/service": [
            "game_service.go",
            "leaderboard_service.go",
            "achievement_service.go"
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
            "restore.sh"
        ],
        
        "infrastructure/mongodb": [
            "docker-compose.yml",
            "init.js",
            "README.md"
        ],
        
        "infrastructure/redis": [
            "docker-compose.yml",
            "redis.conf",
            "README.md"
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
        
        "infrastructure/nginx": [
            "nginx.conf",
            "Dockerfile",
            "README.md"
        ],
        "infrastructure/nginx/conf.d": [
            "gateway.conf",
            "ssl.conf"
        ],
        
        "infrastructure/prometheus": [
            "prometheus.yml",
            "docker-compose.yml",
            "README.md"
        ],
        "infrastructure/prometheus/rules": [
            "alerts.yml"
        ],
        
        "infrastructure/grafana": [
            "docker-compose.yml",
            "grafana.ini",
            "README.md"
        ],
        "infrastructure/grafana/dashboards": [
            "api-gateway.json",
            "services-overview.json",
            "system-metrics.json"
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
            "media-service.yaml",
            "messaging-service.yaml",
            "notification-service.yaml",
            "api-gateway.yaml"
        ],
        "infrastructure/kubernetes/services": [
            "user-service.yaml",
            "feed-service.yaml",
            "api-gateway.yaml"
        ],
        "infrastructure/kubernetes/ingress": [
            "ingress.yaml",
            "tls-secret.yaml"
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
            "elasticache.tf"
        ],
        
        "shared/proto": [
            "user.proto",
            "post.proto",
            "message.proto",
            "notification.proto",
            "media.proto"
        ],
        
        "shared/events": [
            "user_events.json",
            "post_events.json",
            "message_events.json"
        ],
        
        "scripts": [
            "setup.sh",
            "start-all.sh",
            "stop-all.sh",
            "deploy.sh",
            "migrate.sh",
            "seed-data.sh",
            "health-check.sh"
        ],
        
        "scripts/kafka": [
            "create-topics.sh",
            "consume-events.sh",
            "produce-test-event.sh"
        ],
        
        "scripts/monitoring": [
            "setup-prometheus.sh",
            "setup-grafana.sh"
        ],
        
        "docs": [
            "README.md",
            "ARCHITECTURE.md",
            "API_DOCUMENTATION.md",
            "DEPLOYMENT.md",
            "CONTRIBUTING.md",
            "MICROSERVICES.md",
            "EVENT_DRIVEN.md",
            "SCALING.md"
        ],
        
        "docs/api": [
            "user-service.md",
            "feed-service.md",
            "messaging-service.md",
            "media-service.md"
        ],
        
        "tests/integration": [
            "user_flow_test.py",
            "post_creation_flow_test.py",
            "messaging_flow_test.py"
        ],
        
        "tests/load": [
            "locustfile.py",
            "k6-script.js",
            "README.md"
        ],
        
        ".github/workflows": [
            "ci.yml",
            "cd.yml",
            "user-service.yml",
            "feed-service.yml",
            "media-service.yml"
        ],
        
        "monitoring/alerts": [
            "service-down.yml",
            "high-latency.yml",
            "error-rate.yml"
        ],
        
        "config": [
            "development.env",
            "staging.env",
            "production.env"
        ]
    }
    
    base_dir = "SocialinkBackend"
    
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
    
    print(f"\n✅ Socialink Backend microservices structure created successfully in '{base_dir}/' directory!")
    print(f"🚀 Architecture Overview:")
    print(f"   - Go Services: User, Search, API Gateway, Group, Event, Marketplace, Live Streaming, Payment, Page, Job, Dating, Gaming")
    print(f"   - Python Services: Feed, Recommendation, Story, Analytics, Moderation, Ad")
    print(f"   - Scala Services: Messaging, Notification (Akka-based)")
    print(f"   - Rust Services: Media (high-performance processing)")
    print(f"\n📦 Infrastructure:")
    print(f"   - Kafka for event streaming")
    print(f"   - PostgreSQL & MongoDB for databases")
    print(f"   - Redis for caching")
    print(f"   - Elasticsearch for search")
    print(f"   - MinIO/S3 for object storage")
    print(f"   - Prometheus & Grafana for monitoring")
    print(f"   - Kubernetes manifests & Terraform for IaC")
    print(f"\n💡 Total services: 18 microservices")
    print(f"📁 Total files/folders created: {sum(len(files) for files in structure.values()) + len(root_files)} items")

if __name__ == "__main__":
    create_socialink_backend()
