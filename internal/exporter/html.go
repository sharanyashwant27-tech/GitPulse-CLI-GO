package exporter

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/git"
)

func (e *Exporter) writeHTML(path string, analysis *git.Analysis) error {
	tmplPath := filepath.Join(e.TemplateDir, "report.html")
	var tmpl *template.Template
	var err error

	if _, statErr := os.Stat(tmplPath); statErr == nil {
		tmpl, err = template.ParseFiles(tmplPath)
	} else {
		tmpl, err = template.New("report").Parse(defaultHTMLTemplate)
	}
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, analysis); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

const defaultHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1"/>
<title>GitPulse — {{.Repository.Name}}</title>
<style>
  :root { --bg:#0f1419; --card:#1a2332; --fg:#e7ecf3; --accent:#7aa2f7; --muted:#8b9bb4; --ok:#9ece6a; }
  body { margin:0; font-family: ui-sans-serif, system-ui, sans-serif; background:var(--bg); color:var(--fg); }
  header { padding:2rem 1.5rem; background:linear-gradient(135deg,#1a2332,#0f1419 60%); border-bottom:1px solid #243044; }
  h1 { margin:0; font-size:1.8rem; color:var(--accent); }
  .sub { color:var(--muted); margin-top:.4rem; }
  main { padding:1.5rem; display:grid; gap:1rem; grid-template-columns:repeat(auto-fit,minmax(280px,1fr)); }
  .card { background:var(--card); border:1px solid #243044; border-radius:12px; padding:1rem 1.2rem; }
  .card h2 { margin:0 0 .8rem; font-size:1rem; color:var(--accent); }
  .metric { font-size:1.6rem; font-weight:700; }
  .muted { color:var(--muted); font-size:.9rem; }
  table { width:100%; border-collapse:collapse; font-size:.9rem; }
  td, th { text-align:left; padding:.35rem 0; border-bottom:1px solid #243044; }
  .ok { color:var(--ok); }
</style>
</head>
<body>
<header>
  <h1>GitPulse · {{.Repository.Name}}</h1>
  <div class="sub">{{.Repository.Path}} · {{.Repository.DefaultBranch}} · generated {{.GeneratedAt.Format "2006-01-02 15:04"}}</div>
</header>
<main>
  <section class="card">
    <h2>Repository Overview</h2>
    <div class="metric">{{.Repository.Name}}</div>
    <div class="muted">
      Branch {{.Repository.DefaultBranch}} · Status {{.Repository.Status}} ·
      {{.CommitStats.TotalCommits}} commits · {{.Repository.BranchCount}} branches ·
      {{.Repository.TagCount}} tags · {{.CommitStats.TotalAuthors}} contributors
    </div>
  </section>
  <section class="card">
    <h2>Health</h2>
    <div class="metric ok">{{printf "%.1f" .Health.Score}} <span class="muted">({{.Health.Grade}})</span></div>
    <div class="muted">Streak {{.Streak.Current}}d · best {{.Streak.Longest}}d</div>
  </section>
  <section class="card">
    <h2>Top Contributors</h2>
    <table>
      <tr><th>Name</th><th>Commits</th><th>%</th></tr>
      {{range $i, $c := .Contributors}}{{if lt $i 10}}
      <tr><td>{{$c.Name}}</td><td>{{$c.Commits}}</td><td>{{printf "%.1f" $c.Percentage}}%</td></tr>
      {{end}}{{end}}
    </table>
  </section>
  <section class="card">
    <h2>Languages</h2>
    <table>
      {{range $i, $l := .Languages}}{{if lt $i 8}}
      <tr><td>{{$l.Name}}</td><td>{{printf "%.1f" $l.Percentage}}%</td><td>{{$l.Files}} files</td></tr>
      {{end}}{{end}}
    </table>
  </section>
  <section class="card" style="grid-column:1/-1">
    <h2>Recent Commits</h2>
    <table>
      <tr><th>Hash</th><th>Author</th><th>Message</th><th>Date</th></tr>
      {{range .RecentCommits}}
      <tr><td><code>{{.ShortHash}}</code></td><td>{{.Author}}</td><td>{{.Message}}</td><td>{{.Date.Format "2006-01-02"}}</td></tr>
      {{end}}
    </table>
  </section>
</main>
</body>
</html>`
