package main

import (
	"fmt"
	"log"

	"github.com/phankieuphu/go-aws-example/internal/domain/handlers"
	"github.com/phankieuphu/go-aws-example/internal/domain/services"

	"github.com/phankieuphu/go-aws-example/config"
	"github.com/phankieuphu/go-aws-example/internal/adapters/repositories"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	cfg := config.LoadEnvConfig()
	//db := rds.NewRDSClient(cfg)
	//	db.AutoMigrate(&entities.User{})
	repository := repositories.NewUserRepository(cfg)
	fmt.Println("Repository: ")
	userService := services.NewUserServices(*repository)
	fmt.Println("userService: ")

	r := handlers.SetupRouters(userService)

	// return
	addr := fmt.Sprintf(":%s", "8080")
	log.Printf("Server is running at %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to run server", err)
	}
}
