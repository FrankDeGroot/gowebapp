package mocks

import (
	"todo-app/act"

	"github.com/stretchr/testify/mock"
)

type MockConsumer struct {
	mock.Mock
}

func (m *MockConsumer) Consume() (*act.TaskAction, error) {
	args := m.Called()
	return args.Get(0).(*act.TaskAction), args.Error(1)
}

func (m *MockConsumer) Close() {
	m.Called()
}
