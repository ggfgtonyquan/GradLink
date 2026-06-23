package posts_test

import (
	"net/http"
	"testing"
	"time"

	"gradlink/internal/db"
	"gradlink/internal/testutil"
)

func TestGetPostDetailReturnsPublishedPostWithTagsAndAuthor(t *testing.T) {
	api := testutil.NewAPI(t)
	tagA := seedTag(t, api, "笔试讨论", "written-test", 60)
	tagB := seedTag(t, api, "吐槽避坑", "warnings", 80)
	author := seedUserWithProfile(t, api, "author@example.com", "Author Nick", "Wuhan University", true)
	post := seedPost(t, api, author.ID, "笔试经验分享", "完整笔试复盘正文，包含题型、时间安排和复习建议。", db.StatusPublished, time.Now(), 7)
	attachTags(t, api, post.ID, tagA.ID, tagB.ID)

	rec := api.JSON(t, http.MethodGet, "/api/posts/"+uintString(post.ID), nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected detail success, got %d: %s", rec.Code, rec.Body.String())
	}
	payload := testutil.DecodeJSON(t, rec)
	if payload["id"] != float64(post.ID) || payload["title"] != post.Title || payload["body"] != post.Body {
		t.Fatalf("unexpected detail payload: %s", rec.Body.String())
	}
	if payload["status"] != db.StatusPublished || payload["createdAt"] == "" {
		t.Fatalf("expected status and createdAt, got %s", rec.Body.String())
	}
	authorPayload := payload["author"].(map[string]any)
	if authorPayload["nickname"] != "Author Nick" || authorPayload["school"] != "Wuhan University" {
		t.Fatalf("expected visible author school, got %+v", authorPayload)
	}
	if len(payload["tags"].([]any)) != 2 {
		t.Fatalf("expected two tags, got %s", rec.Body.String())
	}
	counts := payload["counts"].(map[string]any)
	if counts["likes"] != float64(0) || counts["comments"] != float64(0) || counts["favorites"] != float64(0) {
		t.Fatalf("unexpected counts payload: %+v", counts)
	}
}

func TestGetPostDetailReturnsNotFoundForMissingOrNonPublishedPost(t *testing.T) {
	api := testutil.NewAPI(t)
	author := seedUserWithProfile(t, api, "author@example.com", "Author Nick", "Wuhan University", true)
	hidden := seedPost(t, api, author.ID, "隐藏帖子", "不应该可见", db.StatusHidden, time.Now(), 0)

	missing := api.JSON(t, http.MethodGet, "/api/posts/99999", nil)
	if missing.Code != http.StatusNotFound {
		t.Fatalf("expected missing post 404, got %d: %s", missing.Code, missing.Body.String())
	}
	if code := testutil.ErrorCode(t, missing); code != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND, got %s", code)
	}

	nonPublished := api.JSON(t, http.MethodGet, "/api/posts/"+uintString(hidden.ID), nil)
	if nonPublished.Code != http.StatusNotFound {
		t.Fatalf("expected hidden post 404, got %d: %s", nonPublished.Code, nonPublished.Body.String())
	}
	if code := testutil.ErrorCode(t, nonPublished); code != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND, got %s", code)
	}
}
