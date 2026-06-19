package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"epod-service/internal/kafka"
	"epod-service/internal/repository"
	"epod-service/internal/storage"

	"github.com/google/uuid"
)

type EPODService struct {
	repo          repository.EPODRepository
	fileStorage   storage.StorageInterface
	kafkaProducer kafka.ProducerInterface
}

// NewEPODService sekarang menerima 3 argumen sesuai dengan yang dikirim dari main.go
func NewEPODService(
	repo repository.EPODRepository,
	fileStorage storage.StorageInterface,
	kafkaProducer kafka.ProducerInterface,
) *EPODService {
	return &EPODService{
		repo:          repo,
		fileStorage:   fileStorage,
		kafkaProducer: kafkaProducer,
	}
}

// ProcessUpload menangani logika bisnis pengunggahan foto ePOD
func (s *EPODService) ProcessUpload(
	ctx context.Context, 
	awb, receiver string, 
	lat, lng float64, 
	file io.Reader, 
	filename string,
) (*repository.EPOD, error) {
	
	// 1. Amankan nama file agar unik (mencegah file tertimpa jika namanya sama)
	uniqueFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)

	// 2. Upload file ke storage (Lokal / ./uploads)
	photoURL, err := s.fileStorage.Upload(file, uniqueFilename)
	if err != nil {
		return nil, fmt.Errorf("gagal mengunggah foto: %v", err)
	}

	// 3. Buat objek entity EPOD untuk database
	epodData := &repository.EPOD{
		ID:           uuid.New().String(), // Generate UUID otomatis untuk Supabase
		AWB:          awb,
		Receiver:     receiver,
		PhotoURL:     photoURL,
		SignatureURL: "", // Opsional, bisa dikosongkan dulu
		Lat:          lat,
		Lng:          lng,
		Status:       "PENDING_REVIEW",
		CreatedAt:    time.Now(),
	}

	// 4. Simpan data ke database Supabase via Repository
	err = s.repo.Create(epodData)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan ke database: %v", err)
	}

	// 5. Kirim event ke Kafka bahwa ada ePOD baru masuk
	event := kafka.EPODEvent{
		EventName: "EPOD_UPLOADED",
		EPODID:    epodData.ID,
		AWB:       epodData.AWB,
		Status:    epodData.Status,
		Timestamp: time.Now(),
	}
	// Menggunakan background context agar proses HTTP tidak tersendat jika Kafka lambat
	_ = s.kafkaProducer.PublishEvent(context.Background(), "epod-events", event)

	return epodData, nil
}

// GetAll mengambil semua data ePOD dengan limit default 100 data
func (s *EPODService) GetAll(ctx context.Context, status string) ([]repository.EPOD, error) {
	return s.repo.FindAll(status, 100)
}

// GetByID mengambil data ePOD berdasarkan ID internal (UUID)
func (s *EPODService) GetByID(ctx context.Context, id string) (*repository.EPOD, error) {
	return s.repo.FindByID(id)
}

// GetByAWB mengambil data ePOD berdasarkan nomor resi
func (s *EPODService) GetByAWB(ctx context.Context, awb string) (*repository.EPOD, error) {
	return s.repo.FindByAWB(awb)
}

// VerifyEPOD mengubah status ePOD dan menyebarkan hasilnya ke Kafka
func (s *EPODService) VerifyEPOD(ctx context.Context, id string, status string) error {
	// 1. Cari data ePOD terlebih dahulu untuk mendapatkan nomor AWB-nya
	epod, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// 2. Update status di database (VERIFIED / REJECTED)
	err = s.repo.UpdateStatus(id, status)
	if err != nil {
		return fmt.Errorf("gagal memperbarui status: %v", err)
	}

	// 3. Kirim event hasil verifikasi ke Kafka agar sistem lain tahu
	event := kafka.EPODEvent{
		EventName: "EPOD_STATUS_UPDATED",
		EPODID:    id,
		AWB:       epod.AWB,
		Status:    status,
		Timestamp: time.Now(),
	}
	_ = s.kafkaProducer.PublishEvent(context.Background(), "epod-events", event)

	return nil
}