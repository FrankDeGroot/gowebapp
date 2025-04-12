package mocks

import (
	"todo-app/dto"

	"github.com/stretchr/testify/mock"
)

type MockTodoRepo struct {
	mock.Mock
}

func (m *MockTodoRepo) GetAll() (*[]dto.SavedTodo, error) {
	args := m.Called()
	return args.Get(0).(*[]dto.SavedTodo), args.Error(1)
}

func (m *MockTodoRepo) GetOne(id string) (*dto.SavedTodo, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.SavedTodo), args.Error(1)
}

func (m *MockTodoRepo) Insert(todo *dto.Todo) (*dto.SavedTodo, error) {
	args := m.Called(todo)
	return args.Get(0).(*dto.SavedTodo), args.Error(1)
}

func (m *MockTodoRepo) Update(todo *dto.SavedTodo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoRepo) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
