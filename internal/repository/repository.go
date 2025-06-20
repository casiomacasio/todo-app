package repository

import (
	"github.com/casiomacasio/todo-app/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
	"time"
)

const (
	usersTable      = "users"
	todoListsTable  = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
	refreshTokensTable = "refresh_tokens"
)

type Authorization interface{
	CreateUser(user domain.User) (int, error) 
	GetUser(username, password string) (domain.User, error) 
	SaveRefreshToken(hashed_token string, userId int, expires_at time.Time) (uuid.UUID, error)
	GetUserIdAndHashByRefreshTokenId(refreshToken uuid.UUID) (int, string, error) 
	DeleteRefreshToken(refreshToken uuid.UUID) error 
	RevokeRefreshToken(userId int) (bool, error)
}

type TodoList interface{
	Create(userId int, list domain.TodoList) (int, error)
	GetAll(userId int) ([]domain.TodoList, error)
	GetById(userId, listId int) (domain.TodoList, error)
	UpdateById(userId, listId int, title, description *string) error
	DeleteById(userId, listId int) error

}

type TodoItem interface{
	Create(userId, listId int, input domain.TodoItem) (int, error)
	GetAllItems(userId, listId int) ([]domain.TodoItem, error)
	GetById(userId, itemId int) (domain.TodoItem, error)
	UpdateById(userId, itemId int, title, description *string, done *bool) error
	DeleteById(userId, itemId int) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList: NewTodoListPostgres(db),
		TodoItem: NewTodoItemPostgres(db),
	}
}