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

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/tolga-guldutuna/go-todo/internal/docs" // <-- BUNU UNUTMA
)

// @title           Go Todo API
// @version         1.0
// @description     Simple layered todo API in Go.
// @host            localhost:8080
// @BasePath        /
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := config.Load()

	pool, err := db.NewPool(ctx, cfg.DB_DSN)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	defer pool.Close()

	repo := todo.NewRepository(pool)
	svc := todo.NewService(repo)
	handler := todo.NewHandler(svc)

	// TEK mux kullanıyoruz
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Swagger route: DİKKAT → mux.Handle, path "/swagger/" olacak
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("✅ API listening on", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
