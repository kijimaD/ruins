package assets

import "embed"

// FS is the embedded file system containing game assets
//
//go:embed file metadata levels
var FS embed.FS
