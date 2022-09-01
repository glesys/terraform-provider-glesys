# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).
## Unreleased

## [0.4.0] - 2022-09-01
### Added
- Support for IP resources @alxbse
### Changed
- Moved to golangci-lint for static analysis
- Updated test and build versions of Go
## [0.3.2] - 2022-01-11
### Changed
- Added check for empty Object Storage Instance credentials list @norrland (#36)
## [0.3.1] - 2021-11-11
### Changed
- Fixed typo in DNS example - (Thanks to simon @ bamze)
- Bumped glesys-go to v3.0.0
## [0.3.0] - 2021-10-15
### Added
- New handling of KVM templates @larsve
### Changed
- Bump Terraform SDK version
- Fixed go test warning in makefile
- Bump Vagrantfile box version
## [0.2.0] - 2021-02-15
### Added
- resource `glesys_objectstorage_credential` @xevz (#13)
- resource `glesys_objectstorage_instance` @xevz (#13)
- Acceptance testing for `glesys_objectstorage_*` @xevz (#13)
- KVM servers in `glesys_server` now can be created with multiple users,
and multiple keys per user. @norrland (#20)
### Changed
- Moved Makefile to GNUMakefile @xevz (#13)
- server: Ignore case sensitivity on hostname (#22)

## [0.1.0] - 2020-12-18
### Added
- Initial release
- Support for `servers`, `dns`, `loadbalancers`, `networks`,
  `networkadapters`.
