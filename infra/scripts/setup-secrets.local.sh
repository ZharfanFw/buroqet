#!/bin/bash
# =============================================================================
# setup-secrets.sh — Buat semua Kubernetes Secrets untuk Buroqet
# =============================================================================
# Isi semua nilai di bawah sebelum menjalankan script ini.
# File ini TIDAK BOLEH di-commit ke Git!
# =============================================================================

set -e

NAMESPACE="buroqet"

echo "🔐 Buroqet Kubernetes Secrets Setup"
echo "====================================="
echo ""

# ─── ISI NILAI INI SEBELUM JALANKAN ──────────────────────────────────────

# MongoDB Atlas (gratis di cloud.mongodb.com)
# Format: mongodb+srv://USERNAME:PASSWORD@CLUSTER.mongodb.net/DATABASE
MONGO_URI_TRACKING="mongodb+srv://zharfanfaz21_db_user:bsUIL2Yyewp5hfUj@buroqet-tracking.qdsmduw.mongodb.net/tracking_db?retryWrites=true&w=majority"

# Upstash Redis (gratis di upstash.com)
REDIS_ADDR="distinct-dinosaur-128801.upstash.io:6379"
REDIS_PASSWORD="gQAAAAAAAfchAAIgcDEwOGNiOWUwNTA4NmU0MDJkOWI3YjAxYzk4NTE3ZjA0MA"

# Upstash Kafka (gratis di upstash.com/kafka)
KAFKA_BROKER="kafka.buroqet.svc.cluster.local:9092"
KAFKA_USERNAME=""
KAFKA_PASSWORD=""

# PostgreSQL untuk Auth Service (Neon.tech gratis: neon.tech)
# Format: postgresql://user:pass@host/dbname
POSTGRES_AUTH_HOST="aws-1-ap-southeast-1.pooler.supabase.com"
POSTGRES_AUTH_USER="postgres.gvmknmaybejslbghzsov"
POSTGRES_AUTH_PASSWORD="HfuyrF1xd1iUVagY"
POSTGRES_AUTH_DB="postgres"

# JWT Secret (minimal 32 karakter random)
# Generate: openssl rand -hex 32
JWT_SECRET="a3f5b72183e9b1c5c4e09d18e2a3c749b5d3c81e9f0c72a1e3b6d9f8c4e2a1b5"

# ─── TIDAK PERLU EDIT DI BAWAH INI ──────────────────────────────────────

echo "📌 Membuat secret untuk Tracking Service..."
kubectl create secret generic tracking-secret \
  --namespace="$NAMESPACE" \
  --from-literal=mongo_uri="$MONGO_URI_TRACKING" \
  --from-literal=redis_password="$REDIS_PASSWORD" \
  --dry-run=client -o yaml | kubectl apply -f -
echo "✅ tracking-secret"

echo ""
echo "📌 Membuat secret untuk Auth Service..."
kubectl create secret generic auth-secret \
  --namespace="$NAMESPACE" \
  --from-literal=db_user="$POSTGRES_AUTH_USER" \
  --from-literal=db_password="$POSTGRES_AUTH_PASSWORD" \
  --from-literal=jwt_secret="$JWT_SECRET" \
  --dry-run=client -o yaml | kubectl apply -f -
echo "✅ auth-secret"

echo ""
echo "📌 Membuat shared secret (Kafka) untuk semua service..."
kubectl create secret generic kafka-secret \
  --namespace="$NAMESPACE" \
  --from-literal=kafka_username="$KAFKA_USERNAME" \
  --from-literal=kafka_password="$KAFKA_PASSWORD" \
  --dry-run=client -o yaml | kubectl apply -f -
echo "✅ kafka-secret"

echo ""
echo "📌 Update ConfigMap tracking dengan Upstash Kafka + Redis..."
kubectl create configmap tracking-config \
  --namespace="$NAMESPACE" \
  --from-literal=mongo_db="tracking_db" \
  --from-literal=redis_addr="$REDIS_ADDR" \
  --from-literal=kafka_broker="$KAFKA_BROKER" \
  --from-literal=app_port="8080" \
  --dry-run=client -o yaml | kubectl apply -f -
echo "✅ tracking-config"

echo ""
echo "📌 Update ConfigMap auth..."
kubectl create configmap auth-config \
  --namespace="$NAMESPACE" \
  --from-literal=db_host="$POSTGRES_AUTH_HOST" \
  --from-literal=db_port="5432" \
  --from-literal=db_name="$POSTGRES_AUTH_DB" \
  --from-literal=app_port="8080" \
  --dry-run=client -o yaml | kubectl apply -f -
echo "✅ auth-config"

echo ""
echo "📌 Update ConfigMap order..."
kubectl create configmap order-config \
  --namespace="$NAMESPACE" \
  --from-literal=kafka_broker="$KAFKA_BROKER" \
  --from-literal=app_port="8080" \
  --from-literal=pricing_service_url="http://pricing-service:8084" \
  --dry-run=client -o yaml | kubectl apply -f -
echo "✅ order-config"

echo ""
echo "📌 Update Secret order..."
POSTGRES_ORDER_URL="postgresql://$POSTGRES_AUTH_USER:$POSTGRES_AUTH_PASSWORD@$POSTGRES_AUTH_HOST:5432/$POSTGRES_AUTH_DB?sslmode=require&search_path=oms"
kubectl create secret generic order-secret \
  --namespace="$NAMESPACE" \
  --from-literal=database_url="$POSTGRES_ORDER_URL" \
  --dry-run=client -o yaml | kubectl apply -f -
echo "✅ order-secret"

echo ""
echo "📌 Update Secret epod..."
POSTGRES_EPOD_URL="postgresql://$POSTGRES_AUTH_USER:$POSTGRES_AUTH_PASSWORD@$POSTGRES_AUTH_HOST:5432/$POSTGRES_AUTH_DB?sslmode=require&search_path=epod"
kubectl create secret generic epod-secret \
  --namespace="$NAMESPACE" \
  --from-literal=database_url="$POSTGRES_EPOD_URL" \
  --dry-run=client -o yaml | kubectl apply -f -
echo "✅ epod-secret"

echo ""
echo "📌 Update ConfigMap epod..."
kubectl create configmap epod-config \
  --namespace="$NAMESPACE" \
  --from-literal=kafka_broker="$KAFKA_BROKER" \
  --from-literal=app_port="8080" \
  --dry-run=client -o yaml | kubectl apply -f -
echo "✅ epod-config"

echo ""
echo "═══════════════════════════════════════"
echo "✅ Semua secrets berhasil dibuat!"
echo "═══════════════════════════════════════"
echo ""
echo "Verifikasi:"
echo "  kubectl get secrets -n $NAMESPACE"
echo "  kubectl get configmaps -n $NAMESPACE"
