package main

import (
	"context"
	"fantasy/db"
	"fantasy/delivery"
	"fantasy/handler"
	"fantasy/repository"
	"log"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func dbConnection() db.DbInterface {
	log.Printf("cs: %+v, %+v, %+v, %+v", viper.GetString(`database.url`), viper.GetString(`database.user`), viper.GetString(`database.pass`), viper.GetString(`database.name`))
	db, err := db.Connect(context.Background(), viper.GetString(`database.url`), viper.GetString(`database.user`), viper.GetString(`database.pass`), viper.GetString(`database.name`))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return db
}

func main() {

	db := dbConnection()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	userRepository := repository.NewUserRepository(&db)
	refreshTokenRepository := repository.NewRefreshTokenRepository(&db)
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	userHandler := handler.NewUserHandler(userRepository, refreshTokenRepository, timeoutContext, []byte(viper.GetString(`jwt.secret`)))

	delivery.NewAuthHTTPDelivery(e, userHandler)

	// Start server
	e.Logger.Fatal(e.Start(viper.GetString(`server.address`)))
}
