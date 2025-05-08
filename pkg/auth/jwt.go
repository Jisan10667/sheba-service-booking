package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"service-booking/config"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// TokenType defines the type of token
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// JWTClaims custom claims structure
type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Type   string `json:"type"` // New field to distinguish token type
	jwt.RegisteredClaims
}

// TokenConfig holds JWT configuration
type TokenConfig struct {
	SecretKey             string
	AccessTokenDuration   time.Duration
	RefreshTokenDuration  time.Duration
}

// DefaultTokenConfig provides default JWT configuration
func DefaultTokenConfig() TokenConfig {
	return TokenConfig{
		SecretKey:             getSecretKey(),
		AccessTokenDuration:   24 * time.Hour,
		RefreshTokenDuration:  7 * 24 * time.Hour,
	}
}

// getSecretKey retrieves the secret key from environment or config
func getSecretKey() string {
	// First check environment variable
	if secretKey := os.Getenv("JWT_SECRET_KEY"); secretKey != "" {
		return secretKey
	}

	// Fallback to config file
	// Assuming you'll add a jwt section to your config file
	secretKey := config.String("jwt_secret_key")
	if secretKey == "" {
		// Generate a fallback secret key (DO NOT USE IN PRODUCTION)
		return generateFallbackSecretKey()
	}
	return secretKey
}

// generateFallbackSecretKey creates a fallback secret key
func generateFallbackSecretKey() string {
	fmt.Println("WARNING: Using generated fallback secret key. Set JWT_SECRET_KEY in production!")
	return "fallback-secret-key-change-immediately"
}

// GenerateToken creates a JWT token with specified type and duration
func GenerateToken(userID uint, email, name, role string, tokenType TokenType, config TokenConfig) (string, error) {
	var expirationTime time.Time
	switch tokenType {
	case AccessToken:
		expirationTime = time.Now().Add(config.AccessTokenDuration)
	case RefreshToken:
		expirationTime = time.Now().Add(config.RefreshTokenDuration)
	default:
		return "", errors.New("invalid token type")
	}

	claims := JWTClaims{
		UserID: strconv.FormatUint(uint64(userID), 10),
		Role:   role,
		Email:  email,
		Name:   name,
		Type:   string(tokenType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "service-booking-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}

// GenerateAllTokens creates both access and refresh tokens
func GenerateAllTokens(userID uint, email, name, role string) (accessToken, refreshToken string, err error) {
	config := DefaultTokenConfig()

	accessToken, err = GenerateToken(userID, email, name, role, AccessToken, config)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = GenerateToken(userID, email, name, role, RefreshToken, config)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*JWTClaims, error) {
	config := DefaultTokenConfig()

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// HashPassword securely hashes a password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ComparePasswords compares a hashed password with a plain text password
func ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}