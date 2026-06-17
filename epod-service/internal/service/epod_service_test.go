package service_test

import (
	"testing"

	"epod-service/internal/domain"
	"epod-service/mocks"
	"epod-service/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestProcessUpload(t *testing.T) {

	storage := mocks.NewMockStorage()
	kafka := mocks.NewMockKafka()

	epodSvc := service.NewEPODService(
		storage,
		kafka,
	)

	req := domain.UploadRequest{
		AWB:       "AWB-001",
		CourierID: "CR-001",
		Latitude:  -6.2,
		Longitude: 106.8,
		FileName:  "proof.jpg",
	}

	resp, err := epodSvc.ProcessUpload(req)

	assert.NoError(t, err)
	assert.Equal(t, "SUCCESS", resp.Status)
	assert.NotEmpty(t, resp.ImageURL)
}