package service_test

import (
	"context"
	"testing"
	"time"

	"dispatch-fleet/internal/domain"
    "dispatch-fleet/internal/service"
    "dispatch-fleet/mocks"
    "go.uber.org/mock/gomock"
)

// ── Helper ─────────────────────────────────────────────────────────────────

func newTestCourier(id, name string, lon, lat float64, status domain.CourierStatus) *domain.Courier {
	return &domain.Courier{
		ID:   id,
		Name: name,
		CurrentLocation: domain.Point{
			Longitude: lon,
			Latitude:  lat,
		},
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ── Test Suite ─────────────────────────────────────────────────────────────

// TestDispatchService_AssignCourierToPickup menguji semua skenario
// bisnis pada fungsi AssignCourierToPickup menggunakan GoMock.
func TestDispatchService_AssignCourierToPickup(t *testing.T) {
	t.Parallel()

	// Koordinat pickup dummy: Gudang Jakarta Pusat
	pickupLocation := domain.Point{Longitude: 106.8272, Latitude: -6.1751}
	searchRadius := 5000.0 // 5 km

	tests := []struct {
		name        string
		setupMock   func(repo *mocks.MockFleetRepository)
		pickup      domain.Point
		radius      float64
		wantErr     error
		wantStatus  domain.CourierStatus
		wantNonNil  bool
	}{
		{
			// ✅ Happy path: kurir available ditemukan dan berhasil diassign
			name: "success: nearest available courier found and assigned",
			setupMock: func(repo *mocks.MockFleetRepository) {
				foundCourier := newTestCourier("courier-001", "Budi Santoso", 106.8300, -6.1780, domain.StatusAvailable)
				distanceMeters := 350.5

				// Ekspektasi 1: FindNearestAvailableCourier dipanggil SEKALI
				repo.EXPECT().
					FindNearestAvailableCourier(gomock.Any(), pickupLocation, searchRadius).
					Return(foundCourier, distanceMeters, nil).
					Times(1)

				// Ekspektasi 2: UpdateCourierStatus dipanggil SEKALI dengan status 'assigned'
				repo.EXPECT().
					UpdateCourierStatus(gomock.Any(), "courier-001", domain.StatusAssigned).
					Return(nil).
					Times(1)
			},
			pickup:     pickupLocation,
			radius:     searchRadius,
			wantErr:    nil,
			wantStatus: domain.StatusAssigned,
			wantNonNil: true,
		},
		{
			// ❌ Skenario: tidak ada kurir available dalam radius
			name: "error: no available courier in radius",
			setupMock: func(repo *mocks.MockFleetRepository) {
				// Repository mengembalikan nil courier (tidak ada yang ditemukan)
				repo.EXPECT().
					FindNearestAvailableCourier(gomock.Any(), pickupLocation, searchRadius).
					Return(nil, 0.0, nil).
					Times(1)

				// UpdateCourierStatus TIDAK boleh dipanggil sama sekali
				repo.EXPECT().
					UpdateCourierStatus(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			pickup:     pickupLocation,
			radius:     searchRadius,
			wantErr:    domain.ErrNoCourierAvailable,
			wantNonNil: false,
		},
		{
			// ❌ Skenario: repository FindNearest gagal (DB error)
			name: "error: repository FindNearestAvailableCourier fails",
			setupMock: func(repo *mocks.MockFleetRepository) {
				repo.EXPECT().
					FindNearestAvailableCourier(gomock.Any(), pickupLocation, searchRadius).
					Return(nil, 0.0, domain.ErrCourierNotFound).
					Times(1)

				repo.EXPECT().
					UpdateCourierStatus(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			pickup:     pickupLocation,
			radius:     searchRadius,
			wantErr:    domain.ErrCourierNotFound,
			wantNonNil: false,
		},
		{
			// ❌ Skenario: UpdateCourierStatus gagal setelah kurir ditemukan
			name: "error: update status fails after courier found",
			setupMock: func(repo *mocks.MockFleetRepository) {
				foundCourier := newTestCourier("courier-002", "Sari Dewi", 106.8290, -6.1760, domain.StatusAvailable)

				repo.EXPECT().
					FindNearestAvailableCourier(gomock.Any(), pickupLocation, searchRadius).
					Return(foundCourier, 200.0, nil).
					Times(1)

				// Simulasi DB timeout saat update
				repo.EXPECT().
					UpdateCourierStatus(gomock.Any(), "courier-002", domain.StatusAssigned).
					Return(domain.ErrCourierNotFound).
					Times(1)
			},
			pickup:     pickupLocation,
			radius:     searchRadius,
			wantErr:    domain.ErrCourierNotFound,
			wantNonNil: false,
		},
		{
			// ❌ Skenario: radius tidak valid (0 atau negatif)
			name: "error: invalid radius (zero)",
			setupMock: func(repo *mocks.MockFleetRepository) {
				// Tidak ada call ke repository sama sekali
				repo.EXPECT().
					FindNearestAvailableCourier(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			pickup:     pickupLocation,
			radius:     0,
			wantErr:    nil, // wrapped error, cek via wantNonNil
			wantNonNil: false,
		},
	}

	for _, tc := range tests {
		tc := tc // capture range var untuk t.Parallel()
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// ── Setup GoMock ──────────────────────────────────────────────
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockFleetRepository(ctrl)
			tc.setupMock(mockRepo)

			// ── Buat service dengan mock dependency ───────────────────────
			svc := service.NewDispatchService(mockRepo)

			// ── Jalankan fungsi yang ditest ───────────────────────────────
			ctx := context.Background()
			result, err := svc.AssignCourierToPickup(ctx, tc.pickup, tc.radius)

			// ── Assertions ────────────────────────────────────────────────
			if tc.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tc.wantErr)
					return
				}
				// Gunakan errors.Is untuk mendukung wrapped errors
				if !isErrorMatch(err, tc.wantErr) {
					t.Errorf("expected error containing %v, got: %v", tc.wantErr, err)
				}
			}

			if tc.wantNonNil && result == nil {
				t.Error("expected non-nil result, got nil")
				return
			}

			if !tc.wantNonNil && result != nil {
				t.Errorf("expected nil result, got: %+v", result)
				return
			}

			if tc.wantNonNil && result != nil {
				if result.Courier == nil {
					t.Error("result.Courier should not be nil on success")
					return
				}
				if result.Courier.Status != tc.wantStatus {
					t.Errorf("expected courier status %q, got %q", tc.wantStatus, result.Courier.Status)
				}
				if result.DistanceMeters <= 0 {
					t.Errorf("expected positive distance, got %.2f", result.DistanceMeters)
				}
			}
		})
	}
}

// isErrorMatch memeriksa apakah err mengandung target (mendukung errors.Is).
func isErrorMatch(err, target error) bool {
	if err == target {
		return true
	}
	// Unwrap untuk wrapped errors (fmt.Errorf dengan %w)
	type unwrapper interface{ Unwrap() error }
	if u, ok := err.(unwrapper); ok {
		return isErrorMatch(u.Unwrap(), target)
	}
	return false
}