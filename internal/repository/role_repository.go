package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"restapirian/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleRepository interface {
	List(ctx context.Context, query domain.ListQuery) ([]domain.Role, int64, error)
	FindByID(ctx context.Context, id int64) (*domain.Role, error)
	FindByName(ctx context.Context, name string) (*domain.Role, error)
	SeedDefaults(ctx context.Context) error
}

type roleRepository struct{ db *gorm.DB }

var _ RoleRepository = (*roleRepository)(nil)

func NewRoleRepository(db *gorm.DB) RoleRepository { return &roleRepository{db: db} }

func (r *roleRepository) List(ctx context.Context, query domain.ListQuery) ([]domain.Role, int64, error) {
	var roles []domain.Role
	var total int64
	db := r.db.WithContext(ctx).Model(&domain.Role{})
	if query.Search != "" {
		db = db.Where("name LIKE ?", "%"+query.Search+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	column := "id"
	if query.SortBy == "name" || query.SortBy == "created_at" {
		column = query.SortBy
	}
	order := "ASC"
	if strings.EqualFold(query.Order, "desc") {
		order = "DESC"
	}
	err := db.Order(fmt.Sprintf("%s %s", column, order)).
		Limit(query.PerPage).Offset((query.Page - 1) * query.PerPage).Find(&roles).Error
	return roles, total, err
}

func (r *roleRepository) FindByID(ctx context.Context, id int64) (*domain.Role, error) {
	var role domain.Role
	err := r.db.WithContext(ctx).First(&role, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	return &role, err
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	return &role, err
}

func (r *roleRepository) SeedDefaults(ctx context.Context) error {
	roles := []domain.Role{{Name: domain.RoleAdmin}, {Name: domain.RoleUser}}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&roles).Error
}
