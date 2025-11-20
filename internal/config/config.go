// Package config provides configuration loading and management for the application.
package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct{
	Env string `yaml:"env" env-defaut:"dev"`
	HTTPServer `yaml:"http_server"`
	DataBase `yaml:"database"`
}

type HTTPServer struct{
	Address string `yaml:"address" env-defaut:"0.0.0.0:8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DataBase struct{
	Host string `yaml:"host" env-default:"localhost"`
	Port string `yaml:"port" env-default:"5432"`
	User string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	Dbname string `yaml:"dbname" env-default:"transactions"`
	Sslmode string `yaml:"sslmode" env-default:"disable"`
}

func MustLoad() *Config  {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == ""{
		log.Fatal("CONFIG_PATH env is not set")
	}

	if _, err :=os.Stat(configPath); err != nil{
		log.Fatalf("err while oppening config file: %s", err)
	}

	var cfg Config
	
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("err while readig config, %s", err)
	}

	return &cfg
	
}