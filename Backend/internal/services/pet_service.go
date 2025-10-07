package services

import (
	"errors"

	"petmatch/internal/models"
	"petmatch/internal/repositories"
)

var (
	ErrUnauthorizedPetAccess = errors.New("pet does not belong to shelter")
	ErrPetNotFound           = errors.New("pet not found")
	ErrShelterRoleRequired   = errors.New("only shelters can manage pets")
)

type PetService struct {
	pets *repositories.PetRepository
}

type PetFilterInput struct {
	Species   string
	Breed     string
	Location  string
	Status    *models.PetStatus
	ShelterID *uint
	MinAge    *uint
	MaxAge    *uint
}

type CreatePetInput struct {
	Name        string
	Species     string
	Breed       string
	Age         uint
	Description string
	Location    string
	PhotoURL    *string
}

type UpdatePetInput struct {
	Name        string
	Species     string
	Breed       string
	Age         uint
	Description string
	Location    string
	PhotoURL    *string
	Status      models.PetStatus
}

func NewPetService(repo *repositories.PetRepository) *PetService {
	return &PetService{pets: repo}
}

func (s *PetService) List(filter PetFilterInput) ([]models.Pet, error) {
	return s.pets.List(repositories.PetFilter{
		Species:   filter.Species,
		Breed:     filter.Breed,
		Location:  filter.Location,
		Status:    filter.Status,
		ShelterID: filter.ShelterID,
		MinAge:    filter.MinAge,
		MaxAge:    filter.MaxAge,
	})
}

func (s *PetService) GetByID(id uint) (*models.Pet, error) {
	pet, err := s.pets.FindByID(id)
	if err != nil {
		return nil, err
	}
	if pet == nil {
		return nil, ErrPetNotFound
	}
	return pet, nil
}

func (s *PetService) Create(owner *models.User, input CreatePetInput) (*models.Pet, error) {
	if owner.Role != models.RoleShelter {
		return nil, ErrShelterRoleRequired
	}

	pet := &models.Pet{
		ShelterID:   owner.ID,
		Name:        input.Name,
		Species:     input.Species,
		Breed:       input.Breed,
		Age:         input.Age,
		Description: input.Description,
		Location:    input.Location,
		PhotoURL:    input.PhotoURL,
		Status:      models.PetStatusAvailable,
	}

	if err := s.pets.Create(pet); err != nil {
		return nil, err
	}

	return pet, nil
}

func (s *PetService) Update(owner *models.User, id uint, input UpdatePetInput) (*models.Pet, error) {
	if owner.Role != models.RoleShelter {
		return nil, ErrShelterRoleRequired
	}

	pet, err := s.pets.FindByID(id)
	if err != nil {
		return nil, err
	}
	if pet == nil {
		return nil, ErrPetNotFound
	}

	if pet.ShelterID != owner.ID {
		return nil, ErrUnauthorizedPetAccess
	}

	pet.Name = input.Name
	pet.Species = input.Species
	pet.Breed = input.Breed
	pet.Age = input.Age
	pet.Description = input.Description
	pet.Location = input.Location
	pet.PhotoURL = input.PhotoURL
	pet.Status = input.Status

	if err := s.pets.Update(pet); err != nil {
		return nil, err
	}

	return pet, nil
}

func (s *PetService) Delete(owner *models.User, id uint) error {
	if owner.Role != models.RoleShelter {
		return ErrShelterRoleRequired
	}

	pet, err := s.pets.FindByID(id)
	if err != nil {
		return err
	}
	if pet == nil {
		return ErrPetNotFound
	}

	if pet.ShelterID != owner.ID {
		return ErrUnauthorizedPetAccess
	}

	return s.pets.Delete(id)
}
