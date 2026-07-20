package git

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func (a *Analyzer) loadCommits() ([]*object.Commit, error) {
	ref, err := a.repo.Head()
	if err != nil {
		return nil, fmt.Errorf("resolve HEAD: %w", err)
	}

	iter, err := a.repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("walk commits: %w", err)
	}
	defer iter.Close()

	var commits []*object.Commit
	err = iter.ForEach(func(c *object.Commit) error {
		when := c.Author.When
		if !a.opts.Since.IsZero() && when.Before(a.opts.Since) {
			return nil
		}
		if !a.opts.Until.IsZero() && when.After(a.opts.Until) {
			return nil
		}
		commits = append(commits, c)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return commits, nil
}

func (a *Analyzer) commitStats(commits []*object.Commit) CommitStats {
	stats := CommitStats{
		TotalCommits: len(commits),
		CommitsByDay: make(map[string]int),
	}
	authors := map[string]struct{}{}
	var additions, deletions, filesChanged int

	for i, c := range commits {
		authors[strings.ToLower(c.Author.Email)] = struct{}{}
		day := c.Author.When.Format("2006-01-02")
		stats.CommitsByDay[day]++
		stats.CommitsByWeekday[int(c.Author.When.Weekday())]++
		stats.CommitsByHour[c.Author.When.Hour()]++

		if i == 0 {
			stats.LastCommit = c.Author.When
		}
		if i == len(commits)-1 {
			stats.FirstCommit = c.Author.When
		}

		if i < 200 {
			add, del, files := commitDiffStats(c)
			additions += add
			deletions += del
			filesChanged += files
		}
	}

	stats.TotalAuthors = len(authors)
	stats.Additions = additions
	stats.Deletions = deletions
	stats.FilesChanged = filesChanged
	stats.Today, stats.ThisWeek, stats.ThisMonth = activityWindows(stats.CommitsByDay)

	if !stats.FirstCommit.IsZero() && !stats.LastCommit.IsZero() {
		days := stats.LastCommit.Sub(stats.FirstCommit).Hours() / 24
		if days < 1 {
			days = 1
		}
		stats.AvgCommitsPerDay = float64(stats.TotalCommits) / days
	}
	return stats
}

func activityWindows(byDay map[string]int) (today, week, month int) {
	now := time.Now()
	todayKey := now.Format("2006-01-02")
	today = byDay[todayKey]

	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).
		AddDate(0, 0, -(weekday - 1))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	for day, count := range byDay {
		t, err := time.ParseInLocation("2006-01-02", day, now.Location())
		if err != nil {
			continue
		}
		if !t.Before(weekStart) && !t.After(now) {
			week += count
		}
		if !t.Before(monthStart) && !t.After(now) {
			month += count
		}
	}
	return today, week, month
}

func commitDiffStats(c *object.Commit) (additions, deletions, files int) {
	stats, err := c.Stats()
	if err != nil {
		return 0, 0, 0
	}
	for _, s := range stats {
		additions += s.Addition
		deletions += s.Deletion
		files++
	}
	return additions, deletions, files
}

func (a *Analyzer) timeline(commits []*object.Commit, days int) []TimelinePoint {
	end := time.Now().Truncate(24 * time.Hour)
	start := end.AddDate(0, 0, -days+1)
	byDay := map[string]int{}
	for _, c := range commits {
		d := c.Author.When.Truncate(24 * time.Hour)
		if d.Before(start) || d.After(end) {
			continue
		}
		byDay[d.Format("2006-01-02")]++
	}

	points := make([]TimelinePoint, 0, days)
	for i := 0; i < days; i++ {
		d := start.AddDate(0, 0, i)
		key := d.Format("2006-01-02")
		points = append(points, TimelinePoint{Date: key, Commits: byDay[key]})
	}
	return points
}

