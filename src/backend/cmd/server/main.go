package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"backend/internal/cache"
	"backend/internal/db"
	"backend/internal/handler"
	"backend/internal/nats"
)

func main() {
	dsn := getEnv("DB_DSN", "postgres://orders:orders@localhost:5432/orders?sslmode=disable")
	natsURL := getEnv("NATS_URL", "nats://localhost:4222")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// DB
	pg, err := db.New(dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.RunMigrations(dsn); err != nil {
		log.Fatal("migration error:", err)
	}

	// Кэш + восстановление
	c := cache.New()
	restored, err := pg.LoadAll(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, o := range restored {
		c.Set(o)
	}
	log.Printf("restored %d orders from DB", len(restored))

	// NATS
	go func() {
		if err := nats.Start(ctx, natsURL, pg, c); err != nil {
			log.Printf("nats error: %v", err)
		}
	}()

	// HTTP
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.New(c),
	}
	go func() {
		log.Println("HTTP server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(shutdownCtx)
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
