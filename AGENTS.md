# AGENTS.md

## Commands

```sh
go build ./cmd/...          # build
go run ./cmd <c|cxx> [flags] <project-name>  # run
go test ./...               # test (no tests exist yet; passes vacuously)
go vet ./...                # vet/lint
gofmt -w .                  # format
```

No Makefile, taskfile, external dependencies, or CI. `go mod tidy` does nothing useful — the module uses only stdlib.

## go.mod version

`go.mod` declares `go 1.26.0` (pre-release at time of writing). Build will fail on older toolchains. Do not downgrade without auditing stdlib feature usage.

## Architecture

- `cmd/main.go` — dispatches `c` → `cproject`, `cxx` → stub (not implemented)
- `cmd/cproject/` — all live logic: flag parsing, template definitions, project scaffolding
- `internal/folder/folder.go` — `CreateFolder()`: creates a dir, writes templates + copies embedded files
- `cmd/cproject/files/` — embedded static files (`.clang-format`, `.gitignore`, `tests.c`) via `//go:embed`
- `cmd/cxxproject/cxxproject.go` — one-line stub, nothing implemented

## Critical quirks

**CWD mutation**: `folder.CreateFolder()` calls `os.Chdir(folderName)`. The process working directory changes with each folder created. Any new `createXxx()` method must track where CWD is. See how `createTests()` saves/restores root with `os.Chdir(p.rootFolder)`.

**Tests write to filesystem**: `CProject.Run()` creates real directories relative to CWD. Any Go tests for it must use a temp directory.

**`--exe-flags` is broken**: `StringSliceValue.Set()` assigns to a value receiver, so the flag never propagates. It always uses the default `["-Wall", "-Werror", "-pedantic"]`.

**Embed glob**: Adding new static files to `cmd/cproject/files/` requires updating the `//go:embed` directive in `cproject.go` if the new path isn't already covered.

## Planned but unimplemented (see TODOS.md)

Atomic project creation, `--dry-run`, `.clang-tidy`, lib-only folder option, C++ support, Linux kernel module support.
