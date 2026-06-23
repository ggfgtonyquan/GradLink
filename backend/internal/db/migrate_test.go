package db

import (
	"testing"

	"gorm.io/gorm"
)

func TestRequiredTablesExist(t *testing.T) {
	database := testDB(t)

	want := []string{
		"users",
		"student_profiles",
		"posts",
		"tags",
		"post_tags",
		"comments",
		"likes",
		"favorites",
		"reports",
		"moderation_keywords",
		"moderation_queue",
		"sessions",
	}

	for _, table := range want {
		if !tableExists(t, database, table) {
			t.Fatalf("missing table %s", table)
		}
	}
}

func testDB(t *testing.T) *gorm.DB {
	t.Helper()

	database, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := Migrate(database); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	return database
}

func tableExists(t *testing.T, database *gorm.DB, name string) bool {
	t.Helper()

	var count int
	if err := database.Raw(
		"SELECT count(*) FROM sqlite_master WHERE type = ? AND name = ?",
		"table",
		name,
	).Scan(&count).Error; err != nil {
		t.Fatalf("check table %s: %v", name, err)
	}

	return count == 1
}
