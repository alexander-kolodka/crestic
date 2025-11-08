package version

import (
	"runtime/debug"
)

var Version = "" // can be set via -ldflags

func String() string {
	if Version != "" {
		return Version
	}

	if bi, ok := debug.ReadBuildInfo(); ok && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		return bi.Main.Version
	}

	return "dev"
}
