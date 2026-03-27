package utils

// AppError adalah struktur error standar untuk Gothub ERP
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Mengimplementasikan interface error bawaan Golang
func (e *AppError) Error() string {
	return e.Message
}

// Fungsi bantuan untuk mempercepat pembuatan error yang umum
func NewBadRequest(msg string) *AppError {
	return &AppError{Code: 400, Message: msg}
}

func NewUnauthorized(msg string) *AppError {
	return &AppError{Code: 401, Message: msg}
}

func NewForbidden(msg string) *AppError {
	return &AppError{Code: 403, Message: msg}
}

func NewNotFound(msg string) *AppError {
	return &AppError{Code: 404, Message: msg}
}

func NewInternalError(msg string) *AppError {
	return &AppError{Code: 500, Message: msg}
}
