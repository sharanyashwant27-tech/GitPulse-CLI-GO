package dashboard

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/git"
	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/theme"
)

type tickMsg time.Time
type analysisMsg struct {
	analysis *git.Analysis
	err      error
}

// Options configures the interactive dashboard.
type Options struct {
	RepoPath      string
	Theme         string
	Watch         bool
	WatchInterval time.Duration
	Limit         int
}

type model struct {
	opts      Options
	styles    theme.Styles
	analysis  *git.Analysis
	err       error
	loading   bool
	width     int
	height    int
	tab       int
	tabs      []string
	viewport  viewport.Model
	ready     bool
	lastRefresh time.Time
}

// NewDashboard creates a Bubble Tea program for the GitPulse dashboard.
func NewDashboard(opts Options) *tea.Program {
	if opts.WatchInterval <= 0 {
		opts.WatchInterval = 5 * time.Second
	}
	t := theme.Get(opts.Theme)
	m := model{
		opts:   opts,
		styles: theme.NewStyles(t),
		tabs: []string{
			"Dashboard", "Overview", "Commits", "Contributors", "Branches",
			"Languages", "Timeline", "Heatmap", "Health",
		},
		loading: true,
	}
	return tea.NewProgram(m, tea.WithAltScreen())
}

func (m model) Init() tea.Cmd {
	cmds := []tea.Cmd{m.loadAnalysis()}
	if m.opts.Watch {
		cmds = append(cmds, tick(m.opts.WatchInterval))
	}
	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "right", "l":
			m.tab = (m.tab + 1) % len(m.tabs)
			m.viewport.SetContent(m.renderTab())
		case "shift+tab", "left", "h":
			m.tab = (m.tab - 1 + len(m.tabs)) % len(m.tabs)
			m.viewport.SetContent(m.renderTab())
		case "r":
			m.loading = true
			return m, m.loadAnalysis()
		case "j", "down":
			m.viewport.LineDown(1)
		case "k", "up":
			m.viewport.LineUp(1)
		case "pgdown", "f":
			m.viewport.ViewDown()
		case "pgup", "b":
			m.viewport.ViewUp()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		headerHeight := 6
		footerHeight := 2
		verticalMargin := headerHeight + footerHeight
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMargin)
			m.viewport.YPosition = headerHeight
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMargin
		}
		m.viewport.SetContent(m.renderTab())
	case analysisMsg:
		m.loading = false
		m.err = msg.err
		m.analysis = msg.analysis
		m.lastRefresh = time.Now()
		if m.ready {
			m.viewport.SetContent(m.renderTab())
		}
	case tickMsg:
		if m.opts.Watch {
			return m, tea.Batch(m.loadAnalysis(), tick(m.opts.WatchInterval))
		}
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading GitPulse…"
	}

	title := m.styles.Title.Render("⚡ GitPulse")
	if m.opts.Watch {
		live := lipgloss.NewStyle().
			Bold(true).
			Foreground(m.styles.Theme.Background).
			Background(m.styles.Theme.Success).
			Padding(0, 1).
			Render("LIVE")
		title = lipgloss.JoinHorizontal(lipgloss.Center, title, " ", live)
	}
	themeName := m.styles.Muted.Render("theme: " + m.styles.Theme.Name)
	repo := ""
	if m.analysis != nil {
		repo = m.styles.Subtitle.Render(m.analysis.Repository.Name + " · " + m.analysis.Repository.DefaultBranch)
	}

	headerLeft := lipgloss.JoinVertical(lipgloss.Left, title, repo)
	headerRight := themeName
	if m.opts.Watch {
		headerRight = m.styles.Success.Render("Live Mode") + m.styles.Muted.Render(fmt.Sprintf(" · every %s · refreshed %s",
			m.opts.WatchInterval, m.lastRefresh.Format("15:04:05")))
	}
	header := lipgloss.JoinHorizontal(lipgloss.Top,
		headerLeft,
		strings.Repeat(" ", max(1, m.width-lipgloss.Width(headerLeft)-lipgloss.Width(headerRight)-2)),
		headerRight,
	)

	tabBar := m.renderTabs()
	body := m.viewport.View()
	if m.loading && m.analysis == nil {
		body = m.styles.Muted.Render("\n  Analyzing repository…")
	}
	if m.err != nil && m.analysis == nil {
		body = m.styles.Error.Render("\n  Error: " + m.err.Error())
	}

	footer := m.styles.Muted.Render("tab/←→ switch · ↑↓ scroll · r refresh · q quit")
	if m.opts.Watch {
		footer = m.styles.Muted.Render("Live Mode · auto-refresh on · tab/←→ switch · ↑↓ scroll · r refresh now · q quit")
	}
	return lipgloss.JoinVertical(lipgloss.Left, header, tabBar, body, footer)
}

func (m model) renderTabs() string {
	var parts []string
	for i, name := range m.tabs {
		style := m.styles.Muted
		if i == m.tab {
			style = lipgloss.NewStyle().
				Foreground(m.styles.Theme.Background).
				Background(m.styles.Theme.Primary).
				Bold(true).
				Padding(0, 1)
		} else {
			style = style.Padding(0, 1)
		}
		parts = append(parts, style.Render(name))
	}
	return strings.Join(parts, " ")
}

