package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseEnvMissingFileIsNotAnError(t *testing.T) {
	c := &Config{}

	if err := c.ParseEnv("/nonexistent/.env"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseEnvAndReplaceEnv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := "# comment\nFOO=bar\n\nBAZ=qux\n"
	if err := os.WriteFile(envPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	c := &Config{}
	if err := c.ParseEnv(envPath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := c.ReplaceEnv("value is env(FOO) and env(BAZ)")
	want := "value is bar and qux"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestReplaceEnvUnknownVarBecomesEmpty(t *testing.T) {
	c := &Config{}

	got := c.ReplaceEnv("x=env(MISSING)")
	want := "x="
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestReplaceEnvNoPlaceholder(t *testing.T) {
	c := &Config{}

	got := c.ReplaceEnv("plain string")
	if got != "plain string" {
		t.Errorf("got %q, want unchanged string", got)
	}
}
