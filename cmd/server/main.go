package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"samplegadai-backend/internal/config"
	"samplegadai-backend/internal/database"
	"samplegadai-backend/internal/handler"
	"samplegadai-backend/internal/marketprice"
	"samplegadai-backend/internal/notification"
	"samplegadai-backend/internal/repository"
	"samplegadai-backend/internal/route"
	"samplegadai-backend/internal/service"
	"samplegadai-backend/internal/storage"
)

func main() {
	log.Println("Starting Pawn Management System Backend (Sistem Informasi Gadai)...")

	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Init Database & Migrations
	db := database.InitDB(cfg)
	database.RunMigrations(db, "migrations/000001_init_schema.up.sql")

	// 3. Init MinIO Storage & Market Price Engine & Notifications
	minioStorage := storage.NewMinioStorage(cfg)
	marketPriceSvc := marketprice.NewMarketPriceService(db)
	notifSvc := notification.NewNotificationService()

	// 4. Init Layers: Repository -> Service -> Handler
	repo := repository.NewRepository(db)
	svc := service.NewService(repo, cfg, minioStorage, marketPriceSvc, notifSvc)
	h := handler.NewHandler(svc)

	// 5. Setup Router & Graceful Shutdown Server
	router := route.SetupRouter(cfg, h)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server running on port %s (%s mode)", cfg.Port, cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup error: %v", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
