package config

import (
	"os"
	"strings"
)

type Cfg struct {
	Addr        string
	BaseUrl     string
	DbPath      string
	CorsOrigins []string
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

func Load() (Cfg, error) {
	config := Cfg{}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	config.Addr = addr
	checkEnvVal("BASE_URL", "http://localhost:8080", &config.BaseUrl)
	checkEnvVal("DB_PATH", "data/app.db", &config.DbPath)
	config.CorsOrigins = parseCorsOrigins(os.Getenv("CORS_ORIGINS"))

	return config, nil
}
