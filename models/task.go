package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatusChange struct {
	From      string             `bson:"from" json:"from"`
	To        string             `bson:"to" json:"to"`
	ChangedAt time.Time          `bson:"changed_at" json:"changed_at"`
	ChangedBy primitive.ObjectID `bson:"changed_by" json:"changed_by"`
}

type Task struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title         string             `bson:"title" json:"title"`
	Description   string             `bson:"description" json:"description"`
	Status        string             `bson:"status" json:"status"`
	Priority      string             `bson:"priority" json:"priority"`
	DueDate       time.Time          `bson:"due_date" json:"due_date"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	StatusHistory []StatusChange     `bson:"status_history" json:"status_history"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}
