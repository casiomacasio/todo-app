package service

import (
	"github.com/casiomacasio/todo-app/internal/domain"
	"github.com/casiomacasio/todo-app/internal/repository"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(username, password string) (string, error) 
	ParseToken(token string) (int, error)
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