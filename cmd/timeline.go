package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/dashboard"
)

func newTimelineCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "timeline",
		Short: "Show commit activity timeline and weekly heatmap",
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := shared.AnalyzeRepo()
			if err != nil {
				return err
			}
			s := shared.Styles()
			fmt.Println(s.Title.Render("GitPulse · Timeline"))
			fmt.Println(dashboard.MonthlyActivityGraph(a.Monthly, 17))
			fmt.Println()
			fmt.Println(dashboard.WeeklyHeatmap(a.CommitStats.CommitsByWeekday, 10))
			fmt.Println()
			fmt.Println(dashboard.LineChart(a.Timeline, 12, 72))
			fmt.Println()
			fmt.Println(dashboard.HeatmapASCII(a.Heatmap))
			fmt.Println()
			fmt.Println(s.Subtitle.Render("Last 14 days"))
			start := len(a.Timeline) - 14
			if start < 0 {
				start = 0
			}
			for _, p := range a.Timeline[start:] {
				bar := strings.Repeat("█", p.Commits)
				if p.Commits == 0 {
					bar = "·"
				}
				fmt.Printf("  %s  %3d  %s\n", p.Date, p.Commits, s.Bar.Render(bar))
			}
			return nil
		},
	}
}
