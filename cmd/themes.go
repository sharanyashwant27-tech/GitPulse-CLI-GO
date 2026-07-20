package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/theme"
)

func newThemesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "themes",
		Short: "List available color themes",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(theme.TerminalThemes())
			if shared.Cfg != nil {
				fmt.Println()
				fmt.Println("Active:", theme.Get(shared.Cfg.Theme).Name)
				fmt.Println("Use: gitpulse --theme <name>")
			}
			return nil
		},
	}
}
