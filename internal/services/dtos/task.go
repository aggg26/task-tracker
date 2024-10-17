package dtos

type CreateTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsComplete  bool   `json:"is_complete"`
}
