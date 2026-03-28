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
// 1. MOCK SERVICE UNTUK COURSE
// ==========================================
type MockCourseService struct{}

func (m *MockCourseService) BuatCourse(req models.CourseRequest, teacherID uint) (*models.Course, error) {
	return &models.Course{ID: 1, Nama: req.Nama, TeacherID: teacherID}, nil
}

func (m *MockCourseService) AmbilSemuaCourse(page, limit int) ([]models.Course, int64, error) {
	courses := []models.Course{{ID: 1, Nama: "Biologi", TeacherID: 1}}
	return courses, 1, nil
}

func (m *MockCourseService) UpdateCourse(courseID uint, req models.CourseRequest, userID uint, userRole string) (*models.Course, error) {
	// Simulasi Logika Bisnis dari Service
	if userRole != "Admin" && userID != 1 { // Anggap pemilik mapel ini adalah User ID 1
		return nil, utils.NewForbidden("Akses ditolak: Anda bukan pengajar mata pelajaran ini")
	}
	return &models.Course{ID: courseID, Nama: req.Nama}, nil
}

func (m *MockCourseService) HapusCourse(courseID uint, userID uint, userRole string) error {
	return nil
}

// Middleware bantuan untuk menyuntikkan UserID, Email, dan Role
func mockAuthMiddleware(userID uint, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("email", "test@example.com")
		c.Set("role", role)
		c.Next()
	}
}

func setupCourseTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.GlobalErrorHandler())
	return router
}

// ==========================================
// 2. SKENARIO TES
// ==========================================

func TestCreateCourse_Success(t *testing.T) {
	handler := &CourseHandler{Service: &MockCourseService{}}
	router := setupCourseTestRouter()

	router.Use(mockAuthMiddleware(1, "Guru")) // Suntikkan sesi
	router.POST("/courses", handler.CreateCourse)

	jsonBody := []byte(`{"nama": "Biologi Kelas XI", "deskripsi": "Materi"}`)
	req, _ := http.NewRequest("POST", "/courses", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Biologi Kelas XI")
}

func TestUpdateCourse_Forbidden(t *testing.T) {
	handler := &CourseHandler{Service: &MockCourseService{}}
	router := setupCourseTestRouter()

	// Skenario: User ID 99 (Bukan pemilik) dan cuma Guru
	router.Use(mockAuthMiddleware(99, "Guru"))
	router.PUT("/courses/:id", handler.UpdateCourse)

	jsonBody := []byte(`{"nama": "Biologi Update", "deskripsi": "Update materi"}`)
	req, _ := http.NewRequest("PUT", "/courses/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Harus ditolak (403 Forbidden)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "bukan pengajar")
}

func TestUpdateCourse_Admin_Success(t *testing.T) {
	handler := &CourseHandler{Service: &MockCourseService{}}
	router := setupCourseTestRouter()

	// Skenario: User ID 99 (Bukan pemilik) TAPI rolenya Admin
	router.Use(mockAuthMiddleware(99, "Admin"))
	router.PUT("/courses/:id", handler.UpdateCourse)

	jsonBody := []byte(`{"nama": "Biologi Update Admin", "deskripsi": "Update materi"}`)
	req, _ := http.NewRequest("PUT", "/courses/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Harusnya diizinkan (200 OK)
	assert.Equal(t, http.StatusOK, w.Code)
}
