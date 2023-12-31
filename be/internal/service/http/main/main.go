package main

import (
	httpcommon "github.com/chientranthien/goldenpay/internal/common/http"
	"github.com/chientranthien/goldenpay/internal/service/http/config"
	"github.com/chientranthien/goldenpay/internal/service/http/controller"
	userclient "github.com/chientranthien/goldenpay/internal/service/user/client"
	walletclient "github.com/chientranthien/goldenpay/internal/service/wallet/client"
)

func setupHTTPServer() {
	uClient := userclient.NewUserServiceClient(config.Get().UserService.Addr)
	wClient := walletclient.NewWalletServiceClient(config.Get().WalletService.Addr)

	httpcommon.Init(config.Get().HttpService, uClient)

	httpcommon.RegisterPut(
		"api/v1/users/transactions",
		func() httpcommon.Ctl { return controller.NewTransferController(uClient, wClient) },
		&controller.TransferBody{},
		&controller.TransferData{},
	)
	httpcommon.RegisterPut(
		"api/v1/users/topups",
		func() httpcommon.Ctl { return controller.NewTopupController(wClient) },
		&controller.TopupController{},
		&controller.TopupData{},
	)
	httpcommon.RegisterPost(
		"api/v1/users/transactions/_query",
		func() httpcommon.Ctl { return controller.NewGetUserTransactionsController(uClient, wClient) },
		&controller.GetUserTransactionsBody{},
		&controller.GetUserTransactionsData{},
	)
	httpcommon.RegisterGet(
		"api/v1/users/wallets",
		func() httpcommon.Ctl { return controller.NewGetUserWalletController(wClient) },
		&controller.GetUserWalletBody{},
		&controller.GetUserWalletData{},
	)
}

func main() {
	setupHTTPServer()
}
