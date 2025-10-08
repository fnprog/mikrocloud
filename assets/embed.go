package assets

import "embed"

// Embed the entire dist directory from the web build
//
//go:embed all:dist
var FrontendFS embed.FS
