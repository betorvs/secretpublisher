# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.3] - 2021-02-18
### Added
- scan-secrets and scan-configmaps to scan secrets or config maps using a label to search and copy them to secret receiver
- secret-subvalue command to scan secrets and sync just a subkey.


## [0.0.2] - 2021-02-08
### Changed
- Update golang version to 1.15
- Removed travis-ci files

### Added
- Add Labels

## [0.0.1] - 2019-12-10
### Added
- Sub-commands: create, update, check, delete to manage secrets and exist to check if one secret already exist and is up-to-date.
- ENCODING_REQUEST variable to accepted only encoded requests.
