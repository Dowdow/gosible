package runner_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Dowdow/gosible/action"
	"github.com/Dowdow/gosible/runner"
)

// runAndCollect runs cfg through a Runner and collects every event until a
// Done or Failed event is seen (Runner.Run never closes the channel itself;
// that's left to the consumer, same as ui/logs.go does).
func runAndCollect(t *testing.T, cfg *runner.Config) []runner.Event {
	t.Helper()

	r := runner.NewRunner(cfg)
	ch := make(chan runner.Event)

	var events []runner.Event
	done := make(chan struct{})

	go func() {
		for e := range ch {
			events = append(events, e)
			switch e.(type) {
			case runner.Done, runner.Failed:
				close(done)
			}
		}
	}()

	go r.Run(ch)

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for the runner to finish")
	}

	return events
}

func assertNoFailure(t *testing.T, events []runner.Event) {
	t.Helper()

	for _, e := range events {
		if f, ok := e.(runner.Failed); ok {
			t.Fatalf("unexpected failure: %v", f.Err)
		}
	}

	if len(events) == 0 {
		t.Fatal("expected at least one event")
	}
	if _, ok := events[len(events)-1].(runner.Done); !ok {
		t.Fatalf("expected the last event to be Done, got %T", events[len(events)-1])
	}
}

func countEvents[T any](events []runner.Event) int {
	n := 0
	for _, e := range events {
		if _, ok := e.(T); ok {
			n++
		}
	}
	return n
}

func TestRunnerRunSuccess(t *testing.T) {
	srv := newTestSSHServer(t, "")

	cfg := &runner.Config{
		Machine: runner.Machine{Address: srv.addr(), User: "tester", Password: "x"},
		Actions: []runner.Action{
			{Name: "say hi", Args: action.ShellArgs("echo hi")},
			{Name: "make dir", Args: &action.DirArgs{Path: "/tmp/x"}},
		},
	}

	events := runAndCollect(t, cfg)
	assertNoFailure(t, events)

	if n := countEvents[runner.ActionStarted](events); n != 2 {
		t.Errorf("expected 2 ActionStarted events, got %d", n)
	}
	if n := countEvents[runner.ActionCompleted](events); n != 2 {
		t.Errorf("expected 2 ActionCompleted events, got %d", n)
	}
	for _, e := range events {
		if ac, ok := e.(runner.ActionCompleted); ok && !ac.Success {
			t.Errorf("expected all actions to succeed, got %+v", ac)
		}
	}
}

func TestRunnerRunActionFailureStopsEarly(t *testing.T) {
	srv := newTestSSHServer(t, "")

	cfg := &runner.Config{
		Machine: runner.Machine{Address: srv.addr(), User: "tester", Password: "x"},
		Actions: []runner.Action{
			{Name: "boom", Args: action.ShellArgs("false-command")},
			{Name: "never runs", Args: action.ShellArgs("echo hi")},
		},
	}

	events := runAndCollect(t, cfg)

	if n := countEvents[runner.ActionStarted](events); n != 1 {
		t.Fatalf("expected only the first action to start, got %d ActionStarted events", n)
	}
	if _, ok := events[len(events)-1].(runner.Failed); !ok {
		t.Fatalf("expected the last event to be Failed, got %T", events[len(events)-1])
	}
}

