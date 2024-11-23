package parser

import "time"

type Subtask struct {
	ID          int    `json:"id"`
	TodoID      int    `json:"todo_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Done        bool   `json:"done"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
