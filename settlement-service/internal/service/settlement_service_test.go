// settlement-service/internal/service/settlement_service_test.go

package service_test

import (
    "context"
    "errors"
    "testing"

    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"

    "settlement-service/internal/service"
    "settlement-service/mocks"
)

func TestProcessDeliveryCommission(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockSettlementRepository(ctrl)
    mockPricing := mocks.NewMockPricingServiceClient(ctrl)

    svc := service.NewSettlementService(mockRepo, mockPricing)
    ctx := context.Background()

    t.Run("sukses_hitung_komisi_reguler", func(t *testing.T) {
        courierID := "COURIER-001"
        awb := "AWB-DELIVERED-001"
        serviceType := "REGULER"
        commissionRate := 3500.0 // Rp 3.500 per pengiriman

        // Mock pricing service mengembalikan rate
        mockPricing.EXPECT().
            GetCommissionRate(gomock.Any(), serviceType).
            Return(commissionRate, nil).
            Times(1)

        // Mock repo berhasil menyimpan
        mockRepo.EXPECT().
            CreateCommissionLog(gomock.Any(), gomock.Any()).
            DoAndReturn(func(ctx context.Context, log interface{}) error {
                // Kita bisa validasi isi dari argument yang dikirim
                // Ini pattern "DoAndReturn" untuk assertion lebih dalam
                return nil
            }).
            Times(1)

        err := svc.ProcessDeliveryCommission(ctx, courierID, awb, serviceType)

        assert.NoError(t, err)
    })

    t.Run("gagal_courier_id_kosong", func(t *testing.T) {
        // Tidak ada ekspektasi ke repo/pricing
        err := svc.ProcessDeliveryCommission(ctx, "", "AWB-001", "REGULER")
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "courier ID dan AWB tidak boleh kosong")
    })

    t.Run("gagal_pricing_service_error", func(t *testing.T) {
        mockPricing.EXPECT().
            GetCommissionRate(gomock.Any(), "EXPRESS").
            Return(0.0, errors.New("pricing service unavailable")).
            Times(1)

        // Repo tidak boleh dipanggil kalau pricing gagal
        mockRepo.EXPECT().CreateCommissionLog(gomock.Any(), gomock.Any()).Times(0)

        err := svc.ProcessDeliveryCommission(ctx, "COURIER-001", "AWB-001", "EXPRESS")
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "gagal mengambil commission rate")
    })

    t.Run("gagal_rate_tidak_valid", func(t *testing.T) {
        // Rate 0 dianggap tidak valid
        mockPricing.EXPECT().
            GetCommissionRate(gomock.Any(), "GRATIS").
            Return(0.0, nil).
            Times(1)

        mockRepo.EXPECT().CreateCommissionLog(gomock.Any(), gomock.Any()).Times(0)

        err := svc.ProcessDeliveryCommission(ctx, "COURIER-001", "AWB-001", "GRATIS")
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "commission rate tidak valid")
    })
}