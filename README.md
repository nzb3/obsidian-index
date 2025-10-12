# Obsidian Index

A powerful CLI tool for automatically creating comprehensive indexes for your Obsidian vaults. This tool generates markdown files with links to all entries in each directory, processing from leaves to root for optimal organization.

## Features

- **Automatic Index Generation**: Creates index files for each directory in your Obsidian vault
- **Smart Processing**: Processes directories from deepest to shallowest levels
- **Flexible Configuration**: Support for dry-run mode, backups, and directory exclusions
- **Safe Operations**: Atomic file operations to prevent data corruption
- **Cross-Platform**: Works on macOS, Linux, and Windows

## Installation

### Homebrew (macOS)

```bash
brew install nzb3/obsidian-index/obsidian-index
```

### Manual Installation

1. Download the latest release from the [releases page](https://github.com/nzb3/obsidian-index/releases)
2. Extract the binary to a directory in your PATH
3. Make it executable: `chmod +x obsidian-index`

### Build from Source

```bash
git clone https://github.com/nzb3/obsidian-index.git
cd obsidian-index
go build -o obsidian-index ./cmd
```

## Usage

### Basic Usage

```bash
obsidian-index init --dir /path/to/your/obsidian/vault
```

### Advanced Options

```bash
obsidian-index init \
  --dir /path/to/your/obsidian/vault \
  --verbose \
  --dry-run \
  --backup \
  --exclude "templates" \
  --exclude "attachments"
```

### Command Options

- `--dir, -d`: Path to the Obsidian vault directory (required)
- `--verbose, -v`: Enable verbose output for detailed logging
- `--dry-run`: Show what would be done without creating files
- `--backup`: Create backup of existing index files before overwriting
- `--exclude`: Directories to exclude from indexing (can be used multiple times)

## How It Works

1. **Directory Discovery**: Recursively scans your Obsidian vault for directories
2. **Leaf-First Processing**: Processes directories from deepest to shallowest levels
3. **Index Generation**: Creates markdown files with links to all files and subdirectories
4. **Smart Naming**: Index files are named after their parent directory (e.g., `notes.md` for a `notes/` directory)
5. **Link Format**: Uses Obsidian's `[[link]]` format for all generated links

## Example Output

Given a vault structure like:
```
MyVault/
├── notes/
│   ├── project-a/
│   │   ├── meeting-notes.md
│   │   └── todo.md
│   └── project-b/
│       └── research.md
└── templates/
    └── daily-note.md
```

The tool will create:
- `MyVault/index.md` with links to `notes/` and `templates/`
- `MyVault/notes/notes.md` with links to `project-a/` and `project-b/`
- `MyVault/notes/project-a/project-a.md` with links to `meeting-notes.md` and `todo.md`
- `MyVault/notes/project-b/project-b.md` with links to `research.md`

## Safety Features

- **Dry Run Mode**: Test the tool without making changes
- **Backup Support**: Automatically backup existing index files
- **Atomic Operations**: Uses temporary files to prevent corruption
- **Permission Handling**: Gracefully handles permission errors
- **Validation**: Validates vault directory before processing

## Configuration

The tool respects the following configuration options:

- **Exclude Directories**: Skip specific directories from indexing
- **Verbose Logging**: Get detailed information about the indexing process
- **Backup Mode**: Automatically backup existing files before overwriting
- **Dry Run**: Preview changes without modifying files

## Development

### Prerequisites

- Go 1.25 or later
- Git

### Building

```bash
# Clone the repository
git clone https://github.com/nzb3/obsidian-index.git
cd obsidian-index

# Build the binary
go build -o obsidian-index ./cmd

# Run tests
go test ./...
```

### Project Structure

```
obsidian-index/
├── cmd/                    # Main application entry point
├── internal/
│   ├── app/               # Application logic
│   ├── cmd/               # CLI command definitions
│   ├── config/            # Configuration management
│   ├── indexator/         # Core indexing logic
│   └── version/           # Version information
├── Formula/               # Homebrew formula
└── scripts/              # Build scripts
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes and add tests
4. Commit your changes: `git commit -m "feat: add new feature"`
5. Push to the branch: `git push origin feature-name`
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/nzb3/obsidian-index/issues) page
2. Create a new issue with detailed information about your problem
3. Include your operating system, Obsidian version, and any error messages

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes and version history.
