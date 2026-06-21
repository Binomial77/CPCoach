package routes 

import (
  "github.com/gin-gonic/gin"
  "cpcoach/controllers"
  "cpcoach/middleware"
)

func CreateRouter() *gin.Engine{
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.GET("/" , controllers.RootController)
	router.GET("/signup" , controllers.SignupPage)
	router.POST("/signup" , controllers.SignupController)
	router.GET("/login" , controllers.LoginPage)
	router.POST("/login" , controllers.LoginController)
	router.GET("/dashboard" , middleware.Authentication() , controllers.DashboardController)
	router.POST("/postproblem" , middleware.Authentication() , controllers.PostProblemController)
	router.GET("/getguidance" , middleware.Authentication() , controllers.GetGuidanceController)
	router.POST("/updaterating" , middleware.Authentication() , controllers.UpdateRatingController)
	return router
} 