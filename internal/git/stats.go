package git

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func (a *Analyzer) languages() ([]LanguageStat, error) {
	wt, err := a.repo.Worktree()
	root := a.path
	if err == nil && wt != nil {
		root = wt.Filesystem.Root()
	}

	counts := map[string]struct {
		bytes int64
		files int
	}{}
	var total int64

	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		if d.IsDir() {
			if name == ".git" || name == "node_modules" || name == "vendor" || name == "dist" || name == "build" {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(name))
		lang := languageFromExt(ext)
		if lang == "" {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		c := counts[lang]
		c.bytes += info.Size()
		c.files++
		counts[lang] = c
		total += info.Size()
		return nil
	})

	out := make([]LanguageStat, 0, len(counts))
	for name, c := range counts {
		pct := 0.0
		if total > 0 {
			pct = (float64(c.bytes) / float64(total)) * 100
		}
		out = append(out, LanguageStat{
			Name:       name,
			Bytes:      c.bytes,
			Files:      c.files,
			Percentage: pct,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Bytes > out[j].Bytes
	})
	return out, nil
}

func (a *Analyzer) fileChanges(commits []*object.Commit, limit int) []FileChangeStat {
	byPath := map[string]*FileChangeStat{}
	max := 100
	if len(commits) < max {
		max = len(commits)
	}
	for i := 0; i < max; i++ {
		stats, err := commits[i].Stats()
		if err != nil {
			continue
		}
		for _, s := range stats {
			fc, ok := byPath[s.Name]
			if !ok {
				fc = &FileChangeStat{Path: s.Name}
				byPath[s.Name] = fc
			}
			fc.Changes++
			fc.Additions += s.Addition
			fc.Deletions += s.Deletion
		}
	}

	out := make([]FileChangeStat, 0, len(byPath))
	for _, fc := range byPath {
		out = append(out, *fc)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Changes > out[j].Changes
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}

func (a *Analyzer) fileChangeSummary(commits []*object.Commit) FileChangeSummary {
	var summary FileChangeSummary
	max := 150
	if len(commits) < max {
		max = len(commits)
	}
	for i := 0; i < max; i++ {
		c := commits[i]
		var parent *object.Commit
		if c.NumParents() > 0 {
			p, err := c.Parent(0)
			if err != nil {
				continue
			}
			parent = p
		}

		patch, err := c.Patch(parent)
		if err != nil {
			stats, serr := c.Stats()
			if serr != nil {
				continue
			}
			for _, s := range stats {
				switch {
				case s.Addition > 0 && s.Deletion == 0:
					summary.Added++
				case s.Addition == 0 && s.Deletion > 0:
					summary.Deleted++
				default:
					summary.Modified++
				}
			}
			continue
		}

		for _, fp := range patch.FilePatches() {
			from, to := fp.Files()
			switch {
			case from == nil && to != nil:
				summary.Added++
			case from != nil && to == nil:
				summary.Deleted++
			case from != nil && to != nil && from.Path() != to.Path():
				summary.Renamed++
			default:
				summary.Modified++
			}
		}
	}
	return summary
}

func languageFromExt(ext string) string {
	m := map[string]string{
		".go": "Go", ".py": "Python", ".js": "JavaScript", ".ts": "TypeScript",
		".tsx": "TypeScript", ".jsx": "JavaScript", ".rs": "Rust", ".java": "Java",
		".c": "C", ".h": "C", ".cpp": "C++", ".cc": "C++", ".hpp": "C++",
		".cs": "C#", ".rb": "Ruby", ".php": "PHP", ".swift": "Swift",
		".kt": "Kotlin", ".scala": "Scala", ".md": "Markdown", ".html": "HTML",
		".css": "CSS", ".scss": "SCSS", ".json": "JSON", ".yaml": "YAML",
		".yml": "YAML", ".toml": "TOML", ".sh": "Shell", ".bash": "Shell",
		".sql": "SQL", ".vue": "Vue", ".svelte": "Svelte", ".dart": "Dart",
		".r": "R", ".lua": "Lua", ".ex": "Elixir", ".exs": "Elixir",
		".hs": "Haskell", ".zig": "Zig", ".nim": "Nim",
	}
	return m[ext]
}

// FormatBytes formats byte counts for display.
func FormatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
