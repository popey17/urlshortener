package config

import (
	"os"
)

type Cfg struct {
	Addr     string
	BaseUrl  string
	DbPath   string
	CodeSize int
}

func checkEnvVal(target, defaultVal string, finalVal *string) {
	result := os.Getenv(target)

	if result == "" {
		*finalVal = defaultVal
	} else {
		*finalVal = result
	}

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
	config.CodeSize = 7

	return config, nil
}
