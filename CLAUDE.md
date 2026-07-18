# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Keeping this file current

The architecture is actively evolving. Whenever a change alters package boundaries, responsibilities, or the event/data flow described below, update this file in the same session — don't let it drift out of sync with the code.

## What this is

gosible is a terminal (bubbletea) reinterpretation of Ansible for homelab management, driven by a single JSON config file passed as an argument (`gosible config.json`).

## Architecture (strict layering — respect the dependency direction)

- `config` — JSON parsing, validation, `env(VAR)`/relative-path resolution. Imports `action` and `runner`.
- `action` — the `Args` interface (`Validate`, `Prepare`, `Run`) and each action type (copy, dir, docker, file, shell). Must not import `config`: `config` already imports `action`, so a reverse import creates a cycle. Context-dependent resolution (paths, env vars) is threaded in as plain functions via `Prepare(resolvePath, replaceEnv func(string) string)`, not by depending on `config.Config` directly.
- `runner` — SSH orchestration only (`Runner`, `Machine`, `sshExecutor`). Implements `action.Executor` on top of `ssh.Session`. Emits its own `runner.Event` types over a channel and must not import `bubbletea`; `ui` is responsible for translating events into `tea.Msg`.
- `ui` — the bubbletea TUI (tasks → machines → logs screens).

`action.Args`'s `Validate`/`Prepare`/`Run` split matters: `Validate()` is pure structural validation with no external dependency; `Prepare()` applies context (path resolution, env substitution) via the injected functions; `Run()` executes against an `Executor`. Don't merge these back together — that split is what removed the import cycle this project used to have between `config` and `runner`.

## Commands

- Build: `make build` (binary at `build/gosible`), or `go run ./cmd/main.go config.json` for local dev.
- Test: `go test ./...` (add `-race` for the full suite). No external tools/services needed: `action`/`config` tests use in-memory fakes, and `runner` tests spin up a pure-Go in-process fake SSH server (`runner/testserver_test.go`) instead of a real sshd.
- Before considering work done: `go vet ./...` and `gofmt -l .` should both be clean.

## Collaboration

Never commit or push automatically — leave changes as diffs in the working tree for review; the user commits and pushes themselves.
