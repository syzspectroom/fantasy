package repository

import (
	"context"
	"fantasy/db"
	e "fantasy/error"
	"fantasy/model"
	"log"
)

// RefreshTokenRepository repository object
type RefreshTokenRepository struct {
	Conn *db.DbInterface
}

// NewRefreshTokenRepository will create a struct that represent the UseInterface interface
func NewRefreshTokenRepository(Conn *db.DbInterface) RefreshTokenRepositoryInterface {
	return &RefreshTokenRepository{Conn}
}

// Store RefreshToken into db
func (rtr *RefreshTokenRepository) Store(ctx context.Context, rt *model.RefreshToken) error {
	err := (*rtr.Conn).Insert(ctx, "refreshTokens", rt)
	if err != nil {
		return &e.Error{Op: "RefreshTokenRepository.Store", Err: err}
	}

	return nil
}

// Update record with key. new value is update
func (rtr *RefreshTokenRepository) Update(ctx context.Context, key string, update interface{}) error {
	err := (*rtr.Conn).Update(ctx, "refreshTokens", key, update)
	if err != nil {
		return &e.Error{Op: "RefreshTokenRepository.Update", Err: err}
	}

	return nil
}

// GetByToken get RefreshToken struct from the db
func (rtr *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	query := "FOR u IN refreshTokens FILTER u.token == @token LIMIT 1 RETURN u"
	bindVars := map[string]interface{}{
		"token": token,
	}
	refreshToken := &model.RefreshToken{}
	_, err := (*rtr.Conn).Query(ctx, query, bindVars, refreshToken)
	if err != nil {
		return nil, &e.Error{Op: "RefreshTokenRepository.GetByToken", Err: err}
	}
	return refreshToken, nil
}

// ExistsByToken check if refresh token exist
// TODO: add more validations. like ip or useragent check
func (rtr *RefreshTokenRepository) ExistsByToken(ctx context.Context, token string) (bool, error) {
	query := "RETURN LENGTH(FOR d IN user FILTER d.token == @token LIMIT 1 RETURN true) > 0"
	bindVars := map[string]interface{}{
		"token": token,
	}
	var exists bool
	_, err := (*rtr.Conn).Query(ctx, query, bindVars, &exists)
	if err != nil {
		return false, &e.Error{Op: "RefreshTokenRepository.ExistsByToken", Err: err}
	}
	log.Printf("exists: %+v", exists)
	return exists, nil
}

// GetByUserAgentIPAndUserID return RefreshToken
// TODO: create generic GetByValues func
func (rtr *RefreshTokenRepository) GetByUserAgentIPAndUserID(ctx context.Context, userAgent string, ip string, userID string) (*model.RefreshToken, error) {
	query := "FOR r IN refreshTokens FILTER r.user_agent == @user_agent FILTER r.ip == @ip FILTER r.user_id == @userId LIMIT 1 RETURN u"
	bindVars := map[string]interface{}{
		"user_agent": userAgent,
		"ip":         ip,
		"userId":     userID,
	}
	refreshToken := &model.RefreshToken{}
	_, err := (*rtr.Conn).Query(ctx, query, bindVars, refreshToken)
	if err != nil {
		return nil, &e.Error{Op: "RefreshTokenRepository.GetByUserAgentIPAndUserID", Err: err}
	}

	return refreshToken, nil
}
