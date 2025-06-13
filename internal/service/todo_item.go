package service

import (
	"github.com/casiomacasio/todo-app/internal/domain"
	"github.com/casiomacasio/todo-app/internal/repository"
)

type TodoItemService struct {
	repo repository.TodoItem
}

func NewTodoItemService(repo repository.TodoItem) *TodoItemService {
	return &TodoItemService{repo: repo}
}

func (s *TodoItemService) Create(userId, listId int, list domain.TodoItem) (int, error) {
	return s.repo.Create(userId, listId, list)
}

func (s *TodoItemService) GetAllItems(userId, listId int) ([]domain.TodoItem, error) {
	return s.repo.GetAllItems(userId, listId)
}

func (s *TodoItemService) GetById(userId, itemId int) (domain.TodoItem, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *TodoItemService) UpdateById(userId, itemId int, title, description *string, done *bool) error{
	return s.repo.UpdateById(userId, itemId, title, description, done)
}

func (s *TodoItemService) DeleteById(userId, itemId int) error {
	return s.repo.DeleteById(userId, itemId)
}

