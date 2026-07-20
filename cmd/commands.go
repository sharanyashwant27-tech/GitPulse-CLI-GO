package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newCommandsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "commands",
		Short: "List primary CLI commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(CLICommands())
			return nil
		},
	}
}

// CLICommands renders the primary command reference card.
func CLICommands() string {
	cmds := []string{
		"gitpulse",
		"gitpulse stats",
		"gitpulse commits",
		"gitpulse branches",
		"gitpulse contributors",
		"gitpulse timeline",
		"gitpulse graph",
		"gitpulse health",
		"gitpulse export html",
		"gitpulse export pdf",
		"gitpulse watch",
	}
	var b strings.Builder
	b.WriteString("CLI Commands")
	for _, c := range cmds {
		b.WriteString("\n\n")
		b.WriteString(c)
	}
	return b.String()
}
