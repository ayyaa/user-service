// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/generated"
)

type RepositoryInterface interface {
	Insert(ctx context.Context, user generated.RegisterUserRequest) (id int, err error)
}
