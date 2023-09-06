package repository

import (
	"context"
)

const (
	QueryInsertUser = `INSERT INTO users (full_name, phone_number, password) VALUES ($1, $2, $3) RETURNING id`
	QueryGetUserById = `SELECT id, full_name, phone_number FROM users WHERE id = $1`
	QueryGetUserByPhone = `SELECT id, full_name, phone_number, password FROM users WHERE phone_number = $1`
	QueryUpdateUser = `SELECT id, full_name, phone_number, password FROM users WHERE phone_number = $1`
)


// repo insert table users
func (r *Repository) Insert(ctx context.Context, user UserReq) (id int, err error) {
	err = r.DB.QueryRowContext(ctx, QueryInsertUser,
		user.FullName,
		user.PhoneNumber,
		user.Password,
	).Scan(&id)

	// check if any error
	if err != nil {
		return id, err
	}
	return id, err
}


// repo get user by id from table users
func (r *Repository) GetUserByID(ctx context.Context, id int) (user UserRes, err error) {
	err = r.DB.QueryRowContext(ctx, QueryGetUserById,
		id,
	).Scan(&user.Id, &user.FullName, &user.PhoneNumber)

	// check if any error
	if err != nil {
		return user, err
	}
	return user, err
}

// repo update user table users
func (r *Repository) UpdateUser(ctx context.Context, user UserReq, id int) (err error) {
	_, err = r.DB.ExecContext(ctx, QueryGetUserByPhone, user.FullName, user.PhoneNumber, id)

	// check if any error
	if err != nil {
		return err
	}
	return err
}

// repo get user by phone_number from table users
func (r *Repository) GetUserByPhone(ctx context.Context, phoneNumber string) (user UserRes, err error) {
	err = r.DB.QueryRowContext(ctx, QueryUpdateUser, phoneNumber).Scan(&user.Id, &user.FullName, &user.PhoneNumber, &user.Password)

	// check if any error
	if err != nil {
		// check if any error sql no rows
		if err.Error() == "sql: no rows in result set" {
			return user, nil
		}
		return user, err
	}
	return user, nil
}
