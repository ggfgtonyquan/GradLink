package comments_test

import (
	"net/http"
	"strconv"
	"testing"

	"gradlink/internal/db"
	"gradlink/internal/testutil"
)

func TestCreateCommentRequiresAuthAndCompleteProfile(t *testing.T) {
	api := testutil.NewAPI(t)
	author := testutil.SeedUserWithProfile(t, api, "author@example.com", "Author", "Nanjing University", true)
	post := testutil.SeedPost(t, api, author.ID, "面试复盘", "正文内容", db.StatusPublished)

	unauthenticated := api.JSON(t, http.MethodPost, "/api/posts/"+uintString(post.ID)+"/comments", map[string]any{
		"body": "写得很有帮助",
	})
	if unauthenticated.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthenticated comment create to fail with 401, got %d: %s", unauthenticated.Code, unauthenticated.Body.String())
	}
	if code := testutil.ErrorCode(t, unauthenticated); code != "UNAUTHENTICATED" {
		t.Fatalf("expected UNAUTHENTICATED, got %s", code)
	}

	cookie := testutil.RegisterUser(t, api, "commenter@example.com", "Commenter")
	incomplete := api.JSON(t, http.MethodPost, "/api/posts/"+uintString(post.ID)+"/comments", map[string]any{
		"body": "写得很有帮助",
	}, cookie)
	if incomplete.Code != http.StatusForbidden {
		t.Fatalf("expected incomplete profile to fail with 403, got %d: %s", incomplete.Code, incomplete.Body.String())
	}
	if code := testutil.ErrorCode(t, incomplete); code != "PROFILE_REQUIRED" {
		t.Fatalf("expected PROFILE_REQUIRED, got %s", code)
	}
}

func TestCreateCommentOnPublishedPostIncrementsCountsAndHotScore(t *testing.T) {
	api := testutil.NewAPI(t)
	author := testutil.SeedUserWithProfile(t, api, "author@example.com", "Author", "Nanjing University", true)
	post := testutil.SeedPost(t, api, author.ID, "面试复盘", "正文内容", db.StatusPublished)
	cookie := testutil.RegisterUser(t, api, "commenter@example.com", "Commenter")
	testutil.SaveProfile(t, api, cookie, "Commenter", "Fudan University", true)

	rec := api.JSON(t, http.MethodPost, "/api/posts/"+uintString(post.ID)+"/comments", map[string]any{
		"body": " 写得很有帮助 ",
	}, cookie)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected comment create success, got %d: %s", rec.Code, rec.Body.String())
	}

	payload := testutil.DecodeJSON(t, rec)
	if payload["body"] != "写得很有帮助" || payload["status"] != db.StatusPublished {
		t.Fatalf("unexpected comment payload: %s", rec.Body.String())
	}
	if payload["parentCommentId"] != nil {
		t.Fatalf("top-level comment should not have parent: %s", rec.Body.String())
	}
	counts := payload["counts"].(map[string]any)
	if counts["comments"] != float64(1) {
		t.Fatalf("expected comment count 1, got %s", rec.Body.String())
	}

	var stored db.Comment
	if err := api.DB.First(&stored, uint(payload["id"].(float64))).Error; err != nil {
		t.Fatalf("load created comment: %v", err)
	}
	if stored.PostID != post.ID || stored.ParentCommentID != nil || stored.Status != db.StatusPublished {
		t.Fatalf("comment was not saved correctly: %+v", stored)
	}

	var updated db.Post
	if err := api.DB.First(&updated, post.ID).Error; err != nil {
		t.Fatalf("load updated post: %v", err)
	}
	if updated.CommentCount != 1 || updated.HotScore <= 0 {
		t.Fatalf("expected post counts and hot score to update, got %+v", updated)
	}
}

