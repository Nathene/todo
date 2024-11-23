package parser

type TodoGroup struct {
	ID       int         `json:"id"`
	Name     string      `json:"group_name"`
	Username string      `json:"username"`
	TodoList *[]TodoList `json:"todo_list"`
}
