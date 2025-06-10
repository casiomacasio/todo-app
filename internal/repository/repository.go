package repository

import (
	"github.com/casiomacasio/todo-app/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable      = "users"
	todoListsTable  = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

type Authorization interface{
	CreateUser(user domain.User) (int, error) 
	GetUser(username, password string) (domain.User, error) 
}

type TodoList interface{
	Create(userId int, list domain.TodoList) (int, error)
	GetAll(userId int) ([]domain.TodoList, error)
	GetById(userId, listId int) (domain.TodoList, error)
	UpdateById(userId, listId int, title, description string) (domain.TodoList, error)
	DeleteById(userId, listId int) error

}

type TodoItem interface{}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList: NewTodoListPostgres(db),
	}
}