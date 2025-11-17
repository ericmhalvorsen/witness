// +build !darwin

package capture

import "fmt"

// newPlatformCapturer returns an error on unsupported platforms
func newPlatformCapturer(config Config) (Capturer, error) {
	return nil, fmt.Errorf("screen capture is not supported on this platform (only macOS is currently supported)")
}
