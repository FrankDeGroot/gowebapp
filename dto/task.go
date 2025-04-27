package dto

type Task struct {
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type SavedTask struct {
	Id string `json:"id"`
	Task
}
