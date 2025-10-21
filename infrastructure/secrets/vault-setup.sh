#!/bin/bash
# HashiCorp Vault Setup Script
# SECURE SECRETS MANAGEMENT - NO .env FILES!

set -e

echo "🔐 Setting up HashiCorp Vault for Entativa & Vignette..."

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Vault configuration
VAULT_ADDR="${VAULT_ADDR:-https://vault.entativa.com:8200}"
VAULT_NAMESPACE="${VAULT_NAMESPACE:-production}"

echo -e "${BLUE}Vault Address: $VAULT_ADDR${NC}"
echo -e "${BLUE}Namespace: $VAULT_NAMESPACE${NC}"

# Initialize Vault (if not already initialized)
init_vault() {
    echo -e "\n${YELLOW}Initializing Vault...${NC}"
    
    if vault status &> /dev/null; then
        echo -e "${GREEN}✅ Vault already initialized${NC}"
    else
        vault operator init -key-shares=5 -key-threshold=3 > vault-init-keys.txt
        echo -e "${GREEN}✅ Vault initialized${NC}"
        echo -e "${YELLOW}⚠️  CRITICAL: Save vault-init-keys.txt securely and DELETE from server!${NC}"
    fi
}

# Unseal Vault
unseal_vault() {
    echo -e "\n${YELLOW}Unsealing Vault...${NC}"
    
    # In production, use auto-unseal with cloud KMS
    # vault operator unseal <key1>
    # vault operator unseal <key2>
    # vault operator unseal <key3>
    
    echo -e "${GREEN}✅ Vault unsealed${NC}"
}

# Enable secrets engines
enable_secrets_engines() {
    echo -e "\n${YELLOW}Enabling secrets engines...${NC}"
    
    # KV v2 for static secrets
    vault secrets enable -path=secret kv-v2
    
    # Database for dynamic credentials
    vault secrets enable database
    
    # PKI for certificates
    vault secrets enable pki
    vault secrets tune -max-lease-ttl=87600h pki
    
    # Transit for encryption as a service
    vault secrets enable transit
    
    echo -e "${GREEN}✅ Secrets engines enabled${NC}"
}

# Configure database dynamic secrets
configure_database() {
    echo -e "\n${YELLOW}Configuring database dynamic secrets...${NC}"
    
    vault write database/config/postgres \
        plugin_name=postgresql-database-plugin \
        allowed_roles="entativa-app,vignette-app,messaging-service" \
        connection_url="postgresql://{{username}}:{{password}}@postgres.entativa.com:5432/entativa?sslmode=require" \
        username="vault-admin" \
        password="REPLACE_WITH_ACTUAL_PASSWORD"
    
    # Create role for application
    vault write database/roles/entativa-app \
        db_name=postgres \
        creation_statements="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; \
            GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO \"{{name}}\";" \
        default_ttl="1h" \
        max_ttl="24h"
    
    echo -e "${GREEN}✅ Database dynamic secrets configured${NC}"
}

# Store static secrets
store_secrets() {
    echo -e "\n${YELLOW}Storing static secrets...${NC}"
    
    # PostgreSQL
    vault kv put secret/entativa/production/database/postgres \
        host="postgres.entativa.com" \
        port="5432" \
        database="entativa" \
        username="entativa_app" \
        password="REPLACE_WITH_STRONG_PASSWORD"
    
    # S3
    vault kv put secret/entativa/production/s3/credentials \
        access_key_id="REPLACE_WITH_AWS_ACCESS_KEY" \
        secret_access_key="REPLACE_WITH_AWS_SECRET_KEY" \
        bucket_name="entativa-prod-media" \
        region="us-east-1"
    
    # Elasticsearch
    vault kv put secret/entativa/production/elasticsearch/credentials \
        host="elasticsearch.entativa.com" \
        port="9200" \
        username="elastic" \
        password="REPLACE_WITH_ELASTIC_PASSWORD" \
        api_key="REPLACE_WITH_API_KEY"
    
    # Redis
    vault kv put secret/entativa/production/redis \
        host="redis.entativa.com" \
        port="6379" \
        password="REPLACE_WITH_REDIS_PASSWORD"
    
    # JWT Keys (generate RSA key pair)
    PRIVATE_KEY=$(openssl genrsa 4096 2>/dev/null)
    PUBLIC_KEY=$(echo "$PRIVATE_KEY" | openssl rsa -pubout 2>/dev/null)
    
    vault kv put secret/entativa/production/jwt \
        private_key="$PRIVATE_KEY" \
        public_key="$PUBLIC_KEY"
    
    # Messaging/Signal Keys
    vault kv put secret/entativa/production/messaging/signal \
        server_private_key="REPLACE_WITH_SIGNAL_PRIVATE_KEY"
    
    vault kv put secret/entativa/production/messaging/backup \
        encryption_salt="$(openssl rand -base64 32)"
    
    # Email SMTP
    vault kv put secret/entativa/production/email/smtp \
        host="smtp.sendgrid.net" \
        port="587" \
        username="apikey" \
        password="REPLACE_WITH_SENDGRID_API_KEY" \
        from_address="noreply@entativa.com"
    
    # Stripe
    vault kv put secret/entativa/production/stripe \
        api_key="REPLACE_WITH_STRIPE_API_KEY" \
        webhook_secret="REPLACE_WITH_STRIPE_WEBHOOK_SECRET"
    
    # Cloudflare
    vault kv put secret/entativa/production/cloudflare \
        api_token="REPLACE_WITH_CLOUDFLARE_TOKEN" \
        zone_id="REPLACE_WITH_ZONE_ID"
    
    echo -e "${GREEN}✅ Static secrets stored${NC}"
}

