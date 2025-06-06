package repository

import (
	"fmt"

	"github.com/casiomacasio/todo-app/internal/domain"
	"github.com/jmoiron/sqlx"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db:db}
}

func (r *TodoListPostgres) Create(userId int, list domain.TodoList) (int, error){
	tx, err := r.db.Begin()
	if err != nil {
		return 0, nil
	}
	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}
