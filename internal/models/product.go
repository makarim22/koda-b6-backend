package models

type Product struct {
	ID          int    `json:"id" `
	ProductName string `json:"product_name" db:"product_name"`
	Description string `json:"description" db:"description"`
	BasePrice   int    `json:"base_price" db:"base_price"`
	Stock       int    `json:"stock" db:"stock"`
	// VariantId int `json:"variant_id"`
	// SizeId int `json:"size_id"`
}

type ProductDetail struct {
	ID          int               `json:"id"`
	ProductName string            `json:"product_name"`
	Description string            `json:"description"`
	Stock       int               `json:"stock"`
	BasePrice   float64           `json:"base_price"`
	Categories  []ProductCategory `json:"categories"`
	Variants    []Variant         `json:"variants"`
	Sizes       []Size            `json:"sizes"`
	Images      []ProductImage    `json:"images"`
	Rating      float64           `json:"rating"`
	ReviewCount int               `json:"review_count"`
}

type ProductCategoryMap struct {
	ID         int `json:"id" db:"id"`
	ProductID  int `json:"product_id" db:"product_id"`
	CategoryID int `json:"category_id" db:"category_id"`
}

type ProductImage struct {
	ID        int    `json:"id" db:"id"`
	ProductID int    `json:"product_id" db:"product_id"`
	Path      string `json:"path" db:"path"`
	IsPrimary bool   `json:"is_primary" db:"is_primary"`
}

type ProductDiscount struct {
	ID           int    `json:"id" db:"id"`
	ProductID    int    `json:"product_id" db:"product_id"`
	DiscountRate string `json:"discount_rate" db:"discount_rate"`
	Description  string `json:"description" db:"description"`
	IsFlashSale  bool   `json:"is_flash_sale" db:"is_flash_sale"`
}

var products = map[int]Product{}

var nextProductId = 1
