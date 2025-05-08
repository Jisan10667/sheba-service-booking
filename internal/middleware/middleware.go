package middleware

import (
	"net/http"
	"strings"
	"strconv"
	"service-booking/internal/model"
	"service-booking/pkg/auth"

	"github.com/gin-gonic/gin"
)

// JWTAuth middleware for authenticating JWT tokens
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Check the header format
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := headerParts[1]

		// Validate the token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Check if it's an access token
		if claims.Type != string(auth.AccessToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
			c.Abort()
			return
		}

		// Convert UserID from string to uint
		userID, err := strconv.ParseUint(claims.UserID, 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		// Set user information in the context
		c.Set("user_id", uint(userID))
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("name", claims.Name)

		c.Next()
	}
}

// AdminOnly middleware to restrict access to admin routes
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the role from the context (set by JWTAuth middleware)
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		// Check if the role is admin
		if role != string(model.UserRoleAdmin) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admin rights required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RefreshTokenHandler handles token refresh
func RefreshTokenHandler(c *gin.Context) {
	// Get the refresh token from the request
	refreshTokenString := c.GetHeader("X-Refresh-Token")
	if refreshTokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token is missing"})
		return
	}

	// Validate the refresh token
	claims, err := auth.ValidateToken(refreshTokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Check if it's a refresh token
	if claims.Type != string(auth.RefreshToken) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
		return
	}

	// Convert UserID from string to uint
	userID, err := strconv.ParseUint(claims.UserID, 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Generate new access and refresh tokens
	accessToken, refreshToken, err := auth.GenerateAllTokens(
		uint(userID), 
		claims.Email, 
		claims.Name, 
		claims.Role,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new tokens"})
		return
	}

	// Return the new tokens
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}