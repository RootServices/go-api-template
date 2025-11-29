package version

import (
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name        string
		wantBuild   string
		wantBranch  string
		wantErr     bool
		description string
	}{
		{
			name:        "successfully parses embedded version.json",
			wantBuild:   "", // Will be populated from actual version.json
			wantBranch:  "", // Will be populated from actual version.json
			wantErr:     false,
			description: "should successfully parse the embedded version.json file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get()

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// Verify that we got valid version data
			if got.Build == "" {
				t.Error("Get() returned empty Build field")
			}

			if got.Branch == "" {
				t.Error("Get() returned empty Branch field")
			}
		})
	}
}
