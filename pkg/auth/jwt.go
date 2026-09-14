package auth

import (
	"fmt"
	"strconv"
	"time"

	"restapirian/internal/domain"
	"restapirian/internal/service"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret   []byte
	issuer   string
	duration time.Duration
}

type claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTManager(secret, issuer string, duration time.Duration) (*JWTManager, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must contain at least 32 characters")
	}
	if issuer == "" {
		issuer = "restapirian-api"
	}
	if duration <= 0 {
		duration = time.Hour
	}
	return &JWTManager{secret: []byte(secret), issuer: issuer, duration: duration}, nil
}

func (m *JWTManager) Generate(user *domain.User) (string, time.Time, error) {
	if user == nil || user.Role == nil {
		return "", time.Time{}, domain.ErrForbidden
	}
	now := time.Now()
	expiresAt := now.Add(m.duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		UserID: user.ID, Username: user.Username, Role: user.Role.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: m.issuer, Subject: strconv.FormatInt(user.ID, 10),
			IssuedAt: jwt.NewNumericDate(now), NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	})
	signed, err := token.SignedString(m.secret)
	return signed, expiresAt, err
}

func (m *JWTManager) Parse(rawToken string) (*service.AuthClaims, error) {
	parsedClaims := &claims{}
	token, err := jwt.ParseWithClaims(rawToken, parsedClaims, func(token *jwt.Token) (any, error) {
		return m.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithIssuer(m.issuer), jwt.WithExpirationRequired())
	if err != nil || !token.Valid {
		return nil, domain.ErrUnauthorized
	}
	if parsedClaims.UserID < 1 || (parsedClaims.Role != domain.RoleAdmin && parsedClaims.Role != domain.RoleUser) {
		return nil, domain.ErrUnauthorized
	}
	return &service.AuthClaims{UserID: parsedClaims.UserID, Username: parsedClaims.Username, Role: parsedClaims.Role}, nil
}
