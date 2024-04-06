package middleware

import (
    "basic-trade-api/models/admin"
    "fmt"
    "github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateAdminRegister(adminRequest admin.AdminRegisterRequest) error {
    err := validate.Struct(adminRequest)
    if err != nil {
        validationErrors := make(map[string]string)
        for _, err := range err.(validator.ValidationErrors) {
            validationErrors[err.Field()] = err.Tag()
        }
        return fmt.Errorf("validation errors: %v", validationErrors)
    }
    return nil
}

func ValidateAdminLogin(adminRequest admin.AdminLoginRequest) error {
    err := validate.Struct(adminRequest)
    if err != nil {
        validationErrors := make(map[string]string)
        for _, err := range err.(validator.ValidationErrors) {
            validationErrors[err.Field()] = err.Tag()
        }
        return fmt.Errorf("validation errors: %v", validationErrors)
    }
    return nil
}