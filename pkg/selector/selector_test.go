package selector

import (
	"testing"

	"github.com/ericmhalvorsen/witness/pkg/capture"
)

func TestParseRegionString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *capture.Region
		wantErr bool
	}{
		{
			name:  "valid region",
			input: "100,200,800,600",
			want: &capture.Region{
				X:      100,
				Y:      200,
				Width:  800,
				Height: 600,
			},
			wantErr: false,
		},
		{
			name:  "valid region with zeros",
			input: "0,0,1920,1080",
			want: &capture.Region{
				X:      0,
				Y:      0,
				Width:  1920,
				Height: 1080,
			},
			wantErr: false,
		},
		{
			name:    "invalid format - missing value",
			input:   "100,200,800",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "extra values ignored",
			input: "100,200,800,600,100",
			want: &capture.Region{
				X:      100,
				Y:      200,
				Width:  800,
				Height: 600,
			},
			wantErr: false,
		},
		{
			name:    "invalid format - non-numeric",
			input:   "abc,200,800,600",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid - zero width",
			input:   "100,200,0,600",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid - zero height",
			input:   "100,200,800,0",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid - negative width",
			input:   "100,200,-800,600",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid - negative height",
			input:   "100,200,800,-600",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRegionString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRegionString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("ParseRegionString() returned nil, expected region")
					return
				}
				if got.X != tt.want.X || got.Y != tt.want.Y ||
					got.Width != tt.want.Width || got.Height != tt.want.Height {
					t.Errorf("ParseRegionString() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}

func TestFormatRegionString(t *testing.T) {
	tests := []struct {
		name  string
		input *capture.Region
		want  string
	}{
		{
			name: "valid region",
			input: &capture.Region{
				X:      100,
				Y:      200,
				Width:  800,
				Height: 600,
			},
			want: "100,200,800,600",
		},
		{
			name: "region at origin",
			input: &capture.Region{
				X:      0,
				Y:      0,
				Width:  1920,
				Height: 1080,
			},
			want: "0,0,1920,1080",
		},
		{
			name:  "nil region",
			input: nil,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatRegionString(tt.input)
			if got != tt.want {
				t.Errorf("FormatRegionString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAndFormatRoundTrip(t *testing.T) {
	tests := []string{
		"0,0,1920,1080",
		"100,200,800,600",
		"50,50,100,100",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			region, err := ParseRegionString(tt)
			if err != nil {
				t.Fatalf("ParseRegionString() failed: %v", err)
			}
			formatted := FormatRegionString(region)
			if formatted != tt {
				t.Errorf("Round trip failed: got %v, want %v", formatted, tt)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Message == "" {
		t.Error("DefaultConfig() returned empty message")
	}
	if !config.ShowDimensions {
		t.Error("DefaultConfig() ShowDimensions should be true by default")
	}
}
