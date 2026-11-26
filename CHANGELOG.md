# Changelog

All notable changes to this project are documented in this file. The format
is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this
project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `Option`, `Result` and `Either` with generic `Map`, `AndThen` and `Fold` methods, iterator and JSON support.
- Benchmarks showing that chains cost the same as plain Go (`make bench`).
- A Makefile with the development tasks (`make help`).
- Requires Go 1.27 or later.
