package service

import (
	"context"
	"strings"

	"restapirian/internal/domain"
	"restapirian/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserInput struct {
	Username string `json:"username" binding:"required,min=3,max=80"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	RoleID   int64  `json:"role_id" binding:"required,gt=0"`
}

type UpdateUserInput struct {
	Username string `json:"username" binding:"required,min=3,max=80"`
	Password string `json:"password" binding:"omitempty,min=8,max=72"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	RoleID   int64  `json:"role_id" binding:"required,gt=0"`
}

type UserService interface {
	List(ctx context.Context, query domain.ListQuery) (domain.Page[domain.User], error)
	Get(ctx context.Context, id int64) (*domain.User, error)
	Create(ctx context.Context, input CreateUserInput, actorID int64) (*domain.User, error)
	Update(ctx context.Context, id int64, input UpdateUserInput, actorID int64) (*domain.User, error)
	Delete(ctx context.Context, id, actorID int64) error
}

type userService struct {
	users repository.UserRepository
	roles repository.RoleRepository
}

var _ UserService = (*userService)(nil)

func NewUserService(users repository.UserRepository, roles repository.RoleRepository) UserService {
	return &userService{users: users, roles: roles}
}

func (s *userService) List(ctx context.Context, query domain.ListQuery) (domain.Page[domain.User], error) {
	query = normalizeListQuery(query, "id")
	items, total, err := s.users.List(ctx, query)
	return domain.Page[domain.User]{Data: items, Meta: pageMeta(query, total)}, err
}

func (s *userService) Get(ctx context.Context, id int64) (*domain.User, error) {
	if id < 1 {
		return nil, domain.ErrInvalidInput
	}
	return s.users.FindByID(ctx, id)
}

func normalizeUserIdentity(username, email, name string) (string, string, string) {
	return strings.ToLower(strings.TrimSpace(username)), strings.ToLower(strings.TrimSpace(email)), strings.TrimSpace(name)
}

func (s *userService) Create(ctx context.Context, input CreateUserInput, actorID int64) (*domain.User, error) {
	input.Username, input.Email, input.Name = normalizeUserIdentity(input.Username, input.Email, input.Name)
	if len(input.Password) < 8 || len(input.Password) > 72 {
		return nil, domain.ErrInvalidInput
	}
	exists, err := s.users.IdentityExists(ctx, input.Username, input.Email, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrConflict
	}
	role, err := s.roles.FindByID(ctx, input.RoleID)
	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &domain.User{Username: input.Username, Password: string(hash), Email: input.Email, Name: input.Name, RoleID: &role.ID, CreatedBy: actorID, UpdatedBy: actorID}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return s.users.FindByID(ctx, user.ID)
}

func (s *userService) Update(ctx context.Context, id int64, input UpdateUserInput, actorID int64) (*domain.User, error) {
	if id < 1 || len(input.Password) > 72 || (input.Password != "" && len(input.Password) < 8) {
		return nil, domain.ErrInvalidInput
	}
	input.Username, input.Email, input.Name = normalizeUserIdentity(input.Username, input.Email, input.Name)
	exists, err := s.users.IdentityExists(ctx, input.Username, input.Email, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrConflict
	}
	role, err := s.roles.FindByID(ctx, input.RoleID)
	if err != nil {
		return nil, err
	}
	user := &domain.User{ID: id, Username: input.Username, Email: input.Email, Name: input.Name, RoleID: &role.ID, UpdatedBy: actorID}
	if input.Password != "" {
		hash, hashErr := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return nil, hashErr
		}
		user.Password = string(hash)
	}
	if err := s.users.Update(ctx, user, input.Password != ""); err != nil {
		return nil, err
	}
	return s.users.FindByID(ctx, id)
}

func (s *userService) Delete(ctx context.Context, id, actorID int64) error {
	if id < 1 {
		return domain.ErrInvalidInput
	}
	if id == actorID {
		return domain.ErrForbidden
	}
	return s.users.Delete(ctx, id)
}
