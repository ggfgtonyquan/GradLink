package db

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
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

func TestSeedCreatesDisabledAdminWithoutUsablePasswordByDefault(t *testing.T) {
	database := testDB(t)

	if err := Seed(database); err != nil {
		t.Fatalf("seed test database: %v", err)
	}

	admin := loadAdminUser(t, database)
	if admin.Status != UserStatusDisabled {
		t.Fatalf("expected disabled default admin, got %q", admin.Status)
	}
	if admin.PasswordHash == "" || strings.Contains(admin.PasswordHash, "password") {
		t.Fatalf("expected unusable default admin hash, got %q", admin.PasswordHash)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte("password")); err == nil {
		t.Fatal("default admin password hash must not match a predictable password")
	}
}

func TestSeedWithOptionsActivatesAdminWithBcryptPassword(t *testing.T) {
	database := testDB(t)
	rawPassword := "local-only-admin-password"

	if err := SeedWithOptions(database, SeedOptions{AdminPassword: rawPassword}); err != nil {
		t.Fatalf("seed test database with options: %v", err)
	}

	admin := loadAdminUser(t, database)
	if admin.Status != UserStatusActive {
		t.Fatalf("expected active admin, got %q", admin.Status)
	}
	if admin.Role != UserRoleAdmin {
		t.Fatalf("expected admin role, got %q", admin.Role)
	}
	if strings.Contains(admin.PasswordHash, rawPassword) {
		t.Fatal("admin password hash contains the raw password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(rawPassword)); err != nil {
		t.Fatalf("admin bcrypt hash does not match provided password: %v", err)
	}
}

func TestSeedWithoutPasswordDoesNotOverwriteExistingAdminPassword(t *testing.T) {
	database := testDB(t)
	rawPassword := "local-only-admin-password"

	if err := SeedWithOptions(database, SeedOptions{AdminPassword: rawPassword}); err != nil {
		t.Fatalf("seed test database with options: %v", err)
	}
	originalAdmin := loadAdminUser(t, database)

	if err := Seed(database); err != nil {
		t.Fatalf("reseed test database without password: %v", err)
	}

	admin := loadAdminUser(t, database)
	if admin.PasswordHash != originalAdmin.PasswordHash {
		t.Fatal("default seed overwrote existing admin password hash")
	}
	if admin.Status != UserStatusActive {
		t.Fatalf("default seed changed existing admin status to %q", admin.Status)
	}
}

func TestSeedReconcilesTagsAndModerationKeywords(t *testing.T) {
	database := testDB(t)

	if err := Seed(database); err != nil {
		t.Fatalf("seed test database: %v", err)
	}
	if err := database.Model(&Tag{}).
		Where("slug = ?", InitialTags[0].Slug).
		Updates(map[string]interface{}{"name": "旧标签", "display_order": 999}).Error; err != nil {
		t.Fatalf("mutate seed tag: %v", err)
	}
	if err := database.Model(&ModerationKeyword{}).
		Where("keyword = ?", initialModerationKeywords[0].Keyword).
		Updates(map[string]interface{}{"risk_level": "low", "action": StatusPublished, "enabled": false}).Error; err != nil {
		t.Fatalf("mutate moderation keyword: %v", err)
	}

	if err := Seed(database); err != nil {
		t.Fatalf("reseed test database: %v", err)
	}

	var tag Tag
	if err := database.Where("slug = ?", InitialTags[0].Slug).First(&tag).Error; err != nil {
		t.Fatalf("load reconciled tag: %v", err)
	}
	if tag.Name != InitialTags[0].Name || tag.DisplayOrder != InitialTags[0].DisplayOrder {
		t.Fatalf("tag was not reconciled: %+v", tag)
	}

	var keyword ModerationKeyword
	if err := database.Where("keyword = ?", initialModerationKeywords[0].Keyword).First(&keyword).Error; err != nil {
		t.Fatalf("load reconciled moderation keyword: %v", err)
	}
	if keyword.RiskLevel != initialModerationKeywords[0].RiskLevel ||
		keyword.Action != initialModerationKeywords[0].Action ||
		keyword.Enabled != initialModerationKeywords[0].Enabled {
		t.Fatalf("moderation keyword was not reconciled: %+v", keyword)
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

func loadAdminUser(t *testing.T, database *gorm.DB) User {
	t.Helper()

	var admin User
	if err := database.Where("email = ?", "admin@gradlink.local").First(&admin).Error; err != nil {
		t.Fatalf("load admin user seed: %v", err)
	}

	return admin
}
