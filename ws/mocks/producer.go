package mocks

import (
	"todo-app/act"

	"github.com/stretchr/testify/mock"
)

type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) Produce(taskAction *act.TaskAction) error {
	args := m.Called(taskAction)
	return args.Error(0)
}

func (m *MockProducer) Close() {
	m.Called()
}
