package app

import (
	"context"
	"time"

	"wb_labs_l0_backend/internal/config"
	"wb_labs_l0_backend/internal/db"
	"wb_labs_l0_backend/internal/logger"
	postgresrepo "wb_labs_l0_backend/internal/repository/postgres"
	"wb_labs_l0_backend/internal/service"
	cachepkg "wb_labs_l0_backend/internal/service/cache"
	httptransport "wb_labs_l0_backend/internal/transport/http"
	natstransport "wb_labs_l0_backend/internal/transport/nats"

	"os"
	"os/signal"
	"syscall"
)

func Run() error {
	// ! load config
	cfg := config.Load()

	// ! logger
	log := logger.New(cfg.LogLevel)

	// ! connect to db
	pg, err := db.Connect(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal("db connect failed", err)
	}
	// ! repo
	repo := postgresrepo.New(pg, log)

	// ! cache
	cache := cachepkg.NewInMemoryCache()

	// ! service
	svc := service.NewOrderService(repo, cache, log)

	// ! restore cache from DB
	if err := svc.RestoreCache(context.Background()); err != nil {
		log.Warn("restore cache failed", err)
	}

	// ! NATS subscriber
	sc, sub, err := natstransport.NewSubscriber(cfg, svc, log)
	if err != nil {
		log.Fatal("nats subscriber init failed", err)
	}
	// ! start subscriber
	go func() {
		if err := sub.Start(); err != nil {
			log.Error("nats subscriber stopped with error", err)
		}
	}()

	// ! HTTP server
	server := httptransport.NewServer(cfg, svc, log)
	go func() {
		if err := server.Start(); err != nil {
			log.Fatal("http server error", err)
		}
	}()

	// ! graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutdown initiated")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	server.Stop(ctx)
	sub.Stop()
	sc.Close()
	pg.Close()

	log.Info("shutdown completed")
	return nil
}
