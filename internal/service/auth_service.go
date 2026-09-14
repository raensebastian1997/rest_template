package service

import (
	"context"
	"strings"
	"time"

	"restapirian/internal/domain"
	"restapirian/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthClaims struct {
	UserID   int64
	Username string
	Role     string
}

type TokenManager interface {
	Generate(user *domain.User) (string, time.Time, error)
	Parse(token string) (*AuthClaims, error)
}

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=80"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResult struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresAt   time.Time    `json:"expires_at"`
	User        *domain.User `json:"user"`
}

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*AuthResult, error)
	Login(ctx context.Context, input LoginInput) (*AuthResult, error)
}

type authService struct {
	users             repository.UserRepository
	roles             repository.RoleRepository
	tokens            TokenManager
	dummyPasswordHash []byte
}

var _ AuthService = (*authService)(nil)

func NewAuthService(users repository.UserRepository, roles repository.RoleRepository, tokens TokenManager) AuthService {
	dummyHash, _ := bcrypt.GenerateFromPassword([]byte("invalid-password"), bcrypt.DefaultCost)
	return &authService{users: users, roles: roles, tokens: tokens, dummyPasswordHash: dummyHash}
}

func (s *authService) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {
	input.Username = strings.ToLower(strings.TrimSpace(input.Username))
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Name = strings.TrimSpace(input.Name)
	if input.Username == "" || input.Email == "" || input.Name == "" || len(input.Password) < 8 || len(input.Password) > 72 {
		return nil, domain.ErrInvalidInput
	}
	exists, err := s.users.IdentityExists(ctx, input.Username, input.Email, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrConflict
	}
	role, err := s.roles.FindByName(ctx, domain.RoleUser)
	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &domain.User{Username: input.Username, Password: string(hash), Email: input.Email, Name: input.Name, RoleID: &role.ID}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	user.Role = role
	return s.issueToken(user)
}

func (s *authService) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	user, err := s.users.FindByUsername(ctx, strings.ToLower(strings.TrimSpace(input.Username)))
	passwordHash := s.dummyPasswordHash
	if err == nil && user != nil {
		passwordHash = []byte(user.Password)
	}
	passwordErr := bcrypt.CompareHashAndPassword(passwordHash, []byte(input.Password))
	if err != nil || passwordErr != nil {
		return nil, domain.ErrUnauthorized
	}
	if user.Role == nil {
		return nil, domain.ErrForbidden
	}
	return s.issueToken(user)
}

func (s *authService) issueToken(user *domain.User) (*AuthResult, error) {
	token, expiresAt, err := s.tokens.Generate(user)
	if err != nil {
		return nil, err
	}
	return &AuthResult{AccessToken: token, TokenType: "Bearer", ExpiresAt: expiresAt, User: user}, nil
}
