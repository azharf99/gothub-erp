package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/azharf99/gothub-erp/internal/database"
	"github.com/azharf99/gothub-erp/internal/handler"
	"github.com/azharf99/gothub-erp/internal/middleware"
	"github.com/azharf99/gothub-erp/internal/models"
	"github.com/azharf99/gothub-erp/internal/repository"
	"github.com/azharf99/gothub-erp/internal/routes"
	"github.com/azharf99/gothub-erp/internal/service" // Pastikan import service buatanmu ada di sini
)

func main() {
	// ==========================================
	// 1. LOAD ENVIRONMENT VARIABLES (.env)
	// ==========================================
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, membaca environment variable dari sistem (Docker/GCP)")
	}

	// Atur Gin Mode sesuai environment (release untuk GCP, debug untuk lokal)
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// ==========================================
	// 2. KONFIGURASI DATABASE POSTGRESQL
	// ==========================================
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database PostgreSQL:", err)
	}
	fmt.Println("Berhasil terhubung ke database PostgreSQL!")

	// ==========================================
	// 3. AUTO MIGRATION & SEEDER
	// ==========================================
	err = db.AutoMigrate(&models.User{}, &models.Course{})
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database:", err)
	}
	fmt.Println("Migrasi database berhasil!")

	database.SeedSuperAdmin(db)

	// ==========================================
	// 4. DEPENDENCY INJECTION (CLEAN ARCHITECTURE)
	// ==========================================
	// REPOSITORY LAYER
	userRepo := repository.NewPostgresRepo(db)
	courseRepo := repository.NewCourseRepo(db)

	// SERVICE LAYER (Otak Aplikasi)
	userService := service.NewUserService(userRepo)
	courseService := service.NewCourseService(courseRepo)

	// HANDLER LAYER (Menerima Service, bukan Repository)
	userHandler := &handler.UserHandler{Service: userService}
	courseHandler := &handler.CourseHandler{Service: courseService}

	// ==========================================
	// 5. SETUP ROUTER & KEAMANAN
	// ==========================================
	router := gin.Default()

	// KEAMANAN: Mematikan kepercayaan proxy secara default untuk mencegah spoofing IP.
	// Jika nanti kamu butuh IP asli user di belakang GCP Load Balancer, ubah nil menjadi IP Load Balancer GCP.
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Println("Peringatan gagal mengatur Trusted Proxies:", err)
	}

	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.SetupCORS())
	router.Use(middleware.RateLimiter())

	// Daftarkan semua rute API
	routes.SetupRoutes(router, userHandler, courseHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("ERP Server berjalan dalam mode %s di port %s...\n", gin.Mode(), port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Gagal menyalakan server:", err)
	}
}
