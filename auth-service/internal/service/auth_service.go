package service

import (
	"errors"
	"os"
	"time"

	"auth-service/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(getEnv("JWT_SECRET", "super-secret-key-buroqet-123"))

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

type authService struct {
	userRepo domain.UserRepository
}

func NewAuthService(userRepo domain.UserRepository) domain.AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(name, email, password, role string) error {
	// Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal memproses password")
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}
	return s.userRepo.Create(user)
}

func (s *authService) Login(email, password string) (string, string, *domain.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", "", nil, errors.New("email atau password salah")
	}

	// Bandingkan password plaintext dengan hash di database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", nil, errors.New("email atau password salah")
	}

	// Generate JWT Tokens
	accessToken, refreshToken, err := s.generateTokenPair(user)
	if err != nil {
		return "", "", nil, errors.New("gagal men-generate token")
	}

	return accessToken, refreshToken, user, nil
}

func (s *authService) ValidateToken(tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("klaim token tidak valid")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("email tidak ditemukan di token")
	}

	return s.userRepo.FindByEmail(email)
}

func (s *authService) RefreshToken(tokenString string) (string, string, error) {
	if tokenString == "" {
		return "", "", errors.New("token tidak valid")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("refresh token tidak valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("klaim token tidak valid")
	}

	// Pastikan ini adalah refresh token
	if tokenType, ok := claims["type"].(string); !ok || tokenType != "refresh" {
		return "", "", errors.New("token bukan refresh token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", "", errors.New("email tidak ditemukan di token")
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("user tidak ditemukan")
	}

	return s.generateTokenPair(user)
}

func (s *authService) generateTokenPair(user *domain.User) (string, string, error) {
	// Access Token (Exp 1 jam)
	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"type":    "access",
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh Token (Exp 7 hari)
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}