package cli

import (
	"strings"
	"testing"
	"time"
)

func TestResolveStringOptionPrecedence(t *testing.T) {
	explicit := map[string]bool{"service": true}
	got := resolveStringOption("service", "flag", "env", "cfg", "default", explicit)
	if got != "flag" {
		t.Fatalf("expected flag value, got %q", got)
	}

	explicit = map[string]bool{}
	got = resolveStringOption("service", "flag", "env", "cfg", "default", explicit)
	if got != "env" {
		t.Fatalf("expected env value, got %q", got)
	}

	got = resolveStringOption("service", "flag", "", "cfg", "default", explicit)
	if got != "cfg" {
		t.Fatalf("expected config value, got %q", got)
	}
}

func TestResolveIntOptionParsesEnv(t *testing.T) {
	got, err := resolveIntOption("concurrency", "I18N_CONCURRENCY", 4, "8", 0, 4, map[string]bool{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 8 {
		t.Fatalf("expected 8, got %d", got)
	}
}

func TestResolveIntOptionRejectsInvalidEnv(t *testing.T) {
	_, err := resolveIntOption("concurrency", "I18N_CONCURRENCY", 4, "abc", 0, 4, map[string]bool{})
	if err == nil {
		t.Fatal("expected error for invalid int env")
	}
	if !strings.Contains(err.Error(), "I18N_CONCURRENCY") {
		t.Fatalf("expected env var name in error, got %v", err)
	}
}

func TestResolveDurationOptionParsesConfig(t *testing.T) {
	got, err := resolveDurationOption("timeout", "I18N_TIMEOUT", 0, "", "2m", 0, map[string]bool{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 2*time.Minute {
		t.Fatalf("expected 2m, got %s", got)
	}
}

func TestResolveBoolOptionParsesEnv(t *testing.T) {
	got, err := resolveBoolOption("dry-run", false, "true", false, false, map[string]bool{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Fatal("expected true")
	}
}
