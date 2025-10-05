package dto

import "time"

type CreateTaskRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	DueDate     time.Time `json:"due_date"`
}
