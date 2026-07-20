package exporter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/git"
)

// Format identifies an export target format.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
	FormatHTML Format = "html"
	FormatPDF  Format = "pdf"
)

// Exporter writes analysis reports to disk.
type Exporter struct {
	TemplateDir string
	OutputDir   string
}

// New creates an Exporter.
func New(templateDir, outputDir string) *Exporter {
	return &Exporter{TemplateDir: templateDir, OutputDir: outputDir}
}

// Export writes the analysis in the requested format and returns the output path.
func (e *Exporter) Export(analysis *git.Analysis, format Format, filename string) (string, error) {
	if err := os.MkdirAll(e.OutputDir, 0o755); err != nil {
		return "", err
	}
	if filename == "" {
		filename = fmt.Sprintf("gitpulse-%s-%s", analysis.Repository.Name, time.Now().Format("20060102-150405"))
	}
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))

	switch format {
	case FormatJSON:
		path := filepath.Join(e.OutputDir, filename+".json")
		return path, e.writeJSON(path, analysis)
	case FormatCSV:
		path := filepath.Join(e.OutputDir, filename+".csv")
		return path, e.writeCSV(path, analysis)
	case FormatHTML:
		path := filepath.Join(e.OutputDir, filename+".html")
		return path, e.writeHTML(path, analysis)
	case FormatPDF:
		path := filepath.Join(e.OutputDir, filename+".pdf")
		return path, e.writePDF(path, analysis)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}
