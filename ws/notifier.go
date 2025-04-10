package ws

import "todo-app/dto"

type Notifier interface {
	Add(*dto.SavedToDo)
	Change(*dto.SavedToDo)
	Delete(*dto.SavedToDo)
}
