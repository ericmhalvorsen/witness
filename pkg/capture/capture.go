package capture

import (
	"image"
	"time"
)

// Region defines a rectangular area to capture
type Region struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Config holds configuration for screen capture
type Config struct {
	// Region to capture. If nil, captures full screen
	Region *Region

	// Target frames per second
	FPS int

	// Display ID (for multi-monitor setups). 0 for main display
	DisplayID uint32
}

// Frame represents a single captured frame
type Frame struct {
	Image     *image.RGBA
	Timestamp time.Time
}

// Capturer is the interface for screen capture implementations
type Capturer interface {
	// Start begins the capture process
	Start() error

	// Stop ends the capture process
	Stop() error

	// Frames returns a channel that receives captured frames
	Frames() <-chan *Frame

	// Errors returns a channel for capture errors
	Errors() <-chan error
}

// NewCapturer creates a platform-specific capturer
// This will be implemented per platform (macOS, Linux, etc.)
func NewCapturer(config Config) (Capturer, error) {
	// Platform-specific implementation will be called here
	return newPlatformCapturer(config)
}
