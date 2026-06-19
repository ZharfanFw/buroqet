# Buroqet - Deployment & Fixes To-Do List

Berikut adalah daftar tugas dan error yang masih perlu diselesaikan sebelum sistem siap rilis sepenuhnya ke cloud:

## 🐛 Code Bugs (Golang)
- [ ] **Settlement Service Domain Bug:** URL untuk call ke Pricing Service masih hardcoded ke `http://localhost:8084` di dalam file `settlement.go` (atau file service terkait). Ini harus diubah agar mengambil dari environment variable (misal: `PRICING_SERVICE_URL`) yang sudah disiapkan di `deployment.yaml` (`http://pricing-service:8084`).
- [ ] **Warehouse Service Kafka Placeholder (Bug #9):** Warehouse service saat ini masih mengirim data dummy/placeholder ke Kafka. Perlu diubah agar mempublikasikan event WMS yang sebenarnya.
- [ ] **Auth Service `/refresh` Endpoint (Bug #4):** Memperbaiki implementasi endpoint `/refresh` pada auth service yang bermasalah.

## 🚀 K8s Manifests & Cloud Setup
- [ ] **Ganti Placeholder Image:** Di semua file `deployment.yaml` masing-masing service, cari tulisan `ghcr.io/<your-github-username>/...` dan ganti `<your-github-username>` dengan username GitHub organization atau akun kalian yang sebenarnya.
- [ ] **Persiapan Cloud Accounts:** Buat akun dan setup _Free Tier_ di:
    - MongoDB Atlas (Untuk Tracking)
    - Neon.tech (PostgreSQL untuk Auth, Order, Pricing, dll)
    - Upstash (Redis & Kafka)
- [ ] **Isi Kubernetes Secrets:** Copy `infra/scripts/setup-secrets.sh` ke `setup-secrets.local.sh`, lalu isi variabel yang berawalan `GANTI...` dengan credential asli dari cloud provider di atas. Setelah itu jalankan script-nya.
- [ ] **Setup GitHub Actions Secrets:** Tambahkan `AZURE_CREDENTIALS`, `AKS_RESOURCE_GROUP`, dan `AKS_CLUSTER_NAME` di repository settings untuk mengaktifkan CI/CD pipeline otomatis ke AKS.
- [ ] **Warehouse Service Config:** Pastikan Warehouse service mendengarkan database connection string dari environment variables dengan benar, mengingat sebelumnya datanya di-hardcode.

## 🎨 UI/Frontend
- [ ] **Integrasi API Frontend:** Pastikan URL base API di frontend (`src/store/auth.store.ts` dll) sudah mengarah ke endpoint Ingress Kubernetes, bukan ke `localhost`.
- [ ] **Penyelesaian UI Maps & Tracking:** Memastikan maps UI di frontend menampilkan koordinat baru yang sudah dikirim oleh Backend via Kafka.
