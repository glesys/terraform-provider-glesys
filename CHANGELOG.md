# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).
## Unreleased

## [0.13.0] -
### Added
### Changed
- `glesys_server_disk` - Added `Type` attribute.

## [0.12.0] - 2024-09-23
### Added
- `glesys_privatenetwork` - Manage PrivateNetwork resources.
- `glesys_privatenetwork_segment` - Manage PrivateNetwork segments.

### Changed
- `glesys_networkadapter` - Added Importer function.
- `glesys_server` - `networkadapter` attribute listing attached NICS.

## [0.11.3] - 2024-06-10
### Changed
- Fix goreleaser config

## [0.11.2] - 2024-06-10
### Changed
- Bump dependencies

## [0.11.1] - 2024-03-08
### Changed
- Fix possible race condition when creating multiple `glesys_ip` resources.
- Remove OpenVZ examples
- Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.32.0 to 2.33.0
- Bump go version in tests

## [0.11.0] - 2024-01-24
### Added
- Resource `glesys_server_disk` for managing extra disks in the VMware platform.

### Changed
- Update glesys-go to v8.1.0
- Add 'extra_disks' attribute to `glesys_server`

## [0.10.3] - 2024-01-09
### Changed
- Update dependencies with high vulnerability rating
- Update terraform-plugin-sdk/v2 dependency

## [0.10.2] - 2023-12-19
### Changed
- Update github actions dependencies
- Update golang.org/x/crypto dependencies

## [0.10.1] - 2023-11-23
### Changed
- Updated the minimum go version to 1.20

## [0.10.0] - 2023-11-23
### Changed
- Update docs for network resources.
- Update docs for networkadapter resources.
- `glesys_emailaccount` API no longer accept user defined passwords.
- `glesys_loadbalancer_target` requires `targetip` attribute.
- Fix links to API documentation (@stemid)
- Bump glesys-go to v8

### Added
- Add `glesys_emailalias`
- Add docs for `glesys_objectstorage_*`
- Add docs for `glesys_ip`

## [0.9.0] - 2023-05-03
### Changed
- Update context in API calls.
- Handle error when deleting an already deleted DNSDomain Record.
- Update tests to use Go 1.20.
- Various dependency updates.
- Fixed issue when using UUID templates in KVM.
- Fixed issue with missing DNS Records not being recreated. (#126)

### Added
- Implemented `datasource_glesys_network`.

## [0.8.0] - 2022-11-01
### Changed
- Update workflow releaser to use goreleaser-action v3.2.0
- Change `glesys_server` `cloudconfigparams` type to TypeMap.
- Update glesys-go to v6.1.0

### Added
- `api_endpoint` provider configuration variable for setting base URL for the
  API requests.

## [0.7.1] - 2022-10-26
### Added
- Additional linters
### Changed
- Update resources with `CreateContext/ReadContext/UpdateContext/DeleteContext`
  to support request cancellation and warning diagnostics.
- Fixing linter errors
## [0.7.0] - 2022-10-20
### Added
- Implement DataSource for `glesys_dnsdomain`
### Changed
- Bump dependencies (#80)

## [0.6.0] - 2022-10-06
### Added
- Implement `glesys_emailaccount` @norrland (#75)
### Changed
- Loadbalancer deprecate `blacklist` @norrland (#73)
- Support import in `glesys_server` @norrland (#66)
- Support `cloudconfig` and `cloudconfigparams` in `glesys_server`
- Update glesys-go to v5

## [0.5.0] - 2022-09-21
### Changed
- Update to glesys-go/v4 @norrland #59
- Refactor `glesys_server` Create and Delete to wait for attribute.
- Fix typo in resource_glesys_server @norrland
- Fix wrong header in docs/index.md @norrland
- Update `glesys_domain_record` type to Required

## [0.4.6] - 2022-09-13
### Added
- Extended documentation for glesys_dnsdomain_record)
### Changed
- Build on Go 1.18, 1.19. Bump terraform-plugin-sdk/v2 to v2.22.0 (#54)
## [0.4.5] - 2022-09-12
### Added
- Documentation for glesys_server, glesys_dnsdomain(_record)
- Templates for provider and resources docs
## [0.4.4] - 2022-09-09
### Changed
- Skip tests on docs changes
- Fix typo in docs for provider version constraint
## [0.4.3] - 2022-09-08
### Added
- Updated documentation @norrland #51
### Changed
- glesys_servers: Skip updating bandwidth for KVM platform #50
## [0.4.2] - 2022-09-06
### Changed
- Changed the go version in go.mod #47
## [0.4.1] - 2022-09-02
### Added
- Setup goreleaser for Terraform Registry releases
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
