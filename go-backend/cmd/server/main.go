package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"go-backend/config"
	sqlcdb "go-backend/db/sqlc/generated"
	"go-backend/internal/handler"
	"go-backend/internal/logger"
	"go-backend/internal/middleware"
	"go-backend/internal/repository"
	"go-backend/internal/routes"
	"go-backend/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	// Load configuration.
	godotenv.Load()
	cfg := config.Load()

	// Initialize logger.
	logger.Init(cfg.Env)
	defer logger.Sync()
	log := logger.Get()

	log.Info("starting application",
		zap.String("env", cfg.Env),
		zap.String("port", cfg.ServerPort),
	)

	// Connect to PostgreSQL.
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DBConnString())
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	// Verify connection.
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("failed to ping database", zap.Error(err))
	}
	log.Info("connected to database successfully")

	// Run migration: create the users table if it doesn't exist.
	migrationSQL := `CREATE TABLE IF NOT EXISTS users (
		id   SERIAL PRIMARY KEY,
		name TEXT   NOT NULL,
		dob  DATE   NOT NULL
	);`
	if _, err := pool.Exec(ctx, migrationSQL); err != nil {
		log.Fatal("failed to run migration", zap.Error(err))
	}
	log.Info("database migration applied successfully")

	// Wire dependencies.
	queries := sqlcdb.New(pool)
	userRepo := repository.NewUserRepository(queries)
	userService := service.NewUserService(userRepo, log)
	userHandler := handler.NewUserHandler(userService, log)

	// Initialize Fiber.
	app := fiber.New(fiber.Config{
		AppName: "Ainyx Users API",
	})

	// Apply global middleware.
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger(log))

	// Health check endpoint.
	app.Get("/health", func(c fiber.Ctx) error {
		if err := pool.Ping(c.Context()); err != nil {
			log.Error("health check failed", zap.Error(err))
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy",
			})
		}
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	// Register routes.
	routes.Setup(app, userHandler)

	// Graceful shutdown.
	go func() {
		if err := app.Listen(cfg.ServerAddr()); err != nil {
			log.Fatal("server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Error("error during server shutdown", zap.Error(err))
	}
	log.Info("server stopped")
}
