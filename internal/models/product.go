package models

type Product struct {
	ID int `json:"id" `
	ProductName string`json:"name" db:"product_name`
	Description string `json:"description" db:"description"`
	BasePrice int `json: "price" db:"base_price"`
	Stock int `json:"stock" db:"stock"`
	// VariantId int `json:"variant_id"`
	// SizeId int `json:"size_id"`
}


var products = map[int]Product{}

var nextProductId = 1