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

// CREATE: Menyimpan course baru
func (r *CourseRepo) BuatCourse(course *models.Course) error {
	return r.DB.Create(course).Error
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
