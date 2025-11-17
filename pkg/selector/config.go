package selector

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ericmhalvorsen/witness/pkg/capture"
)

// RegionConfig stores saved regions
type RegionConfig struct {
	Regions map[string]*capture.Region `json:"regions"`
	Default string                      `json:"default,omitempty"`
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "witness")
	configFile := filepath.Join(configDir, "regions.json")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configFile, nil
}

// loadConfig loads the region configuration
func loadConfig() (*RegionConfig, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, return empty config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &RegionConfig{
			Regions: make(map[string]*capture.Region),
		}, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config RegionConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if config.Regions == nil {
		config.Regions = make(map[string]*capture.Region)
	}

	return &config, nil
}

// saveConfig saves the region configuration
func saveConfig(config *RegionConfig) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// SaveRegion saves a named region
func SaveRegion(name string, region *capture.Region) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	config.Regions[name] = region

	return saveConfig(config)
}

// LoadRegion loads a named region
func LoadRegion(name string) (*capture.Region, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	region, exists := config.Regions[name]
	if !exists {
		return nil, fmt.Errorf("region '%s' not found", name)
	}

	return region, nil
}

// ListRegions returns all saved region names
func ListRegions() ([]string, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(config.Regions))
	for name := range config.Regions {
		names = append(names, name)
	}

	return names, nil
}

// GetRegionInfo returns detailed information about a saved region
func GetRegionInfo(name string) (string, error) {
	region, err := LoadRegion(name)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s: %dx%d at (%d,%d)",
		name, region.Width, region.Height, region.X, region.Y), nil
}

// DeleteRegion deletes a named region
func DeleteRegion(name string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	if _, exists := config.Regions[name]; !exists {
		return fmt.Errorf("region '%s' not found", name)
	}

	delete(config.Regions, name)

	return saveConfig(config)
}

// SetDefaultRegion sets the default region to use
func SetDefaultRegion(name string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	if _, exists := config.Regions[name]; !exists {
		return fmt.Errorf("region '%s' not found", name)
	}

	config.Default = name

	return saveConfig(config)
}

// GetDefaultRegion gets the default region
func GetDefaultRegion() (*capture.Region, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	if config.Default == "" {
		return nil, fmt.Errorf("no default region set")
	}

	return LoadRegion(config.Default)
}
