package db

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const seededAdminEmail = "admin@gradlink.local"
const legacySeededAdminHash = "$2a$10$7EqJtq98hPqEX7fNZaFWoOhi/e1tRg3Y5tYk.OwMfe0X3Syq3JqjG"
const legacySeededAdminPassword = "password"

type SeedOptions struct {
	AdminPassword string
}

var InitialTags = []Tag{
	{Name: "面试经验", Slug: "interview-experience", DisplayOrder: 10},
	{Name: "Offer 比较", Slug: "offer-comparison", DisplayOrder: 20},
	{Name: "秋招进度", Slug: "autumn-recruiting", DisplayOrder: 30},
	{Name: "春招进度", Slug: "spring-recruiting", DisplayOrder: 40},
	{Name: "简历求助", Slug: "resume-help", DisplayOrder: 50},
	{Name: "笔试讨论", Slug: "written-test", DisplayOrder: 60},
	{Name: "签约三方", Slug: "tripartite-agreement", DisplayOrder: 70},
	{Name: "吐槽避坑", Slug: "warnings", DisplayOrder: 80},
	{Name: "综合提问", Slug: "general-questions", DisplayOrder: 90},
}

var initialModerationKeywords = []ModerationKeyword{
	{Keyword: "内推收费", RiskLevel: "medium", Action: StatusPending, Enabled: true},
	{Keyword: "加微信", RiskLevel: "medium", Action: StatusPending, Enabled: true},
	{Keyword: "包过", RiskLevel: "high", Action: StatusPending, Enabled: true},
	{Keyword: "代考", RiskLevel: "high", Action: StatusRejected, Enabled: true},
}

func Seed(database *gorm.DB) error {
	return SeedWithOptions(database, SeedOptions{})
}

func SeedWithOptions(database *gorm.DB, options SeedOptions) error {
	return database.Transaction(func(tx *gorm.DB) error {
		for _, tag := range InitialTags {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "slug"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "display_order", "updated_at"}),
			}).Create(&tag).Error; err != nil {
				return err
			}
		}

		if err := seedAdmin(tx, options); err != nil {
			return err
		}

		for _, keyword := range initialModerationKeywords {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "keyword"}},
				DoUpdates: clause.AssignmentColumns([]string{"risk_level", "action", "enabled", "updated_at"}),
			}).Create(&keyword).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func seedAdmin(tx *gorm.DB, options SeedOptions) error {
	adminPassword := strings.TrimSpace(options.AdminPassword)
	if adminPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		admin := User{
			Email:        seededAdminEmail,
			PasswordHash: string(hash),
			Nickname:     "GradLink Admin",
			Role:         UserRoleAdmin,
			Status:       UserStatusActive,
		}
		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "email"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"password_hash": admin.PasswordHash,
				"nickname":      admin.Nickname,
				"role":          admin.Role,
				"status":        admin.Status,
				"updated_at":    gorm.Expr("CURRENT_TIMESTAMP"),
			}),
		}).Create(&admin).Error
	}

	var existing User
	err := tx.Where("email = ?", seededAdminEmail).First(&existing).Error
	if err == nil {
		if legacyAdminPasswordMatches(existing.PasswordHash) {
			hash, err := unusablePasswordHash()
			if err != nil {
				return err
			}
			return tx.Model(&existing).Updates(map[string]interface{}{
				"nickname":      "GradLink Admin",
				"password_hash": hash,
				"role":          UserRoleAdmin,
				"status":        UserStatusDisabled,
			}).Error
		}

		return tx.Model(&existing).Updates(map[string]interface{}{
			"nickname": "GradLink Admin",
			"role":     UserRoleAdmin,
		}).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	hash, err := unusablePasswordHash()
	if err != nil {
		return err
	}

	admin := User{
		Email:        seededAdminEmail,
		PasswordHash: string(hash),
		Nickname:     "GradLink Admin",
		Role:         UserRoleAdmin,
		Status:       UserStatusDisabled,
	}
	return tx.Create(&admin).Error
}

func legacyAdminPasswordMatches(hash string) bool {
	if hash == legacySeededAdminHash {
		return true
	}

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(legacySeededAdminPassword)) == nil
}

func unusablePasswordHash() (string, error) {
	token, err := randomToken(32)
	if err != nil {
		return "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("disabled-admin-"+token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
