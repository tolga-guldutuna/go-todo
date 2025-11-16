package main

import (
	"context"
	"fmt"
	"go-todo/internal/config"
	"go-todo/internal/db"
	"log"
	"time"
)

func main() {
	// context = Java'daki RequestContext/Timeout mantığına benzer.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Config yükle (.env + env var)
	cfg := config.Load()

	// DB pool oluştur
	pool, err := db.NewPool(ctx, cfg.DB_DSN)
	if err != nil {
		log.Fatalf("DB bağlantısı kurulamadı: %v", err)
	}
	defer pool.Close()

	// Basit health check: SELECT 1
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("DB ping FAILED: %v", err)
	}

	fmt.Println("✅ PostgreSQL bağlantısı başarılı!")
}
