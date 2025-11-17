// +build darwin

package selector

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestMacOSSelectorReadLastSelection(t *testing.T) {
	mockCmd := NewMockSystemCommand()

	// Set up mock output that matches macOS defaults output format
	mockOutput := `{
    Height = 600;
    Width = 800;
    X = 100;
    Y = 200;
}`
	mockCmd.SetOutput("defaults", []byte(mockOutput))

	selector := NewMacOSSelectorWithExecutor(mockCmd).(*macOSSelector)

	region, err := selector.readLastSelection()
	if err != nil {
		t.Fatalf("readLastSelection() failed: %v", err)
	}

	if region.X != 100 {
		t.Errorf("X = %d, want 100", region.X)
	}
	if region.Y != 200 {
		t.Errorf("Y = %d, want 200", region.Y)
	}
	if region.Width != 800 {
		t.Errorf("Width = %d, want 800", region.Width)
	}
	if region.Height != 600 {
		t.Errorf("Height = %d, want 600", region.Height)
	}

	// Verify the command was called with correct arguments
	if !mockCmd.WasCalled("defaults", "read", "com.apple.screencapture", "last-selection") {
		t.Error("defaults command was not called with expected arguments")
	}
}

func TestMacOSSelectorReadLastSelectionError(t *testing.T) {
	mockCmd := NewMockSystemCommand()
	mockCmd.SetError("defaults", fmt.Errorf("command failed"))

	selector := NewMacOSSelectorWithExecutor(mockCmd).(*macOSSelector)

	_, err := selector.readLastSelection()
	if err == nil {
		t.Error("readLastSelection() should fail when command fails")
	}
}

func TestMacOSSelectorReadLastSelectionInvalidDimensions(t *testing.T) {
	mockCmd := NewMockSystemCommand()

	// Set up mock output with zero width
	mockOutput := `{
    Height = 600;
    Width = 0;
    X = 100;
    Y = 200;
}`
	mockCmd.SetOutput("defaults", []byte(mockOutput))

	selector := NewMacOSSelectorWithExecutor(mockCmd).(*macOSSelector)

	_, err := selector.readLastSelection()
	if err == nil {
		t.Error("readLastSelection() should fail for invalid dimensions")
	}
}

func TestMacOSSelectorReadLastSelectionMalformedOutput(t *testing.T) {
	tests := []struct {
		name   string
		output string
	}{
		{
			name:   "missing height",
			output: `{Width = 800; X = 100; Y = 200;}`,
		},
		{
			name:   "missing width",
			output: `{Height = 600; X = 100; Y = 200;}`,
		},
		{
			name:   "non-numeric values",
			output: `{Height = abc; Width = def; X = 100; Y = 200;}`,
		},
		{
			name:   "empty output",
			output: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCmd := NewMockSystemCommand()
			mockCmd.SetOutput("defaults", []byte(tt.output))

			selector := NewMacOSSelectorWithExecutor(mockCmd).(*macOSSelector)

			_, err := selector.readLastSelection()
			if err == nil {
				t.Error("readLastSelection() should fail for malformed output")
			}
		})
	}
}

