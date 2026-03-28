package models

type ProductVariant struct {
	ProductID int `json:"product_id" db:"product_id"`
	VariantID int `json:"variant_id" db:"variant_id"`
}
