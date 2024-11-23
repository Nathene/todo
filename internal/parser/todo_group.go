package parser

import "time"

type TodoGroup struct {
	ID        int         `json:"id"`
	Name      string      `json:"group_name"`
	Username  string      `json:"username"`
	TodoList  *[]TodoList `json:"todo_list"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
