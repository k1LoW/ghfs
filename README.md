# ghfs [![build](https://github.com/k1LoW/ghfs/actions/workflows/ci.yml/badge.svg)](https://github.com/k1LoW/ghfs/actions/workflows/ci.yml)

`ghfs` implements the [io/fs](https://pkg.go.dev/io/fs) interfaces for GitHub remote repositories.

The implementation wraps [go-github](https://github.com/google/go-github) client and use [Git Database API](https://docs.github.com/en/rest/reference/git).

## Supported interface

- [fs.FS](https://pkg.go.dev/io/fs#FS)
- [fs.ReadFileFS](https://pkg.go.dev/io/fs#ReadFileFS)

## References

- [johejo/ghfs](https://github.com/johejo/ghfs): Package ghfs wraps the github v3 rest api with io/fs.
