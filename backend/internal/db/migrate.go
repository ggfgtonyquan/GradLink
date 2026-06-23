package db

import "gorm.io/gorm"

func Migrate(database *gorm.DB) error {
	return database.AutoMigrate(
		&User{},
		&StudentProfile{},
		&Post{},
		&Tag{},
		&PostTag{},
		&Comment{},
		&Like{},
		&Favorite{},
		&Report{},
		&ModerationKeyword{},
		&ModerationQueueItem{},
		&Session{},
	)
}
