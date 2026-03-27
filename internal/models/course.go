package models

import "time"

// Model Course untuk Database
type Course struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Nama      string `gorm:"type:varchar(100);not null" json:"nama"`
	Deskripsi string `gorm:"type:text" json:"deskripsi"`

	// Relasi: Course ini milik siapa?
	TeacherID uint `gorm:"not null" json:"teacher_id"`                    // Ini akan jadi Foreign Key
	Teacher   User `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"` // Data guru akan dimuat ke sini

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Payload untuk menerima request pembuatan Course baru
type CourseRequest struct {
	Nama      string `json:"nama" binding:"required"`
	Deskripsi string `json:"deskripsi" binding:"required"`
}

// Kontrak untuk Course Repository
type CourseRepository interface {
	BuatCourse(course *Course) error
	AmbilSemuaCourse() ([]Course, error)
	// (Untuk CRUD lengkap, kamu bisa tambahkan Update dan Delete nanti)
}
