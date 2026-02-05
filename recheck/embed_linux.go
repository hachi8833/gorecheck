//go:build linux

package recheck

import "embed"

//go:embed bin/recheck-linux-amd64
var embeddedBinary embed.FS

const binaryName = "bin/recheck-linux-amd64"
