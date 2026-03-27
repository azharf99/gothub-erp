package middleware

import (
	"net/http"
	"strings"

	"github.com/azharf99/gothub-erp/internal/utils"
	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak, token tidak ditemukan"})
			return
		}

		// 2. Format token harus "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid"})
			return
		}

		tokenString := parts[1]

		// 3. Validasi token menggunakan utility kita
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah kedaluwarsa"})
			return
		}

		// 4. Pastikan ini adalah Access Token, bukan Refresh Token
		if claims.Type != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Harap gunakan Access Token"})
			return
		}

		// 5. Simpan UserID ke dalam Context agar bisa dipakai oleh Handler selanjutnya
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		// Lanjut ke Handler tujuan
		c.Next()
	}
}
