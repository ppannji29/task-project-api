package models

import "time"

type Profile struct {
	Address   string    `bson:"address"`
	Photo     string    `bson:"photo"`
	Bio       string    `bson:"bio"`
	Birthdate time.Time `bson:"birthdate"`
	Age       int       `bson:"age"`
	Gender    string    `bson:"gender"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
