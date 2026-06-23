package reactions_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"gradlink/internal/db"
	"gradlink/internal/testutil"
)

func TestLikePostIsIdempotentAndUpdatesCounts(t *testing.T) {
	api := testutil.NewAPI(t)
	author := testutil.SeedUserWithProfile(t, api, "author@example.com", "Author", "Nanjing University", true)
	post := testutil.SeedPost(t, api, author.ID, "面试复盘", "正文内容", db.StatusPublished)
	cookie := testutil.RegisterUser(t, api, "reader@example.com", "Reader")
	testutil.SaveProfile(t, api, cookie, "Reader", "Fudan University", true)

	first := api.JSON(t, http.MethodPost, "/api/posts/"+uintString(post.ID)+"/like", nil, cookie)
	if first.Code != http.StatusOK {
		t.Fatalf("expected like success, got %d: %s", first.Code, first.Body.String())
	}
	assertReactionCounts(t, first, 1, 0)

	second := api.JSON(t, http.MethodPost, "/api/posts/"+uintString(post.ID)+"/like", nil, cookie)
	if second.Code != http.StatusOK {
		t.Fatalf("expected repeated like success, got %d: %s", second.Code, second.Body.String())
	}
	assertReactionCounts(t, second, 1, 0)

	var likeRows int64
	if err := api.DB.Model(&db.Like{}).Where("target_type = ? AND target_id = ?", "post", post.ID).Count(&likeRows).Error; err != nil {
		t.Fatalf("count likes: %v", err)
	}
	if likeRows != 1 {
		t.Fatalf("expected one like row, got %d", likeRows)
	}

	var liked db.Post
	if err := api.DB.First(&liked, post.ID).Error; err != nil {
		t.Fatalf("load liked post: %v", err)
	}
	if liked.LikeCount != 1 || liked.HotScore <= 0 {
		t.Fatalf("expected like count and hot score to update, got %+v", liked)
	}

	firstDelete := api.JSON(t, http.MethodDelete, "/api/posts/"+uintString(post.ID)+"/like", nil, cookie)
	if firstDelete.Code != http.StatusOK {
		t.Fatalf("expected unlike success, got %d: %s", firstDelete.Code, firstDelete.Body.String())
	}
	assertReactionCounts(t, firstDelete, 0, 0)

	secondDelete := api.JSON(t, http.MethodDelete, "/api/posts/"+uintString(post.ID)+"/like", nil, cookie)
	if secondDelete.Code != http.StatusOK {
		t.Fatalf("expected repeated unlike success, got %d: %s", secondDelete.Code, secondDelete.Body.String())
	}
	assertReactionCounts(t, secondDelete, 0, 0)

	var unliked db.Post
	if err := api.DB.First(&unliked, post.ID).Error; err != nil {
		t.Fatalf("load unliked post: %v", err)
	}
	if unliked.LikeCount != 0 || unliked.HotScore != 0 {
		t.Fatalf("expected unlike to restore counts and hot score, got %+v", unliked)
	}
}

