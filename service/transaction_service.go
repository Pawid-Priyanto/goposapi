package service

import (
	"api-pos/model"
	"api-pos/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []model.CheckoutItem, useLock bool) (*model.Transaction, error) {
	return s.repo.Create(items)
}

func (s *TransactionService) GetReport(startDate, endDate string) (*model.ReportResponse, error) {
	return s.repo.GetReport(startDate, endDate)
}
