package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"restapirian/internal/domain"

	"gorm.io/gorm"
)

type ProductRepository interface {
	List(ctx context.Context, query domain.ListQuery) ([]domain.Product, int64, error)
	FindByID(ctx context.Context, id int32) (*domain.Product, error)
	CodeExists(ctx context.Context, code string, excludeID int32) (bool, error)
	Create(ctx context.Context, product *domain.Product) error
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id int32) error
}

type productRepository struct{ db *gorm.DB }

var _ ProductRepository = (*productRepository)(nil)

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) List(ctx context.Context, query domain.ListQuery) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	db := r.db.WithContext(ctx).Model(&domain.Product{}).
		Joins("LEFT JOIN category ON category.id = product.category_id")
	if query.Search != "" {
		like := "%" + query.Search + "%"
		db = db.Where("product.name LIKE ? OR product.code LIKE ? OR category.name LIKE ?", like, like, like)
	}
	if query.CategoryID != nil {
		db = db.Where("product.category_id = ?", *query.CategoryID)
	}
	if err := db.Distinct("product.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortColumns := map[string]string{
		"id": "product.id", "name": "product.name", "code": "product.code",
		"price": "product.price", "category": "category.name", "created_at": "product.created_at",
	}
	orderColumn := sortColumns[query.SortBy]
	if orderColumn == "" {
		orderColumn = "product.id"
	}
	order := "ASC"
	if strings.EqualFold(query.Order, "desc") {
		order = "DESC"
	}

	err := db.Select("product.*").Preload("Category").
		Order(fmt.Sprintf("%s %s", orderColumn, order)).
		Limit(query.PerPage).Offset((query.Page - 1) * query.PerPage).
		Find(&products).Error
	return products, total, err
}

func (r *productRepository) FindByID(ctx context.Context, id int32) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).Preload("Category").First(&product, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	return &product, err
}

func (r *productRepository) CodeExists(ctx context.Context, code string, excludeID int32) (bool, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&domain.Product{}).Where("code = ?", code)
	if excludeID > 0 {
		db = db.Where("id <> ?", excludeID)
	}
	err := db.Count(&count).Error
	return count > 0, err
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	result := r.db.WithContext(ctx).Model(&domain.Product{}).Where("id = ?", product.ID).Updates(map[string]any{
		"name": product.Name, "code": product.Code, "price": product.Price,
		"category_id": product.CategoryID, "updated_by": product.UpdatedBy,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id int32) error {
	result := r.db.WithContext(ctx).Delete(&domain.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
