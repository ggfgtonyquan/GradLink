package auth_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gradlink/internal/db"
	"gradlink/internal/testutil"
)

func TestRegisterCreatesActiveUserWithHashedPasswordAndSession(t *testing.T) {
	api := testutil.NewAPI(t)

	rec := api.JSON(t, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    " Student@Example.COM ",
		"password": "correct horse battery staple",
		"nickname": "Student One",
	})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	cookie := testutil.SessionCookie(t, rec)
	if !cookie.HttpOnly {
		t.Fatal("session cookie must be HTTP-only")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected SameSite=Lax, got %v", cookie.SameSite)
	}
	if cookie.Value == "" {
		t.Fatal("session cookie value must not be empty")
	}

	var user db.User
	if err := api.DB.Where("email = ?", "student@example.com").First(&user).Error; err != nil {
		t.Fatalf("load registered user: %v", err)
	}
	if user.Status != db.UserStatusActive || user.Role != db.UserRoleUser {
		t.Fatalf("unexpected user role/status: %+v", user)
	}
	if strings.Contains(user.PasswordHash, "correct horse battery staple") {
		t.Fatal("password hash contains raw password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("correct horse battery staple")); err != nil {
		t.Fatalf("password was not stored as bcrypt hash: %v", err)
	}

	var session db.Session
	if err := api.DB.First(&session).Error; err != nil {
		t.Fatalf("load created session: %v", err)
	}
	if session.UserID != user.ID {
		t.Fatalf("session belongs to user %d, want %d", session.UserID, user.ID)
	}
	if session.TokenHash == cookie.Value || strings.Contains(session.TokenHash, cookie.Value) {
		t.Fatal("database must store a token hash, not the raw cookie token")
	}
	if !session.ExpiresAt.After(time.Now()) {
		t.Fatalf("session should expire in the future, got %s", session.ExpiresAt)
	}
}

func TestRegisterRollsBackUserWhenSessionCreationFails(t *testing.T) {
	api := testutil.NewAPI(t)
	if err := api.DB.Migrator().DropTable(&db.Session{}); err != nil {
		t.Fatalf("drop sessions table: %v", err)
	}

	rec := api.JSON(t, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    "rollback@example.com",
		"password": "correct horse battery staple",
		"nickname": "Rollback User",
	})
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected session failure to return 500, got %d: %s", rec.Code, rec.Body.String())
	}

	var count int64
	if err := api.DB.Model(&db.User{}).Where("email = ?", "rollback@example.com").Count(&count).Error; err != nil {
		t.Fatalf("count users: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected registration transaction rollback, found %d user rows", count)
	}
}

func TestRegisterCanIssueSecureSessionCookie(t *testing.T) {
	t.Setenv("GRADLINK_COOKIE_SECURE", "true")
	api := testutil.NewAPI(t)

	rec := api.JSON(t, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    "secure@example.com",
		"password": "correct horse battery staple",
		"nickname": "Secure User",
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	cookie := testutil.SessionCookie(t, rec)
	if !cookie.Secure {
		t.Fatal("expected secure cookie when GRADLINK_COOKIE_SECURE=true")
	}
}

func TestRegisterValidationAndDuplicateEmailUseConsistentErrors(t *testing.T) {
	api := testutil.NewAPI(t)

	rec := api.JSON(t, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    "not-an-email",
		"password": "short",
		"nickname": "",
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
	if code := testutil.ErrorCode(t, rec); code != "VALIDATION_FAILED" {
		t.Fatalf("expected VALIDATION_FAILED, got %s", code)
	}

	first := api.JSON(t, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    "student@example.com",
		"password": "correct horse battery staple",
		"nickname": "Student One",
	})
	if first.Code != http.StatusCreated {
		t.Fatalf("expected first registration to succeed, got %d: %s", first.Code, first.Body.String())
	}

	duplicate := api.JSON(t, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    " STUDENT@example.com ",
		"password": "correct horse battery staple",
		"nickname": "Student Two",
	})
	if duplicate.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 duplicate email, got %d: %s", duplicate.Code, duplicate.Body.String())
	}
	if code := testutil.ErrorCode(t, duplicate); code != "VALIDATION_FAILED" {
		t.Fatalf("expected VALIDATION_FAILED, got %s", code)
	}

	tooLong := api.JSON(t, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    "long-password@example.com",
		"password": strings.Repeat("a", 73),
		"nickname": "Long Password",
	})
	if tooLong.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for too-long password, got %d: %s", tooLong.Code, tooLong.Body.String())
	}
	if code := testutil.ErrorCode(t, tooLong); code != "VALIDATION_FAILED" {
		t.Fatalf("expected VALIDATION_FAILED, got %s", code)
	}
}

