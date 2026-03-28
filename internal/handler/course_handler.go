package handler

import (
	"math"
	"net/http"

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

	userID, _, _, err := utils.GetCurrentUser(c)
	if err != nil {
		c.Error(err)
		return
	}

	newCourse := models.Course{
		Nama:      req.Nama,
		Deskripsi: req.Deskripsi,
		TeacherID: userID,
	}

	if err := h.Repo.BuatCourse(&newCourse); err != nil {
		c.Error(utils.NewInternalError("Gagal menyimpan mata pelajaran ke database"))
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Mata pelajaran berhasil ditambahkan", newCourse)
}

// GET ALL COURSES
func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)

	courses, totalItems, err := h.Repo.AmbilSemuaCourse(page, limit)
	if err != nil {
		c.Error(utils.NewInternalError("Gagal mengambil data mata pelajaran"))
		return
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	meta := utils.PaginationMeta{
		CurrentPage: page,
		Limit:       limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}

	utils.SendPaginatedSuccess(c, http.StatusOK, "Berhasil mengambil daftar mata pelajaran", courses, meta)
}

// UPDATE COURSE
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	courseID, err := utils.GetParamID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req models.CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewBadRequest("Format JSON tidak sesuai atau data tidak lengkap"))
		return
	}

	course, err := h.Repo.AmbilCourseByID(courseID)
	if err != nil {
		c.Error(utils.NewNotFound("Mata pelajaran tidak ditemukan"))
		return
	}

	userID, _, userRole, err := utils.GetCurrentUser(c)
	if err != nil {
		c.Error(err)
		return
	}

	if userRole != "Admin" && course.TeacherID != userID {
		c.Error(utils.NewForbidden("Akses ditolak: Anda bukan pengajar mata pelajaran ini"))
		return
	}

	course.Nama = req.Nama
	course.Deskripsi = req.Deskripsi

	if err := h.Repo.UpdateCourse(course); err != nil {
		c.Error(utils.NewInternalError("Gagal memperbarui mata pelajaran"))
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Mata pelajaran berhasil diperbarui", course)
}

// DELETE COURSE (Versi Standar Industri yang Sangat Bersih)
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	courseID, err := utils.GetParamID(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	userID, _, userRole, err := utils.GetCurrentUser(c)
	if err != nil {
		c.Error(err)
		return
	}

	course, err := h.Repo.AmbilCourseByID(courseID)
	if err != nil {
		c.Error(utils.NewNotFound("Mata pelajaran tidak ditemukan"))
		return
	}

	if userRole != "Admin" && course.TeacherID != userID {
		c.Error(utils.NewForbidden("Akses ditolak: Anda bukan pengajar mata pelajaran ini"))
		return
	}

	if err := h.Repo.HapusCourse(courseID); err != nil {
		c.Error(utils.NewInternalError("Gagal menghapus mata pelajaran"))
		return
	}

	utils.SendSuccess(c, 200, "Mata pelajaran berhasil dihapus", nil)
}
