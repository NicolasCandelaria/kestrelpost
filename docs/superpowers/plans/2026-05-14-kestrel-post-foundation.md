# Kestrel Post foundation — Implementation plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Ship a runnable Go repo with a **pure ending resolver** matching `docs/superpowers/specs/2026-05-14-kestrel-post-ending-evaluator-design.md`, plus a **minimal Wish + Bubble Tea** SSH entrypoint that proves the stack (one session = one TUI).

**Architecture:** `internal/ending` holds `RunState`, tunable `Config`, and `ResolveEnding` with no I/O. `internal/ui` exposes a Bubble Tea `tea.Model` used only for layout/navigation stubs. `cmd/kestrelpost` hosts an SSH server via Wish and attaches the Wish Bubble Tea middleware so each connection runs an isolated program. Narrative and fuel simulation are **out of scope** for this plan; tests drive state by constructing `RunState` directly.

**Tech Stack:** Go 1.22+ (toolchain in `go.mod`). Wish and Bubble Tea **v2** currently ship from Charm’s `charm.land` modules (matches upstream `wish/examples/bubbletea`): `charm.land/wish/v2`, `charm.land/bubbletea/v2`, `charm.land/wish/v2/bubbletea`, `charm.land/wish/v2/activeterm`, `charm.land/wish/v2/logging`, and `github.com/charmbracelet/ssh` for `ssh.ErrServerClosed`. Bubble Tea v2 uses `tea.KeyMsg`, `View() tea.View`, and `tea.NewView(...)` with `v.AltScreen = true` instead of `tea.WithAltScreen()` in many examples.

---

## File structure (new repo)

| Path | Responsibility |
|------|------------------|
| `go.mod` / `go.sum` | Module `kestrelpost`, pins Charm deps. |
| `internal/ending/id.go` | `Ending` type + `String()` + numeric IDs 1–7 aligned with spec. |
| `internal/ending/state.go` | `RunState` (all spec fields), `Config` (`KMax`, `MThreshold`, `TThreshold`, `RelayMinTerminalNight`, `DeadAirExclusiveMaxTerminalNight`). |
| `internal/ending/resolve.go` | `ResolveEnding(cfg Config, s RunState) Ending` — pure cascade per priority list. |
| `internal/ending/resolve_test.go` | Table-driven tests for priority, Broadcast vs Kid, RELAY vs DEAD AIR, fallback. |
| `internal/ui/model.go` | Minimal `tea.Model`: stub `View()` as `tea.NewView(...)`, `v.AltScreen = true`, quit on `tea.KeyMsg` `q` / `ctrl+c`. |
| `cmd/kestrelpost/main.go` | `wish.NewServer`, host key path flag or auto, `:2222`, Bubble Tea middleware factory that returns `ui.NewModel()`. |
| `README.md` | How to run locally: `go run ./cmd/kestrelpost`, `ssh -p 2222 localhost` (known_hosts note from Wish README). |

---

### Task 1: Go module and ending IDs

**Files:**
- Create: `go.mod`
- Create: `internal/ending/id.go`

- [ ] **Step 1: Write `go.mod`**

```go
module kestrelpost

go 1.22
```

- [ ] **Step 2: Write `internal/ending/id.go`**

```go
package ending

import "strconv"

// Ending IDs match docs/superpowers/specs/2026-05-14-kestrel-post-ending-evaluator-design.md
type Ending uint8

const (
	TheRelay         Ending = 1
	DarkFrequency    Ending = 2
	TheKidWasRight   Ending = 3
	FullBroadcast    Ending = 4
	TheConvoy        Ending = 5
	DeadAir          Ending = 6
	Fallback         Ending = 7
)

func (e Ending) String() string {
	switch e {
	case TheRelay:
		return "THE_RELAY"
	case DarkFrequency:
		return "DARK_FREQUENCY"
	case TheKidWasRight:
		return "THE_KID_WAS_RIGHT"
	case FullBroadcast:
		return "FULL_BROADCAST"
	case TheConvoy:
		return "THE_CONVOY"
	case DeadAir:
		return "DEAD_AIR"
	case Fallback:
		return "FALLBACK"
	default:
		return "Ending(" + strconv.Itoa(int(e)) + ")"
	}
}
```

- [ ] **Step 3: Commit**

```bash
git add go.mod internal/ending/id.go
git commit -m "feat(ending): add module and Ending identifiers"
```

