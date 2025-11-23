package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"reward-system/internal/api"
	"reward-system/internal/config"
	"reward-system/internal/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &config.Config{
		DBDSN:      os.Getenv("DB_DSN"),
		WechatToken: os.Getenv("WECHAT_TOKEN"),
		APIToken:   os.Getenv("API_TOKEN"),
		MCPURL:     os.Getenv("MCP_SERVER_URL"),
		Port:       os.Getenv("PORT"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	database, err := db.InitDB(cfg.DBDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Migrate(database); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	router := api.SetupRouter(database, cfg)
	
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}