package handler

import (
	"net/http"
	"strconv"

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

// UPDATE COURSE
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	idParam := c.Param("id")
	courseID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID mata pelajaran tidak valid"})
		return
	}

	var req models.CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap"})
		return
	}

	// Cari mapel di database
	course, err := h.Repo.AmbilCourseByID(uint(courseID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mata pelajaran tidak ditemukan"})
		return
	}

	// ==========================================
	// OTORISASI: Cek apakah user berhak mengubah?
	// ==========================================
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")

	// Jika dia BUKAN Admin, dan dia BUKAN guru pembuat mapel ini -> Tolak!
	if userRole != "Admin" && course.TeacherID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: Anda bukan pengajar mata pelajaran ini"})
		return
	}

	// Update data
	course.Nama = req.Nama
	course.Deskripsi = req.Deskripsi

	if err := h.Repo.UpdateCourse(course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui mata pelajaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mata pelajaran berhasil diperbarui",
		"data":    course,
	})
}

// DELETE COURSE
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	idParam := c.Param("id")
	courseID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID mata pelajaran tidak valid"})
		return
	}

	course, err := h.Repo.AmbilCourseByID(uint(courseID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mata pelajaran tidak ditemukan"})
		return
	}

	// ==========================================
	// OTORISASI KEPEMILIKAN
	// ==========================================
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")

	if userRole != "Admin" && course.TeacherID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: Anda hanya dapat menghapus mata pelajaran yang Anda buat"})
		return
	}

	if err := h.Repo.HapusCourse(uint(courseID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus mata pelajaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mata pelajaran berhasil dihapus"})
}
