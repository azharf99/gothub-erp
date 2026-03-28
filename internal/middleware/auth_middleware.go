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
			c.AbortWithError(http.StatusUnauthorized, utils.NewUnauthorized("Akses ditolak, token tidak ditemukan"))
			return
		}

		// 2. Format token harus "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithError(http.StatusUnauthorized, utils.NewUnauthorized("Format token tidak valid"))
			return
		}

		tokenString := parts[1]

		// 3. Validasi token menggunakan utility kita
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, utils.NewUnauthorized("Token tidak valid atau sudah kedaluwarsa"))
			return
		}

		// 4. Pastikan ini adalah Access Token, bukan Refresh Token
		if claims.Type != "access" {
			c.AbortWithError(http.StatusUnauthorized, utils.NewUnauthorized("Harap gunakan Access Token"))
			return
		}

		// 5. Simpan UserID ke dalam Context agar bisa dipakai oleh Handler selanjutnya
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role) // <<< SIMPAN ROLE KE CONTEXT

		// Lanjut ke Handler tujuan
		c.Next()
	}
}

// RequireRole menerima daftar role yang diizinkan (variadic parameter)
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil role user dari context (yang sebelumnya diisi oleh RequireAuth)
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithError(http.StatusForbidden, utils.NewForbidden("Role tidak ditemukan"))
			return
		}

		// Cek apakah role user ada di dalam daftar role yang diizinkan
		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		// Jika tidak diizinkan, tolak aksesnya!
		if !isAllowed {
			c.AbortWithError(http.StatusForbidden, utils.NewForbidden("Akses ditolak: Anda tidak memiliki izin untuk mengakses resource ini"))
			return
		}

		// Jika lolos, silakan lanjut ke Handler
		c.Next()
	}
}
