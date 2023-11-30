package models

import "gorm.io/gorm"

type Token struct {
	gorm.Model
	Token string `json:"token" form:"token" query:"token" gorm:"not null;unique;varchar(200)"`
}
