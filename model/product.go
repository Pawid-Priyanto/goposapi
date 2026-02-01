package model

type Product struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Price         int    `json:"price"`
	Stock         int    `json:"stock"`
	Category_ID   int    `json:"category_id"`
	Category_Name string `json:"category_name"`
}
