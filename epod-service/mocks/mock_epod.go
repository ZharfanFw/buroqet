package mocks

import (
	"context"
	"io"

	"epod-service/internal/repository"
	"epod-service/internal/kafka"

	"github.com/stretchr/testify/mock"
)

// === 1. MOCK REPOSITORY ===
type EPODRepository struct {
	mock.Mock
}

func (m *EPODRepository) Create(epod *repository.EPOD) error {
	args := m.Called(epod)
	return args.Error(0)
}

func (m *EPODRepository) FindAll(status string, limit int) ([]repository.EPOD, error) {
	args := m.Called(status, limit)
	return args.Get(0).([]repository.EPOD), args.Error(1)
}

func (m *EPODRepository) FindByID(id string) (*repository.EPOD, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.EPOD), args.Error(1)
}

func (m *EPODRepository) FindByAWB(awb string) (*repository.EPOD, error) {
	args := m.Called(awb)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.EPOD), args.Error(1)
}

func (m *EPODRepository) UpdateStatus(id string, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

// === 2. MOCK STORAGE ===
type StorageInterface struct {
	mock.Mock
}

func (m *StorageInterface) Upload(file io.Reader, fileName string) (string, error) {
	args := m.Called(file, fileName)
	return args.String(0), args.Error(1)
}

// === 3. MOCK KAFKA PRODUCER ===
type ProducerInterface struct {
	mock.Mock
}

func (m *ProducerInterface) PublishEvent(ctx context.Context, topic string, event kafka.EPODEvent) error {
	args := m.Called(ctx, topic, event)
	return args.Error(0)
}

func (m *ProducerInterface) Close() error {
	args := m.Called()
	return args.Error(0)
}