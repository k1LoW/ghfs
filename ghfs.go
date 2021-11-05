package ghfs

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strconv"
	"testing/fstest"

	"github.com/google/go-github/v39/github"
	"github.com/k1LoW/go-github-client/v39/factory"
)

var _ fs.FS = (*FS)(nil)
var ctx = context.Background()

type FS struct {
	client *github.Client
	owner  string
	repo   string
	shafs  fs.FS
}

func (fsys *FS) Open(name string) (fs.File, error) {
	f, err := fsys.shafs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		files, err := fs.ReadDir(fsys.shafs, name)
		if err != nil {
			return nil, err
		}
		return &dir{
			path:  name,
			files: files,
		}, nil
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	sha := string(b)
	blob, _, err := fsys.client.Git.GetBlob(ctx, fsys.owner, fsys.repo, sha)
	if err != nil {
		return nil, err
	}

	encoding := blob.GetEncoding()
	var data string

	switch encoding {
	case "base64":
		if blob.Content == nil {
			return nil, errors.New("malformed response: base64 encoding of null content")
		}
		c, err := base64.StdEncoding.DecodeString(blob.GetContent())
		if err != nil {
			return nil, err
		}
		data = string(c)
	case "":
		if blob.Content == nil {
			data = ""
		}
		data = blob.GetContent()
	default:
		return nil, fmt.Errorf("unsupported content encoding: %v", encoding)
	}

	return &file{
		path: name,
		data: data,
		size: blob.GetSize(),
		fi:   fi,
	}, nil
}

func New(owner, repo string) (*FS, error) {
	client, err := factory.NewGithubClient()
	if err != nil {
		return nil, err
	}
	return NewWithGithubClient(client, owner, repo)
}

func NewWithGithubClient(client *github.Client, owner, repo string) (*FS, error) {
	ctx := context.Background()

	c, _, err := client.Repositories.ListCommits(ctx, owner, repo, &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: 1,
			Page:    1,
		},
	})
	if err != nil {
		return nil, err
	}
	if len(c) == 0 {
		return nil, fmt.Errorf("repository not found: %s/%s", owner, repo)
	}

	sha := c[0].GetSHA()
	t, _, err := client.Git.GetTree(ctx, owner, repo, sha, true)
	if err != nil {
		return nil, err
	}

	fsys := fstest.MapFS{}
	for _, e := range t.Entries {
		if e.GetType() == "blob" {
			m, err := filemode(e.GetMode())
			if err != nil {
				return nil, err
			}
			fsys[e.GetPath()] = &fstest.MapFile{
				Data: []byte(e.GetSHA()),
				Mode: m,
				Sys:  e.GetSize(),
			}
		}
	}

	return &FS{
		client: client,
		owner:  owner,
		repo:   repo,
		shafs:  fsys,
	}, nil
}

func filemode(s string) (fs.FileMode, error) {
	n, err := strconv.ParseUint(s, 8, 32)
	if err != nil {
		return fs.FileMode(0), err
	}
	return fs.FileMode(n), nil
}
