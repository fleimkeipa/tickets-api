package main

import (
	"fmt"
	"log"

	"github.com/fleimkeipa/tickets-api/config"
	_ "github.com/fleimkeipa/tickets-api/docs" // which is the generated folder after swag init
	"github.com/fleimkeipa/tickets-api/pkg"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	swagger "github.com/swaggo/echo-swagger"
)

func main() {
	// Load environment configuration
	if err := config.LoadEnv(""); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Create a new Echo instance
	var e = echo.New()

	// Configure Echo settings
	configureEcho(e)

	// Configure CORS middleware
	configureCORS(e)

	// Start the Echo application
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", viper.GetInt("api_service.port"))))
}

// Configures the Echo instance
func configureEcho(e *echo.Echo) {
	e.HideBanner = true
	e.HidePort = true

	// Add Swagger documentation route
	e.GET("/swagger/*", swagger.WrapHandler)

	// Add Recover middleware
	e.Use(middleware.Recover())

	// Create a new validator instance
	e.Validator = pkg.NewValidator()
}

// Configures CORS settings
func configureCORS(e *echo.Echo) {
	var corsConfig = middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	})

	e.Use(corsConfig)
}
