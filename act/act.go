package act

import "todo-app/dto"

type Action string

const (
	Add    Action = "A"
	Change Action = "C"
	Delete Action = "D"
)

type TodoAction struct {
	Action Action `json:"action"`
	dto.SavedTodo
}

type Notifier func(*TodoAction)

func Make(action Action, savedTodo *dto.SavedTodo) *TodoAction {
	return &TodoAction{
		Action:    action,
		SavedTodo: *savedTodo,
	}
}
