package repository

import (
	"auth-service/internal/domain"
	"gorm.io/gorm"
)

type userPostgres struct {
	db *gorm.DB
}

func NewUserPostgres(db *gorm.DB) domain.UserRepository {
	return &userPostgres{db: db}
}

func (r *userPostgres) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userPostgres) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userPostgres) GetActiveDrivers() ([]domain.User, error) {
	var drivers []domain.User
	// Contoh query sederhana untuk mendapatkan kurir
	err := r.db.Where("role = ?", "kurir").Find(&drivers).Error
	return drivers, err
}