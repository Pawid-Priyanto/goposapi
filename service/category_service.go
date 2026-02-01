package service

import (
	"api-pos/model"
	"api-pos/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAllCategories() ([]model.Category, error) {
	return s.repo.GetAll()

}
func (s *CategoryService) GetCategoryByID(id int) (*model.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) CreateCategory(category model.Category) (*model.Category, error) {

	return s.repo.Create(category)
}

func (s *CategoryService) UpdateCategory(id int, category model.Category) (*model.Category, error) {
	return s.repo.Update(id, category)
}

func (s *CategoryService) DeleteCategory(id int) error {
	return s.repo.Delete(id)
}
