package service

import (
	"context"
	"fmt"
	"strings"

	"restapirian/internal/domain"
	"restapirian/internal/repository"
)

type ProductInput struct {
	Name       string  `json:"name" binding:"required"`
	Code       string  `json:"code" binding:"required"`
	Price      float64 `json:"price" binding:"gte=0"`
	CategoryID int32   `json:"category_id" binding:"required,gt=0"`
}

type ProductService interface {
	List(ctx context.Context, query domain.ListQuery) (domain.Page[domain.Product], error)
	Get(ctx context.Context, id int32) (*domain.Product, error)
	Create(ctx context.Context, input ProductInput, actorID int64) (*domain.Product, error)
	Update(ctx context.Context, id int32, input ProductInput, actorID int64) (*domain.Product, error)
	Delete(ctx context.Context, id int32) error
}

type productService struct {
	products   repository.ProductRepository
	categories repository.CategoryRepository
}

var _ ProductService = (*productService)(nil)

func NewProductService(products repository.ProductRepository, categories repository.CategoryRepository) ProductService {
	return &productService{products: products, categories: categories}
}

func (s *productService) List(ctx context.Context, query domain.ListQuery) (domain.Page[domain.Product], error) {
	query = normalizeListQuery(query, "id")
	items, total, err := s.products.List(ctx, query)
	return domain.Page[domain.Product]{Data: items, Meta: pageMeta(query, total)}, err
}

func (s *productService) Get(ctx context.Context, id int32) (*domain.Product, error) {
	if id < 1 {
		return nil, domain.ErrInvalidInput
	}
	return s.products.FindByID(ctx, id)
}

func (s *productService) validate(ctx context.Context, input *ProductInput, excludeID int32) error {
	input.Name = strings.TrimSpace(input.Name)
	input.Code = strings.ToUpper(strings.TrimSpace(input.Code))
	if input.Name == "" || input.Code == "" || input.CategoryID < 1 || input.Price < 0 {
		return domain.ErrInvalidInput
	}
	if _, err := s.categories.FindByID(ctx, input.CategoryID); err != nil {
		return fmt.Errorf("category: %w", err)
	}
	exists, err := s.products.CodeExists(ctx, input.Code, excludeID)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrConflict
	}
	return nil
}

func (s *productService) Create(ctx context.Context, input ProductInput, actorID int64) (*domain.Product, error) {
	if err := s.validate(ctx, &input, 0); err != nil {
		return nil, err
	}
	categoryID := input.CategoryID
	product := &domain.Product{
		Name: input.Name, Code: input.Code, Price: input.Price, CategoryID: &categoryID, CreatedBy: actorID, UpdatedBy: actorID,
	}
	if err := s.products.Create(ctx, product); err != nil {
		return nil, err
	}
	return s.products.FindByID(ctx, product.ID)
}

func (s *productService) Update(ctx context.Context, id int32, input ProductInput, actorID int64) (*domain.Product, error) {
	if id < 1 {
		return nil, domain.ErrInvalidInput
	}
	if err := s.validate(ctx, &input, id); err != nil {
		return nil, err
	}
	categoryID := input.CategoryID
	product := &domain.Product{ID: id, Name: input.Name, Code: input.Code, Price: input.Price, CategoryID: &categoryID, UpdatedBy: actorID}
	if err := s.products.Update(ctx, product); err != nil {
		return nil, err
	}
	return s.products.FindByID(ctx, id)
}

func (s *productService) Delete(ctx context.Context, id int32) error {
	if id < 1 {
		return domain.ErrInvalidInput
	}
	return s.products.Delete(ctx, id)
}
