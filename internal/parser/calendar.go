package parser

import "time"

type Event struct {
	ID          int       `json:"id"`          // Maps to 'id' column
	Name        string    `json:"name"`        // Maps to 'name' column
	Description string    `json:"description"` // Maps to 'description' column
	EventDate   time.Time `json:"event_date"`  // Maps to 'event_date' column
	CreatedAt   time.Time `json:"created_at"`  // Maps to 'created_at' column
	DaysLeft    int       `json:"days_left"`
}
