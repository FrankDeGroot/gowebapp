package dto

type ToDo struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}
