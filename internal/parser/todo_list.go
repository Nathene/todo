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
	GroupID     int       `json:"group_id"`
	Status      string    `json:"status"`
	Done        bool      `json:"done"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
