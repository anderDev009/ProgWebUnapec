package handlers

import (
	"net/http"
	"strconv"

	"petmatch/internal/middleware"
	"petmatch/internal/models"
	"petmatch/internal/services"

	"github.com/gin-gonic/gin"
)

type AdoptionHandler struct {
    adoptions *services.AdoptionService
}

type createAdoptionRequest struct {
	Message string `json:"message"`
}

type updateAdoptionStatusRequest struct {
	Status models.AdoptionStatus `json:"status" binding:"required"`
}

func NewAdoptionHandler(service *services.AdoptionService) *AdoptionHandler {
	return &AdoptionHandler{adoptions: service}
}

func (h *AdoptionHandler) Create(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	petID, err := strconv.Atoi(c.Param("petId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pet id"})
		return
	}

	var req createAdoptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request, err := h.adoptions.Create(user, services.CreateRequestInput{
		PetID:   uint(petID),
		Message: req.Message,
	})
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case services.ErrAdopterRoleRequired:
			status = http.StatusForbidden
		case services.ErrPetNotFound:
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"request": request})
}

func (h *AdoptionHandler) ListForShelter(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	requests, err := h.adoptions.ListForShelter(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"requests": requests})
}

func (h *AdoptionHandler) ListForAdopter(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	requests, err := h.adoptions.ListForAdopter(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"requests": requests})
}

// List returns adoption requests based on the caller role.
// - Shelter: lists requests for their pets
// - Adopter: lists their submitted requests
func (h *AdoptionHandler) List(c *gin.Context) {
    user := middleware.CurrentUser(c)
    if user == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
        return
    }

    switch user.Role {
    case models.RoleShelter:
        requests, err := h.adoptions.ListForShelter(user.ID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"requests": requests})
        return
    case models.RoleAdopter:
        requests, err := h.adoptions.ListForAdopter(user.ID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"requests": requests})
        return
    default:
        c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
        return
    }
}

func (h *AdoptionHandler) UpdateStatus(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateAdoptionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request, err := h.adoptions.UpdateStatus(user, uint(id), services.UpdateRequestInput{
		Status: req.Status,
	})
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case services.ErrRequestNotFound:
			status = http.StatusNotFound
		case services.ErrPetNotFound:
			status = http.StatusNotFound
		case services.ErrShelterOwnership:
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"request": request})
}
