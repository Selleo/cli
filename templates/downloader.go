package templates

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/google/go-github/v64/github"
)

func NewDownloader() (*downloader, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN is required")
	}
	client := github.NewClient(nil).WithAuthToken(token)
	return &downloader{
		client: client,
	}, nil
}

type downloader struct {
	client *github.Client
}

type DownloadInput struct {
	RepoOwner   string
	RepoName    string
	FolderPath  string
	Destination string
	Force       bool
}

func (d *downloader) Download(ctx context.Context, opts *DownloadInput) error {
	files := []*github.RepositoryContent{}
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Prefix = "Downloading files "
	s.Color("green")
	s.Start()
	defer s.Stop()

	err := d.walk(ctx, opts.RepoOwner, opts.RepoName, opts.FolderPath, func(file *github.RepositoryContent) {
		files = append(files, file)
	})

	wg := sync.WaitGroup{}
	wg.Add(len(files))
	defer wg.Wait()

	for _, file := range files {
		full := path.Join(opts.Destination, file.GetPath())
		dir := path.Dir(full)
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			return err
		}
	}

	for _, file := range files {
		go func(file *github.RepositoryContent) {
			defer wg.Done()

			resp, err := http.Get(file.GetDownloadURL())
			if err != nil {
				fmt.Println(err)
				return
			}
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			full := path.Join(opts.Destination, file.GetPath())
			_, err = os.Stat(full)
			exists := !os.IsNotExist(err)
			writeToFile := true
			if exists && !opts.Force {
				writeToFile = false
			}

			if writeToFile {
				f, err := os.Create(full)
				if err != nil {
					fmt.Println(err)
					return
				}
				_, err = f.Write(b)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}(file)
	}

	return err
}

func (d *downloader) walk(ctx context.Context, repoOwner string, repoName string, path string, walkFn func(file *github.RepositoryContent)) error {
	_, dir, _, err := d.client.Repositories.GetContents(
		ctx,
		repoOwner,
		repoName,
		path,
		&github.RepositoryContentGetOptions{Ref: "main"},
	)
	if err != nil {
		return err
	}
	for _, file := range dir {
		if file.GetType() == "file" {
			walkFn(file)
		} else {
			err := d.walk(ctx, repoOwner, repoName, file.GetPath(), walkFn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
