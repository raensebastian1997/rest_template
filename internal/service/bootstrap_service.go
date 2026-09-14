package service

import (
	"context"
	"strings"

	"restapirian/internal/domain"
	"restapirian/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type BootstrapAdminInput struct {
	Username string
	Password string
	Email    string
	Name     string
}

func BootstrapRBAC(ctx context.Context, users repository.UserRepository, roles repository.RoleRepository, admin BootstrapAdminInput) error {
	if err := roles.SeedDefaults(ctx); err != nil {
		return err
	}
	userRole, err := roles.FindByName(ctx, domain.RoleUser)
	if err != nil {
		return err
	}
	if err := users.AssignRoleWhereMissing(ctx, userRole.ID); err != nil {
		return err
	}
	if strings.TrimSpace(admin.Password) == "" {
		return nil
	}
	admin.Username = strings.ToLower(strings.TrimSpace(admin.Username))
	admin.Email = strings.ToLower(strings.TrimSpace(admin.Email))
	admin.Name = strings.TrimSpace(admin.Name)
	if admin.Username == "" {
		admin.Username = "admin"
	}
	if admin.Email == "" {
		admin.Email = "admin@example.local"
	}
	if admin.Name == "" {
		admin.Name = "Administrator"
	}
	exists, err := users.IdentityExists(ctx, admin.Username, admin.Email, 0)
	if err != nil || exists {
		return err
	}
	adminRole, err := roles.FindByName(ctx, domain.RoleAdmin)
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return users.Create(ctx, &domain.User{
		Username: admin.Username, Password: string(hash), Email: admin.Email,
		Name: admin.Name, RoleID: &adminRole.ID,
	})
}
