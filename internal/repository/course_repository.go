package repository

import (
	"github.com/azharf99/gothub-erp/internal/models"
	"gorm.io/gorm"
)

type CourseRepo struct {
	DB *gorm.DB
}

func NewCourseRepo(db *gorm.DB) models.CourseRepository {
	return &CourseRepo{DB: db}
}

// CREATE: Menyimpan course baru dan langsung memuat data Guru-nya
func (r *CourseRepo) BuatCourse(course *models.Course) error {
	// Lakukan insert ke database
	if err := r.DB.Create(course).Error; err != nil {
		return err
	}
	// LANGSUNG MUAT ULANG relasi Teacher setelah insert sukses!
	return r.DB.Preload("Teacher", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Nama", "Email", "Role")
	}).First(course, course.ID).Error
}

// READ: Mengambil semua course beserta data Gurunya
func (r *CourseRepo) AmbilSemuaCourse() ([]models.Course, error) {
	var courses []models.Course
	// Preload("Teacher") menyuruh GORM menarik data User yang terkait
	// Omit("Teacher.Password") memastikan password guru tidak ikut bocor ke response JSON
	err := r.DB.Preload("Teacher", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Nama", "Email", "Role") // Hanya ambil kolom yang aman
	}).Find(&courses).Error

	return courses, err
}

// Find Course by ID, termasuk data Guru-nya
func (r *CourseRepo) AmbilCourseByID(id uint) (*models.Course, error) {
	var course models.Course
	err := r.DB.Preload("Teacher", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Nama", "Email", "Role")
	}).First(&course, id).Error

	if err != nil {
		return nil, err
	}
	return &course, nil
}

// Update Course
func (r *CourseRepo) UpdateCourse(course *models.Course) error {
	return r.DB.Save(course).Error
}

// Delete Course
func (r *CourseRepo) HapusCourse(id uint) error {
	return r.DB.Delete(&models.Course{}, id).Error
}
