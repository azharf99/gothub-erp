package service

import (
	"github.com/azharf99/gothub-erp/internal/models"
	"github.com/azharf99/gothub-erp/internal/utils"
)

type courseService struct {
	repo models.CourseRepository
}

func NewCourseService(repo models.CourseRepository) models.CourseService {
	return &courseService{repo: repo}
}

func (s *courseService) BuatCourse(req models.CourseRequest, teacherID uint) (*models.Course, error) {
	newCourse := models.Course{
		Nama:      req.Nama,
		Deskripsi: req.Deskripsi,
		TeacherID: teacherID,
	}

	if err := s.repo.BuatCourse(&newCourse); err != nil {
		return nil, utils.NewInternalError("Gagal menyimpan mata pelajaran ke database")
	}
	return &newCourse, nil
}

func (s *courseService) AmbilSemuaCourse(page, limit int) ([]models.Course, int64, error) {
	courses, total, err := s.repo.AmbilSemuaCourse(page, limit)
	if err != nil {
		return nil, 0, utils.NewInternalError("Gagal mengambil data mata pelajaran")
	}
	return courses, total, nil
}

func (s *courseService) UpdateCourse(courseID uint, req models.CourseRequest, userID uint, userRole string) (*models.Course, error) {
	course, err := s.repo.AmbilCourseByID(courseID)
	if err != nil {
		return nil, utils.NewNotFound("Mata pelajaran tidak ditemukan")
	}

	// LOGIKA BISNIS: Otorisasi Kepemilikan
	if userRole != "Admin" && course.TeacherID != userID {
		return nil, utils.NewForbidden("Akses ditolak: Anda bukan pengajar mata pelajaran ini")
	}

	course.Nama = req.Nama
	course.Deskripsi = req.Deskripsi

	if err := s.repo.UpdateCourse(course); err != nil {
		return nil, utils.NewInternalError("Gagal memperbarui mata pelajaran")
	}
	return course, nil
}

func (s *courseService) HapusCourse(courseID uint, userID uint, userRole string) error {
	course, err := s.repo.AmbilCourseByID(courseID)
	if err != nil {
		return utils.NewNotFound("Mata pelajaran tidak ditemukan")
	}

	// LOGIKA BISNIS: Otorisasi Kepemilikan
	if userRole != "Admin" && course.TeacherID != userID {
		return utils.NewForbidden("Akses ditolak: Anda hanya dapat menghapus mata pelajaran yang Anda buat")
	}

	if err := s.repo.HapusCourse(courseID); err != nil {
		return utils.NewInternalError("Gagal menghapus mata pelajaran")
	}
	return nil
}
