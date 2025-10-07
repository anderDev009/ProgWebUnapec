package models

import "time"

type PetStatus string

const (
	PetStatusAvailable PetStatus = "available"
	PetStatusAdopted   PetStatus = "adopted"
)

type Pet struct {
	ID          uint      `gorm:"primaryKey"`
	ShelterID   uint      `gorm:"not null"`
	Shelter     User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name        string    `gorm:"size:120;not null"`
	Species     string    `gorm:"size:80;not null"`
	Breed       string    `gorm:"size:120"`
	Age         uint      `gorm:"not null"`
	Description string    `gorm:"type:text"`
	Location    string    `gorm:"size:120"`
	PhotoURL    *string   `gorm:"size:255"`
	Status      PetStatus `gorm:"size:20;default:'available'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
