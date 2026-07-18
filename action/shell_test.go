package action

import (
	"errors"
	"testing"
)

func TestShellArgsRun(t *testing.T) {
	exec := &fakeExecutor{runOut: "out", runErrOut: "err"}
	a := ShellArgs("echo hi")

	stdout, stderr, err := a.Run(exec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stdout != "out" || stderr != "err" {
		t.Errorf("got stdout=%q stderr=%q", stdout, stderr)
	}
	if len(exec.runs) != 1 || exec.runs[0] != "echo hi" {
		t.Errorf("expected command %q to be run, got %v", "echo hi", exec.runs)
	}
}

func TestShellArgsRunError(t *testing.T) {
	exec := &fakeExecutor{runErr: errors.New("boom")}
	a := ShellArgs("false")

	if _, _, err := a.Run(exec); err == nil {
		t.Fatal("expected an error")
	}
}
