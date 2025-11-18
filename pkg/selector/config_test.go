package selector

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ericmhalvorsen/witness/pkg/capture"
)

// Helper function to create a temporary config directory
func setupTestConfig(t *testing.T) (string, func()) {
	t.Helper()

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "witness-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)

	cleanup := func() {
		os.Setenv("HOME", oldHome)
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func TestSaveAndLoadRegion(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	region := &capture.Region{
		X:      100,
		Y:      200,
		Width:  800,
		Height: 600,
	}

	// Test saving
	err := SaveRegion("test-region", region)
	if err != nil {
		t.Fatalf("SaveRegion() failed: %v", err)
	}

	// Test loading
	loaded, err := LoadRegion("test-region")
	if err != nil {
		t.Fatalf("LoadRegion() failed: %v", err)
	}

	if loaded.X != region.X || loaded.Y != region.Y ||
		loaded.Width != region.Width || loaded.Height != region.Height {
		t.Errorf("Loaded region %+v doesn't match saved region %+v", loaded, region)
	}
}

func TestLoadRegionNotFound(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	_, err := LoadRegion("nonexistent")
	if err == nil {
		t.Error("LoadRegion() should fail for nonexistent region")
	}
}

func TestListRegions(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Initially should be empty
	regions, err := ListRegions()
	if err != nil {
		t.Fatalf("ListRegions() failed: %v", err)
	}
	if len(regions) != 0 {
		t.Errorf("Expected 0 regions, got %d", len(regions))
	}

	// Add some regions
	region1 := &capture.Region{X: 0, Y: 0, Width: 100, Height: 100}
	region2 := &capture.Region{X: 10, Y: 10, Width: 200, Height: 200}

	SaveRegion("region1", region1)
	SaveRegion("region2", region2)

	// List should now contain both
	regions, err = ListRegions()
	if err != nil {
		t.Fatalf("ListRegions() failed: %v", err)
	}
	if len(regions) != 2 {
		t.Errorf("Expected 2 regions, got %d", len(regions))
	}

	// Check that both names are present
	names := make(map[string]bool)
	for _, name := range regions {
		names[name] = true
	}
	if !names["region1"] || !names["region2"] {
		t.Errorf("Missing expected region names, got: %v", regions)
	}
}

func TestDeleteRegion(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	region := &capture.Region{X: 0, Y: 0, Width: 100, Height: 100}
	SaveRegion("test-delete", region)

	// Verify it exists
	_, err := LoadRegion("test-delete")
	if err != nil {
		t.Fatalf("Region should exist before delete: %v", err)
	}

	// Delete it
	err = DeleteRegion("test-delete")
	if err != nil {
		t.Fatalf("DeleteRegion() failed: %v", err)
	}

	// Verify it's gone
	_, err = LoadRegion("test-delete")
	if err == nil {
		t.Error("Region should not exist after delete")
	}
}

func TestDeleteRegionNotFound(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	err := DeleteRegion("nonexistent")
	if err == nil {
		t.Error("DeleteRegion() should fail for nonexistent region")
	}
}

func TestSetAndGetDefaultRegion(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	region := &capture.Region{X: 50, Y: 50, Width: 300, Height: 300}
	SaveRegion("my-default", region)

	// Set as default
	err := SetDefaultRegion("my-default")
	if err != nil {
		t.Fatalf("SetDefaultRegion() failed: %v", err)
	}

	// Get default
	defaultRegion, err := GetDefaultRegion()
	if err != nil {
		t.Fatalf("GetDefaultRegion() failed: %v", err)
	}

	if defaultRegion.X != region.X || defaultRegion.Y != region.Y ||
		defaultRegion.Width != region.Width || defaultRegion.Height != region.Height {
		t.Errorf("Default region %+v doesn't match expected %+v", defaultRegion, region)
	}
}

func TestSetDefaultRegionNotFound(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	err := SetDefaultRegion("nonexistent")
	if err == nil {
		t.Error("SetDefaultRegion() should fail for nonexistent region")
	}
}

func TestGetDefaultRegionNotSet(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	_, err := GetDefaultRegion()
	if err == nil {
		t.Error("GetDefaultRegion() should fail when no default is set")
	}
}

func TestGetRegionInfo(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	region := &capture.Region{X: 100, Y: 200, Width: 800, Height: 600}
	SaveRegion("test-info", region)

	info, err := GetRegionInfo("test-info")
	if err != nil {
		t.Fatalf("GetRegionInfo() failed: %v", err)
	}

	// Check that info contains the important data
	expected := "test-info: 800x600 at (100,200)"
	if info != expected {
		t.Errorf("GetRegionInfo() = %q, want %q", info, expected)
	}
}

func TestConfigFilePersistence(t *testing.T) {
	tmpDir, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create and save a region
	region := &capture.Region{X: 10, Y: 20, Width: 100, Height: 200}
	SaveRegion("persistent", region)

	// Verify the file was created
	configPath := filepath.Join(tmpDir, ".config", "witness", "regions.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Read the file and verify JSON structure
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var config RegionConfig
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to parse config JSON: %v", err)
	}

	if len(config.Regions) != 1 {
		t.Errorf("Expected 1 region in config, got %d", len(config.Regions))
	}

	if savedRegion, ok := config.Regions["persistent"]; !ok {
		t.Error("Region 'persistent' not found in config")
	} else {
		if savedRegion.X != 10 || savedRegion.Y != 20 ||
			savedRegion.Width != 100 || savedRegion.Height != 200 {
			t.Errorf("Saved region data incorrect: %+v", savedRegion)
		}
	}
}

func TestMultipleRegionsManagement(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create multiple regions
	regions := map[string]*capture.Region{
		"fullscreen": {X: 0, Y: 0, Width: 1920, Height: 1080},
		"window":     {X: 100, Y: 100, Width: 800, Height: 600},
		"corner":     {X: 0, Y: 0, Width: 400, Height: 400},
	}

	// Save all regions
	for name, region := range regions {
		if err := SaveRegion(name, region); err != nil {
			t.Fatalf("Failed to save region %s: %v", name, err)
		}
	}

	// Verify all can be loaded
	for name, expectedRegion := range regions {
		loadedRegion, err := LoadRegion(name)
		if err != nil {
			t.Fatalf("Failed to load region %s: %v", name, err)
		}
		if loadedRegion.X != expectedRegion.X || loadedRegion.Y != expectedRegion.Y ||
			loadedRegion.Width != expectedRegion.Width || loadedRegion.Height != expectedRegion.Height {
			t.Errorf("Region %s mismatch: got %+v, want %+v", name, loadedRegion, expectedRegion)
		}
	}

	// Verify list contains all regions
	list, err := ListRegions()
	if err != nil {
		t.Fatalf("ListRegions() failed: %v", err)
	}
	if len(list) != len(regions) {
		t.Errorf("Expected %d regions in list, got %d", len(regions), len(list))
	}
}

func TestOverwriteExistingRegion(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	region1 := &capture.Region{X: 0, Y: 0, Width: 100, Height: 100}
	region2 := &capture.Region{X: 50, Y: 50, Width: 200, Height: 200}

	// Save initial region
	SaveRegion("overwrite-test", region1)

	// Overwrite with new region
	SaveRegion("overwrite-test", region2)

	// Load and verify it's the new region
	loaded, err := LoadRegion("overwrite-test")
	if err != nil {
		t.Fatalf("LoadRegion() failed: %v", err)
	}

	if loaded.X != region2.X || loaded.Y != region2.Y ||
		loaded.Width != region2.Width || loaded.Height != region2.Height {
		t.Errorf("Expected overwritten region %+v, got %+v", region2, loaded)
	}
}
