package main

import (
	"golang-final-project-user/config"
	"golang-final-project-user/controllers"
	"net/http"
	"os"

	_ "golang-final-project-user/docs"

	midd "golang-final-project-user/middleware"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var PORT = os.Getenv("PORTAPP")

// @Title API Documentation User Service
// @Verison 1
// @Description Golang Final Project User Service, Developed by Moch Nurfaizal, Wasis Witjaksono, Andreas Mangaratua Girsang, M. Irfan Mauluddin, Immanuel Juan Junior Sompotan
// @Host localhost:8088
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	PORT = ":" + os.Getenv("PORTAPP")

	config.ConnectGorm()
	config.StartingApps()

	r := echo.New()

	r.GET("/", func(ctx echo.Context) error {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html	")
	})

	r.GET("/anya", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, config.TextSmugAnya())
	})

	// route for swagger
	r.GET("/swagger/*", echoSwagger.WrapHandler)

	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)

	r.GET("/logout", controllers.LogoutUser)

	groupAdmin := r.Group("/admin")
	groupAdmin.Use(midd.AuthenticationAdmin)
	groupAdmin.POST("/create-user", controllers.RegisterUser)
	groupAdmin.PUT("/update-user/:id", controllers.UpdateUser)
	groupAdmin.GET("/detail-user/:id", controllers.DetailUser)

	groupAdmin.GET("/list-user", controllers.ListUser)
	groupAdmin.DELETE("/delete-user/:id", controllers.DeleteUser)

	groupCustomer := r.Group("/customer")
	groupCustomer.Use(midd.AuthenticationCustomer)
	groupCustomer.GET("/profile", controllers.DataUser)

	r.Logger.Fatal(r.Start(PORT))
}
