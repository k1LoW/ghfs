package ghfs

import (
	"io/fs"
)

var _ fs.DirEntry = (*dent)(nil)

type dent struct {
	de fs.DirEntry
}

func (de *dent) Name() string {
	return de.de.Name()
}

func (de *dent) IsDir() bool {
	return de.de.IsDir()
}

func (de *dent) Type() fs.FileMode {
	return de.de.Type()
}

func (de *dent) Info() (fs.FileInfo, error) {
	fi, err := de.de.Info()
	if err != nil {
		return nil, err
	}

	if de.de.IsDir() {
		return &dinfo{
			name: de.de.Name(),
		}, nil
	}

	return &finfo{
		name: fi.Name(),
		size: int64(fi.Sys().(int)),
		mode: fi.Mode(),
	}, nil
}
