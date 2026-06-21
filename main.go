package main

import (
	"cpcoach/database"
	"cpcoach/initializers"
	"cpcoach/models"
	"cpcoach/routes"
	"os"
	"fmt"
)

func init() {
	initializers.LoadEnv()
	database.ConnectDB()
}

func main() {

	err := database.DB.AutoMigrate(
		&models.User{},
		&models.UserRating{},
		&models.ProblemStat{},
	)

	if err != nil {
		panic(err)
	}

	router := routes.CreateRouter()
	fmt.Print("Server started running...")
	router.Run(":"+os.Getenv("PORT"))
}