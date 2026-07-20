package theme

import (
	"strings"
	"testing"
)

func TestGetKnownThemes(t *testing.T) {
	for _, name := range Names() {
		th := Get(name)
		if th.Name == "" {
			t.Fatalf("theme %q has empty display name", name)
		}
		if th.Primary == "" {
			t.Fatalf("theme %q missing primary color", name)
		}
	}
}

func TestGetUnknownFallsBackToDefault(t *testing.T) {
	th := Get("not-a-real-theme")
	if th.Name != "Default" {
		t.Fatalf("expected Default fallback, got %s", th.Name)
	}
}

func TestGetAliases(t *testing.T) {
	if Get("tokyo night").Name != "Tokyo Night" {
		t.Fatalf("got %s", Get("tokyo night").Name)
	}
	if Get("Solarized").Name != "Solarized" {
		t.Fatalf("got %s", Get("Solarized").Name)
	}
}

func TestTerminalThemes(t *testing.T) {
	out := TerminalThemes()
	for _, want := range []string{
		"Terminal Themes",
		"Default",
		"Dracula",
		"Nord",
		"Catppuccin",
		"Tokyo Night",
		"Gruvbox",
		"Solarized",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}

func TestProgressBar(t *testing.T) {
	s := NewStyles(Get("nord"))
	bar := s.ProgressBar(50, 10)
	if bar == "" {
		t.Fatal("expected non-empty progress bar")
	}
}
