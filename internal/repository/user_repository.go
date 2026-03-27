package repository

import (
	"github.com/azharf99/gothub-erp/internal/models"
	"gorm.io/gorm"
)

type PostgresRepo struct {
	DB *gorm.DB // Menggunakan pointer ke koneksi GORM
}

// Konstruktor ini akan dipanggil di main.go nanti
func NewPostgresRepo(db *gorm.DB) models.UserRepository {
	return &PostgresRepo{DB: db}
}

// 1. Fungsi Menyimpan User Baru (Register)
func (r *PostgresRepo) SimpanUser(user *models.User) error {
	// GORM otomatis menerjemahkan ini menjadi: INSERT INTO users (...) VALUES (...)
	return r.DB.Create(user).Error
}

// 2. Fungsi Mencari User (Login)
func (r *PostgresRepo) CariBerdasarkanEmail(email string) (*models.User, error) {
	var user models.User
	// GORM otomatis menerjemahkan ini menjadi: SELECT * FROM users WHERE email = ? LIMIT 1
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
