package git

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// Calculate computes a repository health score from analysis inputs.
func Calculate(
	stats CommitStats,
	contributors []Contributor,
	branches []BranchInfo,
	streak Streak,
) HealthScore {
	details := map[string]string{}

	freq := scoreCommitFrequency(stats.AvgCommitsPerDay)
	details["commit_frequency"] = fmt.Sprintf("%.2f commits/day → %.0f/100", stats.AvgCommitsPerDay, freq)

	diversity := scoreContributorDiversity(len(contributors), stats.TotalCommits)
	details["contributor_diversity"] = fmt.Sprintf("%d contributors → %.0f/100", len(contributors), diversity)

	hygiene := scoreBranchHygiene(branches)
	details["branch_hygiene"] = fmt.Sprintf("%d branches → %.0f/100", len(branches), hygiene)

	churn := scoreCodeChurn(stats.Additions, stats.Deletions, stats.TotalCommits)
	details["code_churn"] = fmt.Sprintf("+%d/-%d → %.0f/100", stats.Additions, stats.Deletions, churn)

	recency := scoreActivityRecency(stats.LastCommit)
	details["activity_recency"] = fmt.Sprintf("last commit %s → %.0f/100", formatAge(stats.LastCommit), recency)

	score := freq*0.25 + diversity*0.20 + hygiene*0.15 + churn*0.15 + recency*0.25
	if streak.Current > 0 {
		bonus := math.Min(10, float64(streak.Current)*0.5)
		score = math.Min(100, score+bonus)
		details["streak_bonus"] = fmt.Sprintf("+%.1f for %d-day streak", bonus, streak.Current)
	}

	return HealthScore{
		Score:                round1(score),
		Grade:                grade(score),
		CommitFrequency:      round1(freq),
		ContributorDiversity: round1(diversity),
		BranchHygiene:        round1(hygiene),
		CodeChurn:            round1(churn),
		ActivityRecency:      round1(recency),
		Details:              details,
	}
}

// CalculateProductivity computes the Code Quality productivity score.
func CalculateProductivity(
	stats CommitStats,
	branches []BranchInfo,
	commitTypes []CommitTypeStat,
	languages []LanguageStat,
	repo Repository,
	healthScore HealthScore,
) ProductivityScore {
	freq := healthScore.CommitFrequency
	if freq == 0 {
		freq = scoreCommitFrequency(stats.AvgCommitsPerDay)
	}

	branchScore := scoreBranchUsage(branches)
	docsScore := scoreDocumentation(commitTypes, languages, repo)
	testsScore := scoreTests(commitTypes, languages)
	repoHealth := healthScore.Score

	factors := []ProductivityFactor{
		{Name: "Commit frequency", Score: round1(freq), Passed: freq >= 40},
		{Name: "Branch usage", Score: round1(branchScore), Passed: branchScore >= 40},
		{Name: "Documentation", Score: round1(docsScore), Passed: docsScore >= 40},
		{Name: "Tests", Score: round1(testsScore), Passed: testsScore >= 40},
		{Name: "Repository health", Score: round1(repoHealth), Passed: repoHealth >= 40},
	}

	sum := 0.0
	for _, f := range factors {
		sum += f.Score
	}
	score := sum / float64(len(factors))

	return ProductivityScore{
		Label:   "Code Quality",
		Score:   round1(score),
		Factors: factors,
	}
}

func scoreCommitFrequency(avg float64) float64 {
	switch {
	case avg >= 3:
		return 100
	case avg >= 1:
		return 70 + (avg-1)*15
	case avg >= 0.3:
		return 40 + ((avg - 0.3) / 0.7 * 30)
	case avg > 0:
		return avg / 0.3 * 40
	default:
		return 0
	}
}

func scoreContributorDiversity(authors, commits int) float64 {
	if commits == 0 {
		return 0
	}
	ratio := float64(authors) / math.Max(1, math.Sqrt(float64(commits)))
	score := math.Min(100, ratio*40)
	if authors >= 5 {
		score = math.Min(100, score+15)
	} else if authors >= 2 {
		score = math.Min(100, score+8)
	}
	return score
}

