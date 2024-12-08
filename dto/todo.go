package dto

type ToDo struct {
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type SavedToDo struct {
	Id int `json:"id"`
	ToDo
}
