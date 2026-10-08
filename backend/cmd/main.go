package main

import (
	"os"
	"tmaster/internal/api/handlers"
	"tmaster/internal/api/router"
	"tmaster/internal/auth"
	"tmaster/internal/config"
	"tmaster/internal/logger"
	"tmaster/internal/repository/postgres"
	"tmaster/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	appLogger := logger.New()

	err := godotenv.Load()
	if err != nil {
		appLogger.Error(
			"godotenv error:",
			"error", err,
		)
		os.Exit(1)
	}

	configEnv := ""
	env := os.Getenv("ENV")
	switch env {
	case "local":
		configEnv = "CONFIG_PATH_BACKEND"
	case "dev":
		configEnv = "CONFIG_PATH_DOCKER"
	default:
		appLogger.Error(
			"wrong work env",
			"error", "wrong config_env",
			"env", env,
		)
	}

	path := os.Getenv(configEnv)
	cfg, err := config.InitConfig(path)

	if err != nil {
		appLogger.Error(
			"failed init config with error",
			"error", err,
			"config path", path,
		)
		os.Exit(1)
	}

	db, err := postgres.InitDB(cfg)

	if err != nil {
		appLogger.Error(
			"failed init to db with error",
			"error", err,
		)
		os.Exit(1)
	}

	appLogger.Info(
		"db init and connected",
		"port", cfg.Server.Port,
	)

	secret := os.Getenv("JWT_SECRET")
	jwtManager := auth.NewJWTManager([]byte(secret))

	repo := postgres.NewRepo(db)
	service := service.NewService(repo, jwtManager)
	handler := handlers.NewHandler(service, appLogger)
	router := router.NewRouter(handler, jwtManager, cfg, appLogger)

	if err := router.Run(cfg.Server.Port); err != nil {
		appLogger.Error(
			"failed start service",
			"error", err,
		)
		os.Exit(1)
	}
}
