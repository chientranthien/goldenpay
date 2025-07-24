package main

import (
	"github.com/chientranthien/goldenpay/internal/common"
	httpcommon "github.com/chientranthien/goldenpay/internal/common/http"
	"github.com/chientranthien/goldenpay/internal/service/http/config"
	"github.com/chientranthien/goldenpay/internal/service/http/controller"
	userclient "github.com/chientranthien/goldenpay/internal/service/user/client"
	walletclient "github.com/chientranthien/goldenpay/internal/service/wallet/client"
)

func setupHTTPServer() {
	common.L().Infow("userService", "config", config.Get().UserService)
	common.L().Infow("walletService", "config", config.Get().WalletService)
	uClient := userclient.NewUserServiceClient(config.Get().UserService.Addr)
	wClient := walletclient.NewWalletServiceClient(config.Get().WalletService.Addr)

	httpcommon.Init(config.Get().HttpService, uClient)

	httpcommon.RegisterPut(httpcommon.PutEndpointInfo{
		EP:       "api/v1/users/transactions",
		NewCtlFn: func() httpcommon.Ctl { return controller.NewTransferController(uClient, wClient) },
		Req:      &controller.TransferBody{},
		Resp:     &controller.TransferData{},
	})
	httpcommon.RegisterPut(httpcommon.PutEndpointInfo{
		EP:       "api/v1/users/topups",
		NewCtlFn: func() httpcommon.Ctl { return controller.NewTopupController(wClient) },
		Req:      &controller.TopupBody{},
		Resp:     &controller.TopupData{},
	})
	httpcommon.RegisterPost(httpcommon.PostEndpointInfo{
		EP:       "api/v1/users/transactions/_query",
		NewCtlFn: func() httpcommon.Ctl { return controller.NewGetUserTransactionsController(uClient, wClient) },
		Req:      &controller.GetUserTransactionsBody{},
		Resp:     &controller.GetUserTransactionsData{},
	})
	httpcommon.RegisterPost(httpcommon.PostEndpointInfo{
		EP:       "api/v1/users/contacts/_query",
		NewCtlFn: func() httpcommon.Ctl { return controller.NewGetContactsController(uClient) },
		Req:      &controller.GetContactsBody{},
		Resp:     &controller.GetContactsData{},
	})
	httpcommon.RegisterGet(httpcommon.GetEndpointInfo{
		EP:       "api/v1/users/wallets",
		NewCtlFn: func() httpcommon.Ctl { return controller.NewGetUserWalletController(wClient) },
		Req:      nil,
		Resp:     &controller.GetUserWalletData{},
	})

	httpcommon.Run()
}

func main() {
	setupHTTPServer()
}
