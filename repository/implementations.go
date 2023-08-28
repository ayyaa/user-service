package repository

import (
	"context"
	"time"

	"github.com/SawitProRecruitment/UserService/generated"
)

func (r *Repository) Insert(ctx context.Context, user generated.RegisterUserRequest) (id int, err error) {
	now := time.Now()
	query := `INSERT INTO users (name, phone, password, updated_at, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err = r.DB.QueryRowContext(ctx, query,
		user.Name,
		user.Phone,
		user.Password,
		now,
		now,
	).Scan(&id)

	if err != nil {
		return id, err
	}
	return id, err
}
