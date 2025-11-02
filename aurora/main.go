package main

import (
	"context"
	"ex-aurora/config"
	"ex-aurora/internal/adapters/repositories"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	cfg := config.LoadEnvConfig()
	//db := rds.NewRDSClient(cfg)
	//	db.AutoMigrate(&entities.User{})
	println(cfg.RDS.HOST)
	repository := repositories.NewUserRepository(cfg)
	users, err := repository.GetUser(context.Background())
	if err != nil {
	}
	fmt.Println(users)
	return
}
