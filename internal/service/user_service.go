package service

import (
	"github.com/azharf99/gothub-erp/internal/models"
	"github.com/azharf99/gothub-erp/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo models.UserRepository
}

func NewUserService(repo models.UserRepository) models.UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(req models.RegisterRequest) (*models.User, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, utils.NewInternalError("Gagal memproses password")
	}

	if validationErr := req.ValidateCustomBusinessLogic(); validationErr != nil {
		return nil, utils.NewBadRequest(validationErr.Error())
	}

	newUser := models.User{
		Nama:     req.Nama,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "Siswa",
	}
	return &newUser, nil
}

func (s *userService) LoginUser(req models.LoginRequest) (string, string, error) {
	user, err := s.repo.CariBerdasarkanEmail(req.Email)
	if err != nil || user == nil {
		return "", "", utils.NewUnauthorized("Email atau password salah")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", "", utils.NewUnauthorized("Email atau password salah")
	}
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		return "", "", utils.NewInternalError("Gagal membuat token autentikasi")
	}
	return accessToken, refreshToken, nil
}

func (s *userService) CreateUserFromDashboard(req models.RegisterRequest, currentUserRole string) (*models.User, error) {
	existingUser, _ := s.repo.CariBerdasarkanEmail(req.Email)
	if existingUser != nil {
		return nil, utils.NewBadRequest("Email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, utils.NewInternalError("Gagal memproses password")
	}

	newUser := models.User{
		Nama:     req.Nama,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     currentUserRole,
	}

	if err := s.repo.SimpanUser(&newUser); err != nil {
		return nil, utils.NewInternalError("Gagal menyimpan data pengguna")
	}
	return &newUser, nil
}

func (s *userService) GetSemuaUser(page, limit int) ([]models.User, int64, error) {
	users, totalItems, err := s.repo.AmbilSemuaUser(page, limit)
	if err != nil {
		return nil, 0, utils.NewInternalError("Gagal mengambil data pengguna")
	}
	return users, totalItems, nil
}

func (s *userService) UpdateDataUser(id uint, req models.UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.AmbilUserByID(id)
	if err != nil {
		return nil, utils.NewNotFound("Pengguna tidak ditemukan")
	}

	user.Nama = req.Nama
	user.Email = req.Email
	if req.Role != "" {
		user.Role = req.Role
	}

	if err := s.repo.UpdateUser(user); err != nil {
		return nil, utils.NewInternalError("Gagal memperbarui data pengguna")
	}

	return user, nil
}

func (s *userService) HapusDataUser(id uint) error {
	err := s.repo.HapusUser(id)
	if err != nil {
		return utils.NewInternalError("Gagal menghapus pengguna")
	}
	return nil
}
