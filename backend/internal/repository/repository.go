package repository

import (
	"time"

	"github.com/casiomacasio/todo-app/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
	CreateUser(user domain.CreateUserRequest) (int, error) 
	GetUser(username, password string) (domain.User, error) 
	SaveRefreshToken(hashed_token string, userId int, expires_at time.Time) (uuid.UUID, error)
	GetUserIdAndHashByRefreshTokenId(refreshToken uuid.UUID) (int, string, error) 
	DeleteRefreshToken(refreshToken uuid.UUID) error 
	RevokeRefreshToken(tokenUUID uuid.UUID) (bool, error)
	RevokeRefreshTokenByUserId(userId int) (bool, error)

}

type TodoList interface{
	Create(userId int, list domain.CreateListRequest) (int, error)
	GetAll(userId int) ([]domain.TodoList, error)
	GetById(userId, listId int) (domain.TodoList, error)
	UpdateById(userId, listId int, title, description *string) error
	DeleteById(userId, listId int) error

}

type TodoItem interface{
	Create(userId, listId int, input domain.CreateItemRequest) (int, error)
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