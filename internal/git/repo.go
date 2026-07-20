package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Options controls analysis scope.
type Options struct {
	Path  string
	Since time.Time
	Until time.Time
	Limit int
}

// Analyzer inspects a local Git repository using go-git.
type Analyzer struct {
	repo *git.Repository
	path string
	opts Options
}

// Open opens a repository at path and returns an Analyzer.
func Open(opts Options) (*Analyzer, error) {
	if opts.Path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		opts.Path = cwd
	}
	if opts.Limit <= 0 {
		opts.Limit = 20
	}
	if opts.Until.IsZero() {
		opts.Until = time.Now()
	}

	repo, err := git.PlainOpenWithOptions(opts.Path, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, fmt.Errorf("open repository at %s: %w", opts.Path, err)
	}

	return &Analyzer{repo: repo, path: opts.Path, opts: opts}, nil
}

// Analyze performs a full repository analysis.
func (a *Analyzer) Analyze() (*Analysis, error) {
	commits, err := a.loadCommits()
	if err != nil {
		return nil, err
	}

	repoMeta, err := a.repositoryInfo(commits)
	if err != nil {
		return nil, err
	}

	stats := a.commitStats(commits)
	contributors := a.contributors(commits)
	branches, err := a.branches()
	if err != nil {
		return nil, err
	}
	languages, err := a.languages()
	if err != nil {
		return nil, err
	}
	if len(languages) > 0 {
		repoMeta.PrimaryLanguage = languages[0].Name
	}
	repoMeta.License = detectLicense(repoMeta.Path)
	repoMeta.Stars = "N/A"

	timeline := a.timeline(commits, 90)
	monthly := a.monthlyActivity(commits)
	authorTimelines := a.authorTimelines(commits, 5)
	commitTypes := classifyCommitTypes(commits)
	heatmap := a.heatmap(commits, 52)
	recent := a.recentCommits(commits, a.opts.Limit)
	files := a.fileChanges(commits, a.opts.Limit)
	summary := a.fileChangeSummary(commits)
	streak := a.streak(commits)
	h := Calculate(stats, contributors, branches, streak)
	prod := CalculateProductivity(stats, branches, commitTypes, languages, *repoMeta, h)

	return &Analysis{
		Repository:        *repoMeta,
		CommitStats:       stats,
		Contributors:      contributors,
		Branches:          branches,
		Languages:         languages,
		Timeline:          timeline,
		Monthly:           monthly,
		AuthorTimelines:   authorTimelines,
		CommitTypes:       commitTypes,
		Heatmap:           heatmap,
		RecentCommits:     recent,
		FileChanges:       files,
		FileChangeSummary: summary,
		Health:            h,
		Productivity:      prod,
		Streak:            streak,
		GeneratedAt:       time.Now().UTC(),
	}, nil
}

func (a *Analyzer) repositoryInfo(commits []*object.Commit) (*Repository, error) {
	name := filepath.Base(a.path)
	wt, err := a.repo.Worktree()
	root := a.path
	if err == nil && wt != nil {
		root = wt.Filesystem.Root()
		name = filepath.Base(root)
	}

	var remoteURL string
	if remotes, err := a.repo.Remotes(); err == nil && len(remotes) > 0 {
		cfg := remotes[0].Config()
		if len(cfg.URLs) > 0 {
			remoteURL = cfg.URLs[0]
		}
	}

	head, err := a.repo.Head()
	defaultBranch := "main"
	if err == nil {
		defaultBranch = head.Name().Short()
	}

	var first, last time.Time
	if len(commits) > 0 {
		last = commits[0].Author.When
		first = commits[len(commits)-1].Author.When
	}

	size, _ := dirSize(root)
	branchCount, tagCount := a.countRefs()
	status := a.worktreeStatus()

	return &Repository{
		Path:          root,
		Name:          name,
		RemoteURL:     remoteURL,
		DefaultBranch: defaultBranch,
		Status:        status,
		SizeBytes:     size,
		BranchCount:   branchCount,
		TagCount:      tagCount,
		Stars:         "N/A",
		CreatedAt:     first,
		LastCommitAt:  last,
	}, nil
}

func (a *Analyzer) worktreeStatus() string {
	wt, err := a.repo.Worktree()
	if err != nil || wt == nil {
		return "Clean"
	}
	st, err := wt.Status()
	if err != nil {
		return "Unknown"
	}
	if st.IsClean() {
		return "Clean"
	}
	return "Dirty"
}

func (a *Analyzer) countRefs() (branches, tags int) {
	iter, err := a.repo.References()
	if err != nil {
		return 0, 0
	}
	defer iter.Close()
	_ = iter.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name()
		switch {
		case name.IsBranch():
			branches++
		case name.IsTag():
			tags++
		}
		return nil
	})
	return branches, tags
}

func detectLicense(root string) string {
	candidates := []string{
		"LICENSE", "LICENSE.md", "LICENSE.txt", "LICENCE", "LICENCE.md",
		"COPYING", "COPYING.md", "license", "License",
	}
	for _, name := range candidates {
		path := filepath.Join(root, name)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		text := string(data)
		upper := strings.ToUpper(text)
		switch {
		case strings.Contains(upper, "MIT LICENSE") || strings.Contains(upper, "PERMISSION IS HEREBY GRANTED, FREE OF CHARGE"):
			return "MIT"
		case strings.Contains(upper, "APACHE LICENSE") && strings.Contains(upper, "VERSION 2"):
			return "Apache-2.0"
		case strings.Contains(upper, "GNU GENERAL PUBLIC LICENSE") && strings.Contains(upper, "VERSION 3"):
			return "GPL-3.0"
		case strings.Contains(upper, "GNU GENERAL PUBLIC LICENSE") && strings.Contains(upper, "VERSION 2"):
			return "GPL-2.0"
		case strings.Contains(upper, "BSD 3-CLAUSE") || (strings.Contains(upper, "REDISTRIBUTION AND USE") && strings.Contains(upper, "NEITHER THE NAME")):
			return "BSD-3-Clause"
		case strings.Contains(upper, "BSD 2-CLAUSE") || strings.Contains(upper, "SIMPLIFIED BSD"):
			return "BSD-2-Clause"
		case strings.Contains(upper, "MOZILLA PUBLIC LICENSE"):
			return "MPL-2.0"
		case strings.Contains(upper, "ISC LICENSE"):
			return "ISC"
		case strings.Contains(upper, "UNLICENSE"):
			return "Unlicense"
		default:
			for _, line := range strings.Split(text, "\n") {
				line = strings.TrimSpace(line)
				if line != "" {
					if len(line) > 40 {
						line = line[:40] + "…"
					}
					return line
				}
			}
		}
	}
	return "N/A"
}

func dirSize(root string) (int64, error) {
	var size int64
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		size += info.Size()
		return nil
	})
	return size, err
}
