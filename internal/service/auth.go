package service

import (
	"github.com/casiomacasio/todo-app/internal/domain"
	"github.com/casiomacasio/todo-app/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const salt = "sjfkla;sjfkjaefe4324jifjewajfjefiowejf;weijf"
type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user domain.User) (int,error) {
	user.Password, _ = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func generatePasswordHash(password string) (string,error) {
	password += salt	
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}