package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetParamID mengambil ID dari URL parameter dan mengembalikannya sebagai uint.
// Jika gagal, otomatis melempar AppError.
func GetParamID(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, NewBadRequest("Parameter ID tidak valid")
	}
	return uint(id), nil
}

// GetCurrentUser mengekstrak UserID, Email, dan Role dari token JWT yang diset Middleware.
func GetCurrentUser(c *gin.Context) (userID uint, email string, role string, err error) {
	idVal, idExists := c.Get("userID")
	emailVal, emailExists := c.Get("email")
	roleVal, roleExists := c.Get("role")

	if !idExists || !roleExists || !emailExists {
		return 0, "", "", NewUnauthorized("Sesi tidak valid, silakan login ulang")
	}

	return idVal.(uint), emailVal.(string), roleVal.(string), nil
}

// GetPaginationParams mengekstrak page dan limit dari URL, mengembalikan nilai default yang aman.
func GetPaginationParams(c *gin.Context) (page int, limit int) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	p, errPage := strconv.Atoi(pageStr)
	l, errLimit := strconv.Atoi(limitStr)

	if errPage != nil || p < 1 {
		p = 1
	}
	if errLimit != nil || l < 1 {
		l = 10
	}

	return p, l
}
