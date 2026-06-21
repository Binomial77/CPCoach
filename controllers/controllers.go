package controllers

import (
	"net/http"
	"cpcoach/models"
	"cpcoach/database"
	"strings"
	"cpcoach/utils"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
)

func RootController(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", nil)
}

func SignupPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}

func SignupController(c *gin.Context) {

	var req models.SignUpOrLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email cannot be empty",
		})
		return
	}

	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "password cannot be empty",
		})
		return
	}

	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "password must be at least 8 characters",
		})
		return
	}

	var existingUser models.User

	err := database.DB.
		Where("email = ?", req.Email).
		First(&existingUser).Error

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "user already exists",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	rating := models.UserRating{
		UserID:        user.ID,
		CurrentRating: 0,
	}

	if err := database.DB.Create(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to initialize rating",
		})
		return
	}

	stats := models.ProblemStat{
		UserID: user.ID,
	}

	if err := database.DB.Create(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to initialize statistics",
		})
		return
	}

	tokenString, err := utils.GenerateJWT(
		user.ID,
		user.Email,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.SetCookie(
		"jwt",
		tokenString,
		86400,
		"/",
		"",
		false,
		true,
	)
	c.Redirect(
		http.StatusFound,
		"/dashboard",
	)
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func LoginController(c *gin.Context) {

	var req models.SignUpOrLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email cannot be empty",
		})
		return
	}

	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "password cannot be empty",
		})
		return
	}

	var user models.User

	err := database.DB.
		Where("email = ?", req.Email).
		First(&user).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	tokenString, err := utils.GenerateJWT(
		user.ID,
		user.Email,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.SetCookie(
		"jwt",
		tokenString,
		86400,
		"/",
		"",
		false,
		true,
	)

	c.Redirect(
		http.StatusFound,
		"/dashboard",
	)
}

func DashboardController(c *gin.Context) {

	userID := c.MustGet("userID").(uint)

	var rating models.UserRating
	var stats models.ProblemStat

	if err := database.DB.
		Where("user_id = ?", userID).
		First(&rating).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "rating record not found",
		})
		return
	}

	if err := database.DB.
		Where("user_id = ?", userID).
		First(&stats).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "problem statistics not found",
		})
		return
	}

	c.HTML(
		http.StatusOK,
		"dashboard.html",
		gin.H{
			"CurrentRating":     rating.CurrentRating,
			"TotalSolved":       stats.TotalSolved,
			"GraphCount":        stats.GraphCount,
			"TreeCount":         stats.TreeCount,
			"GreedyCount":       stats.GreedyCount,
			"DynamicProgramming": stats.DynamicProgramming,
			"BinarySearch":      stats.BinarySearch,
			"NumberTheory":      stats.NumberTheory,
			"Type1":             stats.Type1,
			"Type2":             stats.Type2,
			"Type3":             stats.Type3,
			"Type4":             stats.Type4,
		},
	)
}

func PostProblemController(c *gin.Context) {

	userID := c.MustGet("userID").(uint)

	var req models.PostProblemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	var stat models.ProblemStat

	err := database.DB.
		Where("user_id = ?", userID).
		First(&stat).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "statistics record not found",
		})
		return
	}

	if req.GraphCount {
		stat.GraphCount++
	}

	if req.TreeCount {
		stat.TreeCount++
	}

	if req.GreedyCount {
		stat.GreedyCount++
	}

	if req.DynamicProgramming {
		stat.DynamicProgramming++
	}

	if req.BinarySearch {
		stat.BinarySearch++
	}

	if req.NumberTheory {
		stat.NumberTheory++
	}

	stat.TotalSolved++

	rating := req.ProblemRating

	if rating <= 1200 {
		stat.Type1++
	} else if rating <= 1500 {
		stat.Type2++
	} else if rating <= 1900 {
		stat.Type3++
	} else {
		stat.Type4++
	}

	if err := database.DB.Save(&stat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update statistics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "problem logged successfully",
	})
}

func UpdateRatingController(c *gin.Context) {

	userID := c.MustGet("userID").(uint)

	var req models.UpdateRatingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	var rating models.UserRating

	err := database.DB.
		Where("user_id = ?", userID).
		First(&rating).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "rating record not found",
		})
		return
	}

	rating.CurrentRating = req.CurrentRating

	if err := database.DB.Save(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update rating",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "rating updated successfully",
		"current_rating": rating.CurrentRating,
	})
}


