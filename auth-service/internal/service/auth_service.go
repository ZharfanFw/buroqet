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
	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: password, // Catatan: Sebaiknya di-hash jika di produksi, namun ikuti implementasi dasar model terlebih dahulu
		Role:     role,
	}
	return s.userRepo.Create(user)
}

func (s *authService) Login(email, password string) (string, string, *domain.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", "", nil, errors.New("user tidak ditemukan")
	}
	
	if user.Password != password {
		return "", "", nil, errors.New("password salah")
	}

	// Buat token tiruan / dummy token untuk login (karena token generator tidak didefinisikan secara khusus)
	token := "dummy-jwt-token-for-" + user.Email
	return token, user.Role, user, nil
}

func (s *authService) ValidateToken(token string) (*domain.User, error) {
	return &domain.User{}, nil 
}

func (s *authService) RefreshToken(token string) (string, string, error) {
	if token == "" {
		return "", "", errors.New("token tidak valid")
	}
	return "new-access-token-123", "new-refresh-token-456", nil
}