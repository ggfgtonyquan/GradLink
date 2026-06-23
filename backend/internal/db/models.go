package db

import "time"

const (
	StatusPublished = "published"
	StatusPending   = "pending"
	StatusHidden    = "hidden"
	StatusRejected  = "rejected"
	StatusDeleted   = "deleted"

	UserRoleAdmin = "admin"
	UserRoleUser  = "user"

	UserStatusActive   = "active"
	UserStatusDisabled = "disabled"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"size:255;not null;uniqueIndex"`
	PasswordHash string `gorm:"size:255;not null"`
	Nickname     string `gorm:"size:80;not null"`
	Role         string `gorm:"size:32;not null;default:user;index"`
	Status       string `gorm:"size:32;not null;default:active;index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type StudentProfile struct {
	ID             uint   `gorm:"primaryKey"`
	UserID         uint   `gorm:"not null;uniqueIndex"`
	School         string `gorm:"size:120;not null"`
	EducationLevel string `gorm:"size:40;not null"`
	Major          string `gorm:"size:120;not null"`
	GraduationYear int    `gorm:"not null;index"`
	SchoolVisible  bool   `gorm:"not null;default:true"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	User           User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Post struct {
	ID            uint      `gorm:"primaryKey"`
	AuthorID      uint      `gorm:"not null;index"`
	Title         string    `gorm:"size:160;not null"`
	Body          string    `gorm:"type:text;not null"`
	Status        string    `gorm:"size:32;not null;default:published;index:idx_posts_status_created"`
	CommentCount  int       `gorm:"not null;default:0"`
	LikeCount     int       `gorm:"not null;default:0"`
	FavoriteCount int       `gorm:"not null;default:0"`
	HotScore      int       `gorm:"not null;default:0;index"`
	CreatedAt     time.Time `gorm:"index:idx_posts_status_created"`
	UpdatedAt     time.Time
	Author        User `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type Tag struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:80;not null;uniqueIndex"`
	Slug         string `gorm:"size:120;not null;uniqueIndex"`
	DisplayOrder int    `gorm:"not null;default:0;index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PostTag struct {
	PostID    uint `gorm:"primaryKey"`
	TagID     uint `gorm:"primaryKey;index"`
	CreatedAt time.Time
	Post      Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tag       Tag  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Comment struct {
	ID              uint   `gorm:"primaryKey"`
	PostID          uint   `gorm:"not null;index:idx_comments_post_parent"`
	AuthorID        uint   `gorm:"not null;index"`
	ParentCommentID *uint  `gorm:"index:idx_comments_post_parent"`
	Body            string `gorm:"type:text;not null"`
	Status          string `gorm:"size:32;not null;default:published;index"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Post            Post     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Author          User     `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ParentComment   *Comment `gorm:"foreignKey:ParentCommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Like struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"not null;uniqueIndex:idx_likes_user_target"`
	TargetType string `gorm:"size:32;not null;uniqueIndex:idx_likes_user_target"`
	TargetID   uint   `gorm:"not null;uniqueIndex:idx_likes_user_target"`
	CreatedAt  time.Time
	User       User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Favorite struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"not null;uniqueIndex:idx_favorites_user_post"`
	PostID    uint `gorm:"not null;uniqueIndex:idx_favorites_user_post"`
	CreatedAt time.Time
	User      User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Post      Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Report struct {
	ID         uint   `gorm:"primaryKey"`
	ReporterID uint   `gorm:"not null;index"`
	TargetType string `gorm:"size:32;not null;index:idx_reports_target"`
	TargetID   uint   `gorm:"not null;index:idx_reports_target"`
	Reason     string `gorm:"size:80;not null"`
	Detail     string `gorm:"type:text"`
	Status     string `gorm:"size:32;not null;default:pending;index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Reporter   User `gorm:"foreignKey:ReporterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type ModerationKeyword struct {
	ID        uint   `gorm:"primaryKey"`
	Keyword   string `gorm:"size:120;not null;uniqueIndex"`
	RiskLevel string `gorm:"size:32;not null"`
	Action    string `gorm:"size:32;not null"`
	Enabled   bool   `gorm:"not null;default:true;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ModerationQueueItem struct {
	ID              uint   `gorm:"primaryKey"`
	TargetType      string `gorm:"size:32;not null;index:idx_moderation_queue_target"`
	TargetID        uint   `gorm:"not null;index:idx_moderation_queue_target"`
	MatchedKeywords string `gorm:"type:text;not null"`
	Status          string `gorm:"size:32;not null;default:pending;index"`
	ReviewerID      *uint  `gorm:"index"`
	Decision        string `gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Reviewer        *User `gorm:"foreignKey:ReviewerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (ModerationQueueItem) TableName() string {
	return "moderation_queue"
}

type Session struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	TokenHash string    `gorm:"size:255;not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
