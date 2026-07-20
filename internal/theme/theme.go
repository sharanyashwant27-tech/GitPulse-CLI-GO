package theme

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Theme defines a color palette for the TUI and CLI output.
type Theme struct {
	Name       string
	Background lipgloss.Color
	Foreground lipgloss.Color
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Error      lipgloss.Color
	Muted      lipgloss.Color
	Border     lipgloss.Color
	Highlight  lipgloss.Color
}

var themes = map[string]Theme{
	"default": {
		Name:       "Default",
		Background: "#0F1419",
		Foreground: "#E6EDF3",
		Primary:    "#58A6FF",
		Secondary:  "#79C0FF",
		Accent:     "#FFA657",
		Success:    "#3FB950",
		Warning:    "#D29922",
		Error:      "#F85149",
		Muted:      "#8B949E",
		Border:     "#30363D",
		Highlight:  "#A5D6FF",
	},
	"dracula": {
		Name:       "Dracula",
		Background: "#282A36",
		Foreground: "#F8F8F2",
		Primary:    "#BD93F9",
		Secondary:  "#8BE9FD",
		Accent:     "#FF79C6",
		Success:    "#50FA7B",
		Warning:    "#FFB86C",
		Error:      "#FF5555",
		Muted:      "#6272A4",
		Border:     "#44475A",
		Highlight:  "#F1FA8C",
	},
	"nord": {
		Name:       "Nord",
		Background: "#2E3440",
		Foreground: "#ECEFF4",
		Primary:    "#88C0D0",
		Secondary:  "#81A1C1",
		Accent:     "#B48EAD",
		Success:    "#A3BE8C",
		Warning:    "#EBCB8B",
		Error:      "#BF616A",
		Muted:      "#4C566A",
		Border:     "#3B4252",
		Highlight:  "#EBCB8B",
	},
	"catppuccin": {
		Name:       "Catppuccin",
		Background: "#1E1E2E",
		Foreground: "#CDD6F4",
		Primary:    "#CBA6F7",
		Secondary:  "#89B4FA",
		Accent:     "#F5C2E7",
		Success:    "#A6E3A1",
		Warning:    "#F9E2AF",
		Error:      "#F38BA8",
		Muted:      "#6C7086",
		Border:     "#313244",
		Highlight:  "#94E2D5",
	},
	"tokyo-night": {
		Name:       "Tokyo Night",
		Background: "#1A1B26",
		Foreground: "#C0CAF5",
		Primary:    "#7AA2F7",
		Secondary:  "#BB9AF7",
		Accent:     "#FF9E64",
		Success:    "#9ECE6A",
		Warning:    "#E0AF68",
		Error:      "#F7768E",
		Muted:      "#565F89",
		Border:     "#24283B",
		Highlight:  "#7DCFFF",
	},
	"gruvbox": {
		Name:       "Gruvbox",
		Background: "#282828",
		Foreground: "#EBDBB2",
		Primary:    "#FABD2F",
		Secondary:  "#83A598",
		Accent:     "#D3869B",
		Success:    "#B8BB26",
		Warning:    "#FE8019",
		Error:      "#FB4934",
		Muted:      "#928374",
		Border:     "#3C3836",
		Highlight:  "#8EC07C",
	},
	"solarized": {
		Name:       "Solarized",
		Background: "#002B36",
		Foreground: "#839496",
		Primary:    "#268BD2",
		Secondary:  "#2AA198",
		Accent:     "#CB4B16",
		Success:    "#859900",
		Warning:    "#B58900",
		Error:      "#DC322F",
		Muted:      "#586E75",
		Border:     "#073642",
		Highlight:  "#EEE8D5",
	},
}

// Get returns a theme by name, defaulting to Default.
func Get(name string) Theme {
	key := normalize(name)
	if t, ok := themes[key]; ok {
		return t
	}
	return themes["default"]
}

// Names returns theme identifiers in display order.
func Names() []string {
	return []string{
		"default",
		"dracula",
		"nord",
		"catppuccin",
		"tokyo-night",
		"gruvbox",
		"solarized",
	}
}

// DisplayNames returns human-readable theme names in display order.
func DisplayNames() []string {
	out := make([]string, 0, len(Names()))
	for _, id := range Names() {
		out = append(out, themes[id].Name)
	}
	return out
}

// TerminalThemes renders the spaced theme list card.
//
//	Terminal Themes
//	Default
//
//	Dracula
func TerminalThemes() string {
	var b strings.Builder
	b.WriteString("Terminal Themes")
	for _, name := range DisplayNames() {
		b.WriteString("\n\n")
		b.WriteString(name)
	}
	return b.String()
}

func normalize(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.Join(strings.Fields(name), "-")
	switch name {
	case "tokyo", "tokyonight":
		return "tokyo-night"
	case "":
		return "default"
	default:
		return name
	}
}

// Styles provides pre-built lipgloss styles for a theme.
type Styles struct {
	Theme     Theme
	Title     lipgloss.Style
	Subtitle  lipgloss.Style
	Panel     lipgloss.Style
	Label     lipgloss.Style
	Value     lipgloss.Style
	Success   lipgloss.Style
	Warning   lipgloss.Style
	Error     lipgloss.Style
	Muted     lipgloss.Style
	Highlight lipgloss.Style
	TableHead lipgloss.Style
	Bar       lipgloss.Style
}

// NewStyles builds lipgloss styles for the given theme.
func NewStyles(t Theme) Styles {
	return Styles{
		Theme: t,
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Primary).
			MarginBottom(1),
		Subtitle: lipgloss.NewStyle().
			Foreground(t.Secondary).
			Italic(true),
		Panel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Border).
			Padding(0, 1).
			MarginRight(1).
			MarginBottom(1),
		Label: lipgloss.NewStyle().
			Foreground(t.Muted).
			Width(16),
		Value: lipgloss.NewStyle().
			Foreground(t.Foreground).
			Bold(true),
		Success: lipgloss.NewStyle().Foreground(t.Success),
		Warning: lipgloss.NewStyle().Foreground(t.Warning),
		Error:   lipgloss.NewStyle().Foreground(t.Error),
		Muted:   lipgloss.NewStyle().Foreground(t.Muted),
		Highlight: lipgloss.NewStyle().
			Foreground(t.Highlight).
			Bold(true),
		TableHead: lipgloss.NewStyle().
			Foreground(t.Accent).
			Bold(true),
		Bar: lipgloss.NewStyle().Foreground(t.Primary),
	}
}

// ProgressBar renders a colored ASCII progress bar.
func (s Styles) ProgressBar(pct float64, width int) string {
	if width <= 0 {
		width = 20
	}
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	filled := int(pct / 100 * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return s.Bar.Render(bar)
}
