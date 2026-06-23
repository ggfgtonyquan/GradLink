package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gradlink/internal/db"
	apphttp "gradlink/internal/http"
)

type TestAPI struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func NewAPI(t *testing.T) *TestAPI {
	t.Helper()

	gin.SetMode(gin.TestMode)

	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	router := gin.New()
	apphttp.RegisterRoutes(router, &apphttp.Deps{DB: database})

	return &TestAPI{Router: router, DB: database}
}

func (api *TestAPI) JSON(t *testing.T, method, path string, body any, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	t.Helper()

	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		reader = bytes.NewReader(payload)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	rec := httptest.NewRecorder()
	api.Router.ServeHTTP(rec, req)
	return rec
}

func DecodeJSON(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode JSON response %q: %v", rec.Body.String(), err)
	}

	return payload
}

func ErrorCode(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()

	payload := DecodeJSON(t, rec)
	errPayload, ok := payload["error"].(map[string]any)
	if !ok {
		t.Fatalf("missing error object in response: %s", rec.Body.String())
	}
	code, ok := errPayload["code"].(string)
	if !ok {
		t.Fatalf("missing error code in response: %s", rec.Body.String())
	}

	return code
}

func SessionCookie(t *testing.T, rec *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()

	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "gradlink_session" {
			return cookie
		}
	}

	t.Fatalf("missing session cookie in response: %s", rec.Body.String())
	return nil
}

func RegisterUser(t *testing.T, api *TestAPI, email, nickname string) *http.Cookie {
	t.Helper()

	rec := api.JSON(t, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    email,
		"password": "correct horse battery staple",
		"nickname": nickname,
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("register user: got %d: %s", rec.Code, rec.Body.String())
	}
	return SessionCookie(t, rec)
}

func SaveProfile(t *testing.T, api *TestAPI, cookie *http.Cookie, nickname, school string, schoolVisible bool) {
	t.Helper()

	rec := api.JSON(t, http.MethodPut, "/api/me/profile", map[string]any{
		"nickname":       nickname,
		"school":         school,
		"educationLevel": "master",
		"major":          "Computer Science",
		"graduationYear": 2027,
		"schoolVisible":  schoolVisible,
	}, cookie)
	if rec.Code != http.StatusOK {
		t.Fatalf("save profile: got %d: %s", rec.Code, rec.Body.String())
	}
}

func SeedUserWithProfile(t *testing.T, api *TestAPI, email, nickname, school string, schoolVisible bool) db.User {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte("correct horse battery staple"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	user := db.User{
		Email:        email,
		PasswordHash: string(hash),
		Nickname:     nickname,
		Role:         db.UserRoleUser,
		Status:       db.UserStatusActive,
	}
	if err := api.DB.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	profile := db.StudentProfile{
		UserID:         user.ID,
		School:         school,
		EducationLevel: "master",
		Major:          "Computer Science",
		GraduationYear: 2027,
		SchoolVisible:  schoolVisible,
	}
	if err := api.DB.Create(&profile).Error; err != nil {
		t.Fatalf("create profile: %v", err)
	}
	return user
}

func SeedPost(t *testing.T, api *TestAPI, authorID uint, title, body, status string) db.Post {
	t.Helper()

	now := time.Now()
	post := db.Post{
		AuthorID:  authorID,
		Title:     title,
		Body:      body,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := api.DB.Create(&post).Error; err != nil {
		t.Fatalf("create post: %v", err)
	}
	return post
}
