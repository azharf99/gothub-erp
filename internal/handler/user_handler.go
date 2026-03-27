package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/azharf99/gothub-erp/internal/models"
	"github.com/azharf99/gothub-erp/internal/utils"
)

type UserHandler struct {
	Repo models.UserRepository
}

// ==========================================
// REGISTER LOGIC (Hash Password dan Validasi Bisnis)
// ==========================================
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hashing Password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password"})
		return
	}

	// Validasi Bisnis Custom (Contoh: Nama tidak boleh "admin")
	validationErr := req.ValidateCustomBusinessLogic()
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	// 2. Siapkan model User untuk disimpan
	newUser := models.User{
		Nama:     req.Nama,
		Email:    req.Email,
		Password: string(hashedPassword), // Simpan versi hash-nya, BUKAN versi aslinya
		Role:     req.Role,
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
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Email, user.Role)
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
	newAccessToken, newRefreshToken, err := utils.GenerateTokens(claims.UserID, claims.Email, claims.Role)
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
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"message": "Selamat datang di area terlarang!",
		"data": gin.H{
			"user_id": userID,
			"email":   email,
			"role":    role,
		},
	})
}

// ==========================================
// LOGOUT
// ==========================================
func (h *UserHandler) Logout(c *gin.Context) {
	// Di sistem stateless murni, kita cukup mengirimkan respons sukses
	// dan menginstruksikan Frontend untuk menghapus token di sisi mereka.
	// Jika ingin level enterprise (Strict), token yang dikirim di header bisa dicegat
	// lalu dimasukkan ke tabel "TokenBlacklist" di database menggunakan Repo.

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout berhasil. Silakan hapus token di sisi klien.",
	})
}

// ==========================================
// GET ALL USERS (Hanya Admin)
// ==========================================
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.Repo.AmbilSemuaUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// ==========================================
// UPDATE USER
// ==========================================
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Mengambil ID dari URL parameter (misal: /users/2)
	idParam := c.Param("id")
	targetID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pengguna tidak valid"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari user yang mau diupdate
	user, err := h.Repo.AmbilUserByID(uint(targetID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	// Update data yang diizinkan
	user.Nama = req.Nama
	user.Email = req.Email
	if req.Role != "" {
		user.Role = req.Role // Hanya izinkan update role jika dikirim di JSON
	}

	if err := h.Repo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data pengguna berhasil diperbarui",
		"data":    user,
	})
}

// ==========================================
// DELETE USER
// ==========================================
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	targetID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pengguna tidak valid"})
		return
	}

	if err := h.Repo.HapusUser(uint(targetID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dihapus"})
}