func TestMacOSSelectorSelectWithName(t *testing.T) {
	tmpDir, cleanup := setupTestConfig(t)
	defer cleanup()

	mockCmd := NewMockSystemCommand()

	// Mock the screencapture command (should succeed)
	// We need to create the temp file to simulate successful capture
	mockCmd.RunInteractive = func(name string, args ...string) error {
		if name == "screencapture" {
			// Create the temporary file to simulate successful selection
			tmpFile := args[len(args)-1] // Last argument is the file path
			return os.WriteFile(tmpFile, []byte("fake image data"), 0644)
		}
		return nil
	}

	// Mock the defaults read command
	mockOutput := `{
    Height = 600;
    Width = 800;
    X = 100;
    Y = 200;
}`
	mockCmd.SetOutput("defaults", []byte(mockOutput))

	selector := NewMacOSSelectorWithExecutor(mockCmd)

	// Select and save a region
	region, err := selector.SelectWithName("test-region")
	if err != nil {
		t.Fatalf("SelectWithName() failed: %v", err)
	}

	if region == nil {
		t.Fatal("SelectWithName() returned nil region")
	}

	// Verify the region was saved
	configPath := filepath.Join(tmpDir, ".config", "witness", "regions.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Verify we can load the region
	loaded, err := LoadRegion("test-region")
	if err != nil {
		t.Fatalf("Failed to load saved region: %v", err)
	}

	if loaded.X != region.X || loaded.Y != region.Y ||
		loaded.Width != region.Width || loaded.Height != region.Height {
		t.Errorf("Loaded region %+v doesn't match selected region %+v", loaded, region)
	}
}

func TestMacOSSelectorSelectCanceled(t *testing.T) {
	mockCmd := NewMockSystemCommand()
	mockCmd.SetError("screencapture", fmt.Errorf("user canceled"))

	selector := NewMacOSSelectorWithExecutor(mockCmd)

	_, err := selector.Select()
	if err == nil {
		t.Error("Select() should fail when user cancels")
	}
}

func TestMacOSSelectorParseDifferentFormats(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		wantX    int
		wantY    int
		wantW    int
		wantH    int
		wantErr  bool
	}{
		{
			name: "standard format",
			output: `{
    Height = 480;
    Width = 640;
    X = 50;
    Y = 100;
}`,
			wantX: 50,
			wantY: 100,
			wantW: 640,
			wantH: 480,
			wantErr: false,
		},
		{
			name: "compact format",
			output: `{Height = 480; Width = 640; X = 50; Y = 100;}`,
			wantX: 50,
			wantY: 100,
			wantW: 640,
			wantH: 480,
			wantErr: false,
		},
		{
			name: "with decimal values",
			output: `{
    Height = 480.5;
    Width = 640.7;
    X = 50.2;
    Y = 100.9;
}`,
			wantX: 50,
			wantY: 100,
			wantW: 640,
			wantH: 480,
			wantErr: false,
		},
		{
			name: "large values",
			output: `{
    Height = 2160;
    Width = 3840;
    X = 0;
    Y = 0;
}`,
			wantX: 0,
			wantY: 0,
			wantW: 3840,
			wantH: 2160,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCmd := NewMockSystemCommand()
			mockCmd.SetOutput("defaults", []byte(tt.output))

			selector := NewMacOSSelectorWithExecutor(mockCmd).(*macOSSelector)

			region, err := selector.readLastSelection()
			if (err != nil) != tt.wantErr {
				t.Errorf("readLastSelection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if region.X != tt.wantX {
					t.Errorf("X = %d, want %d", region.X, tt.wantX)
				}
				if region.Y != tt.wantY {
					t.Errorf("Y = %d, want %d", region.Y, tt.wantY)
				}
				if region.Width != tt.wantW {
					t.Errorf("Width = %d, want %d", region.Width, tt.wantW)
				}
				if region.Height != tt.wantH {
					t.Errorf("Height = %d, want %d", region.Height, tt.wantH)
				}
			}
		})
	}
}

func TestSystemCommandMock(t *testing.T) {
	mockCmd := NewMockSystemCommand()

	// Test setting and getting output
	mockCmd.SetOutput("test", []byte("output"))
	output, err := mockCmd.Run("test")
	if err != nil {
		t.Errorf("Run() failed: %v", err)
	}
	if string(output) != "output" {
		t.Errorf("Output = %q, want %q", output, "output")
	}

	// Test setting and getting error
	testErr := fmt.Errorf("test error")
	mockCmd.SetError("failing-cmd", testErr)
	_, err = mockCmd.Run("failing-cmd")
	if err == nil {
		t.Error("Run() should return error for failing command")
	}

	// Test call counting
	mockCmd.Reset()
	mockCmd.Run("count-test", "arg1", "arg2")
	mockCmd.Run("count-test", "arg3")
	mockCmd.Run("other-cmd")

	if count := mockCmd.GetCallCount("count-test"); count != 2 {
		t.Errorf("GetCallCount() = %d, want 2", count)
	}

	// Test WasCalled
	if !mockCmd.WasCalled("count-test", "arg1", "arg2") {
		t.Error("WasCalled() should return true for called command")
	}
	if mockCmd.WasCalled("count-test", "wrong", "args") {
		t.Error("WasCalled() should return false for wrong arguments")
	}

	// Test Reset
	mockCmd.Reset()
	if count := mockCmd.GetCallCount("count-test"); count != 0 {
		t.Errorf("After Reset(), GetCallCount() = %d, want 0", count)
	}
}
