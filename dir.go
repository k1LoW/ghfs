package ghfs

import (
	"errors"
	"io"
	"io/fs"
	"path/filepath"
)

var _ fs.File = (*dir)(nil)
var _ fs.ReadDirFile = (*dir)(nil)

type dir struct {
	path   string
	files  []fs.DirEntry
	offset int
}

func (d *dir) Stat() (fs.FileInfo, error) {
	return &dinfo{name: filepath.Base(d.path)}, nil
}

func (d *dir) Close() error {
	return nil
}

func (d *dir) Read(b []byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.path, Err: errors.New("is a directory")}
}

func (d *dir) ReadDir(count int) ([]fs.DirEntry, error) {
	n := len(d.files) - d.offset
	if count > 0 && n > count {
		n = count
	}
	if n == 0 {
		if count <= 0 {
			return nil, nil
		}
		return nil, io.EOF
	}
	list := make([]fs.DirEntry, n)
	for i := range list {
		list[i] = &dent{
			de: d.files[d.offset+i],
		}
	}
	d.offset += n
	return list, nil
}
