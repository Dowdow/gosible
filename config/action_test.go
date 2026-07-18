package config

import (
	"encoding/json"
	"testing"

	"github.com/Dowdow/gosible/action"
)

func TestActionUnmarshalShell(t *testing.T) {
	var a Action
	if err := json.Unmarshal([]byte(`{"name":"n","type":"shell","args":"echo hi"}`), &a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	shellArgs, ok := a.Args.(*action.ShellArgs)
	if !ok {
		t.Fatalf("expected *action.ShellArgs, got %T", a.Args)
	}
	if string(*shellArgs) != "echo hi" {
		t.Errorf("got %q, want %q", *shellArgs, "echo hi")
	}
}

func TestActionUnmarshalCopy(t *testing.T) {
	var a Action
	if err := json.Unmarshal([]byte(`{"type":"copy","args":{"src":"a","dest":"b"}}`), &a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	copyArgs, ok := a.Args.(*action.CopyArgs)
	if !ok {
		t.Fatalf("expected *action.CopyArgs, got %T", a.Args)
	}
	if copyArgs.Src != "a" || copyArgs.Dest != "b" {
		t.Errorf("got %+v", copyArgs)
	}
}

func TestActionUnmarshalDir(t *testing.T) {
	var a Action
	if err := json.Unmarshal([]byte(`{"type":"dir","args":{"path":"/tmp/x","mod":"755"}}`), &a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dirArgs, ok := a.Args.(*action.DirArgs)
	if !ok {
		t.Fatalf("expected *action.DirArgs, got %T", a.Args)
	}
	if dirArgs.Path != "/tmp/x" || dirArgs.Mod != "755" {
		t.Errorf("got %+v", dirArgs)
	}
}

func TestActionUnmarshalDocker(t *testing.T) {
	var a Action
	if err := json.Unmarshal([]byte(`{"type":"docker","args":{"src":"./img","dest":"/x","image":"img:latest"}}`), &a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dockerArgs, ok := a.Args.(*action.DockerArgs)
	if !ok {
		t.Fatalf("expected *action.DockerArgs, got %T", a.Args)
	}
	if dockerArgs.Image != "img:latest" {
		t.Errorf("got %+v", dockerArgs)
	}
}

func TestActionUnmarshalFile(t *testing.T) {
	var a Action
	if err := json.Unmarshal([]byte(`{"type":"file","args":{"dest":"/x","content":["a","b"]}}`), &a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fileArgs, ok := a.Args.(*action.FileArgs)
	if !ok {
		t.Fatalf("expected *action.FileArgs, got %T", a.Args)
	}
	if len(fileArgs.Content) != 2 {
		t.Errorf("got %+v", fileArgs)
	}
}

func TestActionUnmarshalUnknownType(t *testing.T) {
	var a Action
	if err := json.Unmarshal([]byte(`{"type":"bogus"}`), &a); err == nil {
		t.Fatal("expected an error for an unknown action type")
	}
}

func TestActionUnmarshalIdOnlyReference(t *testing.T) {
	var a Action
	if err := json.Unmarshal([]byte(`{"id":"shared-action"}`), &a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Args != nil {
		t.Errorf("expected nil Args for an id-only reference, got %+v", a.Args)
	}
}
