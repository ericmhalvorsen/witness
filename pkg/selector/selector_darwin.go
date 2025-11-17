// +build darwin

package selector

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ericmhalvorsen/witness/pkg/capture"
)

// macOSSelector uses macOS built-in tools for region selection
type macOSSelector struct {
	config Config
}

// newPlatformSelector creates a macOS selector
func newPlatformSelector() (Selector, error) {
	return &macOSSelector{
		config: DefaultConfig(),
	}, nil
}

// Select launches an interactive region selector
func (s *macOSSelector) Select() (*capture.Region, error) {
	fmt.Println("📐 Select a screen region...")
	fmt.Println("   - Click and drag to select the capture area")
	fmt.Println("   - Press ESC to cancel")
	fmt.Println()

	// Create a temporary file for the screenshot
	// We don't actually need the screenshot, just the selection coordinates
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "witness-selection-tmp.png")
	defer os.Remove(tmpFile) // Clean up

	// Use screencapture with interactive selection
	// -i: interactive mode (click and drag)
	// -x: no sound
	cmd := exec.Command("screencapture", "-i", "-x", tmpFile)

	// Run the command and wait for user selection
	if err := cmd.Run(); err != nil {
		// User likely canceled (ESC)
		return nil, fmt.Errorf("selection canceled")
	}

	// Check if file was created (user completed selection)
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("no region selected")
	}

	// Read the last selection from macOS preferences
	region, err := s.readLastSelection()
	if err != nil {
		return nil, fmt.Errorf("failed to read selection coordinates: %w", err)
	}

	fmt.Printf("✓ Selected region: %dx%d at (%d,%d)\n",
		region.Width, region.Height, region.X, region.Y)

	return region, nil
}

// SelectWithName selects a region and saves it with a name
func (s *macOSSelector) SelectWithName(name string) (*capture.Region, error) {
	region, err := s.Select()
	if err != nil {
		return nil, err
	}

	// Save the region with the name
	if err := SaveRegion(name, region); err != nil {
		return nil, fmt.Errorf("failed to save region: %w", err)
	}

	fmt.Printf("✓ Saved region '%s'\n", name)
	return region, nil
}

// readLastSelection reads the last selection coordinates from macOS preferences
func (s *macOSSelector) readLastSelection() (*capture.Region, error) {
	// Read the screencapture preferences
	cmd := exec.Command("defaults", "read",
		"com.apple.screencapture", "last-selection")

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to read last-selection: %w", err)
	}

	// Parse the output
	// Format is like: {
	//     Height = 480;
	//     Width = 640;
	//     X = 100;
	//     Y = 200;
	// }
	output := out.String()
	region := &capture.Region{}

	// Parse each line
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			valueStr := strings.TrimSpace(strings.TrimSuffix(parts[1], ";"))

			value, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				continue
			}

			switch key {
			case "X":
				region.X = int(value)
			case "Y":
				region.Y = int(value)
			case "Width":
				region.Width = int(value)
			case "Height":
				region.Height = int(value)
			}
		}
	}

	// Validate the region
	if region.Width <= 0 || region.Height <= 0 {
		return nil, fmt.Errorf("invalid region dimensions: %dx%d", region.Width, region.Height)
	}

	return region, nil
}
