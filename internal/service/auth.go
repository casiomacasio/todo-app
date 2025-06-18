package service

import (
	"errors"
	"fmt"
	"crypto/rand"
	"os"
	"time"

	"github.com/casiomacasio/todo-app/internal/domain"
	"github.com/casiomacasio/todo-app/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenTTL = 15 * time.Minute
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user domain.User) (int, error) {
	hashedPassword, err := GeneratePasswordHash(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hashedPassword
	id, err := s.repo.CreateUser(user)
	if err != nil {
		if err == repository.ErrUsernameExists {
			return 0, repository.ErrUsernameExists
		}
		return 0, err
	}
	return id, nil
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, password)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return "", repository.ErrUserNotFound
		}
		if err == repository.ErrInvalidPassword {
			return "", repository.ErrInvalidPassword
		}
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	signingKey := []byte(os.Getenv("signingKey"))
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *AuthService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(os.Getenv("signingKey")), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return 0, ErrInvalidToken
	}
	return claims.UserId, nil
}

func GeneratePasswordHash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+os.Getenv("salt")), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}