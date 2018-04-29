package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// SetDatabase func
func SetDatabase(database *gorm.DB) {
	db = database
}

// SeedDatabase func
func SeedDatabase(db *gorm.DB) {
	migrateDB()
	seedAll()
	fmt.Println("Seed running")
}

// RunMigrations func
func RunMigrations(db *gorm.DB) {
	migrateDB()
	fmt.Println("Migration running")
}
