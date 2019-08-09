package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/whywaita/grasshopper/storage"
)

func main() {
	args := os.Args

	if len(args) == 1 {
		log.Fatal("grasshopper")
	}

	os.Exit(run(args[1]))
}

func run(filePath string) int {
	path, err := getFullPath(filePath)
	if err != nil {
		log.Fatal(err)
	}

	client := storage.NewGitHubClient()
	err = client.Put(path)
	if errors.Cause(err) == storage.ErrNoChange {
		// no change
		return 1
	} else if err != nil {
		log.Printf("%+v\n", err)
		return 1
	}

	log.Printf("detect to change! backup is done!\nfile: %s\n", filePath)
	return 0
}

func getFullPath(path string) (string, error) {
	if !filepath.IsAbs(path) {
		p, err := filepath.Abs(path)
		if err != nil {
			return "", errors.Wrap(err, "failed to get absolute path")
		}

		return p, nil
	}

	return path, nil
}
