package functional

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"dispatch-fleet/internal/domain"
	"dispatch-fleet/internal/repository"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// ── Database Setup ─────────────────────────────────────────────────────────

// testDB menyimpan koneksi database yang dibagikan selama test suite berjalan.
var testDB *sql.DB

// TestMain adalah entry point untuk functional test suite.
// Setup koneksi DB dilakukan sekali, lalu semua test berjalan, lalu cleanup.
func TestMain(m *testing.M) {
	db, err := setupTestDB()
	if err != nil {
		fmt.Printf("FATAL: cannot connect to test database: %v\n", err)
		os.Exit(1)
	}
	testDB = db
	defer testDB.Close()

	if err := migrateTestSchema(testDB); err != nil {
		fmt.Printf("FATAL: cannot migrate test schema: %v\n", err)
		os.Exit(1)
	}

	// Jalankan semua test, lalu exit dengan kode dari test runner
	exitCode := m.Run()
	os.Exit(exitCode)
}

// setupTestDB membuat koneksi ke PostgreSQL test container.
// DSN dibaca dari env variable agar fleksibel di CI/CD.
func setupTestDB() (*sql.DB, error) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		// Default: gunakan container dari docker-compose.test.yml
		dsn = "host=127.0.0.1 port=5433 user=user password=pass dbname=dispatch_db sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening db connection: %w", err)
	}

	// Retry koneksi sampai DB ready (berguna saat container baru start)
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		if err := db.PingContext(context.Background()); err == nil {
			fmt.Println("✅ Test database connected.")
			return db, nil
		}
		fmt.Printf("⏳ Waiting for test DB... attempt %d/%d\n", i+1, maxRetries)
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("test database not reachable after %d retries", maxRetries)
}

// migrateTestSchema membuat tabel yang dibutuhkan di test database.
// Ini menggantikan migration tool agar test bisa berdiri sendiri.
func migrateTestSchema(db *sql.DB) error {
	// language=sql
	schema := `
		CREATE EXTENSION IF NOT EXISTS postgis;

		DROP TABLE IF EXISTS couriers;

		CREATE TABLE couriers (
			id               TEXT PRIMARY KEY,
			name             TEXT NOT NULL,
			current_location GEOGRAPHY(POINT, 4326) NOT NULL,
			status           TEXT NOT NULL CHECK (status IN ('available', 'assigned', 'on_delivery', 'offline')),
			created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		-- Spatial index untuk ST_DWithin agar query efisien
		CREATE INDEX idx_couriers_location ON couriers USING GIST(current_location);
		CREATE INDEX idx_couriers_status   ON couriers (status);
	`

	if _, err := db.ExecContext(context.Background(), schema); err != nil {
		return fmt.Errorf("running schema migration: %w", err)
	}

	fmt.Println("✅ Test schema migrated.")
	return nil
}

// cleanTable truncates tabel sebelum setiap test untuk isolasi data.
func cleanTable(t *testing.T, db *sql.DB, tableName string) {
	t.Helper()
	if _, err := db.ExecContext(context.Background(), "TRUNCATE TABLE "+tableName+" CASCADE"); err != nil {
		t.Fatalf("cleaning table %s: %v", tableName, err)
	}
}

// ── Functional Tests ───────────────────────────────────────────────────────

