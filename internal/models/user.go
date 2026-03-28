package models

import (
	"errors"
	"time"
)

type User struct {
	ID        uint     `gorm:"primaryKey"`
	Nama      string   `gorm:"type:varchar(100);not null"`
	Email     string   `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string   `gorm:"not null"`
	Role      string   `gorm:"type:varchar(20);not null;default:'Siswa'"`
	Courses   []Course `gorm:"foreignKey:TeacherID" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RegisterRequest struct {
	Nama     string `json:"nama" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"omitempty,oneof=Admin Guru Siswa Karyawan"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Nama  string `json:"nama" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"omitempty,oneof=Admin Guru Siswa Karyawan"`
}

func (req RegisterRequest) ValidateCustomBusinessLogic() error {
	if req.Nama == "admin" {
		return errors.New("nama 'admin' tidak boleh digunakan")
	}
	if req.Role == "" {
		req.Role = "Siswa"
	}
	return nil
}

type UserRepository interface {
	SimpanUser(user *User) error
	CariBerdasarkanEmail(email string) (*User, error)
	AmbilSemuaUser(page int, limit int) ([]User, int64, error)
	AmbilUserByID(id uint) (*User, error)
	UpdateUser(user *User) error
	HapusUser(id uint) error
}

type UserService interface {
	RegisterUser(req RegisterRequest) (*User, error)
	LoginUser(req LoginRequest) (string, string, error) // Return: accessToken, refreshToken, error
	CreateUserFromDashboard(req RegisterRequest, currentUserRole string) (*User, error)
	GetSemuaUser(page, limit int) ([]User, int64, error)
	UpdateDataUser(id uint, req UpdateUserRequest) (*User, error)
	HapusDataUser(id uint) error
}
