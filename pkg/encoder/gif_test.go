package encoder

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ericmhalvorsen/witness/pkg/capture"
)

// Helper function to create a test frame with a solid color
func createTestFrame(width, height int, c color.Color) *capture.Frame {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with color
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}

	return &capture.Frame{
		Image:     img,
		Timestamp: time.Now(),
	}
}

// Helper function to create a test frame with a gradient pattern
func createGradientFrame(width, height int) *capture.Frame {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create a gradient based on position
			r := uint8(x * 255 / width)
			g := uint8(y * 255 / height)
			b := uint8((x + y) * 255 / (width + height))
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	return &capture.Frame{
		Image:     img,
		Timestamp: time.Now(),
	}
}

func TestNewGIFEncoder(t *testing.T) {
	tests := []struct {
		name        string
		fps         int
		quality     GIFQuality
		wantDelay   int
	}{
		{
			name:      "15 FPS medium quality",
			fps:       15,
			quality:   QualityMedium,
			wantDelay: 6, // 100/15 = 6.66... rounds to 6
		},
		{
			name:      "10 FPS low quality",
			fps:       10,
			quality:   QualityLow,
			wantDelay: 10,
		},
		{
			name:      "30 FPS high quality",
			fps:       30,
			quality:   QualityHigh,
			wantDelay: 3, // 100/30 = 3.33... rounds to 3
		},
		{
			name:      "very high FPS caps at minimum delay",
			fps:       200,
			quality:   QualityMedium,
			wantDelay: 1, // Minimum delay is 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := NewGIFEncoder("test.gif", tt.fps, tt.quality)

			if encoder == nil {
				t.Fatal("NewGIFEncoder() returned nil")
			}
			if encoder.quality != tt.quality {
				t.Errorf("quality = %v, want %v", encoder.quality, tt.quality)
			}
			if encoder.delay != tt.wantDelay {
				t.Errorf("delay = %v, want %v", encoder.delay, tt.wantDelay)
			}
			if encoder.outputPath != "test.gif" {
				t.Errorf("outputPath = %v, want %v", encoder.outputPath, "test.gif")
			}
		})
	}
}

func TestAddFrame(t *testing.T) {
	encoder := NewGIFEncoder("test.gif", 15, QualityMedium)

	// Add a valid frame
	frame := createTestFrame(100, 100, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	err := encoder.AddFrame(frame)
	if err != nil {
		t.Errorf("AddFrame() failed for valid frame: %v", err)
	}

	if encoder.FrameCount() != 1 {
		t.Errorf("FrameCount() = %d, want 1", encoder.FrameCount())
	}

	// Add nil frame should error
	err = encoder.AddFrame(nil)
	if err == nil {
		t.Error("AddFrame() should fail for nil frame")
	}

	// Add frame with nil image should error
	badFrame := &capture.Frame{Image: nil, Timestamp: time.Now()}
	err = encoder.AddFrame(badFrame)
	if err == nil {
		t.Error("AddFrame() should fail for frame with nil image")
	}
}

func TestAddMultipleFrames(t *testing.T) {
	encoder := NewGIFEncoder("test.gif", 15, QualityMedium)

	colors := []color.Color{
		color.RGBA{R: 255, G: 0, B: 0, A: 255},   // Red
		color.RGBA{R: 0, G: 255, B: 0, A: 255},   // Green
		color.RGBA{R: 0, G: 0, B: 255, A: 255},   // Blue
		color.RGBA{R: 255, G: 255, B: 0, A: 255}, // Yellow
	}

	for i, c := range colors {
		frame := createTestFrame(100, 100, c)
		err := encoder.AddFrame(frame)
		if err != nil {
			t.Fatalf("AddFrame() failed for frame %d: %v", i, err)
		}
	}

	if encoder.FrameCount() != len(colors) {
		t.Errorf("FrameCount() = %d, want %d", encoder.FrameCount(), len(colors))
	}
}

func TestEncode(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "witness-encoder-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	outputPath := filepath.Join(tmpDir, "test.gif")
	encoder := NewGIFEncoder(outputPath, 15, QualityMedium)

	// Add some frames
	for i := 0; i < 5; i++ {
		frame := createTestFrame(100, 100, color.RGBA{
			R: uint8(i * 50),
			G: 100,
			B: 200,
			A: 255,
		})
		encoder.AddFrame(frame)
	}

	// Encode to file
	err = encoder.Encode()
	if err != nil {
		t.Fatalf("Encode() failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	// Verify file is not empty
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Failed to stat output file: %v", err)
	}
	if info.Size() == 0 {
		t.Error("Output file is empty")
	}
}

func TestEncodeNoFrames(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "witness-encoder-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	outputPath := filepath.Join(tmpDir, "empty.gif")
	encoder := NewGIFEncoder(outputPath, 15, QualityMedium)

	// Try to encode without adding frames
	err = encoder.Encode()
	if err == nil {
		t.Error("Encode() should fail when no frames have been added")
	}
}

func TestQualityLevels(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "witness-encoder-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	qualities := []struct {
		name    string
		quality GIFQuality
	}{
		{"low", QualityLow},
		{"medium", QualityMedium},
		{"high", QualityHigh},
	}

	// Create a test frame with gradient (to test color palette)
	frame := createGradientFrame(200, 200)

	for _, q := range qualities {
		t.Run(q.name, func(t *testing.T) {
			outputPath := filepath.Join(tmpDir, q.name+".gif")
			encoder := NewGIFEncoder(outputPath, 15, q.quality)

			// Add the same frame multiple times
			for i := 0; i < 3; i++ {
				if err := encoder.AddFrame(frame); err != nil {
					t.Fatalf("AddFrame() failed: %v", err)
				}
			}

			// Encode
			if err := encoder.Encode(); err != nil {
				t.Fatalf("Encode() failed: %v", err)
			}

			// Verify file exists
			if _, err := os.Stat(outputPath); os.IsNotExist(err) {
				t.Errorf("Output file not created for %s quality", q.name)
			}
		})
	}
}

func TestGetPalette(t *testing.T) {
	tests := []struct {
		quality     GIFQuality
		minColors   int
		description string
	}{
		{QualityLow, 64, "low quality should have 64 colors"},
		{QualityMedium, 256, "medium quality should have 256 colors"},
		{QualityHigh, 216, "high quality should have 216 colors (WebSafe)"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			encoder := NewGIFEncoder("test.gif", 15, tt.quality)
			palette := encoder.getPalette()

			if len(palette) < tt.minColors {
				t.Errorf("palette size = %d, want at least %d", len(palette), tt.minColors)
			}
		})
	}
}

