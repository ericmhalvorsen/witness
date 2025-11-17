package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ericmhalvorsen/witness/pkg/capture"
	"github.com/ericmhalvorsen/witness/pkg/selector"
)

const version = "0.1.0-dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "select":
		handleSelect(os.Args[2:])
	case "regions":
		handleRegions(os.Args[2:])
	case "gif":
		handleGif(os.Args[2:])
	case "video":
		handleVideo(os.Args[2:])
	case "help", "--help", "-h":
		printUsage()
	case "version", "--version", "-v":
		fmt.Printf("Witness version %s\n", version)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handleSelect(args []string) {
	fs := flag.NewFlagSet("select", flag.ExitOnError)
	name := fs.String("name", "", "Save the selected region with a name")
	setDefault := fs.Bool("default", false, "Set this region as the default")

	fs.Usage = func() {
		fmt.Println("Usage: witness select [options]")
		fmt.Println("\nLaunch an interactive region selector")
		fmt.Println("\nOptions:")
		fs.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  witness select                    # Select a region")
		fmt.Println("  witness select -name demo         # Select and save as 'demo'")
		fmt.Println("  witness select -name demo -default # Select, save, and set as default")
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	// Create selector
	sel, err := selector.NewSelector()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Select region
	var region *capture.Region
	if *name != "" {
		region, err = sel.SelectWithName(*name)
	} else {
		region, err = sel.Select()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Set as default if requested
	if *setDefault && *name != "" {
		if err := selector.SetDefaultRegion(*name); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to set default region: %v\n", err)
		} else {
			fmt.Printf("✓ Set '%s' as default region\n", *name)
		}
	}

	// Print region info
	if *name == "" {
		fmt.Println("\nTo use this region in capture:")
		fmt.Printf("  witness gif -r %s\n", selector.FormatRegionString(region))
		fmt.Println("\nOr save it for later use:")
		fmt.Printf("  witness select -name myregion\n")
	}
}

func handleRegions(args []string) {
	fs := flag.NewFlagSet("regions", flag.ExitOnError)
	delete := fs.String("delete", "", "Delete a saved region")
	setDefault := fs.String("default", "", "Set a region as default")

	fs.Usage = func() {
		fmt.Println("Usage: witness regions [options]")
		fmt.Println("\nManage saved screen regions")
		fmt.Println("\nOptions:")
		fs.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  witness regions                    # List all saved regions")
		fmt.Println("  witness regions -delete demo       # Delete 'demo' region")
		fmt.Println("  witness regions -default demo      # Set 'demo' as default")
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	// Handle delete
	if *delete != "" {
		if err := selector.DeleteRegion(*delete); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Deleted region '%s'\n", *delete)
		return
	}

	// Handle set default
	if *setDefault != "" {
		if err := selector.SetDefaultRegion(*setDefault); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Set '%s' as default region\n", *setDefault)
		return
	}

	// Handle list (default action)
	names, err := selector.ListRegions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(names) == 0 {
		fmt.Println("No saved regions")
		fmt.Println("\nCreate one with: witness select -name myregion")
		return
	}

	fmt.Println("Saved regions:")
	for _, name := range names {
		info, err := selector.GetRegionInfo(name)
		if err != nil {
			continue
		}
		fmt.Printf("  %s\n", info)
	}
}

func handleGif(args []string) {
	fs := flag.NewFlagSet("gif", flag.ExitOnError)
	output := fs.String("o", "", "Output file path")
	regionStr := fs.String("r", "", "Capture region (x,y,w,h)")
	regionName := fs.String("region", "", "Use a saved region by name")
	fps := fs.Int("f", 15, "Frames per second")
	quality := fs.String("q", "medium", "Quality level (low, medium, high)")

	fs.Usage = func() {
		fmt.Println("Usage: witness gif [options]")
		fmt.Println("\nRecord screen and save as GIF")
		fmt.Println("\nOptions:")
		fs.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  witness gif -o demo.gif")
		fmt.Println("  witness gif -o demo.gif -f 10 -q low")
		fmt.Println("  witness gif -region demo -o capture.gif")
		fmt.Println("  witness gif -r 0,0,800,600 -o capture.gif")
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	// TODO: Implement GIF recording
	fmt.Println("GIF recording not yet implemented")
	fmt.Printf("Output: %s\n", *output)
	fmt.Printf("Region: %s\n", *regionStr)
	fmt.Printf("Region name: %s\n", *regionName)
	fmt.Printf("FPS: %d\n", *fps)
	fmt.Printf("Quality: %s\n", *quality)
}

func handleVideo(args []string) {
	fs := flag.NewFlagSet("video", flag.ExitOnError)
	output := fs.String("o", "", "Output file path")
	regionStr := fs.String("r", "", "Capture region (x,y,w,h)")
	regionName := fs.String("region", "", "Use a saved region by name")
	fps := fs.Int("f", 30, "Frames per second")
	quality := fs.String("q", "medium", "Quality level (low, medium, high)")

	fs.Usage = func() {
		fmt.Println("Usage: witness video [options]")
		fmt.Println("\nRecord screen and save as MP4")
		fmt.Println("\nOptions:")
		fs.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  witness video -o tutorial.mp4")
		fmt.Println("  witness video -o tutorial.mp4 -f 30 -q high")
		fmt.Println("  witness video -region demo -o capture.mp4")
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	// TODO: Implement video recording
	fmt.Println("Video recording not yet implemented")
	fmt.Printf("Output: %s\n", *output)
	fmt.Printf("Region: %s\n", *regionStr)
	fmt.Printf("Region name: %s\n", *regionName)
	fmt.Printf("FPS: %d\n", *fps)
	fmt.Printf("Quality: %s\n", *quality)
}

func printUsage() {
	usage := `Witness - Screen Capture Tool
Version: ` + version + `

Usage: witness <command> [options]

Commands:
  select     Launch interactive region selector
  regions    Manage saved regions
  gif        Record and save as GIF
  video      Record and save as MP4 (coming soon)
  help       Show this help message
  version    Show version information

Quick Start:
  1. Select a capture region:
     witness select -name demo

  2. Record a GIF:
     witness gif -region demo -o demo.gif

  3. Or use a one-time region:
     witness gif -r 0,0,800,600 -o demo.gif

For detailed help on a command:
  witness <command> -help

Examples:
  witness select                      # Select region interactively
  witness select -name demo -default  # Select, save, and set as default
  witness regions                     # List all saved regions
  witness gif -o demo.gif -f 15       # Record GIF at 15 FPS
  witness gif -region demo -o out.gif # Use saved region
`
	fmt.Println(usage)
}
