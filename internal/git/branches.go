package git

import (
	"sort"

	"github.com/go-git/go-git/v5/plumbing"
)

func (a *Analyzer) branches() ([]BranchInfo, error) {
	head, _ := a.repo.Head()
	current := ""
	if head != nil {
		current = head.Name().Short()
	}

	iter, err := a.repo.References()
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var branches []BranchInfo
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name()
		if !name.IsBranch() && !name.IsRemote() {
			return nil
		}
		info := BranchInfo{
			Name:      name.Short(),
			IsRemote:  name.IsRemote(),
			IsCurrent: name.IsBranch() && name.Short() == current,
		}
		commit, err := a.repo.CommitObject(ref.Hash())
		if err == nil {
			info.LastCommit = commit.Author.When
			info.Author = commit.Author.Name
		}
		branches = append(branches, info)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(branches, func(i, j int) bool {
		if branches[i].IsCurrent != branches[j].IsCurrent {
			return branches[i].IsCurrent
		}
		return branches[i].LastCommit.After(branches[j].LastCommit)
	})
	return branches, nil
}
