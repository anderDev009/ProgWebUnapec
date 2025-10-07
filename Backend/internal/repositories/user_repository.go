package repositories

import (
	"errors"

	"petmatch/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type UserFilter struct {
	Role     *models.UserRole
	Approved *bool
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) List(filter UserFilter) ([]models.User, error) {
	query := r.db.Model(&models.User{})

	if filter.Role != nil {
		query = query.Where("role = ?", *filter.Role)
	}

	if filter.Approved != nil {
		query = query.Where("is_approved = ?", *filter.Approved)
	}

	var users []models.User
	if err := query.Order("created_at desc").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
