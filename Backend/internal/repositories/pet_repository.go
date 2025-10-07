package repositories

import (
	"errors"

	"petmatch/internal/models"

	"gorm.io/gorm"
)

type PetRepository struct {
	db *gorm.DB
}

type PetFilter struct {
	Species   string
	Breed     string
	Location  string
	Status    *models.PetStatus
	ShelterID *uint
	MinAge    *uint
	MaxAge    *uint
}

func NewPetRepository(db *gorm.DB) *PetRepository {
	return &PetRepository{db: db}
}

func (r *PetRepository) Create(pet *models.Pet) error {
	return r.db.Create(pet).Error
}

func (r *PetRepository) Update(pet *models.Pet) error {
	return r.db.Save(pet).Error
}

func (r *PetRepository) Delete(id uint) error {
	return r.db.Delete(&models.Pet{}, id).Error
}

func (r *PetRepository) FindByID(id uint) (*models.Pet, error) {
	var pet models.Pet
	if err := r.db.Preload("Shelter").First(&pet, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pet, nil
}

func (r *PetRepository) List(filter PetFilter) ([]models.Pet, error) {
	query := r.db.Preload("Shelter").Model(&models.Pet{})

	if filter.Species != "" {
		query = query.Where("species = ?", filter.Species)
	}

	if filter.Breed != "" {
		query = query.Where("breed = ?", filter.Breed)
	}

	if filter.Location != "" {
		query = query.Where("location LIKE ?", "%"+filter.Location+"%")
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.ShelterID != nil {
		query = query.Where("shelter_id = ?", *filter.ShelterID)
	}

	if filter.MinAge != nil {
		query = query.Where("age >= ?", *filter.MinAge)
	}

	if filter.MaxAge != nil {
		query = query.Where("age <= ?", *filter.MaxAge)
	}

	var pets []models.Pet
	if err := query.Order("created_at desc").Find(&pets).Error; err != nil {
		return nil, err
	}

	return pets, nil
}
