# Dusty

**D**irectory **U**sage **ST**atistics s**Y**stem

A beautiful terminal-based directory size analyzer, similar to Linux's `du` command or Windows TreeSize Free. Built with Go and features an interactive TUI (Terminal User Interface).

## Example

```
 Dusty - /Users/username/projects
Total Size: 2.5 GB | Sort: size

      1.2 GB  48.0% [█████████████████████░] node_modules/
    850.3 MB  34.0% [███████████████████░░░] build/
    250.0 MB  10.0% [██████████░░░░░░░░░░░░] .git/
    125.5 MB   5.0% [█████░░░░░░░░░░░░░░░░░] dist/
     50.2 MB   2.0% [██░░░░░░░░░░░░░░░░░░░░] docs/
      8.5 MB   0.3% [░░░░░░░░░░░░░░░░░░░░░░] src/
    512.0 KB   0.0% [░░░░░░░░░░░░░░░░░░░░░░] README.md

  ↑/↓: Navigate | ←/→: Enter/Exit | s: Sort | Home: Root | q: Quit
```

## Features

- **Recursive Directory Scanning**: Calculates cumulative storage usage for all folders and files
- **Interactive Navigation**: Browse through your directory tree with keyboard controls
- **Visual Size Bars**: See relative sizes at a glance with progress bars
- **Multiple Sort Options**: Sort by size or name
- **Colorful Display**: Syntax-highlighted output with distinct colors for directories and files
- **Human-Readable Sizes**: Automatic formatting (B, KB, MB, GB, TB, PB, EB)
- **Percentage Display**: Shows what percentage each item takes of its parent directory
- **Fast Performance**: Concurrent scanning for quick results

## Installation

### Download Pre-built Binary (Recommended)

Download the latest release for your platform from the [Releases page](https://github.com/twiddles/dusty/releases).

**Linux/macOS:**
```bash
# Download and install (replace URL with actual release URL and platform)
curl -L https://github.com/twiddles/dusty/releases/latest/download/dusty-linux-amd64 -o dusty
chmod +x dusty
sudo mv dusty /usr/local/bin/
```

**Windows:**
Download the `.exe` file and add it to your PATH.

### Install with Go

```bash
go install github.com/twiddles/dusty@latest
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/twiddles/dusty.git
cd dusty

# Build the binary
go build -o dusty

# Optional: Install to your PATH
go install
```

## Usage

### Basic Usage

Scan the current directory:
```bash
./dusty
```

Scan a specific directory:
```bash
./dusty /path/to/directory
```

Or if installed:
```bash
dusty ~/Documents
```

### Keyboard Controls

| Key | Action |
|-----|--------|
| `↑` / `k` | Move cursor up |
| `↓` / `j` | Move cursor down |
| `→` / `Enter` / `l` | Enter selected directory |
| `←` / `Backspace` / `h` | Go back to parent directory |
| `d` | Delete highlighted file or folder |
| `s` | Toggle sort (size/name) |
| `Home` | Return to root directory |
| `q` / `Ctrl+C` | Quit application |

## How It Works

1. **Scanning Phase**: The application recursively walks through the directory tree, calculating the size of each file and aggregating sizes up the tree
2. **Interactive Display**: Uses Bubble Tea framework for a smooth TUI experience
3. **Navigation**: You can drill down into directories to see their contents and sizes
4. **Sorting**: Toggle between sorting by size (largest first) or alphabetically by name

## Technical Details

- **Language**: Go 1.21+
- **UI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions for terminal output
- **Architecture**: Clean separation between data model and view

## Dependencies

```
github.com/charmbracelet/bubbletea v1.3.10
github.com/charmbracelet/lipgloss v1.1.0
```

## Project Structure

```
dusty/
├── main.go                        # Main application code
├── go.mod                         # Go module definition
├── go.sum                         # Dependency checksums
├── Makefile                       # Build automation
├── .goreleaser.yml               # GoReleaser configuration
├── .github/workflows/            # GitHub Actions for releases
├── LICENSE                        # MIT License
└── README.md                     # This file
```

## Development

### Build Commands

```bash
make build        # Build for current platform
make build-all    # Build for all platforms (Linux, macOS, Windows)
make install      # Install to $GOPATH/bin
make test         # Run tests
make clean        # Remove build artifacts
make help         # Show all available commands
```

**Note for Windows users:** If you don't have `make`, you can build directly with:
```bash
go build -o dusty.exe
```

## Performance Tips

- For very large directory trees (millions of files), the initial scan may take some time
- The application skips directories it doesn't have permission to access
- Symbolic links are followed, which may lead to duplicate counting if they point within the scanned tree

## Limitations

- Does not detect hard links (files may be counted multiple times)
- Symbolic links are treated as regular files/directories
- Requires read permissions for all directories you want to scan

## Contributing

Feel free to open issues or submit pull requests with improvements!

## License

MIT License - Feel free to use this project for any purpose.

## Credits

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) by Charm
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) by Charm
