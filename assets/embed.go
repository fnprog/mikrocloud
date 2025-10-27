// Package assets Imports the frontend
package assets

import "embed"

// FrontendFS The entire dist directory from the web build
//
//go:embed all:dist
var FrontendFS embed.FS