// TestFleetRepository_FindNearestAvailableCourier_PostGIS memverifikasi bahwa
// query ST_Distance berjalan dengan benar di PostgreSQL + PostGIS nyata.
func TestFleetRepository_FindNearestAvailableCourier_PostGIS(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	cleanTable(t, testDB, "couriers")

	repo := repository.NewPostgresFleetRepository(testDB)
	ctx := context.Background()

	// ── Seed data kurir dummy ─────────────────────────────────────────────
	//
	// Skenario geospasial:
	//   Pickup: Gedung Sarinah, Jakarta (-6.1892, 106.8233)
	//   Kurir A (available):  ~300m dari pickup  → HARUS terpilih
	//   Kurir B (available):  ~1200m dari pickup → di-skip karena lebih jauh
	//   Kurir C (assigned):   ~100m dari pickup  → di-skip karena status bukan available
	//   Kurir D (offline):    ~50m dari pickup   → di-skip karena offline
	//   Kurir E (available):  ~8km dari pickup   → di luar radius 5km

	pickupLocation := domain.Point{Longitude: 106.8233, Latitude: -6.1892}

	couriers := []struct {
		courier      *domain.Courier
		description  string
		expectPicked bool
	}{
		{
			description:  "Kurir A: available, ~300m → HARUS dipilih",
			expectPicked: true,
			courier: &domain.Courier{
				ID: "courier-A", Name: "Andi Prasetyo",
				CurrentLocation: domain.Point{Longitude: 106.8260, Latitude: -6.1905}, // ~300m
				Status:          domain.StatusAvailable, CreatedAt: time.Now(), UpdatedAt: time.Now(),
			},
		},
		{
			description:  "Kurir B: available, ~1200m → skip (lebih jauh dari A)",
			expectPicked: false,
			courier: &domain.Courier{
				ID: "courier-B", Name: "Budi Hartono",
				CurrentLocation: domain.Point{Longitude: 106.8340, Latitude: -6.1950}, // ~1200m
				Status:          domain.StatusAvailable, CreatedAt: time.Now(), UpdatedAt: time.Now(),
			},
		},
		{
			description:  "Kurir C: assigned, ~100m → skip (status bukan available)",
			expectPicked: false,
			courier: &domain.Courier{
				ID: "courier-C", Name: "Citra Lestari",
				CurrentLocation: domain.Point{Longitude: 106.8240, Latitude: -6.1900}, // ~100m
				Status:          domain.StatusAssigned, CreatedAt: time.Now(), UpdatedAt: time.Now(),
			},
		},
		{
			description:  "Kurir D: offline, ~50m → skip (status offline)",
			expectPicked: false,
			courier: &domain.Courier{
				ID: "courier-D", Name: "Dedi Maulana",
				CurrentLocation: domain.Point{Longitude: 106.8236, Latitude: -6.1895}, // ~50m
				Status:          domain.StatusOffline, CreatedAt: time.Now(), UpdatedAt: time.Now(),
			},
		},
		{
			description:  "Kurir E: available, ~8km → skip (luar radius 5km)",
			expectPicked: false,
			courier: &domain.Courier{
				ID: "courier-E", Name: "Eka Saputra",
				CurrentLocation: domain.Point{Longitude: 106.9050, Latitude: -6.2300}, // ~8km
				Status:          domain.StatusAvailable, CreatedAt: time.Now(), UpdatedAt: time.Now(),
			},
		},
	}

	for _, tc := range couriers {
		if err := repo.InsertCourier(ctx, tc.courier); err != nil {
			t.Fatalf("seeding %s: %v", tc.description, err)
		}
	}

	// ── Jalankan fungsi yang ditest ───────────────────────────────────────
	radiusMeters := 5000.0
	foundCourier, distanceMeters, err := repo.FindNearestAvailableCourier(ctx, pickupLocation, radiusMeters)

	// ── Assertions ────────────────────────────────────────────────────────
	if err != nil {
		t.Fatalf("FindNearestAvailableCourier returned unexpected error: %v", err)
	}
	if foundCourier == nil {
		t.Fatal("expected a courier to be found, got nil")
	}
	if foundCourier.ID != "courier-A" {
		t.Errorf("expected courier-A to be selected (nearest available), got: %s", foundCourier.ID)
	}
	if foundCourier.Status != domain.StatusAvailable {
		t.Errorf("expected status 'available', got: %q", foundCourier.Status)
	}
	if distanceMeters <= 0 || distanceMeters > radiusMeters {
		t.Errorf("distance %.2fm is out of expected range (0, %.2f]", distanceMeters, radiusMeters)
	}

	t.Logf("✅ Kurir terpilih: %s (%s), jarak: %.2f meter", foundCourier.ID, foundCourier.Name, distanceMeters)
}

