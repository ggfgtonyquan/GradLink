package moderation_test

import (
	"net/http"
	"testing"

	"gradlink/internal/db"
	"gradlink/internal/testutil"
)

func TestCreateReportForPublishedPost(t *testing.T) {
	api := testutil.NewAPI(t)
	author := testutil.SeedUserWithProfile(t, api, "author@example.com", "Author", "Nanjing University", true)
	post := testutil.SeedPost(t, api, author.ID, "面试复盘", "正文内容", db.StatusPublished)
	cookie := testutil.RegisterUser(t, api, "reporter@example.com", "Reporter")
	testutil.SaveProfile(t, api, cookie, "Reporter", "Fudan University", true)

	rec := api.JSON(t, http.MethodPost, "/api/reports", map[string]any{
		"targetType": "post",
		"targetId":   post.ID,
		"reason":     "spam",
		"detail":     "重复发布广告内容",
	}, cookie)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected report create success, got %d: %s", rec.Code, rec.Body.String())
	}

	payload := testutil.DecodeJSON(t, rec)
	if payload["targetType"] != "post" || payload["status"] != db.StatusPending {
		t.Fatalf("unexpected report payload: %s", rec.Body.String())
	}

	var stored db.Report
	if err := api.DB.First(&stored, uint(payload["id"].(float64))).Error; err != nil {
		t.Fatalf("load created report: %v", err)
	}
	if stored.TargetType != "post" || stored.TargetID != post.ID || stored.Reason != "spam" || stored.Status != db.StatusPending {
		t.Fatalf("report was not saved correctly: %+v", stored)
	}
}

func TestCreateReportForPublishedComment(t *testing.T) {
	api := testutil.NewAPI(t)
	author := testutil.SeedUserWithProfile(t, api, "author@example.com", "Author", "Nanjing University", true)
	post := testutil.SeedPost(t, api, author.ID, "面试复盘", "正文内容", db.StatusPublished)
	comment := db.Comment{
		PostID:   post.ID,
		AuthorID: author.ID,
		Body:     "评论内容",
		Status:   db.StatusPublished,
	}
	if err := api.DB.Create(&comment).Error; err != nil {
		t.Fatalf("create comment: %v", err)
	}
	cookie := testutil.RegisterUser(t, api, "reporter@example.com", "Reporter")
	testutil.SaveProfile(t, api, cookie, "Reporter", "Fudan University", true)

	rec := api.JSON(t, http.MethodPost, "/api/reports", map[string]any{
		"targetType": "comment",
		"targetId":   comment.ID,
		"reason":     "abuse",
		"detail":     "攻击性评论",
	}, cookie)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected report create success, got %d: %s", rec.Code, rec.Body.String())
	}
	payload := testutil.DecodeJSON(t, rec)
	if payload["targetType"] != "comment" || payload["targetId"] != float64(comment.ID) {
		t.Fatalf("unexpected report payload: %s", rec.Body.String())
	}
}

func TestCreateReportRequiresAuthProfileAndValidVisibleTarget(t *testing.T) {
	api := testutil.NewAPI(t)
	author := testutil.SeedUserWithProfile(t, api, "author@example.com", "Author", "Nanjing University", true)
	hidden := testutil.SeedPost(t, api, author.ID, "隐藏帖", "正文内容", db.StatusHidden)

	unauthenticated := api.JSON(t, http.MethodPost, "/api/reports", map[string]any{
		"targetType": "post",
		"targetId":   hidden.ID,
		"reason":     "spam",
	})
	if unauthenticated.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthenticated report to fail with 401, got %d: %s", unauthenticated.Code, unauthenticated.Body.String())
	}
	if code := testutil.ErrorCode(t, unauthenticated); code != "UNAUTHENTICATED" {
		t.Fatalf("expected UNAUTHENTICATED, got %s", code)
	}

	cookie := testutil.RegisterUser(t, api, "reporter@example.com", "Reporter")
	incomplete := api.JSON(t, http.MethodPost, "/api/reports", map[string]any{
		"targetType": "post",
		"targetId":   hidden.ID,
		"reason":     "spam",
	}, cookie)
	if incomplete.Code != http.StatusForbidden {
		t.Fatalf("expected incomplete profile to fail with 403, got %d: %s", incomplete.Code, incomplete.Body.String())
	}
	if code := testutil.ErrorCode(t, incomplete); code != "PROFILE_REQUIRED" {
		t.Fatalf("expected PROFILE_REQUIRED, got %s", code)
	}

	testutil.SaveProfile(t, api, cookie, "Reporter", "Fudan University", true)
	invalid := api.JSON(t, http.MethodPost, "/api/reports", map[string]any{
		"targetType": "user",
		"targetId":   hidden.ID,
		"reason":     "spam",
	}, cookie)
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("expected validation failure, got %d: %s", invalid.Code, invalid.Body.String())
	}
	if code := testutil.ErrorCode(t, invalid); code != "VALIDATION_FAILED" {
		t.Fatalf("expected VALIDATION_FAILED, got %s", code)
	}

	notFound := api.JSON(t, http.MethodPost, "/api/reports", map[string]any{
		"targetType": "post",
		"targetId":   hidden.ID,
		"reason":     "spam",
	}, cookie)
	if notFound.Code != http.StatusNotFound {
		t.Fatalf("expected hidden post to be not found, got %d: %s", notFound.Code, notFound.Body.String())
	}
	if code := testutil.ErrorCode(t, notFound); code != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND, got %s", code)
	}

	blankReason := api.JSON(t, http.MethodPost, "/api/reports", map[string]any{
		"targetType": "post",
		"targetId":   hidden.ID,
		"reason":     "   ",
	}, cookie)
	if blankReason.Code != http.StatusBadRequest {
		t.Fatalf("expected blank reason validation failure, got %d: %s", blankReason.Code, blankReason.Body.String())
	}
}
