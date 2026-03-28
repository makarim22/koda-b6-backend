package models

type Variant struct {
	ID              int     `json:"id" db:"id"`
	Name            string  `json:"name" db:"name"`
	AdditionalPrice float64 `json:"additional_price" db:"additional_price"`
}
