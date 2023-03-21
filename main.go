package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"wayshub/database"
)

func main() {
	e := echo.New()
	mysql.DatabaseInit()
	database.RunMigration()

	routes.RouteInit(e.Group("/api/v1"))

	fmt.Println("server running localhost:5000")
	e.Logger.Fatal(e.Start("localhost:5000"))

}
