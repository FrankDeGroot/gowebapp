package mocks

import (
	"todo-app/act"

	"github.com/stretchr/testify/mock"
)

type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) Produce(todoAction *act.TodoAction) error {
	args := m.Called(todoAction)
	return args.Error(0)
}

func (m *MockProducer) Close() {
	m.Called()
}
