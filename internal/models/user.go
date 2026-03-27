package models

import (
	"errors"
	"time"
)

// Tambahkan tag `gorm` untuk PostgreSQL
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Nama      string `gorm:"type:varchar(100);not null"`
	Email     string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string `gorm:"not null"` // Tambahan untuk JWT nanti
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Struct untuk Payload Request
type RegisterRequest struct {
	Nama  string `json:"nama" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// Method untuk Validasi
func (req RegisterRequest) ValidateCustomBusinessLogic() error {
	if req.Nama == "admin" {
		return errors.New("nama 'admin' tidak boleh digunakan")
	}
	return nil
}

// Interface: Kontrak yang harus dipatuhi oleh folder 'repository'
type UserRepository interface {
	SimpanUser(user *User) error
	CariBerdasarkanEmail(email string) (*User, error) // Tambahan method untuk fitur Login
}
