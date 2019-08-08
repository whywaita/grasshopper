package main

import (
	"fmt"
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

	path, err := getFullPath(args[1])
	if err != nil {
		log.Fatal(err)
	}

	client := storage.NewGitHubClient()
	err = client.Put(path)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	log.Println("file upload is done.")
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