func TestCreateCommentValidatesPayloadAndPublishedPost(t *testing.T) {
	api := testutil.NewAPI(t)
	author := testutil.SeedUserWithProfile(t, api, "author@example.com", "Author", "Nanjing University", true)
	hidden := testutil.SeedPost(t, api, author.ID, "隐藏帖", "正文内容", db.StatusHidden)
	cookie := testutil.RegisterUser(t, api, "commenter@example.com", "Commenter")
	testutil.SaveProfile(t, api, cookie, "Commenter", "Fudan University", true)

	invalid := api.JSON(t, http.MethodPost, "/api/posts/"+uintString(hidden.ID)+"/comments", map[string]any{
		"body": "   ",
	}, cookie)
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("expected validation failure, got %d: %s", invalid.Code, invalid.Body.String())
	}
	if code := testutil.ErrorCode(t, invalid); code != "VALIDATION_FAILED" {
		t.Fatalf("expected VALIDATION_FAILED, got %s", code)
	}

	notFound := api.JSON(t, http.MethodPost, "/api/posts/"+uintString(hidden.ID)+"/comments", map[string]any{
		"body": "隐藏帖不能评论",
	}, cookie)
	if notFound.Code != http.StatusNotFound {
		t.Fatalf("expected hidden post to be not found, got %d: %s", notFound.Code, notFound.Body.String())
	}
	if code := testutil.ErrorCode(t, notFound); code != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND, got %s", code)
	}
}

func TestCreateReplyAllowsOnlyOneReplyLevel(t *testing.T) {
	api := testutil.NewAPI(t)
	author := testutil.SeedUserWithProfile(t, api, "author@example.com", "Author", "Nanjing University", true)
	post := testutil.SeedPost(t, api, author.ID, "面试复盘", "正文内容", db.StatusPublished)
	cookie := testutil.RegisterUser(t, api, "commenter@example.com", "Commenter")
	testutil.SaveProfile(t, api, cookie, "Commenter", "Fudan University", true)

	commentRec := api.JSON(t, http.MethodPost, "/api/posts/"+uintString(post.ID)+"/comments", map[string]any{
		"body": "先占个楼",
	}, cookie)
	if commentRec.Code != http.StatusCreated {
		t.Fatalf("expected comment create success, got %d: %s", commentRec.Code, commentRec.Body.String())
	}
	parentID := uint(testutil.DecodeJSON(t, commentRec)["id"].(float64))

	replyRec := api.JSON(t, http.MethodPost, "/api/comments/"+uintString(parentID)+"/replies", map[string]any{
		"body": "回复楼主",
	}, cookie)
	if replyRec.Code != http.StatusCreated {
		t.Fatalf("expected reply create success, got %d: %s", replyRec.Code, replyRec.Body.String())
	}
	replyPayload := testutil.DecodeJSON(t, replyRec)
	if replyPayload["parentCommentId"] != float64(parentID) {
		t.Fatalf("expected reply parent id, got %s", replyRec.Body.String())
	}
	counts := replyPayload["counts"].(map[string]any)
	if counts["comments"] != float64(2) {
		t.Fatalf("expected comment count 2 after reply, got %s", replyRec.Body.String())
	}

	replyID := uint(replyPayload["id"].(float64))
	nested := api.JSON(t, http.MethodPost, "/api/comments/"+uintString(replyID)+"/replies", map[string]any{
		"body": "不能嵌套回复",
	}, cookie)
	if nested.Code != http.StatusBadRequest {
		t.Fatalf("expected nested reply validation failure, got %d: %s", nested.Code, nested.Body.String())
	}
	if code := testutil.ErrorCode(t, nested); code != "VALIDATION_FAILED" {
		t.Fatalf("expected VALIDATION_FAILED, got %s", code)
	}

	var updated db.Post
	if err := api.DB.First(&updated, post.ID).Error; err != nil {
		t.Fatalf("load updated post: %v", err)
	}
	if updated.CommentCount != 2 {
		t.Fatalf("nested reply should not increment comments, got %+v", updated)
	}
}

func uintString(value uint) string {
	return strconv.FormatUint(uint64(value), 10)
}
