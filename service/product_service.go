package service

import (
	"api-pos/model"
	"api-pos/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAllProducts() ([]model.Product, error) {
	return s.repo.GetAll()

}
func (s *ProductService) GetProductByID(id int) (*model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) CreateProduct(product model.Product) (*model.Product, error) {

	return s.repo.Create(product)
}

func (s *ProductService) UpdateProduct(id int, product model.Product) (*model.Product, error) {
	return s.repo.Update(id, product)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.Delete(id)
}
