package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"teststore/backend/internal/auth"
	"teststore/backend/internal/config"
	"teststore/backend/internal/db"
	"teststore/backend/internal/handlers"
	"teststore/backend/internal/models"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("gagal konek database: %v", err)
	}
	defer pool.Close()

	userRepo := models.NewUserRepo(pool)
	assetRepo := models.NewAssetRepo(pool)
	productRepo := models.NewProductRepo(pool)
	dressStyleRepo := models.NewDressStyleRepo(pool)
	tokens := auth.NewTokenManager(cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(userRepo, tokens)
	maxUploadBytes := cfg.MaxUploadMB * 1024 * 1024
	assetHandler := handlers.NewAssetHandler(assetRepo, cfg.UploadDir, maxUploadBytes)
	catalogHandler := handlers.NewCatalogHandler(productRepo, dressStyleRepo)

	router := gin.Default()
	router.MaxMultipartMemory = maxUploadBytes

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.AllowOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Static("/uploads", cfg.UploadDir)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.GET("/me", handlers.RequireAuth(tokens), authHandler.Me)
		}

		assetGroup := api.Group("/assets")
		{
			assetGroup.GET("", assetHandler.List)
			assetGroup.GET("/:key", assetHandler.Get)
			assetGroup.POST("/upload", handlers.RequireAuth(tokens), assetHandler.Upload)
			assetGroup.DELETE("/:key", handlers.RequireAuth(tokens), assetHandler.Delete)
		}

		api.GET("/products", catalogHandler.ListProducts)
		api.GET("/dress-styles", catalogHandler.ListDressStyles)
	}

	log.Printf("TestStore API listening on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
