package variant

import "time"

type VariantResponse struct {
	ID          int       `json:"id"`
	UUID        string    `json:"uuid"`
	VariantName string    `json:"variantName"`
	Quantity    int       `json:"quantity"`
	ProductID   int       `json:"productId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
