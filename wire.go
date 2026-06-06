//go:build wireinject
// +build wireinject

package main

import (
	"backend-crowdfunding/campaign"
	"backend-crowdfunding/config"
	"backend-crowdfunding/database"
	"backend-crowdfunding/handler"
	"backend-crowdfunding/storage"
	"backend-crowdfunding/transaction"
	"backend-crowdfunding/user"

	"github.com/google/wire"
)

func InitializeApp() (*App, error) {
	wire.Build(
		config.Load,
		database.NewConnection,
		storage.New,
		provideAuthService,
		providePaymentService,
		user.ProviderSet,
		campaign.ProviderSet,
		transaction.ProviderSet,
		handler.ProviderSet,
		newRouter,
		newApp,
	)

	return nil, nil
}
