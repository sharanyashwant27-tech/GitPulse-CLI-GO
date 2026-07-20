package dashboard

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/git"
	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/theme"
)

// Render returns the repository overview block matching the GitPulse feature layout.
func Render(a *git.Analysis, s theme.Styles) string {
	if a == nil {
		return ""
	}
	repo := a.Repository
	status := repo.Status
	if status == "" {
		status = "Unknown"
	}

	statusStyle := s.Success
	if !strings.EqualFold(status, "Clean") {
		statusStyle = s.Warning
	}

	row := func(label, value string) string {
		return s.Muted.Render(fmt.Sprintf("%-12s: ", label)) + value
	}

	lines := []string{
		s.Title.Render("Repository Overview") + s.Muted.Render(" · ") +
			s.Muted.Render("Repository: ") + s.Value.Render(repo.Name),
		"",
		row("Branch", s.Highlight.Render(repo.DefaultBranch)),
		row("Status", statusStyle.Render(status)),
		row("Commits", s.Value.Render(formatInt(a.CommitStats.TotalCommits))),
		row("Branches", s.Value.Render(formatInt(repo.BranchCount))),
		row("Tags", s.Value.Render(formatInt(repo.TagCount))),
		row("Contributors", s.Value.Render(formatInt(a.CommitStats.TotalAuthors))),
		row("Created", s.Value.Render(formatCreated(repo.CreatedAt))),
		row("Last Commit", s.Value.Render(formatRelative(repo.LastCommitAt))),
		"",
		s.ProgressBar(a.Health.Score, 26),
		row("Health Score", healthValue(a.Health, s)),
	}
	return strings.Join(lines, "\n")
}

// RenderBeautiful renders the centered GIT PULSE dashboard card.
func RenderBeautiful(a *git.Analysis, s theme.Styles, width int) string {
	if a == nil {
		return ""
	}
	if width < 44 {
		width = 44
	}
	inner := width - 2
	if inner < 40 {
		inner = 40
	}

	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Theme.Primary).
		Width(inner).
		Padding(0, 1)

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(s.Theme.Primary).
		Align(lipgloss.Center).
		Width(inner - 2).
		Render("GIT PULSE")

	divider := s.Muted.Render(strings.Repeat("─", inner-2))

	lang := a.Repository.PrimaryLanguage
	if lang == "" {
		lang = "N/A"
	}
	license := a.Repository.License
	if license == "" {
		license = "N/A"
	}
	stars := a.Repository.Stars
	if stars == "" {
		stars = "N/A"
	}

	tree := func(label string, items []string) string {
		var b strings.Builder
		b.WriteString(s.Subtitle.Render(label))
		b.WriteByte('\n')
		for i, item := range items {
			prefix := " ├── "
			if i == len(items)-1 {
				prefix = " └── "
			}
			b.WriteString(s.Muted.Render(prefix))
			b.WriteString(item)
			b.WriteByte('\n')
		}
		return strings.TrimRight(b.String(), "\n")
	}

	repoSection := tree(" Repository", []string{
		s.Value.Render(a.Repository.Name),
		"Size: " + s.Value.Render(git.FormatBytes(a.Repository.SizeBytes)),
		"Language: " + s.Highlight.Render(lang),
		"Stars: " + s.Muted.Render(stars),
		"License: " + s.Value.Render(license),
	})

	activitySection := tree(" Activity", []string{
		fmt.Sprintf("Today's commits : %s", s.Value.Render(fmt.Sprintf("%d", a.CommitStats.Today))),
		fmt.Sprintf("This week       : %s", s.Value.Render(fmt.Sprintf("%d", a.CommitStats.ThisWeek))),
		fmt.Sprintf("This month      : %s", s.Value.Render(fmt.Sprintf("%d", a.CommitStats.ThisMonth))),
	})

	healthBar := s.ProgressBar(a.Health.Score, 20)
	healthPct := healthValue(a.Health, s)
	healthSection := strings.Join([]string{
		s.Subtitle.Render(" Health"),
		" " + healthBar + " " + healthPct,
	}, "\n")

	body := strings.Join([]string{
		title,
		divider,
		"",
		repoSection,
		"",
		activitySection,
		"",
		healthSection,
		"",
		divider,
	}, "\n")

	return border.Render(body)
}

// LanguageBars renders language breakdown bars (legacy alias).
func LanguageBars(langs []git.LanguageStat, limit int) string {
	return LanguageStatistics(langs, limit)
}

