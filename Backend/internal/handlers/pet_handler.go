package handlers

import (
	"net/http"
	"strconv"

	"petmatch/internal/middleware"
	"petmatch/internal/models"
	"petmatch/internal/services"

	"github.com/gin-gonic/gin"
)

type PetHandler struct {
	pets *services.PetService
}

type createPetRequest struct {
	Name        string  `json:"name" binding:"required"`
	Species     string  `json:"species" binding:"required"`
	Breed       string  `json:"breed"`
	Age         uint    `json:"age" binding:"required"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	PhotoURL    *string `json:"photoUrl"`
}

type updatePetRequest struct {
	Name        string           `json:"name" binding:"required"`
	Species     string           `json:"species" binding:"required"`
	Breed       string           `json:"breed"`
	Age         uint             `json:"age" binding:"required"`
	Description string           `json:"description"`
	Location    string           `json:"location"`
	PhotoURL    *string          `json:"photoUrl"`
	Status      models.PetStatus `json:"status" binding:"required"`
}

func NewPetHandler(pets *services.PetService) *PetHandler {
	return &PetHandler{pets: pets}
}

func (h *PetHandler) List(c *gin.Context) {
	var status *models.PetStatus
	if raw := c.Query("status"); raw != "" {
		s := models.PetStatus(raw)
		status = &s
	}

	var minAge *uint
	if raw := c.Query("minAge"); raw != "" {
		if value, err := strconv.Atoi(raw); err == nil && value >= 0 {
			v := uint(value)
			minAge = &v
		}
	}

	var maxAge *uint
	if raw := c.Query("maxAge"); raw != "" {
		if value, err := strconv.Atoi(raw); err == nil && value >= 0 {
			v := uint(value)
			maxAge = &v
		}
	}

	pets, err := h.pets.List(services.PetFilterInput{
		Species:  c.Query("species"),
		Breed:    c.Query("breed"),
		Location: c.Query("location"),
		Status:   status,
		MinAge:   minAge,
		MaxAge:   maxAge,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pets": pets})
}

func (h *PetHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	pet, err := h.pets.GetByID(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrPetNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pet": pet})
}

func (h *PetHandler) Create(c *gin.Context) {
	var req createPetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := middleware.CurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	pet, err := h.pets.Create(user, services.CreatePetInput{
		Name:        req.Name,
		Species:     req.Species,
		Breed:       req.Breed,
		Age:         req.Age,
		Description: req.Description,
		Location:    req.Location,
		PhotoURL:    req.PhotoURL,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrShelterRoleRequired {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"pet": pet})
}

func (h *PetHandler) Update(c *gin.Context) {
	var req updatePetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	pet, err := h.pets.Update(user, uint(id), services.UpdatePetInput{
		Name:        req.Name,
		Species:     req.Species,
		Breed:       req.Breed,
		Age:         req.Age,
		Description: req.Description,
		Location:    req.Location,
		PhotoURL:    req.PhotoURL,
		Status:      req.Status,
	})
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case services.ErrShelterRoleRequired:
			status = http.StatusForbidden
		case services.ErrPetNotFound:
			status = http.StatusNotFound
		case services.ErrUnauthorizedPetAccess:
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pet": pet})
}

func (h *PetHandler) Delete(c *gin.Context) {
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

	err = h.pets.Delete(user, uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case services.ErrShelterRoleRequired:
			status = http.StatusForbidden
		case services.ErrPetNotFound:
			status = http.StatusNotFound
		case services.ErrUnauthorizedPetAccess:
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
