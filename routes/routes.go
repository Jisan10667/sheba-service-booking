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
	categoryRepo := repository.NewCategoryRepository(db)

	// Initialize services
	serviceService := service.NewServiceService(serviceRepo, categoryRepo)
	bookingService := service.NewBookingService(bookingRepo, serviceRepo, userRepo)
	authService := service.NewAuthService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	// Initialize handlers
	serviceHandler := handler.NewServiceHandler(serviceService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Middleware for rate limiting and CORS
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Public routes
	v1 := router.Group("/api/v1")
	{
		// Service routes
		v1.GET("/services", serviceHandler.GetServices)
		v1.GET("/services/:id", serviceHandler.GetServiceByID)
		v1.GET("/services/featured", serviceHandler.GetFeaturedServices)

		// Category routes
		v1.GET("/categories", categoryHandler.GetCategories)
		v1.GET("/categories/:id", categoryHandler.GetCategoryByID)
		v1.GET("/categories/:id/subcategories", categoryHandler.GetSubCategories)

		// Booking routes
		v1.POST("/bookings", bookingHandler.CreateBooking)
		v1.GET("/bookings/:id", bookingHandler.GetBookingByID)
		v1.GET("/bookings/reference/:code", bookingHandler.GetBookingByReferenceCode)

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
		protected.POST("/profile/change-password", authHandler.ChangePassword)

		// User bookings
		protected.GET("/bookings", bookingHandler.GetBookings)
		protected.PUT("/bookings/:id/status", bookingHandler.UpdateBookingStatus)
		protected.DELETE("/bookings/:id", bookingHandler.CancelBooking)
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

		// Admin category routes
		admin.POST("/categories", categoryHandler.CreateCategory)
		admin.PUT("/categories/:id", categoryHandler.UpdateCategory)
		admin.DELETE("/categories/:id", categoryHandler.DeleteCategory)
	}

	return router
}