package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/azharf99/gothub-erp/internal/database"
	"github.com/azharf99/gothub-erp/internal/handler"
	"github.com/azharf99/gothub-erp/internal/models"
	"github.com/azharf99/gothub-erp/internal/repository"
	"github.com/azharf99/gothub-erp/internal/routes"
	"github.com/azharf99/gothub-erp/internal/service"
	"github.com/gin-contrib/cors"
)

func main() {
	// ==========================================
	// 1. LOAD ENVIRONMENT VARIABLES (.env)
	// ==========================================
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment system")
	}

	// ==========================================
	// 2. KONFIGURASI DATABASE POSTGRESQL
	// ==========================================
	// Merangkai Data Source Name (DSN) dari variabel .env
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Membuka koneksi ke PostgreSQL menggunakan GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database PostgreSQL:", err)
	}
	fmt.Println("Berhasil terhubung ke database PostgreSQL!")

	// ==========================================
	// 3. AUTO MIGRATION (Keajaiban GORM)
	// ==========================================
	// GORM akan membaca struct User dan otomatis membuatkan tabel 'users'
	// lengkap dengan kolom, tipe data, dan primary key-nya jika belum ada.
	err = db.AutoMigrate(&models.User{}, &models.Course{})
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database:", err)
	}
	fmt.Println("Migrasi database berhasil!")

	// 3. 2 JALANKAN SEEDER DI SINI
	database.SeedSuperAdmin(db)
	fmt.Println("Seeding data berhasil!")

	// ==========================================
	// 4. DEPENDENCY INJECTION
	// ==========================================
	// Suntikkan koneksi DB 'db' ke dalam PostgresRepo
	// ==========================================
	// INIT REPOSITORY
	// ==========================================
	userRepo := repository.NewPostgresRepo(db)
	courseRepo := repository.NewCourseRepo(db)

	// ==========================================
	// INIT SERVICE (Lapisan Baru)
	// ==========================================
	// Impor package service buatanmu di atas: "github.com/azharf99/gothub-erp/internal/service"
	userService := service.NewUserService(userRepo)
	courseService := service.NewCourseService(courseRepo)

	// ==========================================
	// INIT HANDLER
	// ==========================================
	// Handler sekarang menerima Service, bukan Repository
	userHandler := &handler.UserHandler{Service: userService}
	courseHandler := &handler.CourseHandler{Service: courseService}

	// ==========================================
	// 5. SETUP ROUTER & START SERVER
	// ==========================================
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},    // Izinkan frontend lokalmu
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},  // Izinkan semua metode CRUD
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"}, // Izinkan header token
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Browser tidak perlu repot tanya OPTIONS lagi selama 12 jam
	}))

	// Daftarkan semua rute API
	routes.SetupRoutes(router, userHandler, courseHandler)

	// Ambil port dari .env, default ke 8080 jika kosong
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("ERP Server berjalan di port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Gagal menyalakan server:", err)
	}
}
