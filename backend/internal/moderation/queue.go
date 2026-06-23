package moderation

import (
	"strings"

	"gorm.io/gorm"
	"gradlink/internal/db"
)

type ReviewDecision struct {
	Status          string
	MatchedKeywords []string
}

func ReviewText(tx *gorm.DB, content string) (ReviewDecision, error) {
	var keywords []db.ModerationKeyword
	if err := tx.Where("enabled = ?", true).Order("id ASC").Find(&keywords).Error; err != nil {
		return ReviewDecision{}, err
	}

	normalizedContent := strings.ToLower(content)
	matches := make([]string, 0)
	status := db.StatusPublished
	for _, keyword := range keywords {
		cleanKeyword := strings.TrimSpace(keyword.Keyword)
		if cleanKeyword == "" || !strings.Contains(normalizedContent, strings.ToLower(cleanKeyword)) {
			continue
		}

		matches = append(matches, cleanKeyword)
		if keyword.Action == db.StatusRejected {
			status = db.StatusRejected
			continue
		}
		if status != db.StatusRejected && keyword.Action == db.StatusPending {
			status = db.StatusPending
		}
	}

	return ReviewDecision{Status: status, MatchedKeywords: matches}, nil
}

func EnqueueIfNeeded(tx *gorm.DB, targetType string, targetID uint, decision ReviewDecision) error {
	if len(decision.MatchedKeywords) == 0 || decision.Status == db.StatusPublished {
		return nil
	}

	item := db.ModerationQueueItem{
		TargetType:      targetType,
		TargetID:        targetID,
		MatchedKeywords: strings.Join(decision.MatchedKeywords, ","),
		Status:          decision.Status,
	}
	return tx.Create(&item).Error
}
