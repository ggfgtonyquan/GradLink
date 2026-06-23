package tags_test

import (
	"net/http"
	"testing"

	"gradlink/internal/db"
	"gradlink/internal/testutil"
)

func TestGetTagsReturnsSeededTagsInDisplayOrder(t *testing.T) {
	api := testutil.NewAPI(t)
	if err := db.Seed(api.DB); err != nil {
		t.Fatalf("seed database: %v", err)
	}

	rec := api.JSON(t, http.MethodGet, "/api/tags", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected tags success, got %d: %s", rec.Code, rec.Body.String())
	}

	payload := testutil.DecodeJSON(t, rec)
	rawTags, ok := payload["tags"].([]any)
	if !ok {
		t.Fatalf("expected tags array, got %s", rec.Body.String())
	}
	if len(rawTags) != len(db.InitialTags) {
		t.Fatalf("expected %d tags, got %d: %s", len(db.InitialTags), len(rawTags), rec.Body.String())
	}

	for i, raw := range rawTags {
		tag, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("expected tag object at %d, got %T", i, raw)
		}
		want := db.InitialTags[i]
		displayOrder, ok := tag["displayOrder"].(float64)
		if !ok {
			t.Fatalf("expected displayOrder number at %d, got %+v", i, tag)
		}
		if tag["name"] != want.Name || tag["slug"] != want.Slug || int(displayOrder) != want.DisplayOrder {
			t.Fatalf("tag %d not in seeded order: got %+v want %+v", i, tag, want)
		}
		if tag["id"].(float64) == 0 {
			t.Fatalf("expected tag id at %d, got %+v", i, tag)
		}
	}
}
