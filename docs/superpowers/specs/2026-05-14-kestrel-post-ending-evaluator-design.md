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
| `maren_hub_support` | int or float | Cumulative support for Maren as hub (routing, honest intel, not undermining when it matters); compared to `M_THRESHOLD` for THE RELAY. |

Constants `K_MAX` and `M_THRESHOLD` are content-tuned; this spec only requires they exist.

## Ending definitions (logic)

- **THE CONVOY (5)** — Extraction via betrayal: `convoy_betrayal` is true and resolves successfully for this run.
- **FULL BROADCAST (4)** — `osei_full_release` is true.
- **DARK FREQUENCY (2)** — `harrow_dark_plan` is true for the sustained window defined by content, and run does not end on a higher-priority ending first.
- **THE KID WAS RIGHT (3)** — `kid_investigation_stage >= K_MAX` and **not** `osei_full_release`.
- **DEAD AIR (6)** — Early failure: fuel exhausted “too soon” (e.g. `fuel <= 0` and `terminal_dark_night < 7`; exact cutoff is pacing), and **not** convoy extraction.
- **THE RELAY (1)** — Bittersweet default: late runway (e.g. fuel hits 0 at `night >= 9` or end-of-night-9 per discrete model), `maren_hub_support >= M_THRESHOLD`, and none of the exclusive higher-priority endings apply.

**RELAY vs DEAD AIR:** Same mechanism (fuel hits zero); split by **when** it happens plus **hub support** and absence of other ending flags.

## Resolution order (first match wins)

Irreversible / player-chosen outcomes trump ambient mystery and default fate.

1. **THE CONVOY** — `convoy_betrayal` resolved successfully.
2. **FULL BROADCAST** — `osei_full_release`.
3. **DARK FREQUENCY** — sustained `harrow_dark_plan` (and not already matched above).
4. **THE KID WAS RIGHT** — `kid_investigation_stage >= K_MAX` **and** not `osei_full_release`.
5. **DEAD AIR** — early fuel collapse per cutoff, not convoy.
6. **THE RELAY** — late collapse + hub threshold, not matched above.
7. **Fallback** — If nothing matches: implement as extra content, or map to DEAD AIR or a muted RELAY per product decision (must be deterministic).

## Locked design choice: Broadcast vs Kid

**Option 1 (approved):** If both **FULL BROADCAST** and **THE KID WAS RIGHT** could be true, **FULL BROADCAST wins**. Truth is the irreversible act; step 2 removes step 4 whenever `osei_full_release` is true.

## Implementation note

Implement as a pure function `ResolveEnding(state) -> EndingID` with no side effects, so tests can assert priority and edge cases.

## Self-review

- No contradictory priority: Kid is evaluated after Broadcast, so Broadcast always suppresses Kid when both flags true.
- Convoy is highest priority so extraction cannot be overwritten by broadcast/dark/kid/failure epilogues if betrayal succeeded.
- `terminal_dark_night < 7` for DEAD AIR is a placeholder aligned with original pitch (“night six”); tune in content pass without changing evaluator shape.
