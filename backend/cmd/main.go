package main

import (
	"log"
	"net/http"
	"os"
	"tmaster/internal/api/handlers"
	"tmaster/internal/api/router"
	"tmaster/internal/auth"
	"tmaster/internal/config"
	"tmaster/internal/repository/postgres"
	"tmaster/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("godotenv error:", err)
	}

	config_env := ""
	switch os.Getenv("ENV") {
	case "local":
		config_env = "CONFIG_PATH_BACKEND"
	case "dev":
		config_env = "CONFIG_PATH_DOCKER"
	default:
		log.Fatal("wrong work env")
	}

	path := os.Getenv(config_env)
	cfg, err := config.InitConfig(path)

	if err != nil {
		log.Fatalf("failed init config with error %v\n", err)
	}

	db, err := postgres.InitDB(cfg)

	if err != nil {
		log.Fatalf("failed init to db with error %v\n", err)
	}

	secret := os.Getenv("JWT_SECRET")
	JWTManager := auth.NewJWTManager([]byte(secret))

	repo := postgres.NewRepo(db)
	service := service.NewService(repo, JWTManager)
	handler := handlers.NewHandler(service)
	router := router.NewRouter(handler, JWTManager, cfg)

	http.ListenAndServe(cfg.Server.Port, router)
}
