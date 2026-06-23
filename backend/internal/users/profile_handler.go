package users

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gradlink/internal/auth"
	"gradlink/internal/db"
)

type ProfileHandler struct {
	DB *gorm.DB
}

func NewProfileHandler(database *gorm.DB) *ProfileHandler {
	return &ProfileHandler{DB: database}
}

func (h *ProfileHandler) Me(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		auth.AbortWithError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Authentication required")
		return
	}

	profile, complete, err := h.profileForUser(user.ID)
	if err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	payload := gin.H{
		"id":              user.ID,
		"email":           user.Email,
		"nickname":        user.Nickname,
		"role":            user.Role,
		"status":          user.Status,
		"profileComplete": complete,
	}
	if complete {
		payload["profile"] = profileResponse(profile)
	}

	c.JSON(http.StatusOK, payload)
}

func (h *ProfileHandler) PutProfile(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		auth.AbortWithError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Authentication required")
		return
	}

	var input struct {
		Nickname       string `json:"nickname"`
		School         string `json:"school"`
		EducationLevel string `json:"educationLevel"`
		Major          string `json:"major"`
		GraduationYear int    `json:"graduationYear"`
		SchoolVisible  bool   `json:"schoolVisible"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		auth.AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid profile payload")
		return
	}

	nickname := strings.TrimSpace(input.Nickname)
	school := strings.TrimSpace(input.School)
	educationLevel := strings.TrimSpace(input.EducationLevel)
	major := strings.TrimSpace(input.Major)
	if nickname == "" || len(nickname) > 80 ||
		school == "" || len(school) > 120 ||
		educationLevel == "" || len(educationLevel) > 40 ||
		major == "" || len(major) > 120 ||
		input.GraduationYear < 1900 || input.GraduationYear > 2100 {
		auth.AbortWithError(c, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid profile payload")
		return
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&db.User{}).Where("id = ?", user.ID).Update("nickname", nickname).Error; err != nil {
			return err
		}

		values := map[string]interface{}{
			"user_id":         user.ID,
			"school":          school,
			"education_level": educationLevel,
			"major":           major,
			"graduation_year": input.GraduationYear,
			"school_visible":  input.SchoolVisible,
		}

		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"school":          school,
				"education_level": educationLevel,
				"major":           major,
				"graduation_year": input.GraduationYear,
				"school_visible":  input.SchoolVisible,
				"updated_at":      gorm.Expr("CURRENT_TIMESTAMP"),
			}),
		}).Model(&db.StudentProfile{}).Create(values).Error
	})
	if err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	profile, _, err := h.profileForUser(user.ID)
	if err != nil {
		auth.AbortWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"profileComplete": true,
		"profile":         profileResponse(profile),
	})
}

func (h *ProfileHandler) profileForUser(userID uint) (*db.StudentProfile, bool, error) {
	var profile db.StudentProfile
	err := h.DB.Where("user_id = ?", userID).First(&profile).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return &profile, true, nil
}

func profileResponse(profile *db.StudentProfile) gin.H {
	return gin.H{
		"school":         profile.School,
		"educationLevel": profile.EducationLevel,
		"major":          profile.Major,
		"graduationYear": profile.GraduationYear,
		"schoolVisible":  profile.SchoolVisible,
	}
}
