package encoder

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"

	"github.com/ericmhalvorsen/witness/pkg/capture"
)

// GIFQuality defines the quality level for GIF encoding
type GIFQuality int

const (
	// QualityLow uses aggressive palette reduction for smallest files
	QualityLow GIFQuality = iota
	// QualityMedium balances file size and visual quality
	QualityMedium
	// QualityHigh preserves more colors for better quality
	QualityHigh
)

// GIFEncoder encodes captured frames as an animated GIF
type GIFEncoder struct {
	quality    GIFQuality
	delay      int  // Delay between frames in 100ths of a second
	outputPath string
	frames     []*image.Paletted
	delays     []int
}

// NewGIFEncoder creates a new GIF encoder
func NewGIFEncoder(outputPath string, fps int, quality GIFQuality) *GIFEncoder {
	// Convert FPS to delay (in 100ths of a second)
	// delay = 100 / fps
	delay := 100 / fps
	if delay < 1 {
		delay = 1 // Minimum delay
	}

	return &GIFEncoder{
		quality:    quality,
		delay:      delay,
		outputPath: outputPath,
		frames:     make([]*image.Paletted, 0),
		delays:     make([]int, 0),
	}
}

// AddFrame adds a frame to the GIF
func (e *GIFEncoder) AddFrame(frame *capture.Frame) error {
	if frame == nil || frame.Image == nil {
		return fmt.Errorf("invalid frame")
	}

	// Convert RGBA to Paletted image
	palettedImg := e.convertToPaletted(frame.Image)

	e.frames = append(e.frames, palettedImg)
	e.delays = append(e.delays, e.delay)

	return nil
}

// Encode writes all frames to the output file as an animated GIF
func (e *GIFEncoder) Encode() error {
	if len(e.frames) == 0 {
		return fmt.Errorf("no frames to encode")
	}

	// Create output file
	outFile, err := os.Create(e.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Create GIF
	anim := &gif.GIF{
		Image: e.frames,
		Delay: e.delays,
	}

	// Encode to file
	if err := gif.EncodeAll(outFile, anim); err != nil {
		return fmt.Errorf("failed to encode GIF: %w", err)
	}

	return nil
}

// FrameCount returns the number of frames currently buffered
func (e *GIFEncoder) FrameCount() int {
	return len(e.frames)
}

// convertToPaletted converts an RGBA image to a paletted image
func (e *GIFEncoder) convertToPaletted(img *image.RGBA) *image.Paletted {
	bounds := img.Bounds()
	palettedImg := image.NewPaletted(bounds, e.getPalette())

	// Draw the RGBA image onto the paletted image
	// This will automatically handle color quantization
	draw.FloydSteinberg.Draw(palettedImg, bounds, img, image.Point{})

	return palettedImg
}

// getPalette returns the color palette based on quality setting
func (e *GIFEncoder) getPalette() color.Palette {
	switch e.quality {
	case QualityLow:
		// Use a reduced palette (64 colors) for smaller file size
		return palette.Plan9[:64]
	case QualityMedium:
		// Use Plan9 palette (256 colors)
		return palette.Plan9
	case QualityHigh:
		// Use WebSafe palette (216 colors) - better color accuracy
		return palette.WebSafe
	default:
		return palette.Plan9
	}
}

// EstimateSize provides a rough estimate of the output file size
func (e *GIFEncoder) EstimateSize() int64 {
	if len(e.frames) == 0 {
		return 0
	}

	// Rough estimate: header + (frame_size * num_frames)
	// This is very approximate
	frameSize := e.frames[0].Bounds().Dx() * e.frames[0].Bounds().Dy()
	estimatedSize := int64(frameSize * len(e.frames) / 4) // GIF compression ~4x

	return estimatedSize
}
