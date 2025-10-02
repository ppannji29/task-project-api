package models

import "time"

type Otp struct {
	ID        string    `bson:"_id,omitempty"`
	UserEmail string    `bson:"user_email"`
	OtpCode   string    `bson:"otp_code"`
	OtpExpiry time.Time `bson:"otp_expiry"`
	IsClaimed bool      `bson:"is_claimed"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
