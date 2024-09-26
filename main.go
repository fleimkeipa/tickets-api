package main

import (
	"fmt"
	"log"

	"github.com/fleimkeipa/tickets-api/config"
	"github.com/fleimkeipa/tickets-api/controller"
	_ "github.com/fleimkeipa/tickets-api/docs" // which is the generated folder after swag init
	"github.com/fleimkeipa/tickets-api/pkg"
	"github.com/fleimkeipa/tickets-api/repositories"
	"github.com/fleimkeipa/tickets-api/uc"

	"github.com/go-pg/pg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	swagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

func main() {
	// Load environment configuration
	if err := config.LoadEnv(""); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Create a new Echo instance
	e := echo.New()

	// Configure Echo settings
	configureEcho(e)

	// Configure CORS middleware
	configureCORS(e)

	// Configure the logger
	sugar := configureLogger(e)
	defer func() {
		if err := sugar.Sync(); err != nil { // Clean up logger at the end
			log.Fatal(err)
		}
	}()

	// Initialize PostgreSQL client
	dbClient := initDB()
	defer func() {
		if err := dbClient.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	validator := pkg.NewValidator()

	// Create Ticket handlers and related components
	ticketRepo := repositories.NewTicketRepository(dbClient)
	ticketUC := uc.NewTicketUC(ticketRepo, validator)
	ticketHandler := controller.NewTicketHandler(ticketUC)

	// Define Ticket routes
	ticketsRoutes := e.Group("/tickets")
	ticketsRoutes.POST("", ticketHandler.CreateTicket)
	ticketsRoutes.GET("/:id", ticketHandler.GetByID)
	ticketsRoutes.POST("/:id/purchases", ticketHandler.PurchaseTicket)

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
}

// Configures CORS settings
func configureCORS(e *echo.Echo) {
	corsConfig := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	})

	e.Use(corsConfig)
}

// Configures the logger and adds it as middleware
func configureLogger(e *echo.Echo) *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	e.Use(pkg.ZapLogger(logger))

	sugar := logger.Sugar()
	loggerHandler := controller.NewLogger(sugar)
	e.Use(loggerHandler.LoggerMiddleware)

	return sugar
}

// Initializes the PostgreSQL client
func initDB() *pg.DB {
	psqlDB := pkg.NewPSQLClient()
	if psqlDB == nil {
		log.Fatal("Failed to initialize PostgreSQL client")
	}

	log.Println("PostgreSQL client initialized successfully")
	return psqlDB
}
