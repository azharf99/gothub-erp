package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/azharf99/gothub-erp/internal/models"
	"github.com/azharf99/gothub-erp/internal/utils"
)

type UserHandler struct {
	Repo models.UserRepository
}

// ==========================================
// REGISTER LOGIC (Hash Password sebelum simpan)
// ==========================================
func (h *UserHandler) Register(c *gin.Context) {
	// Pastikan struct RegisterRequest di models memiliki field Password
	var req struct {
		Nama     string `json:"nama" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Hashing Password menggunakan bcrypt
	// GenerateFromPassword menerima byte array password dan tingkat kerumitan (Cost)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password"})
		return
	}

	// 2. Siapkan model User untuk disimpan
	newUser := models.User{
		Nama:     req.Nama,
		Email:    req.Email,
		Password: string(hashedPassword), // Simpan versi hash-nya, BUKAN versi aslinya
	}

	// 3. Simpan ke database
	if err := h.Repo.SimpanUser(&newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan user ke database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil, silakan login"})
}

// ==========================================
// LOGIN LOGIC (Verifikasi Hash & Buat JWT)
// ==========================================
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest // Pastikan LoginRequest ada di models (Email & Password)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}

	// 1. Cari user di database berdasarkan Email
	user, err := h.Repo.CariBerdasarkanEmail(req.Email)
	if err != nil || user == nil {
		// Gunakan pesan error yang umum demi keamanan (jangan beritahu apakah email atau password yang salah)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	// 2. Bandingkan password asli dari request dengan password hash dari database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	// 3. Jika password cocok, Generate KEDUA Token
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token autentikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Login berhasil",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// ==========================================
// FUNGSI UNTUK REFRESH TOKEN
// ==========================================
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token dibutuhkan"})
		return
	}

	// Validasi refresh token
	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token tidak valid atau kedaluwarsa. Silakan login ulang."})
		return
	}

	// Pastikan tipe tokennya benar-benar "refresh"
	if claims.Type != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token yang diberikan bukan refresh token"})
		return
	}

	// Buat pasangan token baru
	newAccessToken, newRefreshToken, err := utils.GenerateTokens(claims.UserID, claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Token berhasil diperbarui",
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

// ==========================================
// CONTOH ENDPOINT YANG DILINDUNGI
// ==========================================
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Mengambil data yang diselipkan oleh Middleware ke dalam Context
	userID, _ := c.Get("userID")
	email, _ := c.Get("email")

	c.JSON(http.StatusOK, gin.H{
		"message": "Selamat datang di area terlarang!",
		"data": gin.H{
			"user_id": userID,
			"email":   email,
		},
	})
}
