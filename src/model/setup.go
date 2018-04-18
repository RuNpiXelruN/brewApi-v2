package model

import "github.com/jinzhu/gorm"

func SetDatabase(database *gorm.DB) {
	db = database
}

func SeedDatabase(db *gorm.DB) {
	migrateDB()
	seedAll()
}
