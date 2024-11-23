package parser

import "time"

type TodoList struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	Description string    `json:"description"`
	Tasks       []string  `json:"tasks"`
	Urgent      bool      `json:"urgent"`
	Priority    string    `json:"priority"`
	Subtasks    []Subtask `json:"subtasks"`

	Status string `json:"status"`
	Done   bool   `json:"done"`

	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
