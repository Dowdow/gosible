package action

import (
	"reflect"
	"strings"
	"testing"
)

func TestFileArgsPrepareReplacesEnv(t *testing.T) {
	a := &FileArgs{Content: []string{"A=env(FOO)", "B=bar"}}

	replaceEnv := func(s string) string { return strings.ReplaceAll(s, "env(FOO)", "baz") }
	if err := a.Prepare(identityResolvePath, replaceEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"A=baz", "B=bar"}
	if !reflect.DeepEqual(a.Content, want) {
		t.Errorf("got %v, want %v", a.Content, want)
	}
}

func TestFileArgsRunCopiesGeneratedFileContent(t *testing.T) {
	exec := &fakeExecutor{}
	a := &FileArgs{Dest: "/remote/path/file.txt", Content: []string{"line1", "line2"}}

	if _, _, err := a.Run(exec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(exec.starts) != 1 {
		t.Fatalf("expected exactly one Start call, got %d", len(exec.starts))
	}
	if got := exec.stdin.String(); !strings.Contains(got, "line1\nline2") {
		t.Errorf("expected written content to contain %q, got %q", "line1\nline2", got)
	}
}
