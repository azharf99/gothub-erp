package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azharf99/gothub-erp/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 1. KITA BUAT REPOSITORY PALSU (MOCK)
type MockUserRepository struct{}

func (m *MockUserRepository) SimpanUser(user *models.User) error {
	// Pura-pura sukses menyimpan dan langsung memberikan ID
	user.ID = 1
	return nil
}

func (m *MockUserRepository) CariBerdasarkanEmail(email string) (*models.User, error) {
	return nil, nil // Dilewati dulu untuk tes register
}

// 2. FUNGSI TESTING UTAMA
func TestRegister_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Suntikkan Database Palsu ke Handler
	mockRepo := &MockUserRepository{}
	handler := &UserHandler{Repo: mockRepo}

	router := gin.Default()
	router.POST("/register", handler.Register)

	// Siapkan data JSON palsu
	jsonBody := []byte(`{"nama": "Azhar", "email": "azhar@example.com", "password": "password123"}`)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Jalankan request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validasi hasilnya menggunakan assert
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Registrasi berhasil")
}
