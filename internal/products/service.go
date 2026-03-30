package products

import (
	"context"
	"ecom/internal/store"
)

type IService interface {
	GetProducts(ctx context.Context) ([]store.Product, error)
	GetProductByID(ctx context.Context, id int32) (store.Product, error)
	CreateProduct(ctx context.Context, name string, price float64) (store.Product, error)
	DeleteProduct(ctx context.Context, id int32) error
}

func NewService(db *store.Queries) IService {
	return &Service{
		db: db,
	}
}

type Service struct {
	db *store.Queries
}

func (service *Service) GetProducts(ctx context.Context) ([]store.Product, error) {
	products, err := service.db.ListProducts(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (service *Service) CreateProduct(ctx context.Context, name string, price float64) (store.Product, error) {
	// Validate and normalize input
	normalizedName, err := ValidateProductInput(name, price)
	if err != nil {
		return store.Product{}, err
	}

	arg := store.CreateProductParams{
		Name:  normalizedName,
		Price: price,
	}

	product, err := service.db.CreateProduct(ctx, arg)
	if err != nil {
		return store.Product{}, err
	}

	return product, nil
}

func (service *Service) GetProductByID(ctx context.Context, id int32) (store.Product, error) {
	product, err := service.db.GetProduct(ctx, id)
	if err != nil {
		return store.Product{}, err
	}
	return product, nil
}

func (service *Service) DeleteProduct(ctx context.Context, id int32) error {
	return service.db.DeleteProduct(ctx, id)
}
