// This file contains types that are used in the repository layer.
package repository

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

// models user request
type UserReq struct {
	FullName      string    `json:"name,omitempty" db:"name"`
	PhoneNumber     string    `json:"phone_number,omitempty" db:"phone_number"`
	Password  string    `json:"password,omitempty" db:"password"`
}

// models user response
type UserRes struct {
	Id        int         `json:"id,omitempty" db:"id"`
	FullName      null.String `json:"name,omitempty" db:"name"`
	PhoneNumber     null.String `json:"phone_number,omitempty" db:"phone_number"`
	Password  string      `json:"password,omitempty" db:"password"`
	CreatedAt time.Time   `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at,omitempty" db:"updated_at"`
}
