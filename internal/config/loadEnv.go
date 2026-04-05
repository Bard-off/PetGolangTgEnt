package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type ConfigDB struct {
	DB_URL string `env:"DB_URL"`
	BOT_KEY string `env:"BOT_KEY"`
}

func LoadConfig() ConfigDB {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	bot_key := os.Getenv("BOT_KEY")
	log.Println("Получена информация с .env файла")
	db_url := os.Getenv("DB_URL")
	cfg := ConfigDB{DB_URL: db_url, BOT_KEY: bot_key}
	return cfg
}