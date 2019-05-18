package main

import (
	"Go-iris-boilerplate/config"
	"Go-iris-boilerplate/controllers"
	"Go-iris-boilerplate/database"
	"Go-iris-boilerplate/models"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func newApp() (app *iris.Application) {
	app = iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	// app.Logger().SetLevel("debug")

	// app.Use(logger.New())

	// app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
	// 	ctx.JSON(controllers.ApiResource(false, nil, "404 Not Found"))
	// })
	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		ctx.WriteString("Oups something went wrong, try again")

	})
	database.DB.AutoMigrate(
		&models.User{},
		&models.RequestToken{},
	)

	iris.RegisterOnInterrupt(func() {
		database.DB.Close()
	})

	auth := app.Party("/auth").AllowMethods(iris.MethodOptions)
	{
		auth.Post("/signup", controllers.UserSignup)
		auth.Post("/verify-email", controllers.VerifyEmail)
		auth.Post("/login", controllers.UserLogin)
		auth.Post("/forgot-password", controllers.ForgotPassword)
		auth.Post("/reset-password", controllers.ResetPassword)
		// auth.Get("/email", controllers.Email)
	}

	v1 := app.Party("/v1").AllowMethods(iris.MethodOptions)
	{
		v1.Use(jwtHandler.Serve)
		v1.PartyFunc("/users", func(users router.Party) {
			users.Get("/{id:uint}", controllers.GetUser)
			users.Get("/", controllers.GetAllUsers)
			// users.Post("/", controllers.CreateUser)
			users.Put("/{id:uint}", controllers.UpdateUser)
		})
	}

	return
}

func main() {
	app := newApp()
	fmt.Printf("Environment: %s", os.Getenv("ENV_MODE"))

	addr := config.Conf.Get("app.addr").(string)
	app.Run(iris.Addr(addr))
}
