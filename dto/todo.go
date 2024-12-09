package dto

type ToDo struct {
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type SavedToDo struct {
	Id string `json:"id"`
	ToDo
}
