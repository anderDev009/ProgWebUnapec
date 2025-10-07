package models

import "time"

type AdoptionStatus string

const (
	AdoptionStatusPending  AdoptionStatus = "pending"
	AdoptionStatusApproved AdoptionStatus = "approved"
	AdoptionStatusRejected AdoptionStatus = "rejected"
)

type AdoptionRequest struct {
	ID        uint           `gorm:"primaryKey"`
	PetID     uint           `gorm:"not null"`
	Pet       Pet            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AdopterID uint           `gorm:"not null"`
	Adopter   User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Message   string         `gorm:"type:text"`
	Status    AdoptionStatus `gorm:"size:20;default:'pending'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
