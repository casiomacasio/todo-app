package service

import (
	"github.com/casiomacasio/todo-app/internal/domain"
	"github.com/casiomacasio/todo-app/internal/repository"
	"github.com/google/uuid"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GetUserByRefreshTokenAndRefreshTokenId(refresh_token string, refreshTokenUUID uuid.UUID) (int, error)
	ParseToken(token string) (int, error)
	GetUser(username, password string) (domain.User, error)
	GenerateToken(userId int) (string, error)
	GenerateRefreshToken(userId int) (string, string, error)
	RevokeRefreshToken(userId int) error
}

type TodoList interface {
	Create(userId int, list domain.TodoList) (int, error)
	GetAll(userId int) ([]domain.TodoList, error)
	GetById(userId, listId int) (domain.TodoList, error)
	UpdateById(userId, listId int, title, description *string) error
	DeleteById(userId, listId int) error
}

type TodoItem interface {
	Create(userId, listId int, input domain.TodoItem) (int, error)
	GetAllItems(userId, listId int) ([]domain.TodoItem, error)
	GetById(userId, itemId int) (domain.TodoItem, error)
	UpdateById(userId, itemId int, title, description *string, done *bool) error
	DeleteById(userId, itemId int) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList: NewTodoListService(repos.TodoList),
		TodoItem: NewTodoItemService(repos.TodoItem),
	}
}