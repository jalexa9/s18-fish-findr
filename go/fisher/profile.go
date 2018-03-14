package fisher

import (
	"time"
)

// Profile type represents a users profile.
type Profile struct {
	ID          int        `db:"user_id" json:"user_id"`
	FirstName   string     `db:"first_name" json:"first_name"`
	LastName    string     `db:"last_name" json:"last_name"`
	UserName    string     `db:"user_name" json:"user_name"`
	Password    string     `db:"password" json:"password"`
	PhoneNumber string     `db:"phone_number" json:"phone_number"`
	Email       string     `db:"email_address" json:"email_address"`
	Facebook    string     `db:"facebook_profile_link" json:"facebook_profile"`
	Bio         string     `db:"bio" json:"bio"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	Interest    []Interest `json:"interest"`
}

// Interest represents users interest.
type Interest struct {
	ID   int    `db:"interest_id" json:"interest_id"`
	Type string `db:"interest_type" json:"interest_type"`
}
