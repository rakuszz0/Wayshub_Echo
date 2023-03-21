package database

import (
	"fmt"
	"wayshub/models"
	"wayshub/pkg/mysql"
)

func RunMigration() {
	err := mysql.DB.AutoMigrate(
		&models.channel{},
		&models.comment{},
		&models.Subscribe{},
		&models.video{})

	if err != nil {
		fmt.Println(err)
		panic("Migration Failed")
	}
	fmt.Println("Migration Success")
}
