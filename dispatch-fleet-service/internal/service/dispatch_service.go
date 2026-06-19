package service

import (
	"context"
	"fmt"

	"dispatch-fleet/internal/domain"
)

// dispatchService adalah implementasi konkret dari domain.DispatchService.
type dispatchService struct {
	fleetRepo domain.FleetRepository
}

// NewDispatchService membuat instance baru dispatchService.
func NewDispatchService(fleetRepo domain.FleetRepository) domain.DispatchService {
	return &dispatchService{
		fleetRepo: fleetRepo,
	}
}

// AssignCourierToPickup mencari kurir terdekat yang tersedia dan mengassign-nya.
//
// Alur bisnis:
//  1. Cari kurir berstatus 'available' terdekat dalam radius dari pickup location.
//  2. Jika tidak ada → kembalikan ErrNoCourierAvailable.
//  3. Update status kurir tersebut menjadi 'assigned'.
//  4. Kembalikan data kurir beserta jarak tempuh.
//
// TODO: Implementasi lengkap belum selesai — test akan FAIL sampai logika diisi.
func (s *dispatchService) AssignCourierToPickup(
	ctx context.Context,
	pickupLocation domain.Point,
	radiusMeters float64,
) (*domain.AssignCourierResult, error) {
	// --- LANGKAH 1: Validasi input ---
	if radiusMeters <= 0 {
		return nil, fmt.Errorf("radius must be greater than 0, got %.2f", radiusMeters)
	}

	// --- LANGKAH 2: Cari kurir terdekat ---
	// TODO: Implementasi pencarian menggunakan PostGIS ST_Distance
	courier, distanceMeters, err := s.fleetRepo.FindNearestAvailableCourier(ctx, pickupLocation, radiusMeters)
	if err != nil {
		return nil, fmt.Errorf("finding nearest courier: %w", err)
	}
	if courier == nil {
		return nil, domain.ErrNoCourierAvailable
	}

	// --- LANGKAH 3: Update status kurir → assigned ---
	// TODO: Tambahkan optimistic locking / distributed lock di sini
	if err := s.fleetRepo.UpdateCourierStatus(ctx, courier.ID, domain.StatusAssigned); err != nil {
		return nil, fmt.Errorf("updating courier status: %w", err)
	}

	// --- LANGKAH 4: Kembalikan hasil ---
	courier.Status = domain.StatusAssigned
	return &domain.AssignCourierResult{
		Courier:        courier,
		DistanceMeters: distanceMeters,
	}, nil
}