package parser

import "time"

type Note struct {
	ID        int       `json:"id"`         // Maps to the 'id' column
	Username  string    `json:"username"`   // Maps to the 'username' column
	Title     string    `json:"title"`      // Maps to the 'title' column
	Content   string    `json:"content"`    // Maps to the 'content' column
	CreatedAt time.Time `json:"created_at"` // Maps to the 'created_at' column
	UpdatedAt time.Time `json:"updated_at"` // Maps to the 'updated_at' column
}
