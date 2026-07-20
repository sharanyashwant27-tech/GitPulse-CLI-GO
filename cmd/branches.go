package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/dashboard"
)

func newBranchesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "branches",
		Short: "List local and remote branches",
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := shared.AnalyzeRepo()
			if err != nil {
				return err
			}
			s := shared.Styles()
			fmt.Println(dashboard.BranchViewer(a.Branches, 0))
			fmt.Println()
			fmt.Println(s.Subtitle.Render("Branch details"))
			fmt.Println(s.TableHead.Render(fmt.Sprintf("%-32s %-8s %-12s %-20s",
				"NAME", "TYPE", "LAST COMMIT", "AUTHOR")))
			for _, b := range a.Branches {
				kind := "local"
				if b.IsRemote {
					kind = "remote"
				}
				name := b.Name
				if b.IsCurrent {
					name = "* " + name
				}
				date := "-"
				if !b.LastCommit.IsZero() {
					date = b.LastCommit.Format("2006-01-02")
				}
				line := fmt.Sprintf("%-32s %-8s %-12s %-20s", truncate(name, 32), kind, date, truncate(b.Author, 20))
				if b.IsCurrent {
					fmt.Println(s.Success.Render(line))
				} else {
					fmt.Println(line)
				}
			}
			return nil
		},
	}
}
