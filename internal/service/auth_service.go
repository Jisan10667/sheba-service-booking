package service

import (
	"errors"
	"fmt"
	"strings"

	"service-booking/internal/model"
	"service-booking/internal/repository"
	"service-booking/pkg/auth"
)

// AuthService interface defines the methods for authentication and user management
type AuthService interface {
	Register(user *model.User) (*model.User, error)
	Authenticate(email, password string) (*model.User, error)
	GetUserByID(id uint) (*model.User, error)
	UpdateProfile(userID uint, updateData interface{}) (*model.User, error)
	ChangePassword(userID uint, currentPassword, newPassword string) error
}

// authService implements AuthService
type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo}
}

// Register handles user registration
func (s *authService) Register(user *model.User) (*model.User, error) {
	// Validate email
	if user.Email == "" {
		return nil, errors.New("email is required")
	}

	// Normalize email
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))

	// Check if email already exists
	existingUser, _ := s.userRepo.FindByEmail(user.Email)
	if existingUser != nil {
		return nil, errors.New("email already in use")
	}

	// Set default role if not specified
	if user.Role == "" {
		user.Role = model.UserRoleUser
	}

	// Create the user
	err := s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

// Authenticate validates user credentials and returns the user
func (s *authService) Authenticate(email, password string) (*model.User, error) {
	// Normalize email
	email = strings.TrimSpace(strings.ToLower(email))

	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if !auth.ComparePasswords(user.Password, password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID
func (s *authService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return user, nil
}

// UpdateProfile updates user profile information
func (s *authService) UpdateProfile(userID uint, updateData interface{}) (*model.User, error) {
	// Find the existing user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Type assert and update fields based on the input
	switch data := updateData.(type) {
	case map[string]interface{}:
		// Handle map input (from JSON)
		if name, ok := data["name"].(string); ok && name != "" {
			user.Name = name
		}
		if email, ok := data["email"].(string); ok && email != "" {
			// Normalize and validate email
			normalizedEmail := strings.TrimSpace(strings.ToLower(email))
			
			// Check if email is already in use by another user
			existingUser, _ := s.userRepo.FindByEmail(normalizedEmail)
			if existingUser != nil && existingUser.ID != userID {
				return nil, errors.New("email already in use")
			}
			
			user.Email = normalizedEmail
		}
	case *struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}:
		// Handle struct input
		if data.Name != "" {
			user.Name = data.Name
		}
		if data.Email != "" {
			// Normalize and validate email
			normalizedEmail := strings.TrimSpace(strings.ToLower(data.Email))
			
			// Check if email is already in use by another user
			existingUser, _ := s.userRepo.FindByEmail(normalizedEmail)
			if existingUser != nil && existingUser.ID != userID {
				return nil, errors.New("email already in use")
			}
			
			user.Email = normalizedEmail
		}
	default:
		return nil, errors.New("invalid update data type")
	}

	// Update the user
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %v", err)
	}

	return user, nil
}

// ChangePassword handles password change for a user
func (s *authService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	// Find the user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify current password
	if !auth.ComparePasswords(user.Password, currentPassword) {
		return errors.New("current password is incorrect")
	}

	// Validate new password
	if len(newPassword) < 8 {
		return errors.New("new password must be at least 8 characters long")
	}

	// Hash the new password
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %v", err)
	}

	// Update the password
	user.Password = hashedPassword
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	return nil
}