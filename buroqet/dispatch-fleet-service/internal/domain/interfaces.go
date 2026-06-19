package domain

import "context"

// FleetRepository mendefinisikan kontrak akses data untuk entitas kurir.
// Implementasi konkret ada di internal/repository/.
//
//go:generate go run go.uber.org/mock/mockgen -source=interfaces.go -destination=../../mocks/mock_fleet_repository.go -package=mocks
type FleetRepository interface {
	// FindNearestAvailableCourier mencari kurir berstatus 'available' terdekat
	// dari titik pickup dalam radius (meter) tertentu menggunakan PostGIS ST_Distance.
	FindNearestAvailableCourier(ctx context.Context, pickup Point, radiusMeters float64) (*Courier, float64, error)

	// UpdateCourierStatus mengubah status kurir berdasarkan ID.
	UpdateCourierStatus(ctx context.Context, courierID string, status CourierStatus) error

	// InsertCourier menyimpan data kurir baru (digunakan untuk seeding test data).
	InsertCourier(ctx context.Context, courier *Courier) error

	// GetCourierByID mengambil data kurir berdasarkan ID.
	GetCourierByID(ctx context.Context, courierID string) (*Courier, error)
}

// DispatchService mendefinisikan kontrak bisnis untuk proses dispatching kurir.
// Implementasi konkret ada di internal/service/.
type DispatchService interface {
	// AssignCourierToPickup mencari kurir terdekat yang tersedia dan mengassign-nya.
	// Mengembalikan error jika tidak ada kurir tersedia dalam radius yang diberikan.
	AssignCourierToPickup(ctx context.Context, pickupLocation Point, radiusMeters float64) (*AssignCourierResult, error)
}