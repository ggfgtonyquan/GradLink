package http

import "gorm.io/gorm"

type Deps struct {
	DB *gorm.DB
}
