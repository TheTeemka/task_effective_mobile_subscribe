package main

import (
	"log/slog"

	_ "github.com/TheTeemka/task_effective_mobile_subscribe/docs"
	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/config"
	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/database"
	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/handlers"
	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/repo"
	"github.com/TheTeemka/task_effective_mobile_subscribe/internal/services"
	"github.com/TheTeemka/task_effective_mobile_subscribe/pkg/logging"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	slog.Info("Loading configuration")
	cfg := config.LoadConfig()
	slog.Info("Configuration loaded successfully", "logLevel", cfg.LogLevel, "port", cfg.Port)

	logging.SetSlog(cfg.LogLevel)

	db := database.NewPSQLConnection(cfg.PSQLSource)

	repo := repo.NewSubscriptionRepo(db)
	svc := services.NewSubscriptionService(repo)
	handler := handlers.NewSubscriptionHandler(svc)

	r := gin.Default()

	api := r.Group("/api/subscriptions")
	{
		api.GET("/", handler.ListSubscriptions)
		api.GET("/:id", handler.GetSubscription)
		api.GET("/sum", handler.GetSum)
		api.POST("/", handler.CreateSubscription)
		api.PATCH("/:id", handler.UpdateSubscription)
		api.DELETE("/:id", handler.DeleteSubscription)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	slog.Info("Starting server", "port", cfg.Port)

	r.Run(cfg.Port)
}
