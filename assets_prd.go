//go:build !dev
package main

import (
	"embed"
)

//go:embed content/projects.json
var contentFS embed.FS

//go:embed all:static
var staticFS embed.FS

//go:embed all:templates
var templateFS embed.FS

//go:embed all:images
var imageFS embed.FS
