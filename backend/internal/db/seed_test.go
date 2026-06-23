package db

import (
	"testing"

	"gorm.io/gorm"
)

func TestSeedCreatesTagsAdminAndKeywords(t *testing.T) {
	database := testDB(t)

	if err := Seed(database); err != nil {
		t.Fatalf("seed test database: %v", err)
	}

	if countRows(t, database, "tags") != int64(len(InitialTags)) {
		t.Fatalf("expected %d seed tags", len(InitialTags))
	}
	if !adminUserExists(t, database) {
		t.Fatal("expected admin user seed")
	}
	if countRows(t, database, "moderation_keywords") == 0 {
		t.Fatal("expected moderation keyword seeds")
	}
}

func countRows(t *testing.T, database *gorm.DB, table string) int64 {
	t.Helper()

	var count int64
	if err := database.Raw("SELECT count(*) FROM " + table).Scan(&count).Error; err != nil {
		t.Fatalf("count %s: %v", table, err)
	}

	return count
}

func adminUserExists(t *testing.T, database *gorm.DB) bool {
	t.Helper()

	var count int64
	if err := database.Raw(
		"SELECT count(*) FROM users WHERE email = ? AND role = ? AND password_hash <> ?",
		"admin@gradlink.local",
		"admin",
		"password",
	).Scan(&count).Error; err != nil {
		t.Fatalf("check admin user seed: %v", err)
	}

	return count == 1
}
