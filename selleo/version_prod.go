//go:build prod
// +build prod

package selleo

import (
	_ "embed"
)

//go:embed version.txt
var version string
