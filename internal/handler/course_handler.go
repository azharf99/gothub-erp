package handler

import (
	"math"
	"net/http"

	"github.com/azharf99/gothub-erp/internal/models"
	"github.com/azharf99/gothub-erp/internal/utils"
	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	// UBAH: Sekarang Handler bergantung pada Service, bukan lagi Repository
	Service models.CourseService
}

// GET ALL COURSES
func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)

	courses, totalItems, err := h.Service.AmbilSemuaCourse(page, limit)
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

	// Lemparkan ke Service untuk diproses
	newCourse, err := h.Service.BuatCourse(req, userID)
	if err != nil {
		c.Error(err)
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Mata pelajaran berhasil ditambahkan", newCourse)
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
		c.Error(utils.NewBadRequest("Format JSON tidak sesuai"))
		return
	}

	userID, _, userRole, err := utils.GetCurrentUser(c)
	if err != nil {
		c.Error(err)
		return
	}

	// Service yang akan pusing memikirkan apakah user ini berhak edit atau tidak
	updatedCourse, err := h.Service.UpdateCourse(courseID, req, userID, userRole)
	if err != nil {
		c.Error(err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Mata pelajaran berhasil diperbarui", updatedCourse)
}

// DELETE COURSE
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

	// Lemparkan ke Service
	if err := h.Service.HapusCourse(courseID, userID, userRole); err != nil {
		c.Error(err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Mata pelajaran berhasil dihapus", nil)
}
