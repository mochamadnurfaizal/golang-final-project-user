package models

type CreateUsers struct {
	Username string `json:"username" form:"username" query:"username" gorm:"not null;unique;varchar(200)" swagger:"description(Username)" example:"Juann" valid:"required~Username harus diisi"`
	Fullname string `json:"fullname" form:"fullname" query:"fullname" gorm:"not null;varchar(200)" swagger:"description(Full Name)" example:"Juann" valid:"required~Full name harus diisi"`
	Password string `json:"password" form:"password" query:"password" gorm:"not null;varchar(200)" swagger:"description(Password)" example:"Juann" valid:"required~Password harus diisi"`
}

type LoginUsers struct {
	Username string `json:"username" form:"username" query:"username" gorm:"not null;unique;varchar(200)" swagger:"description(Username)" example:"Juann" valid:"required~Username harus diisi"`
	Password string `json:"password" form:"password" query:"password" gorm:"not null;varchar(200)" swagger:"description(Password)" example:"Juann" valid:"required~Password harus diisi"`
}
