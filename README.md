# Witness

A lightweight, efficient screen capture tool for macOS that records screen content and saves as GIF or MP4 files.

## Overview

Witness is a command-line screen recorder designed to replace clunky GUI tools like Giphy Capture. It focuses on creating small, efficient files optimized for conveying concepts - perfect for documentation, tutorials, and bug reports.

## Features

- **Efficient Formats**: Export as optimized GIF or MP4
- **Quality Control**: Multiple compression levels from maximum compression to high quality
- **Flexible Capture**: Full screen or specific regions
- **Command-line Driven**: Fast and scriptable
- **macOS Native**: Uses Core Graphics for high-performance capture

## Installation

### Prerequisites

- macOS 10.12 or later
- Go 1.21 or later
- Xcode Command Line Tools
- [Mise](https://mise.jdx.dev/) (recommended) or Make

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Install Mise (recommended)
curl https://mise.run | sh
```

### Build from Source

```bash
git clone https://github.com/ericmhalvorsen/witness.git
cd witness

# Using Mise (recommended)
mise run build

# Or using Make
make build

# Or directly with Go
go build -o witness ./cmd/witness
```

## Usage

### Quick Start

1. **Select a capture region interactively:**
```bash
witness select -name demo
```
This launches macOS's native selection tool - just click and drag to select your capture area!

2. **Record a GIF using your saved region:**
```bash
witness gif -region demo -o demo.gif
```

### Region Selection

Witness makes it easy to select and reuse screen regions:

```bash
# Interactive selection (click and drag)
witness select

# Save the selection for later use
witness select -name myarea

# Save and set as default
witness select -name myarea -default

# List all saved regions
witness regions

# Delete a saved region
witness regions -delete myarea
```

### GIF Recording

```bash
# Record using a saved region
witness gif -region demo -o demo.gif

# Record using manual coordinates
witness gif -r 0,0,800,600 -o demo.gif

# Record at lower FPS for smaller files
witness gif -region demo -o demo.gif -f 10

# Record with different quality levels
witness gif -region demo -o demo.gif -q low   # Smallest files
witness gif -region demo -o demo.gif -q high  # Best quality
```

### Video Recording (Coming Soon)

```bash
# Record as MP4
witness video -region demo -o tutorial.mp4

# High quality recording
witness video -region demo -o tutorial.mp4 -q high
```

### Command Reference

**Selection Commands:**
- `witness select` - Launch interactive region selector
- `witness select -name <name>` - Select and save region
- `witness select -name <name> -default` - Select, save, and set as default

**Region Management:**
- `witness regions` - List all saved regions
- `witness regions -delete <name>` - Delete a saved region
- `witness regions -default <name>` - Set a region as default

**Recording Commands:**
- `witness gif -o <file>` - Record GIF
  - `-region <name>` - Use a saved region
  - `-r <x,y,w,h>` - Use manual coordinates
  - `-f <fps>` - Frames per second (default: 15)
  - `-q <quality>` - Quality level: low, medium, high (default: medium)

## Development

This project uses [Mise](https://mise.jdx.dev/) for task management and tool versioning.

### Available Tasks

View all available tasks:
```bash
mise tasks
```

Common tasks:
```bash
mise run build           # Build the binary
mise run test            # Run all tests
mise run test:verbose    # Run tests with verbose output
mise run test:coverage   # Show test coverage
mise run test:race       # Run tests with race detection
mise run clean           # Remove build artifacts
mise run fmt             # Format code
mise run lint            # Run linter
mise run install         # Install to GOPATH/bin
mise run run             # Build and run
```

Shortcuts:
```bash
mise run b               # Same as 'build'
mise run t               # Same as 'test'
mise run c               # Same as 'clean'
```

### Testing

Comprehensive test suite with >90% coverage on core packages. See [TESTING.md](TESTING.md) for details.

```bash
# Run all tests
mise run test

# Run with coverage report
mise run test:coverage

# Run with race detection
mise run test:race
```

**Alternative**: You can also use Make if you prefer:
```bash
make test
make test-coverage
make build
```

## Architecture

```
witness/
├── cmd/
│   └── witness/          # Main CLI application
├── pkg/
│   ├── capture/          # Screen capture interface
│   ├── encoder/          # GIF and video encoders
│   └── selector/         # Interactive region selection
└── internal/
    └── macos/            # macOS-specific capture implementation
```

### Key Components

- **Capture Package**: Platform-agnostic interface for screen capture
- **Encoder Package**: Handles GIF and video encoding
- **Selector Package**: Interactive region selection and management
- **macOS Package**: Core Graphics integration via CGo

## Technical Details

### macOS Screen Capture

Witness uses Core Graphics APIs for screen capture:
- `CGDisplayCreateImage` for simple single-frame capture
- `CGDisplayStream` (future) for efficient continuous capture

### Region Selection

Interactive region selection leverages macOS's native screenshot tool:
- Uses `screencapture -i` for familiar click-and-drag selection
- Reads selection coordinates from system preferences
- Stores regions in `~/.config/witness/regions.json` for reuse
- Future: Custom overlay using DarwinKit for enhanced UX

### GIF Encoding

Uses Go's standard `image/gif` library with optimizations:
- Floyd-Steinberg dithering for smooth color reduction
- Configurable color palettes (64-256 colors)
- Frame deduplication (planned)

### Video Encoding

Planned integration with:
- `x264-go` for H.264 encoding
- Configurable CRF (Constant Rate Factor) for quality control
- Multiple presets for file size optimization

## Development Status

**Current Version**: 0.1.0-dev

See [PROGRESS.md](PROGRESS.md) for detailed development progress and roadmap.

### Completed
- ✅ Project structure
- ✅ Basic capture interface
- ✅ GIF encoder implementation
- ✅ macOS CGDisplayCreateImage integration
- ✅ Interactive region selection
- ✅ Region persistence and management
- ✅ CLI command parsing
- ✅ Comprehensive test suite with mocking
- ✅ Mise task runner configuration

### In Progress
- 🔄 GIF recording integration (connecting capture + encoder)
- 🔄 Testing on actual macOS system

### Planned
- ⏳ MP4/H.264 encoding
- ⏳ Advanced compression options
- ⏳ Native region selector overlay (using DarwinKit)
- ⏳ Linux support

## Contributing

This is a personal learning project, but suggestions and feedback are welcome!

## License

MIT License - See LICENSE file for details

## Acknowledgments

Built as an alternative to Giphy Capture, inspired by the need for a fast, scriptable screen recording tool.
