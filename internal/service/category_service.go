package service

import (
	"context"
	"strings"

	"restapirian/internal/domain"
	"restapirian/internal/repository"
)

type CategoryInput struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type CategoryService interface {
	List(ctx context.Context, query domain.ListQuery) (domain.Page[domain.Category], error)
	Get(ctx context.Context, id int32) (*domain.Category, error)
	Create(ctx context.Context, input CategoryInput) (*domain.Category, error)
	Update(ctx context.Context, id int32, input CategoryInput) (*domain.Category, error)
	Delete(ctx context.Context, id int32) error
}

type categoryService struct{ repository repository.CategoryRepository }

var _ CategoryService = (*categoryService)(nil)

func NewCategoryService(repository repository.CategoryRepository) CategoryService {
	return &categoryService{repository: repository}
}

func (s *categoryService) List(ctx context.Context, query domain.ListQuery) (domain.Page[domain.Category], error) {
	query = normalizeListQuery(query, "id")
	items, total, err := s.repository.List(ctx, query)
	return domain.Page[domain.Category]{Data: items, Meta: pageMeta(query, total)}, err
}

func (s *categoryService) Get(ctx context.Context, id int32) (*domain.Category, error) {
	if id < 1 {
		return nil, domain.ErrInvalidInput
	}
	return s.repository.FindByID(ctx, id)
}

func (s *categoryService) validate(ctx context.Context, input *CategoryInput, excludeID int32) error {
	input.Name = strings.TrimSpace(input.Name)
	input.Code = strings.ToUpper(strings.TrimSpace(input.Code))
	if input.Name == "" || input.Code == "" {
		return domain.ErrInvalidInput
	}
	exists, err := s.repository.CodeExists(ctx, input.Code, excludeID)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrConflict
	}
	return nil
}

func (s *categoryService) Create(ctx context.Context, input CategoryInput) (*domain.Category, error) {
	if err := s.validate(ctx, &input, 0); err != nil {
		return nil, err
	}
	category := &domain.Category{Name: input.Name, Code: input.Code}
	if err := s.repository.Create(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryService) Update(ctx context.Context, id int32, input CategoryInput) (*domain.Category, error) {
	if id < 1 {
		return nil, domain.ErrInvalidInput
	}
	if err := s.validate(ctx, &input, id); err != nil {
		return nil, err
	}
	category := &domain.Category{ID: id, Name: input.Name, Code: input.Code}
	if err := s.repository.Update(ctx, category); err != nil {
		return nil, err
	}
	return s.repository.FindByID(ctx, id)
}

func (s *categoryService) Delete(ctx context.Context, id int32) error {
	if id < 1 {
		return domain.ErrInvalidInput
	}
	return s.repository.Delete(ctx, id)
}
