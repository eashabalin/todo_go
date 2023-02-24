package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
	"todoListAPI/model"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listID int, item model.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemID int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)

	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listItemsTable)
	_, err = tx.Exec(createListItemQuery, listID, itemID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemID, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(userID, listID int) ([]model.TodoItem, error) {
	var items []model.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li ON ti.id=li.item_id
									INNER JOIN %s ul ON ul.list_id=li.list_id WHERE li.list_id=$1 AND ul.user_id=$2`,
		todoItemsTable, listItemsTable, userListsTable)
	err := r.db.Select(&items, query, listID, userID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemPostgres) GetByID(userID, itemID int) (model.TodoItem, error) {
	var item model.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li ON ti.id=li.item_id
									INNER JOIN %s ul ON ul.list_id=li.list_id WHERE ti.id=$1 AND ul.user_id=$2`,
		todoItemsTable, listItemsTable, userListsTable)
	err := r.db.Get(&item, query, itemID, userID)
	if err != nil {
		return item, err
	}

	return item, nil
}

func (r *TodoItemPostgres) Delete(userID, itemID int) error {
	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul WHERE ti.id=li.item_id
									AND li.list_id=ul.list_id AND ul.user_id=$1 AND ti.id=$2`,
		todoItemsTable, listItemsTable, userListsTable)
	_, err := r.db.Exec(query, userID, itemID)
	return err
}

func (r *TodoItemPostgres) Update(userID, itemID int, input model.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argID))
		args = append(args, *input.Title)
		argID++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argID))
		args = append(args, *input.Description)
		argID++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argID))
		args = append(args, *input.Done)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul
								WHERE ti.id=li.item_id AND li.list_id=ul.list_id AND ul.user_id=$%d AND ti.id=$%d`,
		todoItemsTable, setQuery, listItemsTable, userListsTable, argID, argID+1)

	args = append(args, userID, itemID)

	logrus.Debugf("updateQuery: %s\n", query)
	logrus.Debugf("args: %s\n", args)

	_, err := r.db.Exec(query, args...)

	return err
}