// LanguageStatistics renders language composition with percentage bars.
func LanguageStatistics(langs []git.LanguageStat, limit int) string {
	if len(langs) == 0 {
		return "Language Statistics\n(no language data)"
	}
	if limit <= 0 {
		limit = 4
	}

	display := make([]git.LanguageStat, 0, limit+1)
	othersPct := 0.0
	othersCount := 0
	for i, l := range langs {
		if i < limit {
			display = append(display, l)
			continue
		}
		othersPct += l.Percentage
		othersCount += l.Files
	}
	if othersPct > 0 || othersCount > 0 {
		display = append(display, git.LanguageStat{
			Name:       "Others",
			Percentage: othersPct,
			Files:      othersCount,
		})
	}

	const maxWidth = 11
	var b strings.Builder
	b.WriteString("Language Statistics\n")
	for i, l := range display {
		barLen := int(l.Percentage/100*float64(maxWidth) + 0.5)
		if l.Percentage > 0 && barLen < 1 {
			barLen = 1
		}
		fmt.Fprintf(&b, "%-7s %s %.0f%%", l.Name, strings.Repeat("█", barLen), l.Percentage)
		if i < len(display)-1 {
			b.WriteString("\n\n")
		}
	}
	return b.String()
}

// TopContributors renders a clean name/commits leaderboard.
func TopContributors(contributors []git.Contributor, limit int) string {
	if len(contributors) == 0 {
		return "Top Contributors\n(no contributors)"
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > len(contributors) {
		limit = len(contributors)
	}

	const nameWidth = 20
	var b strings.Builder
	b.WriteString("Top Contributors\n")
	fmt.Fprintf(&b, "%-*s %s\n\n", nameWidth, "Name", "Commits")
	for i := 0; i < limit; i++ {
		c := contributors[i]
		name := c.Name
		if name == "" {
			name = c.Email
		}
		if len(name) > nameWidth {
			name = name[:nameWidth-1] + "…"
		}
		fmt.Fprintf(&b, "%-*s %d\n", nameWidth, name, c.Commits)
	}
	return strings.TrimRight(b.String(), "\n")
}

// BranchViewer renders a clean spaced list of local branch names.
func BranchViewer(branches []git.BranchInfo, limit int) string {
	names := branchViewerNames(branches)
	if len(names) == 0 {
		return "Branch Viewer\n(no branches)"
	}
	if limit > 0 && len(names) > limit {
		names = names[:limit]
	}

	var b strings.Builder
	b.WriteString("Branch Viewer\n")
	for i, name := range names {
		b.WriteString(name)
		if i < len(names)-1 {
			b.WriteString("\n\n")
		}
	}
	return b.String()
}

// RecentCommitsList renders recent commit subjects with checkmarks.
func RecentCommitsList(commits []git.RecentCommit, limit int) string {
	if len(commits) == 0 {
		return "Recent Commits\n(no commits)"
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > len(commits) {
		limit = len(commits)
	}

	var b strings.Builder
	b.WriteString("Recent Commits\n")
	for i := 0; i < limit; i++ {
		msg := strings.TrimSpace(commits[i].Message)
		if msg == "" {
			msg = commits[i].ShortHash
		}
		msg = humanizeCommitSubject(msg)
		fmt.Fprintf(&b, "✓ %s", msg)
		if i < limit-1 {
			b.WriteString("\n\n")
		}
	}
	return b.String()
}

// FileChangeStatistics renders Added/Modified/Deleted/Renamed totals.
func FileChangeStatistics(summary git.FileChangeSummary) string {
	rows := []struct {
		label string
		value int
	}{
		{"Added", summary.Added},
		{"Modified", summary.Modified},
		{"Deleted", summary.Deleted},
		{"Renamed", summary.Renamed},
	}

	var b strings.Builder
	b.WriteString("File Change Statistics\n")
	for i, row := range rows {
		fmt.Fprintf(&b, "%-10s : %d", row.label, row.value)
		if i < len(rows)-1 {
			b.WriteString("\n\n")
		}
	}
	return b.String()
}

// ProductivityScore renders the Code Quality productivity card.
func ProductivityScore(p git.ProductivityScore) string {
	if p.Label == "" {
		p.Label = "Code Quality"
	}
	barWidth := 20
	filled := int(p.Score/100*float64(barWidth) + 0.5)
	if filled > barWidth {
		filled = barWidth
	}
	if p.Score > 0 && filled < 1 {
		filled = 1
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	var b strings.Builder
	b.WriteString("Productivity Score\n")
	b.WriteString(p.Label)
	b.WriteString("\n\n")
	b.WriteString(bar)
	b.WriteString("\n\n")
	fmt.Fprintf(&b, "%.0f%%\n\n", p.Score)
	b.WriteString("Based on\n")
	for _, f := range p.Factors {
		mark := "✓"
		if !f.Passed {
			mark = "○"
		}
		fmt.Fprintf(&b, "\n%s %s", mark, f.Name)
	}
	return b.String()
}

// GitStreak renders the current commit streak card.
func GitStreak(streak git.Streak) string {
	const barWidth = 12
	filled := streak.Current
	if filled > barWidth {
		filled = barWidth
	}
	if streak.Current > 0 && filled < 1 {
		filled = 1
	}
	bar := strings.Repeat("█", filled)
	if filled == 0 {
		bar = strings.Repeat("░", barWidth)
	}

	daysLabel := "Days"
	if streak.Current == 1 {
		daysLabel = "Day"
	}

	var b strings.Builder
	b.WriteString("Git Streak\n")
	b.WriteString("Current Streak\n\n")
	b.WriteString(bar)
	b.WriteString("\n\n")
	fmt.Fprintf(&b, "%d %s", streak.Current, daysLabel)
	if streak.Longest > 0 && streak.Longest != streak.Current {
		fmt.Fprintf(&b, "\n\nLongest: %d Days", streak.Longest)
	}
	return b.String()
}

// CommitTypes renders categorized commit-message shares.
func CommitTypes(types []git.CommitTypeStat, maxWidth int) string {
	if len(types) == 0 {
		return "Commit Types\n(no data)"
	}
	if maxWidth <= 0 {
		maxWidth = 12
	}

	var b strings.Builder
	b.WriteString("Commit Types\n")
	for i, t := range types {
		barLen := int(t.Percentage/100*float64(maxWidth) + 0.5)
		if t.Percentage > 0 && barLen < 1 {
			barLen = 1
		}
		fmt.Fprintf(&b, "%-12s %s %.0f%%", t.Name, strings.Repeat("█", barLen), t.Percentage)
		if i < len(types)-1 {
			b.WriteString("\n\n")
		} else {
			b.WriteByte('\n')
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func healthValue(h git.HealthScore, s theme.Styles) string {
	pct := fmt.Sprintf("%.0f%%", h.Score)
	style := s.Success
	if h.Score < 70 {
		style = s.Warning
	}
	if h.Score < 50 {
		style = s.Error
	}
	return style.Render(pct)
}

func formatInt(n int) string {
	s := strconv.Itoa(n)
	if n < 1000 {
		return s
	}
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	return b.String()
}

func formatCreated(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("Jan 2006")
}

func formatRelative(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		n := int(d.Minutes())
		if n == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", n)
	case d < 48*time.Hour:
		n := int(d.Hours())
		if n == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", n)
	case d < 30*24*time.Hour:
		n := int(d.Hours() / 24)
		if n == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", n)
	case d < 365*24*time.Hour:
		n := int(d.Hours() / 24 / 30)
		if n <= 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", n)
	default:
		n := int(d.Hours() / 24 / 365)
		if n <= 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", n)
	}
}

func branchViewerNames(branches []git.BranchInfo) []string {
	var local []git.BranchInfo
	for _, b := range branches {
		if !b.IsRemote {
			local = append(local, b)
		}
	}
	source := local
	if len(source) == 0 {
		source = branches
	}

	seen := map[string]bool{}
	names := make([]string, 0, len(source))
	for _, b := range source {
		if b.IsCurrent && !seen[b.Name] {
			names = append(names, b.Name)
			seen[b.Name] = true
		}
	}
	for _, b := range source {
		if seen[b.Name] {
			continue
		}
		names = append(names, b.Name)
		seen[b.Name] = true
	}
	return names
}

func humanizeCommitSubject(msg string) string {
	msg = strings.TrimSpace(msg)
	if i := strings.Index(msg, ":"); i > 0 && i < 40 {
		prefix := strings.ToLower(strings.TrimSpace(msg[:i]))
		if !strings.Contains(prefix, " ") && len(prefix) <= 20 {
			rest := strings.TrimSpace(msg[i+1:])
			if rest != "" {
				return capitalizeFirst(rest)
			}
		}
	}
	return msg
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	if r[0] >= 'a' && r[0] <= 'z' {
		r[0] = r[0] - ('a' - 'A')
	}
	return string(r)
}
