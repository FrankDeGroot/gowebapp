package dto

type ToDo struct {
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type SavedToDo struct {
	Id string `json:"id"`
	ToDo
}

const ActionAdd = "A"
const ActionChg = "C"
const ActionDel = "D"

type ToDoEvent struct {
	Action string `json:"action"`
	SavedToDo
}
