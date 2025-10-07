package models

import "time"

type UserRole string

const (
	RoleAdopter UserRole = "adopter"
	RoleShelter UserRole = "shelter"
	RoleAdmin   UserRole = "admin"
)

type User struct {
	ID           uint     `gorm:"primaryKey"`
	Name         string   `gorm:"size:120;not null"`
	Email        string   `gorm:"size:180;not null;uniqueIndex"`
	PasswordHash string   `gorm:"size:255;not null"`
	Role         UserRole `gorm:"size:20;not null"`
	ShelterName  *string  `gorm:"size:150"`
	Phone        *string  `gorm:"size:40"`
	City         *string  `gorm:"size:80"`
	IsApproved   bool     `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Pets         []Pet `gorm:"foreignKey:ShelterID"`
}
