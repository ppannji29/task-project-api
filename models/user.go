package models

import "time"

type User struct {
	UserID    string    `bson:"_id,omitempty" json:"user_id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email"`
	Phone     string    `bson:"phone" json:"phone"`
	IsEnabled bool      `bson:"is_enabled" json:"is_enabled"`
	Profile   *Profile  `bson:"profile,omitempty" json:"profile,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
