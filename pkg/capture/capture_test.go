package capture

import (
	"image"
	"testing"
	"time"
)

func TestRegion(t *testing.T) {
	tests := []struct {
		name   string
		region Region
	}{
		{
			name: "standard region",
			region: Region{
				X:      100,
				Y:      200,
				Width:  800,
				Height: 600,
			},
		},
		{
			name: "region at origin",
			region: Region{
				X:      0,
				Y:      0,
				Width:  1920,
				Height: 1080,
			},
		},
		{
			name: "small region",
			region: Region{
				X:      50,
				Y:      50,
				Width:  100,
				Height: 100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.region

			if r.X < 0 || r.Y < 0 {
				t.Error("Region coordinates should not be negative")
			}
			if r.Width <= 0 || r.Height <= 0 {
				t.Error("Region dimensions should be positive")
			}
		})
	}
}

func TestConfig(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "default config",
			config: Config{
				Region:    nil,
				FPS:       15,
				DisplayID: 0,
			},
		},
		{
			name: "config with region",
			config: Config{
				Region: &Region{
					X:      100,
					Y:      100,
					Width:  800,
					Height: 600,
				},
				FPS:       30,
				DisplayID: 0,
			},
		},
		{
			name: "config with secondary display",
			config: Config{
				Region:    nil,
				FPS:       24,
				DisplayID: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.config

			if c.FPS <= 0 {
				t.Error("FPS should be positive")
			}
			if c.Region != nil {
				if c.Region.Width <= 0 || c.Region.Height <= 0 {
					t.Error("Region dimensions should be positive")
				}
			}
		})
	}
}

func TestFrame(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	timestamp := time.Now()

	frame := Frame{
		Image:     img,
		Timestamp: timestamp,
	}

	if frame.Image == nil {
		t.Error("Frame.Image should not be nil")
	}
	if frame.Timestamp.IsZero() {
		t.Error("Frame.Timestamp should not be zero")
	}
	if !frame.Timestamp.Equal(timestamp) {
		t.Error("Frame.Timestamp should match the set timestamp")
	}
}

func TestFrameTimestamp(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	start := time.Now()
	frame := Frame{
		Image:     img,
		Timestamp: start,
	}
	end := time.Now()

	if frame.Timestamp.Before(start) || frame.Timestamp.After(end) {
		t.Error("Frame timestamp should be between start and end times")
	}
}

func TestConfigWithNilRegion(t *testing.T) {
	config := Config{
		Region:    nil, // Full screen capture
		FPS:       15,
		DisplayID: 0,
	}

	if config.Region != nil {
		t.Error("Region should be nil for full screen capture")
	}
	if config.FPS != 15 {
		t.Errorf("FPS = %d, want 15", config.FPS)
	}
	if config.DisplayID != 0 {
		t.Errorf("DisplayID = %d, want 0", config.DisplayID)
	}
}

func TestConfigWithRegion(t *testing.T) {
	region := &Region{
		X:      100,
		Y:      200,
		Width:  800,
		Height: 600,
	}

	config := Config{
		Region:    region,
		FPS:       30,
		DisplayID: 0,
	}

	if config.Region == nil {
		t.Fatal("Region should not be nil")
	}
	if config.Region.X != 100 || config.Region.Y != 200 {
		t.Errorf("Region position = (%d,%d), want (100,200)", config.Region.X, config.Region.Y)
	}
	if config.Region.Width != 800 || config.Region.Height != 600 {
		t.Errorf("Region size = (%d,%d), want (800,600)", config.Region.Width, config.Region.Height)
	}
}

func TestMultipleFramesSequence(t *testing.T) {
	frames := []*Frame{}
	lastTimestamp := time.Time{}

	for i := 0; i < 5; i++ {
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		frame := &Frame{
			Image:     img,
			Timestamp: time.Now(),
		}
		frames = append(frames, frame)
		time.Sleep(1 * time.Millisecond) // Ensure timestamps are different
	}

	// Verify timestamps are in order
	for i, frame := range frames {
		if frame.Timestamp.Before(lastTimestamp) {
			t.Errorf("Frame %d timestamp is before previous frame", i)
		}
		lastTimestamp = frame.Timestamp
	}
}
