package rds

import (
	"fmt"
	"log"

	"github.com/phankieuphu/go-aws-example/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewRDSClient(cfg *config.EnvConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", cfg.RDS.HOST, cfg.RDS.PORT, cfg.RDS.USERNAME, cfg.RDS.PASSWORD, cfg.RDS.DBNAME)
	//fmt.Println("dsn", dsn)
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open conneciton: %v", err)
	}

	fmt.Println("Connected to RDS")
	return client
}
