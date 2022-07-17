package repository

import (
	"fmt"
	"mvc/config"
	"mvc/domain"
)

func FindAll() {
	db := config.Init()

	var musics []domain.Music
	db.Find(&musics)

	fmt.Println(musics)
}

func Create(music domain.Music) {
	db := config.Init()
	db.Create(&music)
}
