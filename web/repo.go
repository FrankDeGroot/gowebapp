package web

import "todo-app/dto"

type TaskRepo interface {
	GetAll() (*[]dto.SavedTask, error)
	GetOne(string) (*dto.SavedTask, error)
	Insert(*dto.Task) (*dto.SavedTask, error)
	Update(*dto.SavedTask) error
	Delete(string) error
}
