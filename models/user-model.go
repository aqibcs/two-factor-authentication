package models

import (
	"github.com/go-playground/validator"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Email        string    `gorm:"uniqueIndex;not null" validate:"email"`
	Password     string    `gorm:"not null"`
	Otp_enabled  bool      `gorm:"default:false"`
	Otp_verified bool      `gorm:"default:false"`
	Otp_secret   string
	Otp_auth_url string
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (user *User) BeforeCreate(*gorm.DB) error {
	user.ID = uuid.NewV4()

	return nil
}

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" bindinig:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginUserInput struct {
	Email    string `json:"email" bindinig:"required"`
	Password string `json:"password" binding:"required"`
}

type OTPInput struct {
	UserId string `json:"user_id"`
	Otp    string `json:"otp"`
}
