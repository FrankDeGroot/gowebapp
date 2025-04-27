package mocks

import (
	"todo-app/dto"

	"github.com/stretchr/testify/mock"
)

type MockTaskRepo struct {
	mock.Mock
}

func (m *MockTaskRepo) GetAll() (*[]dto.SavedTask, error) {
	args := m.Called()
	return args.Get(0).(*[]dto.SavedTask), args.Error(1)
}

func (m *MockTaskRepo) GetOne(id string) (*dto.SavedTask, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.SavedTask), args.Error(1)
}

func (m *MockTaskRepo) Insert(task *dto.Task) (*dto.SavedTask, error) {
	args := m.Called(task)
	return args.Get(0).(*dto.SavedTask), args.Error(1)
}

func (m *MockTaskRepo) Update(task *dto.SavedTask) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepo) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
