package handler

import (
	"context"
	"time"

	e "fantasy/error"
	"fantasy/model"
	"fantasy/repository"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/mold.v2/modifiers"
	validator "gopkg.in/go-playground/validator.v9"
)

type authHandler struct {
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

// NewAuthHandler creates new authHandler struct
func NewAuthHandler(ur repository.UserRepositoryInterface, rtr repository.RefreshTokenRepositoryInterface, timeout time.Duration, jwtSecret []byte) AuthHandlerInterface {
	return &authHandler{
		userRepo:         ur,
		refreshTokenRepo: rtr,
		jwtSecret:        jwtSecret,
		contextTimeout:   timeout,
	}
}

func (uh *authHandler) Login(c context.Context, lm *model.LoginInput) (*model.LoginResponse, error) {
	const op = "authHandler.Login"
	conform := modifiers.New()
	if err := conform.Struct(c, lm); err != nil {
		return nil, &e.Error{Op: op, Err: err}
	}
	validate := validator.New()
	if err := validate.Struct(lm); err != nil {
		return nil, &e.Error{Op: op, Err: err}
	}
	//TODO: handle user not found
	dbUser, err := uh.userRepo.GetByEmail(c, lm.Email)
	if err != nil {
		return nil, &e.Error{Op: op, Err: err}
	}
	//TODO: handle passwords do not match
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.HashedPassword), []byte(lm.Password))
	if err != nil {
		return nil, &e.Error{Code: e.EINVALID, Op: op, Message: err.Error()}
	}

	tokenString, err := uh.generateJWTforUser(dbUser)
	if err != nil {
		return nil, &e.Error{Op: op, Err: err}
	}
	lm.Meta.UserID = dbUser.ID
	refreshToken := &model.RefreshToken{}
	refreshToken.FillWithMeta(lm.Meta)
	err = uh.refreshTokenRepo.Store(c, refreshToken)
	if err != nil {
		return nil, &e.Error{Op: op, Err: err}
	}
	loginResponse := &model.LoginResponse{
		Token:   tokenString,
		Refresh: refreshToken.Token,
	}
	return loginResponse, nil
}

func (uh *authHandler) Register(c context.Context, rm *model.Registration) error {
	const op = "authHandler.Register"
	// ctx, cancel := context.WithTimeout(c, uh.contextTimeout)
	// defer cancel()
	conform := modifiers.New()
	if err := conform.Struct(c, rm); err != nil {
		return &e.Error{Op: op, Err: err}
	}
	validate := validator.New()
	if err := validate.Struct(rm); err != nil {
		return &e.Error{Code: e.EINVALID, Op: op, Message: err.Error()}
	}

	user := &model.User{}
	user.Email = rm.Email

	if err := user.HashPassword(rm.Password); err != nil {
		return &e.Error{Op: op, Err: err}
	}

	// TODO: add custom error for uniq validation
	userExist, err := uh.userRepo.ExistsByEmail(c, user.Email)
	if err != nil {
		return &e.Error{Op: op, Err: err}
	}

	if userExist {
		return &e.Error{Op: op, Code: e.EINVALID, Message: "Email is already taken"}
	}
	if err := uh.userRepo.Store(c, user); err != nil {
		return &e.Error{Op: op, Err: err}
	}

	return nil
}

func (uh *authHandler) RefreshByToken(c context.Context, refreshInput *model.RefreshInputStruct) (*model.LoginResponse, error) {
	const op = "authHandler.RefreshByToken"
	user, refreshToken, err := uh.userRepo.GetByRefreshToken(c, refreshInput.RefreshToken)
	if err != nil {
		return nil, &e.Error{Op: op, Err: err}
	}

	tokenString, err := uh.generateJWTforUser(user)
	if err != nil {
		return nil, &e.Error{Op: op, Err: err}
	}
	refreshInput.Meta.UserID = user.ID
	refreshToken.FillWithMeta(refreshInput.Meta)
	err = uh.refreshTokenRepo.Update(c, refreshToken.ID, refreshToken)
	if err != nil {
		return nil, &e.Error{Op: op, Err: err}
	}
	loginResponse := &model.LoginResponse{
		Token:   tokenString,
		Refresh: refreshToken.Token,
	}

	return loginResponse, nil
}

func (uh *authHandler) generateJWTforUser(dbUser *model.User) (string, error) {
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
		return "", &e.Error{Op: "authHandler.generateJWTforUser", Err: err}
	}

	return tokenString, nil
}
