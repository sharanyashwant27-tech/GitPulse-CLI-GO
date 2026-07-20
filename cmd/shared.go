package cmd

import (
	"fmt"
	"time"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/config"
	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/git"
	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/theme"
	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/utils"
)

// Shared holds global CLI state shared across commands.
type Shared struct {
	Cfg     *config.Config
	CfgFile string
}

// AnalyzeRepo runs analysis for the configured repository.
func (s *Shared) AnalyzeRepo() (*git.Analysis, error) {
	opts := git.Options{
		Path:  s.Cfg.RepoPath,
		Limit: s.Cfg.Limit,
	}
	if s.Cfg.Since != "" {
		t, err := time.Parse("2006-01-02", s.Cfg.Since)
		if err != nil {
			return nil, fmt.Errorf("invalid --since date (use YYYY-MM-DD): %w", err)
		}
		opts.Since = t
	}
	if s.Cfg.Until != "" {
		t, err := time.Parse("2006-01-02", s.Cfg.Until)
		if err != nil {
			return nil, fmt.Errorf("invalid --until date (use YYYY-MM-DD): %w", err)
		}
		opts.Until = t
	}

	progress := utils.StartAnalysisProgress("Analyzing repository")
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				progress.Tick()
			}
		}
	}()

	az, err := git.Open(opts)
	if err != nil {
		close(done)
		progress.Finish()
		return nil, err
	}
	analysis, err := az.Analyze()
	close(done)
	progress.Finish()
	if err != nil {
		return nil, err
	}
	return analysis, nil
}

// Styles returns themed lipgloss styles for the current config.
func (s *Shared) Styles() theme.Styles {
	return theme.NewStyles(theme.Get(s.Cfg.Theme))
}
