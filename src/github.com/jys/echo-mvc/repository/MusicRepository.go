package repository

import (
	"mvc/config"
	"mvc/domain"
)

func Find(id uint) *domain.Music {
	db := config.Init()

	music := &domain.Music{}

	db.Where("ID = ?", id).Find(music)

	return music
}

func FindAll() *[]domain.Music {
	db := config.Init()

	musics := &[]domain.Music{}
	db.Find(musics)

	return musics
}

func Create(music domain.Music) {
	db := config.Init()
	db.Create(&music)
}
