package game

// ActForNight maps nights 1–9 to act / “level” 1–3.
func ActForNight(night int) int {
	switch {
	case night <= 3:
		return 1
	case night <= 6:
		return 2
	default:
		return 3
	}
}

// ActTitle is a short banner for the act (terminal “level”).
func ActTitle(act int) string {
	switch act {
	case 1:
		return "LEVEL I · STOCK AND SURGE"
	case 2:
		return "LEVEL II · STATIC AND SHADOW"
	case 3:
		return "LEVEL III · CONVOY WINDOW"
	default:
		return "LEVEL"
	}
}

// Choice is one keyed response for a night (1–3).
type Choice struct {
	Label        string
	Fuel         int // subtracted from fuel (positive number = cost)
	Hub          int
	Trust        int
	Kid          int
	SetHarrow    bool
	OseiRelease  bool
	ConvoyBetray bool
}

// NightCard is the incoming beat for a single night.
type NightCard struct {
	Act     int
	Source  string
	Hash    string
	Quote   string
	Choices [3]Choice
}

// NightScript returns the transmission card for the given night (1–9).
// Nights 1–3: Maren / triage. 4–6: Kid + Harrow pressure. 7–9: Cole, Osei, endgame.
func NightScript(night int) NightCard {
	switch night {
	case 1:
		return NightCard{1, "MAREN", "9f2c",
			"We’re holding eight at the community hall. Two fevers, one kid who won’t drink. I need routing—where is the nearest clinic that still has power?",
			[3]Choice{
				{Label: "long medical route + frequencies (−18 fuel · +hub · +trust)", Fuel: 18, Hub: 2, Trust: 1},
				{Label: "short grid packet + IV timing (−10 fuel · +hub · +trust)", Fuel: 10, Hub: 1, Trust: 1},
				{Label: "standby ping only (−3 fuel · −trust)", Fuel: 3, Hub: 0, Trust: -1},
			}}
	case 2:
		return NightCard{1, "MAREN", "9f2c",
			"Generators are coughing. If I run the hall lights past midnight I lose the vaccine fridge. Tell me straight: do I cut heat or light?",
			[3]Choice{
				{Label: "run numbers + recommend heat curfew (−20 fuel · +hub · +trust)", Fuel: 20, Hub: 2, Trust: 1},
				{Label: "one-line priority: fridge first (−10 fuel · +hub · +trust)", Fuel: 10, Hub: 1, Trust: 1},
				{Label: "defer — say you need morning wx (−4 fuel · −trust)", Fuel: 4, Hub: 0, Trust: -1},
			}}
	case 3:
		return NightCard{1, "MAREN", "9f2c",
			"Someone’s scanning our band with a bad hash. Not hostile yet. Do I answer, jam, or go dark until you pattern-match?",
			[3]Choice{
				{Label: "pattern brief + safe CQ window (−22 fuel · +hub · +trust)", Fuel: 22, Hub: 2, Trust: 1},
				{Label: "tight procedure: listen-only 20m (−10 fuel · +hub)", Fuel: 10, Hub: 1, Trust: 0},
				{Label: "tell them to go dark without you (−3 fuel · −trust)", Fuel: 3, Hub: 0, Trust: -1},
			}}
	case 4:
		return NightCard{2, "KID", "—",
			"No name. Quiet voice. “Kestrel, do the geese still cross the river at dawn? I’m counting.” It sounds stupid until you hear the static under it.",
			[3]Choice{
				{Label: "answer honestly + ask why it matters (−16 fuel · +2 kid steps)", Fuel: 16, Hub: 0, Trust: 0, Kid: 2},
				{Label: "short factual: migration window yes/no (−10 fuel · +1 kid step)", Fuel: 10, Hub: 0, Trust: 0, Kid: 1},
				{Label: "clip carrier — refuse the game (−4 fuel · −1 kid step · −trust)", Fuel: 4, Hub: 0, Trust: -1, Kid: -1},
			}}
	case 5:
		return NightCard{2, "MAREN", "9f2c",
			"The hall line is up again. We lost one today. I don’t need poetry—I need whether the north road convoy signature you saw yesterday matches civ or mil.",
			[3]Choice{
				{Label: "full hash compare + caution tape (−22 fuel · +hub · +trust)", Fuel: 22, Hub: 2, Trust: 1},
				{Label: "mil/civ one-liner from your log (−10 fuel · +hub · +trust)", Fuel: 10, Hub: 1, Trust: 1},
				{Label: "withhold — “can’t confirm” (−5 fuel · −trust · +hub)", Fuel: 5, Hub: 1, Trust: -2},
			}}
	case 6:
		return NightCard{2, "HARROW", "b81d",
			"Post Harrow on narrowband. “You’re lighting up the commons like a Christmas tree, Kestrel. Go tight with me—two nets, shared roster—or we both burn.”",
			[3]Choice{
				{Label: "accept shared roster + split bands (−18 fuel · set Harrow plan)", Fuel: 18, Hub: 1, Trust: 0, SetHarrow: true},
				{Label: "counter: shared listen, you keep broadcast (−10 fuel)", Fuel: 10, Hub: 1, Trust: 0},
				{Label: "decline cold (−4 fuel · −trust)", Fuel: 4, Hub: 0, Trust: -1},
			}}
	case 7:
		return NightCard{3, "COLE", "44aa",
			"Ex-mil cadence. Twelve bodies, one truck. “I need convoy corridor truth. You give me a corridor, I give you silence on your little hall project.”",
			[3]Choice{
				{Label: "trade corridor sketch (stale) (−20 fuel · +hub · −trust)", Fuel: 20, Hub: 2, Trust: -1},
				{Label: "misdirect + burn Cole’s time (−10 fuel · +trust)", Fuel: 10, Hub: 0, Trust: 1},
				{Label: "hard refuse (−6 fuel · −hub)", Fuel: 6, Hub: -1, Trust: 0},
			}}
	case 8:
		return NightCard{3, "DR. OSEI", "loop",
			"Recorded loop, drifting in fragments: “…solar class … undisclosed … cascade …” You could patch it clean and push it wide—or keep it in-house.",
			[3]Choice{
				{Label: "clean + FULL BROADCAST on open net (−24 fuel · Osei release)", Fuel: 24, Hub: 0, Trust: 0, OseiRelease: true},
				{Label: "archive locally + send Maren a redacted slice (−10 fuel · +hub · +trust)", Fuel: 10, Hub: 2, Trust: 1},
				{Label: "log only — no share (−5 fuel)", Fuel: 5, Hub: 0, Trust: 0},
			}}
	case 9:
		return NightCard{3, "MAREN + SKY", "9f2c",
			"Final night pressure: hall is full, Cole’s truck hash is on the road, Harrow’s carrier winks. “Pick who gets your last hour of fuel—me, them, or yourself.”",
			[3]Choice{
				{Label: "all-in for the hall + relay map (−26 fuel · +hub · +trust)", Fuel: 26, Hub: 3, Trust: 2},
				{Label: "split: fifteen minutes each net (−10 fuel · +hub · +trust)", Fuel: 10, Hub: 2, Trust: 1},
				{Label: "sell Maren grid to convoy for extraction (−8 fuel · CONVOY)", Fuel: 8, Hub: -2, Trust: -3, ConvoyBetray: true},
			}}
	default:
		// Fallback: treat as night 1 beat if out of range.
		return NightScript(1)
	}
}
