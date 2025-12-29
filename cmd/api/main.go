package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"koalbot_api/internal/config"
	"koalbot_api/internal/db"
	"koalbot_api/internal/handler"
	"koalbot_api/internal/migrate"
	"koalbot_api/internal/repository"
	"koalbot_api/internal/router"
	"koalbot_api/internal/seed"
	"koalbot_api/internal/service"
	"koalbot_api/internal/stockity"
)

func main() {
	cfg := config.Load()

	database, err := db.Open(cfg.DBDSN, cfg.DBMaxOpenConns, cfg.DBMaxIdleConns, cfg.DBConnMaxLifetime)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := migrate.Run(database, cfg.MigrationPath); err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(database)
	masterPenggunaRepo := repository.NewMasterPenggunaRepository(database)
	penggunaDetailRepo := repository.NewPenggunaDetailRepository(database)
	tokenRepo := repository.NewTokenRepository(database)
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	masterPenggunaService := service.NewMasterPenggunaService(masterPenggunaRepo)
	tokenService := service.NewTokenService(cfg.JWTSecret, tokenRepo)
	stockityClient := stockity.NewClient(cfg.StockityBaseURL, 15*time.Second)
	authHandler := handler.NewAuthHandler(authService, tokenService)
	userHandler := handler.NewUserHandler(userService)
	masterPenggunaHandler := handler.NewMasterPenggunaHandler(masterPenggunaService)
	v1LoginHandler := handler.NewV1LoginHandler(stockityClient, masterPenggunaRepo, penggunaDetailRepo, cfg.JWTSecret, cfg.StockityBaseURL)
	streamHandler := handler.NewStreamHandler()
	healthHandler := handler.NewHealthHandler(database)

	if err := seed.Users(database); err != nil {
		log.Fatal(err)
	}

	r := router.New(authHandler, v1LoginHandler, userHandler, masterPenggunaHandler, streamHandler, healthHandler, router.Options{
		MaxBodyBytes:         cfg.MaxBodyBytes,
		RateLimitRPS:         cfg.RateLimitRPS,
		RateLimitBurst:       cfg.RateLimitBurst,
		LoginRateLimitRPS:    cfg.LoginRateLimitRPS,
		LoginRateLimitBurst:  cfg.LoginRateLimitBurst,
		TrustedProxies:       cfg.TrustedProxies,
		JWTSecret:            cfg.JWTSecret,
		CORSAllowedOrigins:   cfg.CORSAllowedOrigins,
		CORSAllowCredentials: cfg.CORSAllowCredentials,
	})

	httpSrv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	go func() {
		log.Printf("listening on :%s", cfg.Port)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
