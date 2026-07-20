package utils

import (
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

// AnalysisProgress shows a spinner while analysis runs.
type AnalysisProgress struct {
	bar *progressbar.ProgressBar
}

// StartAnalysisProgress begins an indeterminate progress spinner on stderr.
func StartAnalysisProgress(description string) *AnalysisProgress {
	if description == "" {
		description = "Analyzing repository"
	}
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionThrottle(80*time.Millisecond),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetRenderBlankState(true),
	)
	return &AnalysisProgress{bar: bar}
}

// Tick advances the spinner once.
func (p *AnalysisProgress) Tick() {
	if p == nil || p.bar == nil {
		return
	}
	_ = p.bar.Add(1)
}

// Finish clears the spinner.
func (p *AnalysisProgress) Finish() {
	if p == nil || p.bar == nil {
		return
	}
	_ = p.bar.Finish()
}
