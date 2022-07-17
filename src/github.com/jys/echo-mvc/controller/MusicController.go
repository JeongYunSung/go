package MusicController

import (
	"github.com/labstack/echo/v4"
	"mvc/domain"
	"mvc/repository"
	"net/http"
)

func Index(c echo.Context) error {
	repository.FindAll()
	return c.String(http.StatusOK, "NoContent")
}

func Show(c echo.Context) error {

	return nil
}

func Create(c echo.Context) error {
	title := c.FormValue("title")
	artist := c.FormValue("artist")

	music := domain.Music{Title: title, Artist: artist}
	repository.Create(music)
	return c.NoContent(http.StatusNoContent)
}

func Update(c echo.Context) error {

	return nil
}

func Delete(c echo.Context) error {

	return nil
}
