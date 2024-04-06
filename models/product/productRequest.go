package product

import (
	"mime/multipart"

	"github.com/go-playground/validator/v10"
)

type ProductRequest struct {
	Name      string                `form:"name" binding:"required,min=3,max=100" validate:"required,min=3,max=100"`
	ImageFile *multipart.FileHeader `form:"file"`
	ImageURL  string                `form:"imageUrl"` // Add this field
}

var Validate = validator.New()
