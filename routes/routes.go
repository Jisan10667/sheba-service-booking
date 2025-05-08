package routes

import (
	"service-booking/db"
	"service-booking/internal/handler"
	"service-booking/internal/middleware"
	"service-booking/internal/repository"
	"service-booking/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the Gin router with all necessary routes
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Initialize database connection
	db := db.GetDB()

	// Initialize repositories
	serviceRepo := repository.NewServiceRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	serviceService := service.NewServiceService(serviceRepo)
	bookingService := service.NewBookingService(bookingRepo, serviceRepo)
	authService := service.NewAuthService(userRepo)

	// Initialize handlers
	serviceHandler := handler.NewServiceHandler(serviceService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	authHandler := handler.NewAuthHandler(authService)

	// Middleware for rate limiting and CORS
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Public routes
	v1 := router.Group("/api/v1")
	{
		// Service routes
		v1.GET("/services", serviceHandler.GetServices)
		v1.GET("/services/:id", serviceHandler.GetServiceByID)

		// Booking routes
		v1.POST("/bookings", bookingHandler.CreateBooking)
		v1.GET("/bookings/:id", bookingHandler.GetBookingByID)

		// Auth routes
		v1.POST("/auth/register", authHandler.Register)
		v1.POST("/auth/login", authHandler.Login)
		
		// Token refresh route (public but requires a valid refresh token)
		v1.POST("/auth/refresh", middleware.RefreshTokenHandler)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.JWTAuth())
	{
		// User profile routes
		protected.GET("/profile", authHandler.GetProfile)
		protected.PUT("/profile", authHandler.UpdateProfile)
	}

	// Admin routes (protected)
	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.JWTAuth())
	admin.Use(middleware.AdminOnly())
	{
		// Admin service routes
		admin.POST("/services", serviceHandler.CreateService)
		admin.PUT("/services/:id", serviceHandler.UpdateService)
		admin.DELETE("/services/:id", serviceHandler.DeleteService)

		// Admin booking routes
		admin.GET("/bookings", bookingHandler.GetBookings)
		admin.PUT("/bookings/:id/status", bookingHandler.UpdateBookingStatus)
	}

	return router
}