package dashboard

import (
	"fmt"
	"strings"

	"github.com/guptarohit/asciigraph"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/git"
)

// LineChart renders a commit timeline as an ASCII line chart.
func LineChart(points []git.TimelinePoint, height, width int) string {
	if len(points) == 0 {
		return "(no data)"
	}
	data := make([]float64, len(points))
	for i, p := range points {
		data[i] = float64(p.Commits)
	}
	if height <= 0 {
		height = 8
	}
	if width <= 0 {
		width = 60
	}
	return asciigraph.Plot(data,
		asciigraph.Height(height),
		asciigraph.Width(width),
		asciigraph.Caption("Commit activity"),
	)
}

// BarChart renders a horizontal bar chart for labeled values.
func BarChart(labels []string, values []float64, maxWidth int) string {
	if len(labels) == 0 || len(labels) != len(values) {
		return "(no data)"
	}
	if maxWidth <= 0 {
		maxWidth = 40
	}
	max := 0.0
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	if max == 0 {
		max = 1
	}

	var b strings.Builder
	labelWidth := 0
	for _, l := range labels {
		if len(l) > labelWidth {
			labelWidth = len(l)
		}
	}
	if labelWidth > 18 {
		labelWidth = 18
	}

	for i, label := range labels {
		if len(label) > labelWidth {
			label = label[:labelWidth-1] + "…"
		}
		barLen := int((values[i] / max) * float64(maxWidth))
		bar := strings.Repeat("█", barLen)
		fmt.Fprintf(&b, "%-*s │%s %.1f\n", labelWidth, label, bar, values[i])
	}
	return strings.TrimRight(b.String(), "\n")
}

// ContributorBars renders contributor commit bars.
func ContributorBars(contributors []git.Contributor, limit int) string {
	if len(contributors) == 0 {
		return "(no contributors)"
	}
	if limit > 0 && len(contributors) > limit {
		contributors = contributors[:limit]
	}
	labels := make([]string, len(contributors))
	values := make([]float64, len(contributors))
	for i, c := range contributors {
		labels[i] = c.Name
		values[i] = float64(c.Commits)
	}
	return BarChart(labels, values, 30)
}

// AuthorTimeline renders per-author monthly commit bars.
func AuthorTimeline(timelines []git.AuthorTimeline, maxWidth int) string {
	if len(timelines) == 0 {
		return "Author Timeline\n(no data)"
	}
	if maxWidth <= 0 {
		maxWidth = 9
	}

	var b strings.Builder
	b.WriteString("Author Timeline")
	for _, tl := range timelines {
		b.WriteString("\n")
		b.WriteString(tl.Name)
		maxCommits := 0
		for _, p := range tl.Monthly {
			if p.Commits > maxCommits {
				maxCommits = p.Commits
			}
		}
		for _, p := range tl.Monthly {
			barLen := 0
			if maxCommits > 0 && p.Commits > 0 {
				barLen = int(float64(p.Commits)/float64(maxCommits)*float64(maxWidth) + 0.5)
				if barLen < 1 {
					barLen = 1
				}
			}
			fmt.Fprintf(&b, "\n\n%-3s %s", p.Month, strings.Repeat("█", barLen))
		}
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n")
}

// HeatmapASCII renders a compact weekly contribution heatmap.
func HeatmapASCII(cells []git.HeatmapCell) string {
	if len(cells) == 0 {
		return "(no heatmap data)"
	}

	max := 0
	for _, c := range cells {
		if c.Count > max {
			max = c.Count
		}
	}

	levels := []rune{'░', '▒', '▓', '█'}
	weeks := 0
	for _, c := range cells {
		if c.Week > weeks {
			weeks = c.Week
		}
	}
	weeks++

	grid := make([][]int, 7)
	for i := range grid {
		grid[i] = make([]int, weeks)
	}
	for _, c := range cells {
		if c.Week >= 0 && c.Week < weeks && c.Weekday >= 0 && c.Weekday < 7 {
			grid[c.Weekday][c.Week] = c.Count
		}
	}

	days := []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
	var b strings.Builder
	b.WriteString("Weekly contribution heatmap\n")
	for d := 0; d < 7; d++ {
		b.WriteString(days[d])
		b.WriteString(" ")
		for w := 0; w < weeks; w++ {
			count := grid[d][w]
			if max == 0 || count == 0 {
				b.WriteRune('·')
				continue
			}
			lvl := (count * len(levels)) / (max + 1)
			if lvl >= len(levels) {
				lvl = len(levels) - 1
			}
			b.WriteRune(levels[lvl])
		}
		b.WriteByte('\n')
	}
	b.WriteString("   less · ░ ▒ ▓ █ more")
	return strings.TrimRight(b.String(), "\n")
}

// WeekdayBars renders commits by weekday.
func WeekdayBars(byWeekday [7]int) string {
	labels := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	values := make([]float64, 7)
	for i, v := range byWeekday {
		values[i] = float64(v)
	}
	return BarChart(labels, values, 30)
}

// WeeklyHeatmap renders weekday commit intensity as horizontal bars (Mon→Sun).
func WeeklyHeatmap(byWeekday [7]int, maxWidth int) string {
	if maxWidth <= 0 {
		maxWidth = 10
	}

	order := []struct {
		label string
		idx   int
	}{
		{"Mon", 1},
		{"Tue", 2},
		{"Wed", 3},
		{"Thu", 4},
		{"Fri", 5},
		{"Sat", 6},
		{"Sun", 0},
	}

	maxCommits := 0
	for _, o := range order {
		if byWeekday[o.idx] > maxCommits {
			maxCommits = byWeekday[o.idx]
		}
	}

	var b strings.Builder
	b.WriteString("Weekly Heatmap\n")
	for _, o := range order {
		count := byWeekday[o.idx]
		barLen := 0
		if maxCommits > 0 && count > 0 {
			barLen = int(float64(count)/float64(maxCommits)*float64(maxWidth) + 0.5)
			if barLen < 1 {
				barLen = 1
			}
		}
		fmt.Fprintf(&b, "%-3s %s\n", o.label, strings.Repeat("█", barLen))
	}
	return strings.TrimRight(b.String(), "\n")
}

// MonthlyActivityGraph renders the year-to-date commit activity graph.
func MonthlyActivityGraph(points []git.MonthlyPoint, maxWidth int) string {
	if len(points) == 0 {
		return "(no monthly activity)"
	}
	if maxWidth <= 0 {
		maxWidth = 17
	}

	maxCommits := 0
	for _, p := range points {
		if p.Commits > maxCommits {
			maxCommits = p.Commits
		}
	}

	var b strings.Builder
	b.WriteString("Commit Activity Graph\n")
	for _, p := range points {
		barLen := 0
		if maxCommits > 0 && p.Commits > 0 {
			barLen = int(float64(p.Commits)/float64(maxCommits)*float64(maxWidth) + 0.5)
			if barLen < 1 {
				barLen = 1
			}
		}
		bar := strings.Repeat("█", barLen)
		fmt.Fprintf(&b, "%-3s %s\n", p.Month, bar)
	}
	return strings.TrimRight(b.String(), "\n")
}
