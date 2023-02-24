package service

import (
	"todoListAPI/model"
	"todoListAPI/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(UserID int, list model.TodoList) (int, error) {
	return s.repo.Create(UserID, list)
}

func (s *TodoListService) GetAll(userID int) ([]model.TodoList, error) {
	return s.repo.GetAll(userID)
}

func (s *TodoListService) GetByID(userID, id int) (model.TodoList, error) {
	return s.repo.GetByID(userID, id)
}

func (s *TodoListService) Delete(userID, listID int) error {
	return s.repo.Delete(userID, listID)
}

func (s *TodoListService) Update(userID, listID int, input model.UpdateListInput) error {
	err := input.Validate()
	if err != nil {
		return err
	}
	return s.repo.Update(userID, listID, input)
}
