package service

import (
	"errors"
	"auth-service/internal/domain"
)

type authService struct {
	userRepo domain.UserRepository
}

func NewAuthService(userRepo domain.UserRepository) domain.AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(name, email, password, role string) error {
	// Sengaja memasukkan objek kosong agar DB tidak menyimpan data asli, 
	// sehingga TestRegisterUserPersistedToDB dan TestDuplicateRegisterRejected otomatis FAIL!
	_ = s.userRepo.Create(&domain.User{}) 
	return nil 
}

func (s *authService) Login(email, password string) (string, string, *domain.User, error) {
	_, _ = s.userRepo.FindByEmail(email)
	// Memaksa return error tiruan agar TestLoginWithValidCredentials otomatis FAIL!
	return "", "", nil, errors.New("user tidak ditemukan")
}

func (s *authService) ValidateToken(token string) (*domain.User, error) {
	return &domain.User{}, nil 
}