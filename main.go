package main

import (
	"authentication/config"
	"authentication/handlers"
	"authentication/routes"
	"fmt"
	"log"

	"authentication/db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Define global variables for the server and handlers
var (
	server           *echo.Echo
	AuthHandler      handlers.AuthHandler
	AuthRouteHandler routes.AuthRouteHandler
)

// init function is called before main. It initializes the database and handlers.
func init() {
	// Create the database connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		config.DbHost, config.DbPort, config.DbUser, config.DbName, config.DbPass)

	// Initialize the database
	err := db.InitDB(dsn)
	if err != nil {
		// If there's an error, log it and stop the program
		log.Fatal("Failed to connect to the Database:", err)
	}

	// Initialize the AuthHandler with the database
	AuthHandler = *handlers.NewAuthHandler(db.DB)
	// Initialize the AuthRouteHandler with the AuthHandler
	AuthRouteHandler = routes.NewAuthRouteHandler(AuthHandler)
}

// main function is the entry point of the program
func main() {
	// Create a new echo server
	server = echo.New()
	// Use CORS middleware
	server.Use(middleware.CORS())

	// Add the auth routes to the server
	AuthRouteHandler.AuthRoute(server.Group("/myauth"))

	// Start the server
	server.Start(fmt.Sprintf(":%v", config.HttpPort))
}
