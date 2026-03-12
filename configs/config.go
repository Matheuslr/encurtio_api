package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Cassandra CassandraConfig
	API       APIConfig
}

type CassandraConfig struct {
	Hosts    []string
	Keyspace string
}

type APIConfig struct {
	Port string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system env")
	}

	return &Config{
		Cassandra: CassandraConfig{
			Hosts:    []string{getEnv("CASSANDRA_HOST", "127.0.0.1")},
			Keyspace: getEnv("CASSANDRA_KEYSPACE", "encurtio"),
		},
		API: APIConfig{
			Port: getEnv("APP_PORT", "8080"),
		},
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
