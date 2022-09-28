package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title  string `gorm:"type:varchar(100)"`
	Slug   string `gorm:"type:varchar(100);unique"`
	Desc   string `gorm:"type:text"`
	Tag    string
	UserID uint
}
