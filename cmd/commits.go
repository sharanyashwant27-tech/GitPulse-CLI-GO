package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/dashboard"
)

func newCommitsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "commits",
		Short: "Show recent commits, commit types, and file change stats",
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := shared.AnalyzeRepo()
			if err != nil {
				return err
			}
			fmt.Println(dashboard.RecentCommitsList(a.RecentCommits, shared.Cfg.Limit))
			fmt.Println()
			fmt.Println(dashboard.CommitTypes(a.CommitTypes, 12))
			fmt.Println()
			fmt.Println(dashboard.FileChangeStatistics(a.FileChangeSummary))
			fmt.Println()
			fmt.Println(dashboard.MonthlyActivityGraph(a.Monthly, 17))
			return nil
		},
	}
}
