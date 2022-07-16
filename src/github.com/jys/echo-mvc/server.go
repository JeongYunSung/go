package main

import (
	"github.com/joho/godotenv"
	"mvc/config"
)

func main() {
	_ = godotenv.Load()
	config.Init()
}
