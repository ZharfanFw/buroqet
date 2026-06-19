package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"warehouse-service/internal/domain"
	"warehouse-service/internal/service"
	mock_domain "warehouse-service/mocks"
)

// TestProcessInbound menguji semua skenario fungsi ProcessInbound
func TestProcessInbound(t *testing.T) {

	// --- SETUP MOCK ---
	// gomock.Controller adalah "hakim" yang memastikan semua ekspektasi terpenuhi
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // Cek di akhir test: apakah semua yang diharapkan benar-benar dipanggil?

	// Buat tiruan (mock) dari repository dan kafka
	mockRepo := mock_domain.NewMockWarehouseRepository(ctrl)
	mockKafka := mock_domain.NewMockKafkaProducer(ctrl)

	// Inject mock ke dalam service
	svc := service.NewWarehouseService(mockRepo, mockKafka)

	ctx := context.Background()

	// =========================================================
	// TEST CASE 1: Skenario SUKSES — paket baru masuk gudang
	// =========================================================
	t.Run("sukses_inbound_paket_baru", func(t *testing.T) {
		awb := "AWB-TEST-001"
		hubID := "HUB-JKT-01"

		// Setup ekspektasi:
		// "Saya EXPECT bahwa GetPackageByAWB akan dipanggil SEKALI
		//  dengan argumen apapun untuk ctx dan awb=AWB-TEST-001,
		//  dan akan mengembalikan nil, error (artinya paket belum ada)"
		mockRepo.EXPECT().
			GetPackageByAWB(gomock.Any(), awb).
			Return(nil, errors.New("not found")).
			Times(1)

		// "Saya EXPECT bahwa SavePackage akan dipanggil SEKALI
		//  dengan package apapun, dan akan berhasil (return nil error)"
		mockRepo.EXPECT().
			SavePackage(gomock.Any(), gomock.Any()).
			Return(nil).
			Times(1)

		// "Saya EXPECT bahwa PublishEvent akan dipanggil SEKALI untuk topic package.arrived"
		mockKafka.EXPECT().
			PublishEvent(gomock.Any(), "package.arrived", awb, gomock.Any()).
			Return(nil).
			Times(1)

		// Jalankan fungsi yang ditest
		result, err := svc.ProcessInbound(ctx, awb, hubID)

		// Assertion: apa yang kita harapkan dari hasilnya?
		assert.NoError(t, err)                    // tidak boleh ada error
		assert.NotNil(t, result)                  // hasil tidak boleh nil
		assert.Equal(t, awb, result.AWB)          // AWB harus sesuai
		assert.Equal(t, "INBOUND", result.Status) // status harus INBOUND
		assert.Equal(t, hubID, result.HubID)      // hub ID harus sesuai
	})

	// =========================================================
	// TEST CASE 2: AWB sudah terdaftar sebelumnya
	// =========================================================
	t.Run("gagal_awb_sudah_ada", func(t *testing.T) {
		awb := "AWB-DUPLIKAT-001"
		hubID := "HUB-JKT-01"

		// Mock mengembalikan paket yang sudah ada
		existingPkg := &domain.Package{AWB: awb, Status: "INBOUND"}
		mockRepo.EXPECT().
			GetPackageByAWB(gomock.Any(), awb).
			Return(existingPkg, nil). // paket ditemukan!
			Times(1)

		// SavePackage dan PublishEvent TIDAK boleh dipanggil (Times(0))
		// Kalau dipanggil, test otomatis gagal — ini adalah kekuatan mock!
		mockRepo.EXPECT().SavePackage(gomock.Any(), gomock.Any()).Times(0)
		mockKafka.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		result, err := svc.ProcessInbound(ctx, awb, hubID)

		assert.Error(t, err)                               // harus ada error
		assert.Nil(t, result)                              // result harus nil
		assert.Contains(t, err.Error(), "sudah terdaftar") // pesan error harus mengandung ini
	})

	// =========================================================
	// TEST CASE 3: AWB kosong — validasi input
	// =========================================================
	t.Run("gagal_awb_kosong", func(t *testing.T) {
		// Tidak ada ekspektasi apapun ke repo/kafka
		// Karena validasi harusnya gagal sebelum sampai ke sana

		result, err := svc.ProcessInbound(ctx, "", "HUB-JKT-01")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "AWB tidak boleh kosong")
	})

	// =========================================================
	// TEST CASE 4: Database error saat SavePackage
	// =========================================================
	t.Run("gagal_database_error", func(t *testing.T) {
		awb := "AWB-DB-ERROR"
		hubID := "HUB-JKT-01"

		mockRepo.EXPECT().
			GetPackageByAWB(gomock.Any(), awb).
			Return(nil, errors.New("not found")).
			Times(1)

		// Simulasikan database sedang down
		mockRepo.EXPECT().
			SavePackage(gomock.Any(), gomock.Any()).
			Return(errors.New("connection refused")).
			Times(1)

		// Kafka tidak boleh dipanggil kalau DB gagal
		mockKafka.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		result, err := svc.ProcessInbound(ctx, awb, hubID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "gagal menyimpan paket")
	})

	// =========================================================
	// TEST CASE 5: Kafka gagal — tapi operasi tetap sukses
	// (eventual consistency pattern)
	// =========================================================
	t.Run("sukses_meski_kafka_gagal", func(t *testing.T) {
		awb := "AWB-KAFKA-FAIL"
		hubID := "HUB-JKT-01"

		mockRepo.EXPECT().
			GetPackageByAWB(gomock.Any(), awb).
			Return(nil, errors.New("not found")).
			Times(1)

		mockRepo.EXPECT().
			SavePackage(gomock.Any(), gomock.Any()).
			Return(nil).
			Times(1)

		// Kafka gagal
		mockKafka.EXPECT().
			PublishEvent(gomock.Any(), "package.arrived", awb, gomock.Any()).
			Return(errors.New("kafka broker unavailable")).
			Times(1)

		// Meski kafka gagal, fungsi harus tetap return sukses
		result, err := svc.ProcessInbound(ctx, awb, hubID)

		assert.NoError(t, err) // tetap sukses!
		assert.NotNil(t, result)
		assert.Equal(t, "INBOUND", result.Status)
	})
}
