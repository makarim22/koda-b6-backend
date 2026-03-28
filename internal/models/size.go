package models

type Size struct {
	ID              int     `json:"id" db:"id"`
	Name            string  `json:"name" db:"name"`
	AdditionalPrice float64 `json:"additional_price" db:"additional_price"`
}

type ProductSize struct {
	ProductID int `json:"product_id" db:"product_id"`
	SizeID    int `json:"size_id" db:"size_id"`
}
