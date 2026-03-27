package middleware

import (
	"errors"
	"net/http"

	"github.com/azharf99/gothub-erp/internal/utils"
	"github.com/gin-gonic/gin"
)

// GlobalErrorHandler akan mencegat semua error yang dilempar oleh Handler
func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lanjutkan eksekusi ke Handler tujuan
		c.Next()

		// Setelah Handler selesai, cek apakah ada error yang dikumpulkan
		if len(c.Errors) > 0 {
			// Ambil error yang paling terakhir dilempar
			err := c.Errors.Last().Err

			var appErr *utils.AppError

			// Cek apakah ini adalah custom AppError buatan kita?
			if errors.As(err, &appErr) {
				// Jika ya, gunakan kode dan pesan dari AppError kita
				c.JSON(appErr.Code, gin.H{
					"success": false,
					"message": appErr.Message,
				})
			} else {
				// Jika ini error bawaan Golang yang tidak terduga (misal: panic, database putus)
				// Kembalikan 500 Internal Server Error untuk keamanan agar jeroan sistem tidak bocor
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Terjadi kesalahan internal pada server",
				})
			}
		}
	}
}