---

### Task 2: RunState, Config, and ResolveEnding

**Files:**
- Create: `internal/ending/state.go`
- Create: `internal/ending/resolve.go`

- [ ] **Step 1: Write `internal/ending/state.go`**

```go
package ending

// Config holds tunable thresholds; content can override per build or load later.
type Config struct {
	KMax int // kid_investigation_stage >= KMax counts as Kid payoff

	MThreshold int // maren_hub_support >= for THE_RELAY
	TThreshold int // maren_trust >= for THE_RELAY

	// RelayMinTerminalNight: e.g. 9 — fuel ran out on this night or later.
	RelayMinTerminalNight int

	// DeadAirExclusiveMaxTerminalNight: DEAD AIR when terminal_dark_night < this value (spec placeholder uses 7).
	DeadAirExclusiveMaxTerminalNight int
}

func DefaultConfig() Config {
	return Config{
		KMax:                             3,
		MThreshold:                       5,
		TThreshold:                       3,
		RelayMinTerminalNight:            9,
		DeadAirExclusiveMaxTerminalNight: 7,
	}
}

// RunState is the input to ResolveEnding; narrative systems mutate this over a run.
type RunState struct {
	Night               int
	Fuel                int
	TerminalDarkNight   int // night index when fuel first hit 0; 0 if not yet ended that way
	HarrowDarkPlan      bool
	KidInvestigation    int
	OseiFullRelease     bool
	ConvoyBetrayal      bool
	MarenHubSupport     int
	MarenTrust          int
}
```

- [ ] **Step 2: Write failing test file skeleton** (skip if you prefer test-after; here minimal compile check)

Create empty `internal/ending/resolve.go` with:

```go
package ending

func ResolveEnding(cfg Config, s RunState) Ending {
	return Fallback
}
```

- [ ] **Step 3: Run `go test ./internal/ending/...`**

Run: `go test ./internal/ending/...`

Expected: PASS (no tests yet) or compile OK.

- [ ] **Step 4: Replace `resolve.go` with full implementation**

```go
package ending

// ResolveEnding picks exactly one ending; first matching rule wins (spec order).
func ResolveEnding(cfg Config, s RunState) Ending {
	if s.ConvoyBetrayal {
		return TheConvoy
	}
	if s.OseiFullRelease {
		return FullBroadcast
	}
	if s.HarrowDarkPlan {
		return DarkFrequency
	}
	if s.KidInvestigation >= cfg.KMax && !s.OseiFullRelease {
		return TheKidWasRight
	}
	if s.Fuel <= 0 && s.TerminalDarkNight > 0 && s.TerminalDarkNight < cfg.DeadAirExclusiveMaxTerminalNight && !s.ConvoyBetrayal {
		return DeadAir
	}
	if s.Fuel <= 0 &&
		s.TerminalDarkNight >= cfg.RelayMinTerminalNight &&
		s.MarenHubSupport >= cfg.MThreshold &&
		s.MarenTrust >= cfg.TThreshold {
		return TheRelay
	}
	return Fallback
}
```

- [ ] **Step 5: Commit**

```bash
git add internal/ending/state.go internal/ending/resolve.go
git commit -m "feat(ending): add RunState, Config, and ResolveEnding cascade"
```

---

### Task 3: Table-driven tests (spec coverage)

**Files:**
- Create: `internal/ending/resolve_test.go`

- [ ] **Step 1: Write `internal/ending/resolve_test.go`**

