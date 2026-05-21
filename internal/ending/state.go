package ending

// Config holds tunable thresholds; content can override per build or load later.
type Config struct {
	KMax int // kid_investigation_stage >= KMax counts as Kid payoff

	MThreshold int // maren_hub_support >= for THE_RELAY
	TThreshold int // maren_trust >= for THE_RELAY
	HThreshold int // sustained nights on Harrow plan for DARK_FREQUENCY

	// RelayMinTerminalNight: fuel ran out on this night or later.
	RelayMinTerminalNight int

	// DeadAirExclusiveMaxTerminalNight: GONE_DARK when terminal_dark_night < this value.
	DeadAirExclusiveMaxTerminalNight int
}

func DefaultConfig() Config {
	return Config{
		KMax:                             6,
		MThreshold:                       14,
		TThreshold:                       8,
		HThreshold:                       5,
		RelayMinTerminalNight:            16,
		DeadAirExclusiveMaxTerminalNight: 16,
	}
}

// RunState is the input to ResolveEnding; narrative systems mutate this over a run.
type RunState struct {
	Night             int
	Fuel              int
	TerminalDarkNight int // night index when fuel first hit 0; 0 if not yet ended that way
	HarrowDarkPlan    bool
	HarrowDarkNights  int
	KidInvestigation  int
	OseiFullRelease   bool
	ConvoyBetrayal    bool
	MarenHubSupport   int
	MarenTrust        int
	DogAlive          bool
	DogName           string
	DogHunger         int
	ThreadsPinned     []string
	FrequenciesSeen   map[string]bool
	Weather           string
	NightFaults       []string
}