func TestFavoritePostIsIdempotentAndUpdatesCounts(t *testing.T) {
	api := testutil.NewAPI(t)
	author := testutil.SeedUserWithProfile(t, api, "author@example.com", "Author", "Nanjing University", true)
	post := testutil.SeedPost(t, api, author.ID, "面试复盘", "正文内容", db.StatusPublished)
	cookie := testutil.RegisterUser(t, api, "reader@example.com", "Reader")
	testutil.SaveProfile(t, api, cookie, "Reader", "Fudan University", true)

	first := api.JSON(t, http.MethodPost, "/api/posts/"+uintString(post.ID)+"/favorite", nil, cookie)
	if first.Code != http.StatusOK {
		t.Fatalf("expected favorite success, got %d: %s", first.Code, first.Body.String())
	}
	assertReactionCounts(t, first, 0, 1)

	second := api.JSON(t, http.MethodPost, "/api/posts/"+uintString(post.ID)+"/favorite", nil, cookie)
	if second.Code != http.StatusOK {
		t.Fatalf("expected repeated favorite success, got %d: %s", second.Code, second.Body.String())
	}
	assertReactionCounts(t, second, 0, 1)

	var favoriteRows int64
	if err := api.DB.Model(&db.Favorite{}).Where("post_id = ?", post.ID).Count(&favoriteRows).Error; err != nil {
		t.Fatalf("count favorites: %v", err)
	}
	if favoriteRows != 1 {
		t.Fatalf("expected one favorite row, got %d", favoriteRows)
	}

	firstDelete := api.JSON(t, http.MethodDelete, "/api/posts/"+uintString(post.ID)+"/favorite", nil, cookie)
	if firstDelete.Code != http.StatusOK {
		t.Fatalf("expected unfavorite success, got %d: %s", firstDelete.Code, firstDelete.Body.String())
	}
	assertReactionCounts(t, firstDelete, 0, 0)

	secondDelete := api.JSON(t, http.MethodDelete, "/api/posts/"+uintString(post.ID)+"/favorite", nil, cookie)
	if secondDelete.Code != http.StatusOK {
		t.Fatalf("expected repeated unfavorite success, got %d: %s", secondDelete.Code, secondDelete.Body.String())
	}
	assertReactionCounts(t, secondDelete, 0, 0)
}

func TestReactionsRequireAuthProfileAndPublishedPost(t *testing.T) {
	api := testutil.NewAPI(t)
	author := testutil.SeedUserWithProfile(t, api, "author@example.com", "Author", "Nanjing University", true)
	hidden := testutil.SeedPost(t, api, author.ID, "隐藏帖", "正文内容", db.StatusHidden)

	unauthenticated := api.JSON(t, http.MethodPost, "/api/posts/"+uintString(hidden.ID)+"/like", nil)
	if unauthenticated.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthenticated like to fail with 401, got %d: %s", unauthenticated.Code, unauthenticated.Body.String())
	}
	if code := testutil.ErrorCode(t, unauthenticated); code != "UNAUTHENTICATED" {
		t.Fatalf("expected UNAUTHENTICATED, got %s", code)
	}

	cookie := testutil.RegisterUser(t, api, "reader@example.com", "Reader")
	incomplete := api.JSON(t, http.MethodPost, "/api/posts/"+uintString(hidden.ID)+"/favorite", nil, cookie)
	if incomplete.Code != http.StatusForbidden {
		t.Fatalf("expected incomplete profile to fail with 403, got %d: %s", incomplete.Code, incomplete.Body.String())
	}
	if code := testutil.ErrorCode(t, incomplete); code != "PROFILE_REQUIRED" {
		t.Fatalf("expected PROFILE_REQUIRED, got %s", code)
	}

	testutil.SaveProfile(t, api, cookie, "Reader", "Fudan University", true)
	notFound := api.JSON(t, http.MethodPost, "/api/posts/"+uintString(hidden.ID)+"/like", nil, cookie)
	if notFound.Code != http.StatusNotFound {
		t.Fatalf("expected hidden post to be not found, got %d: %s", notFound.Code, notFound.Body.String())
	}
	if code := testutil.ErrorCode(t, notFound); code != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND, got %s", code)
	}
}

func assertReactionCounts(t *testing.T, rec *httptest.ResponseRecorder, likes, favorites int) {
	t.Helper()

	payload := testutil.DecodeJSON(t, rec)
	counts, ok := payload["counts"].(map[string]any)
	if !ok {
		t.Fatalf("expected counts object, got %s", rec.Body.String())
	}
	if counts["likes"] != float64(likes) || counts["favorites"] != float64(favorites) {
		t.Fatalf("expected likes=%d favorites=%d, got %s", likes, favorites, rec.Body.String())
	}
}

func uintString(value uint) string {
	return strconv.FormatUint(uint64(value), 10)
}