func (a *Analyzer) monthlyActivity(commits []*object.Commit) []MonthlyPoint {
	now := time.Now()
	year := now.Year()
	byMonth := map[string]int{}
	for _, c := range commits {
		t := c.Author.When.In(now.Location())
		if t.Year() != year {
			continue
		}
		key := t.Format("2006-01")
		byMonth[key]++
	}

	points := make([]MonthlyPoint, 0, int(now.Month()))
	for m := time.January; m <= now.Month(); m++ {
		t := time.Date(year, m, 1, 0, 0, 0, 0, now.Location())
		key := t.Format("2006-01")
		points = append(points, MonthlyPoint{
			Month:     t.Format("Jan"),
			YearMonth: key,
			Commits:   byMonth[key],
		})
	}
	return points
}

func (a *Analyzer) authorTimelines(commits []*object.Commit, limit int) []AuthorTimeline {
	if limit <= 0 {
		limit = 5
	}
	now := time.Now()
	year := now.Year()

	type authorKey struct {
		name  string
		email string
	}
	totals := map[string]int{}
	names := map[string]authorKey{}
	byAuthorMonth := map[string]map[string]int{}

	for _, c := range commits {
		t := c.Author.When.In(now.Location())
		if t.Year() != year {
			continue
		}
		email := strings.ToLower(c.Author.Email)
		if email == "" {
			email = strings.ToLower(c.Author.Name)
		}
		names[email] = authorKey{name: c.Author.Name, email: c.Author.Email}
		totals[email]++
		if byAuthorMonth[email] == nil {
			byAuthorMonth[email] = map[string]int{}
		}
		byAuthorMonth[email][t.Format("2006-01")]++
	}

	emails := make([]string, 0, len(totals))
	for email := range totals {
		emails = append(emails, email)
	}
	sort.Slice(emails, func(i, j int) bool {
		return totals[emails[i]] > totals[emails[j]]
	})
	if len(emails) > limit {
		emails = emails[:limit]
	}

	out := make([]AuthorTimeline, 0, len(emails))
	for _, email := range emails {
		info := names[email]
		display := info.name
		if display == "" {
			display = info.email
		}
		if parts := strings.Fields(display); len(parts) > 0 {
			display = parts[0]
		}
		monthly := make([]MonthlyPoint, 0, int(now.Month()))
		for m := time.January; m <= now.Month(); m++ {
			t := time.Date(year, m, 1, 0, 0, 0, 0, now.Location())
			key := t.Format("2006-01")
			monthly = append(monthly, MonthlyPoint{
				Month:     t.Format("Jan"),
				YearMonth: key,
				Commits:   byAuthorMonth[email][key],
			})
		}
		out = append(out, AuthorTimeline{
			Name:    display,
			Email:   info.email,
			Monthly: monthly,
		})
	}
	return out
}

var (
	conventionalTypeRe = regexp.MustCompile(`(?i)^(feat|feature|fix|bugfix|refactor|docs?|doc|test|tests|chore|ci|perf|style|build)(\(.+\))?!?:\s*`)
	keywordTypeRe      = regexp.MustCompile(`(?i)\b(feat(ure)?s?|fix(es|ed)?|bug\s*fix(es)?|refactor(ed|ing)?|docs?|documentation|tests?|chore|ci|perf|performance)\b`)
)

func classifyCommitTypes(commits []*object.Commit) []CommitTypeStat {
	order := []string{"Features", "Fixes", "Refactor", "Docs", "Tests", "Others"}
	counts := map[string]int{}
	for _, name := range order {
		counts[name] = 0
	}

	for _, c := range commits {
		msg := strings.SplitN(c.Message, "\n", 2)[0]
		counts[commitTypeOf(msg)]++
	}

	total := len(commits)
	out := make([]CommitTypeStat, 0, len(order))
	for _, name := range order {
		pct := 0.0
		if total > 0 {
			pct = (float64(counts[name]) / float64(total)) * 100
		}
		out = append(out, CommitTypeStat{
			Name:       name,
			Count:      counts[name],
			Percentage: pct,
		})
	}
	return out
}

