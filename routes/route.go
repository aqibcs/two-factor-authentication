package routes

import (
	"authentication/handlers"

	"github.com/labstack/echo"
)

// AuthRouteHandler struct that holds the authHandler
type AuthRouteHandler struct {
	authHandler handlers.AuthHandler
}

// NewAuthRouteHandler is a constructor function that returns a new AuthRouteHandler
func NewAuthRouteHandler(authHandler handlers.AuthHandler) AuthRouteHandler {
	return AuthRouteHandler{authHandler}
}

// AuthRoute is a method on AuthRouteHandler that sets up the routes for authentication
func (rh *AuthRouteHandler) AuthRoute(rg *echo.Group) {
	router := rg.Group("/authentication") // Creating a new group of routes under "/authentication"

	router.POST("/register", rh.authHandler.SignUpUser)      // Route for user registration
	router.POST("/login", rh.authHandler.LoginUser)          // Route for user login
	router.POST("/otp/generate", rh.authHandler.GenerateOTP) // Routes for generate otp
	router.POST("/otp/verify", rh.authHandler.VerifyOTP)
	router.POST("/otp/validate", rh.authHandler.ValidateOTP)
	router.POST("/otp/disable", rh.authHandler.DisableOTP)
}
