package admin

import (
	"github.com/go-playground/validator/v10"
)

type AdminRegisterRequest struct {
	Name     string `json:"name" binding:"required,min=3,max=100" validate:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100" validate:"required,min=6,max=100"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email,min=3,max=100"`
	Password string `json:"password" binding:"required,min=3,max=100" validate:"required,min=3,max=100"`
}

var Validate = validator.New()