# Changelog

## [v1.4.0](https://github.com/k1LoW/ghfs/compare/v1.3.2...v1.4.0) - 2025-02-28
- Bump golang.org/x/crypto from 0.21.0 to 0.31.0 in /example by @dependabot in https://github.com/k1LoW/ghfs/pull/32
- chore: update Go version in tests and go.mod to 1.23 by @k1LoW in https://github.com/k1LoW/ghfs/pull/36
- [breaking changes] Upgrade go-github and go-github-client by @wreulicke in https://github.com/k1LoW/ghfs/pull/34
- Bump golang.org/x/net from 0.23.0 to 0.33.0 in /example by @dependabot in https://github.com/k1LoW/ghfs/pull/35

## [v1.3.2](https://github.com/k1LoW/ghfs/compare/v1.3.1...v1.3.2) - 2024-11-27
- Bump github.com/golang-jwt/jwt/v4 from 4.5.0 to 4.5.1 by @dependabot in https://github.com/k1LoW/ghfs/pull/29
- Bump github.com/cli/go-gh/v2 from 2.6.0 to 2.11.1 by @dependabot in https://github.com/k1LoW/ghfs/pull/31

## [v1.3.1](https://github.com/k1LoW/ghfs/compare/v1.3.0...v1.3.1) - 2024-04-19
- Bump golang.org/x/net from 0.17.0 to 0.23.0 in /example by @dependabot in https://github.com/k1LoW/ghfs/pull/27

## [v1.3.0](https://github.com/k1LoW/ghfs/compare/v1.2.0...v1.3.0) - 2024-04-16
- Bump google.golang.org/protobuf from 1.28.0 to 1.33.0 in /example by @dependabot in https://github.com/k1LoW/ghfs/pull/24
- Update go-github-client ( and go version ) by @k1LoW in https://github.com/k1LoW/ghfs/pull/26

## [v1.1.1](https://github.com/k1LoW/ghfs/compare/v1.1.0...v1.1.1) - 2024-01-26
- Bump golang.org/x/net from 0.12.0 to 0.17.0 by @dependabot in https://github.com/k1LoW/ghfs/pull/17
- Bump golang.org/x/net from 0.7.0 to 0.17.0 in /example by @dependabot in https://github.com/k1LoW/ghfs/pull/18
- Bump golang.org/x/crypto from 0.14.0 to 0.17.0 by @dependabot in https://github.com/k1LoW/ghfs/pull/21
- Bump golang.org/x/crypto from 0.14.0 to 0.17.0 in /example by @dependabot in https://github.com/k1LoW/ghfs/pull/20
- Bump github.com/cloudflare/circl from 1.3.3 to 1.3.7 by @dependabot in https://github.com/k1LoW/ghfs/pull/22
- Update pkgs by @k1LoW in https://github.com/k1LoW/ghfs/pull/23

## [v1.0.2](https://github.com/k1LoW/ghfs/compare/v1.0.1...v1.0.2) - 2023-07-08
- Update pkgs by @k1LoW in https://github.com/k1LoW/ghfs/pull/14
- Setup tagpr by @k1LoW in https://github.com/k1LoW/ghfs/pull/15

## [v1.0.1](https://github.com/k1LoW/ghfs/compare/v1.0.0...v1.0.1) (2023-02-19)

* Fix go.* [#13](https://github.com/k1LoW/ghfs/pull/13) ([k1LoW](https://github.com/k1LoW))

## [v1.0.0](https://github.com/k1LoW/ghfs/compare/v0.9.0...v1.0.0) (2023-02-19)

* Update pkgs [#12](https://github.com/k1LoW/ghfs/pull/12) ([k1LoW](https://github.com/k1LoW))
* Update go and pkgs [#11](https://github.com/k1LoW/ghfs/pull/11) ([k1LoW](https://github.com/k1LoW))
* Bump golang.org/x/net from 0.5.0 to 0.7.0 [#10](https://github.com/k1LoW/ghfs/pull/10) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v0.9.0](https://github.com/k1LoW/ghfs/compare/v0.8.0...v0.9.0) (2023-02-12)

* Bump up go-github-client version [#9](https://github.com/k1LoW/ghfs/pull/9) ([k1LoW](https://github.com/k1LoW))

## [v0.8.0](https://github.com/k1LoW/ghfs/compare/v0.7.0...v0.8.0) (2022-12-23)

* Bump up go and pkgs version [#8](https://github.com/k1LoW/ghfs/pull/8) ([k1LoW](https://github.com/k1LoW))

## [v0.7.0](https://github.com/k1LoW/ghfs/compare/v0.6.0...v0.7.0) (2022-08-22)

* Add octocov [#7](https://github.com/k1LoW/ghfs/pull/7) ([k1LoW](https://github.com/k1LoW))
* Bump up go and go-github version [#6](https://github.com/k1LoW/ghfs/pull/6) ([k1LoW](https://github.com/k1LoW))

## [v0.6.0](https://github.com/k1LoW/ghfs/compare/v0.5.0...v0.6.0) (2021-11-14)

* Return an empty FS if the repository is a zero-commit repository and a default branch has been specified. [#5](https://github.com/k1LoW/ghfs/pull/5) ([k1LoW](https://github.com/k1LoW))

## [v0.5.0](https://github.com/k1LoW/ghfs/compare/v0.4.0...v0.5.0) (2021-11-08)

* Support fs.SubFS [#4](https://github.com/k1LoW/ghfs/pull/4) ([k1LoW](https://github.com/k1LoW))

## [v0.4.0](https://github.com/k1LoW/ghfs/compare/v0.3.0...v0.4.0) (2021-11-07)

* Add Branch and Tag options [#3](https://github.com/k1LoW/ghfs/pull/3) ([k1LoW](https://github.com/k1LoW))

## [v0.3.0](https://github.com/k1LoW/ghfs/compare/v0.2.1...v0.3.0) (2021-11-07)

* [BREAKING] Use functional option pattern [#2](https://github.com/k1LoW/ghfs/pull/2) ([k1LoW](https://github.com/k1LoW))

## [v0.2.1](https://github.com/k1LoW/ghfs/compare/v0.2.0...v0.2.1) (2021-11-07)


## [v0.2.0](https://github.com/k1LoW/ghfs/compare/v0.1.0...v0.2.0) (2021-11-06)

* Support fs.ReadDirFS [#1](https://github.com/k1LoW/ghfs/pull/1) ([k1LoW](https://github.com/k1LoW))

## [v0.1.0](https://github.com/k1LoW/ghfs/compare/a4b05ac393a8...v0.1.0) (2021-11-06)
