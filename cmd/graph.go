package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/dashboard"
)

func newGraphCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "graph",
		Short: "Render ASCII charts for activity, languages, and contributors",
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := shared.AnalyzeRepo()
			if err != nil {
				return err
			}
			s := shared.Styles()
			fmt.Println(s.Title.Render("GitPulse · Graphs"))
			fmt.Println(s.Subtitle.Render("Commit Activity Graph"))
			fmt.Println(dashboard.MonthlyActivityGraph(a.Monthly, 17))
			fmt.Println()
			fmt.Println(s.Subtitle.Render("Commit activity (90 days)"))
			fmt.Println(dashboard.LineChart(a.Timeline, 10, 70))
			fmt.Println()
			fmt.Println(dashboard.LanguageStatistics(a.Languages, 4))
			fmt.Println()
			fmt.Println(dashboard.TopContributors(a.Contributors, 10))
			fmt.Println()
			fmt.Println(dashboard.AuthorTimeline(a.AuthorTimelines, 9))
			fmt.Println()
			fmt.Println(s.Subtitle.Render("Commit share"))
			fmt.Println(dashboard.ContributorBars(a.Contributors, 10))
			fmt.Println()
			fmt.Println(dashboard.WeeklyHeatmap(a.CommitStats.CommitsByWeekday, 10))
			fmt.Println()
			fmt.Println(dashboard.CommitTypes(a.CommitTypes, 12))
			fmt.Println()
			fmt.Println(dashboard.BranchViewer(a.Branches, 12))
			fmt.Println()
			fmt.Println(dashboard.RecentCommitsList(a.RecentCommits, 8))
			fmt.Println()
			fmt.Println(dashboard.FileChangeStatistics(a.FileChangeSummary))
			fmt.Println()
			fmt.Println(dashboard.ProductivityScore(a.Productivity))
			fmt.Println()
			fmt.Println(dashboard.GitStreak(a.Streak))
			fmt.Println()
			fmt.Println(s.Subtitle.Render("Contribution calendar"))
			fmt.Println(dashboard.HeatmapASCII(a.Heatmap))
			fmt.Println()
			fmt.Println(s.Subtitle.Render("Commits by weekday"))
			fmt.Println(dashboard.WeekdayBars(a.CommitStats.CommitsByWeekday))
			return nil
		},
	}
}
