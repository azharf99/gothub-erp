package routes

import (
	"github.com/azharf99/gothub-erp/internal/handler"
	"github.com/azharf99/gothub-erp/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler *handler.UserHandler) {
	api := router.Group("/api/v1")
	{
		// 🔓 RUTE PUBLIK
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		// 🔒 RUTE TERLINDUNGI (Wajib Login)
		protected := api.Group("/")
		protected.Use(middleware.RequireAuth())
		{
			// Semua yang sudah login bisa lihat profil sendiri
			protected.GET("/profile", userHandler.GetProfile)

			// 👨‍🏫 GRUP GURU (Hanya Guru dan Admin yang bisa masuk)
			guruRoutes := protected.Group("/grades")
			guruRoutes.Use(middleware.RequireRole("Guru", "Admin"))
			{
				// Contoh: POST /api/v1/grades (Input nilai siswa)
				// guruRoutes.POST("/", gradeHandler.InputGrade)
			}

			// 👨‍🎓 GRUP SISWA (Siswa, Guru, Admin bisa akses)
			siswaRoutes := protected.Group("/schedules")
			siswaRoutes.Use(middleware.RequireRole("Siswa", "Guru", "Admin"))
			{
				// Contoh: GET /api/v1/schedules (Lihat jadwal pelajaran)
				// siswaRoutes.GET("/", scheduleHandler.GetSchedules)
			}

			// 👑 GRUP SUPER ADMIN (Hanya Admin yang bisa masuk)
			adminRoutes := protected.Group("/teachers-data")
			adminRoutes.Use(middleware.RequireRole("Admin"))
			{
				// Contoh: GET /api/v1/teachers-data (Lihat data kepegawaian guru)
				// adminRoutes.GET("/", adminHandler.GetAllTeachers)
			}
		}
	}
}
