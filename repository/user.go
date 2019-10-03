package repository

import (
	"context"
	"fantasy/db"
	"fantasy/model"
	"log"
	"strings"
)

// UserRepository repository object
type UserRepository struct {
	Conn *db.DbInterface
}

// NewUserRepository will create a struct that represent the UseInterface interface
func NewUserRepository(Conn *db.DbInterface) UserRepositoryInterface {
	return &UserRepository{Conn}
}

// Store store user model in db
func (ur *UserRepository) Store(ctx context.Context, u *model.User) error {
	err := (*ur.Conn).Insert(ctx, "user", u)
	if err != nil {
		return err
	}

	return nil
}

// GetByEmail get user struct byg1	et user struct by email \ email
func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := "FOR u IN user FILTER u.email == @email LIMIT 1 RETURN u"
	bindVars := map[string]interface{}{
		"email": strings.TrimSpace(strings.ToLower(email)),
	}
	user := &model.User{}
	_, err := (*ur.Conn).Query(ctx, query, bindVars, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByRefreshToken get User and RefreshToken structs from token
func (ur *UserRepository) GetByRefreshToken(ctx context.Context, token string) (*model.User, *model.RefreshToken, error) {
	query := `FOR rt IN refreshTokens
	FILTER rt.token == @token LIMIT 0, 1
	FOR u IN user
		FILTER u._key == rt.user_id
		RETURN {user: u, refresh: rt}`
	bindVars := map[string]interface{}{
		"token": token,
	}
	type res struct {
		User    *model.User         `json:"user"`
		Refresh *model.RefreshToken `json:"refresh"`
	}
	resStruct := &res{}
	_, err := (*ur.Conn).Query(ctx, query, bindVars, resStruct)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("r str: %+v", resStruct)
	return resStruct.User, resStruct.Refresh, nil
}

// ExistsByEmail check if user exists by email
func (ur *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := "RETURN LENGTH(FOR d IN user FILTER d.email == @email LIMIT 1 RETURN true) > 0"
	bindVars := map[string]interface{}{
		"email": strings.TrimSpace(strings.ToLower(email)),
	}
	var exists bool
	_, err := (*ur.Conn).Query(ctx, query, bindVars, &exists)
	if err != nil {
		log.Printf("exists error: %+v", err)
		return false, err
	}
	log.Printf("exists: %+v", exists)
	return exists, nil
}
