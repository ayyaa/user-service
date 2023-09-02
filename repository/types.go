// This file contains types that are used in the repository layer.
package repository

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type GetTestByIdOutput struct {
	Name string
}

type UserReq struct {
	Name      string    `json:"name,omitempty" db:"name"`
	Phone     string    `json:"phone,omitempty" db:"phone"`
	Password  string    `json:"password,omitempty" db:"password"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type UserRes struct {
	Id        int         `json:"id,omitempty" db:"id"`
	Name      null.String `json:"name,omitempty" db:"name"`
	Phone     null.String `json:"phone,omitempty" db:"phone"`
	Password  string      `json:"password,omitempty" db:"password"`
	CreatedAt time.Time   `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at,omitempty" db:"updated_at"`
}
