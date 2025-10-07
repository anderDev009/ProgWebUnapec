package repositories

import (
	"errors"

	"petmatch/internal/models"

	"gorm.io/gorm"
)

type AdoptionRepository struct {
	db *gorm.DB
}

func NewAdoptionRepository(db *gorm.DB) *AdoptionRepository {
	return &AdoptionRepository{db: db}
}

func (r *AdoptionRepository) Create(req *models.AdoptionRequest) error {
	return r.db.Create(req).Error
}

func (r *AdoptionRepository) Update(req *models.AdoptionRequest) error {
	return r.db.Save(req).Error
}

func (r *AdoptionRepository) FindByID(id uint) (*models.AdoptionRequest, error) {
	var request models.AdoptionRequest
	if err := r.db.Preload("Pet").Preload("Adopter").First(&request, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &request, nil
}

func (r *AdoptionRepository) ListByShelter(shelterID uint) ([]models.AdoptionRequest, error) {
	var requests []models.AdoptionRequest
	if err := r.db.
		Joins("JOIN pets ON pets.id = adoption_requests.pet_id").
		Where("pets.shelter_id = ?", shelterID).
		Preload("Pet").
		Preload("Adopter").
		Order("adoption_requests.created_at desc").
		Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *AdoptionRepository) ListByAdopter(adopterID uint) ([]models.AdoptionRequest, error) {
	var requests []models.AdoptionRequest
	if err := r.db.Where("adopter_id = ?", adopterID).
		Preload("Pet").
		Preload("Adopter").
		Order("created_at desc").
		Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}
