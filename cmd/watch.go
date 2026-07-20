package cmd

import (
	"github.com/spf13/cobra"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/dashboard"
)

func runDashboard(watch bool) error {
	p := dashboard.NewDashboard(dashboard.Options{
		RepoPath:      shared.Cfg.RepoPath,
		Theme:         shared.Cfg.Theme,
		Watch:         watch,
		WatchInterval: shared.Cfg.WatchInterval,
		Limit:         shared.Cfg.Limit,
	})
	_, err := p.Run()
	return err
}

func newWatchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "watch",
		Short: "Live Mode — auto-refreshing interactive dashboard",
		Long: `Live Mode

  gitpulse watch

Opens the interactive Bubble Tea dashboard and re-analyzes the repository
automatically on an interval (default: 5s).

Examples:
  gitpulse watch
  gitpulse watch --interval 3s
  gitpulse watch --theme nord -r /path/to/repo`,
		Example: `  gitpulse watch
  gitpulse watch --interval 3s
  gitpulse watch --theme catppuccin`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDashboard(true)
		},
	}
}
