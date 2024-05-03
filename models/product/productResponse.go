package product

import (
	"basic-trade-api/models/variant"
	"mime/multipart"
	"time"
)

type ProductResponse struct {
	ID              int                   `json:"id"`
	UUID            string                `json:"uuid"`
	Name            string                `json:"name"`
	ImageURL        string                `json:"imageUrl"`
	ImageFileHeader *multipart.FileHeader `json:"-"`
	AdminID         int                   `json:"adminId"`
	Variants        []variant.VariantResponse
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
