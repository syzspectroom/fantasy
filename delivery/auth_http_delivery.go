package delivery

import (
	"context"
	e "fantasy/error"
	"fantasy/handler"
	"fantasy/model"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// AuthHTTPDelivery struct
type AuthHTTPDelivery struct {
	AuthHandler handler.AuthHandlerInterface
}

// NewAuthHTTPDelivery returns UserHTTPDelivery struct
func NewAuthHTTPDelivery(e *echo.Echo, uh handler.AuthHandlerInterface) {
	AuthHTTPDelivery := &AuthHTTPDelivery{
		AuthHandler: uh,
	}
	e.POST("/login", AuthHTTPDelivery.Login)
	e.POST("/register", AuthHTTPDelivery.Register)
	e.POST("/refresh_token", AuthHTTPDelivery.RefreshToken)
	// e.GET("/users", AuthHTTPDelivery.Find)
	// e.PUT("/users", AuthHTTPDelivery.Update)
	// e.DELETE("/users", AuthHTTPDelivery.Delete)
}

// Register create new user http action
func (ud *AuthHTTPDelivery) Register(c echo.Context) error {
	var reg model.Registration

	err := c.Bind(&reg)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	err = ud.AuthHandler.Register(ctx, &reg)
	if err != nil {
		log.Printf("register error: %+v| user message: %+v", err, e.ErrorMessage(err))
		return echo.NewHTTPError(http.StatusUnprocessableEntity, e.ErrorMessage(err))
	}

	return c.JSON(http.StatusOK, "OK")
}

// Login action
func (ud *AuthHTTPDelivery) Login(c echo.Context) error {
	var login model.LoginInput

	err := c.Bind(&login)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	login.Meta = &model.LoginMeta{
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	loginResponse, err := ud.AuthHandler.Login(ctx, &login)
	if err != nil {
		log.Printf("login error: %+v| user message: %+v", err, e.ErrorMessage(err))
		return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
	}

	return c.JSON(http.StatusOK, loginResponse)
}

// RefreshToken action format: {refresh_token: "token"}
func (ud *AuthHTTPDelivery) RefreshToken(c echo.Context) error {

	refresh := &model.RefreshInputStruct{}
	err := c.Bind(refresh)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	refresh.Meta = &model.LoginMeta{
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	refreshResponse, err := ud.AuthHandler.RefreshByToken(ctx, refresh)
	if err != nil {
		log.Printf("refresh error: %+v| user message: %+v", err, e.ErrorMessage(err))
		return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
	}

	return c.JSON(http.StatusOK, refreshResponse)
}
