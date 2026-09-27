package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Cfg struct {
	Addr        string
	BaseUrl     string
	DbPath      string
	CorsOrigins []string
	DatabaseURL string
}

func checkEnvVal(target, defaultVal string, finalVal *string) {
	result := os.Getenv(target)

	if result == "" {
		*finalVal = defaultVal
	} else {
		*finalVal = result
	}

}

func parseCorsOrigins(raw string) []string {
	defaults := []string{
		"http://localhost:8000",
		"http://localhost:5173",
		"https://url.amk.dev",
	}
	if strings.TrimSpace(raw) == "" {
		return defaults
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if o := strings.TrimSpace(p); o != "" {
			origins = append(origins, o)
		}
	}
	if len(origins) == 0 {
		return defaults
	}
	return origins
}

func Load() (*Cfg, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("loading .env: %w", err)
	}

	config := &Cfg{}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	config.Addr = addr
	checkEnvVal("BASE_URL", "http://localhost:8080", &config.BaseUrl)
	checkEnvVal("DB_PATH", "data/app.db", &config.DbPath)
	config.CorsOrigins = parseCorsOrigins(os.Getenv("CORS_ORIGINS"))
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	config.DatabaseURL = url

	return config, nil
}
