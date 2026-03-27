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
