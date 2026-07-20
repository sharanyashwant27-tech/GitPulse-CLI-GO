package cmd

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/config"
	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/utils"
)

var (
	cfgFile string
	shared  = &Shared{}
	version = "1.0.0"
)

// Execute runs the root command.
func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "gitpulse",
		Short:         "Interactive Git repository analytics for your terminal",
		Long: `GitPulse is a production-grade CLI that analyzes local Git repositories and renders interactive dashboards, charts, and exportable reports.

CLI Commands

  gitpulse
  gitpulse stats
  gitpulse commits
  gitpulse branches
  gitpulse contributors
  gitpulse timeline
  gitpulse graph
  gitpulse health
  gitpulse export html
  gitpulse export pdf
  gitpulse watch`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return err
			}

			if v, err := cmd.Flags().GetString("repo"); err == nil && cmd.Flags().Changed("repo") {
				abs, absErr := filepath.Abs(v)
				if absErr != nil {
					return absErr
				}
				cfg.RepoPath = abs
			}
			if v, err := cmd.Flags().GetString("theme"); err == nil && cmd.Flags().Changed("theme") {
				cfg.Theme = v
			}
			if v, err := cmd.Flags().GetString("log-level"); err == nil && cmd.Flags().Changed("log-level") {
				cfg.LogLevel = v
			}
			if v, err := cmd.Flags().GetInt("limit"); err == nil && cmd.Flags().Changed("limit") {
				cfg.Limit = v
			}
			if v, err := cmd.Flags().GetString("since"); err == nil && cmd.Flags().Changed("since") {
				cfg.Since = v
			}
			if v, err := cmd.Flags().GetString("until"); err == nil && cmd.Flags().Changed("until") {
				cfg.Until = v
			}
			if v, err := cmd.Flags().GetDuration("interval"); err == nil && cmd.Flags().Changed("interval") {
				cfg.WatchInterval = v
			}
			if v, err := cmd.Flags().GetString("export-dir"); err == nil && cmd.Flags().Changed("export-dir") {
				cfg.ExportDir = v
			}
			if v, err := cmd.Flags().GetBool("verbose"); err == nil && cmd.Flags().Changed("verbose") {
				cfg.Verbose = v
				if v {
					cfg.LogLevel = "debug"
				}
			}

			if err := config.ValidateTheme(cfg.Theme); err != nil {
				return err
			}

			shared.Cfg = cfg
			shared.CfgFile = cfgFile

			dev := cfg.Verbose || cfg.LogLevel == "debug"
			if err := utils.Init(cfg.LogLevel, dev); err != nil {
				return err
			}
			utils.L().Debug("config loaded",
				zap.String("repo", cfg.RepoPath),
				zap.String("theme", cfg.Theme),
			)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDashboard(false)
		},
	}

	root.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ./gitpulse.yaml or ~/.config/gitpulse/gitpulse.yaml)")
	root.PersistentFlags().StringP("repo", "r", "", "path to git repository (default: current directory)")
	root.PersistentFlags().StringP("theme", "t", "default", "color theme: default, dracula, nord, catppuccin, tokyo-night, gruvbox, solarized")
	root.PersistentFlags().String("log-level", "info", "log level: debug, info, warn, error")
	root.PersistentFlags().IntP("limit", "n", 20, "limit for lists (commits, contributors, files)")
	root.PersistentFlags().String("since", "", "only include commits after YYYY-MM-DD")
	root.PersistentFlags().String("until", "", "only include commits before YYYY-MM-DD")
	root.PersistentFlags().Duration("interval", 5*time.Second, "watch refresh interval")
	root.PersistentFlags().String("export-dir", "reports", "directory for exported reports")
	root.PersistentFlags().BoolP("verbose", "v", false, "enable verbose logging")

	root.AddCommand(
		newStatsCmd(),
		newCommitsCmd(),
		newBranchesCmd(),
		newContributorsCmd(),
		newTimelineCmd(),
		newGraphCmd(),
		newHealthCmd(),
		newExportCmd(),
		newWatchCmd(),
		newThemesCmd(),
		newCommandsCmd(),
	)

	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)
	return root
}
