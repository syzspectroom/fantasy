package repository

import (
	"context"
	"fantasy/model"
)

// UserRepositoryInterface - interface for user repository
type UserRepositoryInterface interface {
	Store(context.Context, *model.User) error
	GetByEmail(context.Context, string) (*model.User, error)
	ExistsByEmail(context.Context, string) (bool, error)
	GetByRefreshToken(context.Context, string) (*model.User, *model.RefreshToken, error)
}
