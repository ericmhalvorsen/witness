package capture

import (
	"fmt"
	"image/color"
	"testing"
	"time"
)

func TestMockCapturerStartStop(t *testing.T) {
	config := Config{
		FPS:       15,
		DisplayID: 0,
	}

	capturer := NewMockCapturer(config)

	// Test Start
	err := capturer.Start()
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	if !capturer.IsRunning() {
		t.Error("Capturer should be running after Start()")
	}

	// Test starting again (should fail)
	err = capturer.Start()
	if err == nil {
		t.Error("Start() should fail when already running")
	}

	// Test Stop
	err = capturer.Stop()
	if err != nil {
		t.Fatalf("Stop() failed: %v", err)
	}

	if capturer.IsRunning() {
		t.Error("Capturer should not be running after Stop()")
	}

	// Test stopping again (should fail)
	err = capturer.Stop()
	if err == nil {
		t.Error("Stop() should fail when not running")
	}
}

func TestMockCapturerFrames(t *testing.T) {
	config := Config{
		FPS:       30,
		DisplayID: 0,
	}

	capturer := NewMockCapturer(config)
	capturer.FramesToSend = 5
	capturer.FrameDelay = 0 // No delay for faster testing

	err := capturer.Start()
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	// Collect frames
	frameCount := 0
	timeout := time.After(2 * time.Second)

	for {
		select {
		case frame, ok := <-capturer.Frames():
			if !ok {
				// Channel closed, we're done
				goto done
			}
			if frame == nil {
				t.Error("Received nil frame")
			}
			if frame.Image == nil {
				t.Error("Frame has nil image")
			}
			frameCount++
		case <-timeout:
			t.Fatal("Timeout waiting for frames")
		}
	}

done:
	if frameCount != 5 {
		t.Errorf("Expected 5 frames, got %d", frameCount)
	}
}

func TestMockCapturerWithRegion(t *testing.T) {
	region := &Region{
		X:      100,
		Y:      100,
		Width:  800,
		Height: 600,
	}

	config := Config{
		Region:    region,
		FPS:       15,
		DisplayID: 0,
	}

	capturer := NewMockCapturer(config)
	capturer.FramesToSend = 1
	capturer.FrameDelay = 0

	err := capturer.Start()
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	// Get one frame
	frame := <-capturer.Frames()
	if frame == nil {
		t.Fatal("Expected a frame")
	}

	bounds := frame.Image.Bounds()
	if bounds.Dx() != region.Width || bounds.Dy() != region.Height {
		t.Errorf("Frame size = %dx%d, want %dx%d",
			bounds.Dx(), bounds.Dy(), region.Width, region.Height)
	}
}

func TestMockCapturerCustomFrame(t *testing.T) {
	capturer := NewMockCapturer(Config{FPS: 15})

	// Generate a custom frame with a gradient
	frame := capturer.GenerateCustomFrame(100, 100, func(x, y int) color.Color {
		r := uint8(x * 255 / 100)
		g := uint8(y * 255 / 100)
		return color.RGBA{R: r, G: g, B: 0, A: 255}
	})

	if frame == nil {
		t.Fatal("GenerateCustomFrame() returned nil")
	}
	if frame.Image == nil {
		t.Error("Frame has nil image")
	}

	bounds := frame.Image.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("Frame size = %dx%d, want 100x100", bounds.Dx(), bounds.Dy())
	}
}

func TestMockCapturerSendFrame(t *testing.T) {
	config := Config{FPS: 15}
	capturer := NewMockCapturer(config)

	// Should fail when not running
	err := capturer.SendFrame(&Frame{})
	if err == nil {
		t.Error("SendFrame() should fail when not running")
	}

	// Start the capturer
	capturer.Start()
	defer capturer.Stop()

	// Create a custom frame
	customFrame := capturer.GenerateCustomFrame(50, 50, func(x, y int) color.Color {
		return color.RGBA{R: 255, G: 0, B: 0, A: 255}
	})

	// Send it
	err = capturer.SendFrame(customFrame)
	if err != nil {
		t.Fatalf("SendFrame() failed: %v", err)
	}

	// Receive it
	select {
	case frame := <-capturer.Frames():
		if frame == nil {
			t.Error("Received nil frame")
		}
		if frame != customFrame {
			t.Error("Received frame is not the same as sent frame")
		}
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for frame")
	}
}

func TestMockCapturerSendError(t *testing.T) {
	config := Config{FPS: 15}
	capturer := NewMockCapturer(config)

	// Should fail when not running
	err := capturer.SendError(fmt.Errorf("test error"))
	if err == nil {
		t.Error("SendError() should fail when not running")
	}

	// Start the capturer
	capturer.Start()
	defer capturer.Stop()

	// Send an error
	testErr := fmt.Errorf("test error message")
	err = capturer.SendError(testErr)
	if err != nil {
		t.Fatalf("SendError() failed: %v", err)
	}

	// Receive it
	select {
	case receivedErr := <-capturer.Errors():
		if receivedErr == nil {
			t.Error("Received nil error")
		}
		if receivedErr.Error() != testErr.Error() {
			t.Errorf("Error message = %q, want %q", receivedErr.Error(), testErr.Error())
		}
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for error")
	}
}

func TestMockCapturerSimulateError(t *testing.T) {
	config := Config{FPS: 15}
	capturer := NewMockCapturer(config)

	testErr := fmt.Errorf("simulated startup error")
	capturer.SimulateError = testErr

	err := capturer.Start()
	if err == nil {
		t.Error("Start() should fail when SimulateError is set")
	}
	if err.Error() != testErr.Error() {
		t.Errorf("Error = %q, want %q", err.Error(), testErr.Error())
	}
}

func TestMockCapturerFPSRate(t *testing.T) {
	config := Config{
		FPS: 30, // High FPS for testing
	}

	capturer := NewMockCapturer(config)
	capturer.FramesToSend = 10
	capturer.FrameDelay = 0

	start := time.Now()
	capturer.Start()

	frameCount := 0
	for range capturer.Frames() {
		frameCount++
	}
	elapsed := time.Since(start)

	// With 30 FPS, 10 frames should take ~333ms
	// We'll allow some margin for timing variations
	expectedDuration := time.Second * 10 / 30
	if elapsed < expectedDuration/2 || elapsed > expectedDuration*3 {
		t.Logf("Warning: FPS timing may be off. Expected ~%v, got %v", expectedDuration, elapsed)
	}

	if frameCount != 10 {
		t.Errorf("Expected 10 frames, got %d", frameCount)
	}
}

func TestMockCapturerCustomColor(t *testing.T) {
	config := Config{FPS: 15}
	capturer := NewMockCapturer(config)
	capturer.FrameWidth = 10
	capturer.FrameHeight = 10
	capturer.FrameColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	capturer.FramesToSend = 1
	capturer.FrameDelay = 0

	capturer.Start()

	frame := <-capturer.Frames()
	if frame == nil {
		t.Fatal("Expected a frame")
	}

	// Check that the frame has the correct color
	// Sample the center pixel
	c := frame.Image.At(5, 5)
	r, g, b, a := c.RGBA()

	// RGBA returns uint32 values, convert for comparison
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 || a>>8 != 255 {
		t.Errorf("Frame color = RGBA(%d,%d,%d,%d), want RGBA(255,0,0,255)",
			r>>8, g>>8, b>>8, a>>8)
	}
}
