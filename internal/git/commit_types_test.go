package git

import "testing"

func TestCommitTypeOf(t *testing.T) {
	cases := map[string]string{
		"feat: add dashboard":           "Features",
		"feat(ui): polish cards":        "Features",
		"fix: null pointer":            "Fixes",
		"bugfix: race in watcher":      "Fixes",
		"refactor: extract analyzer":    "Refactor",
		"docs: update README":           "Docs",
		"doc: clarify install":          "Docs",
		"test: cover commit types":      "Tests",
		"chore: bump deps":              "Others",
		"Add login flow":                "Features",
		"Fix flaky CI":                  "Fixes",
		"random WIP":                    "Others",
	}
	for msg, want := range cases {
		if got := commitTypeOf(msg); got != want {
			t.Fatalf("commitTypeOf(%q)=%s want %s", msg, got, want)
		}
	}
}

func TestClassifyCommitTypesPercentages(t *testing.T) {
	// Use synthetic messages through commitTypeOf aggregation helper path.
	msgs := []string{
		"feat: a", "feat: b", "feat: c", "feat: d", // 4
		"fix: a", "fix: b", "fix: c", // 3
		"refactor: a", // 1
		"docs: a",     // 1
		"other stuff", // 1
	}
	counts := map[string]int{}
	for _, m := range msgs {
		counts[commitTypeOf(m)]++
	}
	if counts["Features"] != 4 || counts["Fixes"] != 3 {
		t.Fatalf("unexpected counts: %#v", counts)
	}
}

