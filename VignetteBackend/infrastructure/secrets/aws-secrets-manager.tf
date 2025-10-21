# AWS Secrets Manager Configuration
# Alternative to Vault for AWS-native deployments
# Automatically rotates secrets, integrates with RDS, encrypts at rest

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

variable "aws_region" {
  description = "AWS region"
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment (production, staging, development)"
  default     = "production"
}

# KMS key for encryption
resource "aws_kms_key" "secrets" {
  description             = "KMS key for Entativa secrets encryption"
  deletion_window_in_days = 30
  enable_key_rotation     = true

  tags = {
    Name        = "entativa-secrets-key"
    Environment = var.environment
  }
}

resource "aws_kms_alias" "secrets" {
  name          = "alias/entativa-secrets"
  target_key_id = aws_kms_key.secrets.key_id
}

# PostgreSQL credentials
resource "aws_secretsmanager_secret" "postgres" {
  name                    = "entativa/${var.environment}/database/postgres"
  description             = "PostgreSQL database credentials"
  kms_key_id              = aws_kms_key.secrets.arn
  recovery_window_in_days = 30

  tags = {
    Environment = var.environment
    Service     = "database"
  }
}

resource "aws_secretsmanager_secret_version" "postgres" {
  secret_id = aws_secretsmanager_secret.postgres.id
  secret_string = jsonencode({
    host     = "postgres.entativa.com"
    port     = 5432
    database = "entativa"
    username = "entativa_app"
    password = random_password.postgres.result
  })
}

# S3 credentials (use IAM roles instead in production)
resource "aws_secretsmanager_secret" "s3" {
  name                    = "entativa/${var.environment}/s3/credentials"
  description             = "S3 access credentials"
  kms_key_id              = aws_kms_key.secrets.arn
  recovery_window_in_days = 30

  tags = {
    Environment = var.environment
    Service     = "media"
  }
}

resource "aws_secretsmanager_secret_version" "s3" {
  secret_id = aws_secretsmanager_secret.s3.id
  secret_string = jsonencode({
    access_key_id     = aws_iam_access_key.s3_user.id
    secret_access_key = aws_iam_access_key.s3_user.secret
    bucket_name       = aws_s3_bucket.media.id
    region            = var.aws_region
  })
}

# Elasticsearch credentials
resource "aws_secretsmanager_secret" "elasticsearch" {
  name                    = "entativa/${var.environment}/elasticsearch/credentials"
  description             = "Elasticsearch credentials"
  kms_key_id              = aws_kms_key.secrets.arn
  recovery_window_in_days = 30

  tags = {
    Environment = var.environment
    Service     = "search"
  }
}

resource "aws_secretsmanager_secret_version" "elasticsearch" {
  secret_id = aws_secretsmanager_secret.elasticsearch.id
  secret_string = jsonencode({
    host     = aws_elasticsearch_domain.search.endpoint
    port     = 443
    username = "elastic"
    password = random_password.elasticsearch.result
    api_key  = random_password.elasticsearch_api.result
  })
}

# Redis credentials
resource "aws_secretsmanager_secret" "redis" {
  name                    = "entativa/${var.environment}/redis"
  description             = "Redis credentials"
  kms_key_id              = aws_kms_key.secrets.arn
  recovery_window_in_days = 30

  tags = {
    Environment = var.environment
    Service     = "cache"
  }
}

resource "aws_secretsmanager_secret_version" "redis" {
  secret_id = aws_secretsmanager_secret.redis.id
  secret_string = jsonencode({
    host     = aws_elasticache_replication_group.redis.primary_endpoint_address
    port     = 6379
    password = random_password.redis.result
  })
}

# JWT signing keys
resource "aws_secretsmanager_secret" "jwt" {
  name                    = "entativa/${var.environment}/jwt"
  description             = "JWT signing keys"
  kms_key_id              = aws_kms_key.secrets.arn
  recovery_window_in_days = 30

  tags = {
    Environment = var.environment
    Service     = "auth"
  }
}

resource "tls_private_key" "jwt" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_secretsmanager_secret_version" "jwt" {
  secret_id = aws_secretsmanager_secret.jwt.id
  secret_string = jsonencode({
    private_key = tls_private_key.jwt.private_key_pem
    public_key  = tls_private_key.jwt.public_key_pem
  })
}

# Email SMTP credentials
resource "aws_secretsmanager_secret" "email" {
  name                    = "entativa/${var.environment}/email/smtp"
  description             = "Email SMTP credentials"
  kms_key_id              = aws_kms_key.secrets.arn
  recovery_window_in_days = 30

  tags = {
    Environment = var.environment
    Service     = "email"
  }
}

