package service

import (
	"fmt"
	"epod-service/internal/domain"
)

type Storage interface {
	Upload(fileName string) (string, error)
}

type KafkaProducer interface {
	Publish(topic string, message string) error
}

type EPODService struct {
	storage Storage
	kafka   KafkaProducer
}

func NewEPODService(
	storage Storage,
	kafka KafkaProducer,
) *EPODService {
	return &EPODService{
		storage: storage,
		kafka:   kafka,
	}
}

func (s *EPODService) ProcessUpload(
	req domain.UploadRequest,
) (domain.UploadResponse, error) {

	url, err := s.storage.Upload(req.FileName)
	if err != nil {
		return domain.UploadResponse{}, fmt.Errorf("gagal mengunggah foto bukti delivery: %w", err)
	}

	err = s.kafka.Publish("package.delivered", req.AWB)
	if err != nil {
		return domain.UploadResponse{}, fmt.Errorf("gagal mengirimkan event ke kafka: %w", err)
	}

	return domain.UploadResponse{
		Status:   "SUCCESS",
		ImageURL: url,
	}, nil
}