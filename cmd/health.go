package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/dashboard"
)

func newHealthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Compute repository health score and git streak",
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := shared.AnalyzeRepo()
			if err != nil {
				return err
			}
			s := shared.Styles()
			h := a.Health
			fmt.Println(dashboard.ProductivityScore(a.Productivity))
			fmt.Println()
			fmt.Println(dashboard.GitStreak(a.Streak))
			fmt.Println()
			fmt.Println(s.Title.Render("GitPulse · Health"))
			scoreStyle := s.Success
			if h.Score < 70 {
				scoreStyle = s.Warning
			}
			if h.Score < 50 {
				scoreStyle = s.Error
			}
			fmt.Println(scoreStyle.Render(fmt.Sprintf("Score: %.1f / 100   Grade: %s", h.Score, h.Grade)))
			fmt.Println()
			fmt.Printf("%s %s %.0f\n", s.Label.Render("Frequency"), s.ProgressBar(h.CommitFrequency, 28), h.CommitFrequency)
			fmt.Printf("%s %s %.0f\n", s.Label.Render("Diversity"), s.ProgressBar(h.ContributorDiversity, 28), h.ContributorDiversity)
			fmt.Printf("%s %s %.0f\n", s.Label.Render("Branches"), s.ProgressBar(h.BranchHygiene, 28), h.BranchHygiene)
			fmt.Printf("%s %s %.0f\n", s.Label.Render("Churn"), s.ProgressBar(h.CodeChurn, 28), h.CodeChurn)
			fmt.Printf("%s %s %.0f\n", s.Label.Render("Recency"), s.ProgressBar(h.ActivityRecency, 28), h.ActivityRecency)
			fmt.Println()
			if !a.Streak.LastActive.IsZero() {
				fmt.Printf("%s %s\n", s.Label.Render("Last active"), a.Streak.LastActive.Format(time.RFC822))
			}
			fmt.Println()
			fmt.Println(s.Subtitle.Render("Details"))
			for k, v := range h.Details {
				fmt.Printf("  %s: %s\n", k, v)
			}
			return nil
		},
	}
}
