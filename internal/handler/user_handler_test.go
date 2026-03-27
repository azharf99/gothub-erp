package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/azharf99/gothub-erp/internal/models"
)

// ==========================================
// 1. DATABASE PALSU (MOCK REPOSITORY)
// ==========================================
type MockUserRepository struct{}

func (m *MockUserRepository) SimpanUser(user *models.User) error {
	user.ID = 1
	return nil
}

// Simulasi pencarian User di Database
func (m *MockUserRepository) CariBerdasarkanEmail(email string) (*models.User, error) {
	if email == "azhar@example.com" {
		// Kita harus membuat hash dari "password123" seolah-olah data ini diambil dari DB
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		return &models.User{
			ID:       1,
			Nama:     "Azhar",
			Email:    "azhar@example.com",
			Password: string(hashedPassword), // Mengirimkan versi hash ke Handler
		}, nil
	}
	return nil, errors.New("user tidak ditemukan")
}

func (m *MockUserRepository) AmbilSemuaUser() ([]models.User, error) {
	return []models.User{
		{ID: 1, Nama: "Azhar", Role: "Admin"},
		{ID: 2, Nama: "Budi", Role: "Guru"},
	}, nil
}

func (m *MockUserRepository) AmbilUserByID(id uint) (*models.User, error) {
	if id == 1 {
		return &models.User{ID: 1, Nama: "Azhar", Role: "Admin"}, nil
	}
	return nil, errors.New("user tidak ditemukan")
}

func (m *MockUserRepository) UpdateUser(user *models.User) error {
	return nil
}

func (m *MockUserRepository) HapusUser(id uint) error {
	return nil
}

// ==========================================
// 2. TEST REGISTER (Sudah ada sebelumnya)
// ==========================================
func TestRegister_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &MockUserRepository{}
	handler := &UserHandler{Repo: mockRepo}
	router := gin.Default()
	router.POST("/register", handler.Register)

	jsonBody := []byte(`{"nama": "Azhar", "email": "azhar@example.com", "password": "password123", "role": "Guru"}`)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Registrasi berhasil")
}

// ==========================================
// 3. TEST LOGIN - SKENARIO SUKSES
// ==========================================
func TestLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &MockUserRepository{}
	handler := &UserHandler{Repo: mockRepo}
	router := gin.Default()
	router.POST("/login", handler.Login)

	// Mengirim email yang benar dan password asli ("password123")
	jsonBody := []byte(`{"email": "azhar@example.com", "password": "password123"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validasi bahwa statusnya 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Validasi bahwa respons JSON mengandung "access_token" dan "refresh_token"
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.NoError(t, err)
	assert.Contains(t, response, "access_token", "Login sukses harus mengembalikan access_token")
	assert.Contains(t, response, "refresh_token", "Login sukses harus mengembalikan refresh_token")
}

// ==========================================
// 4. TEST LOGIN - SKENARIO PASSWORD SALAH
// ==========================================
func TestLogin_WrongPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &MockUserRepository{}
	handler := &UserHandler{Repo: mockRepo}
	router := gin.Default()
	router.POST("/login", handler.Login)

	// Sengaja mengirim password yang salah
	jsonBody := []byte(`{"email": "azhar@example.com", "password": "password_salah"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validasi bahwa statusnya 401 Unauthorized karena password salah
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Email atau password salah")
}

// ==========================================
// 5. TEST GET ALL USERS (SKENARIO ADMIN)
// ==========================================
func TestGetAllUsers_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &UserHandler{Repo: &MockUserRepository{}}
	router := gin.Default()

	// Kita tidak perlu menguji Middleware JWT di sini karena itu tugas tes integrasi,
	// kita langsung arahkan rute ke fungsi Handler.
	router.GET("/users", handler.GetAllUsers)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Azhar")
	assert.Contains(t, w.Body.String(), "Budi")
}

// ==========================================
// 6. TEST UPDATE USER (SKENARIO ADMIN)
// ==========================================
func TestUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &MockUserRepository{}
	handler := &UserHandler{Repo: mockRepo}
	router := gin.Default()
	router.PUT("/users/:id", handler.UpdateUser)

	// Simulasi update user dengan ID 1
	jsonBody := []byte(`{"nama": "Azhar Faturohman Ahidin", "email": "azhar.faturohman@example.com", "role": "Admin"}`)
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validasi bahwa user sudah berubah namanya sesuai dengan logika di MockUserRepository.UpdateUser
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Azhar Faturohman Ahidin")
}

// ==========================================
// 7. TEST DELETE USER (SKENARIO ADMIN)
// ==========================================
func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &MockUserRepository{}
	handler := &UserHandler{Repo: mockRepo}
	router := gin.Default()
	router.DELETE("/users/:id", handler.DeleteUser)

	// Simulasi delete user dengan ID 1
	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validasi bahwa user sudah dihapus sesuai dengan logika di MockUserRepository.DeleteUser
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Pengguna berhasil dihapus") // Karena ID dihapus, kita bisa cek bahwa ID sekarang 0 atau tidak ada
}
