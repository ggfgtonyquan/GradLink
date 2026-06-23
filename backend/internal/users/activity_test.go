package users_test

import (
	"net/http"
	"testing"

	"gradlink/internal/db"
	"gradlink/internal/testutil"
)

func TestMePostsReturnsCurrentUsersPostsIncludingPending(t *testing.T) {
	api := testutil.NewAPI(t)
	cookie := testutil.RegisterUser(t, api, "student@example.com", "Student")
	testutil.SaveProfile(t, api, cookie, "Student", "Fudan University", true)

	var currentUser db.User
	if err := api.DB.Where("email = ?", "student@example.com").First(&currentUser).Error; err != nil {
		t.Fatalf("load current user: %v", err)
	}
	other := testutil.SeedUserWithProfile(t, api, "other@example.com", "Other", "Nanjing University", true)
	testutil.SeedPost(t, api, currentUser.ID, "我的已发布帖子", "正文", db.StatusPublished)
	testutil.SeedPost(t, api, currentUser.ID, "我的待审核帖子", "正文", db.StatusPending)
	testutil.SeedPost(t, api, other.ID, "别人的帖子", "正文", db.StatusPublished)

	rec := api.JSON(t, http.MethodGet, "/api/me/posts", nil, cookie)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected own posts success, got %d: %s", rec.Code, rec.Body.String())
	}

	payload := testutil.DecodeJSON(t, rec)
	posts := payload["posts"].([]any)
	if len(posts) != 2 {
		t.Fatalf("expected 2 own posts, got %d: %s", len(posts), rec.Body.String())
	}
}

func TestMeFavoritesReturnsCurrentUsersFavoritePosts(t *testing.T) {
	api := testutil.NewAPI(t)
	cookie := testutil.RegisterUser(t, api, "student@example.com", "Student")
	testutil.SaveProfile(t, api, cookie, "Student", "Fudan University", true)

	var currentUser db.User
	if err := api.DB.Where("email = ?", "student@example.com").First(&currentUser).Error; err != nil {
		t.Fatalf("load current user: %v", err)
	}
	author := testutil.SeedUserWithProfile(t, api, "author@example.com", "Author", "Nanjing University", true)
	favoritePost := testutil.SeedPost(t, api, author.ID, "收藏的帖子", "正文", db.StatusPublished)
	hiddenFavorite := testutil.SeedPost(t, api, author.ID, "隐藏收藏", "正文", db.StatusHidden)
	testutil.SeedPost(t, api, author.ID, "未收藏帖子", "正文", db.StatusPublished)
	if err := api.DB.Create(&db.Favorite{UserID: currentUser.ID, PostID: favoritePost.ID}).Error; err != nil {
		t.Fatalf("favorite post: %v", err)
	}
	if err := api.DB.Create(&db.Favorite{UserID: currentUser.ID, PostID: hiddenFavorite.ID}).Error; err != nil {
		t.Fatalf("favorite hidden post: %v", err)
	}

	rec := api.JSON(t, http.MethodGet, "/api/me/favorites", nil, cookie)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected favorites success, got %d: %s", rec.Code, rec.Body.String())
	}

	payload := testutil.DecodeJSON(t, rec)
	posts := payload["posts"].([]any)
	if len(posts) != 1 {
		t.Fatalf("expected 1 visible favorite, got %d: %s", len(posts), rec.Body.String())
	}
	post := posts[0].(map[string]any)
	if post["title"] != "收藏的帖子" {
		t.Fatalf("expected favorite post title, got %s", rec.Body.String())
	}
}
