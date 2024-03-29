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

var (
	ErrNoChange = errors.New("No Change")
)

type GitHub struct {
	Repo          string
	User          string
	PersonalToken string
}

func NewGitHubClient(repo, user, token string) GitHub {
	g := GitHub{}

	if os.Getenv(envGitHubRepo) != "" {
		g.Repo = os.Getenv(envGitHubRepo)
	} else if repo != "" {
		g.Repo = repo
	} else {
		log.Fatal("must be set GitHub Repository")
	}

	if os.Getenv(envGitHubUser) != "" {
		g.User = os.Getenv(envGitHubUser)
	} else if user != "" {
		g.User = user
	} else {
		log.Println("not set GitHub Username, you can't push Private Repository.")
	}

	if os.Getenv(envGitHubToken) != "" {
		g.PersonalToken = os.Getenv(envGitHubToken)
	} else if token != "" {
		g.PersonalToken = token
	} else {
		log.Println("not set GitHub Personal Token, you can't push Private Repository.")
	}

	return g
}

func (g GitHub) mkDir(treePath string) (string, string, error) {
	// make base directory, ${TMPDIR}/grasshopper-${UUID}
	u := uuid.New()
	ghDir := fmt.Sprintf("%s-%s", "grasshopper", u.String())
	repoBaseDir := filepath.Join(os.TempDir(), ghDir)

	destFilePath := filepath.Join(repoBaseDir, treePath)
	err := os.MkdirAll(filepath.Dir(destFilePath), 0775)
	return repoBaseDir, destFilePath, err
}

func (g GitHub) copyTargetFile(sourcePath, destPath string) error {
	src, err := os.Open(sourcePath)
	if err != nil {
		return errors.Wrap(err, "failed to open source file")
	}
	defer src.Close()

	dest, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE, 0775)
	if err != nil {
		return errors.Wrap(err, "failed to create dest file")
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		return errors.Wrap(err, "failed to copy file")
	}

	return nil
}

func (g GitHub) Put(fp string) error {
	tfp, err := file.ToTreePath(fp)
	if err != nil {
		return errors.Wrap(err, "failed to convert tree path")
	}

	repoBaseDir, destFilePath, err := g.mkDir(tfp)

	r, err := git.PlainClone(repoBaseDir, false, &git.CloneOptions{
		URL: g.Repo,
		Auth: &http.BasicAuth{
			Username: g.User,
			Password: g.PersonalToken,
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to clone GitHub repository")
	}
	defer os.Remove(repoBaseDir)

	err = g.copyTargetFile(fp, destFilePath)
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

	status, err := w.Status()
	if err != nil {
		return errors.Wrap(err, "failed to git status")
	}

	if len(status) == 0 {
		// no change
		return ErrNoChange
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
