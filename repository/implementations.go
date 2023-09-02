package repository

import (
	"context"
	"time"
)

func (r *Repository) Insert(ctx context.Context, user UserReq) (id int, err error) {
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

func (r *Repository) GetUserByID(ctx context.Context, id int) (user UserRes, err error) {
	query := `SELECT id, name, phone FROM users WHERE id = $1`

	err = r.DB.QueryRowContext(ctx, query,
		id,
	).Scan(&user.Id, &user.Name, &user.Phone)
	if err != nil {
		return user, err
	}
	return user, err
}

func (r *Repository) EditUser(ctx context.Context, user UserReq, id int) (err error) {
	query := `UPDATE users SET name = $1, phone = $2 WHERE id = $3`

	_, err = r.DB.ExecContext(ctx, query, user.Name, user.Phone, id)

	if err != nil {
		return err
	}
	return err
}

func (r *Repository) GetUserByPhone(ctx context.Context, phone string) (user UserRes, err error) {
	query := `SELECT id, name, phone, password FROM users WHERE phone = $1 `

	err = r.DB.QueryRowContext(ctx, query,
		phone,
	).Scan(&user.Id, &user.Name, &user.Phone, &user.Password)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return user, nil
		}
		return user, err
	}
	return user, nil
}
