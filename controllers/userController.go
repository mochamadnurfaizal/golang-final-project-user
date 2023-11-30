package controllers

import (
	"errors"
	"fmt"
	"golang-final-project-user/config"
	"golang-final-project-user/helpers"
	"golang-final-project-user/models"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func RegisterUser(ctx echo.Context) error {
	db := config.GetDB()

	registerUser := models.Users{}

	if err := ctx.Bind(&registerUser); err != nil {
		fmt.Println(err)
		return GenerateErrorResponse(ctx, "Format Data Salah")
	}

	existingUser := models.Users{}
	db.Where("LOWER(username) = ?", strings.ToLower(registerUser.Username)).First(&existingUser)

	if (models.Users{}) != existingUser {
		return GenerateErrorResponse(ctx, "Username Telah Digunakan")
	}
	err := db.Create(&registerUser).Error
	if err != nil {
		fmt.Println(err)
		return GenerateErrorResponse(ctx, err.Error())
	}

	return GenerateSuccessResponse(ctx, "Register User Berhasil", registerUser)
}

func UpdateUser(ctx echo.Context) error {
	db := config.GetDB()

	id := ctx.Param("id")
	if id == "" {
		return GenerateErrorResponse(ctx, "Data ID Tidak Ditemukan")
	}

	updateUser := models.Users{}

	if err := ctx.Bind(&updateUser); err != nil {
		fmt.Println(err)
		return ctx.JSON(http.StatusOK, "Format Data Salah")
	}

	iduint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return GenerateErrorResponse(ctx, err.Error())
	}
	existingUser := models.Users{}
	db.First(&existingUser, iduint)

	if (models.Users{}) == existingUser {
		return GenerateErrorResponse(ctx, "Data dengan ID : "+id+" Tidak Ditemukan")
	}

	existingUser = models.Users{}
	db.Where("LOWER(username) = ?", strings.ToLower(updateUser.Username)).Where("id != ?", iduint).First(&existingUser)

	if (models.Users{}) != existingUser {
		return GenerateErrorResponse(ctx, "Username Telah Digunakan")
	}
	updateUser.ID = uint(iduint)

	if updateUser.Roles == "" {
		updateUser.Roles = os.Getenv("ROLES_USER")
	}

	err = db.Save(&updateUser).Error
	if err != nil {
		return GenerateErrorResponse(ctx, err.Error())
	}

	return GenerateSuccessResponse(ctx, "Update User Berhasil", updateUser)
}

func DeleteUser(ctx echo.Context) error {
	db := config.GetDB()

	id := ctx.Param("id")
	if id == "" {
		return GenerateErrorResponse(ctx, "Data ID Tidak Ditemukan")
	}

	iduint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return GenerateErrorResponse(ctx, err.Error())
	}

	existingUser := models.Users{}
	db.First(&existingUser, iduint)

	if (models.Users{}) == existingUser {
		return GenerateErrorResponse(ctx, "Data dengan ID : "+id+" Tidak Ditemukan")
	}

	err = db.First(&existingUser, iduint).Unscoped().Delete(&existingUser).Error

	if err != nil {
		return GenerateErrorResponse(ctx, err.Error())
	}

	return GenerateSuccessResponse(ctx, "Delete User dengan ID "+id+" Berhasil", nil)
}

func DetailUser(ctx echo.Context) error {
	db := config.GetDB()

	id := ctx.Param("id")
	if id == "" {
		return GenerateErrorResponse(ctx, "Data ID Tidak Ditemukan")
	}

	iduint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return GenerateErrorResponse(ctx, err.Error())
	}

	existingUser := models.Users{}
	db.First(&existingUser, iduint)

	if (models.Users{}) == existingUser {
		return GenerateErrorResponse(ctx, "Data dengan ID : "+id+" Tidak Ditemukan")
	}

	return GenerateSuccessResponse(ctx, "Get User dengan ID "+id+" Berhasil", existingUser)
}

func DataUser(ctx echo.Context) error {
	db := config.GetDB()

	userData, ok := ctx.Get("userData").(*helpers.Claims)

	if !ok {
		return GenerateErrorResponse(ctx, "Gagal Mendapatkan Data User")
	}

	existingUser := models.Users{}
	db.First(&existingUser, userData.ID)

	if (models.Users{}) == existingUser {
		return GenerateErrorResponse(ctx, "Gagal Mendapatkan Data User")
	}

	return GenerateSuccessResponse(ctx, "Get User Berhasil", existingUser)
}

// ListUser godoc
// GetListUser return User Data
// GetListUser godoc
// @Summary Get User Data
// @Description Get User Data
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Users "ok"
// @Router /admin/list-user [get]
func ListUser(ctx echo.Context) error {
	db := config.GetDB()

	existingUser := []models.Users{}
	db.Find(&existingUser)

	return GenerateSuccessResponse(ctx, "Get List User Berhasil", existingUser)
}

// @Summary User Login
// @Description Login and get JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param models.Users body models.Users true "Login"
// @Success 200 {string} string "JWT token"
// @Router /login [post]
func LoginUser(ctx echo.Context) error {
	db := config.GetDB()

	user := models.Users{}

	if err := ctx.Bind(&user); err != nil {
		fmt.Println(err)
		return GenerateErrorResponse(ctx, "Format Data Salah")
	}

	password := user.Password

	err := db.Model(&user).Where("username = ?", user.Username).Take(&user).Error
	if err != nil {
		fmt.Println(err)
		return GenerateErrorResponse(ctx, "Username Not Found")
	}

	comparePass := helpers.ComparePass([]byte(user.Password), []byte(password))
	if !comparePass {
		fmt.Println(err)
		return GenerateErrorResponse(ctx, "Unauthorized Username or Password")
	}

	token := helpers.GenerateToken(uint(user.ID), user.Username, user.Fullname, user.Roles)

	existingToken := models.Token{}
	db.First(&existingToken)

	existingToken.ID = uint(1)
	existingToken.Token = token

	err = db.Save(&existingToken).Error
	if err != nil {
		fmt.Println(err)
		return GenerateErrorResponse(ctx, err.Error())
	}

	return GenerateSuccessResponse(ctx, "Login User Berhasil", map[string]string{
		"token": token,
	})
}

func LogoutUser(ctx echo.Context) error {
	db := config.GetDB()
	tokenString := ctx.Request().Header.Get("Authorization")

	claims, err := helpers.VerifyToken(tokenString)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, models.Response{
			Messages: "Invalid token",
			Success:  false,
		})
	}

	errResponse := errors.New("Token Can't be Used")
	bearer := strings.HasPrefix(tokenString, "Bearer")

	if !bearer {
		return errResponse
	}

	stringToken := strings.Split(tokenString, " ")[1]

	helpers.AddInvalidToken(stringToken)

	if !helpers.IsTokenInvalid(stringToken) {
		return GenerateErrorResponse(ctx, "Gagal Invalidate Token")
	}

	existingToken := models.Token{}
	db.First(&existingToken)

	existingToken.ID = uint(1)
	existingToken.Token = ""

	err = db.Save(&existingToken).Error
	if err != nil {
		fmt.Println(err)
		return GenerateErrorResponse(ctx, err.Error())
	}

	fmt.Printf("User %s logged out or token invalidated\n", claims.Username)

	return GenerateSuccessResponse(ctx, "User logged out or token invalidated", nil)
}

func GetTokenFromDB() string {
	db := config.GetDB()
	existingToken := models.Token{}
	db.First(&existingToken)

	if (models.Token{}) == existingToken {
		return ""
	}

	return "Bearer " + existingToken.Token
}
