package routes

import (
	"github.com/azharf99/gothub-erp/internal/handler"
	"github.com/azharf99/gothub-erp/internal/middleware"
	"github.com/gin-gonic/gin"
)

// Fungsi untuk mendaftarkan semua rute
func SetupRoutes(router *gin.Engine, userHandler *handler.UserHandler) {
	api := router.Group("/api/v1")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		// Nanti bisa tambah rute lain di sini: api.GET("/users", userHandler.GetAll)
		api.POST("/refresh-token", userHandler.RefreshToken)

		// Rute Terlindungi (Harus bawa Access Token di header Authorization)
		protected := api.Group("/")
		protected.Use(middleware.RequireAuth()) // Pasang gembok di grup ini!
		{
			protected.GET("/profile", userHandler.GetProfile)
			// Nanti bisa tambah: protected.POST("/products", productHandler.Create) dll.
		}
	}
}
