package admin_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"gradlink/internal/db"
	"gradlink/internal/testutil"
)

func TestAdminCanApprovePendingModerationItem(t *testing.T) {
	api := testutil.NewAPI(t)
	author := testutil.SeedUserWithProfile(t, api, "author@example.com", "Author", "Fudan University", true)
	post := testutil.SeedPost(t, api, author.ID, "有人说可以加微信内推，靠谱吗？", "对方要求加微信并购买培训服务。", db.StatusPending)
	item := db.ModerationQueueItem{
		TargetType:      "post",
		TargetID:        post.ID,
		MatchedKeywords: "加微信",
		Status:          db.StatusPending,
	}
	if err := api.DB.Create(&item).Error; err != nil {
		t.Fatalf("create moderation item: %v", err)
	}
	adminCookie := seedAdminCookie(t, api)

	queue := api.JSON(t, http.MethodGet, "/api/admin/moderation", nil, adminCookie)
	if queue.Code != http.StatusOK {
		t.Fatalf("expected queue 200, got %d: %s", queue.Code, queue.Body.String())
	}
	if !strings.Contains(queue.Body.String(), "加微信") {
		t.Fatalf("expected keyword hit in queue, got %s", queue.Body.String())
	}

	approve := api.JSON(t, http.MethodPost, "/api/admin/moderation/1/approve", nil, adminCookie)
	if approve.Code != http.StatusOK {
		t.Fatalf("expected approve 200, got %d: %s", approve.Code, approve.Body.String())
	}

	var approvedPost db.Post
	if err := api.DB.First(&approvedPost, post.ID).Error; err != nil {
		t.Fatalf("load approved post: %v", err)
	}
	if approvedPost.Status != db.StatusPublished {
		t.Fatalf("expected post published after approval, got %q", approvedPost.Status)
	}
	var approvedItem db.ModerationQueueItem
	if err := api.DB.First(&approvedItem, item.ID).Error; err != nil {
		t.Fatalf("load approved item: %v", err)
	}
	if approvedItem.Status != db.StatusPublished || approvedItem.ReviewerID == nil {
		t.Fatalf("expected approved queue item with reviewer, got %+v", approvedItem)
	}
}

func TestAdminModerationListShowsOnlyPendingContentDetails(t *testing.T) {
	api := testutil.NewAPI(t)
	author := testutil.SeedUserWithProfile(t, api, "author@example.com", "Author", "Fudan University", true)
	pendingPost := testutil.SeedPost(t, api, author.ID, "待审核帖子标题", "待审核帖子正文", db.StatusPending)
	approvedPost := testutil.SeedPost(t, api, author.ID, "已处理帖子标题", "已处理帖子正文", db.StatusPublished)
	pendingItem := db.ModerationQueueItem{TargetType: "post", TargetID: pendingPost.ID, MatchedKeywords: "加微信", Status: db.StatusPending}
	approvedItem := db.ModerationQueueItem{TargetType: "post", TargetID: approvedPost.ID, MatchedKeywords: "包过", Status: db.StatusPublished}
	if err := api.DB.Create(&pendingItem).Error; err != nil {
		t.Fatalf("create pending item: %v", err)
	}
	if err := api.DB.Create(&approvedItem).Error; err != nil {
		t.Fatalf("create approved item: %v", err)
	}
	adminCookie := seedAdminCookie(t, api)

	queue := api.JSON(t, http.MethodGet, "/api/admin/moderation", nil, adminCookie)
	if queue.Code != http.StatusOK {
		t.Fatalf("expected queue 200, got %d: %s", queue.Code, queue.Body.String())
	}
	if !strings.Contains(queue.Body.String(), "待审核帖子标题") || !strings.Contains(queue.Body.String(), "待审核帖子正文") {
		t.Fatalf("expected pending content details in queue, got %s", queue.Body.String())
	}
	if strings.Contains(queue.Body.String(), "已处理帖子标题") {
		t.Fatalf("expected completed moderation item to be omitted, got %s", queue.Body.String())
	}
}

