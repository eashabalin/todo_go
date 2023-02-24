package service

import (
	"todoListAPI/model"
	"todoListAPI/pkg/repository"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userID int, list model.TodoList) (int, error)
	GetAll(userID int) ([]model.TodoList, error)
	GetByID(userID, id int) (model.TodoList, error)
	Delete(userID, listID int) error
	Update(userID, listID int, input model.UpdateListInput) error
}

type TodoItem interface {
	Create(userID, listID int, item model.TodoItem) (int, error)
	GetAll(userID, listID int) ([]model.TodoItem, error)
	GetByID(userID, itemID int) (model.TodoItem, error)
	Delete(userID, itemID int) error
	Update(userID, itemID int, input model.UpdateItemInput) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
