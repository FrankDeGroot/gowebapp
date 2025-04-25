package act

import "todo-app/dto"

type Verb string

const (
	Post   Verb = "Post"
	Put    Verb = "Put"
	Delete Verb = "Delete"
)

type TodoAction struct {
	Verb Verb `json:"verb"`
	dto.SavedTodo
}

type Notifier func(*TodoAction)

func Make(verb Verb, savedTodo *dto.SavedTodo) *TodoAction {
	return &TodoAction{
		Verb:      verb,
		SavedTodo: *savedTodo,
	}
}
