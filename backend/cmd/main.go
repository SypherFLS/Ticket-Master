package main

import (
	"net/http"
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
			"error", 
			err,
		)
	}

	config_env := ""
	switch os.Getenv("ENV") {
	case "local":
		config_env = "CONFIG_PATH_BACKEND"
	case "dev":
		config_env = "CONFIG_PATH_DOCKER"
	default:
		appLogger.Error(
			"wrong work env", 
			"error", 
			"wrong config_env",
		)
	}

	path := os.Getenv(config_env)
	cfg, err := config.InitConfig(path)

	if err != nil {
		appLogger.Error(
			"failed init config with error", 
			"error",
			err,
		)
		os.Exit(1)
	}

	db, err := postgres.InitDB(cfg)

	if err != nil {
		appLogger.Error(
			"failed init to db with error", 
			"error", 
			err,
		)
		os.Exit(1)
	}

	appLogger.Info(
		"server start",
		"port",
		cfg.Server.Port,
	)

	secret := os.Getenv("JWT_SECRET")
	JWTManager := auth.NewJWTManager([]byte(secret))

	repo := postgres.NewRepo(db)
	service := service.NewService(repo, JWTManager)
	handler := handlers.NewHandler(service, appLogger)
	router := router.NewRouter(handler, JWTManager, cfg)

	http.ListenAndServe(cfg.Server.Port, router)
}
