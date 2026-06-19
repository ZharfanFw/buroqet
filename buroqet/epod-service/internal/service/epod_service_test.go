package service_test

import (
	"context"
	"strings"
	"testing"

	"epod-service/internal/service"
	"epod-service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessUpload(t *testing.T) {
	// 1. Inisialisasi semua mock yang dibutuhkan (3 argumen)
	// Catatan: Pastikan struct mock ini sudah di-generate oleh Mockery di folder /mocks
	mockRepo := new(mocks.EPODRepository)
	mockStorage := new(mocks.StorageInterface)
	mockKafka := new(mocks.ProducerInterface)

	// 2. Masukkan 3 mock ke constructor service (Garis merah beres! ✅)
	epodSvc := service.NewEPODService(
		mockRepo,
		mockStorage,
		mockKafka,
	)

	// 3. Siapkan data testing dummy
	ctx := context.Background()
	awb := "AWB-001"
	receiver := "Budi Santoso"
	lat := -6.2000
	lng := 106.8000
	filename := "proof.jpg"
	dummyFileContent := strings.NewReader("fake-image-binary-data") // Pura-pura jadi file gambar

	// 4. Set Ekspektasi (Behavior) dari masing-masing Mock
	// Mock Storage: kalau dipanggil Upload, kembalikan path sukses
	mockStorage.On("Upload", mock.Anything, mock.Anything).Return("/uploads/12345_proof.jpg", nil)

	// Mock Repo: kalau dipanggil Create, anggap sukses simpan ke database (return nil error)
	mockRepo.On("Create", mock.Anything).Return(nil)

	// Mock Kafka: kalau dipanggil PublishEvent, anggap sukses kirim message (return nil error)
	mockKafka.On("PublishEvent", mock.Anything, "epod-events", mock.Anything).Return(nil)

	// 5. Eksekusi fungsi ProcessUpload yang baru
	resp, err := epodSvc.ProcessUpload(ctx, awb, receiver, lat, lng, dummyFileContent, filename)

	// 6. Validasi Hasil (Assertion)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "PENDING_REVIEW", resp.Status) // Sesuai default status di service
	assert.Equal(t, "/uploads/12345_proof.jpg", resp.PhotoURL)
	assert.Equal(t, awb, resp.AWB)

	// 7. Pastikan semua fungsi mock di atas benar-benar terpanggil selama test berjalan
	mockStorage.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockKafka.AssertExpectations(t)
}