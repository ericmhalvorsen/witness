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

```bash
xcode-select --install
```

### Build from Source

```bash
git clone https://github.com/ericmhalvorsen/witness.git
cd witness
go build -o witness ./cmd/witness
```

## Usage

### Basic GIF Recording

```bash
# Record full screen as GIF
witness gif -o demo.gif

# Record at 15 FPS for smaller file size
witness gif -o demo.gif -f 15

# Record specific region
witness gif -o demo.gif -r 0,0,1920,1080
```

### Video Recording

```bash
# Record as MP4
witness video -o tutorial.mp4

# High quality recording
witness video -o tutorial.mp4 -q high

# Lower quality for smaller files
witness video -o tutorial.mp4 -q low
```

### Options

- `-o, --output <file>` - Output file path (default: auto-generated)
- `-r, --region <x,y,w,h>` - Capture region (default: full screen)
- `-f, --fps <number>` - Frames per second (default: 30)
- `-q, --quality <level>` - Quality: low, medium, high (default: medium)

## Architecture

```
witness/
├── cmd/
│   └── witness/          # Main CLI application
├── pkg/
│   ├── capture/          # Screen capture interface
│   └── encoder/          # GIF and video encoders
└── internal/
    └── macos/            # macOS-specific capture implementation
```

### Key Components

- **Capture Package**: Platform-agnostic interface for screen capture
- **Encoder Package**: Handles GIF and video encoding
- **macOS Package**: Core Graphics integration via CGo

## Technical Details

### macOS Screen Capture

Witness uses Core Graphics APIs for screen capture:
- `CGDisplayCreateImage` for simple single-frame capture
- `CGDisplayStream` (future) for efficient continuous capture

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

### In Progress
- 🔄 CLI command parsing
- 🔄 Testing and validation

### Planned
- ⏳ MP4/H.264 encoding
- ⏳ Region selection UI
- ⏳ Advanced compression options
- ⏳ Linux support

## Contributing

This is a personal learning project, but suggestions and feedback are welcome!

## License

MIT License - See LICENSE file for details

## Acknowledgments

Built as an alternative to Giphy Capture, inspired by the need for a fast, scriptable screen recording tool.
