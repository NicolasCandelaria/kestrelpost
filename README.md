# Kestrel Post

SSH narrative game (Wish + Bubble Tea v2). Design docs live under `docs/superpowers/`.

## Prerequisites

- [Go](https://go.dev/dl/) 1.23 or newer (matches `charm.land/wish/v2` toolchain expectations).

## Run (development)

Generate a host key once:

```bash
mkdir -p .ssh
ssh-keygen -t ed25519 -f .ssh/kestrel_ed25519 -N ""
```

On Windows PowerShell:

```powershell
New-Item -ItemType Directory -Force .ssh | Out-Null
ssh-keygen -t ed25519 -f .ssh/kestrel_ed25519 -N '""'
```

Sync modules (writes `go.sum`):

```bash
go mod tidy
```

Start the server:

```bash
go run ./cmd/kestrelpost
```

Connect:

```bash
ssh -p 2222 localhost
```

Optional: add to `~/.ssh/config` (see [Wish README](https://github.com/charmbracelet/wish) for `UserKnownHostsFile` tips during local dev).

## Tests

```bash
go test ./...
```

## Ending resolver

Pure logic in `internal/ending` implements `docs/superpowers/specs/2026-05-14-kestrel-post-ending-evaluator-design.md`.
