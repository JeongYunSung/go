package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"mvc/router"
)

func main() {
	e := echo.New()

	godotenv.Load()
	router.Route(e)

	e.Logger.Fatal(e.Start(":8080"))
}
