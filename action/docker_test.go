package action

import "testing"

func TestDockerArgsValidateDefaultsTarName(t *testing.T) {
	a := &DockerArgs{Image: "myimage:latest"}

	if err := a.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Tar != "image.tar" {
		t.Errorf("got tar=%q, want %q", a.Tar, "image.tar")
	}
}

func TestDockerArgsValidateRejectsNonTarSuffix(t *testing.T) {
	a := &DockerArgs{Tar: "image.zip"}

	if err := a.Validate(); err == nil {
		t.Fatal("expected an error for a non-.tar suffix")
	}
}

func TestDockerArgsValidateKeepsExplicitTarName(t *testing.T) {
	a := &DockerArgs{Tar: "custom.tar"}

	if err := a.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Tar != "custom.tar" {
		t.Errorf("got tar=%q, want %q", a.Tar, "custom.tar")
	}
}

func TestDockerArgsPrepareResolvesSrc(t *testing.T) {
	a := &DockerArgs{Src: "relative/dir"}

	resolvePath := func(p string) string { return "/resolved/" + p }
	if err := a.Prepare(resolvePath, identityReplaceEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "/resolved/relative/dir"
	if a.Src != want {
		t.Errorf("got src=%q, want %q", a.Src, want)
	}
}
