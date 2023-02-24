package repository

import (
	"github.com/jmoiron/sqlx"
	"todoListAPI/model"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(username, password string) (model.User, error)
}

type TodoList interface {
	Create(userID int, list model.TodoList) (int, error)
	GetAll(userID int) ([]model.TodoList, error)
	GetByID(userID, id int) (model.TodoList, error)
	Delete(userID, listID int) error
	Update(userID, listID int, input model.UpdateListInput) error
}

type TodoItem interface {
	Create(listID int, item model.TodoItem) (int, error)
	GetAll(userID, listID int) ([]model.TodoItem, error)
	GetByID(userID, itemID int) (model.TodoItem, error)
	Delete(userID, itemID int) error
	Update(userID, itemID int, input model.UpdateItemInput) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
