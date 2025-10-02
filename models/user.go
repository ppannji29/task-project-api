package models

import "time"

type User struct {
	UserID    string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	Phone     string    `bson:"phone"`
	IsEnabled bool      `bson:"is_enabled"`
	Profile   *Profile  `bson:"profile,omitempty"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