```go
package ending

import "testing"

func TestResolveEnding_priorityConvoyOverBroadcast(t *testing.T) {
	cfg := DefaultConfig()
	s := RunState{
		ConvoyBetrayal:  true,
		OseiFullRelease: true,
		KidInvestigation: cfg.KMax + 1,
	}
	if g := ResolveEnding(cfg, s); g != TheConvoy {
		t.Fatalf("got %v want TheConvoy", g)
	}
}

func TestResolveEnding_broadcastBeatsKid(t *testing.T) {
	cfg := DefaultConfig()
	s := RunState{
		OseiFullRelease:  true,
		KidInvestigation: cfg.KMax,
	}
	if g := ResolveEnding(cfg, s); g != FullBroadcast {
		t.Fatalf("got %v want FullBroadcast", g)
	}
}

func TestResolveEnding_kidWhenNoBroadcast(t *testing.T) {
	cfg := DefaultConfig()
	s := RunState{
		KidInvestigation: cfg.KMax,
	}
	if g := ResolveEnding(cfg, s); g != TheKidWasRight {
		t.Fatalf("got %v want TheKidWasRight", g)
	}
}

func TestResolveEnding_darkFrequencyAfterKidBlock(t *testing.T) {
	cfg := DefaultConfig()
	s := RunState{
		HarrowDarkPlan:   true,
		KidInvestigation: cfg.KMax,
	}
	// Harrow checked before Kid; dark wins here
	if g := ResolveEnding(cfg, s); g != DarkFrequency {
		t.Fatalf("got %v want DarkFrequency", g)
	}
}

func TestResolveEnding_deadAirEarlyTerminalNight(t *testing.T) {
	cfg := DefaultConfig()
	s := RunState{
		Fuel:              0,
		TerminalDarkNight: cfg.DeadAirExclusiveMaxTerminalNight - 1,
	}
	if g := ResolveEnding(cfg, s); g != DeadAir {
		t.Fatalf("got %v want DeadAir", g)
	}
}

func TestResolveEnding_relayLateWithMarenScores(t *testing.T) {
	cfg := DefaultConfig()
	s := RunState{
		Fuel:              0,
		TerminalDarkNight: cfg.RelayMinTerminalNight,
		MarenHubSupport:   cfg.MThreshold,
		MarenTrust:        cfg.TThreshold,
	}
	if g := ResolveEnding(cfg, s); g != TheRelay {
		t.Fatalf("got %v want TheRelay", g)
	}
}

func TestResolveEnding_relayFailsLowTrust(t *testing.T) {
	cfg := DefaultConfig()
	s := RunState{
		Fuel:              0,
		TerminalDarkNight: cfg.RelayMinTerminalNight,
		MarenHubSupport:   cfg.MThreshold,
		MarenTrust:        cfg.TThreshold - 1,
	}
	if g := ResolveEnding(cfg, s); g != Fallback {
		t.Fatalf("got %v want Fallback", g)
	}
}

func TestResolveEnding_relayFailsLowHub(t *testing.T) {
	cfg := DefaultConfig()
	s := RunState{
		Fuel:              0,
		TerminalDarkNight: cfg.RelayMinTerminalNight,
		MarenHubSupport:   cfg.MThreshold - 1,
		MarenTrust:        cfg.TThreshold,
	}
	if g := ResolveEnding(cfg, s); g != Fallback {
		t.Fatalf("got %v want Fallback", g)
	}
}
```

- [ ] **Step 2: Run tests**

Run: `go test ./internal/ending/... -v`

Expected: all PASS.

- [ ] **Step 3: Commit**

```bash
git add internal/ending/resolve_test.go
git commit -m "test(ending): cover priority, broadcast vs kid, relay thresholds"
```

---

### Task 4: Minimal Bubble Tea model

**Files:**
- Create: `internal/ui/model.go`

- [ ] **Step 1: Add Wish (pulls Bubble Tea v2 transitively)**

Run:

```bash
go get charm.land/wish/v2@latest
```

- [ ] **Step 2: Write `internal/ui/model.go`**

```go
package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type Model struct{}

func NewModel() Model { return Model{} }

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	s := "KESTREL POST\n\nRelay automation offline.\n\nPress q to disconnect.\n"
	v := tea.NewView(s)
	v.AltScreen = true
	return v
}

func (m Model) String() string { return fmt.Sprintf("%T", m) }
```

- [ ] **Step 3: Verify build**

Run: `go build -o NUL ./...` on Windows, or `go build ./...` on Unix.

Expected: success.

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum internal/ui/model.go
git commit -m "feat(ui): add minimal Bubble Tea v2 stub model"
```

---

### Task 5: Wish SSH server + middleware

**Files:**
- Create: `cmd/kestrelpost/main.go`

- [ ] **Step 1: Ensure Wish v2 is in `go.mod`**

If Task 4 did not add it yet, run:

```bash
go get charm.land/wish/v2@latest
```

Reference implementation: [wish/examples/bubbletea/main.go](https://github.com/charmbracelet/wish/blob/main/examples/bubbletea/main.go).

- [ ] **Step 2: Write `cmd/kestrelpost/main.go`**

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/charmbracelet/ssh"
	"kestrelpost/internal/ui"
)

func main() {
	hostKeyPath := os.Getenv("KESTREL_HOST_KEY")
	if hostKeyPath == "" {
		hostKeyPath = ".ssh/kestrel_ed25519"
	}

	s, err := wish.NewServer(
		wish.WithAddress(":2222"),
		wish.WithHostKeyPath(hostKeyPath),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("SSH listening on :2222")
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Printf("server stopped: %v", err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Println("stopping SSH server")
	shCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(shCtx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Printf("shutdown: %v", err)
	}
	fmt.Fprintln(os.Stderr, "kestrelpost stopped")
}

func teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	m := ui.NewModel()
	return m, nil
}
```

