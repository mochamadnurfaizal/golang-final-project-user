package models

import (
	"errors"
	"fmt"
	"golang-final-project-user/helpers"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Users struct {
	ID        uint      `gorm:"primaryKey" json:"id" swagger:"description(id)" example:"1"`
	CreatedAt time.Time `json:"created_at" swagger:"description(Created at)" example:""`
	UpdatedAt time.Time `json:"updated_at" swagger:"description(Updated at)" example:""`
	Username  string    `json:"username" form:"username" query:"username" gorm:"not null;unique;varchar(200)" swagger:"description(Username)" example:"Ijal" valid:"required~Username harus diisi"`
	Fullname  string    `json:"fullname" form:"fullname" query:"fullname" gorm:"not null;varchar(200)" swagger:"description(Full Name)" example:"Nurfaizal" valid:"required~Full name harus diisi"`
	Roles     string    `json:"roles" form:"roles" query:"roles" gorm:"not null;varchar(200)"`
	Password  string    `json:"password" form:"password" query:"password" gorm:"not null;varchar(200)" swagger:"description(Password)" example:"Password" valid:"required~Password harus diisi"`
}

func (e *Users) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("Validate " + strconv.Itoa(int(e.ID)) + " Before Create")

	if len(e.Fullname) < 5 {
		err = errors.New("Fullname Minimal 5 Karakter")
		return
	}

	if len(e.Username) < 5 {
		err = errors.New("Username Minimal 5 Karakter")
		return
	}

	_, errCreate := govalidator.ValidateStruct(e)
	if errCreate != nil {
		err = errCreate
		return
	}

	e.Password = helpers.HashPass(e.Password)

	e.Roles = os.Getenv("ROLES_USER")

	err = nil

	return
}

func (e *Users) BeforeUpdate(tx *gorm.DB) (err error) {
	fmt.Println("Validate " + strconv.Itoa(int(e.ID)) + " Before Update")

	if len(e.Fullname) < 5 {
		err = errors.New("Fullname Minimal 5 Karakter")
		return
	}

	if len(e.Username) < 5 {
		err = errors.New("Username Minimal 5 Karakter")
		return
	}

	_, errCreate := govalidator.ValidateStruct(e)
	if errCreate != nil {
		err = errCreate
		return
	}

	e.Password = helpers.HashPass(e.Password)

	if e.Roles == "" {
		e.Roles = os.Getenv("ROLES_USER")
	}

	err = nil

	return
}
