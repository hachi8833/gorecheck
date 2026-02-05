//go:build !linux

package recheck

import "embed"

//go:embed bin/recheck
var embeddedBinary embed.FS

const binaryName = "bin/recheck"
