package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

func SetDatabase(database *gorm.DB) {
	db = database
}

func SeedDatabase(db *gorm.DB) {
	migrateDB()
	seedAll()
	fmt.Println("Seed running")
}

func RunMigrations(db *gorm.DB) {
	migrateDB()
	fmt.Println("Migration running")
}
