package dashboard

import (
	"strings"
	"testing"
	"time"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/git"
	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/theme"
)

func TestRenderBeautiful(t *testing.T) {
	a := &git.Analysis{
		Repository: git.Repository{
			Name:            "my-project",
			SizeBytes:       240 * 1024 * 1024,
			PrimaryLanguage: "Go",
			Stars:           "N/A",
			License:         "MIT",
		},
		CommitStats: git.CommitStats{
			Today:     8,
			ThisWeek:  39,
			ThisMonth: 121,
		},
		Health: git.HealthScore{Score: 82, Grade: "A-"},
	}

	out := RenderBeautiful(a, theme.NewStyles(theme.Get("dracula")), 46)
	for _, want := range []string{
		"GIT PULSE",
		"Repository",
		"my-project",
		"Size:",
		"Language:",
		"Go",
		"Stars:",
		"N/A",
		"License:",
		"MIT",
		"Activity",
		"Today's commits",
		"8",
		"This week",
		"39",
		"This month",
		"121",
		"Health",
		"82%",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in beautiful dashboard:\n%s", want, out)
		}
	}
}

func TestRenderOverviewLayout(t *testing.T) {
	a := &git.Analysis{
		Repository: git.Repository{
			Name:          "my-project",
			DefaultBranch: "main",
			Status:        "Clean",
			BranchCount:   18,
			TagCount:      42,
			CreatedAt:     time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC),
			LastCommitAt:  time.Now().Add(-2 * time.Hour),
		},
		CommitStats: git.CommitStats{
			TotalCommits: 1532,
			TotalAuthors: 11,
		},
		Health: git.HealthScore{Score: 94, Grade: "A+"},
	}

	out := Render(a, theme.NewStyles(theme.Get("dracula")))
	for _, want := range []string{
		"Repository Overview",
		"my-project",
		"Branch",
		"main",
		"Status",
		"Clean",
		"1,532",
		"18",
		"42",
		"11",
		"Jan 2022",
		"2 hours ago",
		"Health Score",
		"94%",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in overview:\n%s", want, out)
		}
	}
}

func TestFormatInt(t *testing.T) {
	if formatInt(1532) != "1,532" {
		t.Fatalf("got %s", formatInt(1532))
	}
	if formatInt(42) != "42" {
		t.Fatalf("got %s", formatInt(42))
	}
}

