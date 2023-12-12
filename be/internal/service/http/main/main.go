package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/service/http/config"
	"github.com/chientranthien/goldenpay/internal/service/http/controller"
	userclient "github.com/chientranthien/goldenpay/internal/service/user/client"
	walletclient "github.com/chientranthien/goldenpay/internal/service/wallet/client"
)

func setupRouter() *gin.Engine {
	server := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "https://goldenpay.chientran.info"}
	corsConfig.AllowCredentials = true
	server.Use(cors.New(corsConfig))
	server.Use(func(ctx *gin.Context) {
		body, _ := io.ReadAll(ctx.Request.Body)
		common.L().Infow(
			"incomingReq",
			"method", ctx.Request.Method,
			"url", ctx.Request.URL,
			"body", string(body),
		)

		ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
	})

	uClient := userclient.NewUserServiceClient(config.Get().UserService.Addr)
	signupController := controller.NewSignupController(uClient)
	loginController := controller.NewLoginController(uClient)
	authzController := controller.NewAuthzController(uClient)
	wClient := walletclient.NewWalletServiceClient(config.Get().WalletService.Addr)
	transferController := controller.NewTransferController(uClient, wClient)
	topupController := controller.NewTopupController(uClient, wClient)
	getUserTransactionsController := controller.NewGetUserTransactionsController(uClient, wClient)
	getUserWalletController := controller.NewGetUserWalletController(uClient, wClient)

	server.POST("api/v1/signup", signupController.Do)
	server.POST("api/v1/login", loginController.Do)
	server.POST("api/v1/authz", authzController.Do)
	server.PUT("api/v1/users/transactions", transferController.Do)
	server.PUT("api/v1/users/topups", topupController.Do)
	server.POST("api/v1/users/transactions/_query", getUserTransactionsController.Do)
	server.GET("api/v1/users/wallets", getUserWalletController.Do)
	// Ping test
	server.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return server
}

func main() {
	r := setupRouter()
	r.Run(config.Get().HttpService.Addr)
}
