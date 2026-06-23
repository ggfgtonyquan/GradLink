package comments_test

import (
	"net/http"
	"testing"

	"gradlink/internal/db"
	"gradlink/internal/testutil"
)

func TestListPostCommentsReturnsPublishedCommentsWithAuthorProfile(t *testing.T) {
	api := testutil.NewAPI(t)
	author := testutil.SeedUserWithProfile(t, api, "author@example.com", "Author", "Fudan University", true)
	commenterCookie := testutil.RegisterUser(t, api, "commenter@example.com", "Commenter")
	testutil.SaveProfile(t, api, commenterCookie, "Commenter", "Shanghai Jiao Tong University", false)
	var commenter db.User
	if err := api.DB.Where("email = ?", "commenter@example.com").First(&commenter).Error; err != nil {
		t.Fatalf("load commenter: %v", err)
	}
	post := testutil.SeedPost(t, api, author.ID, "面试复盘", "正文内容", db.StatusPublished)

	published := db.Comment{PostID: post.ID, AuthorID: commenter.ID, Body: "这条评论应该展示", Status: db.StatusPublished}
	pending := db.Comment{PostID: post.ID, AuthorID: commenter.ID, Body: "待审核不展示", Status: db.StatusPending}
	if err := api.DB.Create(&published).Error; err != nil {
		t.Fatalf("create published comment: %v", err)
	}
	if err := api.DB.Create(&pending).Error; err != nil {
		t.Fatalf("create pending comment: %v", err)
	}

	rec := api.JSON(t, http.MethodGet, "/api/posts/"+uintString(post.ID)+"/comments", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected comments success, got %d: %s", rec.Code, rec.Body.String())
	}

	payload := testutil.DecodeJSON(t, rec)
	comments := payload["comments"].([]any)
	if len(comments) != 1 {
		t.Fatalf("expected 1 published comment, got %d: %s", len(comments), rec.Body.String())
	}
	comment := comments[0].(map[string]any)
	if comment["body"] != "这条评论应该展示" {
		t.Fatalf("expected published comment body, got %s", rec.Body.String())
	}
	authorPayload := comment["author"].(map[string]any)
	if authorPayload["nickname"] != "Commenter" || authorPayload["school"] != nil {
		t.Fatalf("expected masked author profile, got %s", rec.Body.String())
	}
}
