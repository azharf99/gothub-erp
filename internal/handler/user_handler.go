package handler

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/azharf99/gothub-erp/internal/models"
	"github.com/azharf99/gothub-erp/internal/utils"
)

type UserHandler struct {
	Service models.UserService
}

// ==========================================
// REGISTER (Public)
// ==========================================
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewBadRequest(err.Error()))
		return
	}

	newUser, err := h.Service.RegisterUser(req)
	if err != nil {
		c.Error(err)
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Registrasi berhasil, silakan login", newUser)
}

// ==========================================
// CREATE USER (Authenticated Only)
// ==========================================
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewBadRequest(err.Error()))
		return
	}

	finalRole := "Siswa"

	currentUserRole, exists := c.Get("role")

	if exists && currentUserRole == "Admin" && req.Role != "" {
		finalRole = req.Role
	}
	user, err := h.Service.CreateUserFromDashboard(req, finalRole)
	if err != nil {
		c.Error(err)
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Pengguna berhasil ditambahkan", user)
}

// ==========================================
// LOGIN LOGIC (Verifikasi Hash & Buat JWT)
// ==========================================
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewBadRequest("Format JSON tidak sesuai atau data tidak lengkap"))
		return
	}

	accessToken, refreshToken, err := h.Service.LoginUser(req)
	if err != nil {
		c.Error(err)
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
		c.Error(utils.NewBadRequest("Refresh token dibutuhkan"))
		return
	}

	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil {
		c.Error(utils.NewUnauthorized("Refresh token tidak valid atau kedaluwarsa. Silakan login ulang."))
		return
	}

	if claims.Type != "refresh" {
		c.Error(utils.NewUnauthorized("Token yang diberikan bukan refresh token"))
		return
	}

	newAccessToken, newRefreshToken, err := utils.GenerateTokens(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		c.Error(utils.NewInternalError("Gagal membuat token baru"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Token berhasil diperbarui",
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

// ==========================================
// PROFIL USER
// ==========================================
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, email, role, err := utils.GetCurrentUser(c)
	if err != nil {
		c.Error(err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Selamat datang di area terlarang!", gin.H{
		"user_id": userID,
		"email":   email,
		"role":    role,
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
	utils.SendSuccess(c, http.StatusOK, "Logout berhasil. Token akan dihapus di sisi klien.", nil)
}

// ==========================================
// GET ALL USERS (Hanya Admin)
// ==========================================
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)

	users, totalItems, err := h.Service.GetSemuaUser(page, limit)
	if err != nil {
		c.Error(err)
		return
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	meta := utils.PaginationMeta{
		CurrentPage: page,
		Limit:       limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}

	utils.SendPaginatedSuccess(c, http.StatusOK, "Berhasil mengambil daftar pengguna", users, meta)
}

// ==========================================
// UPDATE USER
// ==========================================
func (h *UserHandler) UpdateUser(c *gin.Context) {
	targetID, err := utils.GetParamID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewBadRequest(err.Error()))
		return
	}

	user, err := h.Service.UpdateDataUser(targetID, req)
	if err != nil {
		c.Error(err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Berhasil memperbarui data pengguna", user)
}

// ==========================================
// DELETE USER (Hanya Admin)
// ==========================================
func (h *UserHandler) DeleteUser(c *gin.Context) {
	targetID, err := utils.GetParamID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.Service.HapusDataUser(targetID); err != nil {
		c.Error(utils.NewInternalError("Gagal menghapus pengguna"))
		return
	}
	utils.SendSuccess(c, http.StatusOK, "Pengguna berhasil dihapus", nil)
}
