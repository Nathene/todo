package service

type Subtask struct {
	ID          int    `json:"id"`
	TodoID      int    `json:"todo_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}
