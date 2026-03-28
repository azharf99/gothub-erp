package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azharf99/gothub-erp/internal/middleware"
	"github.com/azharf99/gothub-erp/internal/models"
	"github.com/azharf99/gothub-erp/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ==========================================
// 1. MOCK SERVICE UNTUK USER
// ==========================================
type MockUserService struct{}

func (m *MockUserService) RegisterUser(req models.RegisterRequest) (*models.User, error) {
	return &models.User{ID: 1, Nama: req.Nama, Email: req.Email, Role: "Siswa"}, nil
}

func (m *MockUserService) LoginUser(req models.LoginRequest) (string, string, error) {
	if req.Email == "azhar@example.com" && req.Password == "password123" {
		return "mock_access_token", "mock_refresh_token", nil
	}
	return "", "", utils.NewUnauthorized("Email atau password salah")
}

func (m *MockUserService) CreateUserFromDashboard(req models.RegisterRequest, currentUserRole string) (*models.User, error) {
	return &models.User{ID: 2, Nama: req.Nama, Role: "Siswa"}, nil
}

func (m *MockUserService) GetSemuaUser(page, limit int) ([]models.User, int64, error) {
	users := []models.User{
		{ID: 1, Nama: "Azhar", Role: "Admin"},
		{ID: 2, Nama: "Budi", Role: "Guru"},
	}
	return users, 2, nil
}

func (m *MockUserService) UpdateDataUser(id uint, req models.UpdateUserRequest) (*models.User, error) {
	return &models.User{ID: id, Nama: req.Nama, Email: req.Email}, nil
}

func (m *MockUserService) HapusDataUser(id uint) error {
	return nil
}

// Fungsi bantuan untuk setup router dengan error middleware
func setupUserTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.GlobalErrorHandler())
	return router
}

// ==========================================
// 2. SKENARIO TES
// ==========================================

func TestRegister_Success(t *testing.T) {
	handler := &UserHandler{Service: &MockUserService{}} // Inject Service, bukan Repo
	router := setupUserTestRouter()
	router.POST("/register", handler.Register)

	jsonBody := []byte(`{"nama": "Azhar", "email": "azhar@example.com", "password": "password123"}`)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Azhar")
}

func TestLogin_Success(t *testing.T) {
	handler := &UserHandler{Service: &MockUserService{}}
	router := setupUserTestRouter()
	router.POST("/login", handler.Login)

	jsonBody := []byte(`{"email": "azhar@example.com", "password": "password123"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "mock_access_token")
}

func TestLogin_WrongPassword(t *testing.T) {
	handler := &UserHandler{Service: &MockUserService{}}
	router := setupUserTestRouter()
	router.POST("/login", handler.Login)

	jsonBody := []byte(`{"email": "azhar@example.com", "password": "password_salah"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Harusnya 401 karena Service melempar utils.NewUnauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Email atau password salah")
}

func TestGetAllUsers_Success(t *testing.T) {
	handler := &UserHandler{Service: &MockUserService{}}
	router := setupUserTestRouter()
	router.GET("/users", handler.GetAllUsers)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Azhar")
}

func TestDeleteUser_Success(t *testing.T) {
	handler := &UserHandler{Service: &MockUserService{}}
	router := setupUserTestRouter()
	router.DELETE("/users/:id", handler.DeleteUser)

	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
