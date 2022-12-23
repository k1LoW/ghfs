package ghfs

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strconv"
	"strings"
	"testing/fstest"

	"github.com/google/go-github/v48/github"
	"github.com/k1LoW/go-github-client/v48/factory"
)

var (
	_ fs.FS         = (*FS)(nil)
	_ fs.ReadFileFS = (*FS)(nil)
	_ fs.ReadDirFS  = (*FS)(nil)
	_ fs.SubFS      = (*FS)(nil)

	ctx = context.Background()
)

type FS struct {
	client *github.Client
	owner  string
	repo   string
	shafs  fs.FS
	prefix string
}

func (fsys *FS) Open(name string) (fs.File, error) {
	name = fsys.nameWithPrefix(name)

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
	name = fsys.nameWithPrefix(name)
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
	name = fsys.nameWithPrefix(name)
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

func (fsys *FS) Sub(dir string) (fs.FS, error) {
	fsys.prefix = dir
	return fsys, nil
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
	tag    string
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

func Tag(tag string) Option {
	return func(c *config) error {
		if tag != "" {
			c.tag = tag
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
	if c.tag != "" && c.branch != "" {
		return nil, errors.New("only one of tag and branch can be specified")
	}

	var sha string
	if c.tag != "" {
		page := 1
	L:
		for {
			tags, res, err := c.client.Repositories.ListTags(c.ctx, owner, repo, &github.ListOptions{
				Page:    page,
				PerPage: 100,
			})
			if err != nil {
				return nil, err
			}
			for _, t := range tags {
				if c.tag == t.GetName() {
					sha = t.GetCommit().GetSHA()
					break L
				}
			}
			if res.NextPage == 0 {
				break
			}
			page += 1
		}
		if sha == "" {
			return nil, fmt.Errorf("tag '%s' not fount", c.tag)
		}
	} else {
		r, _, err := c.client.Repositories.Get(c.ctx, owner, repo)
		if err != nil {
			return nil, err
		}

		if c.branch == "" {
			c.branch = r.GetDefaultBranch()
		}

		b, _, err := c.client.Repositories.GetBranch(c.ctx, owner, repo, c.branch, false)
		if err != nil {
			if c.branch == r.GetDefaultBranch() {
				// empty repository
				return &FS{
					client: c.client,
					owner:  owner,
					repo:   repo,
					shafs:  fstest.MapFS{},
				}, nil
			}
			return nil, err
		}
		sha = b.GetCommit().GetSHA()
	}

	t, _, err := c.client.Git.GetTree(c.ctx, owner, repo, sha, true)
	if err != nil {
		return nil, err
	}

	shafs := fstest.MapFS{}
	for _, e := range t.Entries {
		if e.GetType() == "blob" {
			m, err := filemode(e.GetMode())
			if err != nil {
				return nil, err
			}
			shafs[e.GetPath()] = &fstest.MapFile{
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
		shafs:  shafs,
	}, nil
}

func (fsys *FS) nameWithPrefix(name string) string {
	if fsys.prefix == "" {
		return name
	}
	if name == "." {
		return strings.TrimPrefix(strings.TrimSuffix(fsys.prefix, "/"), "/")
	}
	return strings.TrimPrefix(fmt.Sprintf("%s/%s", strings.TrimSuffix(fsys.prefix, "/"), name), "/")
}

func filemode(s string) (fs.FileMode, error) {
	n, err := strconv.ParseUint(s, 8, 32)
	if err != nil {
		return fs.FileMode(0), err
	}
	return fs.FileMode(n), nil
}