func TestAdminApprovingPendingCommentUpdatesPostCounts(t *testing.T) {
	api := testutil.NewAPI(t)
	author := testutil.SeedUserWithProfile(t, api, "author@example.com", "Author", "Fudan University", true)
	commenter := testutil.SeedUserWithProfile(t, api, "commenter@example.com", "Commenter", "Nanjing University", true)
	post := testutil.SeedPost(t, api, author.ID, "面试复盘", "正文内容", db.StatusPublished)
	comment := db.Comment{PostID: post.ID, AuthorID: commenter.ID, Body: "建议加微信聊", Status: db.StatusPending}
	if err := api.DB.Create(&comment).Error; err != nil {
		t.Fatalf("create pending comment: %v", err)
	}
	item := db.ModerationQueueItem{TargetType: "comment", TargetID: comment.ID, MatchedKeywords: "加微信", Status: db.StatusPending}
	if err := api.DB.Create(&item).Error; err != nil {
		t.Fatalf("create moderation item: %v", err)
	}
	adminCookie := seedAdminCookie(t, api)

	approve := api.JSON(t, http.MethodPost, "/api/admin/moderation/"+uintString(item.ID)+"/approve", nil, adminCookie)
	if approve.Code != http.StatusOK {
		t.Fatalf("expected approve 200, got %d: %s", approve.Code, approve.Body.String())
	}

	var updated db.Post
	if err := api.DB.First(&updated, post.ID).Error; err != nil {
		t.Fatalf("load updated post: %v", err)
	}
	if updated.CommentCount != 1 || updated.HotScore != 3 {
		t.Fatalf("expected approved comment to update counts, got %+v", updated)
	}
}

func TestAdminModerationDecisionFailsForMissingTarget(t *testing.T) {
	api := testutil.NewAPI(t)
	item := db.ModerationQueueItem{TargetType: "post", TargetID: 999, MatchedKeywords: "加微信", Status: db.StatusPending}
	if err := api.DB.Create(&item).Error; err != nil {
		t.Fatalf("create moderation item: %v", err)
	}
	adminCookie := seedAdminCookie(t, api)

	approve := api.JSON(t, http.MethodPost, "/api/admin/moderation/"+uintString(item.ID)+"/approve", nil, adminCookie)
	if approve.Code != http.StatusNotFound {
		t.Fatalf("expected missing target 404, got %d: %s", approve.Code, approve.Body.String())
	}
}

func TestAdminEndpointsRequireAdminRole(t *testing.T) {
	api := testutil.NewAPI(t)

	unauthenticated := api.JSON(t, http.MethodGet, "/api/admin/moderation", nil)
	if unauthenticated.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthenticated admin request 401, got %d: %s", unauthenticated.Code, unauthenticated.Body.String())
	}

	userCookie := testutil.RegisterUser(t, api, "student@example.com", "Student")
	forbidden := api.JSON(t, http.MethodGet, "/api/admin/moderation", nil, userCookie)
	if forbidden.Code != http.StatusForbidden {
		t.Fatalf("expected non-admin request 403, got %d: %s", forbidden.Code, forbidden.Body.String())
	}
}

func TestAdminCanListAndResolveReports(t *testing.T) {
	api := testutil.NewAPI(t)
	reporter := testutil.SeedUserWithProfile(t, api, "reporter@example.com", "Reporter", "Fudan University", true)
	author := testutil.SeedUserWithProfile(t, api, "author@example.com", "Author", "Nanjing University", true)
	post := testutil.SeedPost(t, api, author.ID, "广告帖", "疑似广告内容", db.StatusPublished)
	report := db.Report{
		ReporterID: reporter.ID,
		TargetType: "post",
		TargetID:   post.ID,
		Reason:     "广告引流",
		Detail:     "疑似引导加微信",
		Status:     db.StatusPending,
	}
	if err := api.DB.Create(&report).Error; err != nil {
		t.Fatalf("create report: %v", err)
	}
	adminCookie := seedAdminCookie(t, api)

	list := api.JSON(t, http.MethodGet, "/api/admin/reports", nil, adminCookie)
	if list.Code != http.StatusOK {
		t.Fatalf("expected reports 200, got %d: %s", list.Code, list.Body.String())
	}
	if !strings.Contains(list.Body.String(), "广告引流") {
		t.Fatalf("expected report reason in response, got %s", list.Body.String())
	}

	resolve := api.JSON(t, http.MethodPost, "/api/admin/reports/"+uintString(report.ID)+"/resolve", nil, adminCookie)
	if resolve.Code != http.StatusOK {
		t.Fatalf("expected resolve 200, got %d: %s", resolve.Code, resolve.Body.String())
	}

	var resolved db.Report
	if err := api.DB.First(&resolved, report.ID).Error; err != nil {
		t.Fatalf("load resolved report: %v", err)
	}
	if resolved.Status != db.ReportStatusResolved {
		t.Fatalf("expected resolved report status, got %q", resolved.Status)
	}
}

