package handlers

import (
	"authentication/models"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/labstack/echo"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

// HTTP status codes
const (
	StatusBadRequest          = http.StatusBadRequest
	StatusUnauthorized        = http.StatusUnauthorized
	StatusInternalServerError = http.StatusInternalServerError
	StatusConflict            = http.StatusConflict
	StatusOK                  = http.StatusOK
	StatusCreated             = http.StatusCreated
)

// ErrorResponse defines the structure for error responses
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// AuthHandler struct holds the database connection
type AuthHandler struct {
	DB *gorm.DB
}

// NewAuthHandler is a constructor function for AuthHandler
func NewAuthHandler(DB *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: DB}
}

// BindJSONError sends a JSON error response for binding errors
func BindJSONError(ctx echo.Context, err error) error {
	return ctx.JSON(StatusBadRequest, ErrorResponse{Status: "fail", Message: "Invalid input data"})
}

// DatabaseError sends a JSON error response for database errors
func DatabaseError(ctx echo.Context, err error) error {
	return ctx.JSON(StatusInternalServerError, ErrorResponse{Status: "error", Message: "Database operation failed"})
}

// SignUpUser is a method on AuthHandler to handle user sign up
func (ah *AuthHandler) SignUpUser(ctx echo.Context) error {
	payload := new(models.CreateUserInput)
	if err := ctx.Bind(payload); err != nil {
		return BindJSONError(ctx, err)
	}

	newUser := models.User{
		Name:     payload.Name,
		Email:    strings.ToLower(payload.Email),
		Password: payload.Password,
	}
	if err := ah.DB.Create(&newUser).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique") {
			return ctx.JSON(StatusConflict, ErrorResponse{Status: "fail", Message: "Email already exists, please use another email address"})
		}
		return DatabaseError(ctx, err)
	}

	return ctx.JSON(StatusCreated, ErrorResponse{Status: "success", Message: "Registered successfully, please login"})
}

// LoginUser is a method on AuthHandler to handle user login
func (ah *AuthHandler) LoginUser(ctx echo.Context) error {
	payload := new(models.LoginUserInput)
	if err := ctx.Bind(payload); err != nil {
		return BindJSONError(ctx, err)
	}

	var user models.User
	_, err := mail.ParseAddress(payload.Email)
	if err != nil {
		return ctx.JSON(StatusBadRequest, ErrorResponse{Status: "fail", Message: "Invalid email address"})
	}

	if err := ah.DB.First(&user, "email = ?", strings.ToLower(payload.Email)).Error; err != nil {
		return ctx.JSON(StatusUnauthorized, ErrorResponse{Status: "fail", Message: "Invalid email or password"})
	}

	userResponse := map[string]interface{}{
		"id":          fmt.Sprintf("%v", user.ID),
		"name":        user.Name,
		"email":       user.Email,
		"otp_enabled": user.Otp_enabled,
	}

	return ctx.JSON(StatusOK, map[string]interface{}{"status": "success", "user": userResponse})
}

// GenerateOTP is a method that generates OTP
func (ah *AuthHandler) GenerateOTP(ctx echo.Context) error {
	payload := new(models.OTPInput)
	if err := ctx.Bind(payload); err != nil {
		return BindJSONError(ctx, err)
	}

	var user models.User
	if err := ah.DB.First(&user, "id = ?", payload.UserId).Error; err != nil {
		return ctx.JSON(StatusBadRequest, ErrorResponse{Status: "fail", Message: "Invalid user"})
	}

	if user.Otp_enabled {
		return ctx.JSON(StatusBadRequest, ErrorResponse{Status: "fail", Message: "OTP is already enabled for the user"})
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      user.Name,
		AccountName: user.Email,
		SecretSize:  15,
	})
	if err != nil {
		return ctx.JSON(StatusInternalServerError, ErrorResponse{Status: "fail", Message: "Failed to generate OTP"})
	}

	dataToUpdate := models.User{
		Otp_secret:   key.Secret(),
		Otp_auth_url: key.URL(),
	}
	if err := ah.DB.Model(&user).Updates(dataToUpdate).Error; err != nil {
		return DatabaseError(ctx, err)
	}

	otpResponse := map[string]interface{}{
		"secret_key": key.Secret(),
		"otp_url":    key.URL(),
	}
	return ctx.JSON(StatusOK, otpResponse)
}

// VerifyOTP is a method on AuthHandler to verify the OTP provided by the user
func (ah *AuthHandler) VerifyOTP(ctx echo.Context) error {
	payload := new(models.OTPInput)
	if err := ctx.Bind(payload); err != nil {
		return BindJSONError(ctx, err)
	}

	var user models.User
	if err := ah.DB.First(&user, "id = ?", payload.UserId).Error; err != nil {
		return ctx.JSON(StatusBadRequest, ErrorResponse{Status: "fail", Message: "Invalid user"})
	}

	valid := totp.Validate(payload.Otp, user.Otp_secret)
	if !valid {
		return ctx.JSON(StatusBadRequest, ErrorResponse{Status: "fail", Message: "Invalid OTP or user doesn't exist"})
	}

	dataToUpdate := models.User{
		Otp_enabled:  true,
		Otp_verified: true,
	}
	if err := ah.DB.Model(&user).Updates(dataToUpdate).Error; err != nil {
		return DatabaseError(ctx, err)
	}

	userResponse := map[string]interface{}{
		"id":          user.ID.String(),
		"name":        user.Name,
		"email":       user.Email,
		"otp_enabled": user.Otp_enabled,
	}
	return ctx.JSON(StatusOK, map[string]interface{}{"otp_verified": true, "user": userResponse})
}

// ValidateOTP is a method on AuthHandler to validate the OTP provided by the user
func (ah *AuthHandler) ValidateOTP(ctx echo.Context) error {
	payload := new(models.OTPInput)
	if err := ctx.Bind(payload); err != nil {
		return BindJSONError(ctx, err)
	}

	user := new(models.User)
	if err := ah.DB.Where("id = ?", payload.UserId).First(user).Error; err != nil {
		return ctx.JSON(StatusUnauthorized, ErrorResponse{Status: "fail", Message: "Invalid user"})
	}

	valid := totp.Validate(payload.Otp, user.Otp_secret)
	if !valid {
		return ctx.JSON(StatusUnauthorized, ErrorResponse{Status: "fail", Message: "Invalid OTP"})
	}

	return ctx.JSON(StatusOK, map[string]interface{}{"otp_valid": true})
}

// DisableOTP is a method on AuthHandler to disable the OTP for a user
func (ah *AuthHandler) DisableOTP(ctx echo.Context) error {
	payload := new(models.OTPInput)
	if err := ctx.Bind(payload); err != nil {
		return BindJSONError(ctx, err)
	}

	var user models.User
	if err := ah.DB.First(&user, "id = ?", payload.UserId).Error; err != nil {
		return ctx.JSON(StatusBadRequest, ErrorResponse{Status: "fail", Message: "Invalid user"})
	}

	user.Otp_enabled = false
	if err := ah.DB.Save(&user).Error; err != nil {
		return DatabaseError(ctx, err)
	}

	userResponse := map[string]interface{}{
		"id":          user.ID.String(),
		"name":        user.Name,
		"email":       user.Email,
		"otp_enabled": user.Otp_enabled,
	}
	return ctx.JSON(StatusOK, map[string]interface{}{"otp_disabled": true, "user": userResponse})
}
