package config

import (
	"testing"

	"github.com/Dowdow/gosible/action"
)

func TestHasTask(t *testing.T) {
	c := &Config{Tasks: []Task{{}, {}}}

	if !c.HasTask(0) || !c.HasTask(1) {
		t.Error("expected indexes 0 and 1 to be valid")
	}
	if c.HasTask(-1) || c.HasTask(2) {
		t.Error("expected out-of-range indexes to be invalid")
	}
}

func TestHasMachineUser(t *testing.T) {
	c := &Config{
		Inventory: []Machine{{Id: "m1", Users: []User{{User: "u1"}}}},
	}

	if !c.HasMachineUser("m1", "u1") {
		t.Error("expected m1/u1 to exist")
	}
	if c.HasMachineUser("m1", "unknown") || c.HasMachineUser("unknown", "u1") {
		t.Error("expected unknown combos to not exist")
	}
}

func TestValidateValidConfig(t *testing.T) {
	c := &Config{
		Inventory: []Machine{{Id: "m1", Users: []User{{User: "u1"}}}},
		Actions:   []Action{{Id: "a1", Args: action.ShellArgs("echo hi")}},
		Tasks: []Task{{
			Name:     "t1",
			Machines: []string{"m1", "m1.u1"},
			Actions:  []Action{{Id: "a1"}},
		}},
	}

	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateUnknownMachine(t *testing.T) {
	c := &Config{
		Tasks: []Task{{Name: "t1", Machines: []string{"missing"}}},
	}

	if err := c.Validate(); err == nil {
		t.Fatal("expected an error for an unknown machine id")
	}
}

func TestValidateUnknownMachineUserCombo(t *testing.T) {
	c := &Config{
		Inventory: []Machine{{Id: "m1", Users: []User{{User: "u1"}}}},
		Tasks:     []Task{{Name: "t1", Machines: []string{"m1.unknown"}}},
	}

	if err := c.Validate(); err == nil {
		t.Fatal("expected an error for an unknown machine.user combo")
	}
}

func TestValidateUnknownActionId(t *testing.T) {
	c := &Config{
		Tasks: []Task{{Name: "t1", Actions: []Action{{Id: "missing"}}}},
	}

	if err := c.Validate(); err == nil {
		t.Fatal("expected an error for an unknown action id")
	}
}

func TestValidateInvalidActionArgs(t *testing.T) {
	c := &Config{
		Actions: []Action{{Id: "bad", Args: &action.DockerArgs{Tar: "not-a-tar"}}},
	}

	if err := c.Validate(); err == nil {
		t.Fatal("expected an error for invalid action args")
	}
}

func TestValidatePreparesActionArgs(t *testing.T) {
	c := &Config{
		configDir: "/base",
		Actions:   []Action{{Id: "a1", Args: &action.CopyArgs{Src: "rel/file", Dest: "/remote"}}},
	}

	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	copyArgs := c.Actions[0].Args.(*action.CopyArgs)
	want := "/base/rel/file"
	if copyArgs.Src != want {
		t.Errorf("got src=%q, want %q", copyArgs.Src, want)
	}
}

func TestConvertGenericAndSpecificActions(t *testing.T) {
	c := &Config{
		Inventory: []Machine{{Id: "m1", Name: "Machine 1", Address: "1.2.3.4:22", Users: []User{{User: "u1", Password: "p"}}}},
		Actions:   []Action{{Id: "shared", Name: "Shared", Args: action.ShellArgs("echo shared")}},
		Tasks: []Task{{
			Name: "t1",
			Actions: []Action{
				{Id: "shared"},
				{Name: "Inline", Args: action.ShellArgs("echo inline")},
			},
		}},
	}

	rc, err := c.Convert(0, "m1", "u1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rc.Machine.Address != "1.2.3.4:22" || rc.Machine.User != "u1" || rc.Machine.Password != "p" {
		t.Errorf("unexpected machine: %+v", rc.Machine)
	}

	if len(rc.Actions) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(rc.Actions))
	}
	if rc.Actions[0].Name != "Shared" || rc.Actions[1].Name != "Inline" {
		t.Errorf("unexpected action names: %q, %q", rc.Actions[0].Name, rc.Actions[1].Name)
	}
}

func TestConvertInvalidTaskIndex(t *testing.T) {
	c := &Config{Tasks: []Task{{}}}

	if _, err := c.Convert(5, "m1", "u1"); err == nil {
		t.Fatal("expected an error for an out-of-range task index")
	}
}

func TestConvertInvalidMachineUser(t *testing.T) {
	c := &Config{Tasks: []Task{{}}}

	if _, err := c.Convert(0, "missing", "user"); err == nil {
		t.Fatal("expected an error for an unknown machine/user combo")
	}
}

func TestResolvePath(t *testing.T) {
	c := &Config{configDir: "/base/dir"}

	tests := []struct {
		in   string
		want string
	}{
		{"/abs/path", "/abs/path"},
		{"rel/path", "/base/dir/rel/path"},
		{"./rel", "/base/dir/rel"},
	}

	for _, tt := range tests {
		if got := c.ResolvePath(tt.in); got != tt.want {
			t.Errorf("ResolvePath(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
