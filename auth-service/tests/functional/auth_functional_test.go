package functional_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"auth-service/internal/domain"
	"auth-service/internal/repository"
	"auth-service/internal/service"
)

type AuthFunctionalSuite struct {
	suite.Suite
	db      *gorm.DB
	service domain.AuthService 
	ctx     context.Context
}

func (s *AuthFunctionalSuite) SetupSuite() {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost user=postgres password=password dbname=db_test port=5433 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	s.Require().NoError(err, "Gagal koneksi ke database test")

	err = db.AutoMigrate(&domain.User{})
	s.Require().NoError(err, "Gagal migrasi database test")

	s.db = db
	s.ctx = context.Background()

	// SEKARANG SUDAH FIX: Memakai konstruktor yang benar dari user_postgres.go
	repo := repository.NewUserPostgres(db)
	s.service = service.NewAuthService(repo)
}

func (s *AuthFunctionalSuite) SetupTest() {
	s.db.Exec("DELETE FROM users")
}

func (s *AuthFunctionalSuite) TearDownSuite() {
	sqlDB, _ := s.db.DB()
	sqlDB.Close()
}

func (s *AuthFunctionalSuite) TestRegisterUserPersistedToDB() {
	name := "Admin Test"
	email := "admin@test.com"
	password := "securepassword123"
	role := "admin"

	err := s.service.Register(name, email, password, role)
	s.NoError(err)

	var user domain.User
	result := s.db.Where("email = ?", email).First(&user)
	s.NoError(result.Error)
	s.Equal(email, user.Email)
	s.Equal(name, user.Name)
	s.NotEqual(password, user.Password) 
}

func (s *AuthFunctionalSuite) TestLoginWithValidCredentials() {
	name := "User Login"
	email := "user@login.com"
	password := "secretpass"
	role := "pelanggan"

	s.Require().NoError(s.service.Register(name, email, password, role))

	accessToken, refreshToken, user, err := s.service.Login(email, password)
	s.NoError(err)
	s.NotEmpty(accessToken)
	s.NotEmpty(refreshToken)
	s.NotNil(user)
	s.Equal(email, user.Email)
}

func (s *AuthFunctionalSuite) TestLoginWithInvalidPassword() {
	email := "wrongpass@test.com"
	password := "correctpassword"

	s.Require().NoError(s.service.Register("Wrong Pass User", email, password, "pelanggan"))

	_, _, _, err := s.service.Login(email, "wrongpassword")
	s.Error(err, "Harus mengembalikan error karena password tidak cocok")
}

func (s *AuthFunctionalSuite) TestDuplicateRegisterRejected() {
	email := "duplicate@test.com"

	err := s.service.Register("User 1", email, "password123", "pelanggan")
	s.NoError(err)

	err = s.service.Register("User 2", email, "password456", "pelanggan")
	s.Error(err, "Harus mengembalikan error karena email sudah terdaftar")
}

func TestAuthFunctionalSuite(t *testing.T) {
	suite.Run(t, new(AuthFunctionalSuite))
}