func TestAdminCanBanAndUnbanUsers(t *testing.T) {
	api := testutil.NewAPI(t)
	user := testutil.SeedUserWithProfile(t, api, "student@example.com", "Student", "Fudan University", true)
	adminCookie := seedAdminCookie(t, api)

	list := api.JSON(t, http.MethodGet, "/api/admin/users", nil, adminCookie)
	if list.Code != http.StatusOK {
		t.Fatalf("expected users 200, got %d: %s", list.Code, list.Body.String())
	}
	if !strings.Contains(list.Body.String(), "student@example.com") {
		t.Fatalf("expected user email in list, got %s", list.Body.String())
	}

	ban := api.JSON(t, http.MethodPost, "/api/admin/users/"+uintString(user.ID)+"/ban", nil, adminCookie)
	if ban.Code != http.StatusOK {
		t.Fatalf("expected ban 200, got %d: %s", ban.Code, ban.Body.String())
	}
	assertUserStatus(t, api, user.ID, db.UserStatusDisabled)

	unban := api.JSON(t, http.MethodPost, "/api/admin/users/"+uintString(user.ID)+"/unban", nil, adminCookie)
	if unban.Code != http.StatusOK {
		t.Fatalf("expected unban 200, got %d: %s", unban.Code, unban.Body.String())
	}
	assertUserStatus(t, api, user.ID, db.UserStatusActive)
}

func TestAdminCanManageKeywords(t *testing.T) {
	api := testutil.NewAPI(t)
	adminCookie := seedAdminCookie(t, api)

	create := api.JSON(t, http.MethodPost, "/api/admin/keywords", map[string]any{
		"keyword":   "培训包过",
		"riskLevel": "high",
		"action":    db.StatusPending,
		"enabled":   true,
	}, adminCookie)
	if create.Code != http.StatusCreated {
		t.Fatalf("expected keyword create 201, got %d: %s", create.Code, create.Body.String())
	}
	created := testutil.DecodeJSON(t, create)
	keywordID := uint(created["id"].(float64))

	update := api.JSON(t, http.MethodPut, "/api/admin/keywords/"+uintString(keywordID), map[string]any{
		"keyword":   "培训包过",
		"riskLevel": "high",
		"action":    db.StatusRejected,
		"enabled":   false,
	}, adminCookie)
	if update.Code != http.StatusOK {
		t.Fatalf("expected keyword update 200, got %d: %s", update.Code, update.Body.String())
	}

	list := api.JSON(t, http.MethodGet, "/api/admin/keywords", nil, adminCookie)
	if list.Code != http.StatusOK {
		t.Fatalf("expected keywords 200, got %d: %s", list.Code, list.Body.String())
	}
	if !strings.Contains(list.Body.String(), db.StatusRejected) {
		t.Fatalf("expected updated keyword action in list, got %s", list.Body.String())
	}

	del := api.JSON(t, http.MethodDelete, "/api/admin/keywords/"+uintString(keywordID), nil, adminCookie)
	if del.Code != http.StatusOK {
		t.Fatalf("expected keyword delete 200, got %d: %s", del.Code, del.Body.String())
	}
	var count int64
	if err := api.DB.Model(&db.ModerationKeyword{}).Where("id = ?", keywordID).Count(&count).Error; err != nil {
		t.Fatalf("count keyword: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected keyword deleted, got count %d", count)
	}
}

func TestAdminKeywordCreateDefaultsToEnabled(t *testing.T) {
	api := testutil.NewAPI(t)
	adminCookie := seedAdminCookie(t, api)

	create := api.JSON(t, http.MethodPost, "/api/admin/keywords", map[string]any{
		"keyword":   "内推收费",
		"riskLevel": "medium",
		"action":    db.StatusPending,
	}, adminCookie)
	if create.Code != http.StatusCreated {
		t.Fatalf("expected keyword create 201, got %d: %s", create.Code, create.Body.String())
	}
	payload := testutil.DecodeJSON(t, create)
	if payload["enabled"] != true {
		t.Fatalf("expected omitted enabled to default true, got %s", create.Body.String())
	}
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

func assertUserStatus(t *testing.T, api *testutil.TestAPI, userID uint, status string) {
	t.Helper()

	var user db.User
	if err := api.DB.First(&user, userID).Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	if user.Status != status {
		t.Fatalf("expected user status %q, got %q", status, user.Status)
	}
}

func uintString(value uint) string {
	return strconv.FormatUint(uint64(value), 10)
}
