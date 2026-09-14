package service

import (
	"context"

	"restapirian/internal/domain"
	"restapirian/internal/repository"
)

type RoleService interface {
	List(ctx context.Context, query domain.ListQuery) (domain.Page[domain.Role], error)
}

type roleService struct{ repository repository.RoleRepository }

var _ RoleService = (*roleService)(nil)

func NewRoleService(repository repository.RoleRepository) RoleService {
	return &roleService{repository: repository}
}

func (s *roleService) List(ctx context.Context, query domain.ListQuery) (domain.Page[domain.Role], error) {
	query = normalizeListQuery(query, "id")
	items, total, err := s.repository.List(ctx, query)
	return domain.Page[domain.Role]{Data: items, Meta: pageMeta(query, total)}, err
}
