package middleware

import (
	"errors"
	"golang-final-project-user/controllers"
	"golang-final-project-user/helpers"
	"golang-final-project-user/models"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func AuthenticationCustomer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			tokenString = controllers.GetTokenFromDB()
		}

		verifyToken, err := helpers.VerifyToken(tokenString)

		if err != nil {
			c.JSON(http.StatusOK, models.Response{
				Messages: err.Error(),
				Success:  false,
			})
			return err
		}

		c.Set("userData", verifyToken)
		return next(c)
	}
}

func AuthenticationAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")

		if tokenString == "" {
			tokenString = controllers.GetTokenFromDB()
		}

		verifyToken, err := helpers.VerifyToken(tokenString)

		if err != nil {
			c.JSON(http.StatusOK, models.Response{
				Messages: err.Error(),
				Success:  false,
			})
			return err
		}

		c.Set("userData", verifyToken)

		rolesAdmin := os.Getenv("ROLES_ADMIN")
		userDataConverted := verifyToken.Roles
		if rolesAdmin != userDataConverted {
			c.JSON(http.StatusOK, models.Response{
				Messages: "User Tidak Memiliki Akses",
				Success:  false,
			})
			return errors.New("User Tidak Memiliki Akses")
		}

		return next(c)
	}
}
