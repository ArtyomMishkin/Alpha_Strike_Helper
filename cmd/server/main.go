package main

import (
	"Alpha_Strike_Helper/internal/handler"
	"Alpha_Strike_Helper/internal/middleware"
	"Alpha_Strike_Helper/internal/repository"
	"Alpha_Strike_Helper/internal/service"
	"Alpha_Strike_Helper/pkg/config"
	"Alpha_Strike_Helper/pkg/database"
	"Alpha_Strike_Helper/pkg/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("🚀 Starting Alpha Strike Helper...")

	// Загрузка конфигурации
	cfg := config.Load()
	log.Printf("📋 Config loaded: DB=%s:%s, User=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser)

	// Подключение к базе данных
	db, err := database.NewPostgresDB(database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Миграции
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}

	// Инициализация репозиториев
	cardRepo := repository.NewCardRepository(db)
	userRepo := repository.NewUserRepository(db)
	collectionRepo := repository.NewCollectionRepository(db)
	lanceRepo := repository.NewLanceRepository(db)
	starRepo := repository.NewStarRepository(db)
	log.Println("✓ Repositories initialized")

	// Инициализация сервисов
	cardService := service.NewCardService(cardRepo)
	authService := service.NewAuthService(userRepo, utils.NewJWTService(cfg.JWTSecret))
	collectionService := service.NewCollectionService(collectionRepo, cardRepo)
	log.Println("✓ Services initialized")

	// Инициализация хендлеров
	cardHandler := handler.NewCardHandler(cardService)
	authHandler := handler.NewAuthHandler(authService)
	collectionHandler := handler.NewCollectionHandler(collectionService)
	lanceHandler := handler.NewLanceHandler(lanceRepo)
	starHandler := handler.NewStarHandler(starRepo)
	log.Println("✓ Handlers initialized")

	// Настройка роутера Gin
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s - Status: %d, Latency: %v\n",
			param.TimeStamp.Format("15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	}))

	// Статические файлы и шаблоны
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	// ==================== HTML СТРАНИЦЫ ====================
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{})
	})

	router.GET("/collections", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{})
	})
	router.GET("/lances", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{})
	})
	router.GET("/stars", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{})
	})
	router.GET("/cards", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{})
	})
	router.GET("/api/v1/lances/:id/export", lanceHandler.ExportLance)
	router.GET("/api/v1/stars/:id/export", starHandler.ExportStar)

	// API для загрузки HTML-компонентов
	router.GET("/templates/:page", func(c *gin.Context) {
		page := c.Param("page")
		allowed := map[string]bool{
			"collections": true,
			"lances":      true,
			"stars":       true,
			"cards":       true,
		}
		if !allowed[page] {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			return
		}
		c.HTML(http.StatusOK, page+".html", gin.H{})
	})

	// ==================== API МАРШРУТЫ ====================
	v1 := router.Group("/api/v1")
	router.PUT("/api/v1/lances/:id", lanceHandler.UpdateLance)
	router.PUT("/api/v1/stars/:id", starHandler.UpdateStar)
	// Публичные маршруты аутентификации
	auth := v1.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Публичные маршруты карточек (чтение + пагинация)
	cards := v1.Group("/cards")
	cards.GET("", cardHandler.List) // ?page=1&pagesize=50
	cards.GET("/:id", cardHandler.Get)
	cards.GET("/search", cardHandler.Search)

	// ✅ ЛЭНСЫ И ЗВЁЗДЫ — ПУБЛИЧНЫЕ (БЕЗ AUTH!)
	lanceGroup := v1.Group("/lances")
	lanceGroup.GET("", lanceHandler.GetLances)
	lanceGroup.POST("", lanceHandler.CreateLance)
	lanceGroup.GET("/:id", lanceHandler.GetLance)
	lanceGroup.DELETE("/:id", lanceHandler.DeleteLance)
	lanceGroup.POST("/:id/cards/:cardId", lanceHandler.AddCardToLance)
	lanceGroup.DELETE("/:id/cards/:cardId", lanceHandler.RemoveCardFromLance)

	starGroup := v1.Group("/stars")
	starGroup.GET("", starHandler.GetStars)
	starGroup.POST("", starHandler.CreateStar)
	starGroup.GET("/:id", starHandler.GetStar)
	starGroup.DELETE("/:id", starHandler.DeleteStar)
	starGroup.POST("/:id/cards/:cardId", starHandler.AddCardToStar)
	starGroup.DELETE("/:id/cards/:cardId", starHandler.RemoveCardFromStar)

	// ✅ ЗАЩИЩЁННЫЕ МАРШРУТЫ (только коллекции + админ)
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(authService))

	// Коллекции пользователя
	collections := protected.Group("/collections")
	collections.GET("", collectionHandler.List)
	collections.POST("", collectionHandler.Create)
	collections.GET("/:id", collectionHandler.Get)
	collections.PUT("/:id", collectionHandler.Update)
	collections.DELETE("/:id", collectionHandler.Delete)
	collections.POST("/:id/cards/:cardId", collectionHandler.AddCard)
	collections.DELETE("/:id/cards/:cardId", collectionHandler.RemoveCard)

	// Административные маршруты для карточек
	adminCards := protected.Group("/admin/cards")
	adminCards.POST("", cardHandler.Create)
	adminCards.PUT("/:id", cardHandler.Update)
	adminCards.DELETE("/:id", cardHandler.Delete)

	// Запуск сервера
	log.Printf("✓ Server starting on port %s", cfg.ServerPort)
	log.Printf("📖 API available at http://localhost:%s/api/v1", cfg.ServerPort)
	log.Printf("🌐 Web UI available at http://localhost:%s/", cfg.ServerPort)
	log.Printf("✅ ЛЭНСЫ/ЗВЁЗДЫ ПУБЛИЧНЫЕ — работают без токена!")

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
