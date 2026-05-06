package main

import (
	"Alpha_Strike_Helper/internal/handler"
	"Alpha_Strike_Helper/internal/middleware"
	"Alpha_Strike_Helper/internal/repository"
	"Alpha_Strike_Helper/internal/service"
	"Alpha_Strike_Helper/pkg/config"
	"Alpha_Strike_Helper/pkg/database"
	"Alpha_Strike_Helper/pkg/utils"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresDB(database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  "disable",
	})
	if err != nil {
		panic(err)
	}
	if err := database.AutoMigrate(db); err != nil {
		panic(err)
	}

	cardRepo := repository.NewCardRepository(db)
	userRepo := repository.NewUserRepository(db)
	cardService := service.NewCardService(cardRepo)
	authService := service.NewAuthService(userRepo, utils.NewJWTService(cfg.JWTSecret))
	cardHandler := handler.NewCardHandler(cardService)

	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "cards-service"})
	})

	v1 := router.Group("/api/v1")
	{
		v1.GET("/cards", cardHandler.List)
		v1.GET("/cards/:id", cardHandler.Get)
		v1.GET("/cards/search", cardHandler.Search)

		adminCards := v1.Group("/admin/cards")
		adminCards.Use(middleware.AuthMiddleware(authService))
		adminCards.POST("", cardHandler.Create)
		adminCards.PUT("/:id", cardHandler.Update)
		adminCards.DELETE("/:id", cardHandler.Delete)

		v1.GET("/chassis-sources", func(c *gin.Context) {
			raw, err := os.ReadFile("static/data/chassis_sources.json")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read chassis_sources.json", "details": err.Error()})
				return
			}
			var payload any
			if err := json.Unmarshal(raw, &payload); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse chassis_sources.json", "details": err.Error()})
				return
			}
			c.JSON(http.StatusOK, payload)
		})
	}

	port := "8082"
	if envPort := os.Getenv("CARDS_SERVICE_PORT"); envPort != "" {
		port = envPort
	}

	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}
