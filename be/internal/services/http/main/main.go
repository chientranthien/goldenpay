package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"

	"github.com/chientranthien/goldenpay/internal/services/http/config"
	"github.com/chientranthien/goldenpay/internal/services/http/controller"
	userclient "github.com/chientranthien/goldenpay/internal/services/user/client"
)

func setupRouter() *gin.Engine {
	server := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowCredentials = true
	server.Use(cors.New(corsConfig))

	uClient := userclient.NewUserServiceClient(config.GetDefaultConfig().UserService.Addr)
	signupController := controller.NewSignupController(uClient)
	loginController := controller.NewLoginController(uClient)
	authzController := controller.NewAuthzController(uClient)


	server.POST("api/v1/signup", signupController.Signup)
	server.POST("api/v1/login", loginController.Login)
	server.POST("api/v1/authz", authzController.Authz)
	server.PUT("api/v1/transactions", authzController.Authz)
	// Ping test
	server.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return server
}

func main() {
	r := setupRouter()
	r.Run(config.GetDefaultConfig().HttpService.Addr)
}
