package ws

import "todo-app/act"

type Consumer interface {
	Consume() (*act.TaskAction, error)
	Close()
}
