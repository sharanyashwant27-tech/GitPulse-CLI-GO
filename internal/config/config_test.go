package config

import "testing"

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.Theme != "default" {
		t.Fatalf("default theme = %s", cfg.Theme)
	}
	if cfg.Limit != 20 {
		t.Fatalf("default limit = %d", cfg.Limit)
	}
}

func TestValidateTheme(t *testing.T) {
	if err := ValidateTheme("nord"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateTheme("solarized"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateTheme("tokyo night"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateTheme("nope"); err == nil {
		t.Fatal("expected error for unknown theme")
	}
}

func TestValidThemesCount(t *testing.T) {
	if len(ValidThemes()) != 7 {
		t.Fatalf("expected 7 themes, got %d", len(ValidThemes()))
	}
}
