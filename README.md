# Kestrel Post

SSH narrative game (Wish + Bubble Tea v2). Design docs live under `docs/superpowers/`.

## Prerequisites

- [Go](https://go.dev/dl/) — use a recent stable release; the `go` line in `go.mod` is the minimum version the Charm stack requested when last tidied.

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

### Vertical slice (in-game)

1. Press any key (except `q`) on the intro to start.
2. Each **night**, pick **1**, **2**, or **3** to spend fuel and adjust Maren hub/trust placeholders.
3. When fuel hits zero before night 9, or after you complete the nine-night runway, the run ends and shows the **resolved ending** from `internal/ending` (e.g. `DEAD_AIR` vs `THE_RELAY`).

Optional: add to `~/.ssh/config` (see [Wish README](https://github.com/charmbracelet/wish) for `UserKnownHostsFile` tips during local dev).

## Tests

```bash
go test ./...
```

## Ending resolver

Pure logic in `internal/ending` implements `docs/superpowers/specs/2026-05-14-kestrel-post-ending-evaluator-design.md`.

## Session loop

`internal/game` advances nights and fuel, then calls `ending.ResolveEnding` on game over. Replace the stub transmission and choice deltas with real content when you add narrative data.
