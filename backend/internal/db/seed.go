package db

import "gorm.io/gorm"

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
	{Keyword: "包过", RiskLevel: "high", Action: StatusPending, Enabled: true},
	{Keyword: "代考", RiskLevel: "high", Action: StatusRejected, Enabled: true},
}

func Seed(database *gorm.DB) error {
	return database.Transaction(func(tx *gorm.DB) error {
		for _, tag := range InitialTags {
			if err := tx.Where("slug = ?", tag.Slug).FirstOrCreate(&tag).Error; err != nil {
				return err
			}
		}

		admin := User{
			Email:        "admin@gradlink.local",
			PasswordHash: "$2a$10$7EqJtq98hPqEX7fNZaFWoOhi/e1tRg3Y5tYk.OwMfe0X3Syq3JqjG",
			Nickname:     "GradLink Admin",
			Role:         "admin",
			Status:       "active",
		}
		if err := tx.Where("email = ?", admin.Email).FirstOrCreate(&admin).Error; err != nil {
			return err
		}

		for _, keyword := range initialModerationKeywords {
			if err := tx.Where("keyword = ?", keyword.Keyword).FirstOrCreate(&keyword).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
