package main

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	Content string
	Options []Option
}

type Option struct {
	gorm.Model
	QuestionID uint
	Text       string
	Votes      int `gorm:"default:0"`
}
