package MusicController

import (
	"github.com/labstack/echo/v4"
	"mvc/domain"
	"mvc/repository"
	"net/http"
	"strconv"
)

func Index(c echo.Context) error {
	musics := repository.FindAll()
	return c.JSON(http.StatusOK, musics)
}

func Show(c echo.Context) error {
	id := c.Param("id")
	uid, _ := strconv.ParseUint(id, 10, 32)
	music := repository.Find(uint(uid))
	return c.JSON(http.StatusOK, music)
}

func Create(c echo.Context) error {
	title := c.FormValue("title")
	artist := c.FormValue("artist")

	music := domain.Music{Title: title, Artist: artist}
	repository.Create(music)
	return c.NoContent(http.StatusNoContent)
}
