package posts_test

import (
	"net/http"
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
