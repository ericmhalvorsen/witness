// +build !darwin

package selector

import "fmt"

// newPlatformSelector returns an error on unsupported platforms
func newPlatformSelector() (Selector, error) {
	return nil, fmt.Errorf("interactive region selection is not supported on this platform (only macOS is currently supported)")
}
