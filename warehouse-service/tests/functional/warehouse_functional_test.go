package functional_test

import (
    "context"
    "os"
    "testing"

    "github.com/stretchr/testify/suite"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "warehouse-service/internal/domain"
    "warehouse-service/internal/repository"
    "warehouse-service/internal/service"
)

// WarehouseFunctionalSuite adalah test suite untuk functional test
// Suite pattern memudahkan setup dan teardown database
type WarehouseFunctionalSuite struct {
    suite.Suite
    db      *gorm.DB
    service *service.WarehouseService
    ctx     context.Context
}

// SetupSuite dijalankan SEKALI sebelum semua test dalam suite
func (s *WarehouseFunctionalSuite) SetupSuite() {
    // Baca database URL dari environment variable
    // Saat CI/CD, ini akan diisi oleh Docker Compose
    dbURL := os.Getenv("TEST_DATABASE_URL")
    if dbURL == "" {
        dbURL = "host=localhost user=testuser password=testpass dbname=wms_test port=5432 sslmode=disable"
    }

    db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
    s.Require().NoError(err, "Gagal koneksi ke database test")

    // Auto migrate — manifest dulu agar FK packages.manifest_id valid
    err = db.AutoMigrate(&domain.Manifest{}, &domain.Package{})
    s.Require().NoError(err)

    s.db = db
    s.ctx = context.Background()

    // Gunakan repository NYATA (bukan mock)
    repo := repository.NewWarehouseRepository(db)
    // Untuk Kafka, kita masih bisa mock karena Kafka mungkin tidak ada di test environment
    // Atau gunakan embedded Kafka jika diperlukan
    kafkaMock := &NoOpKafkaProducer{} // Simple no-op implementasi
    s.service = service.NewWarehouseService(repo, kafkaMock)
}

// NoOpKafkaProducer adalah implementasi Kafka yang tidak melakukan apa-apa
// Digunakan di functional test agar tidak perlu Kafka nyata
type NoOpKafkaProducer struct{}

func (k *NoOpKafkaProducer) PublishEvent(ctx context.Context, topic, key string, value []byte) error {
    return nil // Diabaikan
}

// SetupTest dijalankan sebelum SETIAP test — bersihkan data
func (s *WarehouseFunctionalSuite) SetupTest() {
    // Hapus semua data test agar setiap test mulai dari kondisi bersih
    s.db.Exec("DELETE FROM packages")
    s.db.Exec("DELETE FROM manifests")
}

// TearDownSuite dijalankan SEKALI setelah semua test selesai
func (s *WarehouseFunctionalSuite) TearDownSuite() {
    sqlDB, _ := s.db.DB()
    sqlDB.Close()
}

// TestInboundFlow menguji alur lengkap inbound paket dengan database nyata
func (s *WarehouseFunctionalSuite) TestInboundFlow() {
    awb := "FUNC-AWB-001"
    hubID := "HUB-BANDUNG-01"

    // Jalankan proses inbound
    pkg, err := s.service.ProcessInbound(s.ctx, awb, hubID)

    s.NoError(err)
    s.NotNil(pkg)
    s.Equal(awb, pkg.AWB)
    s.Equal("INBOUND", pkg.Status)

    // Verifikasi data benar-benar tersimpan di database
    var savedPkg domain.Package
    result := s.db.Where("awb = ?", awb).First(&savedPkg)
    s.NoError(result.Error)
    s.Equal(awb, savedPkg.AWB)
    s.Equal("INBOUND", savedPkg.Status)
}

// TestDuplicateInbound memastikan duplikasi AWB ditolak.
// LULUS di Jenkins karena ProcessInbound sudah mengecek duplikat (bukan kode belum selesai).
func (s *WarehouseFunctionalSuite) TestDuplicateInbound() {
    awb := "FUNC-AWB-DUPLIKAT"

    _, err := s.service.ProcessInbound(s.ctx, awb, "HUB-01")
    s.NoError(err)

    _, err = s.service.ProcessInbound(s.ctx, awb, "HUB-01")
    s.Error(err)
}

// TestDispatchNonExistentManifestFails menguji dispatch manifest yang tidak ada.
// GAGAL selama service/repo belum mengembalikan error untuk manifest tidak valid.
func (s *WarehouseFunctionalSuite) TestDispatchNonExistentManifestFails() {
    err := s.service.DispatchManifest(s.ctx, "MANIFEST-FUNC-404")
    s.Error(err, "Dispatch manifest yang tidak ada harus gagal")
}

// TestDispatchEmptyManifestFails menguji dispatch manifest tanpa paket harus ditolak.
// GAGAL selama validasi manifest kosong belum ada di service layer.
func (s *WarehouseFunctionalSuite) TestDispatchEmptyManifestFails() {
    manifestID := "MANIFEST-FUNC-EMPTY"
    repo := repository.NewWarehouseRepository(s.db)

    err := repo.CreateManifest(s.ctx, &domain.Manifest{
        ID:          manifestID,
        OriginHubID: "HUB-JKT-01",
        DestHubID:   "HUB-BDG-01",
        Status:      "OPEN",
    })
    s.Require().NoError(err)

    err = s.service.DispatchManifest(s.ctx, manifestID)
    s.Error(err, "Dispatch manifest kosong harus gagal")
}

// Entry point untuk menjalankan suite
func TestWarehouseFunctionalSuite(t *testing.T) {
    suite.Run(t, new(WarehouseFunctionalSuite))
}