func TestRunnerRunCopySingleFile(t *testing.T) {
	srv := newTestSSHServer(t, "")

	srcDir := t.TempDir()
	srcPath := filepath.Join(srcDir, "hello.txt")
	if err := os.WriteFile(srcPath, []byte("hello from test"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &runner.Config{
		Machine: runner.Machine{Address: srv.addr(), User: "tester", Password: "x"},
		Actions: []runner.Action{
			{Name: "copy", Args: &action.CopyArgs{Src: srcPath, Dest: "/remote/uploaded/hello.txt"}},
		},
	}

	events := runAndCollect(t, cfg)
	assertNoFailure(t, events)

	uploaded := filepath.Join(srv.root, "remote", "uploaded", "hello.txt")
	got, err := os.ReadFile(uploaded)
	if err != nil {
		t.Fatalf("expected uploaded file at %s: %v", uploaded, err)
	}
	if string(got) != "hello from test" {
		t.Errorf("got content %q, want %q", got, "hello from test")
	}
}

func TestRunnerRunCopyDirectory(t *testing.T) {
	srv := newTestSSHServer(t, "")

	srcDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(srcDir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "a.txt"), []byte("A"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "sub", "b.txt"), []byte("B"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &runner.Config{
		Machine: runner.Machine{Address: srv.addr(), User: "tester", Password: "x"},
		Actions: []runner.Action{
			{Name: "copy dir", Args: &action.CopyArgs{Src: srcDir, Dest: "/remote/mydir"}},
		},
	}

	events := runAndCollect(t, cfg)
	assertNoFailure(t, events)

	uploadedRoot := filepath.Join(srv.root, "remote", "mydir")
	for path, want := range map[string]string{
		filepath.Join(uploadedRoot, "a.txt"):        "A",
		filepath.Join(uploadedRoot, "sub", "b.txt"): "B",
	} {
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("expected uploaded file at %s: %v", path, err)
		}
		if string(got) != want {
			t.Errorf("got content %q for %s, want %q", got, path, want)
		}
	}
}

func TestRunnerRunWithPasswordAuth(t *testing.T) {
	srv := newTestSSHServer(t, "s3cret")

	cfg := &runner.Config{
		Machine: runner.Machine{Address: srv.addr(), User: "tester", Password: "s3cret"},
		Actions: []runner.Action{{Name: "hi", Args: action.ShellArgs("echo hi")}},
	}

	assertNoFailure(t, runAndCollect(t, cfg))
}

func TestRunnerRunWithWrongPasswordFails(t *testing.T) {
	srv := newTestSSHServer(t, "s3cret")

	cfg := &runner.Config{
		Machine: runner.Machine{Address: srv.addr(), User: "tester", Password: "wrong"},
		Actions: []runner.Action{{Name: "hi", Args: action.ShellArgs("echo hi")}},
	}

	events := runAndCollect(t, cfg)
	if _, ok := events[len(events)-1].(runner.Failed); !ok {
		t.Fatalf("expected the last event to be Failed, got %T", events[len(events)-1])
	}
}

func TestRunnerRunWithSSHKeyAuth(t *testing.T) {
	srv := newTestSSHServer(t, "")

	keyPath := filepath.Join(t.TempDir(), "id_rsa")
	writeTestPrivateKey(t, keyPath)

	cfg := &runner.Config{
		Machine: runner.Machine{Address: srv.addr(), User: "tester", Ssh: keyPath},
		Actions: []runner.Action{{Name: "hi", Args: action.ShellArgs("echo hi")}},
	}

	assertNoFailure(t, runAndCollect(t, cfg))
}

func TestRunnerRunDialFailure(t *testing.T) {
	// A listener that's immediately closed gives us a local address
	// guaranteed to refuse connections.
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	l.Close()

	cfg := &runner.Config{
		Machine: runner.Machine{Address: addr, User: "tester", Password: "x"},
		Actions: []runner.Action{{Name: "hi", Args: action.ShellArgs("echo hi")}},
	}

	events := runAndCollect(t, cfg)
	if len(events) != 1 {
		t.Fatalf("expected exactly one event, got %d: %v", len(events), events)
	}
	if _, ok := events[0].(runner.Failed); !ok {
		t.Fatalf("expected Failed, got %T", events[0])
	}
}

func writeTestPrivateKey(t *testing.T, path string) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generating client key: %v", err)
	}

	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}
	if err := os.WriteFile(path, pem.EncodeToMemory(block), 0o600); err != nil {
		t.Fatalf("writing client key: %v", err)
	}
}
