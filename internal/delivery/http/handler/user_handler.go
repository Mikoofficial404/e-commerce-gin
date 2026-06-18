package handler

import (
	"ecommerce-gin/internal/dto/request"
	"ecommerce-gin/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userServices service.UserService
}

type JsonHandler struct {
	Login    request.LoginRequest
	Register request.RegisterRequest
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{userServices: *svc}
}

func (h *UserHandler) Register(c *gin.Context) {
	var input request.RegisterRequest
	//TODO bind
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//TODO ambil field dari struct
	data, err := h.userServices.Register(input.Email, input.FullName, input.Password, input.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//TODO panggil Register dan return response
	c.JSON(http.StatusCreated, gin.H{"status": "you are succesfull register", "data": data})
	return
}

func (h *UserHandler) Login(c *gin.Context) {
	var input request.LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.userServices.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "you are logged in", "data": data})
	return
}
