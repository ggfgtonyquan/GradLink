package posts_test

import (
	"net/http"
	"testing"

	"gradlink/internal/db"
	"gradlink/internal/testutil"
)

func TestCreatePostRequiresAuthAndCompleteProfile(t *testing.T) {
	api := testutil.NewAPI(t)
	tag := seedTag(t, api, "面试经验", "interview-experience", 10)

	unauthenticated := api.JSON(t, http.MethodPost, "/api/posts", map[string]any{
		"title":  "面试复盘",
		"body":   "今天整理一下技术面试里的重点问题。",
		"tagIds": []uint{tag.ID},
	})
	if unauthenticated.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthenticated post create to fail with 401, got %d: %s", unauthenticated.Code, unauthenticated.Body.String())
	}
	if code := testutil.ErrorCode(t, unauthenticated); code != "UNAUTHENTICATED" {
		t.Fatalf("expected UNAUTHENTICATED, got %s", code)
	}

	cookie := registerUser(t, api, "student@example.com", "Student One")
	incomplete := api.JSON(t, http.MethodPost, "/api/posts", map[string]any{
		"title":  "面试复盘",
		"body":   "今天整理一下技术面试里的重点问题。",
		"tagIds": []uint{tag.ID},
	}, cookie)
	if incomplete.Code != http.StatusForbidden {
		t.Fatalf("expected incomplete profile to fail with 403, got %d: %s", incomplete.Code, incomplete.Body.String())
	}
	if code := testutil.ErrorCode(t, incomplete); code != "PROFILE_REQUIRED" {
		t.Fatalf("expected PROFILE_REQUIRED, got %s", code)
	}
}

func TestCreatePostValidatesPayloadAndTagIDs(t *testing.T) {
	api := testutil.NewAPI(t)
	tag := seedTag(t, api, "Offer 比较", "offer-comparison", 20)
	cookie := registerUser(t, api, "student@example.com", "Student One")
	saveProfile(t, api, cookie, "Student One", "Fudan University", true)

	cases := []struct {
		name string
		body map[string]any
	}{
		{
			name: "missing title",
			body: map[string]any{"title": "", "body": "body long enough", "tagIds": []uint{tag.ID}},
		},
		{
			name: "missing body",
			body: map[string]any{"title": "valid title", "body": "", "tagIds": []uint{tag.ID}},
		},
		{
			name: "zero tags",
			body: map[string]any{"title": "valid title", "body": "body long enough", "tagIds": []uint{}},
		},
		{
			name: "more than three tags",
			body: map[string]any{"title": "valid title", "body": "body long enough", "tagIds": []uint{tag.ID, tag.ID + 1, tag.ID + 2, tag.ID + 3}},
		},
		{
			name: "unknown tag",
			body: map[string]any{"title": "valid title", "body": "body long enough", "tagIds": []uint{99999}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := api.JSON(t, http.MethodPost, "/api/posts", tc.body, cookie)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected validation failure, got %d: %s", rec.Code, rec.Body.String())
			}
			if code := testutil.ErrorCode(t, rec); code != "VALIDATION_FAILED" {
				t.Fatalf("expected VALIDATION_FAILED, got %s", code)
			}
		})
	}
}

func TestCreatePostPublishesWithCurrentUserAndTags(t *testing.T) {
	api := testutil.NewAPI(t)
	tagA := seedTag(t, api, "秋招进度", "autumn-recruiting", 30)
	tagB := seedTag(t, api, "简历求助", "resume-help", 50)
	cookie := registerUser(t, api, "student@example.com", "Student One")
	saveProfile(t, api, cookie, "Student One", "Shanghai Jiao Tong University", false)

	rec := api.JSON(t, http.MethodPost, "/api/posts", map[string]any{
		"title":  " 秋招时间线整理 ",
		"body":   " 从网申到终面，我把每个阶段的准备事项写下来。 ",
		"tagIds": []uint{tagA.ID, tagB.ID},
	}, cookie)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected post create success, got %d: %s", rec.Code, rec.Body.String())
	}

	payload := testutil.DecodeJSON(t, rec)
	if payload["title"] != "秋招时间线整理" || payload["status"] != db.StatusPublished {
		t.Fatalf("unexpected post payload: %s", rec.Body.String())
	}
	if payload["body"] != "从网申到终面，我把每个阶段的准备事项写下来。" {
		t.Fatalf("expected trimmed body, got %s", rec.Body.String())
	}
	author := payload["author"].(map[string]any)
	if author["nickname"] != "Student One" || author["school"] != nil {
		t.Fatalf("expected semi-anonymous author without hidden school, got %+v", author)
	}
	tags := payload["tags"].([]any)
	if len(tags) != 2 {
		t.Fatalf("expected two tags, got %s", rec.Body.String())
	}

	var stored db.Post
	if err := api.DB.First(&stored, uint(payload["id"].(float64))).Error; err != nil {
		t.Fatalf("load created post: %v", err)
	}
	if stored.Status != db.StatusPublished || stored.AuthorID == 0 {
		t.Fatalf("post not published with author: %+v", stored)
	}
	var tagCount int64
	if err := api.DB.Model(&db.PostTag{}).Where("post_id = ?", stored.ID).Count(&tagCount).Error; err != nil {
		t.Fatalf("count post tags: %v", err)
	}
	if tagCount != 2 {
		t.Fatalf("expected two post tags, got %d", tagCount)
	}
}
