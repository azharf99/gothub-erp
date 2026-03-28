package handler

import (
	"math"
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
// REGISTER (Public)
// ==========================================
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewBadRequest(err.Error()))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Error(utils.NewInternalError("Gagal memproses password"))
		return
	}

	if validationErr := req.ValidateCustomBusinessLogic(); validationErr != nil {
		c.Error(utils.NewBadRequest(validationErr.Error()))
		return
	}

	newUser := models.User{
		Nama:     req.Nama,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "Siswa",
	}

	if err := h.Repo.SimpanUser(&newUser); err != nil {
		c.Error(utils.NewInternalError("Gagal menyimpan user ke database"))
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

	existingUser, _ := h.Repo.CariBerdasarkanEmail(req.Email)
	if existingUser != nil {
		c.Error(utils.NewBadRequest("Email sudah terdaftar"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Error(utils.NewInternalError("Gagal memproses password"))
		return
	}

	finalRole := "Siswa"

	currentUserRole, exists := c.Get("role")

	if exists && currentUserRole == "Admin" && req.Role != "" {
		finalRole = req.Role
	}

	user := models.User{
		Nama:     req.Nama,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     finalRole,
	}

	if err := h.Repo.SimpanUser(&user); err != nil {
		c.Error(utils.NewInternalError("Gagal menyimpan data pengguna"))
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

	user, err := h.Repo.CariBerdasarkanEmail(req.Email)
	if err != nil || user == nil {
		c.Error(utils.NewUnauthorized("Email atau password salah"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.Error(utils.NewUnauthorized("Email atau password salah"))
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		c.Error(utils.NewInternalError("Gagal membuat token autentikasi"))
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

	users, totalItems, err := h.Repo.AmbilSemuaUser(page, limit)
	if err != nil {
		c.Error(utils.NewInternalError("Gagal mengambil data pengguna"))
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

	user, err := h.Repo.AmbilUserByID(targetID)
	if err != nil {
		c.Error(utils.NewNotFound("Pengguna tidak ditemukan"))
		return
	}

	user.Nama = req.Nama
	user.Email = req.Email
	if req.Role != "" {
		user.Role = req.Role
	}

	if err := h.Repo.UpdateUser(user); err != nil {
		c.Error(utils.NewInternalError("Gagal memperbarui data pengguna"))
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

	if err := h.Repo.HapusUser(targetID); err != nil {
		c.Error(utils.NewInternalError("Gagal menghapus pengguna"))
		return
	}
	utils.SendSuccess(c, http.StatusOK, "Pengguna berhasil dihapus", nil)
}
