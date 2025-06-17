package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"url-shortener_/internal/config"
)

func main() {
	loadEnv()

	cfg := config.LoadConfig()
	fmt.Println(cfg)
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
}
