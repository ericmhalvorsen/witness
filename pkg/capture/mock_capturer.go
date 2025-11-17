package capture

import (
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"
)

// MockCapturer is a mock implementation of the Capturer interface for testing
type MockCapturer struct {
	config    Config
	frames    chan *Frame
	errors    chan error
	stopChan  chan struct{}
	isRunning bool
	mu        sync.Mutex

	// Configuration options for the mock
	FrameWidth     int
	FrameHeight    int
	FrameColor     color.Color
	FramesToSend   int
	SimulateError  error
	FrameDelay     time.Duration
}

// NewMockCapturer creates a new mock capturer for testing
func NewMockCapturer(config Config) *MockCapturer {
	return &MockCapturer{
		config:       config,
		frames:       make(chan *Frame, 10),
		errors:       make(chan error, 10),
		stopChan:     make(chan struct{}),
		FrameWidth:   640,
		FrameHeight:  480,
		FrameColor:   color.RGBA{R: 128, G: 128, B: 128, A: 255},
		FramesToSend: -1, // -1 means infinite
		FrameDelay:   time.Millisecond * 10,
	}
}

// Start begins the mock capture process
func (m *MockCapturer) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isRunning {
		return fmt.Errorf("capturer already running")
	}

	// Simulate an error if configured
	if m.SimulateError != nil {
		return m.SimulateError
	}

	m.isRunning = true
	go m.captureLoop()

	return nil
}

// Stop ends the mock capture process
func (m *MockCapturer) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return fmt.Errorf("capturer not running")
	}

	close(m.stopChan)
	m.isRunning = false

	return nil
}

// Frames returns the channel for captured frames
func (m *MockCapturer) Frames() <-chan *Frame {
	return m.frames
}

// Errors returns the channel for errors
func (m *MockCapturer) Errors() <-chan error {
	return m.errors
}

// IsRunning returns whether the capturer is currently running
func (m *MockCapturer) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isRunning
}

// captureLoop generates mock frames at the configured FPS
func (m *MockCapturer) captureLoop() {
	ticker := time.NewTicker(time.Second / time.Duration(m.config.FPS))
	defer ticker.Stop()
	defer close(m.frames)
	defer close(m.errors)

	frameCount := 0

	for {
		select {
		case <-m.stopChan:
			return
		case <-ticker.C:
			// Check if we've sent enough frames
			if m.FramesToSend >= 0 && frameCount >= m.FramesToSend {
				return
			}

			// Apply frame delay if configured
			if m.FrameDelay > 0 {
				time.Sleep(m.FrameDelay)
			}

			// Generate a mock frame
			frame := m.generateFrame()
			m.frames <- frame
			frameCount++
		}
	}
}

// generateFrame creates a mock frame with the configured properties
func (m *MockCapturer) generateFrame() *Frame {
	width := m.FrameWidth
	height := m.FrameHeight

	// Use region dimensions if specified
	if m.config.Region != nil {
		width = m.config.Region.Width
		height = m.config.Region.Height
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with the configured color
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, m.FrameColor)
		}
	}

	return &Frame{
		Image:     img,
		Timestamp: time.Now(),
	}
}

// GenerateCustomFrame allows creating a custom frame for testing
func (m *MockCapturer) GenerateCustomFrame(width, height int, fillFunc func(x, y int) color.Color) *Frame {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, fillFunc(x, y))
		}
	}

	return &Frame{
		Image:     img,
		Timestamp: time.Now(),
	}
}

// SendFrame manually sends a frame to the frames channel (useful for controlled testing)
func (m *MockCapturer) SendFrame(frame *Frame) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return fmt.Errorf("capturer not running")
	}

	select {
	case m.frames <- frame:
		return nil
	case <-time.After(time.Second):
		return fmt.Errorf("timeout sending frame")
	}
}

// SendError manually sends an error to the errors channel
func (m *MockCapturer) SendError(err error) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return fmt.Errorf("capturer not running")
	}

	select {
	case m.errors <- err:
		return nil
	case <-time.After(time.Second):
		return fmt.Errorf("timeout sending error")
	}
}
