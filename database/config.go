package database

import(
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"fmt"
	"log"
	"os"
)

var DB *gorm.DB
func ConnectDB(){
	var err error
	dbname := os.Getenv("DB_NAME")
	DB, err = gorm.Open(sqlite.Open(dbname), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Database connection successfully established!")
}