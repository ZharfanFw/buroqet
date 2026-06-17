package domain

import "time"

type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"size:100;not null"`
    Email     string    `json:"email" gorm:"size:100;uniqueIndex;not null"`
    Password  string    `json:"-" gorm:"not null"` // "-" agar password tidak bocor di JSON
    Role      string    `json:"role" gorm:"size:20;not null"` // pelanggan, kurir, admin
    CreatedAt time.Time `json:"created_at"`
}

// Interface untuk abstraksi sesuai struktur image_8ebfc0.png
type UserRepository interface {
    Create(user *User) error
    FindByEmail(email string) (*User, error)
    GetActiveDrivers() ([]User, error)
}

type AuthService interface {
    Register(name, email, password, role string) error
    Login(email, password string) (accessToken, refreshToken string, user *User, err error)
    ValidateToken(token string) (*User, error)
} 