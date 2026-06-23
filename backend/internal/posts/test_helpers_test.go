package posts_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gradlink/internal/db"
	"gradlink/internal/testutil"
)

func registerUser(t *testing.T, api *testutil.TestAPI, email, nickname string) *http.Cookie {
	t.Helper()

	rec := api.JSON(t, http.MethodPost, "/api/auth/register", map[string]any{
		"email":    email,
		"password": "correct horse battery staple",
		"nickname": nickname,
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("register user: got %d: %s", rec.Code, rec.Body.String())
	}
	return testutil.SessionCookie(t, rec)
}

func saveProfile(t *testing.T, api *testutil.TestAPI, cookie *http.Cookie, nickname, school string, schoolVisible bool) {
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

func seedUserWithProfile(t *testing.T, api *testutil.TestAPI, email, nickname, school string, schoolVisible bool) db.User {
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

func seedTag(t *testing.T, api *testutil.TestAPI, name, slug string, displayOrder int) db.Tag {
	t.Helper()

	tag := db.Tag{Name: name, Slug: slug, DisplayOrder: displayOrder}
	if err := api.DB.Create(&tag).Error; err != nil {
		t.Fatalf("create tag: %v", err)
	}
	return tag
}

func seedPost(t *testing.T, api *testutil.TestAPI, authorID uint, title, body, status string, createdAt time.Time, hotScore int) db.Post {
	t.Helper()

	post := db.Post{
		AuthorID:  authorID,
		Title:     title,
		Body:      body,
		Status:    status,
		HotScore:  hotScore,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
	if err := api.DB.Create(&post).Error; err != nil {
		t.Fatalf("create post: %v", err)
	}
	return post
}

func attachTags(t *testing.T, api *testutil.TestAPI, postID uint, tagIDs ...uint) {
	t.Helper()

	for _, tagID := range tagIDs {
		if err := api.DB.Create(&db.PostTag{PostID: postID, TagID: tagID}).Error; err != nil {
			t.Fatalf("attach tag %d to post %d: %v", tagID, postID, err)
		}
	}
}

func responsePosts(t *testing.T, rec *httptest.ResponseRecorder) []map[string]any {
	t.Helper()

	payload := testutil.DecodeJSON(t, rec)
	rawPosts, ok := payload["posts"].([]any)
	if !ok {
		t.Fatalf("expected posts array, got %s", rec.Body.String())
	}
	posts := make([]map[string]any, 0, len(rawPosts))
	for i, raw := range rawPosts {
		post, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("expected post object at %d, got %T", i, raw)
		}
		posts = append(posts, post)
	}
	return posts
}

func uintString(value uint) string {
	return strconv.FormatUint(uint64(value), 10)
}