// TestFleetRepository_FindNearestAvailableCourier_NoCourierInRadius memverifikasi
// bahwa nil dikembalikan ketika tidak ada kurir dalam radius.
func TestFleetRepository_FindNearestAvailableCourier_NoCourierInRadius(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	cleanTable(t, testDB, "couriers")

	repo := repository.NewPostgresFleetRepository(testDB)
	ctx := context.Background()

	// Seed satu kurir jauh di luar radius
	courier := &domain.Courier{
		ID: "courier-far", Name: "Farida",
		CurrentLocation: domain.Point{Longitude: 107.6191, Latitude: -6.9175}, // Bandung
		Status:          domain.StatusAvailable, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := repo.InsertCourier(ctx, courier); err != nil {
		t.Fatalf("seeding courier: %v", err)
	}

	// Pickup di Jakarta, radius 5km
	pickupLocation := domain.Point{Longitude: 106.8233, Latitude: -6.1892}
	found, dist, err := repo.FindNearestAvailableCourier(ctx, pickupLocation, 5000.0)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != nil {
		t.Errorf("expected nil courier (none in radius), got: %+v (distance: %.2fm)", found, dist)
	}

	t.Log("✅ Tidak ada kurir dalam radius — fungsi mengembalikan nil seperti yang diharapkan.")
}

// TestFleetRepository_UpdateCourierStatus memverifikasi bahwa status kurir
// ter-update dengan benar di database.
func TestFleetRepository_UpdateCourierStatus(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	cleanTable(t, testDB, "couriers")

	repo := repository.NewPostgresFleetRepository(testDB)
	ctx := context.Background()

	// Seed kurir dengan status available
	courier := &domain.Courier{
		ID: "courier-update-test", Name: "Gilang",
		CurrentLocation: domain.Point{Longitude: 106.8260, Latitude: -6.1905},
		Status:          domain.StatusAvailable, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := repo.InsertCourier(ctx, courier); err != nil {
		t.Fatalf("seeding courier: %v", err)
	}

	// Update status → assigned
	if err := repo.UpdateCourierStatus(ctx, courier.ID, domain.StatusAssigned); err != nil {
		t.Fatalf("UpdateCourierStatus returned error: %v", err)
	}

	// Verifikasi perubahan status di DB
	updated, err := repo.GetCourierByID(ctx, courier.ID)
	if err != nil {
		t.Fatalf("GetCourierByID returned error: %v", err)
	}
	if updated.Status != domain.StatusAssigned {
		t.Errorf("expected status %q after update, got %q", domain.StatusAssigned, updated.Status)
	}

	t.Logf("✅ Status kurir berhasil diupdate: %s → %s", domain.StatusAvailable, updated.Status)
}

// TestFleetRepository_UpdateCourierStatus_CourierNotFound memverifikasi bahwa
// error ErrCourierNotFound dikembalikan jika ID tidak ada.
func TestFleetRepository_UpdateCourierStatus_CourierNotFound(t *testing.T) {
	if testDB == nil {
		t.Skip("test database not available")
	}
	cleanTable(t, testDB, "couriers")

	repo := repository.NewPostgresFleetRepository(testDB)
	ctx := context.Background()

	err := repo.UpdateCourierStatus(ctx, "non-existent-id", domain.StatusAssigned)
	if err == nil {
		t.Fatal("expected error for non-existent courier, got nil")
	}
	if !isErrorMatch(err, domain.ErrCourierNotFound) {
		t.Errorf("expected ErrCourierNotFound, got: %v", err)
	}

	t.Log("✅ ErrCourierNotFound dikembalikan dengan benar untuk ID yang tidak ada.")
}

// isErrorMatch helper untuk cek wrapped errors.
func isErrorMatch(err, target error) bool {
	if err == target {
		return true
	}
	type unwrapper interface{ Unwrap() error }
	if u, ok := err.(unwrapper); ok {
		return isErrorMatch(u.Unwrap(), target)
	}
	return false
}
