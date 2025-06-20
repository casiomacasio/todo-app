package service

import (
	"errors"
	"crypto/rand"
	"os"
	"time"
	"github.com/google/uuid"
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
	ErrTokenExpired = errors.New("token expired")
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

func (s *AuthService) GenerateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
	})
	signingKey := []byte(os.Getenv("signingKey"))
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *AuthService) RevokeRefreshToken(userId int) error {
	_, err := s.repo.RevokeRefreshToken(userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) GetUserByRefreshTokenAndRefreshTokenId(refresh_token string, refreshTokenUUID uuid.UUID) (int, error) {
	userId, hashedToken, err := s.repo.GetUserIdAndHashByRefreshTokenId(refreshTokenUUID)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(refresh_token))
	if err != nil {
		return 0, err 
	}

	return userId, nil
}

func (s *AuthService) GenerateRefreshToken(userId int) (string, string, error) {
	_, err := s.repo.RevokeRefreshToken(userId)
	if err != nil {
		return "", "", err
	}

	b := make([]byte, 16)
	_, err = rand.Read(b)
	if err != nil {
		return "", "", err
	}

	tokenUUID, err := uuid.FromBytes(b)
	if err != nil {
		return "", "", err
	}

	hashedToken, err := GenerateTokenHash(tokenUUID)
	if err != nil {
		return "", "", err
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	id, err := s.repo.SaveRefreshToken(hashedToken, userId, expiresAt)
	if err != nil {
		return "", "", err
	}

	return id.String(), tokenUUID.String(), nil
}


func (s *AuthService) GetUser(username, password string) (domain.User, error) {
	user, err := s.repo.GetUser(username, password)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return domain.User{}, repository.ErrUserNotFound
		}
		if err == repository.ErrInvalidPassword {
			return domain.User{}, repository.ErrInvalidPassword
		}
		return domain.User{}, err
	}
	return user, nil
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
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors&jwt.ValidationErrorExpired != 0 {
			return 0, ErrTokenExpired 
		}
		return 0, ErrInvalidToken
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
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func GenerateTokenHash(refreshToken uuid.UUID) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken.String()), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}