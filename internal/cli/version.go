package cli

import (
	"runtime/debug"
	"strings"
)

// Version is set at build time via -ldflags "-X github.com/wowblvck/i18n-translator/internal/cli.Version=vX.Y.Z".
var Version string

var readBuildInfo = debug.ReadBuildInfo

func appVersion() string {
	if v := strings.TrimSpace(Version); v != "" {
		return v
	}

	if info, ok := readBuildInfo(); ok && info != nil {
		if v := strings.TrimSpace(info.Main.Version); v != "" && v != "(devel)" {
			return v
		}
	}

	return "dev"
}
