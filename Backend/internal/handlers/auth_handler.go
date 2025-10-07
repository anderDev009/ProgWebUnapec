package handlers

import (
	"net/http"

	"petmatch/internal/middleware"
	"petmatch/internal/models"
	"petmatch/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *services.AuthService
}

type registerRequest struct {
	Name        string  `json:"name" binding:"required,min=2"`
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=6"`
	Role        string  `json:"role" binding:"required"`
	ShelterName *string `json:"shelterName"`
	Phone       *string `json:"phone"`
	City        *string `json:"city"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func NewAuthHandler(auth *services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.auth.Register(services.RegisterInput{
		Name:        req.Name,
		Email:       req.Email,
		Password:    req.Password,
		Role:        req.Role,
		ShelterName: req.ShelterName,
		Phone:       req.Phone,
		City:        req.City,
	})
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case services.ErrEmailInUse, services.ErrUnsupportedRole:
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":          user.ID,
			"name":        user.Name,
			"email":       user.Email,
			"role":        user.Role,
			"city":        user.City,
			"phone":       user.Phone,
			"isApproved":  user.IsApproved,
			"shelterName": user.ShelterName,
		},
		"message": messageForRole(user.Role),
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		status := http.StatusUnauthorized
		switch err {
		case services.ErrShelterNotApproved:
			status = http.StatusForbidden
		case services.ErrInvalidCredentials:
			status = http.StatusUnauthorized
		default:
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": result.Token,
		"user": gin.H{
			"id":          result.User.ID,
			"name":        result.User.Name,
			"email":       result.User.Email,
			"role":        result.User.Role,
			"city":        result.User.City,
			"phone":       result.User.Phone,
			"isApproved":  result.User.IsApproved,
			"shelterName": result.User.ShelterName,
		},
	})
}

func messageForRole(role models.UserRole) string {
	if role == models.RoleShelter {
		return "Tu cuenta de refugio será revisada por un administrador."
	}
	return "Registro completado."
}

func CurrentUserHandler(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":          user.ID,
			"name":        user.Name,
			"email":       user.Email,
			"role":        user.Role,
			"city":        user.City,
			"phone":       user.Phone,
			"isApproved":  user.IsApproved,
			"shelterName": user.ShelterName,
		},
	})
}
