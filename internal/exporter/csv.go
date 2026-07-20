package exporter

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/git"
)

func (e *Exporter) writeCSV(path string, analysis *git.Analysis) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	_ = w.Write([]string{"section", "key", "value"})
	_ = w.Write([]string{"repository", "name", analysis.Repository.Name})
	_ = w.Write([]string{"repository", "path", analysis.Repository.Path})
	_ = w.Write([]string{"repository", "branch", analysis.Repository.DefaultBranch})
	_ = w.Write([]string{"stats", "total_commits", fmt.Sprintf("%d", analysis.CommitStats.TotalCommits)})
	_ = w.Write([]string{"stats", "total_authors", fmt.Sprintf("%d", analysis.CommitStats.TotalAuthors)})
	_ = w.Write([]string{"stats", "additions", fmt.Sprintf("%d", analysis.CommitStats.Additions)})
	_ = w.Write([]string{"stats", "deletions", fmt.Sprintf("%d", analysis.CommitStats.Deletions)})
	_ = w.Write([]string{"file_changes", "added", fmt.Sprintf("%d", analysis.FileChangeSummary.Added)})
	_ = w.Write([]string{"file_changes", "modified", fmt.Sprintf("%d", analysis.FileChangeSummary.Modified)})
	_ = w.Write([]string{"file_changes", "deleted", fmt.Sprintf("%d", analysis.FileChangeSummary.Deleted)})
	_ = w.Write([]string{"file_changes", "renamed", fmt.Sprintf("%d", analysis.FileChangeSummary.Renamed)})
	_ = w.Write([]string{"health", "score", fmt.Sprintf("%.1f", analysis.Health.Score)})
	_ = w.Write([]string{"health", "grade", analysis.Health.Grade})
	_ = w.Write([]string{"productivity", "label", analysis.Productivity.Label})
	_ = w.Write([]string{"productivity", "score", fmt.Sprintf("%.0f", analysis.Productivity.Score)})
	_ = w.Write([]string{"streak", "current", fmt.Sprintf("%d", analysis.Streak.Current)})
	_ = w.Write([]string{"streak", "longest", fmt.Sprintf("%d", analysis.Streak.Longest)})

	for _, c := range analysis.Contributors {
		_ = w.Write([]string{"contributor", c.Name, fmt.Sprintf("%d commits (%.1f%%)", c.Commits, c.Percentage)})
	}
	for _, b := range analysis.Branches {
		kind := "local"
		if b.IsRemote {
			kind = "remote"
		}
		_ = w.Write([]string{"branch", b.Name, kind})
	}
	for _, l := range analysis.Languages {
		_ = w.Write([]string{"language", l.Name, fmt.Sprintf("%.1f%%", l.Percentage)})
	}
	for _, t := range analysis.CommitTypes {
		_ = w.Write([]string{"commit_type", t.Name, fmt.Sprintf("%d (%.0f%%)", t.Count, t.Percentage)})
	}
	for _, rc := range analysis.RecentCommits {
		_ = w.Write([]string{"commit", rc.ShortHash, rc.Message})
	}
	return w.Error()
}
