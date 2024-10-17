package models

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"desc"`
	IsCompleted bool   `json:"is_completed"`
	CreateAt    string `json:"create_at"`
}
