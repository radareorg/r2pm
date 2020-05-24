package fetchers

import (
	"context"
	"errors"
)

type Git struct {
	Repo string
	Ref  string
}

func NewGit(repo, ref string) *Git {
	return &Git{
		Repo: repo,
		Ref:  ref,
	}
}

func (g *Git) Fetch(ctx context.Context, dir string) error {
	return errors.New("not implemented")
}
