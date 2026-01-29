package services

import (
	"apikasir/models"
	"apikasir/repositories"
)

type ProductService struct {
	Repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{Repo: repo}
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.Repo.GetAll()
}

func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	return s.Repo.GetByID(id)
}

func (s *ProductService) CreateProduct(p *models.Product) error {
	return s.Repo.Create(p)
}

func (s *ProductService) UpdateProduct(p *models.Product) error {
	return s.Repo.Update(p)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.Repo.Delete(id)
}
