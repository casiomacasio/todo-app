package repository

import (
	"fmt"
	"strings"
	"github.com/casiomacasio/todo-app/internal/domain"
	"github.com/jmoiron/sqlx"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db:db}
}

func (r *TodoListPostgres) Create(userId int, list domain.CreateListRequest) (int, error){
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery,userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]domain.TodoList, error) {
	var lists []domain.TodoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1", todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)
	return lists, err
}

func (r *TodoListPostgres) GetById(userId, listId int) (domain.TodoList, error) {
	var list domain.TodoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul ON tl.id = ul.list_id WHERE ul.user_id = $1 AND tl.id = $2", todoListsTable, usersListsTable)

	err := r.db.Get(&list, query, userId, listId)

	if err != nil {
		return list, fmt.Errorf("failed to do a request: %w", err)
	}

	return list, nil
}


func (r *TodoListPostgres) UpdateById(userId, listId int, title, description *string) error {
	query := fmt.Sprintf("UPDATE %s tl SET", todoListsTable)
	args := []interface{}{}
	argIdx := 1

	if title != nil {
		query += fmt.Sprintf(" title = $%d,", argIdx)
		args = append(args, *title)
		argIdx++
	}
	if description != nil {
		query += fmt.Sprintf(" description = $%d,", argIdx)
		args = append(args, *description)
		argIdx++
	}

	query = strings.TrimSuffix(query, ",")

	query += fmt.Sprintf(" FROM %s ul WHERE tl.id = ul.list_id AND ul.user_id = $%d AND tl.id = $%d", usersListsTable, argIdx, argIdx+1)
	args = append(args, userId, listId)

	_, err := r.db.Exec(query, args...)
	return err
}



func (r *TodoListPostgres) DeleteById(userId, listId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	deleteUserListQuery := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND list_id = $2", usersListsTable)
	if _, err := tx.Exec(deleteUserListQuery, userId, listId); err != nil {
		tx.Rollback()
		return err
	}

	deleteListQuery := fmt.Sprintf("DELETE FROM %s WHERE id = $1", todoListsTable)
	if _, err := tx.Exec(deleteListQuery, listId); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}