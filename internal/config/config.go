package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config = Java'daki AppConfig / Properties holder gibi.
type Config struct {
	DB_DSN string // PostgreSQL connection string
	Addr   string // HTTP adresi (ileride kullanacağız)
}

// Load, .env varsa okur, yoksa direkt env'den çeker.
func Load() Config {
	// .env dosyası varsa yükle (hata verse bile ignore ediyoruz).
	_ = godotenv.Load()

	dsn := os.Getenv("TODO_DB_DSN")
	if dsn == "" {
		log.Fatal("TODO_DB_DSN environment variable is required")
	}

	addr := os.Getenv("TODO_HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	return Config{
		DB_DSN: dsn,
		Addr:   addr,
	}
}
