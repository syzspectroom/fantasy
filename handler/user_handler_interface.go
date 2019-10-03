package handler

import (
	"context"
	"fantasy/model"
)

// UserHandlerInterface interface for user handler
type UserHandlerInterface interface {
	Store(context.Context, *model.User) error
	Register(context.Context, *model.Registration) error
	Login(context.Context, *model.LoginInput) (*model.LoginResponse, error)
	RefreshByToken(context.Context, *model.RefreshInputStruct) (*model.LoginResponse, error)
}
