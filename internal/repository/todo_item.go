package repository

import (
	"strings"
	"github.com/casiomacasio/todo-app/internal/domain"
	"github.com/jmoiron/sqlx"
	"fmt"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db:db}
}

func (r *TodoItemPostgres) Create(userId, listId int, input domain.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var exists bool
	checkIfExist := fmt.Sprintf("SELECT 1 FROM %s WHERE user_id = $1 AND list_id = $2", usersListsTable)
	err = tx.QueryRow(checkIfExist, userId, listId).Scan(&exists)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if !exists {
		tx.Rollback()
		return 0, err
	}
	var id int
	createItemsQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)
	row := tx.QueryRow(createItemsQuery, input.Title, input.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	createUsersItemQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemsTable)
	_, err = tx.Exec(createUsersItemQuery,listId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}

func (r *TodoItemPostgres) GetAllItems(userId, listId int) ([]domain.TodoItem, error) {
	var items []domain.TodoItem
	getAllItemsQuery := fmt.Sprintf("SELECT ti.id, ti.title, ti.description FROM %s ti INNER JOIN %s li ON ti.id = li.item_id INNER JOIN %s ul ON li.list_id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2", todoItemsTable, listsItemsTable, usersListsTable)
	err := r.db.Select(&items, getAllItemsQuery, userId, listId)
	if err !=nil {
		return nil, err
	}
	return items, nil
}

func (r *TodoItemPostgres) GetById(userId, itemId int) (domain.TodoItem, error) {
	var item domain.TodoItem
	getItemByIdQuery := fmt.Sprintf("SELECT ti.id, ti.title, ti.description FROM %s ti INNER JOIN %s li ON ti.id = li.item_id INNER JOIN %s ul ON li.list_id = ul.list_id WHERE ul.user_id = $1 AND ti.id = $2", todoItemsTable, listsItemsTable, usersListsTable)
	err := r.db.Get(&item, getItemByIdQuery, userId, itemId)
	if err !=nil {
		return domain.TodoItem{}, err
	}
	return item, nil
}

func (r *TodoItemPostgres) UpdateById(userId, itemId int, title, description *string, done *bool) error {
	query := fmt.Sprintf("UPDATE %s ti SET ", todoItemsTable)
	args := []interface{}{}
	argIdx := 1

	setParts := []string{}
	if title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIdx))
		args = append(args, *title)
		argIdx++
	}
	if description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, *description)
		argIdx++
	}
	if done != nil {
		setParts = append(setParts, fmt.Sprintf("done = $%d", argIdx))
		args = append(args, *done)
		argIdx++
	}

	query += strings.Join(setParts, ", ")

	query += fmt.Sprintf(" FROM %s li JOIN %s ul ON li.list_id = ul.list_id WHERE ti.id = li.item_id AND ul.user_id = $%d AND ti.id = $%d", listsItemsTable, usersListsTable, argIdx, argIdx+1)

	args = append(args, userId, itemId)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r TodoItemPostgres) DeleteById(userId, itemId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	deleteListsItemQuery := fmt.Sprintf("DELETE FROM %s li USING %s ul WHERE li.list_id = ul.list_id AND li.item_id = $1 AND ul.user_id = $2", listsItemsTable, usersListsTable)
	if _, err := tx.Exec(deleteListsItemQuery, itemId, userId); err != nil {
		tx.Rollback()
		return err
	}

	deleteItemQuery := fmt.Sprintf("DELETE FROM %s WHERE id = $1", todoItemsTable)
	if _, err := tx.Exec(deleteItemQuery, itemId); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}