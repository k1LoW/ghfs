package ghfs

import (
	"io"
	"io/fs"
	"path/filepath"
)

var _ fs.File = (*file)(nil)

type file struct {
	path   string
	data   string
	size   int
	fi     fs.FileInfo
	offset int64
}

func (f *file) Stat() (fs.FileInfo, error) {
	return &finfo{
		name: filepath.Base(f.path),
		size: int64(f.size),
		mode: f.fi.Mode(),
	}, nil
}

func (f *file) Close() error {
	return nil
}

func (f *file) Read(b []byte) (int, error) {
	if f.offset >= int64(len(f.data)) {
		return 0, io.EOF
	}
	if f.offset < 0 {
		return 0, &fs.PathError{Op: "read", Path: f.path, Err: fs.ErrInvalid}
	}
	n := copy(b, f.data[f.offset:])
	f.offset += int64(n)
	return n, nil
}
