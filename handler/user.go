package handler

import (
	"context"
	"strings"
	"time"

	"fantasy/model"
	"fantasy/repository"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/go-playground/validator.v9"
)

type userHandler struct {
	userRepo         repository.UserRepositoryInterface
	refreshTokenRepo repository.RefreshTokenRepositoryInterface
	contextTimeout   time.Duration
	jwtSecret        []byte
}

type jwtClaim struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

// NewUserHandler creates new userHandler struct
func NewUserHandler(ur repository.UserRepositoryInterface, rtr repository.RefreshTokenRepositoryInterface, timeout time.Duration, jwtSecret []byte) UserHandlerInterface {
	return &userHandler{
		userRepo:         ur,
		refreshTokenRepo: rtr,
		jwtSecret:        jwtSecret,
		contextTimeout:   timeout,
	}
}

func (uh *userHandler) Store(c context.Context, um *model.User) error {
	ctx, cancel := context.WithTimeout(c, uh.contextTimeout)
	defer cancel()

	if err := uh.userRepo.Store(ctx, um); err != nil {
		return err
	}
	return nil
}

func (uh *userHandler) Login(c context.Context, lm *model.LoginInput) (*model.LoginResponse, error) {
	validate := validator.New()
	if err := validate.Struct(lm); err != nil {
		return nil, err
	}
	//TODO: handle user not found
	dbUser, err := uh.userRepo.GetByEmail(c, lm.Email)
	if err != nil {
		return nil, err
	}
	//TODO: handle passwords do not match
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.HashedPassword), []byte(lm.Password))
	if err != nil {
		return nil, err
	}

	tokenString, err := uh.generateJWTforUser(dbUser)
	if err != nil {
		return nil, err
	}
	lm.Meta.UserID = dbUser.ID
	refreshToken := &model.RefreshToken{}
	refreshToken.FillWithMeta(lm.Meta)
	err = uh.refreshTokenRepo.Store(c, refreshToken)
	if err != nil {
		return nil, err
	}
	loginResponse := &model.LoginResponse{
		Token:   tokenString,
		Refresh: refreshToken.Token,
	}
	return loginResponse, nil
}

func (uh *userHandler) Register(c context.Context, rm *model.Registration) error {
	// ctx, cancel := context.WithTimeout(c, uh.contextTimeout)
	// defer cancel()

	validate := validator.New()
	if err := validate.Struct(rm); err != nil {
		return err
	}

	user := &model.User{}
	user.Email = strings.TrimSpace(strings.ToLower(rm.Email))

	if err := user.HashPassword(rm.Password); err != nil {
		return err
	}

	// TODO: add custom error for uniq validation
	_, err := uh.userRepo.ExistsByEmail(c, user.Email)
	if err != nil {
		return err
	}

	if err := uh.userRepo.Store(c, user); err != nil {
		return err
	}

	return nil
}

func (uh *userHandler) RefreshByToken(c context.Context, refreshInput *model.RefreshInputStruct) (*model.LoginResponse, error) {
	user, refreshToken, err := uh.userRepo.GetByRefreshToken(c, refreshInput.RefreshToken)
	if err != nil {
		return nil, err
	}

	tokenString, err := uh.generateJWTforUser(user)
	if err != nil {
		return nil, err
	}
	refreshInput.Meta.UserID = user.ID
	refreshToken.FillWithMeta(refreshInput.Meta)
	err = uh.refreshTokenRepo.Update(c, refreshToken.ID, refreshToken)
	if err != nil {
		return nil, err
	}
	loginResponse := &model.LoginResponse{
		Token:   tokenString,
		Refresh: refreshToken.Token,
	}

	return loginResponse, nil
}

func (uh *userHandler) generateJWTforUser(dbUser *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim{
		dbUser.ID,
		dbUser.Email,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		},
	})

	tokenString, err := token.SignedString(uh.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
