package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/dashboard"
)

func newStatsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show repository overview and commit statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := shared.AnalyzeRepo()
			if err != nil {
				return err
			}
			s := shared.Styles()
			fmt.Println(dashboard.RenderBeautiful(a, s, 46))
			fmt.Println()
			fmt.Println(dashboard.Render(a, s))
			fmt.Println()
			fmt.Println(s.Subtitle.Render("Commit statistics"))
			fmt.Printf("%s %d\n", s.Label.Render("Total commits"), a.CommitStats.TotalCommits)
			fmt.Printf("%s %d\n", s.Label.Render("Authors"), a.CommitStats.TotalAuthors)
			fmt.Printf("%s %.2f\n", s.Label.Render("Avg/day"), a.CommitStats.AvgCommitsPerDay)
			fmt.Printf("%s %s / %s\n", s.Label.Render("Diff"),
				s.Success.Render(fmt.Sprintf("+%d", a.CommitStats.Additions)),
				s.Error.Render(fmt.Sprintf("-%d", a.CommitStats.Deletions)),
			)
			fmt.Printf("%s %d\n", s.Label.Render("Files changed"), a.CommitStats.FilesChanged)
			if !a.CommitStats.FirstCommit.IsZero() {
				fmt.Printf("%s %s → %s\n", s.Label.Render("Span"),
					a.CommitStats.FirstCommit.Format("2006-01-02"),
					a.CommitStats.LastCommit.Format("2006-01-02"),
				)
			}
			fmt.Println()
			fmt.Println(dashboard.CommitTypes(a.CommitTypes, 12))
			fmt.Println()
			fmt.Println(dashboard.RecentCommitsList(a.RecentCommits, shared.Cfg.Limit))
			fmt.Println()
			fmt.Println(dashboard.FileChangeStatistics(a.FileChangeSummary))
			fmt.Println()
			fmt.Println(s.Subtitle.Render("Top changed files"))
			for i, f := range a.FileChanges {
				if i >= 10 {
					break
				}
				fmt.Printf("  %3d  %s  (+%d/-%d)\n", f.Changes, f.Path, f.Additions, f.Deletions)
			}
			return nil
		},
	}
}
