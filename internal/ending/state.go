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
	Night             int
	Fuel              int
	TerminalDarkNight int // night index when fuel first hit 0; 0 if not yet ended that way
	HarrowDarkPlan    bool
	KidInvestigation  int
	OseiFullRelease   bool
	ConvoyBetrayal    bool
	MarenHubSupport   int
	MarenTrust        int
}
