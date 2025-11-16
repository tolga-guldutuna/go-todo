package todo

import "time"

// Todo = Java'daki @Entity / DTO gibi düşün.
type Todo struct {
	ID        int       `json:"id"`        // SERIAL → int
	Title     string    `json:"title"`     // TEXT
	IsDone    bool      `json:"isDone"`    // BOOLEAN
	CreatedAt time.Time `json:"createdAt"` // TIMESTAMPTZ
	UpdatedAt time.Time `json:"updatedAt"` // TIMESTAMPTZ
}