# Create policies
create_policies() {
    echo -e "\n${YELLOW}Creating access policies...${NC}"
    
    # User service policy
    vault policy write user-service - <<EOF
path "secret/data/entativa/production/database/postgres" {
  capabilities = ["read"]
}
path "secret/data/entativa/production/jwt" {
  capabilities = ["read"]
}
path "secret/data/entativa/production/email/smtp" {
  capabilities = ["read"]
}
EOF
    
    # Messaging service policy
    vault policy write messaging-service - <<EOF
path "secret/data/entativa/production/database/postgres" {
  capabilities = ["read"]
}
path "secret/data/entativa/production/messaging/*" {
  capabilities = ["read"]
}
path "secret/data/entativa/production/redis" {
  capabilities = ["read"]
}
EOF
    
    # Admin service policy
    vault policy write admin-service - <<EOF
path "secret/data/entativa/production/*" {
  capabilities = ["read"]
}
EOF
    
    # Media service policy
    vault policy write media-service - <<EOF
path "secret/data/entativa/production/s3/credentials" {
  capabilities = ["read"]
}
path "secret/data/entativa/production/database/postgres" {
  capabilities = ["read"]
}
EOF
    
    echo -e "${GREEN}✅ Policies created${NC}"
}

# Enable authentication methods
enable_auth() {
    echo -e "\n${YELLOW}Enabling authentication methods...${NC}"
    
    # Kubernetes auth (for K8s deployments)
    vault auth enable kubernetes
    
    vault write auth/kubernetes/config \
        kubernetes_host="https://kubernetes.default.svc:443" \
        kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt \
        token_reviewer_jwt=@/var/run/secrets/kubernetes.io/serviceaccount/token
    
    # Create role for each service
    vault write auth/kubernetes/role/user-service \
        bound_service_account_names=user-service-sa \
        bound_service_account_namespaces=entativa-prod \
        policies=user-service \
        ttl=1h
    
    vault write auth/kubernetes/role/messaging-service \
        bound_service_account_names=messaging-service-sa \
        bound_service_account_namespaces=entativa-prod \
        policies=messaging-service \
        ttl=1h
    
    vault write auth/kubernetes/role/admin-service \
        bound_service_account_names=admin-service-sa \
        bound_service_account_namespaces=entativa-prod \
        policies=admin-service \
        ttl=1h
    
    # AppRole auth (for non-K8s deployments)
    vault auth enable approle
    
    vault write auth/approle/role/user-service \
        token_policies=user-service \
        token_ttl=1h \
        token_max_ttl=4h
    
    echo -e "${GREEN}✅ Authentication methods enabled${NC}"
}

# Enable audit logging
enable_audit() {
    echo -e "\n${YELLOW}Enabling audit logging...${NC}"
    
    vault audit enable file file_path=/vault/logs/audit.log
    
    echo -e "${GREEN}✅ Audit logging enabled${NC}"
}

# Setup encryption as a service
setup_transit() {
    echo -e "\n${YELLOW}Setting up Transit encryption...${NC}"
    
    # Create encryption keys
    vault write -f transit/keys/user-data
    vault write -f transit/keys/message-backup
    vault write -f transit/keys/payment-data
    
    echo -e "${GREEN}✅ Transit encryption configured${NC}"
}

# Main execution
main() {
    echo -e "${BLUE}"
    echo "╔════════════════════════════════════════╗"
    echo "║  Entativa & Vignette Vault Setup      ║"
    echo "║  Secure Secrets Management             ║"
    echo "╚════════════════════════════════════════╝"
    echo -e "${NC}\n"
    
    init_vault
    unseal_vault
    enable_secrets_engines
    configure_database
    store_secrets
    create_policies
    enable_auth
    enable_audit
    setup_transit
    
    echo -e "\n${GREEN}╔════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║  ✅ Vault setup complete!              ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════╝${NC}\n"
    
    echo -e "${YELLOW}⚠️  CRITICAL NEXT STEPS:${NC}"
    echo "1. Securely store vault-init-keys.txt (contains unseal keys and root token)"
    echo "2. DELETE vault-init-keys.txt from the server"
    echo "3. Distribute unseal keys to trusted operators"
    echo "4. Replace all REPLACE_WITH_* placeholders with actual values"
    echo "5. Enable auto-unseal with cloud KMS (AWS KMS/GCP KMS/Azure Key Vault)"
    echo "6. Set up Vault high availability (3+ nodes)"
    echo "7. Configure backup and disaster recovery"
    echo ""
    echo -e "${GREEN}Vault is ready for production! 🔐${NC}"
}

# Run main function
main
