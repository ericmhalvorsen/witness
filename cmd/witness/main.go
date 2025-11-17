package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Witness - Screen Capture Tool")
	fmt.Println("Version: 0.1.0-dev")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// TODO: Implement command parsing
	fmt.Printf("Command: %s\n", os.Args[1])
}

func printUsage() {
	usage := `
Usage: witness <command> [options]

Commands:
  record    Start recording the screen
  gif       Record and save as GIF
  video     Record and save as MP4
  help      Show this help message

Options:
  -o, --output <file>     Output file path (default: capture_<timestamp>.gif|mp4)
  -r, --region <x,y,w,h>  Capture region (default: full screen)
  -f, --fps <number>      Frames per second (default: 30)
  -q, --quality <level>   Quality level: low, medium, high (default: medium)

Examples:
  witness gif -o demo.gif -f 15
  witness video -o tutorial.mp4 -q high
  witness record -r 0,0,1920,1080 -o capture.gif
`
	fmt.Println(usage)
}
