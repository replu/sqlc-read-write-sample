package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/replu/sqlc-read-write-sample/internal/handler"
	"github.com/replu/sqlc-read-write-sample/internal/repository"
	"github.com/replu/sqlc-read-write-sample/internal/service"
	"github.com/replu/sqlc-read-write-sample/internal/util/database"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger = logger.With("app", "sqlc-read-write-sample")
	slog.SetDefault(logger)

	conf, err := LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbConn := database.Connect(fmt.Sprintf(database.DataSourceFormat,
		conf.PrimaryDBUsername,
		conf.PrimaryDBPassword,
		conf.PrimaryDBHost,
		conf.PrimaryDBPort,
		conf.DBDatabase,
	))

	dba := database.NewAccessor(dbConn)

	repo := repository.NewRepository(dba, dbConn)
	srv := service.NewService(dba, repo)

	mux := routing(handler.NewHandler(srv))
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 20 * time.Second,
	}

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Shutdown: %v", err)
		}
		if err := dbConn.Close(); err != nil {
			log.Fatalf("Primary DB Close: %v", err)
		}
	}()

	// Start server
	log.Printf("Server started on :%d\n", conf.ServerPort)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe: %v", err)
	} else {
		logger.Info("Server closed")
	}
}
