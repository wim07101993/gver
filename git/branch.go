package git

import (
	"github.com/go-git/go-git/v6"
	"github.com/pkg/errors"
	"log/slog"
)

func BranchName(repo *git.Repository) (string, error) {
	head, err := repo.Head()
	if err != nil {
		return "", errors.Wrap(err, "failed getting the head of the repo")
	}
	branch := head.Name().Short()
	slog.Debug("BranchName: Found branch name", slog.String("branch", branch))
	return branch, nil
}
