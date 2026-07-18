package action

import "testing"

func TestDirArgsRun(t *testing.T) {
	exec := &fakeExecutor{}
	a := &DirArgs{Path: "/tmp/foo"}

	if _, _, err := a.Run(exec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "mkdir -p /tmp/foo"
	if len(exec.runs) != 1 || exec.runs[0] != want {
		t.Errorf("got %v, want [%q]", exec.runs, want)
	}
}

func TestDirArgsRunWithMod(t *testing.T) {
	exec := &fakeExecutor{}
	a := &DirArgs{Path: "/tmp/foo", Mod: "755"}

	if _, _, err := a.Run(exec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "mkdir -p /tmp/foo -m 755"
	if len(exec.runs) != 1 || exec.runs[0] != want {
		t.Errorf("got %v, want [%q]", exec.runs, want)
	}
}
