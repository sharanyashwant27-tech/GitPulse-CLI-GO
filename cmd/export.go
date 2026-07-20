package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/exporter"
	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/utils"
)

func newExportCmd() *cobra.Command {
	var (
		format   string
		filename string
		outDir   string
	)
	cmd := &cobra.Command{
		Use:   "export [html|pdf|csv|json]",
		Short: "Export analysis to HTML, JSON, CSV, or PDF",
		Long: `Export Reports

  gitpulse export html
  gitpulse export pdf
  gitpulse export csv
  gitpulse export json

You can also use --format / -f for the same formats.`,
		Example: `  gitpulse export html
  gitpulse export pdf
  gitpulse export csv
  gitpulse export json
  gitpulse export html -o my-report
  gitpulse export -f json --dir reports`,
		Args:      cobra.MaximumNArgs(1),
		ValidArgs: []string{"html", "pdf", "csv", "json"},
		RunE: func(cmd *cobra.Command, args []string) error {
			selected := strings.ToLower(format)
			if len(args) == 1 {
				selected = strings.ToLower(args[0])
			}
			if selected == "" {
				selected = "html"
			}
			switch selected {
			case "html", "pdf", "csv", "json":
			default:
				return fmt.Errorf("unsupported format %q (use: html, pdf, csv, json)", selected)
			}

			a, err := shared.AnalyzeRepo()
			if err != nil {
				return err
			}
			dir := shared.Cfg.ExportDir
			if outDir != "" {
				dir = outDir
			}
			exp := exporter.New(utils.FindTemplatesDir(), dir)
			path, err := exp.Export(a, exporter.Format(selected), filename)
			if err != nil {
				return err
			}
			utils.Successf("✓ Exported report to %s\n", path)
			return nil
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "", "export format: html, json, csv, pdf (or pass as argument)")
	cmd.Flags().StringVarP(&filename, "output", "o", "", "output filename (without extension)")
	cmd.Flags().StringVar(&outDir, "dir", "", "output directory (overrides --export-dir)")
	return cmd
}
