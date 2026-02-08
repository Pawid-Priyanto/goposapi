package repositories

import (
	"api-pos/model"
	"fmt"
	"strings"

	"database/sql"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create - Tambah produk baru
func (r *TransactionRepository) Create(items []model.CheckoutItem) (*model.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	totalAmount := 0
	details := make([]model.TransactionDetail, 0)
	// 1. Validasi Produk, Hitung Total, dan Update Stok
	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products where id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}

		if err != nil {
			return nil, err
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("Insufficient stock for product %s (remaining: %d)", productName, stock)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, model.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)

	if err != nil {
		return nil, err
	}

	if len(details) > 0 {
		query := "INSERT INTO transaction_details(transaction_id, product_id, quantity, subtotal) VALUES "
		args := []interface{}{}
		tempData := []string{}

		for i, detail := range details {
			details[i].TransactionID = transactionID
			base := i * 4
			p := fmt.Sprintf("($%d, $%d, $%d, $%d)", base+1, base+2, base+3, base+4)
			tempData = append(tempData, p)
			args = append(args, transactionID, detail.ProductID, detail.Quantity, detail.Subtotal)

		}
		query += strings.Join(tempData, ",")
		_, err := tx.Exec(query, args...)

		if err != nil {
			return nil, fmt.Errorf("batch insert failed: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &model.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

// get report
func (r *TransactionRepository) GetReport(startDate, endDate string) (*model.ReportResponse, error) {
	var report model.ReportResponse

	query := `SELECT COALESCE(SUM(total_amount), 0), COUNT(id) FROM transactions WHERE created_at::date = $1::date OR (created_at::date >= $1::date AND created_at::date <= $2::date)`

	err := r.db.QueryRow(query, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaction)

	if err != nil {
		return nil, err
	}

	// Query 2: Find the Best Selling Product using a JOIN
	queryTopSales := `
		SELECT p.name, SUM(td.quantity) as total_quantity
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE t.created_at::date >= $1 AND t.created_at::date <= $2
		GROUP BY p.name
		ORDER BY total_quantity DESC
		LIMIT 1`

	err = r.db.QueryRow(queryTopSales, startDate, endDate).Scan(&report.TopProduct.Name, &report.TopProduct.QtySold)
	if err == sql.ErrNoRows {
		report.TopProduct.Name = "No sales yet"
		report.TopProduct.QtySold = 0
	} else if err != nil {
		return nil, fmt.Errorf("error fetching top product: %v", err)
	}

	return &report, nil
}
