package version

import (
	_ "embed"
	"encoding/json"
)

//go:embed version.json
var versionBytes []byte

// Version represents the build version information
type Version struct {
	Build  string `json:"build"`
	Branch string `json:"branch"`
}

// Get returns the parsed version information
func Get() (Version, error) {
	var v Version
	err := json.Unmarshal(versionBytes, &v)
	return v, err
}
