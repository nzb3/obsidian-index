# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of obsidian-index CLI tool
- Automatic index generation for Obsidian vaults
- Leaf-first directory processing
- Dry-run mode for safe testing
- Backup functionality for existing index files
- Directory exclusion support
- Verbose logging option
- Atomic file operations for data safety
- Cross-platform support (macOS, Linux, Windows)
- Homebrew installation support

### Features
- **Core Functionality**: Automatically creates markdown index files for each directory in an Obsidian vault
- **Smart Processing**: Processes directories from deepest to shallowest levels for optimal organization
- **Safety Features**: 
  - Dry-run mode to preview changes
  - Backup existing files before overwriting
  - Atomic file operations to prevent corruption
- **Flexibility**: 
  - Exclude specific directories from indexing
  - Verbose output for detailed logging
  - Support for various vault structures
- **CLI Interface**: 
  - Simple command-line interface with clear options
  - Comprehensive help and usage examples
  - Error handling and validation

### Technical Details
- Built with Go 1.25+
- Uses Cobra CLI framework for command handling
- Implements structured logging with slog
- Follows Go best practices and idiomatic patterns
- Comprehensive error handling and validation
- Cross-platform file system operations

## [1.0.0] - 2024-12-19

### Added
- Initial release
- Basic indexing functionality
- CLI interface with init command
- Configuration management
- Directory exclusion support
- Dry-run and backup modes
- Verbose logging
- Homebrew formula for easy installation
- Comprehensive documentation
- MIT License

### Changed
- N/A (initial release)

### Deprecated
- N/A (initial release)

### Removed
- N/A (initial release)

### Fixed
- N/A (initial release)

### Security
- N/A (initial release)

---

## Release Notes

### Version 1.0.0
This is the initial release of obsidian-index, a CLI tool for automatically creating comprehensive indexes for Obsidian vaults. The tool provides a safe and efficient way to generate markdown index files that link to all files and subdirectories within each directory of your vault.

**Key Features:**
- Automatic index generation with smart directory processing
- Safety features including dry-run mode and backup functionality
- Flexible configuration with directory exclusion support
- Cross-platform compatibility
- Easy installation via Homebrew

**Getting Started:**
```bash
# Install via Homebrew
brew install nzb3/obsidian-index/obsidian-index

# Basic usage
obsidian-index init --dir /path/to/your/vault

# Advanced usage with options
obsidian-index init --dir /path/to/your/vault --verbose --dry-run --backup
```

**Documentation:**
- Comprehensive README with usage examples
- CLI help with `obsidian-index --help`
- MIT License for open source usage
