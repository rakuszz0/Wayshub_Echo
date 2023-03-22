package database

import (
	"fmt"
	"wayshub/models"
	"wayshub/pkg/mysql"
)

func RunMigration() {
	err := mysql.DB.AutoMigrate(
		&models.Channel{},
		&models.Video{},
		&models.Comments{},
		&models.Subscribe{})

	if err != nil {
		fmt.Println(err)
		panic("Migration Failed")
	}

	fmt.Println("Migration Success")
}
