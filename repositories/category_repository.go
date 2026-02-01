package repositories

import (
	"api-pos/model"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
)

type CategoryRepository struct {
	// db *pgxpool.Pool
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetAll - Ambil semua category dari database
func (r *CategoryRepository) GetAll() ([]model.Category, error) {
	query := `SELECT id, name, description FROM categories ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetByID - Ambil category by ID
func (r *CategoryRepository) GetByID(id int) (*model.Category, error) {
	query := `SELECT id, name, description FROM categories WHERE id = $1`

	var category model.Category
	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return &category, nil
}

// Create - Tambah category baru
func (r *CategoryRepository) Create(category model.Category) (*model.Category, error) {
	query := `
        INSERT INTO categories (name, description) 
        VALUES ($1, $2) 
        RETURNING id
    `

	err := r.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// Update - Update Category by ID
func (r *CategoryRepository) Update(id int, category model.Category) (*model.Category, error) {
	query := `
        UPDATE categories
        SET name = $1, description = $2 
        WHERE id = $3
    `

	commandTag, err := r.db.Exec(query, category.Name, category.Description, id)
	if err != nil {
		return nil, err
	}

	rows, err := commandTag.RowsAffected()

	if rows == 0 {
		return nil, errors.New("category not found")
	}

	category.ID = id
	return &category, nil
}

// Delete - Hapus category by ID
func (r *CategoryRepository) Delete(id int) error {
	query := `DELETE FROM categories WHERE id = $1`

	commandTag, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := commandTag.RowsAffected()
	if rows == 0 {
		return errors.New("category not found")
	}

	return nil
}
