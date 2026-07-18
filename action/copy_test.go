package action

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyArgsPrepareResolvesSrc(t *testing.T) {
	a := &CopyArgs{Src: "relative/file"}

	resolvePath := func(p string) string { return "/resolved/" + p }
	if err := a.Prepare(resolvePath, identityReplaceEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "/resolved/relative/file"
	if a.Src != want {
		t.Errorf("got src=%q, want %q", a.Src, want)
	}
}

func TestCopyArgsRunSingleFile(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "hello.txt")
	if err := os.WriteFile(srcPath, []byte("hello world"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(srcPath, 0o644); err != nil {
		t.Fatal(err)
	}

	exec := &fakeExecutor{}
	a := &CopyArgs{Src: srcPath, Dest: "/remote/hello.txt"}

	if _, _, err := a.Run(exec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantCmd := "scp -tr /remote"
	if len(exec.starts) != 1 || exec.starts[0] != wantCmd {
		t.Fatalf("got starts=%v, want [%q]", exec.starts, wantCmd)
	}

	got := exec.stdin.String()
	if !strings.HasPrefix(got, "C0644 11 hello.txt\n") {
		t.Errorf("expected scp header for hello.txt, got %q", got)
	}
	if !strings.Contains(got, "hello world") {
		t.Errorf("expected file content in stream, got %q", got)
	}
}

func TestCopyArgsRunDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("A"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sub", "b.txt"), []byte("B"), 0o644); err != nil {
		t.Fatal(err)
	}

	exec := &fakeExecutor{}
	a := &CopyArgs{Src: dir, Dest: "/remote/mydir"}

	if _, _, err := a.Run(exec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := exec.stdin.String()
	for _, want := range []string{"mydir", "sub", "a.txt", "b.txt", "A", "B", "E\n"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected stream to contain %q, got %q", want, got)
		}
	}
}

func TestCopyArgsRunMissingSrc(t *testing.T) {
	exec := &fakeExecutor{}
	a := &CopyArgs{Src: "/does/not/exist", Dest: "/remote/x"}

	if _, _, err := a.Run(exec); err == nil {
		t.Fatal("expected an error for a missing source")
	}
}

func TestCopyArgsRunStartError(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "hello.txt")
	os.WriteFile(srcPath, []byte("hello"), 0o644)

	exec := &fakeExecutor{startErr: errors.New("boom")}
	a := &CopyArgs{Src: srcPath, Dest: "/remote/hello.txt"}

	if _, _, err := a.Run(exec); err == nil {
		t.Fatal("expected an error when Start fails")
	}
}

func TestCopyArgsRunWaitError(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "hello.txt")
	os.WriteFile(srcPath, []byte("hello"), 0o644)

	exec := &fakeExecutor{waitErr: errors.New("boom")}
	a := &CopyArgs{Src: srcPath, Dest: "/remote/hello.txt"}

	if _, _, err := a.Run(exec); err == nil {
		t.Fatal("expected an error when Wait fails")
	}
}