func TestRegisterDatabaseFailureIsInternalError(t *testing.T) {
	api := testutil.NewAPI(t)
	if err := api.DB.Migrator().DropTable(&db.User{}); err != nil {
		t.Fatalf("drop users table: %v", err)
	}

	rec := api.JSON(t, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    "broken@example.com",
		"password": "correct horse battery staple",
		"nickname": "Broken User",
	})
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected real database failure to return 500, got %d: %s", rec.Code, rec.Body.String())
	}
	if code := testutil.ErrorCode(t, rec); code != "INTERNAL_ERROR" {
		t.Fatalf("expected INTERNAL_ERROR, got %s", code)
	}
}

func TestLoginCreatesSessionForActiveUserAndRejectsDisabledUsers(t *testing.T) {
	api := testutil.NewAPI(t)
	hash, err := bcrypt.GenerateFromPassword([]byte("active-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	active := db.User{
		Email:        "active@example.com",
		PasswordHash: string(hash),
		Nickname:     "Active",
		Role:         db.UserRoleUser,
		Status:       db.UserStatusActive,
	}
	disabled := db.User{
		Email:        "disabled@example.com",
		PasswordHash: string(hash),
		Nickname:     "Disabled",
		Role:         db.UserRoleUser,
		Status:       db.UserStatusDisabled,
	}
	if err := api.DB.Create(&active).Error; err != nil {
		t.Fatalf("create active user: %v", err)
	}
	if err := api.DB.Create(&disabled).Error; err != nil {
		t.Fatalf("create disabled user: %v", err)
	}

	rec := api.JSON(t, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    " ACTIVE@example.com ",
		"password": "active-password",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	testutil.SessionCookie(t, rec)

	disabledRec := api.JSON(t, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "disabled@example.com",
		"password": "active-password",
	})
	if disabledRec.Code != http.StatusForbidden {
		t.Fatalf("expected disabled user forbidden, got %d: %s", disabledRec.Code, disabledRec.Body.String())
	}
	if code := testutil.ErrorCode(t, disabledRec); code != "FORBIDDEN" {
		t.Fatalf("expected FORBIDDEN, got %s", code)
	}
}

func TestSeededAdminLoginRequiresActiveAdmin(t *testing.T) {
	api := testutil.NewAPI(t)

	if err := db.Seed(api.DB); err != nil {
		t.Fatalf("seed default admin: %v", err)
	}
	defaultAdmin := api.JSON(t, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "admin@gradlink.local",
		"password": "password",
	})
	if defaultAdmin.Code != http.StatusUnauthorized {
		t.Fatalf("expected disabled default admin login to fail, got %d: %s", defaultAdmin.Code, defaultAdmin.Body.String())
	}

	api = testutil.NewAPI(t)
	if err := db.SeedWithOptions(api.DB, db.SeedOptions{AdminPassword: "admin-secret"}); err != nil {
		t.Fatalf("seed active admin: %v", err)
	}
	activeAdmin := api.JSON(t, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "admin@gradlink.local",
		"password": "admin-secret",
	})
	if activeAdmin.Code != http.StatusOK {
		t.Fatalf("expected active admin login to succeed, got %d: %s", activeAdmin.Code, activeAdmin.Body.String())
	}
}

func TestMeAndLogoutRequireValidActiveSession(t *testing.T) {
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

	me := api.JSON(t, http.MethodGet, "/api/me", nil, cookie)
	if me.Code != http.StatusOK {
		t.Fatalf("expected /api/me success, got %d: %s", me.Code, me.Body.String())
	}
	payload := testutil.DecodeJSON(t, me)
	if payload["profileComplete"] != false {
		t.Fatalf("expected incomplete profile, got %s", me.Body.String())
	}

	var user db.User
	if err := api.DB.Where("email = ?", "student@example.com").First(&user).Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	if err := api.DB.Model(&user).Update("status", db.UserStatusDisabled).Error; err != nil {
		t.Fatalf("disable user: %v", err)
	}
	disabledMe := api.JSON(t, http.MethodGet, "/api/me", nil, cookie)
	if disabledMe.Code != http.StatusForbidden {
		t.Fatalf("expected disabled session forbidden, got %d: %s", disabledMe.Code, disabledMe.Body.String())
	}
	if code := testutil.ErrorCode(t, disabledMe); code != "FORBIDDEN" {
		t.Fatalf("expected FORBIDDEN, got %s", code)
	}

	if err := api.DB.Model(&user).Update("status", db.UserStatusActive).Error; err != nil {
		t.Fatalf("reactivate user: %v", err)
	}
	logout := api.JSON(t, http.MethodPost, "/api/auth/logout", nil, cookie)
	if logout.Code != http.StatusOK {
		t.Fatalf("expected logout success, got %d: %s", logout.Code, logout.Body.String())
	}
	afterLogout := api.JSON(t, http.MethodGet, "/api/me", nil, cookie)
	if afterLogout.Code != http.StatusUnauthorized {
		t.Fatalf("expected logged out session unauthorized, got %d: %s", afterLogout.Code, afterLogout.Body.String())
	}
	if code := testutil.ErrorCode(t, afterLogout); code != "UNAUTHENTICATED" {
		t.Fatalf("expected UNAUTHENTICATED, got %s", code)
	}
}
