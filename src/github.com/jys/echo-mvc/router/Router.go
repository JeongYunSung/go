package router

import (
	"github.com/labstack/echo/v4"
	"mvc/controller"
)

func Route(e *echo.Echo) {
	e.GET("/api/musics", MusicController.Index)
	e.GET("/api/musics/:id", MusicController.Show)
	e.POST("/api/musics", MusicController.Create)
}
