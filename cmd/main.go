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

	_ "github.com/dharmasaputraa/cinema-api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/dharmasaputraa/cinema-api/internal/infrastructure/cache"
	"github.com/dharmasaputraa/cinema-api/internal/infrastructure/config"
	"github.com/dharmasaputraa/cinema-api/internal/infrastructure/db"

	cinemaHTTP "github.com/dharmasaputraa/cinema-api/internal/cinema/delivery/http"
	cinemaRepo "github.com/dharmasaputraa/cinema-api/internal/cinema/repository/postgres"
	cinemaUC "github.com/dharmasaputraa/cinema-api/internal/cinema/usecase"

	"github.com/dharmasaputraa/cinema-api/pkg/middleware"
)

// @title Cinema API
// @version 1.0
// @description This is the API documentation for the Cinema management system.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// ======================
	// TEMP LOGGER
	// ======================
	tmpLog, _ := zap.NewProduction()

	// ======================
	// CONFIG
	// ======================
	cfg, err := config.Load()
	if err != nil {
		tmpLog.Fatal("invalid configuration", zap.Error(err))
	}

	// ======================
	// REAL LOGGER
	// ======================
	log := initLogger(cfg)

	zap.ReplaceGlobals(log)

	defer func() {
		if err := log.Sync(); err != nil {
			fmt.Println("failed to sync logger:", err)
		}
	}()

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
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)

	if cfg.App.AutoMigrate {
		log.Info("running migrations...")

		if err := db.RunMigrations(dsn, log); err != nil {
			log.Fatal("migration failed", zap.Error(err))
		}

		log.Info("migration completed")
	} else {
		log.Info("auto migration disabled")
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
		middleware.RequestID(),
		middleware.CORS(cfg.App.CORSOrigins),
		middleware.Logger(log),
		middleware.ErrorHandler(log),
	)

	// ======================
	// SWAGGER ROUTE
	// ======================
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
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

func initLogger(cfg *config.Config) *zap.Logger {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)

		log, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		return log
	}

	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return log
}
