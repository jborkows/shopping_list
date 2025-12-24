package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shopping/internal/domain/admin"
	"shopping/internal/domain/products"
	"shopping/internal/infrastructure/config"
	"shopping/internal/infrastructure/logging"
	"shopping/internal/infrastructure/oidc"
	"shopping/internal/infrastructure/persistence/sqlite"
	"shopping/internal/migrator"
	"shopping/internal/web"
)

func main() {
	logger := logging.New()
	slog.SetDefault(logger)

	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	conn, err := sqlite.Open(cfg.DBDSN)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer conn.Close()

	if err := migrator.Up(conn); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	repo := sqlite.NewRepo(conn)
	var productsQueries products.Queries = repo
	productsService := products.NewService(repo)
	var adminMaintenance admin.Maintenance = repo

	authenticator, err := oidc.New(cfg)
	if err != nil {
		log.Fatalf("auth: %v", err)
	}

	srv := web.NewServer(cfg, productsQueries, productsService, adminMaintenance, authenticator)

	httpServer := &http.Server{
		Addr:              cfg.Addr,
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", cfg.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(ctx)
}
