package main 

import (
	"tmaster/internal/config"
	"github.com/joho/godotenv"
	"os"
	"log"
	"fmt"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("godotenv error:", err)
	}

	config_env := ""
	switch os.Getenv("ENV"){
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
	
	fmt.Println(cfg)
}