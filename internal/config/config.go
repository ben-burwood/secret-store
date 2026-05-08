package config

import (
	"log"
	"os"
)

type Config struct {
	Port          string
	DBPath        string
	EncryptionKey string
	DashboardUser string
	DashboardPass string
	StaticDir     string
}

func Load() *Config {
	c := &Config{
		Port:          getenv("PORT", "8080"),
		DBPath:        getenv("DB_PATH", "/data/secret-store.db"),
		EncryptionKey: os.Getenv("ENCRYPTION_KEY"),
		DashboardUser: os.Getenv("DASHBOARD_USER"),
		DashboardPass: os.Getenv("DASHBOARD_PASSWORD"),
		StaticDir:     resolveStaticDir(os.Getenv("STATIC_DIR")),
	}
	if c.EncryptionKey == "" {
		log.Fatal("ENCRYPTION_KEY environment variable is required")
	}
	if c.DashboardUser == "" || c.DashboardPass == "" {
		log.Fatal("DASHBOARD_USER and DASHBOARD_PASSWORD environment variables are required")
	}
	return c
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func resolveStaticDir(override string) string {
	var candidates []string
	if override != "" {
		candidates = []string{override}
	} else {
		candidates = []string{"/frontend/dist", "frontend/dist"}
	}
	for _, p := range candidates {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return p
		}
	}
	log.Printf("warning: no static dir found (tried %v); SPA disabled", candidates)
	return ""
}
