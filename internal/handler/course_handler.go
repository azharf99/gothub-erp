package handler

import (
	"net/http"

	"github.com/azharf99/gothub-erp/internal/models"
	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	Repo models.CourseRepository
}

// CREATE COURSE
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req models.CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap"})
		return
	}

	// MENGAMBIL ID DARI TOKEN JWT (Diset oleh middleware RequireAuth)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}

	// Konversi tipe data sesuai kebutuhan
	teacherID := userIDStr.(uint)

	// Bentuk model untuk disimpan
	newCourse := models.Course{
		Nama:      req.Nama,
		Deskripsi: req.Deskripsi,
		TeacherID: teacherID, // Otomatis mengikat mata pelajaran ke guru yang sedang login!
	}

	if err := h.Repo.BuatCourse(&newCourse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan mata pelajaran"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Mata pelajaran berhasil ditambahkan",
		"data":    newCourse,
	})
}

// GET ALL COURSES
func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	courses, err := h.Repo.AmbilSemuaCourse()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data mata pelajaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil",
		"data":    courses,
	})
}
