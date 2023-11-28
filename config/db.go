package config

import (
	"fmt"
	"golang-final-project-user/models"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DbGorm *gorm.DB
var err error

var host = os.Getenv("PGHOST")
var port = os.Getenv("PGPORT")
var user = os.Getenv("PGUSER")
var password = os.Getenv("PGPASSWORD")
var dbname = os.Getenv("PGDBNAME")

func ConnectGorm() {
	host = os.Getenv("PGHOST")
	port = os.Getenv("PGPORT")
	user = os.Getenv("PGUSER")
	password = os.Getenv("PGPASSWORD")
	dbname = os.Getenv("PGDBNAME")

	portDB, err := strconv.Atoi(os.Getenv("PGPORT"))
	if err != nil {
		panic(err)
	}
	psqlInfo := fmt.Sprintf(`
	host=%s
	port=%d
	user=%s`+`
	password=%s
	dbname=%s
	sslmode=disable`, host, portDB, user, password, dbname)

	DbGorm, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DbGorm.AutoMigrate(models.Users{})
	fmt.Println("Sukses Konek DB")

}

func GetDB() *gorm.DB {
	return DbGorm
}
