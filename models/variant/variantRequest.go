package variant

import (
	"github.com/go-playground/validator/v10"
)

type VariantRequest struct {
    VariantName string `json:"variantName" binding:"required,min=3,max=100" validate:"required,min=3,max=100"`
    Quantity    int    `json:"quantity" binding:"required" validate:"required,gte=0"`
    ProductID   int    `json:"productId" binding:"required" validate:"required"`
}

var Validate = validator.New()