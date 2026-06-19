package service

import (
	"ecommerce-gin/internal/domain"
	"ecommerce-gin/internal/pkg/jwt"
	"ecommerce-gin/internal/pkg/rabbitmq"
	"ecommerce-gin/internal/repository/postgres"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserService struct {
	user       postgres.UserRepository
	rabbitConn *amqp.Connection
}

func NewUserService(repo *postgres.UserRepository, rabbitConn *amqp.Connection) *UserService {
	return &UserService{user: *repo, rabbitConn: rabbitConn}
}

func (s *UserService) Register(email string, nama string, password string, phone string) (*domain.User, error) {
	//TODO  Hash password dulu sebelum simpan
	JwtPass, err := jwt.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		Email:       email,
		FullName:    nama,
		Password:    JwtPass,
		PhoneNumber: phone,
	}
	//TODO Lempar ke repository
	isResult, err := s.user.CreateUser(&user)
	if err != nil {
		return nil, err
	}
	rabbitmq.PublishMessage(s.rabbitConn, "email_queue", email)
	return isResult, nil
}

func (s *UserService) Login(email string, pasword string) (string, error) {

	//TODO Cari user menggunakan repo kalau tidak ketemu salah
	isUser, err := s.user.FindByEmail(&domain.User{
		Email: email,
	})
	if err != nil {
		return "", fmt.Errorf("email not registered")
	}
	// TODO Pangill jwt.checkpassword cocokin pakai password dari user
	jwtCheck, err := jwt.CheckPasswordHash(pasword, isUser.Password)
	if err != nil {
		return "", fmt.Errorf("password check failed")
	}
	if !jwtCheck {
		return "", fmt.Errorf("incorrect password")
	}

	parseUUID, err := uuid.Parse(isUser.ID)
	if err != nil {
		return "", fmt.Errorf("failed to parse UUID")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	// TODO buat JWTToken
	jwtMake, err := jwt.MakeJWT(parseUUID, jwtSecret, time.Hour*24)
	if err != nil {
		return "", fmt.Errorf("failed to create JWT token")
	}
	return jwtMake, nil
}
