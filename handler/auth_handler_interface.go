package handler

import (
	"context"
	"fantasy/model"
)

// AuthHandlerInterface interface for user handler
type AuthHandlerInterface interface {
	Register(context.Context, *model.Registration) error
	Login(context.Context, *model.LoginInput) (*model.LoginResponse, error)
	RefreshByToken(context.Context, *model.RefreshInputStruct) (*model.LoginResponse, error)
}