**Notes:**

1. **`activeterm.Middleware()`** is required for Bubble Tea over SSH in the upstream example (PTY).
2. **`logging.Middleware()`** uses Wish’s logging helper; if you prefer std `log` only, you may remove it later—keep it for parity with Charm’s example first.
3. Host key: document `ssh-keygen -t ed25519 -f .ssh/kestrel_ed25519 -N ""` in README if the file is missing.

- [ ] **Step 3: Run server briefly**

Run: `go run ./cmd/kestrelpost`

In another terminal: `ssh -p 2222 -o StrictHostKeyChecking=no localhost` (dev only).

Expected: TUI shows stub; `q` exits session.

- [ ] **Step 4: Commit**

```bash
git add cmd/kestrelpost/main.go go.mod go.sum
git commit -m "feat(cmd): add Wish SSH server with Bubble Tea middleware"
```

---

### Task 6: README and host key note

**Files:**
- Create: `README.md`

- [ ] **Step 1: Write `README.md`**

```markdown
# Kestrel Post

SSH narrative game (Wish + Bubble Tea). Design specs live under `docs/superpowers/`.

## Run (development)

Generate a host key once:

```bash
mkdir -p .ssh
ssh-keygen -t ed25519 -f .ssh/kestrel_ed25519 -N ""
```

Start the server:

```bash
go run ./cmd/kestrelpost
```

Connect:

```bash
ssh -p 2222 localhost
```

Optional: add to `~/.ssh/config`:

```
Host kestrel-local
  HostName localhost
  Port 2222
  UserKnownHostsFile /dev/null
  StrictHostKeyChecking no
```

## Tests

```bash
go test ./...
```

## Ending resolver

Pure logic in `internal/ending` matches `docs/superpowers/specs/2026-05-14-kestrel-post-ending-evaluator-design.md`.
```

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add README for local SSH dev and tests"
```

---

## Plan self-review

**1. Spec coverage**

| Spec requirement | Task |
|------------------|------|
| State fields (`night`, `fuel`, `terminal_dark_night`, flags, Maren hub + trust) | Task 2 `RunState` |
| Priority list 1–7 | Task 2 `ResolveEnding` order |
| Broadcast beats Kid | Task 3 tests + step 2 order |
| THE RELAY needs hub **and** trust | Task 3 `TestResolveEnding_relayLateWithMarenScores` + low hub/trust tests |
| DEAD AIR early terminal night | Task 3 `TestResolveEnding_deadAirEarlyTerminalNight` |
| Pure function, no side effects | `ResolveEnding` only reads args |
| Fallback | Default return + test for low trust relay miss |

**Gap (intentional for this slice):** `night` is stored on `RunState` but not used inside `ResolveEnding` yet — narrative sim will use it later. `harrow_dark_plan` is a single bool; “sustained window” is content’s responsibility to set the bool only when appropriate. **Priority edge case:** if `KidInvestigation >= KMax` and early `TerminalDarkNight` (would look like DEAD AIR timing), the resolver still returns **THE KID WAS RIGHT** because Kid is checked before DEAD AIR—matches the written spec order; add a narrative test later if you want to forbid that combination in the sim.

**2. Placeholder scan:** No TBD steps; pins follow current upstream `charm.land/wish/v2` + `charm.land/bubbletea/v2` (re-run `go get` if Charm re-tags).

**3. Type consistency:** `Ending` constants, `RunState` field names, and `Config` fields match across `state.go`, `resolve.go`, and tests.

---

**Plan complete and saved to `docs/superpowers/plans/2026-05-14-kestrel-post-foundation.md`. Two execution options:**

**1. Subagent-Driven (recommended)** — Dispatch a fresh subagent per task, review between tasks, fast iteration.

**2. Inline Execution** — Execute tasks in this session using executing-plans, batch execution with checkpoints.

**Which approach do you want?**
