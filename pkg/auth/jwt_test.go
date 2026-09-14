package auth

import (
	"testing"
	"time"

	"restapirian/internal/domain"
)

func TestJWTManagerGenerateAndParse(t *testing.T) {
	manager, err := NewJWTManager("test-secret-that-is-at-least-32-characters", "test-api", time.Hour)
	if err != nil {
		t.Fatalf("NewJWTManager returned error: %v", err)
	}
	role := &domain.Role{ID: 1, Name: domain.RoleAdmin}
	user := &domain.User{ID: 42, Username: "admin", Role: role}

	token, _, err := manager.Generate(user)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	claims, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if claims.UserID != user.ID || claims.Role != domain.RoleAdmin {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestJWTManagerRejectsShortSecret(t *testing.T) {
	if _, err := NewJWTManager("short", "test-api", time.Hour); err == nil {
		t.Fatal("expected short secret to be rejected")
	}
}
