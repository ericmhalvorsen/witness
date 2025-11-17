# Screen Capture Tool - Progress Log

## Project Overview
Building a lightweight, efficient screen capture tool for macOS that can record screen content and save as GIF or MP4 files. Goal is to replace clunky tools like Giphy Capture with a command-line driven solution optimized for creating small, efficient files for conveying concepts.

**Tech Stack:** Go (Golang)
**Target Platform:** macOS (with potential Linux support later)

---

## Research Findings

### macOS Screen Capture APIs

#### ScreenCaptureKit (Recommended - Modern)
- Apple's newest, high-performance screen capture framework
- Built from the ground up to replace older APIs
- Superior performance with GPU-accelerated color conversion and scaling
- **Requirement:** macOS 12.3+ (Monterey)
- **Status:** Best option for modern macOS systems

#### CGDisplayStream (Stable Alternative)
- Lower-level API, promoted in WWDC 2012
- GPU-based color conversion and scaling
- Options for cursor capture and minimum frame duration
- **Known Issues:** Hanging issues reported on macOS 15 when multiple processes use it
- **Status:** Good fallback for older systems

#### AVCaptureScreenInput (Legacy)
- Part of AVFoundation framework
- Still functional but being phased out
- **Status:** Not recommended for new projects

### Go Libraries for Screen Capture

#### kbinani/screenshot
- Cross-platform Go library for desktop screenshots
- Supports macOS, Windows, Linux, FreeBSD, OpenBSD, NetBSD
- Good for static screenshots, would need extension for video capture
- **GitHub:** https://github.com/kbinani/screenshot

#### go-scrap
- Go wrapper around Rust scrap library
- Cross-platform with reasonable performance
- **GitHub:** https://github.com/cretz/go-scrap

**Decision:** Will likely need to write CGo bindings to macOS APIs directly for optimal performance and control.

### Video Encoding Libraries

#### For MP4/H.264:

1. **u2takey/ffmpeg-go** (Recommended for simplicity)
   - High-level API (port of ffmpeg-python)
   - Requires FFmpeg installed on system
   - Clean, idiomatic Go API
   - **GitHub:** https://github.com/u2takey/ffmpeg-go

2. **asticode/go-astiav**
   - Comprehensive FFmpeg C bindings
   - Compatible with FFmpeg n7.0
   - More control but more complex
   - **GitHub:** https://github.com/asticode/go-astiav

3. **gen2brain/x264-go**
   - Direct x264 encoder bindings
   - Includes C source code (easier install)
   - Good for H.264 specifically
   - **GitHub:** https://github.com/gen2brain/x264-go

**Decision:** Start with x264-go for direct control, fall back to ffmpeg-go if needed.

#### For GIF:

- **image/gif** (Go Standard Library)
- Built-in, no dependencies!
- `gif.EncodeAll()` for animated GIFs
- Create frames as `[]*image.Paletted`
- Configure delay times between frames
- **Decision:** Use standard library - perfect for our needs

---

## Architecture Plan

### Phase 1: Basic Screen Capture (Current)
- [x] Initialize Go module and project structure
- [x] Create CGo bindings for CGDisplayStream or ScreenCaptureKit
- [x] Implement continuous frame capture
- [ ] Create selection UI for capture area (or start with full screen)
- [ ] Add start/stop recording controls
- [ ] Test on actual macOS system

### Phase 2: GIF Encoding
- [x] Capture frames to `image.Image` format
- [x] Implement frame rate control
- [x] Convert frames to paletted images
- [x] Use `image/gif` to encode to GIF
- [x] Add basic compression controls (color palette reduction)

### Phase 3: MP4/Video Encoding
- [ ] Integrate x264-go for H.264 encoding
- [ ] Implement frame buffer to encoder pipeline
- [ ] Add compression level controls
- [ ] Optimize for file size (adjust bitrate, CRF values)

### Phase 4: Optimization & Polish
- [ ] Add various compression presets (high quality, balanced, maximum compression)
- [ ] Implement smart color palette generation for GIFs
- [ ] Add progress indicators
- [ ] Memory optimization for long recordings
- [ ] Error handling and recovery

### Phase 5: Future Enhancements
- [ ] Linux support
- [ ] Audio capture
- [ ] Real-time preview
- [ ] Hotkey support for start/stop
- [ ] Configurable frame rates

---

## Technical Challenges & Considerations

### Challenge 1: Bridging macOS APIs with Go
- **Solution:** Use CGo to create bindings for Core Graphics/ScreenCaptureKit
- **Trade-off:** CGo introduces complexity but necessary for system APIs

### Challenge 2: File Size Optimization
- **For GIF:** Palette reduction, frame deduplication, dithering
- **For MP4:** CRF (Constant Rate Factor), preset tuning, resolution options
- **Strategy:** Provide presets + advanced options for power users

### Challenge 3: Performance
- **Concern:** Capturing, encoding in real-time without dropped frames
- **Solution:** Buffered pipeline, goroutines for parallel processing
- **Target:** 30 FPS for smooth captures, configurable down to 10 FPS for smaller files

### Challenge 4: User Experience
- **CLI Interface:** Need intuitive commands
- **Feedback:** Progress bars, estimated file size, frame count
- **Interrupts:** Graceful handling of Ctrl+C, save partial recordings

---

## Dependencies to Install

```bash
# For development
go get github.com/gen2brain/x264-go  # H.264 encoding
# image/gif is standard library - no install needed

# System requirements (macOS)
# Xcode Command Line Tools (for CGo)
xcode-select --install
```

---

## Next Steps
1. Initialize Go module
2. Create basic project structure (cmd/, pkg/, internal/)
3. Start with simple CGDisplayStream capture example
4. Test frame capture and display rate
5. Implement GIF encoding first (simpler than video)

---

## Progress Log

### 2025-11-17

#### Research & Planning
- ✅ Researched macOS screen capture APIs (ScreenCaptureKit, CGDisplayStream, AVCaptureScreenInput)
- ✅ Evaluated Go libraries for video encoding (x264-go, ffmpeg-go, go-astiav)
- ✅ Confirmed Go's built-in image/gif package for GIF encoding
- ✅ Created PROGRESS.md with findings and architecture plan

#### Project Setup
- ✅ Initialized Go module (github.com/ericmhalvorsen/witness)
- ✅ Created project structure (cmd/, pkg/, internal/, examples/)
- ✅ Added .gitignore and Makefile
- ✅ Built successfully on Linux (with platform stubs)

#### Core Implementation
- ✅ Implemented capture package with platform-agnostic interface
- ✅ Created macOS capture implementation using CGDisplayCreateImage via CGo
- ✅ Implemented GIF encoder with quality levels (low, medium, high)
- ✅ Added Floyd-Steinberg dithering for smooth color reduction
- ✅ Created CLI skeleton with command structure

#### Documentation
- ✅ Created comprehensive README.md
- ✅ Added example program (examples/simple_gif.go)
- ✅ Documented architecture and usage

#### Next Steps
- 🔄 Test actual screen capture on macOS system
- 🔄 Implement CLI command parsing and flag handling
- 🔄 Add region selection functionality
- 🔄 Test GIF output quality and file sizes
- 🔄 Begin MP4/H.264 encoding integration
