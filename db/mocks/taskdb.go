package mocks

import (
	"todo-app/dto"

	"github.com/stretchr/testify/mock"
)

type MockTaskDb struct {
	mock.Mock
}

func (m *MockTaskDb) GetAll() (*[]dto.SavedTask, error) {
	args := m.Called()
	return args.Get(0).(*[]dto.SavedTask), args.Error(1)
}

func (m *MockTaskDb) GetOne(id string) (*dto.SavedTask, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.SavedTask), args.Error(1)
}

func (m *MockTaskDb) Insert(task *dto.Task) (*dto.SavedTask, error) {
	args := m.Called(task)
	return args.Get(0).(*dto.SavedTask), args.Error(1)
}

func (m *MockTaskDb) Update(task *dto.SavedTask) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskDb) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
