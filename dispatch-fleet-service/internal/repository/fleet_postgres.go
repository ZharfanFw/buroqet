package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"dispatch-fleet/internal/domain"
)

// postgresFleetRepository adalah implementasi domain.FleetRepository menggunakan PostgreSQL + PostGIS.
type postgresFleetRepository struct {
	db *sql.DB
}

// NewPostgresFleetRepository membuat instance baru postgresFleetRepository.
func NewPostgresFleetRepository(db *sql.DB) domain.FleetRepository {
	return &postgresFleetRepository{db: db}
}

// FindNearestAvailableCourier mencari kurir 'available' terdekat menggunakan ST_Distance PostGIS.
//
// Query menggunakan:
//   - ST_Distance: menghitung jarak dalam meter (dengan geography cast)
//   - ST_DWithin: filter awal efisien menggunakan spatial index
//   - ST_MakePoint: membangun geometry dari lon/lat
//
// TODO: Implementasi query SQL lengkap.
func (r *postgresFleetRepository) FindNearestAvailableCourier(
	ctx context.Context,
	pickup domain.Point,
	radiusMeters float64,
) (*domain.Courier, float64, error) {
	// language=sql
	query := `
		SELECT
			id,
			name,
			ST_X(current_location::geometry) AS longitude,
			ST_Y(current_location::geometry) AS latitude,
			status,
			created_at,
			updated_at,
			ST_Distance(
				current_location,
				ST_MakePoint($1, $2)::geography
			) AS distance_meters
		FROM couriers
		WHERE
			status = 'available'
			AND ST_DWithin(
				current_location,
				ST_MakePoint($1, $2)::geography,
				$3
			)
		ORDER BY distance_meters ASC
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, pickup.Longitude, pickup.Latitude, radiusMeters)

	var (
		c               domain.Courier
		distanceMeters  float64
		statusStr       string
	)

	err := row.Scan(
		&c.ID,
		&c.Name,
		&c.CurrentLocation.Longitude,
		&c.CurrentLocation.Latitude,
		&statusStr,
		&c.CreatedAt,
		&c.UpdatedAt,
		&distanceMeters,
	)
	if err == sql.ErrNoRows {
		return nil, 0, nil // tidak ada kurir ditemukan
	}
	if err != nil {
		return nil, 0, fmt.Errorf("querying nearest courier: %w", err)
	}

	c.Status = domain.CourierStatus(statusStr)
	return &c, distanceMeters, nil
}

// UpdateCourierStatus mengupdate status kurir berdasarkan ID.
func (r *postgresFleetRepository) UpdateCourierStatus(
	ctx context.Context,
	courierID string,
	status domain.CourierStatus,
) error {
	// language=sql
	query := `
		UPDATE couriers
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, string(status), time.Now(), courierID)
	if err != nil {
		return fmt.Errorf("updating courier status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrCourierNotFound
	}

	return nil
}

// InsertCourier menyimpan data kurir baru ke database.
// Digunakan terutama untuk seeding data pada functional test.
func (r *postgresFleetRepository) InsertCourier(
	ctx context.Context,
	courier *domain.Courier,
) error {
	// language=sql
	query := `
		INSERT INTO couriers (id, name, current_location, status, created_at, updated_at)
		VALUES (
			$1,
			$2,
			ST_MakePoint($3, $4)::geography,
			$5,
			$6,
			$7
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		courier.ID,
		courier.Name,
		courier.CurrentLocation.Longitude,
		courier.CurrentLocation.Latitude,
		string(courier.Status),
		courier.CreatedAt,
		courier.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting courier: %w", err)
	}

	return nil
}

// GetCourierByID mengambil data kurir berdasarkan ID.
func (r *postgresFleetRepository) GetCourierByID(
	ctx context.Context,
	courierID string,
) (*domain.Courier, error) {
	// language=sql
	query := `
		SELECT
			id,
			name,
			ST_X(current_location::geometry) AS longitude,
			ST_Y(current_location::geometry) AS latitude,
			status,
			created_at,
			updated_at
		FROM couriers
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, courierID)

	var (
		c         domain.Courier
		statusStr string
	)

	err := row.Scan(
		&c.ID,
		&c.Name,
		&c.CurrentLocation.Longitude,
		&c.CurrentLocation.Latitude,
		&statusStr,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrCourierNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("querying courier by id: %w", err)
	}

	c.Status = domain.CourierStatus(statusStr)
	return &c, nil
}