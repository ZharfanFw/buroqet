#!/bin/bash
# =============================================================================
# setup-aks.sh — Setup Azure AKS untuk Buroqet
# =============================================================================
# PRASYARAT:
#   1. Azure CLI sudah terinstall: https://docs.microsoft.com/cli/azure/install-azure-cli
#   2. kubectl sudah terinstall
#   3. Sudah login ke Azure: az login
#
# CARA PAKAI:
#   chmod +x setup-aks.sh
#   ./setup-aks.sh
# =============================================================================

set -e  # stop jika ada error

# ─── KONFIGURASI — GANTI SESUAI KEBUTUHAN ──────────────────────────────────
RESOURCE_GROUP="buroqet-rg"
LOCATION="southeastasia"          # region terdekat Indonesia
CLUSTER_NAME="buroqet-aks"
NODE_COUNT=2
NODE_SIZE="Standard_B2s"         # 2 vCPU, 4 GB RAM, ~$33/bln per node
NAMESPACE="buroqet"
# ────────────────────────────────────────────────────────────────────────────

echo "🚀 Buroqet AKS Setup Script"
echo "=============================="
echo "Resource Group : $RESOURCE_GROUP"
echo "Location       : $LOCATION"
echo "Cluster        : $CLUSTER_NAME"
echo "Nodes          : $NODE_COUNT x $NODE_SIZE"
echo ""

# ─── STEP 1: Login & pilih subscription ────────────────────────────────────
echo "📌 Step 1: Verifikasi Azure login..."
az account show --query "{subscription:name, id:id}" -o table
echo ""
echo "⚠️  Pastikan subscription di atas adalah Azure for Students!"
echo "    Jika salah, jalankan: az account set --subscription <SUBSCRIPTION_ID>"
echo ""
read -p "Lanjut? (y/n): " -n 1 -r
echo ""
[[ ! $REPLY =~ ^[Yy]$ ]] && exit 1

# ─── STEP 2: Buat Resource Group ──────────────────────────────────────────
echo "📌 Step 2: Membuat Resource Group '$RESOURCE_GROUP'..."
az group create \
  --name "$RESOURCE_GROUP" \
  --location "$LOCATION" \
  --output table

# ─── STEP 3: Buat AKS Cluster ─────────────────────────────────────────────
echo ""
echo "📌 Step 3: Membuat AKS Cluster (ini ~5-10 menit)..."
az aks create \
  --resource-group "$RESOURCE_GROUP" \
  --name "$CLUSTER_NAME" \
  --node-count "$NODE_COUNT" \
  --node-vm-size "$NODE_SIZE" \
  --generate-ssh-keys \
  --enable-cluster-autoscaler \
  --min-count 1 \
  --max-count 3 \
  --output table

echo "✅ AKS Cluster berhasil dibuat!"

# ─── STEP 4: Ambil credentials untuk kubectl ──────────────────────────────
echo ""
echo "📌 Step 4: Mengambil kubeconfig..."
az aks get-credentials \
  --resource-group "$RESOURCE_GROUP" \
  --name "$CLUSTER_NAME" \
  --overwrite-existing

echo "✅ kubectl sudah terkonfigurasi ke cluster '$CLUSTER_NAME'"

# ─── STEP 5: Verifikasi cluster ──────────────────────────────────────────
echo ""
echo "📌 Step 5: Verifikasi cluster..."
kubectl get nodes
echo ""

# ─── STEP 6: Install Nginx Ingress Controller ─────────────────────────────
echo "📌 Step 6: Install Nginx Ingress Controller..."
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.4/deploy/static/provider/cloud/deploy.yaml

echo "⏳ Menunggu Ingress Controller ready (30 detik)..."
sleep 30
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s

echo "✅ Nginx Ingress Controller siap!"

# ─── STEP 7: Buat Namespace ───────────────────────────────────────────────
echo ""
echo "📌 Step 7: Membuat namespace '$NAMESPACE'..."
kubectl apply -f infra/kubernetes/namespace.yaml
echo "✅ Namespace '$NAMESPACE' berhasil dibuat"

# ─── STEP 8: Tampilkan IP Ingress ─────────────────────────────────────────
echo ""
echo "📌 Step 8: IP publik Ingress Load Balancer..."
echo "⏳ Menunggu IP tersedia (bisa 1-2 menit)..."
kubectl get svc -n ingress-nginx ingress-nginx-controller --watch &
PID=$!
sleep 60
kill $PID 2>/dev/null || true

INGRESS_IP=$(kubectl get svc -n ingress-nginx ingress-nginx-controller -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo ""
echo "═══════════════════════════════════════════════════════"
echo "✅ SETUP SELESAI!"
echo "═══════════════════════════════════════════════════════"
echo ""
echo "📍 IP Publik Ingress: $INGRESS_IP"
echo ""
echo "LANGKAH SELANJUTNYA:"
echo "  1. Jalankan: ./setup-secrets.sh (isi credentials dulu)"
echo "  2. Deploy semua service: kubectl apply -f . -n buroqet -R"
echo "  3. Akses via: http://$INGRESS_IP"
echo ""
echo "Untuk monitoring:"
echo "  kubectl get pods -n buroqet"
echo "  kubectl get svc -n buroqet"
echo "  kubectl logs -n buroqet deployment/tracking-service -f"
echo "═══════════════════════════════════════════════════════"

# ─── Output untuk GitHub Actions secrets ─────────────────────────────────
echo ""
echo "📋 GITHUB ACTIONS SECRETS yang perlu ditambahkan:"
echo "  AKS_RESOURCE_GROUP = $RESOURCE_GROUP"
echo "  AKS_CLUSTER_NAME   = $CLUSTER_NAME"
echo "  AZURE_CREDENTIALS  = (jalankan perintah di bawah)"
echo ""
echo "Untuk generate AZURE_CREDENTIALS:"
SUBSCRIPTION_ID=$(az account show --query id -o tsv)
echo "  az ad sp create-for-rbac \\"
echo "    --name buroqet-github-actions \\"
echo "    --role contributor \\"
echo "    --scopes /subscriptions/$SUBSCRIPTION_ID/resourceGroups/$RESOURCE_GROUP \\"
echo "    --sdk-auth"
echo ""
echo "Copy output JSON tersebut sebagai nilai AZURE_CREDENTIALS di GitHub Secrets."
