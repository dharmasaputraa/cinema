package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/dharmasaputraa/cinema-api/internal/infrastructure/cache"
	"github.com/dharmasaputraa/cinema-api/internal/infrastructure/config"
	"github.com/dharmasaputraa/cinema-api/internal/infrastructure/db"

	cinemaHTTP "github.com/dharmasaputraa/cinema-api/internal/cinema/delivery/http"
	cinemaRepo "github.com/dharmasaputraa/cinema-api/internal/cinema/repository/postgres"
	cinemaUC "github.com/dharmasaputraa/cinema-api/internal/cinema/usecase"

	"github.com/dharmasaputraa/cinema-api/pkg/middleware"
)

func main() {
	// ======================
	// LOGGER
	// ======================
	var log *zap.Logger
	log, _ = zap.NewProduction()

	// ======================
	// CONFIG
	// ======================
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("invalid configuration", zap.Error(err))
	}

	if cfg.App.Env == "production" {
		log, _ = zap.NewProduction()
		gin.SetMode(gin.ReleaseMode)
	} else {
		log, _ = zap.NewDevelopment()
	}
	defer log.Sync()

	log.Info("config loaded", zap.String("env", cfg.App.Env))

	// ======================
	// DATABASE
	// ======================
	gormDB, err := db.NewPostgres(cfg, log)
	if err != nil {
		log.Fatal("failed to connect db", zap.Error(err))
	}

	// ======================
	// MIGRATION
	// ======================
	dsn := fmt.Sprintf("%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.SSLMode)

	if err := db.RunMigrations(dsn, log); err != nil {
		log.Fatal("migration failed", zap.Error(err))
	}

	// ======================
	// REDIS
	// ======================
	_, err = cache.NewRedis(cfg, log)
	if err != nil {
		log.Fatal("failed to connect redis", zap.Error(err))
	}

	// ======================
	// ROUTER
	// ======================
	r := gin.New()

	// Middleware stack (ORDER PENTING)
	r.Use(
		middleware.CORS(cfg.App.CORSOrigins),
		middleware.Logger(log),
		middleware.ErrorHandler(log),
	)

	// ======================
	// HEALTH CHECK
	// ======================
	r.GET("/health", func(c *gin.Context) {
		requestID, _ := c.Get("request_id")

		c.JSON(http.StatusOK, gin.H{
			"status":     "ok",
			"request_id": requestID,
		})
	})

	api := r.Group("/api/v1")

	// ======================
	// WIRING (CINEMA)
	// ======================
	cinemaRepository := cinemaRepo.NewCinemaRepository(gormDB)
	screenRepository := cinemaRepo.NewScreenRepository(gormDB)
	seatRepository := cinemaRepo.NewSeatRepository(gormDB)

	cinemaUsecase := cinemaUC.NewCinemaUsecase(
		cinemaRepository,
		screenRepository,
		seatRepository,
		gormDB,
	)

	cinemaHandler := cinemaHTTP.NewCinemaHandler(cinemaUsecase)
	cinemaHTTP.RegisterRoutes(api, cinemaHandler)

	// ======================
	// SERVER
	// ======================
	srv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: r,
	}

	go func() {
		log.Info("server starting", zap.String("port", cfg.App.Port))

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	// ======================
	// GRACEFUL SHUTDOWN
	// ======================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}

	log.Info("server exited properly")
}
