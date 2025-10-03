package models

import "time"

type Profile struct {
	Address   string    `bson:"address" json:"address"`
	Photo     string    `bson:"photo" json:"photo"`
	Bio       string    `bson:"bio" json:"bio"`
	Birthdate time.Time `bson:"birthdate" json:"birthdate"`
	Age       int       `bson:"age" json:"age"`
	Gender    string    `bson:"gender" json:"gender"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
