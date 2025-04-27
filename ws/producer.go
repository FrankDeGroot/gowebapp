package ws

import "todo-app/act"

type Producer interface {
	Produce(*act.TaskAction) error
	Close()
}
