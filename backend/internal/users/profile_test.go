package users_test

import (
	"net/http"
	"testing"

	"gradlink/internal/db"
	"gradlink/internal/testutil"
)

func TestPutProfileSavesStudentProfileAndMeIncludesIt(t *testing.T) {
	api := testutil.NewAPI(t)
	register := api.JSON(t, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    "student@example.com",
		"password": "correct horse battery staple",
		"nickname": "Student One",
	})
	if register.Code != http.StatusCreated {
		t.Fatalf("expected register success, got %d: %s", register.Code, register.Body.String())
	}
	cookie := testutil.SessionCookie(t, register)

	rec := api.JSON(t, http.MethodPut, "/api/me/profile", map[string]any{
		"nickname":       "Student Prime",
		"school":         "Zhejiang University",
		"educationLevel": "master",
		"major":          "Computer Science",
		"graduationYear": 2027,
		"schoolVisible":  false,
	}, cookie)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected profile save success, got %d: %s", rec.Code, rec.Body.String())
	}

	var user db.User
	if err := api.DB.Where("email = ?", "student@example.com").First(&user).Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	if user.Nickname != "Student Prime" {
		t.Fatalf("expected nickname update, got %q", user.Nickname)
	}
	var profile db.StudentProfile
	if err := api.DB.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		t.Fatalf("load profile: %v", err)
	}
	if profile.School != "Zhejiang University" ||
		profile.EducationLevel != "master" ||
		profile.Major != "Computer Science" ||
		profile.GraduationYear != 2027 ||
		profile.SchoolVisible {
		t.Fatalf("profile was not saved correctly: %+v", profile)
	}

	me := api.JSON(t, http.MethodGet, "/api/me", nil, cookie)
	if me.Code != http.StatusOK {
		t.Fatalf("expected /api/me success, got %d: %s", me.Code, me.Body.String())
	}
	payload := testutil.DecodeJSON(t, me)
	if payload["profileComplete"] != true {
		t.Fatalf("expected complete profile, got %s", me.Body.String())
	}
	profilePayload, ok := payload["profile"].(map[string]any)
	if !ok {
		t.Fatalf("expected profile object, got %s", me.Body.String())
	}
	if profilePayload["schoolVisible"] != false {
		t.Fatalf("expected schoolVisible false, got %s", me.Body.String())
	}
}

func TestPutProfileRequiresAuthAndValidPayload(t *testing.T) {
	api := testutil.NewAPI(t)

	unauthenticated := api.JSON(t, http.MethodPut, "/api/me/profile", map[string]any{
		"nickname":       "Student Prime",
		"school":         "Zhejiang University",
		"educationLevel": "master",
		"major":          "Computer Science",
		"graduationYear": 2027,
		"schoolVisible":  true,
	})
	if unauthenticated.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized, got %d: %s", unauthenticated.Code, unauthenticated.Body.String())
	}
	if code := testutil.ErrorCode(t, unauthenticated); code != "UNAUTHENTICATED" {
		t.Fatalf("expected UNAUTHENTICATED, got %s", code)
	}

	register := api.JSON(t, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    "student@example.com",
		"password": "correct horse battery staple",
		"nickname": "Student One",
	})
	if register.Code != http.StatusCreated {
		t.Fatalf("expected register success, got %d: %s", register.Code, register.Body.String())
	}
	cookie := testutil.SessionCookie(t, register)

	invalid := api.JSON(t, http.MethodPut, "/api/me/profile", map[string]any{
		"nickname":       "",
		"school":         "",
		"educationLevel": "",
		"major":          "",
		"graduationYear": 1800,
		"schoolVisible":  true,
	}, cookie)
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("expected validation failure, got %d: %s", invalid.Code, invalid.Body.String())
	}
	if code := testutil.ErrorCode(t, invalid); code != "VALIDATION_FAILED" {
		t.Fatalf("expected VALIDATION_FAILED, got %s", code)
	}
}

func TestProfileGateReturnsProfileRequiredForIncompleteProfiles(t *testing.T) {
	api := testutil.NewAPI(t)
	register := api.JSON(t, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    "student@example.com",
		"password": "correct horse battery staple",
		"nickname": "Student One",
	})
	if register.Code != http.StatusCreated {
		t.Fatalf("expected register success, got %d: %s", register.Code, register.Body.String())
	}
	cookie := testutil.SessionCookie(t, register)

	blocked := api.JSON(t, http.MethodPost, "/api/test/profile-gated", nil, cookie)
	if blocked.Code != http.StatusForbidden {
		t.Fatalf("expected profile gate forbidden, got %d: %s", blocked.Code, blocked.Body.String())
	}
	if code := testutil.ErrorCode(t, blocked); code != "PROFILE_REQUIRED" {
		t.Fatalf("expected PROFILE_REQUIRED, got %s", code)
	}

	saveProfile := api.JSON(t, http.MethodPut, "/api/me/profile", map[string]any{
		"nickname":       "Student Prime",
		"school":         "Zhejiang University",
		"educationLevel": "master",
		"major":          "Computer Science",
		"graduationYear": 2027,
		"schoolVisible":  true,
	}, cookie)
	if saveProfile.Code != http.StatusOK {
		t.Fatalf("expected profile save success, got %d: %s", saveProfile.Code, saveProfile.Body.String())
	}

	allowed := api.JSON(t, http.MethodPost, "/api/test/profile-gated", nil, cookie)
	if allowed.Code != http.StatusOK {
		t.Fatalf("expected profile-gated probe success, got %d: %s", allowed.Code, allowed.Body.String())
	}
}
