// +build darwin

// This example demonstrates basic screen capture and GIF encoding
// Note: This only works on macOS
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ericmhalvorsen/witness/pkg/capture"
	"github.com/ericmhalvorsen/witness/pkg/encoder"
)

func main() {
	fmt.Println("Witness - Simple GIF Example")
	fmt.Println("Recording will start in 3 seconds...")
	fmt.Println("Press Ctrl+C to stop recording")

	// Wait 3 seconds before starting
	time.Sleep(3 * time.Second)

	// Configure capture
	config := capture.Config{
		FPS:       15,      // 15 FPS for smaller file size
		DisplayID: 0,       // Main display
		Region:    nil,     // Full screen
	}

	// Create capturer
	capturer, err := capture.NewCapturer(config)
	if err != nil {
		log.Fatalf("Failed to create capturer: %v", err)
	}

	// Create GIF encoder
	outputPath := "output.gif"
	gifEncoder := encoder.NewGIFEncoder(outputPath, config.FPS, encoder.QualityMedium)

	// Start capture
	fmt.Println("Starting capture...")
	if err := capturer.Start(); err != nil {
		log.Fatalf("Failed to start capture: %v", err)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Capture loop
	frameCount := 0
	maxFrames := 300 // Maximum 20 seconds at 15 FPS

	go func() {
		for frame := range capturer.Frames() {
			if err := gifEncoder.AddFrame(frame); err != nil {
				log.Printf("Failed to add frame: %v", err)
				continue
			}

			frameCount++
			if frameCount%15 == 0 {
				fmt.Printf("Captured %d frames (%.1f seconds)\n",
					frameCount, float64(frameCount)/float64(config.FPS))
			}

			if frameCount >= maxFrames {
				fmt.Println("Maximum frame count reached")
				sigChan <- os.Interrupt
				break
			}
		}
	}()

	// Handle errors
	go func() {
		for err := range capturer.Errors() {
			log.Printf("Capture error: %v", err)
		}
	}()

	// Wait for interrupt
	<-sigChan

	// Stop capture
	fmt.Println("\nStopping capture...")
	if err := capturer.Stop(); err != nil {
		log.Printf("Error stopping capture: %v", err)
	}

	// Encode GIF
	fmt.Printf("Encoding %d frames to GIF...\n", gifEncoder.FrameCount())
	if err := gifEncoder.Encode(); err != nil {
		log.Fatalf("Failed to encode GIF: %v", err)
	}

	// Get file size
	info, err := os.Stat(outputPath)
	if err != nil {
		log.Printf("Warning: Could not stat output file: %v", err)
	} else {
		fmt.Printf("GIF saved to %s (%.2f MB)\n", outputPath, float64(info.Size())/1024/1024)
	}

	fmt.Println("Done!")
}
