// +build darwin

package capture

import (
	"github.com/ericmhalvorsen/witness/internal/macos"
)

// newPlatformCapturer creates a macOS-specific capturer
func newPlatformCapturer(config Config) (Capturer, error) {
	return macos.NewDisplayCapturer(config)
}
