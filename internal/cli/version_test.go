package cli

import (
	"runtime/debug"
	"testing"
)

func TestAppVersionPrefersLdflagsVersion(t *testing.T) {
	oldVersion := Version
	oldRead := readBuildInfo
	t.Cleanup(func() {
		Version = oldVersion
		readBuildInfo = oldRead
	})

	Version = "v9.9.9"
	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{
			Main: debug.Module{Version: "v1.0.0"},
		}, true
	}

	if got := appVersion(); got != "v9.9.9" {
		t.Fatalf("expected ldflags version, got %q", got)
	}
}

func TestAppVersionUsesBuildInfo(t *testing.T) {
	oldVersion := Version
	oldRead := readBuildInfo
	t.Cleanup(func() {
		Version = oldVersion
		readBuildInfo = oldRead
	})

	Version = ""
	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{
			Main: debug.Module{Version: "v1.2.3"},
		}, true
	}

	if got := appVersion(); got != "v1.2.3" {
		t.Fatalf("expected build-info version, got %q", got)
	}
}

func TestAppVersionFallbackDev(t *testing.T) {
	oldVersion := Version
	oldRead := readBuildInfo
	t.Cleanup(func() {
		Version = oldVersion
		readBuildInfo = oldRead
	})

	Version = ""
	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{
			Main: debug.Module{Version: "(devel)"},
		}, true
	}

	if got := appVersion(); got != "dev" {
		t.Fatalf("expected dev fallback, got %q", got)
	}
}
