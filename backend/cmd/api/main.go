package main

import (
	"database/sql"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vietlx426/tripsearch/db/sqlc"
	"github.com/vietlx426/tripsearch/internal/auth"
	"github.com/vietlx426/tripsearch/internal/user"
	"github.com/vietlx426/tripsearch/pkg/cache"
	"github.com/vietlx426/tripsearch/pkg/logger"
	"github.com/vietlx426/tripsearch/pkg/middleware"
)

func main() {
	env := getEnv("APP_ENV", "development")
	logger.Init(env)
	log := logger.Get()

	dbURL := getEnv("DATABASE_URL", "postgres://tripsearch:tripsearch@localhost:5432/tripsearch?sslmode=disable")
	sqlDB, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open database")
	}
	defer sqlDB.Close()

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	redisURL := getEnv("REDIS_URL", "redis://localhost:6379")
	if err := cache.Init(redisURL); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to redis")
	}

	queries := db.New(sqlDB)

	userRepo    := user.NewRepository(queries)
	sessionRepo := auth.NewSessionRepository(queries)
	jwtSecret   := getEnv("JWT_SECRET", "change-me-in-production")
	authSvc     := auth.NewService(userRepo, sessionRepo, jwtSecret)
	authHandler := auth.NewHandler(authSvc)
	userSvc     := user.NewService(userRepo)
	userHandler := user.NewHandler(userSvc)

	r := gin.New()
	r.Use(
		middleware.RequestID(),
		middleware.Logger(log),
		middleware.Recovery(log),
		middleware.RateLimit(cache.Get(), 100, time.Minute),
	)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/logout", authHandler.Logout)
			authRoutes.POST("/refresh", authHandler.Refresh)
		}

		userRoutes := v1.Group("/users")
		userRoutes.Use(middleware.JWT(jwtSecret))
		{
			userRoutes.GET("/me", userHandler.GetMe)
		}
	}

	addr := getEnv("PORT", ":8080")
	log.Info().Str("addr", addr).Msg("starting server")
	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
