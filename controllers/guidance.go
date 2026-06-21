package controllers
import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cpcoach/database"
	"cpcoach/models"

	"github.com/gin-gonic/gin"
	"google.golang.org/genai"
)

func GetGuidanceController(c *gin.Context) {

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

	prompt := fmt.Sprintf(`
You are an expert competitive programming coach.

Current Codeforces Rating: %d

Solved Problem Statistics:

Total Solved: %d

Graph: %d
Tree: %d
Greedy: %d
Dynamic Programming: %d
Binary Search: %d
Number Theory: %d

Difficulty Distribution:
<=1200: %d
1201-1500: %d
1501-1900: %d
>=1901: %d

Provide:
1. Strengths
2. Weaknesses
3. Topics needing more practice
4. Recommended rating range
5. A practical 2-week training plan

Keep the response concise.
`,
		rating.CurrentRating,
		stats.TotalSolved,
		stats.GraphCount,
		stats.TreeCount,
		stats.GreedyCount,
		stats.DynamicProgramming,
		stats.BinarySearch,
		stats.NumberTheory,
		stats.Type1,
		stats.Type2,
		stats.Type3,
		stats.Type4,
	)

	client, err := genai.NewClient(
		context.Background(),
		&genai.ClientConfig{
			APIKey: os.Getenv("API_KEY"),
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to initialize Gemini client",
		})
		return
	}

	resp, err := client.Models.GenerateContent(
		context.Background(),
		"gemini-3.5-flash",
		genai.Text(prompt),
		nil,
	)

	if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "error": "Internal server error",
    })
    return
}

	c.JSON(http.StatusOK, gin.H{
		"guidance": resp.Text(),
	})
}