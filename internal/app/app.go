package app

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/qnocks/blacklist-user-service/internal/config"
	"github.com/qnocks/blacklist-user-service/internal/repository"
	"github.com/qnocks/blacklist-user-service/internal/service"
	"github.com/qnocks/blacklist-user-service/internal/transport/rest"
	"github.com/qnocks/blacklist-user-service/pkg/db/postgres"
	"github.com/qnocks/blacklist-user-service/pkg/httpserver"
	"github.com/qnocks/blacklist-user-service/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Blacklist user service API
// @version 1.0
// @description The purpose of an application is to store information about users who have been added to the blacklist

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

const timeout = 5 * time.Second

func Run() {
	cfg := config.GetConfig()

	db, err := postgres.NewPostgres(postgres.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		logger.Errorf("failed to connect to postgres: %s\n", err.Error())
		return
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handler := rest.NewHandler(services)
	srv := httpserver.NewServer(cfg.Server.Port, handler.Init())

	go func() {
		if err := srv.Run(); err != nil {
			logger.Errorf("error running http server: %s\n", err.Error())
			return
		}
	}()

	graceful()
	shutdown(srv, db)
}

func graceful() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func shutdown(srv *httpserver.Server, db *sqlx.DB) {
	ctx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
	defer shutdownCancel()

	if err := db.Close(); err != nil {
		logger.Errorf("failed to close postgres connection: %s\n", err.Error())
		return
	}

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("failed to stop server: %s\n", err.Error())
		return
	}
}
