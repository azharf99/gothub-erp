package routes

import (
	"github.com/azharf99/gothub-erp/internal/handler"
	"github.com/azharf99/gothub-erp/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler *handler.UserHandler, courseHandler *handler.CourseHandler) {
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
			// Endpoint Logout (Siapapun yang login bisa logout)
			protected.POST("/logout", userHandler.Logout)
			protected.GET("/courses", courseHandler.GetAllCourses) // <<< READ

			// 👨‍🏫 GRUP KHUSUS PENGAJAR & ADMIN
			teacherRoutes := protected.Group("/courses")
			teacherRoutes.Use(middleware.RequireRole("Guru", "Admin"))
			{
				teacherRoutes.POST("/", courseHandler.CreateCourse)      // <<< CREATE
				teacherRoutes.PUT("/:id", courseHandler.UpdateCourse)    // <<< TAMBAHAN UPDATE
				teacherRoutes.DELETE("/:id", courseHandler.DeleteCourse) // <<< TAMBAHAN DELETE
			}

			// 👑 GRUP MANAJEMEN USER (Hanya Super Admin)
			adminRoutes := protected.Group("/users")
			adminRoutes.Use(middleware.RequireRole("Admin"))
			{
				adminRoutes.GET("/", userHandler.GetAllUsers)      // READ ALL
				adminRoutes.PUT("/:id", userHandler.UpdateUser)    // UPDATE
				adminRoutes.DELETE("/:id", userHandler.DeleteUser) // DELETE
			}

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
		}
	}
}
