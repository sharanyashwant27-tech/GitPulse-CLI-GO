package exporter

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/git"
)

func (e *Exporter) writePDF(path string, analysis *git.Analysis) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("GitPulse Report — "+analysis.Repository.Name, false)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 10, "GitPulse Repository Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 11)
	writePDFLine(pdf, "Repository", analysis.Repository.Name)
	writePDFLine(pdf, "Path", analysis.Repository.Path)
	writePDFLine(pdf, "Branch", analysis.Repository.DefaultBranch)
	writePDFLine(pdf, "Status", analysis.Repository.Status)
	writePDFLine(pdf, "Remote", analysis.Repository.RemoteURL)
	writePDFLine(pdf, "Branches", fmt.Sprintf("%d", analysis.Repository.BranchCount))
	writePDFLine(pdf, "Tags", fmt.Sprintf("%d", analysis.Repository.TagCount))
	writePDFLine(pdf, "Total commits", fmt.Sprintf("%d", analysis.CommitStats.TotalCommits))
	writePDFLine(pdf, "Contributors", fmt.Sprintf("%d", analysis.CommitStats.TotalAuthors))
	writePDFLine(pdf, "Health score", fmt.Sprintf("%.0f%% (%s)", analysis.Health.Score, analysis.Health.Grade))
	writePDFLine(pdf, "Current streak", fmt.Sprintf("%d days", analysis.Streak.Current))
	writePDFLine(pdf, "Longest streak", fmt.Sprintf("%d days", analysis.Streak.Longest))
	writePDFLine(pdf, "Generated", analysis.GeneratedAt.Format(time.RFC3339))

	pdf.Ln(6)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "Top Contributors")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	limit := 10
	if len(analysis.Contributors) < limit {
		limit = len(analysis.Contributors)
	}
	for i := 0; i < limit; i++ {
		c := analysis.Contributors[i]
		pdf.Cell(0, 6, fmt.Sprintf("%d. %s — %d commits (%.1f%%)", i+1, c.Name, c.Commits, c.Percentage))
		pdf.Ln(6)
	}

	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "Languages")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	limit = 8
	if len(analysis.Languages) < limit {
		limit = len(analysis.Languages)
	}
	for i := 0; i < limit; i++ {
		l := analysis.Languages[i]
		pdf.Cell(0, 6, fmt.Sprintf("%s — %.1f%% (%d files)", l.Name, l.Percentage, l.Files))
		pdf.Ln(6)
	}

	return pdf.OutputFileAndClose(path)
}

func writePDFLine(pdf *gofpdf.Fpdf, label, value string) {
	if value == "" {
		value = "-"
	}
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(45, 7, label+":")
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 7, value)
	pdf.Ln(7)
}
