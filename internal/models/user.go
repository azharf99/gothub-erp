package models

import (
	"errors"
	"time"
)

// Menambahkan field Role
type User struct {
	ID        uint     `gorm:"primaryKey"`
	Nama      string   `gorm:"type:varchar(100);not null"`
	Email     string   `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string   `gorm:"not null"`
	Role      string   `gorm:"type:varchar(20);not null;default:'Siswa'"` // Admin, Guru, Siswa, Karyawan
	Courses   []Course `gorm:"foreignKey:TeacherID" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RegisterRequest struct {
	Nama     string `json:"nama" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=Admin Guru Siswa Karyawan"` // Validasi ketat Gin
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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