func (m model) renderTab() string {
	if m.analysis == nil {
		if m.err != nil {
			return m.styles.Error.Render(m.err.Error())
		}
		return "Loading…"
	}
	a := m.analysis
	s := m.styles

	switch m.tab {
	case 0:
		w := m.width - 6
		if w > 56 {
			w = 56
		}
		return RenderBeautiful(a, s, w)
	case 1:
		return m.renderOverview(a, s)
	case 2:
		return m.renderCommits(a, s)
	case 3:
		body := strings.Join([]string{
			TopContributors(a.Contributors, 12),
			"",
			AuthorTimeline(a.AuthorTimelines, 9),
		}, "\n")
		return s.Panel.Width(m.width - 4).Render(body)
	case 4:
		return m.renderBranches(a, s)
	case 5:
		return s.Panel.Width(m.width-4).Render(LanguageStatistics(a.Languages, 4))
	case 6:
		w := m.width - 10
		if w < 40 {
			w = 40
		}
		body := strings.Join([]string{
			s.Title.Render("Commit Activity Graph"),
			MonthlyActivityGraph(a.Monthly, 17),
			"",
			AuthorTimeline(a.AuthorTimelines, 9),
			"",
			s.Title.Render("Daily Timeline (90 days)"),
			LineChart(a.Timeline, 10, w),
		}, "\n")
		return s.Panel.Width(m.width - 4).Render(body)
	case 7:
		body := strings.Join([]string{
			WeeklyHeatmap(a.CommitStats.CommitsByWeekday, 10),
			"",
			s.Title.Render("Contribution Calendar"),
			HeatmapASCII(a.Heatmap),
		}, "\n")
		return s.Panel.Width(m.width - 4).Render(body)
	case 8:
		return m.renderHealth(a, s)
	default:
		return ""
	}
}

func (m model) renderOverview(a *git.Analysis, s theme.Styles) string {
	overviewBlock := Render(a, s)
	recent := RecentCommitsList(a.RecentCommits, 8)

	left := s.Panel.Width(max(42, m.width/2-3)).Render(overviewBlock)
	right := s.Panel.Width(max(40, m.width/2-3)).Render(recent)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m model) renderCommits(a *git.Analysis, s theme.Styles) string {
	stats := strings.Join([]string{
		s.Title.Render("Commit Statistics"),
		fmt.Sprintf("Total: %s", s.Value.Render(fmt.Sprintf("%d", a.CommitStats.TotalCommits))),
		fmt.Sprintf("Avg/day: %s", s.Value.Render(fmt.Sprintf("%.2f", a.CommitStats.AvgCommitsPerDay))),
		fmt.Sprintf("Additions: %s", s.Success.Render(fmt.Sprintf("+%d", a.CommitStats.Additions))),
		fmt.Sprintf("Deletions: %s", s.Error.Render(fmt.Sprintf("-%d", a.CommitStats.Deletions))),
		fmt.Sprintf("Files touched: %s", s.Value.Render(fmt.Sprintf("%d", a.CommitStats.FilesChanged))),
		"",
		CommitTypes(a.CommitTypes, 12),
	}, "\n")

	files := strings.Join([]string{
		FileChangeStatistics(a.FileChangeSummary),
		"",
		s.Title.Render("Top Changed Files"),
	}, "\n")
	for i, f := range a.FileChanges {
		if i >= 10 {
			break
		}
		files += fmt.Sprintf("\n%s  %s (+%d/-%d)",
			s.Highlight.Render(fmt.Sprintf("%3d", f.Changes)),
			truncate(f.Path, 40),
			f.Additions, f.Deletions,
		)
	}

	left := s.Panel.Width(max(40, m.width/2-3)).Render(stats)
	right := s.Panel.Width(max(40, m.width/2-3)).Render(files)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m model) renderBranches(a *git.Analysis, s theme.Styles) string {
	return s.Panel.Width(m.width - 4).Render(BranchViewer(a.Branches, 0))
}

func (m model) renderHealth(a *git.Analysis, s theme.Styles) string {
	h := a.Health
	scoreColor := s.Success
	if h.Score < 70 {
		scoreColor = s.Warning
	}
	if h.Score < 50 {
		scoreColor = s.Error
	}

	prod := ProductivityScore(a.Productivity)
	streakCard := GitStreak(a.Streak)

	body := strings.Join([]string{
		s.Title.Render("Repository Health"),
		scoreColor.Render(fmt.Sprintf("Score: %.1f / 100  Grade: %s", h.Score, h.Grade)),
		"",
		fmt.Sprintf("%s %s %.0f", s.Label.Render("Frequency"), s.ProgressBar(h.CommitFrequency, 24), h.CommitFrequency),
		fmt.Sprintf("%s %s %.0f", s.Label.Render("Diversity"), s.ProgressBar(h.ContributorDiversity, 24), h.ContributorDiversity),
		fmt.Sprintf("%s %s %.0f", s.Label.Render("Branches"), s.ProgressBar(h.BranchHygiene, 24), h.BranchHygiene),
		fmt.Sprintf("%s %s %.0f", s.Label.Render("Churn"), s.ProgressBar(h.CodeChurn, 24), h.CodeChurn),
		fmt.Sprintf("%s %s %.0f", s.Label.Render("Recency"), s.ProgressBar(h.ActivityRecency, 24), h.ActivityRecency),
	}, "\n")

	left := s.Panel.Width(max(40, m.width/2-3)).Render(prod + "\n\n" + streakCard)
	right := s.Panel.Width(max(40, m.width/2-3)).Render(body)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m model) loadAnalysis() tea.Cmd {
	opts := m.opts
	return func() tea.Msg {
		az, err := git.Open(git.Options{
			Path:  opts.RepoPath,
			Limit: opts.Limit,
		})
		if err != nil {
			return analysisMsg{err: err}
		}
		a, err := az.Analyze()
		return analysisMsg{analysis: a, err: err}
	}
}

func tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n <= 1 {
		return s[:n]
	}
	return s[:n-1] + "…"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
