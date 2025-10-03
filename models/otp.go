package models

import "time"

type Otp struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	UserEmail string    `bson:"user_email" json:"user_email"`
	OtpCode   string    `bson:"otp_code" json:"otp_code"`
	OtpExpiry time.Time `bson:"otp_expiry" json:"otp_expiry"`
	IsClaimed bool      `bson:"is_claimed" json:"is_claimed"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
