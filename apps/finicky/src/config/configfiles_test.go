package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestGetConfigPathsDefaultsToDotConfig(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")

	cfw := &ConfigFileWatcher{}
	paths := cfw.GetConfigPaths()

	if len(paths) != 6 {
		t.Fatalf("expected 6 config paths, got %d: %v", len(paths), paths)
	}

	for _, p := range paths {
		if strings.Contains(p, "~") {
			t.Errorf("path %q was not expanded", p)
		}
	}

	if !strings.HasSuffix(paths[2], filepath.Join(".config", "finicky.js")) {
		t.Errorf("expected fallback to ~/.config, got %q", paths[2])
	}
	if !strings.HasSuffix(paths[4], filepath.Join(".config", "finicky", "finicky.js")) {
		t.Errorf("expected fallback to ~/.config/finicky, got %q", paths[4])
	}
}

func TestGetConfigPathsHonorsXDGConfigHome(t *testing.T) {
	xdgHome := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdgHome)

	cfw := &ConfigFileWatcher{}
	paths := cfw.GetConfigPaths()

	expected := []string{
		filepath.Join(xdgHome, "finicky.js"),
		filepath.Join(xdgHome, "finicky.ts"),
		filepath.Join(xdgHome, "finicky", "finicky.js"),
		filepath.Join(xdgHome, "finicky", "finicky.ts"),
	}

	for i, want := range expected {
		got := paths[i+2]
		if got != want {
			t.Errorf("paths[%d] = %q, want %q", i+2, got, want)
		}
	}
}

func TestNearestExistingDir(t *testing.T) {
	base := t.TempDir()

	if got := nearestExistingDir(base); got != base {
		t.Errorf("nearestExistingDir(%q) = %q, want the dir itself", base, got)
	}

	missing := filepath.Join(base, "finicky", "nested")
	if got := nearestExistingDir(missing); got != base {
		t.Errorf("nearestExistingDir(%q) = %q, want nearest ancestor %q", missing, got, base)
	}
}

func TestIsAncestorOfAny(t *testing.T) {
	folders := map[string]bool{
		"/home/user/.config/finicky": true,
	}

	cases := []struct {
		path string
		want bool
	}{
		{"/home/user/.config/finicky", true},
		{"/home/user/.config", true},
		{"/home/user", true},
		{"/home/user/.config/finick", false},
		{"/home/other", false},
	}

	for _, c := range cases {
		if got := isAncestorOfAny(c.path, folders); got != c.want {
			t.Errorf("isAncestorOfAny(%q) = %v, want %v", c.path, got, c.want)
		}
	}
}

func TestGetConfigPathsCustomPathWins(t *testing.T) {
	cfw := &ConfigFileWatcher{customConfigPath: "/tmp/custom-finicky.js"}
	paths := cfw.GetConfigPaths()

	if len(paths) != 1 || paths[0] != "/tmp/custom-finicky.js" {
		t.Errorf("expected only the custom path, got %v", paths)
	}
}
