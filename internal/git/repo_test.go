package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestOpenMissingRepo(t *testing.T) {
	dir := t.TempDir()
	_, err := Open(Options{Path: dir})
	if err == nil {
		t.Fatal("expected error for non-git directory")
	}
}

func TestAnalyzeTempRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test User",
			"GIT_AUTHOR_EMAIL=test@example.com",
			"GIT_COMMITTER_NAME=Test User",
			"GIT_COMMITTER_EMAIL=test@example.com",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}

	run("init")
	run("config", "user.email", "test@example.com")
	run("config", "user.name", "Test User")

	file := filepath.Join(dir, "main.go")
	if err := os.WriteFile(file, []byte("package main\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "initial commit")

	if err := os.WriteFile(file, []byte("package main\n\nfunc main() {\n\tprintln(\"hi\")\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "second commit")

	az, err := Open(Options{Path: dir, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	a, err := az.Analyze()
	if err != nil {
		t.Fatal(err)
	}
	if a.CommitStats.TotalCommits < 2 {
		t.Fatalf("expected >=2 commits, got %d", a.CommitStats.TotalCommits)
	}
	if len(a.Contributors) == 0 {
		t.Fatal("expected contributors")
	}
	if len(a.RecentCommits) == 0 {
		t.Fatal("expected recent commits")
	}
	if a.Repository.Name == "" {
		t.Fatal("expected repo name")
	}
	foundGo := false
	for _, l := range a.Languages {
		if l.Name == "Go" {
			foundGo = true
			break
		}
	}
	if !foundGo {
		t.Fatal("expected Go in language breakdown")
	}
	if a.GeneratedAt.After(time.Now().Add(time.Minute)) {
		t.Fatal("generated_at in the future")
	}
}

func TestFormatBytes(t *testing.T) {
	if FormatBytes(500) != "500 B" {
		t.Fatalf("unexpected: %s", FormatBytes(500))
	}
	if FormatBytes(2048) == "2048 B" {
		t.Fatal("expected unit scaling")
	}
}