func TestEstimateSize(t *testing.T) {
	encoder := NewGIFEncoder("test.gif", 15, QualityMedium)

	// Should be 0 for no frames
	if size := encoder.EstimateSize(); size != 0 {
		t.Errorf("EstimateSize() = %d for empty encoder, want 0", size)
	}

	// Add a frame
	frame := createTestFrame(100, 100, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	encoder.AddFrame(frame)

	// Should have a non-zero estimate now
	if size := encoder.EstimateSize(); size <= 0 {
		t.Error("EstimateSize() should be positive after adding frames")
	}

	// Add more frames, size should increase
	firstEstimate := encoder.EstimateSize()
	encoder.AddFrame(frame)
	encoder.AddFrame(frame)
	secondEstimate := encoder.EstimateSize()

	if secondEstimate <= firstEstimate {
		t.Error("EstimateSize() should increase when adding more frames")
	}
}

func TestDifferentFrameSizes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "witness-encoder-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	sizes := []struct {
		width  int
		height int
	}{
		{100, 100},
		{200, 150},
		{400, 300},
		{800, 600},
	}

	for _, size := range sizes {
		t.Run("", func(t *testing.T) {
			outputPath := filepath.Join(tmpDir, "size_test.gif")
			encoder := NewGIFEncoder(outputPath, 15, QualityMedium)

			frame := createTestFrame(size.width, size.height, color.RGBA{R: 128, G: 128, B: 128, A: 255})
			if err := encoder.AddFrame(frame); err != nil {
				t.Fatalf("AddFrame() failed for %dx%d: %v", size.width, size.height, err)
			}

			if err := encoder.Encode(); err != nil {
				t.Fatalf("Encode() failed for %dx%d: %v", size.width, size.height, err)
			}

			// Clean up for next iteration
			os.Remove(outputPath)
		})
	}
}

func TestFrameCount(t *testing.T) {
	encoder := NewGIFEncoder("test.gif", 15, QualityMedium)

	if count := encoder.FrameCount(); count != 0 {
		t.Errorf("Initial FrameCount() = %d, want 0", count)
	}

	frame := createTestFrame(100, 100, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	for i := 1; i <= 10; i++ {
		encoder.AddFrame(frame)
		if count := encoder.FrameCount(); count != i {
			t.Errorf("FrameCount() after %d additions = %d, want %d", i, count, i)
		}
	}
}

func TestConvertToPaletted(t *testing.T) {
	encoder := NewGIFEncoder("test.gif", 15, QualityMedium)

	// Create a test RGBA image
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}

	// Convert to paletted
	paletted := encoder.convertToPaletted(img)

	if paletted == nil {
		t.Fatal("convertToPaletted() returned nil")
	}

	// Check dimensions are preserved
	bounds := paletted.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 50 {
		t.Errorf("convertToPaletted() changed dimensions: got %dx%d, want 50x50",
			bounds.Dx(), bounds.Dy())
	}

	// Check that it uses the correct palette
	if paletted.Palette == nil {
		t.Error("convertToPaletted() produced image with nil palette")
	}
}

func TestEncodeInvalidPath(t *testing.T) {
	// Try to write to an invalid path
	encoder := NewGIFEncoder("/invalid/path/that/does/not/exist/test.gif", 15, QualityMedium)

	frame := createTestFrame(100, 100, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	encoder.AddFrame(frame)

	err := encoder.Encode()
	if err == nil {
		t.Error("Encode() should fail for invalid output path")
	}
}
