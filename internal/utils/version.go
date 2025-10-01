package utils

import (
	_ "embed"
)

//go:embed VERSION
var version string

func GetVersion() string {
	if version == "" {
		return "unknown"
	}
	return version
}

func IsDevelopment() bool {
	return version == "development"
}
