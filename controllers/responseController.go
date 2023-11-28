package controllers

import (
	"golang-final-project-user/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GenerateErrorResponse(ctx echo.Context, message string) error {
	return ctx.JSON(http.StatusOK, models.Response{
		Messages: message,
		Success:  false,
	})
}

func GenerateSuccessResponse(ctx echo.Context, message string, data interface{}) error {
	return ctx.JSON(http.StatusOK, models.Response{
		Messages: message,
		Success:  true,
		Data:     data,
	})
}
