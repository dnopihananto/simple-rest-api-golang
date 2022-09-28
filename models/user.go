package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string    `gorm:"type:varchar(100)"`
	FullName string    `gorm:"type:varchar(100)"`
	Email    string    `gorm:"type:varchar(100);unique"`
	SocialId string    `gorm:"type:varchar(100)"`
	Provider string    `gorm:"type:varchar(100)"`
	Avatar   string    `gorm:"type:varchar(100)"`
	Role     bool      `gorm:"default:0`
	Article  []Article `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
