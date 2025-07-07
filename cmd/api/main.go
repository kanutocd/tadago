package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/kanutocd/tada/docs"
	"github.com/kanutocd/tada/internal/config"
	"github.com/kanutocd/tada/internal/database"
	"github.com/kanutocd/tada/internal/handler"
	"github.com/kanutocd/tada/internal/middleware"
	"github.com/kanutocd/tada/internal/repository"
	"github.com/kanutocd/tada/internal/service"
)

// @title           Tada API
// @version         1.0
// @description     A simple Todo API application
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io
// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database
	db, err := database.Connect(cfg.Database.URL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate tables
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed database
	if err := database.Seed(db); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	tadaRepo := repository.NewTadaRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	tadaService := service.NewTadaService(tadaRepo, userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	tadaHandler := handler.NewTadaHandler(tadaService)

	// Setup router
	router := setupRouter(userHandler, tadaHandler)

	// Setup server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func setupRouter(userHandler *handler.UserHandler, tadaHandler *handler.TadaHandler) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	v1 := router.Group("/api/v1")
	{
		// User routes
		users := v1.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Tada routes
		tadas := v1.Group("/tadas")
		{
			tadas.GET("", tadaHandler.GetTadas)
			tadas.POST("", tadaHandler.CreateTada)
			tadas.GET("/:id", tadaHandler.GetTada)
			tadas.PUT("/:id", tadaHandler.UpdateTada)
			tadas.DELETE("/:id", tadaHandler.DeleteTada)
		}
	}

	return router
}
