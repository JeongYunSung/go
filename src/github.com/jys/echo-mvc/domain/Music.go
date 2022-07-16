package domain

import "github.com/jinzhu/gorm"

type Music struct {
	gorm.Model
	Title   string
	Artist  string
	Comment *string
}
