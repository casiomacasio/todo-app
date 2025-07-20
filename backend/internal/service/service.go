package service

import (
	"context"
	"github.com/casiomacasio/todo-app/backend/internal/domain"
	"github.com/casiomacasio/todo-app/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

var ctx = context.Background()

type Authorization interface {
	CreateUser(user domain.CreateUserRequest) (int, error)
	GetUserByRefreshTokenAndRefreshTokenId(refresh_token string, refreshTokenUUID uuid.UUID) (int, error)
	ParseToken(token string) (int, error)
	GetUser(username, password string) (domain.User, error)
	GenerateToken(userId int) (string, error)
	GenerateRefreshToken(userId int) (string, string, error)
	RevokeRefreshToken(uuid.UUID) error
}

type TodoList interface {
	Create(ctx context.Context, userId int, list domain.CreateListRequest) (int, error)
	GetAll(ctx context.Context, userId int) ([]domain.TodoList, error)
	GetById(ctx context.Context, userId, listId int) (domain.TodoList, error)
	UpdateById(ctx context.Context, userId, listId int, title, description *string) error
	DeleteById(ctx context.Context, userId, listId int) error
}

type TodoItem interface {
	Create(userId, listId int, input domain.CreateItemRequest) (int, error)
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

func NewService(repos *repository.Repository, redis *redis.Client) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, redis),
		TodoList: NewTodoListService(repos.TodoList, redis),
		TodoItem: NewTodoItemService(repos.TodoItem, redis),
	}
}