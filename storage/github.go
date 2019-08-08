package storage

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/whywaita/grasshopper/file"

	"github.com/google/uuid"

	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const (
	envGitHubRepo  = "GH_GITHUB_REPO"
	envGitHubUser  = "GH_GITHUB_USER"
	envGitHubToken = "GH_GITHUB_TOKEN"
)

type GitHub struct {
	Repo          string
	User          string
	PersonalToken string
}

func NewGitHubClient() GitHub {
	g := GitHub{}

	if os.Getenv(envGitHubRepo) != "" {
		g.Repo = os.Getenv(envGitHubRepo)
	} else if DefaultGitHubRepository != "" {
		g.Repo = DefaultGitHubRepository
	} else {
		log.Fatal("must be set GitHub Repository")
	}

	if os.Getenv(envGitHubUser) != "" {
		g.User = os.Getenv(envGitHubUser)
	} else if DefaultGitHubUser != "" {
		g.User = DefaultGitHubUser
	} else {
		log.Println("not set GitHub Username, you can't push Private Repository.")
	}

	if os.Getenv(envGitHubToken) != "" {
		g.PersonalToken = os.Getenv(envGitHubToken)
	} else if DefaultGitHubToken != "" {
		g.PersonalToken = DefaultGitHubToken
	} else {
		log.Println("not set GitHub Personal Token, you can't push Private Repository.")
	}

	return g
}

func (g GitHub) Put(fp string) error {
	tfp, err := file.ToTreePath(fp)
	if err != nil {
		return errors.Wrap(err, "failed to convert tree path")
	}

	u := uuid.New()
	ghDir := fmt.Sprintf("%s-%s", "grasshopper", u.String())
	dir := filepath.Join(os.TempDir(), ghDir)
	destFile := filepath.Join(dir, tfp)
	err = os.MkdirAll(filepath.Dir(destFile), 0775)

	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: g.Repo,
		Auth: &http.BasicAuth{
			Username: g.User,
			Password: g.PersonalToken,
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to clone GitHub repository")
	}

	src, err := os.Open(fp)
	if err != nil {
		return errors.Wrap(err, "failed to open source file")
	}
	defer src.Close()

	dest, err := os.Create(destFile)
	if err != nil {
		return errors.Wrap(err, "failed to create dest file")
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		return errors.Wrap(err, "failed to copy file")
	}

	w, err := r.Worktree()
	if err != nil {
		return errors.Wrap(err, "failed to get git worktree")
	}

	_, err = w.Add(tfp)
	if err != nil {
		return errors.Wrap(err, "failed to git add file")
	}

	_, err = w.Commit("backup file by grasshopper", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "grasshopper",
			Email: "grasshopper@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to git commit")
	}

	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: g.User,
			Password: g.PersonalToken,
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to git push")
	}

	return nil
}
