package middleware

import (
	"net/http"

	"cpcoach/utils"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {

	return func(c *gin.Context) {

		tokenString, err := c.Cookie("jwt")

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required",
			})
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		userIDValue, ok := claims["user_id"]

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user id missing in token",
			})
			c.Abort()
			return
		}

		userIDFloat, ok := userIDValue.(float64)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user id format",
			})
			c.Abort()
			return
		}

		userID := uint(userIDFloat)

		c.Set("userID", userID)

		c.Next()
	}
}