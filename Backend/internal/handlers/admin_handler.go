package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"petmatch/internal/models"
	"petmatch/internal/repositories"
	"petmatch/internal/services"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	users *repositories.UserRepository
	auth  *services.AuthService
}

func NewAdminHandler(users *repositories.UserRepository, auth *services.AuthService) *AdminHandler {
	return &AdminHandler{
		users: users,
		auth:  auth,
	}
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	var roleFilter *models.UserRole
	if role := c.Query("role"); role != "" {
		r := models.UserRole(strings.ToLower(role))
		roleFilter = &r
	}

	var approvedFilter *bool
	if approved := c.Query("approved"); approved != "" {
		value := strings.EqualFold(approved, "true")
		approvedFilter = &value
	}

	users, err := h.users.List(repositories.UserFilter{
		Role:     roleFilter,
		Approved: approvedFilter,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *AdminHandler) ApproveShelter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := h.users.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if user.Role != models.RoleShelter {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not a shelter"})
		return
	}

	if user.IsApproved {
		c.JSON(http.StatusOK, gin.H{"message": "shelter already approved"})
		return
	}

	if err := h.auth.ApproveShelter(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "shelter approved",
		"user":    user,
	})
}
