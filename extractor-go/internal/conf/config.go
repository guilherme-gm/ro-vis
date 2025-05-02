package conf

import (
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	DbUrl string
}

var Config config

func Load() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	Config = config{
		DbUrl: os.Getenv("DB_URL"),
	}
}
