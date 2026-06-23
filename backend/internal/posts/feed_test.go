package posts_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"gradlink/internal/db"
	"gradlink/internal/testutil"
)

func TestFeedDefaultsToLatestPublishedPosts(t *testing.T) {
	api := testutil.NewAPI(t)
	tag := seedTag(t, api, "综合提问", "general-questions", 90)
	author := seedUserWithProfile(t, api, "author@example.com", "Author", "Nanjing University", true)

	old := seedPost(t, api, author.ID, "较早帖子", "这是一篇较早发布的内容", db.StatusPublished, time.Now().Add(-2*time.Hour), 0)
	newer := seedPost(t, api, author.ID, "较新帖子", "这是一篇较新发布的内容", db.StatusPublished, time.Now().Add(-1*time.Hour), 0)
	hidden := seedPost(t, api, author.ID, "隐藏帖子", "这篇内容不应该出现在列表", db.StatusHidden, time.Now(), 999)
	attachTags(t, api, old.ID, tag.ID)
	attachTags(t, api, newer.ID, tag.ID)
	attachTags(t, api, hidden.ID, tag.ID)

	rec := api.JSON(t, http.MethodGet, "/api/posts", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected feed success, got %d: %s", rec.Code, rec.Body.String())
	}
	posts := responsePosts(t, rec)
	if len(posts) != 2 {
		t.Fatalf("expected two published posts, got %d: %s", len(posts), rec.Body.String())
	}
	if posts[0]["title"] != newer.Title || posts[1]["title"] != old.Title {
		t.Fatalf("expected latest order, got %s", rec.Body.String())
	}
	if _, ok := posts[0]["excerpt"].(string); !ok {
		t.Fatalf("expected excerpt in feed item, got %s", rec.Body.String())
	}
	if _, hasBody := posts[0]["body"]; hasBody {
		t.Fatalf("feed item should not include full body: %s", rec.Body.String())
	}
}

func TestFeedSupportsHotSortTagFilterAndQuery(t *testing.T) {
	api := testutil.NewAPI(t)
	interview := seedTag(t, api, "面试经验", "interview-experience", 10)
	resume := seedTag(t, api, "简历求助", "resume-help", 50)
	author := seedUserWithProfile(t, api, "author@example.com", "Author", "Zhejiang University", true)

	hot := seedPost(t, api, author.ID, "字节面试复盘", "算法题和项目追问都比较深入", db.StatusPublished, time.Now().Add(-3*time.Hour), 50)
	warm := seedPost(t, api, author.ID, "腾讯面试记录", "项目经历被问得很细", db.StatusPublished, time.Now().Add(-1*time.Hour), 10)
	resumePost := seedPost(t, api, author.ID, "简历修改求建议", "希望大家帮我看看实习描述", db.StatusPublished, time.Now(), 99)
	attachTags(t, api, hot.ID, interview.ID)
	attachTags(t, api, warm.ID, interview.ID)
	attachTags(t, api, resumePost.ID, resume.ID)

	hotRec := api.JSON(t, http.MethodGet, "/api/posts?sort=hot&tagId="+uintString(interview.ID), nil)
	if hotRec.Code != http.StatusOK {
		t.Fatalf("expected hot feed success, got %d: %s", hotRec.Code, hotRec.Body.String())
	}
	hotPosts := responsePosts(t, hotRec)
	if len(hotPosts) != 2 {
		t.Fatalf("expected two interview posts, got %d: %s", len(hotPosts), hotRec.Body.String())
	}
	if hotPosts[0]["title"] != hot.Title || hotPosts[1]["title"] != warm.Title {
		t.Fatalf("expected hot score order within tag, got %s", hotRec.Body.String())
	}

	queryRec := api.JSON(t, http.MethodGet, "/api/posts?q=实习描述", nil)
	if queryRec.Code != http.StatusOK {
		t.Fatalf("expected query feed success, got %d: %s", queryRec.Code, queryRec.Body.String())
	}
	queryPosts := responsePosts(t, queryRec)
	if len(queryPosts) != 1 || queryPosts[0]["title"] != resumePost.Title {
		t.Fatalf("expected body query to find resume post, got %s", queryRec.Body.String())
	}
}

