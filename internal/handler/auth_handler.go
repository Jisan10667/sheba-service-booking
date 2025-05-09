package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"service-booking/internal/model"
	"service-booking/internal/service"
	"service-booking/pkg/auth"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user input
	validate := validator.New()
	if err := validate.Struct(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}


	user.Password = hashedPassword

	// Set default role if not specified
	if user.Role == "" {
		user.Role = model.UserRoleUser
	}

	

	// Register the user
	registeredUser, err := h.authService.Register(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := auth.GenerateAllTokens(
		registeredUser.ID, 
		registeredUser.Email, 
		registeredUser.Name, 
		string(registeredUser.Role),
	)

	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate tokens"})
		return
	}

	// Prepare response (remove sensitive info)
	userResponse := gin.H{
		"id":    registeredUser.ID,
		"email": registeredUser.Email,
		"name":  registeredUser.Name,
		"role":  registeredUser.Role,
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":          userResponse,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Login handles user authentication
// In internal/handler/auth_handler.go
// In internal/handler/auth_handler.go
func (h *AuthHandler) Login(c *gin.Context) {
    var loginRequest struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&loginRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }


    // Authenticate user
    user, err := h.authService.Authenticate(loginRequest.Email, loginRequest.Password)
    if err != nil {
        // Log the specific error
        fmt.Printf("Authentication error for %s: %v\n", loginRequest.Email, err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Generate JWT tokens
    accessToken, refreshToken, err := auth.GenerateAllTokens(
        user.ID, 
        user.Email, 
        user.Name, 
        string(user.Role),
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate tokens"})
        return
    }

    // Prepare response (remove sensitive info)
    userResponse := gin.H{
        "id":    user.ID,
        "email": user.Email,
        "name":  user.Name,
        "role":  user.Role,
    }

    c.JSON(http.StatusOK, gin.H{
        "user":          userResponse,
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    })
}

// GetProfile retrieves the authenticated user's profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Retrieve user ID from the context (set by JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Fetch user profile
	user, err := h.authService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}

	// Prepare response (remove sensitive info)
	profileResponse := gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
	}

	c.JSON(http.StatusOK, profileResponse)
}

// UpdateProfile updates the authenticated user's profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// Retrieve user ID from the context (set by JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Prepare update request struct
	var updateRequest struct {
		Name  string `json:"name" validate:"omitempty,min=2,max=100"`
		Email string `json:"email" validate:"omitempty,email"`
	}

	// Bind and validate input
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the update request
	validate := validator.New()
	if err := validate.Struct(updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update profile
	updatedUser, err := h.authService.UpdateProfile(userID.(uint), &updateRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prepare response (remove sensitive info)
	profileResponse := gin.H{
		"id":    updatedUser.ID,
		"email": updatedUser.Email,
		"name":  updatedUser.Name,
		"role":  updatedUser.Role,
	}

	c.JSON(http.StatusOK, profileResponse)
}

// ChangePassword handles password change for authenticated users
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	// Retrieve user ID from the context (set by JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Prepare password change request struct
	var passwordChangeRequest struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}

	// Bind input
	if err := c.ShouldBindJSON(&passwordChangeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Change password
	err := h.authService.ChangePassword(
		userID.(uint), 
		passwordChangeRequest.CurrentPassword, 
		passwordChangeRequest.NewPassword,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}