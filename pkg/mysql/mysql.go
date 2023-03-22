package mysql

import (
	"fmt"
	// "os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DatabaseInit() {
	var err error
	dsn := "root@tcp(127.0.0.1:3306)/wayshub?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// var DB_HOST = os.Getenv("DB_HOST")
	// var DB_USER = os.Getenv("DB_USER")
	// var DB_PASSWORD = os.Getenv("DB_PASSWORD")
	// var DB_NAME = os.Getenv("DB_NAME")
	// var DB_PORT = os.Getenv("DB_PORT")

	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to Database")
}
