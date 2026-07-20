# GitPulse

**Interactive Git repository analytics for your terminal.**

GitPulse is a production-grade Go CLI that inspects local Git repositories with [go-git](https://github.com/go-git/go-git), renders colorful dashboards with [Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Lip Gloss](https://github.com/charmbracelet/lipgloss), and exports reports to HTML, JSON, CSV, and PDF.

```
┌──────────────────────────────────────────┐
│                GIT PULSE                 │
├──────────────────────────────────────────┤
│                                          │
│  Repository                              │
│  ├── my-project                          │
│  ├── Size: 240 MB                        │
│  ├── Language: Go                        │
│  ├── Stars: N/A                          │
│  └── License: MIT                        │
│                                          │
│  Activity                                │
│  ├── Today's commits : 8                 │
│  ├── This week       : 39                │
│  └── This month      : 121               │
│                                          │
│  Health                                  │
│  ████████████████░░░░ 82%                │
│                                          │
├──────────────────────────────────────────┤
└──────────────────────────────────────────┘
```

```
┌─ Repository Overview · Repository: my-project ─┐  ┌─ Recent Commits ──────────────────────────┐
│ Branch      : main                             │  │ a1b2c3d  2026-07-18  feat: add dashboard  │
│ Status      : Clean                            │  │ d4e5f6a  2026-07-17  fix: heatmap scale   │
│ Commits     : 1,532                            │  │ 9087b6c  2026-07-16  docs: update README  │
│ Branches    : 18                               │  │                                           │
│ Tags        : 42                               │  │                                           │
│ Contributors: 11                               │  │                                           │
│ Created     : Jan 2022                         │  │                                           │
│ Last Commit : 2 hours ago                      │  │                                           │
│ ██████████████████████████                     │  │                                           │
│ Health Score: 94%                              │  │                                           │
└────────────────────────────────────────────────┘  └───────────────────────────────────────────┘
```

---

## Features

| Area | What you get |
|------|----------------|
| Repository overview | Name, branch, Clean/Dirty status, commits, branches, tags, contributors, created, last commit, health bar |
| Beautiful dashboard | Centered GIT PULSE card: size, language, stars, license, today/week/month activity, health bar |
| Commit statistics | Totals, avg/day, additions/deletions, weekday distribution |
| Contributor leaderboard | Top Contributors table (`Name` / `Commits`) |
| Author timeline | Per-author monthly bars (`John` → `Jan ███`) |
| Branch information | Branch Viewer list of local branches (current first) |
| Language breakdown | Language Statistics bars with top languages + Others |
| Commit activity timeline | 90-day ASCII line chart ([asciigraph](https://github.com/guptarohit/asciigraph)) |
| Commit activity graph | Year-to-date monthly bars (`Jan ███` … `Jul ██████████`) |
| Commit types | Conventional/keyword classification: Features, Fixes, Refactor, Docs, Tests, Others |
| Weekly heatmap | Mon→Sun intensity bars (`Mon ██` … `Sun █`) plus contribution calendar |
| Recent commits | Checkmarked subject list (`✓ Fix login timeout`) |
| File change statistics | Added / Modified / Deleted / Renamed event counts |
| Repository health score | Weighted grade (A+–F) with factor breakdown |
| Productivity score | Code Quality rating from frequency, branches, docs, tests, health |
| Git streak | Current streak bar card (`16 Days`) |
| Themes | Default, Dracula, Nord, Catppuccin, Tokyo Night, Gruvbox, Solarized |
| Live Mode | `gitpulse watch` — auto-refreshing Bubble Tea dashboard |
| Export | HTML · JSON · CSV · PDF |

---

## Advanced Features

| Feature | Status |
|---------|--------|
| Repository health scoring | Shipped (`gitpulse health`, dashboard Health tab) |
| License detection | Shipped (LICENSE heuristics in `internal/git`) |
| GitHub integration (stars, forks, pull requests, issues) | Planned |
| GitLab integration | Planned |
| Bitbucket support | Planned |
| Commit sentiment analysis using local AI models | Planned |
| Code ownership insights | Planned |
| Release timeline | Planned |
| Changelog generation | Planned |
| Git hooks manager | Planned |
| Multi-repository dashboard | Planned |
| CI/CD pipeline status | Planned |
| Security scan summary | Planned |
| Dependency analysis | Planned |
| Code churn analysis | Planned |
| Large file detection | Planned |
| Interactive fuzzy search | Planned |
| Keyboard shortcuts | Planned |
| Auto-update checker | Planned |

---

## Requirements

- Go **1.25+**
- A local Git repository (analysis uses go-git; no `git` subprocess for stats)
- Optional: Docker (docs site on port **8098**)

---

## Installation

### From source

```bash
git clone https://github.com/sharanyashwant27-tech/GitPulse-CLI-GO.git
cd GitPulse-CLI-GO
go install .
```

Or build a local binary:

```bash
make build
./bin/gitpulse --help
```

### With Go

```bash
go install github.com/sharanyashwant27-tech/GitPulse-CLI-GO@latest
```

### Docker (CLI)

```bash
docker build --target runtime -t gitpulse:cli .
docker run --rm -v "$PWD":/repo -w /repo gitpulse:cli stats -r /repo
```

---

## Quick start

```bash
# Interactive dashboard (default)
gitpulse

# Or point at any repo
gitpulse --repo /path/to/repo --theme nord

# One-shot reports
gitpulse stats
gitpulse commits
gitpulse branches
gitpulse contributors
gitpulse timeline
gitpulse graph
gitpulse health

# Live Mode
gitpulse watch

# Export
gitpulse export html
gitpulse export pdf
```

---

## CLI Commands

```
CLI Commands
gitpulse

gitpulse stats

gitpulse commits

gitpulse branches

gitpulse contributors

gitpulse timeline

gitpulse graph

gitpulse health

gitpulse export html

gitpulse export pdf

gitpulse watch
```

| Command | Description |
|---------|-------------|
| `gitpulse` | Open the interactive dashboard |
| `gitpulse stats` | Repository overview + commit / file stats |
| `gitpulse commits` | Recent commits, types, and file change stats |
| `gitpulse branches` | Branch Viewer list |
| `gitpulse contributors` | Top Contributors + Author Timeline |
| `gitpulse timeline` | Monthly activity + heatmap |
| `gitpulse graph` | Full ASCII chart suite |
| `gitpulse health` | Productivity, health, and streak |
| `gitpulse export html` | Export HTML report |
| `gitpulse export pdf` | Export PDF report |
| `gitpulse export csv` | Export CSV report |
| `gitpulse export json` | Export JSON report |
| `gitpulse watch` | Live Mode — auto-refreshing dashboard |
| `gitpulse themes` | Terminal Themes list |
| `gitpulse commands` | Print this CLI Commands card |

### Global flags

```
--repo, -r          Path to git repository (default: cwd)
--theme, -t         default | dracula | nord | catppuccin | tokyo-night | gruvbox | solarized
--limit, -n         List limit (default: 20)
--since / --until   Filter commits (YYYY-MM-DD)
--interval          Watch refresh interval (default: 5s)
--export-dir        Export output directory (default: reports)
--log-level         debug | info | warn | error
--config            Path to YAML config
--verbose, -v       Debug logging
```

Environment variables use the `GITPULSE_` prefix (e.g. `GITPULSE_THEME=nord`).

---

## Configuration

Create `gitpulse.yaml` in the project root or `~/.config/gitpulse/gitpulse.yaml`:

```yaml
repo: .
theme: default
log_level: info
watch_interval: 5s
export_dir: reports
limit: 20
```

---

## Themes

```
Terminal Themes
Default

Dracula

Nord

Catppuccin

Tokyo Night

Gruvbox

Solarized
```

```bash
gitpulse themes
gitpulse stats --theme catppuccin
gitpulse --theme "tokyo night"
```

---

## Screenshots (ASCII)

### Top Contributors

```
Top Contributors
Name                 Commits

John Doe             421
Alice                301
David                192
Sarah                166
```

### Author Timeline

```
Author Timeline
John

Jan ███

Feb ██████

Mar █████████

Apr ███████
```

### Commit Types

```
Commit Types
Features     ████████████ 41%

Fixes        ████████ 29%

Refactor     █████ 14%

Docs         ██ 7%

Tests        █ 4%

Others       █ 5%
```

### Branch Viewer

```
Branch Viewer
main

feature/login

feature/payment

release/v2

hotfix/session
```

### Recent Commits

```
Recent Commits
✓ Fix login timeout

✓ Add payment API

✓ Improve dashboard

✓ Refactor services

✓ Update README
```

### File Change Statistics

```
File Change Statistics
Added      : 312

Modified   : 121

Deleted    : 38

Renamed    : 14
```

### Language Statistics

```
Language Statistics
Go      ███████████ 72%

HTML    ███ 14%

CSS     ██ 8%

YAML    █ 4%

Others  █ 2%
```

### Productivity Score

```
Productivity Score
Code Quality

██████████████████░░

87%

Based on

✓ Commit frequency

✓ Branch usage

✓ Documentation

✓ Tests

✓ Repository health
```

### Git Streak

```
Git Streak
Current Streak

████████████

16 Days
```

### Live Mode

```
Live Mode
gitpulse watch
```

Opens the interactive dashboard and re-analyzes the repository on an interval
(default `5s`). A green **LIVE** badge appears in the header while watching.

```bash
gitpulse watch
gitpulse watch --interval 3s
gitpulse watch --theme nord
```

### Health score

```
GitPulse · Health
Score: 86.5 / 100   Grade: A

Frequency      ████████████████████░░░░░░░░  78
Diversity      ████████████████████████░░░░  85
Branches       ██████████████████████████░░  92
Churn          ████████████████████████░░░░  88
Recency        ████████████████████████████ 100

Git streak
Current        4 days
Longest        21 days
```

### Commit Activity Graph

```
Commit Activity Graph
Jan ███
Feb ████████
Mar ███████████
Apr ██████████████
May ███████
Jun █████████████████
Jul ██████████
```

### Timeline graph

```
 3.00 ┤  ╭╮
 2.00 ┤╭─╯╰╮  ╭─╮
 1.00 ┤╯   ╰──╯ ╰─
 0.00 ┼────────────
      Commit activity
```

### Weekly heatmap

```
Weekly Heatmap
Mon ██
Tue █████
Wed ████████
Thu ██████████
Fri ██████
Sat ██
Sun █
```

### Contribution calendar

```
Weekly contribution heatmap
Su ·····░··········
Mo ··▒█▓░··········
Tu ·░▓█▒···········
We ···▒░···········
Th ··░▒▓█··········
Fr ·····░··········
Sa ················
   less · ░ ▒ ▓ █ more
```

---

## Export formats

```
Export Reports
gitpulse export html

gitpulse export pdf

gitpulse export csv

gitpulse export json
```

```bash
gitpulse export html   # reports/gitpulse-<repo>-<ts>.html
gitpulse export json
gitpulse export csv
gitpulse export pdf
```

HTML uses `templates/report.html`. JSON mirrors the full `Analysis` model for downstream tooling.

---

## Architecture

### Suggested Project Structure

```
main.go                    # package main, calls cmd.Execute()
cmd/
  root.go                  # package cmd — Execute, root cobra, flags
  shared.go                # Shared + AnalyzeRepo
  stats.go
  commits.go
  branches.go
  contributors.go
  timeline.go
  graph.go
  export.go
  health.go
  watch.go
  themes.go
  commands.go              # CLICommands list helper + commands subcommand
internal/
  git/
    types.go               # domain types + Analysis
    repo.go                # Open, Analyze, repositoryInfo, refs, license, dirSize
    commits.go             # loadCommits, commitStats, recent, timeline, monthly, authorTimelines, streak, classifyCommitTypes
    branches.go            # branches()
    contributors.go        # contributors()
    stats.go               # languages, fileChanges, fileChangeSummary, FormatBytes
    health.go              # Calculate + CalculateProductivity
  dashboard/
    renderer.go            # Bubble Tea dashboard
    widgets.go             # overview cards + terminal widgets
    charts.go              # ASCII line/bar/heatmap/monthly/weekly/author timeline
  exporter/
    export.go              # Exporter type + Export dispatcher
    html.go
    json.go
    csv.go
    pdf.go
  theme/                   # Lip Gloss palettes
  utils/
    logger.go
    helpers.go             # FindTemplatesDir
  config/
pkg/gitpulse/              # Public Go API (re-exports internal/git)
templates/
assets/
reports/
```

Clean layering: `cmd` → `internal/git` + `internal/dashboard` + `internal/exporter`; TUI and exporters consume `git.Analysis`.

### Public API

```go
import "github.com/sharanyashwant27-tech/GitPulse-CLI-GO/pkg/gitpulse"

analysis, err := gitpulse.Analyze(gitpulse.Options{Path: ".", Limit: 20})
```

---

## Recommended Go Libraries

| Purpose        | Library     |
| -------------- | ----------- |
| CLI            | [Cobra](https://github.com/spf13/cobra) |
| Terminal UI    | [Bubble Tea](https://github.com/charmbracelet/bubbletea) |
| Styling        | [Lip Gloss](https://github.com/charmbracelet/lipgloss) |
| Tables         | [Bubbles](https://github.com/charmbracelet/bubbles) |
| Git operations | [go-git](https://github.com/go-git/go-git) |
| Charts         | [asciigraph](https://github.com/guptarohit/asciigraph) |
| Progress bars  | [progressbar](https://github.com/schollz/progressbar) |
| Colors         | [fatih/color](https://github.com/fatih/color) |
| Config         | [Viper](https://github.com/spf13/viper) |
| Logging        | [Zap](https://github.com/uber-go/zap) |

GitPulse uses these libraries throughout `cmd/`, `internal/git`, `internal/dashboard`, `internal/exporter`, `internal/config`, and `internal/utils`.

---

## Development

```bash
make deps      # go mod download / tidy
make test      # unit tests
make build     # bin/gitpulse
make lint      # go vet
make run       # interactive dashboard
```

### Tests & CI

- Unit tests cover git analysis, health, dashboard widgets/charts, themes, config, and exporter
- GitHub Actions (`.github/workflows/ci.yml`): vet, race tests, build, smoke commands, golangci-lint, Docker image build

---

## Docs site on localhost:8098

Serve the updated README via Docker:

```bash
docker compose up --build -d
# open http://localhost:8098
```

This builds the `docs` stage of the Dockerfile (`gitpulse:docs`) — an nginx image that hosts the rendered README and static assets on port **8098**.

```bash
# Rebuild after README changes
docker compose up --build -d gitpulse-docs

# Stop
docker compose down
```

---

## License

MIT
