package analyser

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

// A Cloner uses the executer to clone the root of a repository into the
// current working directory.
type Cloner interface {
	Clone(context.Context, Executer) error
}

// PullRequestCloner is a Cloner for handling cloning the HeadURL at HeadRef
// and also fetches BaseURL at BaseRef. When using this, the new base ref
// should be set to FETCH_HEAD.
type PullRequestCloner struct {
	HeadURL string
	HeadRef string
	BaseURL string
	BaseRef string
}

var _ Cloner = &PullRequestCloner{}

// Clone implements the Cloner interface.
func (c *PullRequestCloner) Clone(ctx context.Context, exec Executer) error {
	args := []string{"git", "clone", "--depth", "1", "--branch", c.HeadRef, "--single-branch", c.HeadURL, "."}
	out, err := exec.Execute(ctx, args)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("could not execute %v: %q", args, out))
	}

	// This is a PR, fetch base as some tools (apicompat) needs to
	// reference it.
	args = []string{"git", "fetch", "--depth", "1", c.BaseURL, c.BaseRef}
	out, err = exec.Execute(ctx, args)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("could not execute %v: %q", args, out))
	}

	return nil
}

// PushCloner is a Cloner for handling cloning of HeadURL and checking out HeadRef.
type PushCloner struct {
	HeadURL string
	HeadRef string
}

var _ Cloner = &PushCloner{}

// Clone implements the Cloner interface.
func (c *PushCloner) Clone(ctx context.Context, exec Executer) error {
	// clone repo, this cannot be shallow and needs access to all commits
	// therefore cannot be shallow (or if it is, would required a very
	// large depth and --no-single-branch).
	args := []string{"git", "clone", c.HeadURL, "."}
	out, err := exec.Execute(ctx, args)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("could not execute %v: %q", args, out))
	}

	// Checkout sha
	args = []string{"git", "checkout", c.HeadRef}
	out, err = exec.Execute(ctx, args)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("could not execute %v: %q", args, out))
	}

	return nil
}
