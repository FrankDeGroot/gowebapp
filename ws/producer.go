package ws

import "todo-app/act"

type Producer interface {
	Produce(*act.TodoAction) error
	Close()
}
