package service

import (
	"context"
	"errors"
	// "fmt"
	customerrors "koda-b6-backend/internal/errors"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
	"strconv"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (p *ProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	products, err := p.productRepo.GetAll(ctx)

	if err != nil {
		return nil, errors.New("failed to retrieve products")
	}

	if len(products) == 0 {
		return []models.Product{}, nil
	}

	return products, nil
}

func (p *ProductService) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	product, err := p.productRepo.GetByID(ctx, id)

	if err != nil {
		return nil, errors.New("gagal mengambil product")
	}

	return product, nil
}

func (p *ProductService) CreateProduct(ctx context.Context, product *models.Product) error {

	existingProduct, _ := p.productRepo.GetByID(ctx, product.ID)
	if existingProduct != nil {
		return errors.New("product already exists")
	}

	err := p.productRepo.Create(ctx, product)
	if err != nil {
		return errors.New("failed to create product")
	}

	return nil
}

func (p *ProductService) UpdateProduct(ctx context.Context, product *models.Product) error {

	if product.ID == 0 {
		return errors.New("invalid Product Id")
	}

	if product.ProductName == "" {
		return errors.New("product name is required")
	}

	err := p.productRepo.Update(ctx, product)
	if err != nil {
		return errors.New("failed to update product")
	}

	return nil
}

func (p *ProductService) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("user ID cannot be empty")
	}

	idInt, err := strconv.Atoi(id)

	err = p.productRepo.Delete(ctx, idInt)
	if err != nil {
		return errors.New("failed to delete product")
	}

	return nil
}

func (p *ProductService) MostReviewedProduct(ctx context.Context) (*[]models.Product, error) {
	products, err := p.productRepo.MostReview(ctx)
	if err != nil {
		return nil, errors.New("gagal mengambil product")
	}

	return products, nil
}

func (p *ProductService) UpdateStock(ctx context.Context, id, quantity int) error {
	if id <= 0 {
		return customerrors.NewValidationError("product_id", "must be positive")
	}
	if quantity <= 0 {
		return customerrors.NewValidationError("quantity", "must be positive")
	}

	err := p.productRepo.UpdateStock(ctx, id, quantity)
	if err != nil {
		if customerrors.IsNotFoundError(err) {
			return err
		}
		return customerrors.NewDatabaseError("update stock", err)
	}
	return nil
}
