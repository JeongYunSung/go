package router

import (
	"github.com/labstack/echo/v4"
	"mvc/controller"
)

func Route(e *echo.Echo) {
	e.GET("/api/music", MusicController.Index)
	e.GET("/api/music/:id", MusicController.Show)
	e.POST("/api/music", MusicController.Create)
	e.PUT("/api/music/:id", MusicController.Update)
	e.DELETE("/api/music/:id", MusicController.Delete)
}
