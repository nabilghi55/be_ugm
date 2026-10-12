package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	AppEnv              string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBSchema            string
	JWTSecret           string
	JWTRefreshSecret    string
	JWTAccessExpiryMin  int
	JWTRefreshExpiryDay int
	MinioEndpoint       string
	MinioAccessKey      string
	MinioSecretKey      string
	MinioBucketName     string
	MinioUseSSL         bool
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:                getEnv("PORT", "8080"),
		AppEnv:              getEnv("APP_ENV", "development"),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBUser:              getEnv("DB_USER", "postgres"),
		DBPassword:          getEnv("DB_PASSWORD", "postgres"),
		DBName:              getEnv("DB_NAME", "samplegadai"),
		DBSchema:            getEnv("DB_SCHEMA", "public"),
		JWTSecret:           getEnv("JWT_SECRET", "super-secret-access-key-2026"),
		JWTRefreshSecret:    getEnv("JWT_REFRESH_SECRET", "super-secret-refresh-key-2026"),
		JWTAccessExpiryMin:  15,
		JWTRefreshExpiryDay: 7,
		MinioEndpoint:       getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey:      getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey:      getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinioBucketName:     getEnv("MINIO_BUCKET_NAME", "samplegadai-bucket"),
		MinioUseSSL:         false,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
