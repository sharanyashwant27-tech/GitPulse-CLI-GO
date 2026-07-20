package exporter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/git"
)

func sampleAnalysis() *git.Analysis {
	return &git.Analysis{
		Repository: git.Repository{
			Name:          "demo",
			Path:          "/tmp/demo",
			DefaultBranch: "main",
			RemoteURL:     "https://example.com/demo.git",
		},
		CommitStats: git.CommitStats{
			TotalCommits: 10,
			TotalAuthors: 2,
			Additions:    100,
			Deletions:    20,
		},
		Contributors: []git.Contributor{
			{Name: "Ada", Email: "ada@example.com", Commits: 7, Percentage: 70},
			{Name: "Bob", Email: "bob@example.com", Commits: 3, Percentage: 30},
		},
		Languages: []git.LanguageStat{
			{Name: "Go", Percentage: 80, Files: 5},
		},
		RecentCommits: []git.RecentCommit{
			{ShortHash: "abc1234", Message: "init", Author: "Ada", Date: time.Now()},
		},
		Health: git.HealthScore{Score: 82, Grade: "A-"},
		Streak: git.Streak{Current: 3, Longest: 10},
		GeneratedAt: time.Now().UTC(),
	}
}

func TestExportJSON(t *testing.T) {
	dir := t.TempDir()
	exp := New("", dir)
	path, err := exp.Export(sampleAnalysis(), FormatJSON, "test")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var a git.Analysis
	if err := json.Unmarshal(data, &a); err != nil {
		t.Fatal(err)
	}
	if a.Repository.Name != "demo" {
		t.Fatalf("name = %s", a.Repository.Name)
	}
}

func TestExportCSVHTMLPDF(t *testing.T) {
	dir := t.TempDir()
	exp := New("", dir)
	a := sampleAnalysis()

	for _, format := range []Format{FormatCSV, FormatHTML, FormatPDF} {
		path, err := exp.Export(a, format, "test-"+string(format))
		if err != nil {
			t.Fatalf("%s: %v", format, err)
		}
		st, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if st.Size() == 0 {
			t.Fatalf("%s file empty: %s", format, path)
		}
		if filepath.Ext(path) == "" {
			t.Fatalf("missing extension: %s", path)
		}
	}
}

