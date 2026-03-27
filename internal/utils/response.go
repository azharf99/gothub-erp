package utils

import (
	"github.com/gin-gonic/gin"
)

// DefaultSuccessResponse adalah struktur standar envelope JSON kita
type DefaultSuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // omitempty: jika kosong (nil), key "data" tidak akan dimunculkan di JSON
}

// SendSuccess adalah fungsi bantuan untuk mengirimkan JSON balasan sukses
// Parameter 'data' menggunakan tipe interface{} agar bisa menerima bentuk data apa saja (Struct, Slice, Map, dll)
func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, DefaultSuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Struktur Metadata untuk Pagination
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

// Amplop khusus untuk respons yang memiliki Pagination
type PaginatedSuccessResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

// Fungsi bantuan untuk mengirim respons dengan Pagination
func SendPaginatedSuccess(c *gin.Context, statusCode int, message string, data interface{}, meta PaginationMeta) {
	c.JSON(statusCode, PaginatedSuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}
