package models

type Reviews struct {
	Id         int    `json:"id"`
	User_id    int    `json:"user_id"`
	Product_id int    `json:"product_id"`
	Order_id   int    `json:"order_id"`
	Message    string `json:"message"`
	Rating     int    `json:"rating"`
}