resource "aws_secretsmanager_secret_version" "email" {
  secret_id = aws_secretsmanager_secret.email.id
  secret_string = jsonencode({
    host         = "smtp.sendgrid.net"
    port         = 587
    username     = "apikey"
    password     = var.sendgrid_api_key
    from_address = "noreply@entativa.com"
  })

  lifecycle {
    ignore_changes = [secret_string]
  }
}

# Stripe credentials
resource "aws_secretsmanager_secret" "stripe" {
  name                    = "entativa/${var.environment}/stripe"
  description             = "Stripe API credentials"
  kms_key_id              = aws_kms_key.secrets.arn
  recovery_window_in_days = 30

  tags = {
    Environment = var.environment
    Service     = "payment"
  }
}

resource "aws_secretsmanager_secret_version" "stripe" {
  secret_id = aws_secretsmanager_secret.stripe.id
  secret_string = jsonencode({
    api_key        = var.stripe_api_key
    webhook_secret = var.stripe_webhook_secret
  })

  lifecycle {
    ignore_changes = [secret_string]
  }
}

# Messaging/Signal keys
resource "aws_secretsmanager_secret" "messaging" {
  name                    = "entativa/${var.environment}/messaging/signal"
  description             = "Signal Protocol server keys"
  kms_key_id              = aws_kms_key.secrets.arn
  recovery_window_in_days = 30

  tags = {
    Environment = var.environment
    Service     = "messaging"
  }
}

resource "aws_secretsmanager_secret_version" "messaging" {
  secret_id = aws_secretsmanager_secret.messaging.id
  secret_string = jsonencode({
    server_private_key = random_password.signal_key.result
  })
}

# Backup encryption salt
resource "aws_secretsmanager_secret" "backup" {
  name                    = "entativa/${var.environment}/messaging/backup"
  description             = "Message backup encryption salt"
  kms_key_id              = aws_kms_key.secrets.arn
  recovery_window_in_days = 30

  tags = {
    Environment = var.environment
    Service     = "messaging"
  }
}

resource "aws_secretsmanager_secret_version" "backup" {
  secret_id = aws_secretsmanager_secret.backup.id
  secret_string = jsonencode({
    encryption_salt = random_password.backup_salt.result
  })
}

# Random passwords
resource "random_password" "postgres" {
  length  = 32
  special = true
}

resource "random_password" "elasticsearch" {
  length  = 32
  special = true
}

resource "random_password" "elasticsearch_api" {
  length  = 64
  special = false
}

resource "random_password" "redis" {
  length  = 32
  special = true
}

resource "random_password" "signal_key" {
  length  = 64
  special = false
}

resource "random_password" "backup_salt" {
  length  = 32
  special = false
}

# IAM policy for services to access secrets
resource "aws_iam_policy" "secrets_read" {
  name        = "entativa-secrets-read"
  description = "Allow reading Entativa secrets from AWS Secrets Manager"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue",
          "secretsmanager:DescribeSecret"
        ]
        Resource = [
          aws_secretsmanager_secret.postgres.arn,
          aws_secretsmanager_secret.s3.arn,
          aws_secretsmanager_secret.elasticsearch.arn,
          aws_secretsmanager_secret.redis.arn,
          aws_secretsmanager_secret.jwt.arn,
          aws_secretsmanager_secret.email.arn,
          aws_secretsmanager_secret.stripe.arn,
          aws_secretsmanager_secret.messaging.arn,
          aws_secretsmanager_secret.backup.arn
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "kms:Decrypt",
          "kms:DescribeKey"
        ]
        Resource = aws_kms_key.secrets.arn
      }
    ]
  })
}

# IAM role for user service
resource "aws_iam_role" "user_service" {
  name = "entativa-user-service"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ecs.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "user_service_secrets" {
  role       = aws_iam_role.user_service.name
  policy_arn = aws_iam_policy.secrets_read.arn
}

# Automatic secret rotation (for RDS)
resource "aws_secretsmanager_secret_rotation" "postgres" {
  secret_id           = aws_secretsmanager_secret.postgres.id
  rotation_lambda_arn = aws_lambda_function.rotate_postgres.arn

  rotation_rules {
    automatically_after_days = 30
  }
}

# Outputs
output "kms_key_id" {
  description = "KMS key ID for secrets encryption"
  value       = aws_kms_key.secrets.id
}

output "secrets_policy_arn" {
  description = "IAM policy ARN for reading secrets"
  value       = aws_iam_policy.secrets_read.arn
}
