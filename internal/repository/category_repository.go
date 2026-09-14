package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"restapirian/internal/domain"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	List(ctx context.Context, query domain.ListQuery) ([]domain.Category, int64, error)
	FindByID(ctx context.Context, id int32) (*domain.Category, error)
	CodeExists(ctx context.Context, code string, excludeID int32) (bool, error)
	Create(ctx context.Context, category *domain.Category) error
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id int32) error
}

type categoryRepository struct{ db *gorm.DB }

var _ CategoryRepository = (*categoryRepository)(nil)

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) List(ctx context.Context, query domain.ListQuery) ([]domain.Category, int64, error) {
	var categories []domain.Category
	var total int64
	db := r.db.WithContext(ctx).Model(&domain.Category{})
	if query.Search != "" {
		like := "%" + query.Search + "%"
		db = db.Where("name LIKE ? OR code LIKE ?", like, like)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	sortColumns := map[string]string{"id": "id", "name": "name", "code": "code", "created_at": "created_at"}
	column := sortColumns[query.SortBy]
	if column == "" {
		column = "id"
	}
	order := "ASC"
	if strings.EqualFold(query.Order, "desc") {
		order = "DESC"
	}
	err := db.Order(fmt.Sprintf("%s %s", column, order)).
		Limit(query.PerPage).Offset((query.Page - 1) * query.PerPage).Find(&categories).Error
	return categories, total, err
}

func (r *categoryRepository) FindByID(ctx context.Context, id int32) (*domain.Category, error) {
	var category domain.Category
	err := r.db.WithContext(ctx).First(&category, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	return &category, err
}

func (r *categoryRepository) CodeExists(ctx context.Context, code string, excludeID int32) (bool, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&domain.Category{}).Where("code = ?", code)
	if excludeID > 0 {
		db = db.Where("id <> ?", excludeID)
	}
	err := db.Count(&count).Error
	return count > 0, err
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	result := r.db.WithContext(ctx).Model(&domain.Category{}).Where("id = ?", category.ID).
		Updates(map[string]any{"name": category.Name, "code": category.Code})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int32) error {
	result := r.db.WithContext(ctx).Delete(&domain.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
