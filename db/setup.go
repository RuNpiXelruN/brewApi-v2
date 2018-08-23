package db

import (
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	// postgres connection
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	db  *gorm.DB
	err error
)

// Init func
func Init(seed, migrate bool) *gorm.DB {
	for i := 0; i < 10; i++ {
		db, err = gorm.Open("postgres", os.Getenv("DB_PARAMS"))
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		log.Println("Error connecting to DB", err)
	}

	if err = db.DB().Ping(); err != nil {
		log.Println("Error pinging DB", err)
	}

	log.Println("Successfully connected to DB")

	if os.Getenv("SEED") == "true" {
		seedDatabase()
		log.Println("Seeding database..")
		os.Unsetenv("SEED")
	}

	if os.Getenv("MIGRATE") == "true" {
		migrateDatabase()
		log.Println("Migrating database..")
		os.Unsetenv("MIGRATE")
	}

	if os.Getenv("DROP") == "true" {
		dropDatabase()
		log.Println("Dropping and recreating database..")
		os.Unsetenv("DROP")
	}

	if os.Getenv("DROPWITHSEED") == "true" {
		dropWithSeed()
		log.Println("Dropping, recreating, and seeding database..")
		os.Unsetenv("DROPWITHSEED")
	}

	return db
}
