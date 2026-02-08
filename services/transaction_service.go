package services

import (
	"apikasir/models"
	"apikasir/repositories"
)

type TransactionService struct {
	Repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{Repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	return s.Repo.CreateTransaction(items)
}
