package moderation_test

import (
	"net/http"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"gradlink/internal/db"
	"gradlink/internal/testutil"
)

func TestKeywordReviewSendsPostToPendingQueue(t *testing.T) {
	api := testutil.NewAPI(t)
	tag := seedTag(t, api, "综合提问", "general-questions", 90)
	seedKeyword(t, api, "加微信", "medium", db.StatusPending)
	cookie := testutil.RegisterUser(t, api, "risk@example.com", "Risk User")
	testutil.SaveProfile(t, api, cookie, "Risk User", "Fudan University", true)

	rec := api.JSON(t, http.MethodPost, "/api/posts", map[string]any{
		"title":  "有人说可以加微信内推，靠谱吗？",
		"body":   "对方要求加微信并购买培训服务，我不确定是不是广告。",
		"tagIds": []uint{tag.ID},
	}, cookie)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected pending review 202, got %d: %s", rec.Code, rec.Body.String())
	}

	payload := testutil.DecodeJSON(t, rec)
	if payload["status"] != db.StatusPending {
		t.Fatalf("expected pending response, got %s", rec.Body.String())
	}

	var post db.Post
	if err := api.DB.First(&post, uint(payload["id"].(float64))).Error; err != nil {
		t.Fatalf("load post: %v", err)
	}
	if post.Status != db.StatusPending {
		t.Fatalf("expected stored post pending, got %q", post.Status)
	}

	var item db.ModerationQueueItem
	if err := api.DB.Where("target_type = ? AND target_id = ?", "post", post.ID).First(&item).Error; err != nil {
		t.Fatalf("load moderation queue item: %v", err)
	}
	if item.Status != db.StatusPending || !strings.Contains(item.MatchedKeywords, "加微信") {
		t.Fatalf("unexpected moderation queue item: %+v", item)
	}
}

func TestRejectedKeywordRejectsPostAndQueuesReason(t *testing.T) {
	api := testutil.NewAPI(t)
	tag := seedTag(t, api, "综合提问", "general-questions", 90)
	seedKeyword(t, api, "代考", "high", db.StatusRejected)
	cookie := testutil.RegisterUser(t, api, "reject@example.com", "Reject User")
	testutil.SaveProfile(t, api, cookie, "Reject User", "Fudan University", true)

	rec := api.JSON(t, http.MethodPost, "/api/posts", map[string]any{
		"title":  "有人推荐代考服务",
		"body":   "这个代考广告看起来很不靠谱。",
		"tagIds": []uint{tag.ID},
	}, cookie)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected rejected content 422, got %d: %s", rec.Code, rec.Body.String())
	}
	if code := testutil.ErrorCode(t, rec); code != "CONTENT_REJECTED" {
		t.Fatalf("expected CONTENT_REJECTED, got %s", code)
	}

	var post db.Post
	if err := api.DB.Where("title = ?", "有人推荐代考服务").First(&post).Error; err != nil {
		t.Fatalf("load rejected post: %v", err)
	}
	if post.Status != db.StatusRejected {
		t.Fatalf("expected stored post rejected, got %q", post.Status)
	}

	var item db.ModerationQueueItem
	if err := api.DB.Where("target_type = ? AND target_id = ?", "post", post.ID).First(&item).Error; err != nil {
		t.Fatalf("load moderation queue item: %v", err)
	}
	if item.Status != db.StatusRejected || !strings.Contains(item.MatchedKeywords, "代考") {
		t.Fatalf("unexpected moderation queue item: %+v", item)
	}
}

func seedTag(t *testing.T, api *testutil.TestAPI, name, slug string, displayOrder int) db.Tag {
	t.Helper()

	tag := db.Tag{Name: name, Slug: slug, DisplayOrder: displayOrder}
	if err := api.DB.Create(&tag).Error; err != nil {
		t.Fatalf("create tag: %v", err)
	}
	return tag
}

func seedKeyword(t *testing.T, api *testutil.TestAPI, keyword, riskLevel, action string) db.ModerationKeyword {
	t.Helper()

	record := db.ModerationKeyword{Keyword: keyword, RiskLevel: riskLevel, Action: action, Enabled: true}
	if err := api.DB.Create(&record).Error; err != nil {
		t.Fatalf("create keyword: %v", err)
	}
	return record
}

func seedAdminCookie(t *testing.T, api *testutil.TestAPI) *http.Cookie {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte("admin-secret"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash admin password: %v", err)
	}
	admin := db.User{
		Email:        "admin@example.com",
		PasswordHash: string(hash),
		Nickname:     "Admin",
		Role:         db.UserRoleAdmin,
		Status:       db.UserStatusActive,
	}
	if err := api.DB.Create(&admin).Error; err != nil {
		t.Fatalf("create admin: %v", err)
	}

	rec := api.JSON(t, http.MethodPost, "/api/auth/login", map[string]any{
		"email":    "admin@example.com",
		"password": "admin-secret",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("login admin: got %d: %s", rec.Code, rec.Body.String())
	}
	return testutil.SessionCookie(t, rec)
}
