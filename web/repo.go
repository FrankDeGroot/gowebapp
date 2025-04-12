package web

import "todo-app/dto"

type TodoRepo interface {
	GetAll() (*[]dto.SavedTodo, error)
	GetOne(string) (*dto.SavedTodo, error)
	Insert(*dto.Todo) (*dto.SavedTodo, error)
	Update(*dto.SavedTodo) error
	Delete(string) error
}
