package models

type Product struct {
	ID int
	Name string
	Description string
	Stock int
	VariantId int `json:"variant_id"`
	SizeId int `json:"size_id"`
}


var products = map[int]Product{}

var nextProductId = 1