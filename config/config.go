package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Port string
}

type Config struct {
	Env EnvConfig
}

func init() {
	const pathToEnv = "../.env"

	if err := godotenv.Load(pathToEnv); err != nil {
		log.Fatalf("Cannot loading .env file at path: %s\n", pathToEnv)
	}
}

func New() *Config {
	return &Config{
		Env: EnvConfig{
			Port: getEnvByKey("PORT"),
		},
	}
}

func getEnvByKey(key string) string {
	v, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("No field %s in .env file\n", key)
	}

	return v
}
