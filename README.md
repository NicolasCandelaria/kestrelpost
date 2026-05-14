# Kestrel Post

SSH narrative game (Wish + Bubble Tea v2). The in-game UI uses a **terminal.shop–style** chrome (grid header, framed body, shortcut footer); layout is [Lip Gloss](https://github.com/charmbracelet/lipgloss) v2. Aesthetic nod: [terminal.shop clone (Rust)](https://github.com/IsaiahPapa/terminal.shop). Design docs live under `docs/superpowers/`.

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

1. Press almost any key (except `q`) on the intro to start.
2. Each **night**, pick **1**, **2**, or **3** to answer the traffic on the band. Your choices shape what survives the event—without spelling out the machinery on screen.
3. **Nine nights** move from triage through pressure (Kid, Harrow) into the endgame (Cole, Osei, Maren, convoy).
4. When the run ends, you see a **resolved ending** from `internal/ending` with a short epilogue. Reserve is shown only in the top bar (`reserve` tab).

Optional: add to `~/.ssh/config` (see [Wish README](https://github.com/charmbracelet/wish) for `UserKnownHostsFile` tips during local dev).

## Tests

```bash
go test ./...
```

## Ending resolver

Pure logic in `internal/ending` implements `docs/superpowers/specs/2026-05-14-kestrel-post-ending-evaluator-design.md`.

## Session loop

`internal/game` advances nights and fuel using **`script.go`** (nine `NightCard`s across three acts). `Session.ApplyChoice` updates `ending.RunState` flags (Harrow, Osei release, convoy betrayal, kid steps) then calls `ending.ResolveEnding` on game over. Replace or extend `NightScript` to grow the story.
