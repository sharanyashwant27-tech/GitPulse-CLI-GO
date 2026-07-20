package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/dashboard"
	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/git"
	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/theme"
)

func newContributorsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "contributors",
		Short: "Show the contributor leaderboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := shared.AnalyzeRepo()
			if err != nil {
				return err
			}
			s := shared.Styles()
			fmt.Println(dashboard.TopContributors(a.Contributors, shared.Cfg.Limit))
			fmt.Println()
			fmt.Println(renderContributorsTable(a.Contributors, shared.Cfg.Limit, s))
			fmt.Println()
			fmt.Println(dashboard.AuthorTimeline(a.AuthorTimelines, 9))
			fmt.Println()
			fmt.Println(s.Subtitle.Render("Commit share"))
			fmt.Println(dashboard.ContributorBars(a.Contributors, 10))
			return nil
		},
	}
}

func renderContributorsTable(contributors []git.Contributor, limit int, s theme.Styles) string {
	if len(contributors) == 0 {
		return "(no contributors)"
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > len(contributors) {
		limit = len(contributors)
	}

	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Name", Width: 22},
		{Title: "Commits", Width: 10},
		{Title: "%", Width: 8},
	}
	rows := make([]table.Row, 0, limit)
	for i := 0; i < limit; i++ {
		c := contributors[i]
		name := c.Name
		if name == "" {
			name = c.Email
		}
		rows = append(rows, table.Row{
			strconv.Itoa(i + 1),
			name,
			strconv.Itoa(c.Commits),
			fmt.Sprintf("%.1f%%", c.Percentage),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(limit+1),
	)
	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(s.Theme.Border).
		BorderBottom(true).
		Bold(true).
		Foreground(s.Theme.Accent)
	styles.Selected = styles.Selected.
		Foreground(s.Theme.Foreground).
		Background(s.Theme.Border).
		Bold(false)
	styles.Cell = styles.Cell.Foreground(s.Theme.Foreground)
	t.SetStyles(styles)

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(s.Theme.Primary).
		Render(t.View())
}
