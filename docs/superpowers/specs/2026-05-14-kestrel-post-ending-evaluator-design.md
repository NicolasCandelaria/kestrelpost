# Kestrel Post — Ending evaluator (design spec)

**Status:** Approved (2026-05-14)  
**Scope:** Frozen rules for resolving which of six epilogues fires at run end. Implementation is a single evaluator over run state.

## Purpose

At end of run (fuel exhausted, scripted stop, or extraction), the game picks **exactly one** ending. Conditions may overlap; **priority order** defines the winner.

## Minimum state variables

| Symbol | Type | Meaning |
|--------|------|--------|
| `night` | int ≥ 1 | Current or final night index (pacing uses ~9 as horizon). |
| `fuel` | int or float | Remaining generator fuel; at 0 the run ends unless convoy extraction overrides narrative timing. |
| `terminal_dark_night` | int | Night when fuel first hit 0 (or generator lockout), if applicable. |
| `harrow_dark_plan` | bool | Player committed to Harrow’s closed-network / “go dark on broad frequencies” strategy (sticky once true for enough nights—exact threshold left to narrative content). |
| `kid_investigation_stage` | int 0…N | Progress on The Kid’s clue chain; `N = K_MAX` means payoff reached. |
| `osei_full_release` | bool | Dr. Osei recording pushed to **open / wide** channels (not private listen-only). |
| `convoy_betrayal` | bool | Player traded Maren’s settlement + supply cache (or equivalent) for convoy extraction. |
| `maren_hub_support` | int or float | Operational support for Maren’s site as a hub: routing survivors toward her, sharing **accurate** supply / relay intel, actions that grow the settlement’s odds **without** measuring whether she still believes you personally. |
| `maren_trust` | int or float | Credibility and emotional contract with Maren. Rises with honest, costly truths; falls when you lie, omit, or mislead—especially on a **caught** lie beat (must apply a defined penalty in content rules so it is not flavour-only). |

Constants `K_MAX`, `M_THRESHOLD`, and `T_THRESHOLD` are content-tuned; this spec only requires they exist.

### Maren: trust vs hub support

**Problem:** `maren_hub_support` alone can absorb “you lied and she found out” only if every narrative beat explicitly maps to a **documented** change to that same score. Otherwise writers add a powerful scene with no mechanical sink.

**Model (approved):** Use **two** fields always:

- **`maren_hub_support`** — *What you did for the hub*: routing survivors toward her, **accurate** supply / relay intel, actions that grow the settlement’s odds. Not “does she believe you.”
- **`maren_trust`** — *Whether she still accepts you as a reliable operator*. Caught lies, evasions, and withheld medical intel apply **trust** deltas defined in content (e.g. caught lie = −Δ large enough to matter). Hub can rise while trust falls. Use trust for dialogue gates, later-night availability, epilogue tiers, and **THE RELAY** eligibility alongside hub.

**THE RELAY:** Requires **both** `maren_hub_support >= M_THRESHOLD` and `maren_trust >= T_THRESHOLD`. For softer lies-only impact inside the same ending ID, lower `T_THRESHOLD` or add epilogue branching on `maren_trust` tiers without changing the evaluator’s ending ID.

**Convoy:** `convoy_betrayal` implies operational betrayal; set `maren_trust` to minimum (or a `maren_trust_broken` flag) when the deal resolves so downstream systems stay consistent.

## Ending definitions (logic)

- **THE CONVOY (5)** — Extraction via betrayal: `convoy_betrayal` is true and resolves successfully for this run.
- **FULL BROADCAST (4)** — `osei_full_release` is true.
- **DARK FREQUENCY (2)** — `harrow_dark_plan` is true for the sustained window defined by content, and run does not end on a higher-priority ending first.
- **THE KID WAS RIGHT (3)** — `kid_investigation_stage >= K_MAX` and **not** `osei_full_release`.
- **DEAD AIR (6)** — Early failure: fuel exhausted “too soon” (e.g. `fuel <= 0` and `terminal_dark_night < 7`; exact cutoff is pacing), and **not** convoy extraction.
- **THE RELAY (1)** — Bittersweet default: late runway (e.g. fuel hits 0 at `night >= 9` or end-of-night-9 per discrete model), `maren_hub_support >= M_THRESHOLD`, `maren_trust >= T_THRESHOLD`, and none of the exclusive higher-priority endings apply.

**RELAY vs DEAD AIR:** Same mechanism (fuel hits zero); split by **when** it happens plus **hub + trust** thresholds and absence of other ending flags.

## Resolution order (first match wins)

Irreversible / player-chosen outcomes trump ambient mystery and default fate.

1. **THE CONVOY** — `convoy_betrayal` resolved successfully.
2. **FULL BROADCAST** — `osei_full_release`.
3. **DARK FREQUENCY** — sustained `harrow_dark_plan` (and not already matched above).
4. **THE KID WAS RIGHT** — `kid_investigation_stage >= K_MAX` **and** not `osei_full_release`.
5. **DEAD AIR** — early fuel collapse per cutoff, not convoy.
6. **THE RELAY** — late collapse + hub and trust thresholds, not matched above.
7. **Fallback** — If nothing matches: implement as extra content, or map to DEAD AIR or a muted RELAY per product decision (must be deterministic).

## Locked design choice: Broadcast vs Kid

**Option 1 (approved):** If both **FULL BROADCAST** and **THE KID WAS RIGHT** could be true, **FULL BROADCAST wins**. Truth is the irreversible act; step 2 removes step 4 whenever `osei_full_release` is true.

## Locked design choice: Maren

**Separate scores (approved):** `maren_hub_support` and `maren_trust` are both required state; do not collapse into a single composite field for this project.

## Implementation note

Implement as a pure function `ResolveEnding(state) -> EndingID` with no side effects, so tests can assert priority and edge cases.

## Self-review

- No contradictory priority: Kid is evaluated after Broadcast, so Broadcast always suppresses Kid when both flags true.
- Convoy is highest priority so extraction cannot be overwritten by broadcast/dark/kid/failure epilogues if betrayal succeeded.
- `terminal_dark_night < 7` for DEAD AIR is a placeholder aligned with original pitch (“night six”); tune in content pass without changing evaluator shape.
- Maren: `maren_trust` and `maren_hub_support` are both required; caught-lie beat must declare a trust delta so it is mechanically anchored.
