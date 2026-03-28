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

// READ: Ambil Semua User (Sembunyikan password menggunakan Select)
func (r *PostgresRepo) AmbilSemuaUser(page int, limit int) ([]models.User, int64, error) {
	var users []models.User
	var totalItems int64

	// Hitung total keseluruhan data di tabel courses (tanpa limit/offset)
	r.DB.Model(&models.User{}).Count(&totalItems)

	// Hitung Offset (Berapa data yang harus dilewati)
	// Rumus: (Halaman Saat Ini - 1) * Limit
	offset := (page - 1) * limit

	err := r.DB.Select("id", "nama", "email", "role", "created_at", "updated_at").Offset(offset).Limit(limit).Find(&users).Error
	return users, totalItems, err
}

// READ: Ambil Satu User
func (r *PostgresRepo) AmbilUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.Select("id", "nama", "email", "password", "role", "created_at", "updated_at").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UPDATE: Perbarui data User
func (r *PostgresRepo) UpdateUser(user *models.User) error {
	// GORM otomatis mengupdate berdasarkan primary key (user.ID)
	return r.DB.Save(user).Error
}

// DELETE: Hapus User
func (r *PostgresRepo) HapusUser(id uint) error {
	return r.DB.Delete(&models.User{}, id).Error
}
