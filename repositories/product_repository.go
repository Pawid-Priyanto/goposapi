package repositories

import (
	"api-pos/model"
	"errors"

	"database/sql"

	"github.com/jackc/pgx/v5"
)

type ProductRepository struct {
	// db *pgxpool.Pool
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll - Ambil semua produk dari database
func (r *ProductRepository) GetAll() ([]model.Product, error) {
	query := `SELECT p.id, p.name, p.price, p.stock, p.category_id, COALESCE(c.name, '') AS category_name FROM products p LEFT JOIN categories c ON p.category_id = c.id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Category_ID, &product.Category_Name)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// GetByID - Ambil produk by ID
func (r *ProductRepository) GetByID(id int) (*model.Product, error) {
	query := `SELECT p.id, p.name, p.price, p.stock, p.category_id,COALESCE(c.name, 'No Category') as category_name FROM products p LEFT JOIN categories c ON p.category_id = c.id WHERE p.id = $1`

	var product model.Product
	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.Category_ID,
		&product.Category_Name,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

// Create - Tambah produk baru
func (r *ProductRepository) Create(product model.Product) (*model.Product, error) {
	query := `
        INSERT INTO products (name, price, stock, category_id) 
        VALUES ($1, $2, $3,$4) 
        RETURNING id, (SELECT name FROM categories WHERE id = $4)
    `

	err := r.db.QueryRow(query, product.Name, product.Price, product.Stock, product.Category_ID).Scan(&product.ID, &product.Category_Name)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// Update - Update produk by ID
func (r *ProductRepository) Update(id int, product model.Product) (*model.Product, error) {
	query := `
        UPDATE products 
        SET name = $1, price = $2, stock = $3 , category_id = $4
        WHERE id = $5 
		RETURNING id, name, price, stock, category_id, (SELECT name from categories WHERE id = $4)
    `

	err := r.db.QueryRow(query, product.Name, product.Price, product.Stock, product.Category_ID, id).Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Category_ID, &product.Category_Name)
	if err != nil {
		return nil, err
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	product.ID = id
	return &product, nil
}

// Delete - Hapus produk by ID
func (r *ProductRepository) Delete(id int) error {
	query := `DELETE FROM products WHERE id = $1`

	commandTag, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := commandTag.RowsAffected()

	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}
