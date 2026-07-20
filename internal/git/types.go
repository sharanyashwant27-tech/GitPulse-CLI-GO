package git

import "time"

// Repository holds high-level repository metadata.
type Repository struct {
	Path            string    `json:"path"`
	Name            string    `json:"name"`
	RemoteURL       string    `json:"remote_url,omitempty"`
	DefaultBranch   string    `json:"default_branch"`
	Status          string    `json:"status"` // Clean or Dirty
	IsBare          bool      `json:"is_bare"`
	SizeBytes       int64     `json:"size_bytes"`
	BranchCount     int       `json:"branch_count"`
	TagCount        int       `json:"tag_count"`
	PrimaryLanguage string    `json:"primary_language,omitempty"`
	License         string    `json:"license,omitempty"`
	Stars           string    `json:"stars"` // local repos are typically "N/A"
	CreatedAt       time.Time `json:"created_at,omitempty"`
	LastCommitAt    time.Time `json:"last_commit_at,omitempty"`
}

// CommitStats aggregates commit-level statistics.
type CommitStats struct {
	TotalCommits     int            `json:"total_commits"`
	TotalAuthors     int            `json:"total_authors"`
	FirstCommit      time.Time      `json:"first_commit"`
	LastCommit       time.Time      `json:"last_commit"`
	AvgCommitsPerDay float64        `json:"avg_commits_per_day"`
	CommitsByDay     map[string]int `json:"commits_by_day"`
	CommitsByWeekday [7]int         `json:"commits_by_weekday"`
	CommitsByHour    [24]int        `json:"commits_by_hour"`
	Additions        int            `json:"additions"`
	Deletions        int            `json:"deletions"`
	FilesChanged     int            `json:"files_changed"`
	Today            int            `json:"today"`
	ThisWeek         int            `json:"this_week"`
	ThisMonth        int            `json:"this_month"`
}

// Contributor represents a single author contribution profile.
type Contributor struct {
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Commits    int       `json:"commits"`
	Additions  int       `json:"additions"`
	Deletions  int       `json:"deletions"`
	FirstSeen  time.Time `json:"first_seen"`
	LastSeen   time.Time `json:"last_seen"`
	Percentage float64   `json:"percentage"`
}

// BranchInfo describes a local or remote branch.
type BranchInfo struct {
	Name       string    `json:"name"`
	IsRemote   bool      `json:"is_remote"`
	IsCurrent  bool      `json:"is_current"`
	IsMerged   bool      `json:"is_merged"`
	LastCommit time.Time `json:"last_commit"`
	Author     string    `json:"author"`
	Ahead      int       `json:"ahead"`
	Behind     int       `json:"behind"`
}

// LanguageStat represents language composition by estimated bytes.
type LanguageStat struct {
	Name       string  `json:"name"`
	Bytes      int64   `json:"bytes"`
	Files      int     `json:"files"`
	Percentage float64 `json:"percentage"`
}

// TimelinePoint is a single day on the commit activity timeline.
type TimelinePoint struct {
	Date    string `json:"date"`
	Commits int    `json:"commits"`
}

// MonthlyPoint is commit activity for one calendar month.
type MonthlyPoint struct {
	Month     string `json:"month"`      // Jan, Feb, ...
	YearMonth string `json:"year_month"` // 2006-01
	Commits   int    `json:"commits"`
}

// AuthorTimeline is one author's monthly commit activity.
type AuthorTimeline struct {
	Name    string         `json:"name"`
	Email   string         `json:"email"`
	Monthly []MonthlyPoint `json:"monthly"`
}

// CommitTypeStat is a classified commit-message category share.
type CommitTypeStat struct {
	Name       string  `json:"name"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// HeatmapCell is one cell in the weekly contribution heatmap.
type HeatmapCell struct {
	Date    string `json:"date"`
	Weekday int    `json:"weekday"`
	Week    int    `json:"week"`
	Count   int    `json:"count"`
}

// RecentCommit is a compact commit summary for dashboards.
type RecentCommit struct {
	Hash      string    `json:"hash"`
	ShortHash string    `json:"short_hash"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	Email     string    `json:"email"`
	Date      time.Time `json:"date"`
	Files     int       `json:"files"`
	Additions int       `json:"additions"`
	Deletions int       `json:"deletions"`
}

// FileChangeStat tracks churn for a single path.
type FileChangeStat struct {
	Path      string `json:"path"`
	Changes   int    `json:"changes"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

// FileChangeSummary aggregates Added/Modified/Deleted/Renamed file events.
type FileChangeSummary struct {
	Added    int `json:"added"`
	Modified int `json:"modified"`
	Deleted  int `json:"deleted"`
	Renamed  int `json:"renamed"`
}

// HealthScore captures repository health metrics and overall score.
type HealthScore struct {
	Score                float64           `json:"score"`
	Grade                string            `json:"grade"`
	CommitFrequency      float64           `json:"commit_frequency"`
	ContributorDiversity float64           `json:"contributor_diversity"`
	BranchHygiene        float64           `json:"branch_hygiene"`
	CodeChurn            float64           `json:"code_churn"`
	ActivityRecency      float64           `json:"activity_recency"`
	Details              map[string]string `json:"details"`
}

// ProductivityFactor is one pillar of the productivity / code-quality score.
type ProductivityFactor struct {
	Name   string  `json:"name"`
	Passed bool    `json:"passed"`
	Score  float64 `json:"score"`
}

// ProductivityScore is the Code Quality productivity rating.
type ProductivityScore struct {
	Label   string               `json:"label"` // e.g. "Code Quality"
	Score   float64              `json:"score"`
	Factors []ProductivityFactor `json:"factors"`
}

// Streak tracks consecutive commit days.
type Streak struct {
	Current    int       `json:"current"`
	Longest    int       `json:"longest"`
	LastActive time.Time `json:"last_active"`
}

// Analysis is the full repository analysis payload.
type Analysis struct {
	Repository        Repository        `json:"repository"`
	CommitStats         CommitStats       `json:"commit_stats"`
	Contributors        []Contributor     `json:"contributors"`
	Branches            []BranchInfo      `json:"branches"`
	Languages           []LanguageStat    `json:"languages"`
	Timeline            []TimelinePoint   `json:"timeline"`
	Monthly             []MonthlyPoint    `json:"monthly"`
	AuthorTimelines     []AuthorTimeline  `json:"author_timelines"`
	CommitTypes         []CommitTypeStat  `json:"commit_types"`
	Heatmap             []HeatmapCell     `json:"heatmap"`
	RecentCommits       []RecentCommit    `json:"recent_commits"`
	FileChanges         []FileChangeStat  `json:"file_changes"`
	FileChangeSummary   FileChangeSummary `json:"file_change_summary"`
	Health              HealthScore       `json:"health"`
	Productivity        ProductivityScore `json:"productivity"`
	Streak              Streak            `json:"streak"`
	GeneratedAt         time.Time         `json:"generated_at"`
}
