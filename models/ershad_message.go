package models

import (
	"github.com/jinzhu/gorm"
)

type ErshadMessage struct {
	gorm.Model
	EncMessage string
}
