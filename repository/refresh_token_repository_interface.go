package repository

import (
	"context"
	"fantasy/model"
)

// RefreshTokenRepositoryInterface - interface for refresh_token repository
type RefreshTokenRepositoryInterface interface {
	Store(context.Context, *model.RefreshToken) error
	GetByToken(context.Context, string) (*model.RefreshToken, error)
	ExistsByToken(context.Context, string) (bool, error)
	Update(context.Context, string, interface{}) error
}
