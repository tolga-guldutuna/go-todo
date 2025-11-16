package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tolga-guldutuna/go-todo/internal/config"
	"github.com/tolga-guldutuna/go-todo/internal/db"
	"github.com/tolga-guldutuna/go-todo/internal/todo"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := config.Load()

	pool, err := db.NewPool(ctx, cfg.DB_DSN)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	defer pool.Close()

	// katmanlar
	repo := todo.NewRepository(pool)
	svc := todo.NewService(repo)
	handler := todo.NewHandler(svc)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	fmt.Println("✅ API listening on", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
