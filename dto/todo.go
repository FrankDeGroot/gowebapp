package dto

type Todo struct {
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type SavedTodo struct {
	Id string `json:"id"`
	Todo
}