func TestFeedUsesDefaultLimitAndSupportsOffset(t *testing.T) {
	api := testutil.NewAPI(t)
	tag := seedTag(t, api, "综合提问", "general-questions", 90)
	author := seedUserWithProfile(t, api, "author@example.com", "Author", "Nanjing University", true)

	for i := 0; i < 25; i++ {
		post := seedPost(t, api, author.ID, "分页帖子 "+uintString(uint(i)), "这是一篇用于分页测试的内容", db.StatusPublished, time.Now().Add(time.Duration(i)*time.Minute), i)
		attachTags(t, api, post.ID, tag.ID)
	}

	firstPage := api.JSON(t, http.MethodGet, "/api/posts", nil)
	if firstPage.Code != http.StatusOK {
		t.Fatalf("expected feed success, got %d: %s", firstPage.Code, firstPage.Body.String())
	}
	firstPosts := responsePosts(t, firstPage)
	if len(firstPosts) != 20 {
		t.Fatalf("expected default limit 20, got %d", len(firstPosts))
	}

	secondPage := api.JSON(t, http.MethodGet, "/api/posts?limit=10&offset=20", nil)
	if secondPage.Code != http.StatusOK {
		t.Fatalf("expected second page success, got %d: %s", secondPage.Code, secondPage.Body.String())
	}
	secondPosts := responsePosts(t, secondPage)
	if len(secondPosts) != 5 {
		t.Fatalf("expected remaining 5 posts, got %d: %s", len(secondPosts), secondPage.Body.String())
	}

	invalidLimit := api.JSON(t, http.MethodGet, "/api/posts?limit=101", nil)
	if invalidLimit.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid limit failure, got %d: %s", invalidLimit.Code, invalidLimit.Body.String())
	}
	if code := testutil.ErrorCode(t, invalidLimit); code != "VALIDATION_FAILED" {
		t.Fatalf("expected VALIDATION_FAILED, got %s", code)
	}
}

func TestFeedSearchTreatsLikeWildcardsLiterally(t *testing.T) {
	api := testutil.NewAPI(t)
	tag := seedTag(t, api, "综合提问", "general-questions", 90)
	author := seedUserWithProfile(t, api, "author@example.com", "Author", "Nanjing University", true)

	percent := seedPost(t, api, author.ID, "薪资涨幅 20% 怎么看", "正文里也有百分号 20%", db.StatusPublished, time.Now(), 0)
	underscore := seedPost(t, api, author.ID, "岗位 JD_后端", "包含下划线的职位描述", db.StatusPublished, time.Now().Add(-time.Minute), 0)
	ordinary := seedPost(t, api, author.ID, "普通问题", "没有特殊符号", db.StatusPublished, time.Now().Add(-2*time.Minute), 0)
	attachTags(t, api, percent.ID, tag.ID)
	attachTags(t, api, underscore.ID, tag.ID)
	attachTags(t, api, ordinary.ID, tag.ID)

	percentRec := api.JSON(t, http.MethodGet, "/api/posts?q="+url.QueryEscape("%"), nil)
	if percentRec.Code != http.StatusOK {
		t.Fatalf("expected percent query success, got %d: %s", percentRec.Code, percentRec.Body.String())
	}
	percentPosts := responsePosts(t, percentRec)
	if len(percentPosts) != 1 || percentPosts[0]["title"] != percent.Title {
		t.Fatalf("expected literal percent search, got %s", percentRec.Body.String())
	}

	underscoreRec := api.JSON(t, http.MethodGet, "/api/posts?q="+url.QueryEscape("_"), nil)
	if underscoreRec.Code != http.StatusOK {
		t.Fatalf("expected underscore query success, got %d: %s", underscoreRec.Code, underscoreRec.Body.String())
	}
	underscorePosts := responsePosts(t, underscoreRec)
	if len(underscorePosts) != 1 || underscorePosts[0]["title"] != underscore.Title {
		t.Fatalf("expected literal underscore search, got %s", underscoreRec.Body.String())
	}
}
