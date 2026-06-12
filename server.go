package main

import (
	"log/slog"
	"net/http"

	"backend-crowdfunding/auth"
	"backend-crowdfunding/config"
	"backend-crowdfunding/handler"
	"backend-crowdfunding/helper"
	"backend-crowdfunding/payment"
	"backend-crowdfunding/user"

	_ "backend-crowdfunding/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	Router *gin.Engine
	Config *config.Config
}

func newApp(router *gin.Engine, cfg *config.Config) *App {
	return &App{Router: router, Config: cfg}
}

func provideAuthService(cfg *config.Config) auth.Service {
	return auth.NewService(cfg.JWT.SecretKey)
}

func providePaymentService(cfg *config.Config) payment.Service {
	return payment.NewService(cfg.Midtrans.ServerKey, cfg.Midtrans.ClientKey, cfg.Midtrans.Available())
}

func newRouter(
	cfg *config.Config,
	userHandler *handler.UserHandler,
	campaignHandler *handler.CampaignHandler,
	transactionHandler *handler.TransactionHandler,
	authService auth.Service,
	userService user.Service,
) *gin.Engine {
	logger := slog.Default()

	router := gin.New()
	// We don't sit behind a known proxy in development; opt out of Gin's
	// trust-all-proxies default so ClientIP() can't be spoofed.
	_ = router.SetTrustedProxies(nil)
	router.Use(handler.RequestID(), handler.RequestLogger(logger), gin.Recovery())
	router.Use(cors.Default())
	router.Static("/images", "./images")

	authRequired := handler.AuthMiddleware(authService, userService)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Public client configuration. The web app reads this to decide whether to
	// run the live payment flow or show its "contact support" demo modal, so
	// MIDTRANS_ENABLED on the server stays the single source of truth.
	api.GET("/config", func(c *gin.Context) {
		data := gin.H{
			"payments_enabled": cfg.Midtrans.Available(),
			"support_email":    cfg.SupportEmail,
		}
		c.JSON(http.StatusOK, helper.APIResponse("Client configuration", http.StatusOK, "success", data))
	})

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authRequired, userHandler.UploadAvatar)
	api.GET("/users/fetch", authRequired, userHandler.FetchUser)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authRequired, campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authRequired, campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authRequired, campaignHandler.UploadCampaignImage)

	api.GET("/campaigns/:id/transactions", authRequired, transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authRequired, transactionHandler.GetUserTransactions)
	api.POST("/transactions", authRequired, transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	return router
}
