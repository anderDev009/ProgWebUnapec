package services

import (
	"errors"

	"petmatch/internal/models"
	"petmatch/internal/repositories"
)

var (
	ErrAdopterRoleRequired = errors.New("only adopters can submit requests")
	ErrRequestNotFound     = errors.New("adoption request not found")
	ErrShelterOwnership    = errors.New("request does not belong to shelter")
)

type AdoptionService struct {
	adoptions *repositories.AdoptionRepository
	pets      *repositories.PetRepository
}

type CreateRequestInput struct {
	PetID   uint
	Message string
}

type UpdateRequestInput struct {
	Status models.AdoptionStatus
}

func NewAdoptionService(adoptionRepo *repositories.AdoptionRepository, petRepo *repositories.PetRepository) *AdoptionService {
	return &AdoptionService{
		adoptions: adoptionRepo,
		pets:      petRepo,
	}
}

func (s *AdoptionService) Create(adopter *models.User, input CreateRequestInput) (*models.AdoptionRequest, error) {
	if adopter.Role != models.RoleAdopter {
		return nil, ErrAdopterRoleRequired
	}

	pet, err := s.pets.FindByID(input.PetID)
	if err != nil {
		return nil, err
	}
	if pet == nil {
		return nil, ErrPetNotFound
	}

	request := &models.AdoptionRequest{
		PetID:     pet.ID,
		AdopterID: adopter.ID,
		Message:   input.Message,
		Status:    models.AdoptionStatusPending,
	}

	if err := s.adoptions.Create(request); err != nil {
		return nil, err
	}

	return request, nil
}

func (s *AdoptionService) ListForShelter(shelterID uint) ([]models.AdoptionRequest, error) {
	return s.adoptions.ListByShelter(shelterID)
}

func (s *AdoptionService) ListForAdopter(adopterID uint) ([]models.AdoptionRequest, error) {
	return s.adoptions.ListByAdopter(adopterID)
}

func (s *AdoptionService) UpdateStatus(shelter *models.User, requestID uint, input UpdateRequestInput) (*models.AdoptionRequest, error) {
	request, err := s.adoptions.FindByID(requestID)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, ErrRequestNotFound
	}

	pet, err := s.pets.FindByID(request.PetID)
	if err != nil {
		return nil, err
	}
	if pet == nil {
		return nil, ErrPetNotFound
	}

	if pet.ShelterID != shelter.ID {
		return nil, ErrShelterOwnership
	}

	request.Status = input.Status

	if err := s.adoptions.Update(request); err != nil {
		return nil, err
	}

	return request, nil
}
