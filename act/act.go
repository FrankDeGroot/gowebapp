package act

import "todo-app/dto"

type Verb string

const (
	Post   Verb = "Post"
	Put    Verb = "Put"
	Delete Verb = "Delete"
)

type TaskAction struct {
	Verb Verb `json:"verb"`
	dto.SavedTask
}

type Notifier func(*TaskAction)

func Make(verb Verb, savedTask *dto.SavedTask) *TaskAction {
	return &TaskAction{
		Verb:      verb,
		SavedTask: *savedTask,
	}
}
