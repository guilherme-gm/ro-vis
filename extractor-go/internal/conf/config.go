package conf

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type migratorConfig struct {
	DbUrl string
}

var MigratorConfig migratorConfig

type config struct {
	DbUrl string
}

var Config config

type apiConfig struct {
	Port uint16
}

var ApiConfig apiConfig

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func loadConfig() {
	Config = config{
		DbUrl: os.Getenv("DB_URL"),
	}
}

func getIntEnv(name string) int {
	val, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		panic(err)
	}

	return val
}

func loadApiConfig() {
	ApiConfig = apiConfig{
		Port: uint16(getIntEnv("API_PORT")),
	}
}

func LoadMigrator() {
	loadEnv()
	loadConfig()
	MigratorConfig = migratorConfig{
		DbUrl: os.Getenv("MIGRATOR_DB_URL"),
	}
}

func LoadApi() {
	loadEnv()
	loadConfig()
	loadApiConfig()
}

func LoadExtractor() {
	loadEnv()
	loadConfig()
}
