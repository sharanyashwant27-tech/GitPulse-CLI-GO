package git

import (
	"sort"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func (a *Analyzer) contributors(commits []*object.Commit) []Contributor {
	type agg struct {
		c Contributor
	}
	byEmail := map[string]*agg{}

	for i, commit := range commits {
		key := strings.ToLower(commit.Author.Email)
		a, ok := byEmail[key]
		if !ok {
			a = &agg{c: Contributor{
				Name:      commit.Author.Name,
				Email:     commit.Author.Email,
				FirstSeen: commit.Author.When,
				LastSeen:  commit.Author.When,
			}}
			byEmail[key] = a
		}
		a.c.Commits++
		if commit.Author.When.After(a.c.LastSeen) {
			a.c.LastSeen = commit.Author.When
		}
		if commit.Author.When.Before(a.c.FirstSeen) {
			a.c.FirstSeen = commit.Author.When
		}
		if i < 200 {
			add, del, _ := commitDiffStats(commit)
			a.c.Additions += add
			a.c.Deletions += del
		}
	}

	total := float64(len(commits))
	out := make([]Contributor, 0, len(byEmail))
	for _, a := range byEmail {
		if total > 0 {
			a.c.Percentage = (float64(a.c.Commits) / total) * 100
		}
		out = append(out, a.c)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Commits > out[j].Commits
	})
	return out
}
