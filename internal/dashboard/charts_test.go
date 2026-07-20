package dashboard

import (
	"strings"
	"testing"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/git"
)

func TestLineChartEmpty(t *testing.T) {
	if got := LineChart(nil, 5, 20); got != "(no data)" {
		t.Fatalf("unexpected: %s", got)
	}
}

func TestLineChartRenders(t *testing.T) {
	points := []git.TimelinePoint{
		{Date: "2026-01-01", Commits: 1},
		{Date: "2026-01-02", Commits: 3},
		{Date: "2026-01-03", Commits: 2},
	}
	out := LineChart(points, 5, 30)
	if !strings.Contains(out, "Commit activity") {
		t.Fatalf("missing caption: %s", out)
	}
}

func TestLanguageStatistics(t *testing.T) {
	langs := []git.LanguageStat{
		{Name: "Go", Percentage: 72},
		{Name: "HTML", Percentage: 14},
		{Name: "CSS", Percentage: 8},
		{Name: "YAML", Percentage: 4},
		{Name: "JSON", Percentage: 1},
		{Name: "Shell", Percentage: 1},
	}
	out := LanguageStatistics(langs, 4)
	for _, want := range []string{
		"Language Statistics",
		"Go", "72%",
		"HTML", "14%",
		"CSS", "8%",
		"YAML", "4%",
		"Others", "2%",
		"█",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}

func TestLanguageBars(t *testing.T) {
	langs := []git.LanguageStat{
		{Name: "Go", Percentage: 70},
		{Name: "Markdown", Percentage: 30},
	}
	out := LanguageBars(langs, 5)
	if !strings.Contains(out, "Go") {
		t.Fatalf("expected Go in bars: %s", out)
	}
}

func TestAuthorTimeline(t *testing.T) {
	timelines := []git.AuthorTimeline{
		{
			Name: "John",
			Monthly: []git.MonthlyPoint{
				{Month: "Jan", Commits: 3},
				{Month: "Feb", Commits: 6},
				{Month: "Mar", Commits: 9},
				{Month: "Apr", Commits: 7},
			},
		},
	}
	out := AuthorTimeline(timelines, 9)
	for _, want := range []string{
		"Author Timeline",
		"John",
		"Jan",
		"Feb",
		"Mar",
		"Apr",
		"█",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}

func TestGitStreak(t *testing.T) {
	out := GitStreak(git.Streak{Current: 16, Longest: 21})
	for _, want := range []string{
		"Git Streak",
		"Current Streak",
		"████████████",
		"16 Days",
		"Longest: 21 Days",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}

func TestProductivityScore(t *testing.T) {
	out := ProductivityScore(git.ProductivityScore{
		Label: "Code Quality",
		Score: 87,
		Factors: []git.ProductivityFactor{
			{Name: "Commit frequency", Passed: true},
			{Name: "Branch usage", Passed: true},
			{Name: "Documentation", Passed: true},
			{Name: "Tests", Passed: true},
			{Name: "Repository health", Passed: true},
		},
	})
	for _, want := range []string{
		"Productivity Score",
		"Code Quality",
		"87%",
		"Based on",
		"✓ Commit frequency",
		"✓ Branch usage",
		"✓ Documentation",
		"✓ Tests",
		"✓ Repository health",
		"█",
		"░",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}

func TestFileChangeStatistics(t *testing.T) {
	out := FileChangeStatistics(git.FileChangeSummary{
		Added: 312, Modified: 121, Deleted: 38, Renamed: 14,
	})
	for _, want := range []string{
		"File Change Statistics",
		"Added", "312",
		"Modified", "121",
		"Deleted", "38",
		"Renamed", "14",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}

func TestRecentCommitsList(t *testing.T) {
	commits := []git.RecentCommit{
		{Message: "Fix login timeout"},
		{Message: "feat: add payment API"},
		{Message: "Improve dashboard"},
		{Message: "Refactor services"},
		{Message: "Update README"},
	}
	out := RecentCommitsList(commits, 10)
	for _, want := range []string{
		"Recent Commits",
		"✓ Fix login timeout",
		"✓ Add payment API",
		"✓ Improve dashboard",
		"✓ Refactor services",
		"✓ Update README",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}

func TestBranchViewer(t *testing.T) {
	branches := []git.BranchInfo{
		{Name: "feature/login"},
		{Name: "main", IsCurrent: true},
		{Name: "feature/payment"},
		{Name: "release/v2"},
		{Name: "hotfix/session"},
		{Name: "origin/main", IsRemote: true},
	}
	out := BranchViewer(branches, 0)
	for _, want := range []string{
		"Branch Viewer",
		"main",
		"feature/login",
		"feature/payment",
		"release/v2",
		"hotfix/session",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "origin/main") {
		t.Fatalf("should prefer local branches:\n%s", out)
	}
	lines := strings.Split(out, "\n")
	if len(lines) < 2 || lines[1] != "main" {
		t.Fatalf("current branch should be first:\n%s", out)
	}
}

func TestTopContributors(t *testing.T) {
	contributors := []git.Contributor{
		{Name: "John Doe", Commits: 421},
		{Name: "Alice", Commits: 301},
		{Name: "David", Commits: 192},
		{Name: "Sarah", Commits: 166},
	}
	out := TopContributors(contributors, 10)
	for _, want := range []string{
		"Top Contributors",
		"Name",
		"Commits",
		"John Doe",
		"421",
		"Alice",
		"301",
		"David",
		"192",
		"Sarah",
		"166",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}

func TestHeatmapASCII(t *testing.T) {
	cells := []git.HeatmapCell{
		{Date: "2026-01-01", Weekday: 4, Week: 0, Count: 2},
		{Date: "2026-01-02", Weekday: 5, Week: 0, Count: 0},
	}
	out := HeatmapASCII(cells)
	if !strings.Contains(out, "heatmap") {
		t.Fatalf("unexpected heatmap output: %s", out)
	}
}

func TestWeeklyHeatmap(t *testing.T) {
	var byWeekday [7]int
	byWeekday[0] = 1  // Sun
	byWeekday[1] = 2  // Mon
	byWeekday[2] = 5  // Tue
	byWeekday[3] = 8  // Wed
	byWeekday[4] = 10 // Thu
	byWeekday[5] = 6  // Fri
	byWeekday[6] = 2  // Sat

	out := WeeklyHeatmap(byWeekday, 10)
	if !strings.Contains(out, "Weekly Heatmap") {
		t.Fatalf("missing title:\n%s", out)
	}
	lines := strings.Split(out, "\n")
	if len(lines) < 8 {
		t.Fatalf("expected title + 7 days, got:\n%s", out)
	}
	if !strings.HasPrefix(lines[1], "Mon") {
		t.Fatalf("expected Mon first, got %q", lines[1])
	}
	if !strings.HasPrefix(lines[7], "Sun") {
		t.Fatalf("expected Sun last, got %q", lines[7])
	}
	var thu, sun string
	for _, line := range lines {
		if strings.HasPrefix(line, "Thu") {
			thu = line
		}
		if strings.HasPrefix(line, "Sun") {
			sun = line
		}
	}
	if strings.Count(thu, "█") <= strings.Count(sun, "█") {
		t.Fatalf("expected Thu longer than Sun:\n%s\n%s", thu, sun)
	}
}

func TestCommitTypes(t *testing.T) {
	types := []git.CommitTypeStat{
		{Name: "Features", Percentage: 41},
		{Name: "Fixes", Percentage: 29},
		{Name: "Refactor", Percentage: 14},
		{Name: "Docs", Percentage: 7},
		{Name: "Tests", Percentage: 4},
		{Name: "Others", Percentage: 5},
	}
	out := CommitTypes(types, 12)
	for _, want := range []string{
		"Commit Types",
		"Features",
		"41%",
		"Fixes",
		"29%",
		"Refactor",
		"14%",
		"Docs",
		"7%",
		"Tests",
		"4%",
		"Others",
		"5%",
		"█",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}

func TestMonthlyActivityGraph(t *testing.T) {
	points := []git.MonthlyPoint{
		{Month: "Jan", Commits: 3},
		{Month: "Feb", Commits: 8},
		{Month: "Mar", Commits: 11},
		{Month: "Apr", Commits: 14},
		{Month: "May", Commits: 7},
		{Month: "Jun", Commits: 17},
		{Month: "Jul", Commits: 10},
	}
	out := MonthlyActivityGraph(points, 17)
	if !strings.Contains(out, "Commit Activity Graph") {
		t.Fatalf("missing title: %s", out)
	}
	for _, month := range []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul"} {
		if !strings.Contains(out, month) {
			t.Fatalf("missing month %s:\n%s", month, out)
		}
	}
	if !strings.Contains(out, "█") {
		t.Fatalf("expected bars:\n%s", out)
	}
	lines := strings.Split(out, "\n")
	var jun, jan string
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "Jun") {
			jun = line
		}
		if strings.HasPrefix(strings.TrimSpace(line), "Jan") {
			jan = line
		}
	}
	if strings.Count(jun, "█") <= strings.Count(jan, "█") {
		t.Fatalf("expected Jun bar longer than Jan:\n%s\n%s", jan, jun)
	}
}

