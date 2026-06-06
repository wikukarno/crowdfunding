package handler

import (
	"net/http"
	"strings"

	"backend-crowdfunding/auth"
	"backend-crowdfunding/helper"
	"backend-crowdfunding/user"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates the Bearer token, loads the matching user and stores
// it under "currentUser" for the downstream handlers.
func AuthMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			abortUnauthorized(c)
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		token, err := authService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			abortUnauthorized(c)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			abortUnauthorized(c)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			abortUnauthorized(c)
			return
		}

		currentUser, err := userService.GetUserByID(userID)
		if err != nil {
			abortUnauthorized(c)
			return
		}

		c.Set("currentUser", currentUser)
	}
}

func abortUnauthorized(c *gin.Context) {
	response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
	c.AbortWithStatusJSON(http.StatusUnauthorized, response)
}
