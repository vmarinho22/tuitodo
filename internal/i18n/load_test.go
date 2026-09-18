package i18n

import "testing"

func TestFlattenMessages(t *testing.T) {
	nested := map[string]any{
		"error": map[string]any{
			"emptyTitle": "Title cannot be empty.",
			"unexpected": "Something went wrong.",
		},
		"form": map[string]any{
			"title": "Title",
		},
		"tasks": map[string]any{
			"title": "Tasks",
		},
	}
	got := flattenMessages(nested, "")
	if got["error.emptyTitle"] != "Title cannot be empty." {
		t.Fatalf("error.emptyTitle = %q", got["error.emptyTitle"])
	}
	if got["form.title"] != "Title" {
		t.Fatalf("form.title = %q", got["form.title"])
	}
	if got["tasks.title"] != "Tasks" {
		t.Fatalf("tasks.title = %q", got["tasks.title"])
	}
	if len(got) != 4 {
		t.Fatalf("len = %d, want 4", len(got))
	}
}
