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

var (
	_ fs.FS         = (*FS)(nil)
	_ fs.ReadFileFS = (*FS)(nil)
	_ fs.ReadDirFS  = (*FS)(nil)

	ctx = context.Background()
)

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
	defer func() {
		_ = f.Close()
	}()

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

	data, size, err := fsys.readDataFromSHA(sha)
	if err != nil {
		return nil, err
	}

	return &file{
		path: name,
		data: data,
		size: size,
		fi:   fi,
	}, nil
}

func (fsys *FS) ReadFile(name string) ([]byte, error) {
	f, err := fsys.shafs.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	sha := string(b)

	data, _, err := fsys.readDataFromSHA(sha)
	if err != nil {
		return nil, err
	}

	return []byte(data), nil
}

func (fsys *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	f, err := fsys.shafs.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	files, err := fs.ReadDir(fsys.shafs, name)
	if err != nil {
		return nil, err
	}
	dents := []fs.DirEntry{}
	for _, f := range files {
		dents = append(dents, &dent{
			de: f,
		})
	}
	return dents, nil
}

func (fsys *FS) readDataFromSHA(sha string) (string, int, error) {
	blob, _, err := fsys.client.Git.GetBlob(ctx, fsys.owner, fsys.repo, sha)
	if err != nil {
		return "", 0, err
	}

	encoding := blob.GetEncoding()
	var data string

	switch encoding {
	case "base64":
		if blob.Content == nil {
			return "", 0, errors.New("malformed response: base64 encoding of null content")
		}
		c, err := base64.StdEncoding.DecodeString(blob.GetContent())
		if err != nil {
			return "", 0, err
		}
		data = string(c)
	case "":
		if blob.Content != nil {
			data = blob.GetContent()
		}
	default:
		return "", 0, fmt.Errorf("unsupported content encoding: %v", encoding)
	}
	return data, blob.GetSize(), nil
}

type config struct {
	client *github.Client
	ctx    context.Context
	branch string
}

type Option func(*config) error

func Client(client *github.Client) Option {
	return func(c *config) error {
		if client != nil {
			c.client = client
		}
		return nil
	}
}

func Context(ctx context.Context) Option {
	return func(c *config) error {
		if ctx != nil {
			c.ctx = ctx
		}
		return nil
	}
}

func Branch(branch string) Option {
	return func(c *config) error {
		if branch != "" {
			c.branch = branch
		}
		return nil
	}
}

func New(owner, repo string, opts ...Option) (*FS, error) {
	c := &config{}
	for _, o := range opts {
		if err := o(c); err != nil {
			return nil, err
		}
	}
	if c.client == nil {
		client, err := factory.NewGithubClient()
		if err != nil {
			return nil, err
		}
		c.client = client
	}
	if c.ctx == nil {
		c.ctx = context.Background()
	}
	if c.branch == "" {
		r, _, err := c.client.Repositories.Get(c.ctx, owner, repo)
		if err != nil {
			return nil, err
		}
		c.branch = r.GetDefaultBranch()
	}

	b, _, err := c.client.Repositories.GetBranch(c.ctx, owner, repo, c.branch, false)
	if err != nil {
		return nil, err
	}
	sha := b.GetCommit().GetSHA()

	t, _, err := c.client.Git.GetTree(c.ctx, owner, repo, sha, true)
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
		client: c.client,
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
