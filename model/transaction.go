package model

import "time"

type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details"`
}

type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name,omitempty"`
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}

type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

type TopProduct struct {
	Name    string `json:"name"`
	QtySold int    `json:"quantity_sold"`
}

type ReportResponse struct {
	TotalRevenue     int        `json:"total_revenue"`
	TotalTransaction int        `json:"total_transaction"`
	TopProduct       TopProduct `json:"top_product"`
}
