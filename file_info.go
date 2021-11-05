package ghfs

import (
	"io/fs"
	"time"
)

var _ fs.FileInfo = (*finfo)(nil)
var _ fs.FileInfo = (*dinfo)(nil)

type finfo struct {
	name string
	size int64
	mode fs.FileMode
}

func (fi *finfo) Name() string {
	return fi.name
}
func (fi *finfo) Size() int64 {
	return fi.size
}
func (fi *finfo) Mode() fs.FileMode {
	return fi.mode
}
func (fi *finfo) ModTime() time.Time {
	return time.Time{}
}
func (fi *finfo) IsDir() bool {
	return false
}
func (fi *finfo) Sys() interface{} {
	return fi
}

type dinfo struct {
	name string
}

func (d *dinfo) Name() string {
	return d.name
}

func (d *dinfo) Size() int64 {
	return 0
}

func (d *dinfo) Mode() fs.FileMode {
	return fs.ModeDir | 0o555
}

func (d *dinfo) ModTime() time.Time {
	return time.Time{}
}

func (d *dinfo) IsDir() bool {
	return true
}

func (d *dinfo) Sys() interface{} {
	return d.name
}
