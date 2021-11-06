# ghfs [![build](https://github.com/k1LoW/ghfs/actions/workflows/ci.yml/badge.svg)](https://github.com/k1LoW/ghfs/actions/workflows/ci.yml)

`ghfs` provides a GitHub remote repository implementation for Go io/fs interface.

The implementation wraps [go-github](https://github.com/google/go-github) client and use [Git Database API](https://docs.github.com/en/rest/reference/git).

## Supported interface

- [fs.FS](https://pkg.go.dev/io/fs#FS)
- [fs.ReadFileFS](https://pkg.go.dev/io/fs#ReadFileFS)

## References

- [johejo/ghfs](https://github.com/johejo/ghfs): Package ghfs wraps the github v3 rest api with io/fs.
