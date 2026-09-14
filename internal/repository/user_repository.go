package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"restapirian/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	List(ctx context.Context, query domain.ListQuery) ([]domain.User, int64, error)
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	IdentityExists(ctx context.Context, username, email string, excludeID int64) (bool, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User, updatePassword bool) error
	Delete(ctx context.Context, id int64) error
	AssignRoleWhereMissing(ctx context.Context, roleID int64) error
}

type userRepository struct{ db *gorm.DB }

var _ UserRepository = (*userRepository)(nil)

func NewUserRepository(db *gorm.DB) UserRepository { return &userRepository{db: db} }

func (r *userRepository) List(ctx context.Context, query domain.ListQuery) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64
	db := r.db.WithContext(ctx).Model(&domain.User{}).
		Joins("LEFT JOIN role ON role.id = user.role_id")
	if query.Search != "" {
		like := "%" + query.Search + "%"
		db = db.Where("user.username LIKE ? OR user.email LIKE ? OR user.name LIKE ? OR role.name LIKE ?", like, like, like, like)
	}
	if err := db.Distinct("user.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}
	sortColumns := map[string]string{
		"id": "user.id", "username": "user.username", "email": "user.email",
		"name": "user.name", "role": "role.name", "created_at": "user.created_at",
	}
	column := sortColumns[query.SortBy]
	if column == "" {
		column = "user.id"
	}
	order := "ASC"
	if strings.EqualFold(query.Order, "desc") {
		order = "DESC"
	}
	err := db.Select("user.*").Preload("Role").Order(fmt.Sprintf("%s %s", column, order)).
		Limit(query.PerPage).Offset((query.Page - 1) * query.PerPage).Find(&users).Error
	return users, total, err
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Preload("Role").First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	return &user, err
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Preload("Role").Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUnauthorized
	}
	return &user, err
}

func (r *userRepository) IdentityExists(ctx context.Context, username, email string, excludeID int64) (bool, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&domain.User{}).Where("username = ? OR email = ?", username, email)
	if excludeID > 0 {
		db = db.Where("id <> ?", excludeID)
	}
	err := db.Count(&count).Error
	return count > 0, err
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *domain.User, updatePassword bool) error {
	values := map[string]any{
		"username": user.Username, "email": user.Email, "name": user.Name,
		"role_id": user.RoleID, "updated_by": user.UpdatedBy,
	}
	if updatePassword {
		values["password"] = user.Password
	}
	result := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", user.ID).Updates(values)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&domain.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *userRepository) AssignRoleWhereMissing(ctx context.Context, roleID int64) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).
		Where("role_id IS NULL OR role_id = 0").Update("role_id", roleID).Error
}
