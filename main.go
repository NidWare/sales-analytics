package main

import (
	"log"
	"os"
	"path/filepath"
	"sales-count/bot"
	"sales-count/config"
	"sales-count/database"
)

func main() {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal("Error getting executable path:", err)
	}

	execDir := filepath.Dir(execPath)
	configPath := filepath.Join(execDir, "config.yml")

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	db, err := database.ConnectDatabase("database.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	bot.Start(cfg, db)
}
