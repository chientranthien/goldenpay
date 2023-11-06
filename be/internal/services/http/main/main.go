package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chientranthien/goldenpay/internal/services/http/config"
	"github.com/chientranthien/goldenpay/internal/services/http/controller"
	userclient "github.com/chientranthien/goldenpay/internal/services/user/client"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	r := gin.Default()

	uClient := userclient.NewUserServiceClient(config.GetDefaultConfig().UserService.Addr)
	signupController := controller.NewSignupController(uClient)
	loginController := controller.NewLoginController(uClient)

	r.POST("api/v1/signup", signupController.Signup)
	r.POST("api/v1/login", loginController.Login)
	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(config.GetDefaultConfig().HttpService.Addr)
}
