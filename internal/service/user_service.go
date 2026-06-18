package service

import (
	"ecommerce-gin/internal/domain"
	"ecommerce-gin/internal/pkg/jwt"
	"ecommerce-gin/internal/repository/postgres"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	user postgres.UserRepository
}

func NewUserService(repo *postgres.UserRepository) *UserService {
	return &UserService{user: *repo}
}

func (s *UserService) Register(email string, nama string, password string) (*domain.User, error) {
	//TODO  Hash password dulu sebelum simpan
	JwtPass, err := jwt.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		Email:    email,
		FullName: nama,
		Password: JwtPass,
	}
	//TODO Lempar ke repository
	isResult, err := s.user.CreateUser(&user)
	if err != nil {
		return nil, err
	}
	return isResult, nil
}

func (s *UserService) Login(email string, pasword string) (string, error) {

	//TODO Cari user menggunakan repo kalau tidak ketemu salah
	isUser, err := s.user.FindByEmail(&domain.User{
		Email: email,
	})
	if err != nil {
		return "", fmt.Errorf("email tidak terdaftar")
	}
	// TODO Pangill jwt.checkpassword cocokin pakai password dari user
	jwtCheck, err := jwt.CheckPasswordHash(pasword, isUser.Password)
	if err != nil {
		return "", fmt.Errorf("Cek Password Salah")
	}
	if !jwtCheck {
		return "", fmt.Errorf("password yang Anda masukkan salah")
	}

	parseUUID, err := uuid.Parse(isUser.ID)
	if err != nil {
		return "", fmt.Errorf("Parse UUID Gagal")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	// TODO buat JWTToken
	jwtMake, err := jwt.MakeJWT(parseUUID, jwtSecret, time.Hour*24)
	if err != nil {
		return "", fmt.Errorf("Pembuatan JWT gagal")
	}
	return jwtMake, nil
}
