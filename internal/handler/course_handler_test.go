package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azharf99/gothub-erp/internal/middleware"
	"github.com/azharf99/gothub-erp/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ==========================================
// 1. MOCK REPOSITORY UNTUK COURSE
// ==========================================
type MockCourseRepository struct{}

func (m *MockCourseRepository) BuatCourse(course *models.Course) error {
	course.ID = 1
	return nil
}

func (m *MockCourseRepository) AmbilSemuaCourse(page, limit int) ([]models.Course, int64, error) {
	return []models.Course{
		{ID: 1, Nama: "Biologi Kelas XI", TeacherID: 1},
	}, int64(1), nil
}

func (m *MockCourseRepository) AmbilCourseByID(id uint) (*models.Course, error) {
	if id == 1 {
		return &models.Course{ID: 1, Nama: "Biologi Kelas XI", TeacherID: 1}, nil
	}
	return nil, errors.New("course tidak ditemukan")
}

func (m *MockCourseRepository) UpdateCourse(course *models.Course) error { return nil }
func (m *MockCourseRepository) HapusCourse(id uint) error                { return nil }

// Middleware bantuan untuk menyuntikkan UserID dan Role ke dalam Gin Context saat testing
func mockAuthMiddleware(userID uint, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.GlobalErrorHandler()) // Semua tes otomatis pakai penangkap error
	return router
}

// ==========================================
// 2. TEST CREATE COURSE
// ==========================================
func TestCreateCourse_Success(t *testing.T) {
	router := setupTestRouter()
	handler := &CourseHandler{Repo: &MockCourseRepository{}}

	// Suntikkan data bahwa yang sedang login adalah User ID 1 dengan Role Guru
	router.Use(mockAuthMiddleware(1, "Guru"))
	router.POST("/courses", handler.CreateCourse)

	jsonBody := []byte(`{"nama": "Biologi Kelas XI", "deskripsi": "Sistem Peredaran Darah Manusia"}`)
	req, _ := http.NewRequest("POST", "/courses", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Biologi Kelas XI")
}

// ==========================================
// 3. TEST UPDATE COURSE - DITOLAK (Bukan Pemilik)
// ==========================================
func TestUpdateCourse_Forbidden(t *testing.T) {
	router := setupTestRouter()
	handler := &CourseHandler{Repo: &MockCourseRepository{}}
	// Skenario: Yang login adalah User ID 99 (Bukan pemilik mapel ID 1) dan rolenya hanya Guru
	router.Use(mockAuthMiddleware(99, "Guru"))
	router.PUT("/courses/:id", handler.UpdateCourse)

	jsonBody := []byte(`{"nama": "Biologi Update", "deskripsi": "Update materi"}`)
	req, _ := http.NewRequest("PUT", "/courses/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Harusnya ditolak dengan status 403 Forbidden
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "bukan pengajar")
}

// ==========================================
// 4. TEST UPDATE COURSE - SUKSES (Oleh Admin)
// ==========================================
func TestUpdateCourse_Admin_Success(t *testing.T) {
	router := setupTestRouter()
	handler := &CourseHandler{Repo: &MockCourseRepository{}}

	// Skenario: Yang login adalah User ID 99 (Bukan pemilik mapel), TAPI rolenya Admin
	router.Use(mockAuthMiddleware(99, "Admin"))
	router.PUT("/courses/:id", handler.UpdateCourse)

	jsonBody := []byte(`{"nama": "Biologi Update Admin", "deskripsi": "Update materi oleh Admin"}`)
	req, _ := http.NewRequest("PUT", "/courses/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Harusnya diizinkan karena Admin memiliki hak istimewa
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "berhasil diperbarui")
}
