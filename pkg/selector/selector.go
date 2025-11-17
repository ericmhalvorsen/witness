package selector

import (
	"fmt"

	"github.com/ericmhalvorsen/witness/pkg/capture"
)

// Selector provides methods for selecting screen regions
type Selector interface {
	// Select launches an interactive region selector and returns the selected region
	Select() (*capture.Region, error)

	// SelectWithName launches selector and saves the region with a name
	SelectWithName(name string) (*capture.Region, error)
}

// NewSelector creates a platform-specific selector
func NewSelector() (Selector, error) {
	return newPlatformSelector()
}

// Config holds selector configuration
type Config struct {
	// Message to display to user during selection
	Message string

	// Whether to show dimensions during selection
	ShowDimensions bool
}

// DefaultConfig returns the default selector configuration
func DefaultConfig() Config {
	return Config{
		Message:        "Select the screen region to capture",
		ShowDimensions: true,
	}
}

// ParseRegionString parses a region string in format "x,y,w,h"
func ParseRegionString(s string) (*capture.Region, error) {
	var x, y, w, h int
	n, err := fmt.Sscanf(s, "%d,%d,%d,%d", &x, &y, &w, &h)
	if err != nil {
		return nil, fmt.Errorf("invalid region format: %w", err)
	}
	if n != 4 {
		return nil, fmt.Errorf("region must have 4 values (x,y,w,h), got %d", n)
	}
	if w <= 0 || h <= 0 {
		return nil, fmt.Errorf("width and height must be positive")
	}

	return &capture.Region{
		X:      x,
		Y:      y,
		Width:  w,
		Height: h,
	}, nil
}

// FormatRegionString converts a region to string format "x,y,w,h"
func FormatRegionString(r *capture.Region) string {
	if r == nil {
		return ""
	}
	return fmt.Sprintf("%d,%d,%d,%d", r.X, r.Y, r.Width, r.Height)
}
