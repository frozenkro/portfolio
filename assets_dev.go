//go:build dev
package main

import (
	"io/fs"
	"os"
)

// These are all set to "." for interchangeability with assets_prd.go equivalents, 
// which all expect the full paths for embedded artifacts

var contentFS fs.FS = os.DirFS(".")

var staticFS fs.FS = os.DirFS(".")

var templateFS fs.FS = os.DirFS(".")

var imageFS fs.FS = os.DirFS(".")
