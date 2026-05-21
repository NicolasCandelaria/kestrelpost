# Kestrel Post

A story game you play in your terminal over SSH. You are Operator 7 at a northern relay post, answering radio traffic night by night until the run ends.

**Repository:** [github.com/NicolasCandelaria/kestrelpost](https://github.com/NicolasCandelaria/kestrelpost)

---

## What you need

Install these once on the computer that will **run the game**:

| Tool | Get it |
|------|--------|
| **Go** | [go.dev/dl](https://go.dev/dl/) |
| **Git** | [git-scm.com](https://git-scm.com/downloads) |
| **SSH client** | Usually already installed (Terminal on Mac/Linux; Windows 10/11 includes OpenSSH) |

You do **not** need to install anything extra on a second computer if you only **connect** to a game that is already running on another machine (see [Play from another computer on your network](#play-from-another-computer-on-your-network)).

---

## Quick start (same computer)

Use **two terminal windows** on one machine: one runs the server, one plays the game.

### Step 1 — Get the code

**Mac / Linux**

```bash
git clone https://github.com/NicolasCandelaria/kestrelpost.git
cd kestrelpost
go mod tidy
```

**Windows (PowerShell)**

```powershell
git clone https://github.com/NicolasCandelaria/kestrelpost.git
Set-Location kestrelpost
go mod tidy
```

### Step 2 — Create an SSH host key (first time only)

The server needs a key file in the project folder. Run this **once** inside `kestrelpost`:

**Mac / Linux**

```bash
mkdir -p .ssh
ssh-keygen -t ed25519 -f .ssh/kestrel_ed25519 -N ""
```

**Windows (PowerShell)**

```powershell
New-Item -ItemType Directory -Force .ssh | Out-Null
ssh-keygen -t ed25519 -f .ssh/kestrel_ed25519 -N '""'
```

If asked to overwrite, you can say yes only if you are setting up fresh.

### Step 3 — Start the server

Leave this window **open**:

```bash
go run ./cmd/kestrelpost
```

You should see something like: `SSH listening on :2222`

### Step 4 — Play (second terminal)

Open a **new** terminal, then connect:

```bash
ssh -p 2222 localhost
```

The first time, your SSH client may warn about an unknown host key — that is normal for local play. Accept it (type `yes` if prompted).

To quit the game: press **`q`**. To stop the server: go to the server window and press **Ctrl+C**.

---

## Play on a different computer

### Same machine, different day

Repeat **Quick start** on that computer: clone (or pull latest), `go mod tidy`, start the server, then `ssh -p 2222 localhost`.

### Play from another computer on your network

1. On the **host** (the PC running the game), do **Steps 1–3** above and note its local IP address.
   - **Windows:** `ipconfig` → look for **IPv4 Address** (e.g. `192.168.1.25`)
   - **Mac:** System Settings → Network, or `ipconfig getifaddr en0`
2. Allow port **2222** through the host firewall if connections fail.
3. On the **other** computer, open a terminal and run (replace with the host IP):

```bash
ssh -p 2222 192.168.1.25
```

Both computers must be on the same Wi‑Fi or LAN. This does **not** work over the public internet unless you set up port forwarding or a tunnel yourself.

---

## How to play

1. Intro screen: press **`t`** for Tutorial Night 0, or **`enter`** to start the campaign.
2. The campaign runs up to **20 nights** with a fixed night structure:
   - **Pre-shift:** `enter` open radio, `f` feed dog, `p/u` pin or unpin thread, `l` read logbook note
   - **Receive:** `1`, `2`, `3` choose response, `s` open scan
   - **Scan:** arrows tune, `1-4` change band, `enter` lock signal, `r` return to receive
   - **Incident / Logbook:** `enter` continue, `w` write note, `n` advance night
3. On Night 1 pre-shift, name your dog first: **`1` Scout**, **`2` Ash**, or **`3` Bramble**.
4. Save/load snapshot in pre-shift with **`k`** (save) and **`o`** (load).
5. Endings resolve by hidden state (trust, hub support, Kid investigation, Harrow commitment, release/betrayal flags, and fuel runway). Press **`q`** to quit.

The UI keeps resources diegetic: top bar shows gauge state words while the story panel uses rig cues and log lines.

---

## Troubleshooting

| Problem | What to try |
|---------|-------------|
| `go: command not found` | Install Go and open a **new** terminal after installing. |
| `ssh: command not found` (Windows) | Install [OpenSSH Client](https://learn.microsoft.com/en-us/windows-server/administration/openssh/openssh_install_firstuse) or use Git Bash. |
| `Connection refused` on port 2222 | Start the server first (`go run ./cmd/kestrelpost`) and keep that window open. |
| `no such file` for host key | Run **Step 2** (host key) inside the `kestrelpost` folder. |
| Cannot connect from another PC | Check IP address, same network, and firewall on the host for port **2222**. |

---

## Optional settings

| Variable | Default | Purpose |
|----------|---------|---------|
| `KESTREL_HOST_KEY` | `.ssh/kestrel_ed25519` | Path to the server host key file |
| `KESTREL_LISTEN` | `:2222` | Address/port to listen on (e.g. `:22` only if nothing else uses port 22) |

Example (Mac / Linux):

```bash
export KESTREL_LISTEN=":2222"
go run ./cmd/kestrelpost
```

---

## For developers

- **Tests:** `go test ./...`
- **Design:** `docs/superpowers/specs/2026-05-14-kestrel-post-ending-evaluator-design.md`
- **Stack:** [Wish](https://github.com/charmbracelet/wish) + [Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Lip Gloss](https://github.com/charmbracelet/lipgloss); UI chrome inspired by a [terminal.shop-style layout](https://github.com/IsaiahPapa/terminal.shop).
- **Story / nights:** `internal/content/data/*.yaml` (20-night campaign), `internal/game/session.go` (night FSM)
- **Endings:** `internal/ending` (pure resolver + priority tests)
