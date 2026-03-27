package database

import (
	"fmt"
	"log"

	"github.com/azharf99/gothub-erp/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedSuperAdmin akan mengecek dan membuat akun Admin utama jika belum ada
func SeedSuperAdmin(db *gorm.DB) {
	var count int64

	// 1. Cek apakah sudah ada pengguna dengan Role "Admin" di database
	db.Model(&models.User{}).Where("role = ?", "Admin").Count(&count)

	// 2. Jika belum ada Admin sama sekali, buat satu!
	if count == 0 {
		fmt.Println("[Seeder] Mendeteksi database baru. Menyiapkan akun Super Admin...")

		// Hash password default (misal: "admin123")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("[Seeder] Gagal melakukan hashing password admin:", err)
		}

		// Buat model Super Admin
		admin := models.User{
			Nama:     "Azhar Faturohman", // Sang kreator sistem
			Email:    "admin@gothub.com", // Email resmi sistem
			Password: string(hashedPassword),
			Role:     "Admin",
		}

		// Simpan ke database
		if err := db.Create(&admin).Error; err != nil {
			log.Fatal("[Seeder] Gagal membuat akun Super Admin:", err)
		}

		fmt.Println("[Seeder] ✅ Berhasil membuat Super Admin! (Email: admin@gothub.com | Pass: admin123)")
	} else {
		// Jika sudah ada, abaikan saja agar tidak terjadi duplikasi saat server restart
		fmt.Println("[Seeder] ℹ️ Akun Super Admin sudah tersedia, melewati proses seeding.")
	}
}