func scoreBranchHygiene(branches []BranchInfo) float64 {
	local := 0
	stale := 0
	cutoff := time.Now().AddDate(0, -3, 0)
	for _, b := range branches {
		if b.IsRemote {
			continue
		}
		local++
		if !b.LastCommit.IsZero() && b.LastCommit.Before(cutoff) {
			stale++
		}
	}
	if local == 0 {
		return 70
	}
	freshRatio := 1 - float64(stale)/float64(local)
	score := freshRatio * 100
	if local > 20 {
		score -= math.Min(30, float64(local-20)*2)
	}
	return math.Max(0, math.Min(100, score))
}

func scoreCodeChurn(additions, deletions, commits int) float64 {
	if commits == 0 {
		return 50
	}
	total := additions + deletions
	if total == 0 {
		return 60
	}
	avg := float64(total) / float64(math.Min(float64(commits), 200))
	switch {
	case avg >= 20 && avg <= 400:
		return 90
	case avg < 20:
		return 50 + avg*2
	case avg <= 1000:
		return 90 - ((avg - 400) / 600 * 40)
	default:
		return 40
	}
}

func scoreActivityRecency(last time.Time) float64 {
	if last.IsZero() {
		return 0
	}
	days := time.Since(last).Hours() / 24
	switch {
	case days <= 1:
		return 100
	case days <= 7:
		return 90
	case days <= 30:
		return 70
	case days <= 90:
		return 45
	case days <= 180:
		return 25
	default:
		return 10
	}
}

func scoreBranchUsage(branches []BranchInfo) float64 {
	local := 0
	for _, b := range branches {
		if !b.IsRemote {
			local++
		}
	}
	switch {
	case local >= 3 && local <= 15:
		return 95
	case local == 2 || local == 16 || local == 17:
		return 80
	case local == 1:
		return 55
	case local == 0:
		return 30
	default:
		return math.Max(35, 90-float64(local-15)*3)
	}
}

func scoreDocumentation(types []CommitTypeStat, languages []LanguageStat, repo Repository) float64 {
	score := 20.0
	docsPct := typePercentage(types, "Docs")
	score += math.Min(40, docsPct*4)

	mdPct := 0.0
	for _, l := range languages {
		if strings.EqualFold(l.Name, "Markdown") {
			mdPct = l.Percentage
			break
		}
	}
	score += math.Min(25, mdPct*2)

	if repo.License != "" && !strings.EqualFold(repo.License, "N/A") {
		score += 15
	}
	return math.Min(100, score)
}

func scoreTests(types []CommitTypeStat, languages []LanguageStat) float64 {
	score := 15.0
	testPct := typePercentage(types, "Tests")
	score += math.Min(50, testPct*8)

	for _, l := range languages {
		switch strings.ToLower(l.Name) {
		case "go", "python", "javascript", "typescript", "java", "rust":
			score += 10
			return math.Min(100, score)
		}
	}
	return math.Min(100, score)
}

func typePercentage(types []CommitTypeStat, name string) float64 {
	for _, t := range types {
		if t.Name == name {
			return t.Percentage
		}
	}
	return 0
}

func grade(score float64) string {
	switch {
	case score >= 90:
		return "A+"
	case score >= 85:
		return "A"
	case score >= 80:
		return "A-"
	case score >= 75:
		return "B+"
	case score >= 70:
		return "B"
	case score >= 65:
		return "B-"
	case score >= 60:
		return "C+"
	case score >= 55:
		return "C"
	case score >= 50:
		return "C-"
	case score >= 40:
		return "D"
	default:
		return "F"
	}
}

func formatAge(t time.Time) string {
	if t.IsZero() {
		return "never"
	}
	d := time.Since(t)
	switch {
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 48*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

func round1(v float64) float64 {
	return math.Round(v*10) / 10
}
