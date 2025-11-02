package config

import (
	"os"
)

type RDSConfig struct {
	HOST     string `json:"host"'`
	PORT     string `json:"port"`
	USERNAME string `json:"username"`
	PASSWORD string `json:"password"`
	DBNAME   string `json:"dbname"`
}

type AWSConfig struct {
	Region string `json:"region"`
}

type EnvConfig struct {
	AWS AWSConfig
	RDS RDSConfig
}

func LoadEnvConfig() *EnvConfig {
	// port, err := strconv.ItoA(os.Getenv("RDS_PORT"))
	//	if err != nil {
	// set default port if not config from env
	//	port = 5432
	//}
	rdsConfig := RDSConfig{
		HOST:     os.Getenv("RDS_HOST"),
		PORT:     os.Getenv("RDS_PORT"),
		USERNAME: os.Getenv("RDS_USERNAME"),
		PASSWORD: os.Getenv("RDS_PASSWORD"),
		DBNAME:   os.Getenv("RDS_DBNAME"),
	}
	awsConfig := AWSConfig{
		Region: os.Getenv("AWS_REGION"),
	}
	return &EnvConfig{
		AWS: awsConfig,
		RDS: rdsConfig,
	}
}
