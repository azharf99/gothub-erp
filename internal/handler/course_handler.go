package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/azharf99/gothub-erp/internal/models"
	"github.com/azharf99/gothub-erp/internal/utils"
	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	Repo models.CourseRepository
}

// CREATE COURSE
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req models.CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewBadRequest("Format JSON tidak sesuai atau data tidak lengkap"))
		return
	}

	// MENGAMBIL ID DARI TOKEN JWT (Diset oleh middleware RequireAuth)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.Error(utils.NewUnauthorized("Sesi Anda tidak valid, silakan login ulang"))
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
		c.Error(utils.NewInternalError("Gagal menyimpan mata pelajaran ke database"))
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Mata pelajaran berhasil ditambahkan", newCourse)
}

// GET ALL COURSES
func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	// 1. Tangkap parameter dari URL, beri nilai default jika tidak diisi
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// 2. Konversi dari string ke integer
	page, errPage := strconv.Atoi(pageStr)
	limit, errLimit := strconv.Atoi(limitStr)

	// Validasi dasar agar user tidak memasukkan angka minus
	if errPage != nil || page < 1 {
		page = 1
	}
	if errLimit != nil || limit < 1 {
		limit = 10
	}

	// 3. Panggil Repository
	courses, totalItems, err := h.Repo.AmbilSemuaCourse(page, limit)
	if err != nil {
		c.Error(utils.NewInternalError("Gagal mengambil data mata pelajaran"))
		return
	}

	// 4. Hitung Total Halaman (Membulatkan ke atas, misal 11 data / 10 limit = 2 halaman)
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	// 5. Susun Metadata
	meta := utils.PaginationMeta{
		CurrentPage: page,
		Limit:       limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}

	// 6. Kirim menggunakan Amplop Paginated yang baru kita buat!
	utils.SendPaginatedSuccess(c, http.StatusOK, "Berhasil mengambil daftar mata pelajaran", courses, meta)
}

// UPDATE COURSE
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	idParam := c.Param("id")
	courseID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.Error(utils.NewBadRequest("ID mata pelajaran tidak valid"))
		return
	}

	var req models.CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewBadRequest("Format JSON tidak sesuai atau data tidak lengkap"))
		return
	}

	// Cari mapel di database
	course, err := h.Repo.AmbilCourseByID(uint(courseID))
	if err != nil {
		c.Error(utils.NewNotFound("Mata pelajaran tidak ditemukan"))
		return
	}

	// ==========================================
	// OTORISASI: Cek apakah user berhak mengubah?
	// ==========================================
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")

	// Jika dia BUKAN Admin, dan dia BUKAN guru pembuat mapel ini -> Tolak!
	if userRole != "Admin" && course.TeacherID != userID.(uint) {
		c.Error(utils.NewForbidden("Akses ditolak: Anda bukan pengajar mata pelajaran ini"))
		return
	}

	// Update data
	course.Nama = req.Nama
	course.Deskripsi = req.Deskripsi

	if err := h.Repo.UpdateCourse(course); err != nil {
		c.Error(utils.NewInternalError("Gagal memperbarui mata pelajaran"))
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Mata pelajaran berhasil diperbarui", course)
}

// DELETE COURSE
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	idParam := c.Param("id")
	courseID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.Error(utils.NewBadRequest("ID mata pelajaran tidak valid"))
		return
	}

	course, err := h.Repo.AmbilCourseByID(uint(courseID))
	if err != nil {
		c.Error(utils.NewNotFound("Mata pelajaran tidak ditemukan"))
		return
	}

	// ==========================================
	// OTORISASI KEPEMILIKAN
	// ==========================================
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")

	if userRole != "Admin" && course.TeacherID != userID.(uint) {
		c.Error(utils.NewForbidden("Akses ditolak: Anda hanya dapat menghapus mata pelajaran yang Anda buat"))
		return
	}

	if err := h.Repo.HapusCourse(uint(courseID)); err != nil {
		c.Error(utils.NewInternalError("Gagal menghapus mata pelajaran"))
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Mata pelajaran berhasil dihapus", nil)
}
