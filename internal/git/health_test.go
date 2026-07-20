package git

import (
	"testing"
	"time"
)

func TestCalculateHealthyRepo(t *testing.T) {
	stats := CommitStats{
		TotalCommits:     100,
		TotalAuthors:     5,
		AvgCommitsPerDay: 2.5,
		Additions:        5000,
		Deletions:        2000,
		LastCommit:       time.Now().Add(-2 * time.Hour),
	}
	contributors := []Contributor{
		{Name: "A", Commits: 40},
		{Name: "B", Commits: 30},
		{Name: "C", Commits: 20},
		{Name: "D", Commits: 5},
		{Name: "E", Commits: 5},
	}
	branches := []BranchInfo{
		{Name: "main", IsCurrent: true, LastCommit: time.Now()},
		{Name: "feature", LastCommit: time.Now().AddDate(0, -1, 0)},
	}
	streak := Streak{Current: 5, Longest: 12}

	h := Calculate(stats, contributors, branches, streak)
	if h.Score < 50 {
		t.Fatalf("expected healthy score >= 50, got %.1f", h.Score)
	}
	if h.Grade == "F" {
		t.Fatalf("expected non-F grade, got %s", h.Grade)
	}
	if h.ActivityRecency < 80 {
		t.Fatalf("expected high recency for recent commit, got %.1f", h.ActivityRecency)
	}
}

func TestCalculateEmptyRepo(t *testing.T) {
	h := Calculate(CommitStats{}, nil, nil, Streak{})
	if h.Score < 0 || h.Score > 100 {
		t.Fatalf("score out of range: %.1f", h.Score)
	}
	if h.Grade == "" {
		t.Fatal("expected grade")
	}
}

func TestGradeBoundaries(t *testing.T) {
	cases := []struct {
		score float64
		grade string
	}{
		{95, "A+"},
		{87, "A"},
		{72, "B"},
		{55, "C"},
		{42, "D"},
		{10, "F"},
	}
	for _, tc := range cases {
		if g := grade(tc.score); g != tc.grade {
			t.Fatalf("grade(%.0f)=%s want %s", tc.score, g, tc.grade)
		}
	}
}

func TestCalculateProductivity(t *testing.T) {
	stats := CommitStats{
		AvgCommitsPerDay: 2.0,
		LastCommit:       time.Now(),
	}
	branches := []BranchInfo{
		{Name: "main", IsCurrent: true},
		{Name: "feature/a"},
		{Name: "feature/b"},
	}
	types := []CommitTypeStat{
		{Name: "Docs", Percentage: 8},
		{Name: "Tests", Percentage: 6},
	}
	langs := []LanguageStat{
		{Name: "Go", Percentage: 70},
		{Name: "Markdown", Percentage: 10},
	}
	repo := Repository{License: "MIT"}
	h := HealthScore{
		Score:           85,
		CommitFrequency: 80,
	}

	p := CalculateProductivity(stats, branches, types, langs, repo, h)
	if p.Label != "Code Quality" {
		t.Fatalf("label = %s", p.Label)
	}
	if p.Score < 40 || p.Score > 100 {
		t.Fatalf("score out of range: %.1f", p.Score)
	}
	if len(p.Factors) != 5 {
		t.Fatalf("expected 5 factors, got %d", len(p.Factors))
	}
	names := map[string]bool{}
	for _, f := range p.Factors {
		names[f.Name] = true
	}
	for _, want := range []string{
		"Commit frequency", "Branch usage", "Documentation", "Tests", "Repository health",
	} {
		if !names[want] {
			t.Fatalf("missing factor %q", want)
		}
	}
}
