package functional_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"settlement-service/internal/domain"
	"settlement-service/internal/repository"
	"settlement-service/internal/service"
)

// fixedRatePricing menggantikan Pricing Service HTTP agar functional test
// hanya bergantung pada database (sesuai spesifikasi: DB nyata, WS di-mock).
type fixedRatePricing struct {
	rates map[string]float64
}

func (p *fixedRatePricing) GetCommissionRate(_ context.Context, serviceType string) (float64, error) {
	if rate, ok := p.rates[serviceType]; ok {
		return rate, nil
	}
	return 0, nil
}

type SettlementFunctionalSuite struct {
	suite.Suite
	db      *gorm.DB
	service *service.SettlementService
	ctx     context.Context
}

func (s *SettlementFunctionalSuite) SetupSuite() {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost user=testuser password=testpass dbname=settlement_test port=5434 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	s.Require().NoError(err, "Gagal koneksi ke database test")

	err = db.AutoMigrate(&domain.CommissionLog{})
	s.Require().NoError(err)

	s.db = db
	s.ctx = context.Background()

	repo := repository.NewSettlementRepository(db)
	pricing := &fixedRatePricing{
		rates: map[string]float64{
			"REGULER":  3500,
			"EXPRESS":  5000,
			"SAME_DAY": 7500,
		},
	}
	s.service = service.NewSettlementService(repo, pricing)
}

func (s *SettlementFunctionalSuite) SetupTest() {
	s.db.Exec("DELETE FROM commission_logs")
}

func (s *SettlementFunctionalSuite) TearDownSuite() {
	sqlDB, _ := s.db.DB()
	sqlDB.Close()
}

// TestCommissionPersistedToDB menguji alur komisi tersimpan ke PostgreSQL nyata.
func (s *SettlementFunctionalSuite) TestCommissionPersistedToDB() {
	courierID := "COURIER-FUNC-001"
	awb := "AWB-FUNC-001"

	err := s.service.ProcessDeliveryCommission(s.ctx, courierID, awb, "REGULER")
	s.NoError(err)

	var log domain.CommissionLog
	result := s.db.Where("awb = ?", awb).First(&log)
	s.NoError(result.Error)
	s.Equal(courierID, log.CourierID)
	s.Equal(3500.0, log.Amount)
	s.Equal("PENDING", log.Status)
}

// TestCourierEarningsAggregation menguji agregasi saldo dari beberapa pengiriman.
func (s *SettlementFunctionalSuite) TestCourierEarningsAggregation() {
	courierID := "COURIER-FUNC-002"

	s.Require().NoError(s.service.ProcessDeliveryCommission(s.ctx, courierID, "AWB-FUNC-002A", "REGULER"))
	s.Require().NoError(s.service.ProcessDeliveryCommission(s.ctx, courierID, "AWB-FUNC-002B", "EXPRESS"))

	summary, err := s.service.GetCourierEarnings(s.ctx, courierID)
	s.NoError(err)
	s.Equal(2, summary.TotalDeliveries)
	s.Equal(8500.0, summary.TotalAmount)
	s.Equal(8500.0, summary.PendingAmount)
	s.Equal(0.0, summary.PaidAmount)
}

// TestMarkAsPaidClearsPendingBalance menguji penyelesaian komisi PENDING → PAID di database.
func (s *SettlementFunctionalSuite) TestMarkAsPaidClearsPendingBalance() {
	courierID := "COURIER-FUNC-003"

	s.Require().NoError(s.service.ProcessDeliveryCommission(s.ctx, courierID, "AWB-FUNC-003", "REGULER"))

	repo := repository.NewSettlementRepository(s.db)
	err := repo.MarkAsPaid(s.ctx, courierID)
	s.NoError(err)

	summary, err := s.service.GetCourierEarnings(s.ctx, courierID)
	s.NoError(err)
	s.Equal(0.0, summary.PendingAmount)
	s.Equal(3500.0, summary.PaidAmount)
}

// TestDuplicateAWBCommissionRejected memastikan AWB duplikat ditolak di level bisnis+DB.
// Gagal selama validasi duplikat belum ada di service layer (kode belum selesai).
func (s *SettlementFunctionalSuite) TestDuplicateAWBCommissionRejected() {
	courierID := "COURIER-FUNC-004"
	awb := "AWB-FUNC-DUPLIKAT"

	err := s.service.ProcessDeliveryCommission(s.ctx, courierID, awb, "REGULER")
	s.NoError(err)

	err = s.service.ProcessDeliveryCommission(s.ctx, courierID, awb, "REGULER")
	s.Error(err)
}

func TestSettlementFunctionalSuite(t *testing.T) {
	suite.Run(t, new(SettlementFunctionalSuite))
}