func commitTypeOf(message string) string {
	msg := strings.TrimSpace(message)
	if msg == "" {
		return "Others"
	}

	if m := conventionalTypeRe.FindStringSubmatch(msg); len(m) > 1 {
		return mapCommitType(m[1])
	}

	lower := strings.ToLower(msg)
	switch {
	case strings.HasPrefix(lower, "add ") || strings.HasPrefix(lower, "implement ") || strings.HasPrefix(lower, "create "):
		return "Features"
	case strings.HasPrefix(lower, "fix ") || strings.HasPrefix(lower, "hotfix ") || strings.HasPrefix(lower, "patch "):
		return "Fixes"
	case strings.HasPrefix(lower, "refactor ") || strings.HasPrefix(lower, "clean ") || strings.HasPrefix(lower, "cleanup "):
		return "Refactor"
	case strings.HasPrefix(lower, "doc ") || strings.HasPrefix(lower, "docs ") || strings.HasPrefix(lower, "readme"):
		return "Docs"
	case strings.HasPrefix(lower, "test ") || strings.HasPrefix(lower, "tests "):
		return "Tests"
	}

	if m := keywordTypeRe.FindStringSubmatch(msg); len(m) > 1 {
		return mapCommitType(m[1])
	}
	return "Others"
}

func mapCommitType(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "feat", "feature", "features":
		return "Features"
	case "fix", "fixes", "fixed", "bugfix", "bugfixes", "bug fix", "bug fixes":
		return "Fixes"
	case "refactor", "refactored", "refactoring":
		return "Refactor"
	case "doc", "docs", "documentation":
		return "Docs"
	case "test", "tests":
		return "Tests"
	default:
		return "Others"
	}
}

func (a *Analyzer) heatmap(commits []*object.Commit, weeks int) []HeatmapCell {
	end := time.Now().Truncate(24 * time.Hour)
	start := end.AddDate(0, 0, -(weeks*7)+1)
	byDay := map[string]int{}
	for _, c := range commits {
		key := c.Author.When.Format("2006-01-02")
		byDay[key]++
	}

	cells := make([]HeatmapCell, 0, weeks*7)
	for i := 0; i < weeks*7; i++ {
		d := start.AddDate(0, 0, i)
		key := d.Format("2006-01-02")
		cells = append(cells, HeatmapCell{
			Date:    key,
			Weekday: int(d.Weekday()),
			Week:    i / 7,
			Count:   byDay[key],
		})
	}
	return cells
}

func (a *Analyzer) recentCommits(commits []*object.Commit, limit int) []RecentCommit {
	if limit > len(commits) {
		limit = len(commits)
	}
	out := make([]RecentCommit, 0, limit)
	for i := 0; i < limit; i++ {
		c := commits[i]
		msg := strings.SplitN(c.Message, "\n", 2)[0]
		add, del, files := 0, 0, 0
		if i < 50 {
			add, del, files = commitDiffStats(c)
		}
		out = append(out, RecentCommit{
			Hash:      c.Hash.String(),
			ShortHash: c.Hash.String()[:7],
			Message:   msg,
			Author:    c.Author.Name,
			Email:     c.Author.Email,
			Date:      c.Author.When,
			Files:     files,
			Additions: add,
			Deletions: del,
		})
	}
	return out
}

func (a *Analyzer) streak(commits []*object.Commit) Streak {
	days := map[string]bool{}
	var lastActive time.Time
	for _, c := range commits {
		key := c.Author.When.Format("2006-01-02")
		days[key] = true
		if c.Author.When.After(lastActive) {
			lastActive = c.Author.When
		}
	}
	if len(days) == 0 {
		return Streak{}
	}

	today := time.Now().Truncate(24 * time.Hour)
	current := 0
	for d := today; ; d = d.AddDate(0, 0, -1) {
		if days[d.Format("2006-01-02")] {
			current++
			continue
		}
		if d.Equal(today) {
			continue
		}
		break
	}

	longest := 0
	run := 0
	start := today.AddDate(-2, 0, 0)
	for d := start; !d.After(today); d = d.AddDate(0, 0, 1) {
		if days[d.Format("2006-01-02")] {
			run++
			if run > longest {
				longest = run
			}
		} else {
			run = 0
		}
	}

	return Streak{
		Current:    current,
		Longest:    longest,
		LastActive: lastActive,
	}
}
