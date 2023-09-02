// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"
)

type RepositoryInterface interface {
	Insert(ctx context.Context, user UserReq) (id int, err error)
	GetUserByID(ctx context.Context, id int) (user UserRes, err error)
	EditUser(ctx context.Context, user UserReq, id int) (err error)
	GetUserByPhone(ctx context.Context, phone string) (user UserRes, err error)
}
