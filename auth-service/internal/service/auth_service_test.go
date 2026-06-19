package service

import (
	"auth-service/internal/domain"
	"auth-service/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister_Success(t *testing.T) {
	// Setup mock repository
	mockRepo := new(mocks.MockUserRepository)
	svc := NewAuthService(mockRepo)

	// Expectation: Fungsi Create dipanggil dengan data apapun dan return nil (success)
	mockRepo.On("Create", mock.Anything).Return(nil)

	// Action
	err := svc.Register("Imam UPI", "imam@upi.edu", "password123", "kurir")

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := NewAuthService(mockRepo)

	// Expectation: Cari email yang tidak ada, return error
	mockRepo.On("FindByEmail", "salah@upi.edu").Return((*domain.User)(nil), errors.New("not found"))

	// Action
	_, _, _, err := svc.Login("salah@upi.edu", "password123")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "user tidak ditemukan", err.Error())
}