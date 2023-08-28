package repository

import "context"

func (r *Repository) Insert(ctx context.Context, user User) (id int, err error) {
	query := `INSERT INTO users (uuid, slug, name, phone, password) VALUES ($1, $2, $4, $5, $6) RETURNING id`

	err = r.DB.QueryRowContext(ctx, query,
		user.UUID,
		user.Slug,
		user.Name,
		user.Phone,
		user.Password,
	).Scan(&id)

	if err != nil {
		return id, err
	}
	return id, err
}
